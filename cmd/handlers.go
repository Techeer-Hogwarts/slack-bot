package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type FormMessage struct {
	FormID        string   `json:"form_id"`
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
	log.Println(payloadStr)
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	handleBlockActions(payload)
}
func handleBlockActions(payload slack.InteractionCallback) {
	log.Println(payload)
	returnMessage := FormMessage{FormID: "test"}
	log.Println("Received block actions")
	for blockID, actionValues := range payload.View.State.Values {
		for actionID, blockAction := range actionValues {
			switch actionID {
			case "multi_static_select-action":
				log.Println("Received multi_static_select block action")
				static_actions := []string{}
				for _, action := range blockAction.SelectedOptions {
					static_actions = append(static_actions, action.Value)
					log.Printf("Received input block action: %s", action.Value)
				}
				if blockID == "team_role_block" {
					returnMessage.TeamRoles = static_actions
				}
				if blockID == "tech_stack_block" {
					returnMessage.TechStacks = static_actions
				}
			case "team_intro":
				log.Println("Received team intro block action")
				log.Printf("Received input block action: %s", blockAction.Value)
				returnMessage.TeamIntro = blockAction.Value
			case "team_name":
				log.Println("Received team name block action")
				log.Printf("Received input block action: %s", blockAction.Value)
				returnMessage.TeamName = blockAction.Value
			case "multi_users_select-action":
				log.Println("Received multi_users_select block action")
				for _, action := range blockAction.SelectedUsers {
					returnMessage.Members = append(returnMessage.Members, action)
					log.Printf("Received input block action: %s", action)
				}
			case "num_members":
				log.Println("Received num members block action")
				log.Printf("Received input block action: %s", blockAction.Value)
				returnMessage.NumNewMembers = blockAction.Value
			case "plain_text_input-action":
				log.Println("Received plain text input block action")
				log.Printf("Received input block action: %s", blockAction.Value)
				if blockID == "team_desc_block" {
					returnMessage.Description = blockAction.Value
				}
				if blockID == "team_etc_block" {
					returnMessage.Etc = blockAction.Value
				}
			}
		}
	}
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

	if command == "/구인" {
		OpenRecruitmentModal(w, triggerID)
		return
	}

	// Handle other commands or respond to invalid commands
	w.WriteHeader(http.StatusOK)
}
