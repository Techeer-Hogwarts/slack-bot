package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

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

// func getChannelMessages(api *slack.Client, channelID string) ([]string, error) {
// 	var messageTexts []string
// 	historyParams := slack.GetConversationHistoryParameters{
// 		ChannelID: channelID,
// 		Limit:     100,
// 	}

// 	history, err := api.GetConversationHistory(&historyParams)
// 	if err != nil {
// 		return nil, err
// 	}

//		for _, message := range history.Messages {
//			messageTexts = append(messageTexts, message.Text)
//		}
//		return messageTexts, nil
//	}
func getChannelMessages(api *slack.Client, channelID string) (*slack.GetConversationHistoryResponse, error) {
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
	}

	history, err := api.GetConversationHistory(&historyParams)
	if err != nil {
		return nil, err
	}

	return history, nil
}

// func TriggerEvent(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Received a trigger event request")

// 	api := slack.New(botToken)
// 	messages, err := getChannelMessages(api, channelID)
// 	log.Printf("channelID: %v", channelID)
// 	if err != nil {
// 		log.Printf("Failed to retrieve messages: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(messages); err != nil {
// 		log.Printf("Failed to encode messages: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Println("Trigger event processed successfully")
// }

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a trigger event request")

	api := slack.New(botToken)
	history, err := getChannelMessages(api, channelID)
	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("channelID: %v", channelID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Trigger event processed successfully")
}

func getUsernameAndEmail(api *slack.Client, userID string) (string, string, error) {
	user, err := api.GetUserInfo(userID)
	log.Println("user: ", user)
	if err != nil {
		return "", "", err
	}
	return user.Name, user.Profile.Email, nil
}

func constructMessageText(message FormMessage) string {
	return "New recruitment form submitted:\n" +
		"Team Introduction: " + message.TeamIntro + "\n" +
		"Team Name: " + message.TeamName + "\n" +
		"Team Leader: " + message.TeamLeader + "\n" +
		"Roles Needed: " + formatList(message.TeamRoles) + "\n" +
		"Tech Stacks: " + formatList(message.TechStacks) + "\n" +
		"Members: " + formatList(message.Members) + "\n" +
		"Number of New Members: " + message.NumNewMembers + "\n" +
		"Description: " + message.Description + "\n" +
		"Other Details: " + message.Etc
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "None"
	}
	return "- " + strings.Join(items, "\n- ")
}
