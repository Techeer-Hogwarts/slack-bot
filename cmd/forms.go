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
	modal, err := readModalJSON("recruit.json")
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

func sendDMToLeader(api *slack.Client, msg ApplyMessage) error {
	buttonValue := fmt.Sprintf("%s|%s", msg.Applicant, msg.TeamID)
	enrollButton := slack.NewButtonBlockElement(
		"enroll_button", // ActionID
		buttonValue,     // Value
		slack.NewTextBlockObject("plain_text", "수락", false, false),
	)
	actionBlock := slack.NewActionBlock("", enrollButton)

	messageText := fmt.Sprintf(":heavy_exclamation_mark: 새로운 지원자가 있습니다!\n\n*지원자:* <@%s>\n\n*나이:* %s\n\n*학년:* %s\n\n*자기소개:* %s\n\n*희망 직군:* %s\n\n", msg.Applicant, msg.Age, msg.Grade, msg.Pr, msg.Role)
	sectionBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", messageText, false, false),
		nil, nil,
	)

	messageOptions := slack.MsgOptionBlocks(sectionBlock, actionBlock)
	_, _, err := api.PostMessage(msg.Leader, messageOptions)
	return err
}

func sendDMSuccessMessage(api *slack.Client, applicant, message string) error {
	successMessage := slack.MsgOptionText(message, false)
	_, _, err := api.PostMessage(applicant, successMessage)
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
	descInput := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "지원 동기/자기소개", false, false),
		"desc_action",
	)
	descInput.Multiline = true
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
					nil,
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
					nil, // No hint text for this input block
					descInput,
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
