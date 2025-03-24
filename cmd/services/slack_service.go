package services

import "github.com/slack-go/slack"

type SlackService interface {
	// Define methods for Slack service
}

type slackService struct {
	client *slack.Client
}

func NewSlackService(client *slack.Client) SlackService {
	return &slackService{client: client}
}
