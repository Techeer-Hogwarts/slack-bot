package cmd

import (
	"log"
	"net/http"
)

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request to the root path")
	w.Write([]byte("Hello, World!"))
}

// func HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
// 	if err := VerifySlackRequest(r); err != nil {
// 		log.Printf("Invalid request: %v", err)
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	var request slack.SlashCommand
// 	log.Println(r.Body)
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		log.Printf("Failed to decode request body: %v", err)
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	log.Println("Received a request to the slash command path")
// 	if request.Command == "/구인" {
// 		OpenRecruitmentModal(w, request.TriggerID)
// 		return
// 	}

//		w.WriteHeader(http.StatusOK)
//	}
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
