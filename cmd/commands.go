package cmd

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

func SendHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
	if err := VerifySlackRequest(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var s slack.SlashCommand
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if s.Command == "/구인" {
		response := slack.Msg{
			ResponseType: slack.ResponseTypeEphemeral,
			Text:         "This is a private response to the /구인 command!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid command"))
	}
}
