package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a slash command request")

	if err := VerifySlackRequest(r); err != nil {
		log.Printf("Verification failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("Error parsing slash command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cmd.Command == "/구인" {
		response := slack.Msg{
			ResponseType: slack.ResponseTypeEphemeral,
			Text:         "This is a private response to the /구인 command!",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println("Slash command processed successfully")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid command"))
		log.Println("Received an invalid command")
	}
}
