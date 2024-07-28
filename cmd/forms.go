package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func loadModalJSON(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func OpenRecruitmentModal(w http.ResponseWriter, triggerID string) {
	modalJSON, err := loadModalJSON("recruitment_form.json")
	if err != nil {
		log.Printf("Failed to load modal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	api := slack.New(botToken)
	viewRequest := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Blocks: slack.Blocks{},
	}
	err = json.Unmarshal([]byte(modalJSON), &viewRequest)
	if err != nil {
		log.Printf("Failed to unmarshal modal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = api.OpenView(triggerID, viewRequest)
	if err != nil {
		log.Printf("Failed to open modal: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
