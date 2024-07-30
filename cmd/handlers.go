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
	TeamRoles     []string `json:"roles"`
	TechStacks    []string `json:"tech"`
	Members       []string `json:"members"`
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
		log.Println(payload.User)
		log.Println(payload.Type)
		log.Println(payload.ActionID)
		log.Println(payload.Message)
		log.Println(payload.Name)
		log.Println("Received block actions 지원하기")
		for _, action := range payload.ActionCallback.BlockActions {
			if action.ActionID == "apply_button" {
				err := openApplyModal(payload.TriggerID)
				if err != nil {
					log.Printf("Failed to open modal: %v", err)
				}
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else if payload.Type == slack.InteractionTypeViewSubmission {
		log.Println(payload.Type)
		log.Println(payload.User)
		log.Println(payload.ActionID)
		log.Println(payload.Message)
		log.Println(payload.Name)
		log.Println(payload.CallbackID)
		jsonVal := handleBlockActions(payload)
		log.Println(jsonVal)
		if err := postMessageToChannel(channelID, jsonVal); err != nil {
			log.Printf("Failed to post message to channel: %v", err)
			http.Error(w, "Failed to post message to channel", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func postMessageToChannel(channelID string, message FormMessage) error {
	api := slack.New(botToken)
	messageText, err := constructMessageText(message)
	if err != nil {
		return err
	}
	applyButton := slack.NewButtonBlockElement("apply_button", "apply", slack.NewTextBlockObject("plain_text", "지원하기!", false, false))
	actionBlock := slack.NewActionBlock("", applyButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)

	_, _, err = api.PostMessage(channelID, messageBlocks)
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
				if blockID == "team_role_block" {
					returnMessage.TeamRoles = static_actions
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
			case "num_members":
				returnMessage.NumNewMembers = blockAction.Value
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
	case "/지원":
		openApplicationModal(w, triggerID, api)
		w.WriteHeader(http.StatusOK)
	case "/수정":
		openEditModal(w, triggerID, api)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
