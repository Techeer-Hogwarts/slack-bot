package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type FormMessage struct {
	FormID        string   `json:"form_id"`
	TeamLeader    string   `json:"leader"`
	TeamIntro     string   `json:"intro"`
	TeamName      string   `json:"name"`
	TechStacks    []string `json:"tech"`
	Members       []string `json:"members"`
	UxMembers     string   `json:"ux_members"`
	FrontMembers  string   `json:"front_members"`
	BackMembers   string   `json:"back_members"`
	DataMembers   string   `json:"data_members"`
	OpsMembers    string   `json:"ops_members"`
	StudyMembers  string   `json:"study_members"`
	EtcMembers    string   `json:"etc_members"`
	NumNewMembers string   `json:"num_members"`
	Description   string   `json:"description"`
	Etc           string   `json:"etc"`
}

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request to the root path")
	w.Write([]byte("Hello, World!"))
}

func TestEvent(w http.ResponseWriter, r *http.Request) {
	api := slack.New(botToken)
	err := getAllUsers(api)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
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
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.Type == slack.InteractionTypeBlockActions {
		log.Println("Received block actions 지원하기/삭제하기")
		for _, action := range payload.ActionCallback.BlockActions {
			if action.ActionID == "apply_button" {
				err := openApplyModal(payload.TriggerID)
				if err != nil {
					log.Printf("Failed to open modal: %v", err)
				}
				return
			} else if action.ActionID == "delete_button" {
				err := openDeleteModal(payload.TriggerID)
				if err != nil {
					log.Printf("Failed to open modal: %v", err)
				}
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else if payload.Type == slack.InteractionTypeViewSubmission {
		log.Printf("Trigger_id: %s", payload.TriggerID)
		log.Println(payload.User)
		log.Printf("Token: %s", payload.Token)
		log.Printf("Message Time Stamp: %s", payload.MessageTs)
		log.Println(payload.View.CallbackID) // this is the key to distinguish different modals
		log.Println(payload.View.PrivateMetadata)
		log.Printf("Action ID: %v", payload.ActionID)
		if payload.View.CallbackID == "recruitment_form" {
			jsonVal := handleBlockActions(payload)
			if err := postMessageToChannel(channelID, jsonVal); err != nil {
				log.Printf("Failed to post message to channel: %v", err)
				http.Error(w, "Failed to post message to channel", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else if payload.View.CallbackID == "apply_form" {
			log.Printf("Trigger_id: %s", payload.TriggerID)
			log.Printf("TS: %s", payload.ActionTs)
			log.Printf("Message: %s", payload.Message.ClientMsgID)
			log.Printf("Message1: %s", payload.Message.Text)
			log.Printf("Message2: %s", payload.Message.Timestamp)
			log.Printf("Original Message: %s", payload.OriginalMessage.ClientMsgID)
			log.Printf("Original Message1: %s", payload.OriginalMessage.Text)
			log.Printf("Original Message2: %s", payload.OriginalMessage.Timestamp)
			log.Println(payload.BlockID)
			log.Println(payload.Value)
			log.Println("Response URL: ", payload.ResponseURL)
			log.Println(payload)
			jsonBytes, _ := json.Marshal(payload)
			log.Println(string(jsonBytes))
			log.Println("Received view submission 지원하기")
			w.WriteHeader(http.StatusOK)
		} else if payload.View.CallbackID == "delete_form" {
			log.Printf("Trigger_id: %s", payload.TriggerID)
			log.Printf("Message: %s", payload.Message.ClientMsgID)
			log.Printf("Message1: %s", payload.Message.Text)
			log.Printf("Message2: %s", payload.Message.Timestamp)
			log.Printf("Original Message: %s", payload.OriginalMessage.ClientMsgID)
			log.Printf("Original Message1: %s", payload.OriginalMessage.Text)
			log.Printf("Original Message2: %s", payload.OriginalMessage.Timestamp)
			log.Println("Received view submission 삭제하기")
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
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func handleBlockActions(payload slack.InteractionCallback) FormMessage {
	returnMessage := FormMessage{FormID: "test"}
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
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
