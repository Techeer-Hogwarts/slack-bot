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

func sendDMToLeader(api *slack.Client, msg ApplyMessage) error {
	buttonValue := fmt.Sprintf("%s|%s|%s", msg.Applicant, msg.TeamID, msg.Role)
	enrollButton := slack.NewButtonBlockElement(
		"enroll_button", // ActionID
		buttonValue,     // Value
		slack.NewTextBlockObject("plain_text", "수락", false, false),
	)
	denyButton := slack.NewButtonBlockElement(
		"deny_button", // ActionID
		buttonValue,   // Value
		slack.NewTextBlockObject("plain_text", "거절", false, false),
	)
	actionBlock := slack.NewActionBlock("accept_action", enrollButton, denyButton)

	messageText := fmt.Sprintf(":heavy_exclamation_mark: %s 팀의 새로운 지원자가 있습니다!\n\n*지원자:* <@%s>\n\n*나이:* %s\n\n*대학/직장:* %s\n\n*학년:* %s\n\n*자기소개:* %s\n\n*희망 직군:* %s\n\n", msg.TeamName, msg.Applicant, msg.Age, msg.School, msg.Grade, msg.Pr, roleMap[msg.Role])
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

func openApplyModal(api *slack.Client, payload slack.InteractionCallback) error {
	triggerID := payload.TriggerID
	originalMessageTimestmap := payload.Message.Timestamp
	// payloadJson, _ := json.MarshalIndent(payload, "", "  ")
	// log.Printf("Payload: %s", payloadJson)

	// Fetch active teams
	activeTeams, err := db.GetAllRecruitingTeams()
	if err != nil {
		return fmt.Errorf("failed to get active teams: %w", err)
	}

	var defaultTeam *slack.OptionBlockObject

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
		if team.TeamTs == originalMessageTimestmap {
			defaultTeam = option
		}
	}

	// var roleOptions []*slack.OptionBlockObject

	// Create the modal view request
	descInput := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "지원 동기/자기소개", false, false),
		"desc_action",
	)
	descInput.Multiline = true

	teamSelectElement := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject("plain_text", "팀을 골라주세요", false, false),
		"selected_team",
		options...,
	)
	teamSelectElement.InitialOption = defaultTeam

	extraMessage, err := db.GetExtraMessage(originalMessageTimestmap)
	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}
	var roleOptions []*slack.OptionBlockObject
	if extraMessage.BackendWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("backend", slack.NewTextBlockObject("plain_text", "백엔드", false, false), nil))
	}
	if extraMessage.FrontendWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("frontend", slack.NewTextBlockObject("plain_text", "프런트", false, false), nil))
	}
	if extraMessage.UXWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("uxui", slack.NewTextBlockObject("plain_text", "UX/UI 디자이너", false, false), nil))
	}
	if extraMessage.DevopsWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("devops", slack.NewTextBlockObject("plain_text", "데브옵스/SRE", false, false), nil))
	}
	if extraMessage.DataWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("data", slack.NewTextBlockObject("plain_text", "데이터 엔지니어", false, false), nil))
	}
	if extraMessage.StudyWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("study", slack.NewTextBlockObject("plain_text", "스터디", false, false), nil))
	}
	if extraMessage.EtcWant > 0 {
		roleOptions = append(roleOptions, slack.NewOptionBlockObject("etc", slack.NewTextBlockObject("plain_text", "기타", false, false), nil))
	}

	roleSelectElement := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject("plain_text", "희망 직군을 골라주세요", false, false),
		"selected_role",
		roleOptions...,
	)

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
					teamSelectElement,
				),
				slack.NewInputBlock(
					"age_input",
					slack.NewTextBlockObject("plain_text", "나이를 입력 해주세요", false, false),
					slack.NewTextBlockObject("plain_text", "나이 입력", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "나이", false, false), "age_action"),
				),
				slack.NewInputBlock(
					"school_input",
					slack.NewTextBlockObject("plain_text", "학교 이름을 입력 해주세요 (직장인은 회사 이름을 적어주세요)", false, false),
					slack.NewTextBlockObject("plain_text", "학교 입력", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "학교", false, false), "school_action"),
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
					slack.NewTextBlockObject("plain_text", "지원하는 직군", false, false),
					slack.NewTextBlockObject("plain_text", "role", false, false),
					roleSelectElement,
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
