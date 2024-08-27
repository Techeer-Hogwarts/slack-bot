package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	School    string `json:"school"`
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
					sendFailMessage(api, payload.Channel.ID, payload.User.ID, "지원할수 있는 직군이 없습니다")
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "delete_button" {
				log.Println("Delete button clicked")
				err := deleteMessage(api, payload)
				if err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "enroll_button" {
				log.Println("Enroll button clicked")
				log.Printf("Payload Message: %s", action.Value)
				teamTimestamp, err := enrollUser(api, action.Value, payload.Channel.ID)
				if err != nil {
					log.Printf("Failed to enroll user: %v", err)
				}

				err = updateOpenMessageToChannel(api, channelID, teamTimestamp, payload)
				if err != nil {
					log.Printf("Failed to update message: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "deny_button" {
				log.Println("Deny button clicked")
				log.Printf("Payload Message: %s", action.Value)
				// pyalodJsonVal, _ := json.MarshalIndent(payload, "", "  ")
				// log.Printf("Payload: %s", pyalodJsonVal)
				denyChannelID := payload.Channel.ID
				payloadMessageValue := action.Value // "user_id|team_id|role"
				denyMessageTs := payload.Message.Timestamp
				err := deleteApplyMessage(api, denyChannelID, payloadMessageValue, denyMessageTs)
				if err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "close_button" {
				log.Println("Close button clicked")
				err := closeOpenMessageToChannel(api, channelID, payload.Message.Timestamp, payload)
				if err != nil {
					log.Printf("Failed to close recruitment: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "open_button" {
				log.Println("Open button clicked")
				err := reOpenRecruitment(api, channelID, payload.Message.Timestamp, payload)
				if err != nil {
					log.Printf("Failed to reopen recruitment: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				return
			} else if action.ActionID == "fake_apply_button" {
				log.Println("Fake apply button clicked")
				err := sendFailMessage(api, payload.Channel.ID, payload.User.ID, "이미 지원이 마감된 팀입니다.")
				if err != nil {
					log.Printf("Failed to send error message: %v", err)
				}
				w.WriteHeader(http.StatusOK)
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
			return
		} else if payload.View.CallbackID == "apply_form" {
			log.Println("Received view submission 지원하기")
			applicant := payload.User.ID
			applicantName := payload.User.Name
			log.Println("Applicant Name: ", applicantName)
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
				School:    payload.View.State.Values["school_input"]["school_action"].Value,
				Pr:        payload.View.State.Values["desc_input"]["desc_action"].Value,
				Role:      payload.View.State.Values["role_input"]["selected_role"].SelectedOption.Value,
			}

			err = sendDMToLeader(api, appMsg)
			if err != nil {
				log.Printf("Failed to send DM to leader: %v", err)
				http.Error(w, "Failed to send DM to leader", http.StatusInternalServerError)
				return
			}
			msg := fmt.Sprintf("%s 팀에 지원이 완료되었습니다! 팀 리더에게 DM을 보냈습니다.\n\n\n*지원내용*\n\n*나이:* %s\n\n*대학/직장:* %s\n\n*학년:* %s\n\n*자기소개:* %s\n\n*희망 직군:* %s\n", teamObject.TeamName, appMsg.Age, appMsg.School, appMsg.Grade, appMsg.Pr, roleMap[appMsg.Role])
			err = sendDMSuccessMessage(api, applicant, msg)
			if err != nil {
				log.Printf("Failed to send success message: %v", err)
				http.Error(w, "Failed to send success message", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}

func postOpenMessageToChannel(api *slack.Client, channelID string, message FormMessage) error {
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	applyButton := slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false))
	deleteButton := slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	closeButton := slack.NewButtonBlockElement("close_button", "close", slack.NewTextBlockObject("plain_text", ":lock: 모집 닫기", false, false))

	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton, closeButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	_, timestamp, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	err = addTeamToDB(message, timestamp)
	return err
}

func closeOpenMessageToChannel(api *slack.Client, channelID string, timestamp string, payload slack.InteractionCallback) error {
	actionUserID := payload.User.ID
	if timestamp == "" {
		return errors.New("timestamp cannot be empty")
	}
	if api == nil {
		return errors.New("api is nil")
	}
	teamObj, err := db.GetTeam(timestamp)
	if err != nil {
		log.Printf("Failed to get team: %v", err)
		return err
	}
	if teamObj.TeamLeader == actionUserID || actionUserID == "U02AES3BH17" || actionUserID == "U033UTX061X" {
		err := db.DeactivateRecruitTeam(timestamp)
		if err != nil {
			log.Printf("Failed to deactivate team: %v", err)
			return err
		}
		additionalMesssage, err := db.GetExtraMessage(timestamp)
		if err != nil {
			log.Printf("Failed to get extra message: %v", err)
			return err
		}
		teamIDint, _ := strconv.Atoi(teamObj.TeamID)
		techStacks, err := db.GetTagsFromTeam(teamIDint)
		if err != nil {
			log.Printf("Failed to get tags from team: %v", err)
			return err
		}

		message := FormMessage{
			TeamType:          teamObj.TeamType,
			TeamLeader:        teamObj.TeamLeader,
			TeamIntro:         teamObj.TeamIntro,
			TeamName:          teamObj.TeamName,
			TechStacks:        techStacks,
			Members:           nil,
			NumCurrentMembers: teamObj.NumMembers,
			UxMembers:         strconv.Itoa(additionalMesssage.UXWant),
			FrontMembers:      strconv.Itoa(additionalMesssage.FrontendWant),
			BackMembers:       strconv.Itoa(additionalMesssage.BackendWant),
			DataMembers:       strconv.Itoa(additionalMesssage.DataWant),
			OpsMembers:        strconv.Itoa(additionalMesssage.DevopsWant),
			StudyMembers:      strconv.Itoa(additionalMesssage.StudyWant),
			EtcMembers:        strconv.Itoa(additionalMesssage.EtcWant),
			Description:       teamObj.TeamDesc,
			Etc:               teamObj.TeamEtc,
		}
		messageText, err := constructMessageText(message)
		if err != nil {
			log.Println("Failed to construct message text")
			return err
		}

		applyButton := slack.NewButtonBlockElement("fake_apply_button", "apply", slack.NewTextBlockObject("plain_text", ":x: 팀 모집 마감!", false, false))
		deleteButton := slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
		closeButton := slack.NewButtonBlockElement("open_button", "close", slack.NewTextBlockObject("plain_text", ":unlock: 모집 다시 열기", false, false))

		actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton, closeButton)
		section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
		messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

		_, _, _, err = api.UpdateMessage(channelID, timestamp, messageBlocks)
		if err != nil {
			log.Printf("Failed to send message to channel %s: %v", channelID, err)
			return err
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	} else {
		err = sendFailMessage(api, channelID, actionUserID, "팀 리더만 모집을 닫을 수 있습니다.")
		if err != nil {
			log.Printf("Failed to send error message: %v", err)
			return err
		}
	}
	return err
}

func updateOpenMessageToChannel(api *slack.Client, channelID string, timestamp string, payload slack.InteractionCallback) error {
	if timestamp == "" {
		return errors.New("timestamp cannot be empty")
	}
	if api == nil {
		return errors.New("api is nil")
	}
	err := db.DeactivateRecruitTeam(timestamp)
	if err != nil {
		log.Printf("Failed to deactivate team: %v", err)
		return err
	}
	teamObj, err := db.GetTeam(timestamp)
	if err != nil {
		log.Printf("Failed to get team: %v", err)
		return err
	}
	additionalMesssage, err := db.GetExtraMessage(timestamp)
	if err != nil {
		log.Printf("Failed to get extra message: %v", err)
		return err
	}
	teamIDint, _ := strconv.Atoi(teamObj.TeamID)
	techStacks, err := db.GetTagsFromTeam(teamIDint)
	if err != nil {
		log.Printf("Failed to get tags from team: %v", err)
		return err
	}

	message := FormMessage{
		TeamType:          teamObj.TeamType,
		TeamLeader:        teamObj.TeamLeader,
		TeamIntro:         teamObj.TeamIntro,
		TeamName:          teamObj.TeamName,
		TechStacks:        techStacks,
		Members:           nil,
		NumCurrentMembers: teamObj.NumMembers,
		UxMembers:         strconv.Itoa(additionalMesssage.UXWant),
		FrontMembers:      strconv.Itoa(additionalMesssage.FrontendWant),
		BackMembers:       strconv.Itoa(additionalMesssage.BackendWant),
		DataMembers:       strconv.Itoa(additionalMesssage.DataWant),
		OpsMembers:        strconv.Itoa(additionalMesssage.DevopsWant),
		StudyMembers:      strconv.Itoa(additionalMesssage.StudyWant),
		EtcMembers:        strconv.Itoa(additionalMesssage.EtcWant),
		Description:       teamObj.TeamDesc,
		Etc:               teamObj.TeamEtc,
	}
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	applyButton := slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false))
	deleteButton := slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	closeButton := slack.NewButtonBlockElement("close_button", "close", slack.NewTextBlockObject("plain_text", ":lock: 모집 닫기", false, false))

	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton, closeButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

	_, _, _, err = api.UpdateMessage(channelID, timestamp, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

// this function should delete the message, send dm to applicant and leader that the application is denied
func deleteApplyMessage(api *slack.Client, channelID string, value string, timestamp string) error {
	if timestamp == "" {
		return errors.New("timestamp cannot be empty")
	}
	if api == nil {
		return errors.New("api is nil")
	}
	values := strings.Split(value, "|")
	applicantID := values[0]

	teamID, err := strconv.Atoi(values[1])
	if err != nil {
		return err
	}
	teamObj, err := db.GetTeamByID(teamID)
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf("<@%s>님의 %s 팀 가입 신청을 거절하셨습니다.", applicantID, teamObj.TeamName)
	err = sendDMSuccessMessage(api, applicantID, msgText)
	if err != nil {
		return err
	}
	msgText = fmt.Sprintf("<@%s>님의 %s 팀 가입 신청이 거절되었습니다.", applicantID, teamObj.TeamName)
	err = sendDMSuccessMessage(api, teamObj.TeamLeader, msgText)
	if err != nil {
		return err
	}
	_, _, err = api.DeleteMessage(channelID, timestamp)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
		return err
	}
	return err
}

func reOpenRecruitment(api *slack.Client, channelID string, timestamp string, payload slack.InteractionCallback) error {
	actionUserID := payload.User.ID
	if timestamp == "" {
		return errors.New("timestamp cannot be empty")
	}
	if api == nil {
		return errors.New("api is nil")
	}
	teamObj, err := db.GetTeam(timestamp)
	if err != nil {
		log.Printf("Failed to get team: %v", err)
		return err
	}
	if teamObj.TeamLeader == actionUserID || actionUserID == "U02AES3BH17" || actionUserID == "U033UTX061X" {
		err := db.ActivateRecruitTeam(timestamp)
		if err != nil {
			log.Printf("Failed to activate team: %v", err)
			return err
		}
		additionalMesssage, err := db.GetExtraMessage(timestamp)
		if err != nil {
			log.Printf("Failed to get extra message: %v", err)
			return err
		}
		teamIDint, _ := strconv.Atoi(teamObj.TeamID)
		techStacks, err := db.GetTagsFromTeam(teamIDint)
		if err != nil {
			log.Printf("Failed to get tags from team: %v", err)
			return err
		}
		message := FormMessage{
			TeamType:          teamObj.TeamType,
			TeamLeader:        teamObj.TeamLeader,
			TeamIntro:         teamObj.TeamIntro,
			TeamName:          teamObj.TeamName,
			TechStacks:        techStacks,
			Members:           nil,
			NumCurrentMembers: teamObj.NumMembers,
			UxMembers:         strconv.Itoa(additionalMesssage.UXWant),
			FrontMembers:      strconv.Itoa(additionalMesssage.FrontendWant),
			BackMembers:       strconv.Itoa(additionalMesssage.BackendWant),
			DataMembers:       strconv.Itoa(additionalMesssage.DataWant),
			OpsMembers:        strconv.Itoa(additionalMesssage.DevopsWant),
			StudyMembers:      strconv.Itoa(additionalMesssage.StudyWant),
			EtcMembers:        strconv.Itoa(additionalMesssage.EtcWant),
			Description:       teamObj.TeamDesc,
			Etc:               teamObj.TeamEtc,
		}
		messageText, err := constructMessageText(message)
		if err != nil {
			return err
		}
		applyButton := slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false))
		deleteButton := slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
		closeButton := slack.NewButtonBlockElement("close_button", "close", slack.NewTextBlockObject("plain_text", ":lock: 모집 닫기", false, false))

		actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton, closeButton)
		section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
		messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

		_, _, _, err = api.UpdateMessage(channelID, timestamp, messageBlocks)
		if err != nil {
			log.Printf("Failed to send message to channel %s: %v", channelID, err)
			return err
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	} else {
		err = sendFailMessage(api, channelID, actionUserID, "팀 리더만 모집을 다시 열 수 있습니다.")
		if err != nil {
			log.Printf("Failed to send error message: %v", err)
			return err
		}
	}
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
