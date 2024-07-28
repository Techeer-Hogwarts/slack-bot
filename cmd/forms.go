package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func OpenRecruitmentModal(w http.ResponseWriter, triggerID string) {
	api := slack.New(botToken)

	// Read the modal JSON from a file
	modal, err := readModalJSON("recruitment_form.json")
	if err != nil {
		log.Printf("Failed to read modal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = api.OpenView(triggerID, modal)
	if err != nil {
		log.Printf("Failed to open modal: %v", err)
		http.Error(w, "Failed to open modal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Read modal JSON from a file
func readModalJSON(filename string) (slack.ModalViewRequest, error) {
	var modal slack.ModalViewRequest
	data, err := os.ReadFile(filename)
	if err != nil {
		return modal, err
	}

	if err := json.Unmarshal(data, &modal); err != nil {
		return modal, err
	}

	return modal, nil
}
