package services

import "github.com/slack-go/slack"

type AlertService interface {
	SendAlert(message string) error
}

type alertService struct {
	client *slack.Client
}

func NewAlertService(client *slack.Client) *alertService {
	return &alertService{client: client}
}

func (s *alertService) SendAlert(message string) error {
	_, _, err := s.client.PostMessage(message, slack.MsgOptionText(message, false))
	return err
}
