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
	if err := VerifySlackRequest(r); err != nil {
		log.Printf("Failed to verify request: %v", err)
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	log.Println("Received an interaction payload")
	log.Println(r.Body)
	var payload slack.InteractionCallback
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode interaction payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Type == slack.InteractionTypeBlockActions {
		handleBlockActions(payload)
	} else {
		log.Printf("Unhandled interaction type: %s", payload.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func handleBlockActions(payload slack.InteractionCallback) {
	for _, action := range payload.ActionCallback.BlockActions {
		log.Printf("Received action ID: %s, Received action Value: %s", action.ActionID, action.Value)
		// 	switch action.ActionID {
		// 	case "name_input":
		// 		log.Printf("Name: %s", action.Value)
		// 	case "email_input":
		// 		log.Printf("Email: %s", action.Value)
		// 	case "team_select":
		// 		log.Printf("Selected Team: %s", action.SelectedOption.Value)
		// 	default:
		// 		log.Printf("Unhandled action ID: %s", action.ActionID)
		// 	}
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
