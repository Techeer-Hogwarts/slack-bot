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
	jsonVal := handleBlockActions(payload)
	log.Println(jsonVal)
	if err := postMessageToChannel(channelID, jsonVal); err != nil {
		log.Printf("Failed to post message to channel: %v", err)
		http.Error(w, "Failed to post message to channel", http.StatusInternalServerError)
		return
	}
}

func postMessageToChannel(channelID string, message FormMessage) error {
	api := slack.New(botToken)

	// Convert user IDs to usernames and emails
	for i, userID := range message.Members {
		username, email, err := getUsernameAndEmail(api, userID)
		if err != nil {
			log.Printf("Failed to get user info for userID %s: %v", userID, err)
			continue
		}
		message.Members[i] = username + " (" + email + ")"
	}
	message.TeamLeader, _, _ = getUsernameAndEmail(api, message.TeamLeader)

	// Construct the message text
	messageText := constructMessageText(message)

	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(messageText, false))
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
		openRecruitmentModal(w, triggerID)
	case "/지원":
		// Fetch recruitment messages from "bot-testing" channel
		log.Println("checkpoint 0")
		recruitmentMessages, err := getChannelMessages(api, channelID)
		if err != nil {
			log.Printf("Failed to retrieve recruitment messages: %v", err)
			http.Error(w, "Failed to retrieve recruitment messages", http.StatusInternalServerError)
			return
		}
		log.Println("checkpoint 1")

		// Extract team names from recruitment messages (assuming each message has a unique team name)
		var teams []string
		for _, msg := range recruitmentMessages.Messages {
			teams = append(teams, msg.Text) // Modify as per your message structure to extract team names
		}
		log.Println("checkpoint 2")

		// Create a selection form or modal for the user to choose a team
		modalRequest := createModal(teams)
		log.Println("checkpoint 3")
		// Open the modal
		_, err = api.OpenView(triggerID, modalRequest)
		if err != nil {
			http.Error(w, "Failed to open modal", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
