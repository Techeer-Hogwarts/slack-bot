package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request to the root path")
	w.Write([]byte("Hello, World!"))
}

func HandleInteraction(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Extract the payload from the form data
	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		log.Println("Payload not found in form data")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	handleBlockActions(payload)
}
func handleBlockActions(payload slack.InteractionCallback) {
	log.Println("Received block actions")
	for _, actionValues := range payload.View.State.Values {
		for actionID, blockAction := range actionValues {
			switch actionID {
			case "multi_static_select-action":
				log.Println("Received multi_static_select block action")
				for _, action := range blockAction.SelectedOptions {
					log.Printf("Received input block action: %s", action.Value)
				}
			case "team_intro":
				log.Println("Received team intro block action")
				log.Printf("Received input block action: %s", blockAction.Value)
			case "team_name":
				log.Println("Received team name block action")
				log.Printf("Received input block action: %s", blockAction.Value)
			case "multi_users_select-action":
				log.Println("Received multi_users_select block action")
				for _, action := range blockAction.SelectedUsers {
					log.Printf("Received input block action: %s", action)
				}
			case "num_members":
				log.Println("Received num members block action")
				log.Printf("Received input block action: %s", blockAction.SelectedOption.Value)
			case "plain_text_input-action":
				log.Println("Received plain text input block action")
				log.Printf("Received input block action: %s", blockAction.Value)
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

	// Parse form-encoded data
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract relevant fields
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
