package services

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

type ProfileService interface {
	GetProfilePicture(email string) (string, error)
}

type profileService struct {
	client *slack.Client
}

func NewProfileService(client *slack.Client) *profileService {
	return &profileService{client: client}
}

func (s *profileService) GetProfilePicture(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email is empty")
	}
	log.Printf("Getting profile picture for %s", email)
	profile, err := s.client.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	return profile.Profile.ImageOriginal, nil
}
