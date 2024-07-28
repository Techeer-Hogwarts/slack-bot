package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

var (
	signingKey string
	botToken   string
	channelID  string
)

func init() {
	LoadEnv()
	signingKey = GetEnv("SLACK_SIGNING_SECRET", "")
	botToken = GetEnv("SLACK_BOT_TOKEN", "")
	channelID = GetEnv("SLACK_CHANNEL_ID", "")
}

func VerifySlackRequest(req *http.Request) error {
	s, err := slack.NewSecretsVerifier(req.Header, signingKey)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body)) // Reassign body after reading it

	if _, err := s.Write(body); err != nil {
		return err
	}

	if err := s.Ensure(); err != nil {
		return err
	}

	return nil
}

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a trigger event request")

	var payload struct {
		Message string `json:"message"`
	}
	log.Print("Payload: ", payload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("body: ", r.Body)
	api := slack.New(botToken)

	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(payload.Message, false))
	if err != nil {
		log.Printf("Failed posting message: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Message sent: ", payload.Message)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent: %s", payload.Message)
	log.Println("Trigger event processed successfully")
}
