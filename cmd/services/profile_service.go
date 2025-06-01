package services

import "github.com/slack-go/slack"

type ProfileService interface {
}

type profileService struct {
	client *slack.Client
}

func NewProfileService(client *slack.Client) *profileService {
	return &profileService{client: client}
}
