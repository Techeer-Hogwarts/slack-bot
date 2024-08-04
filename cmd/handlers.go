package cmd

import (
	"encoding/json"
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
	requestVal, err := json.Marshal(r.Form)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}
	log.Printf("Received request: %s", string(requestVal))

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		log.Println("Payload not found in form data")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.Type == slack.InteractionTypeBlockActions {
		log.Println("Received block actions 지원하기/삭제하기")
		for _, action := range payload.ActionCallback.BlockActions {
			if action.ActionID == "apply_button" {
				log.Println("Apply button clicked")
				err := openApplyModal(payload.TriggerID)
				if err != nil {
					log.Printf("Failed to open modal: %v", err)
				}
				return
			} else if action.ActionID == "delete_button" {
				log.Println("Delete button clicked")
				err := deleteMessage(payload)
				if err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
				return
			} else if action.ActionID == "enroll_button" {
				log.Println("Enroll button clicked")
				log.Printf("Payload Message: %s", action.Value)
				err := enrollUser(action.Value, payload.Channel.ID)
				if err != nil {
					log.Printf("Failed to enroll user: %v", err)
				}
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else if payload.Type == slack.InteractionTypeViewSubmission {
		if payload.View.CallbackID == "recruitment_form" {
			jsonVal := handleBlockActions(payload)
			if err := postMessageToChannel(channelID, jsonVal); err != nil {
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
			leaderCode := teamObject.TeamLeader
			if err != nil {
				log.Printf("Failed to get team leader code: %v", err)
				http.Error(w, "Failed to get team leader code", http.StatusInternalServerError)
				return
			}
			api := slack.New(botToken)
			appMsg := ApplyMessage{
				TeamID:    selectedTeam,
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

func postMessageToChannel(channelID string, message FormMessage) error {
	api := slack.New(botToken)
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	applyButton := slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", "지원하기!", false, false))
	deleteButton := slack.NewButtonBlockElement("delete_button", "delete", slack.NewTextBlockObject("plain_text", "삭제하기!", false, false))
	actionBlock := slack.NewActionBlock("apply_action", applyButton)
	actionBlock2 := slack.NewActionBlock("delete_action", deleteButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock, actionBlock2)

	_, timestamp, err := api.PostMessage(channelID, messageBlocks)
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
