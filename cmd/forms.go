package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/thomas-and-friends/slack-bot/db"
)

func openRecruitmentModal(w http.ResponseWriter, triggerID string, api *slack.Client) {
	// Read the modal JSON from a file
	modal, err := readModalJSON("recruitment_form.json")
	if err != nil {
		log.Printf("Failed to read modal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = api.OpenView(triggerID, modal)
	if err != nil {
		log.Printf("Failed to open modal: %v", err)
		http.Error(w, "Failed to open modal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Read modal JSON from a file
func readModalJSON(filename string) (slack.ModalViewRequest, error) {
	var modal slack.ModalViewRequest
	data, err := os.ReadFile(filename)
	if err != nil {
		return modal, err
	}

	if err := json.Unmarshal(data, &modal); err != nil {
		return modal, err
	}

	return modal, nil
}

// func sendDMToLeader(api *slack.Client, msg ApplyMessage) error {
// 	// messageText := fmt.Sprintf("You have a new application for your team!\n\nApplicant: <@%s>\n", msg.Applicant)

//		msgJson, err := json.Marshal(msg)
//		if err != nil {
//			return fmt.Errorf("failed to marshal message: %w", err)
//		}
//		messageText := fmt.Sprintf("You have a new application for your team!\n\nApplicant: %s\n", string(msgJson))
//		_, _, err = api.PostMessage(msg.Leader, slack.MsgOptionText(messageText, false))
//		return err
//	}
func sendDMToLeader(api *slack.Client, msg ApplyMessage) error {
	buttonValue := fmt.Sprintf("%s|%s", msg.Applicant, msg.TeamID)
	enrollButton := slack.NewButtonBlockElement(
		"enroll_button", // ActionID
		buttonValue,     // Value
		slack.NewTextBlockObject("plain_text", "수락", false, false),
	)
	actionBlock := slack.NewActionBlock("", enrollButton)

	// messageText := fmt.Sprintf("You have a new application for your team!\n\nApplicant: <@%s>\n", msg.Applicant)
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	messageText := string(msgJson)
	sectionBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", messageText, false, false),
		nil, nil,
	)

	messageOptions := slack.MsgOptionBlocks(sectionBlock, actionBlock)
	_, _, err = api.PostMessage(msg.Leader, messageOptions)
	return err
}

func openApplyModal(triggerID string) error {
	api := slack.New(botToken)

	// Fetch active teams
	activeTeams, err := db.GetAllTeams()
	if err != nil {
		return fmt.Errorf("failed to get active teams: %w", err)
	}

	// Create options from active teams
	var options []*slack.OptionBlockObject
	for _, team := range activeTeams {
		leaderRealName, _, err := db.GetUser(team.TeamLeader)
		if err != nil {
			return fmt.Errorf("failed to get team leader: %w", err)
		}
		option := slack.NewOptionBlockObject(
			team.TeamID, // Assuming TeamID is a unique identifier
			slack.NewTextBlockObject("plain_text", "팀 이름: "+team.TeamName+" - 팀 리더: "+leaderRealName, false, false),
			nil,
		)
		options = append(options, option)
	}

	// Create the modal view request
	modalRequest := slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: "apply_form",
		Title:      slack.NewTextBlockObject("plain_text", "팀에 지원하기", false, false),
		Close:      slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"team_select",
					slack.NewTextBlockObject("plain_text", "Select a Team", false, false),
					slack.NewTextBlockObject("plain_text", "Select a team", false, false),
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeStatic,
						slack.NewTextBlockObject("plain_text", "팀을 골라주세요", false, false),
						"selected_team",
						options...,
					),
				),
				slack.NewInputBlock(
					"age_input",
					slack.NewTextBlockObject("plain_text", "나이를 입력 해주세요", false, false),
					slack.NewTextBlockObject("plain_text", "나이 입력", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "나이", false, false), "age_action"),
				),
				slack.NewInputBlock(
					"grade_input",
					slack.NewTextBlockObject("plain_text", "학년을 입력 해주세요 (졸업 하셨으면 졸업이라고 적어주세요)", false, false),
					slack.NewTextBlockObject("plain_text", "학년 입력", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "학년", false, false), "grade_action"),
				),
				slack.NewInputBlock(
					"desc_input",
					slack.NewTextBlockObject("plain_text", "지원동기/자기소개", false, false),
					slack.NewTextBlockObject("plain_text", "pr", false, false),
					&slack.PlainTextInputBlockElement{
						ActionID:    "desc_action",
						Placeholder: slack.NewTextBlockObject("plain_text", "지원 동기/자기소개", false, false),
						Multiline:   true,
					},
				),
				slack.NewInputBlock(
					"role_input",
					slack.NewTextBlockObject("plain_text", "희망하는 직군", false, false),
					slack.NewTextBlockObject("plain_text", "role", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "ex. 프런트", false, false), "role_action"),
				),
			},
		},
	}

	_, err = api.OpenView(triggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("failed to open modal view: %w", err)
	}

	return nil
}
