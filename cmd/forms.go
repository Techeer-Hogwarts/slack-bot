package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func openRecruitmentModal(w http.ResponseWriter, triggerID string, api *slack.Client) {
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

func openApplicationModal(w http.ResponseWriter, triggerID string, api *slack.Client) {
	log.Print("Opening application modal")
}

func openEditModal(w http.ResponseWriter, triggerID string, api *slack.Client) {
	log.Print("Opening edit modal")
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

// func createModal(teams []string) slack.ModalViewRequest {
// 	options := []*slack.OptionBlockObject{}
// 	for _, team := range teams {
// 		option := slack.NewOptionBlockObject(team, slack.NewTextBlockObject("plain_text", team, false, false), nil)
// 		options = append(options, option)
// 	}

// 	selectBlockElement := slack.NewOptionsSelectBlockElement(
// 		slack.OptTypeStatic,
// 		slack.NewTextBlockObject("plain_text", "Select a team", false, false),
// 		"team_select",
// 		options...,
// 	)

// 	selectBlock := slack.NewInputBlock(
// 		"team_select",
// 		slack.NewTextBlockObject("plain_text", "Select a team", false, false),
// 		nil,
// 		selectBlockElement,
// 	)

// 	blocks := slack.Blocks{
// 		BlockSet: []slack.Block{
// 			slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", "Apply to a Team", false, false), nil, nil),
// 			selectBlock,
// 		},
// 	}

// 	return slack.ModalViewRequest{
// 		Type:       slack.VTModal,
// 		Title:      slack.NewTextBlockObject("plain_text", "Team Application", false, false),
// 		Blocks:     blocks,
// 		CallbackID: "recruitment_form",
// 	}
// }
