package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

var (
	signingKey = os.Getenv("SLACK_SIGNING_SECRET")
	botToken   = os.Getenv("SLACK_BOT_TOKEN")
	channelID  = os.Getenv("CHANNEL_ID")
)

func VerifySlackRequest(req *http.Request) error {
	s, err := slack.NewSecretsVerifier(req.Header, signingKey)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if _, err := s.Write(body); err != nil {
		return err
	}

	if err := s.Ensure(); err != nil {
		return err
	}

	return nil
}

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api := slack.New(botToken)
	channelID := channelID

	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(payload.Message, false))
	if err != nil {
		log.Printf("failed posting message: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent: %s", payload.Message)
}
