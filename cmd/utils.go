package cmd

import (
	"bytes"
	"encoding/json"
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

func getChannelMessages(api *slack.Client, channelID string) ([]slack.Message, error) {
	var messages []slack.Message
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: "#bot-testing",
		Limit:     100, // Adjust the limit as needed
	}

	history, err := api.GetConversationHistory(&historyParams)
	if err != nil {
		return nil, err
	}

	messages = append(messages, history.Messages...)
	return messages, nil
}

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a trigger event request")

	api := slack.New(botToken)
	messages, err := getChannelMessages(api, channelID)
	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(messages)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Failed to encode messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Trigger event processed successfully")
}
