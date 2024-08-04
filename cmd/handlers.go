package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/thomas-and-friends/slack-bot/db"
)

type FormMessage struct {
	TeamType          string   `json:"type"`
	TeamLeader        string   `json:"leader"`
	TeamIntro         string   `json:"intro"`
	TeamName          string   `json:"name"`
	TechStacks        []string `json:"tech"`
	Members           []string `json:"members"`
	NumCurrentMembers int      `json:"current"`
	UxMembers         string   `json:"ux_members"`
	FrontMembers      string   `json:"front_members"`
	BackMembers       string   `json:"back_members"`
	DataMembers       string   `json:"data_members"`
	OpsMembers        string   `json:"ops_members"`
	StudyMembers      string   `json:"study_members"`
	EtcMembers        string   `json:"etc_members"`
	Description       string   `json:"description"`
	Etc               string   `json:"etc"`
}

type ApplyMessage struct {
	TeamID    string `json:"team_id"`
	TeamName  string `json:"team_name"`
	Leader    string `json:"leader"`
	Applicant string `json:"applicant"`
	Age       string `json:"age"`
	Grade     string `json:"grade"`
	Pr        string `json:"pr"`
	Role      string `json:"role"`
}

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request to the root path")
	w.Write([]byte("Hello, World!"))
}

func HandleInteraction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		log.Println("Payload not found in form data")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var payload slack.InteractionCallback
	err := payload.UnmarshalJSON([]byte(payloadStr))
	if err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	api := slack.New(botToken)
	if payload.Type == slack.InteractionTypeBlockActions {
		log.Println("Received block actions 지원하기/삭제하기")
		for _, action := range payload.ActionCallback.BlockActions {
			if action.ActionID == "apply_button" {
				log.Println("Apply button clicked")
				err := openApplyModal(api, payload)
				if err != nil {
					log.Printf("Failed to open modal: %v", err)
				}
				return
			} else if action.ActionID == "delete_button" {
				log.Println("Delete button clicked")
				err := deleteMessage(api, payload)
				if err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
				return
			} else if action.ActionID == "enroll_button" {
				log.Println("Enroll button clicked")
				log.Printf("Payload Message: %s", action.Value)
				err := enrollUser(api, action.Value, payload.Channel.ID)
				if err != nil {
					log.Printf("Failed to enroll user: %v", err)
				}
				return
			} else if action.ActionID == "close_button" {
				log.Println("Close button clicked")
				jsonVal := handleBlockActions(payload)
				err := updateOpenMessageToChannel(api, channelID, jsonVal, payload.Message.Timestamp)
				if err != nil {
					log.Printf("Failed to close recruitment: %v", err)
				}
				return
			} else if action.ActionID == "open_button" {
				log.Println("Open button clicked")
				jsonVal := handleBlockActions(payload)
				err := reOpenRecruitment(api, channelID, jsonVal, payload.Message.Timestamp)
				if err != nil {
					log.Printf("Failed to reopen recruitment: %v", err)
				}
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else if payload.Type == slack.InteractionTypeViewSubmission {
		if payload.View.CallbackID == "recruitment_form" {
			jsonVal := handleBlockActions(payload)
			if err := postOpenMessageToChannel(api, channelID, jsonVal); err != nil {
				log.Printf("Failed to post message to channel: %v", err)
				http.Error(w, "Failed to post message to channel", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else if payload.View.CallbackID == "apply_form" {
			log.Println("Received view submission 지원하기")
			applicant := payload.User.ID
			selectedTeam := payload.View.State.Values["team_select"]["selected_team"].SelectedOption.Value
			teamID, err := strconv.Atoi(selectedTeam)
			if err != nil {
				log.Printf("Failed to convert teamID to int: %v", err)
				http.Error(w, "Failed to convert teamID to int", http.StatusInternalServerError)
				return
			}
			teamObject, err := db.GetTeamByID(teamID)
			if err != nil {
				log.Printf("Failed to get team leader code: %v", err)
				http.Error(w, "Failed to get team leader code", http.StatusInternalServerError)
				return
			}
			leaderCode := teamObject.TeamLeader
			teamName := teamObject.TeamName
			appMsg := ApplyMessage{
				TeamID:    selectedTeam,
				TeamName:  teamName,
				Leader:    leaderCode,
				Applicant: applicant,
				Age:       payload.View.State.Values["age_input"]["age_action"].Value,
				Grade:     payload.View.State.Values["grade_input"]["grade_action"].Value,
				Pr:        payload.View.State.Values["desc_input"]["desc_action"].Value,
				Role:      payload.View.State.Values["role_input"]["role_action"].Value,
			}

			err = sendDMToLeader(api, appMsg)
			if err != nil {
				log.Printf("Failed to send DM to leader: %v", err)
				http.Error(w, "Failed to send DM to leader", http.StatusInternalServerError)
				return
			}
			msg := fmt.Sprintf("%s 팀에 지원이 완료되었습니다! 팀 리더에게 DM을 보냈습니다.", teamObject.TeamName)
			err = sendSuccessMessage(api, channelID, applicant, msg)
			if err != nil {
				log.Printf("Failed to send success message: %v", err)
				http.Error(w, "Failed to send success message", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

func postOpenMessageToChannel(api *slack.Client, channelID string, message FormMessage) error {
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	actionBlock := slack.NewActionBlock(
		"action_block_id",
		slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", true, true)),
		slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", true, true)),
		slack.NewButtonBlockElement("close_button", "close", slack.NewTextBlockObject("plain_text", ":lock: 모집 닫기", true, true)),
	)

	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, true, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

	_, timestamp, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	err = addTeamToDB(message, timestamp)
	return err
}

func updateOpenMessageToChannel(api *slack.Client, channelID string, message FormMessage, timestamp string) error {
	// err := db.DeactivateRecruitTeam(timestamp)
	// if err != nil {
	// 	log.Printf("Failed to deactivate team: %v", err)
	// 	return err
	// }
	// messageText, err := constructMessageText(message)
	// if err != nil {
	// 	return err
	// }
	// actionBlock := slack.NewActionBlock(
	// 	"action_block_id",
	// 	slack.NewButtonBlockElement("open_button", "open", slack.NewTextBlockObject("plain_text", ":unlock: 모집 다시 열기", false, false)),
	// 	slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false)),
	// )

	// section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, true, false), nil, nil)
	// messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

	// _, _, _, err = api.UpdateMessage(channelID, timestamp, messageBlocks)
	// if err != nil {
	// 	log.Printf("Failed to send message to channel %s: %v", channelID, err)
	// 	return err
	// }
	// log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	// err = addTeamToDB(message, timestamp)
	// return err
	button1 := slack.NewButtonBlockElement(
		"actionId-0",   // Action ID for the button
		"click_me_123", // Value for the button
		slack.NewTextBlockObject("plain_text", "Click Me", true, true),
	)

	button2 := slack.NewButtonBlockElement(
		"actionId-1",   // Action ID for the button
		"click_me_123", // Value for the button
		slack.NewTextBlockObject("plain_text", "Click Me", true, true),
	)

	// Create an action block with the buttons
	actionBlock := slack.NewActionBlock(
		"actions_block_id", // Block ID for the action block
		button1,
		button2,
	)

	// Create the message with the action block
	messageBlocks := slack.MsgOptionBlocks(actionBlock)

	// Post the message to Slack
	_, timestamp, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func reOpenRecruitment(api *slack.Client, channelID string, message FormMessage, timestamp string) error {
	err := db.ActivateRecruitTeam(timestamp)
	if err != nil {
		log.Printf("Failed to activate team: %v", err)
		return err
	}
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	actionBlock := slack.NewActionBlock(
		"action_block_id",
		slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false)),
		slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false)),
		slack.NewButtonBlockElement("close_button", "close", slack.NewTextBlockObject("plain_text", ":lock: 모집 닫기", false, false)),
	)

	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, true, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

	_, _, _, err = api.UpdateMessage(channelID, timestamp, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	err = addTeamToDB(message, timestamp)
	return err
}

func handleBlockActions(payload slack.InteractionCallback) FormMessage {
	returnMessage := FormMessage{}
	for blockID, actionValues := range payload.View.State.Values {
		for actionID, blockAction := range actionValues {

			switch actionID {
			case "multi_static_select-action":
				static_actions := []string{}
				for _, action := range blockAction.SelectedOptions {
					static_actions = append(static_actions, action.Value)
				}
				if blockID == "tech_stack_block" {
					returnMessage.TechStacks = static_actions
				}
			case "users_select-action":
				returnMessage.TeamLeader = blockAction.SelectedUser
			case "team_intro":
				returnMessage.TeamIntro = blockAction.Value
			case "team_name":
				returnMessage.TeamName = blockAction.Value
			case "multi_users_select-action":
				returnMessage.Members = append(returnMessage.Members, blockAction.SelectedUsers...)
			case "num_ux_members":
				returnMessage.UxMembers = blockAction.Value
			case "num_front_members":
				returnMessage.FrontMembers = blockAction.Value
			case "num_back_members":
				returnMessage.BackMembers = blockAction.Value
			case "num_data_members":
				returnMessage.DataMembers = blockAction.Value
			case "num_sre_members":
				returnMessage.OpsMembers = blockAction.Value
			case "num_study_members":
				returnMessage.StudyMembers = blockAction.Value
			case "num_etc_members":
				returnMessage.EtcMembers = blockAction.Value
			case "plain_text_input-action":
				if blockID == "team_desc_block" {
					returnMessage.Description = blockAction.Value
				}
				if blockID == "team_etc_block" {
					returnMessage.Etc = blockAction.Value
				}
			}
			returnMessage.NumCurrentMembers = len(returnMessage.Members)
			numStudy, err := strconv.Atoi(returnMessage.StudyMembers)
			if err != nil {
				log.Printf("Failed to convert numStudy to int: %v", err)
				numStudy = 0
			}
			numEtc, err := strconv.Atoi(returnMessage.EtcMembers)
			if err != nil {
				log.Printf("Failed to convert numEtc to int: %v", err)
				numEtc = 0
			}
			if numStudy > 0 {
				returnMessage.TeamType = "study"
			} else if numEtc > 0 {
				returnMessage.TeamType = "etc"
			} else {
				returnMessage.TeamType = "project"
			}
		}
	}
	return returnMessage
}

func HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
	if err := VerifySlackRequest(r); err != nil {
		log.Printf("Invalid request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	command := r.FormValue("command")
	triggerID := r.FormValue("trigger_id")

	log.Printf("Received command: %s", command)
	log.Printf("Trigger ID: %s", triggerID)

	api := slack.New(botToken)

	switch command {
	case "/구인":
		openRecruitmentModal(w, triggerID, api)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
