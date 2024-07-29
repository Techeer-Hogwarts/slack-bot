package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type RichTextElement struct {
	Type     string            `json:"type"`
	Elements []RichTextElement `json:"elements,omitempty"`
	Text     string            `json:"text,omitempty"`
}

type RichTextInputValue struct {
	Type     string            `json:"type"`
	Elements []RichTextElement `json:"elements"`
}

type BlockAction struct {
	ActionID        string          `json:"action_id"`
	SelectedOptions []Option        `json:"selected_options,omitempty"`
	SelectedUsers   []string        `json:"selected_users,omitempty"`
	Value           json.RawMessage `json:"value"` // RawMessage to handle both plain and rich text
	Type            string          `json:"type"`
}

type Option struct {
	Value string `json:"value"`
}

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
	for blockID, actionValues := range payload.View.State.Values {
		for actionID, blockAction := range actionValues {
			if blockID == "tech_stack_block" {
				log.Println("Received tech stack block action")
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.SelectedOptions, blockAction.Type)
				for _, option := range blockAction.SelectedOptions {
					log.Printf("Selected option value: %s", option.Value)
				}
			}
			if blockID == "current_members_block" {
				log.Println("Received current members block action")
				log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.SelectedUsers, blockAction.Type)
				for _, user := range blockAction.SelectedUsers {
					log.Printf("Selected user ID: %s", user)
				}
			}

			if blockID == "rich_text_block" && actionID == "rich_text_input-action" {
				log.Println("Received rich text block action")
				var richTextValue RichTextInputValue
				if err := json.Unmarshal([]byte(blockAction.Value), &richTextValue); err != nil {
					log.Printf("Failed to parse rich text input: %v", err)
					continue
				}
				log.Printf("Block ID: %s, Action ID: %s, Rich Text Value: %+v, Type: %s", blockID, actionID, richTextValue, blockAction.Type)
				for _, element := range richTextValue.Elements {
					processRichTextElement(element)
				}
			} else {
				log.Println("Received plain text block action")
				var plainTextValue string
				if err := json.Unmarshal([]byte(blockAction.Value), &plainTextValue); err == nil {
					log.Printf("Block ID: %s, Action ID: %s, Plain Text Value: %s, Type: %s", blockID, actionID, plainTextValue, blockAction.Type)
				} else {
					log.Printf("Block ID: %s, Action ID: %s, Value: %v, Type: %v", blockID, actionID, blockAction.Value, blockAction.Type)
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
