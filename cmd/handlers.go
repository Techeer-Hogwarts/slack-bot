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
	log.Println(payload.View)
	// for _, action := range payload.ActionCallback.BlockActions {
	// 	log.Printf("Received action ID: %s", action.ActionID)
	// 	// Handle different block types
	// 	switch action.Type {
	// 	case "input":
	// 		// Assuming you have a form submission action
	// 		// For input blocks, you might have to check the parent view's state
	// 		// Note: Slack does not include input block data directly in block actions
	// 		// You might need to refer to `payload.View.State.Values` for input data
	// 		// See below for handling values in `payload.View.State.Values`

	// 	case "static_select":
	// 		// Handle static select (dropdown)
	// 		log.Printf("Selected option value: %s", action.SelectedOption.Value)

	// 	case "multi_static_select":
	// 		// Handle multi-static select
	// 		if action.SelectedOptions != nil {
	// 			for _, option := range action.SelectedOptions {
	// 				log.Printf("Selected option value: %s", option.Value)
	// 			}
	// 		}

	// 	default:
	// 		// Log any unhandled block types
	// 		log.Printf("Unhandled block type: %s", action.Type)
	// 	}
	// }

	// Handle values from input blocks in the view's state
	for blockID, actionValues := range payload.View.State.Values {
		for actionID, blockAction := range actionValues {
			if blockID == "tech_stack_block" {
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.SelectedOptions, blockAction.Type)
				for _, option := range blockAction.SelectedOptions {
					log.Printf("Selected option value: %s", option.Value)
				}
			}
			if blockID == "current_members_block" {
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.SelectedUsers, blockAction.Type)
				for _, user := range blockAction.SelectedUsers {
					log.Printf("Selected user ID: %s", user)
				}
			}
			if actionID == "rich_text_input-action" {
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.Value, blockAction.Type)
			} else {
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.Value, blockAction.Type)
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
