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
	channelID = GetEnv("CHANNEL_ID", "")
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
		ChannelID: channelID,
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
	log.Printf("channelID: %v", channelID)
	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Failed to encode messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Trigger event processed successfully")
}

func listChannels(api *slack.Client) ([]slack.Channel, error) {
	var allChannels []slack.Channel
	cursor := ""
	for {
		params := &slack.GetConversationsParameters{
			Limit:  100, // Adjust the limit as needed
			Cursor: cursor,
		}
		channels, nextCursor, err := api.GetConversations(params)
		if err != nil {
			return nil, err
		}
		allChannels = append(allChannels, channels...)
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	return allChannels, nil
}

func ListChannels(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a list channels request")

	api := slack.New(botToken)
	channels, err := listChannels(api)
	if err != nil {
		log.Printf("Failed to retrieve channels: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(channels); err != nil {
		log.Printf("Failed to encode channels: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("List channels processed successfully")
}
