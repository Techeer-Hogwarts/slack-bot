package services

import "github.com/slack-go/slack"

type SlackService interface {
	DeleteMessage(channelID string, message string) error
}

type slackService struct {
	client *slack.Client
}

func NewSlackService(client *slack.Client) SlackService {
	return &slackService{client: client}
}

func (s *slackService) DeleteMessage(channelID string, messageTimestamp string) error {
	_, _, err := s.client.DeleteMessage(channelID, messageTimestamp)
	return err
}
