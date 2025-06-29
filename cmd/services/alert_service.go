package services

import (
	"errors"
	"log"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/slackmessages"
	"github.com/slack-go/slack"
)

type AlertService interface {
	SendAlert(alertMessage models.AlertMessageSchema) error
	SendAlertToFindMember(FindMemberObject models.FindMemberSchema) error
	SendAlertToUser(UserObject models.UserMessageSchema) error
}

type alertService struct {
	client              *slack.Client
	findMemberChannelID string
}

func NewAlertService(client *slack.Client, findMemberChannelID string) *alertService {
	return &alertService{
		client:              client,
		findMemberChannelID: findMemberChannelID,
	}
}

// SendAlert 유저 메시지 전송 (/alert/message 에서 사용 - 새로운거)
func (s *alertService) SendAlert(alertMessage models.AlertMessageSchema) error {
	switch alertMessage.Type {
	case "user":
		user, err := s.client.GetUserByEmail(alertMessage.Email)
		if err != nil {
			return err
		}
		log.Printf("userID: %s, message: %s", user.ID, alertMessage.Message)
		_, _, err = s.client.PostMessage(user.ID, slack.MsgOptionText(alertMessage.Message, false))
		return err
	case "channel":
		log.Printf("channelID: %s, message: %s", alertMessage.ChannelID, alertMessage.Message)
		_, _, err := s.client.PostMessage(alertMessage.ChannelID, slack.MsgOptionText(alertMessage.Message, false))
		return err
	default:
		return errors.New("invalid type")
	}
}

// SendAlertToFindMember 스터디 팀 공고 메시지 전송 (/alert/find-member 와 /alert/channel 에서 사용)
func (s *alertService) SendAlertToFindMember(FindMemberObject models.FindMemberSchema) error {
	profileIDs := []string{}
	for _, email := range FindMemberObject.Email {
		profile, err := s.client.GetUserByEmail(email)
		if err != nil {
			return err
		}
		profileIDs = append(profileIDs, profile.ID)
	}
	var message slack.MsgOption
	var err error
	switch FindMemberObject.Type {
	case "project":
		message, err = slackmessages.ConstructProjectMessage(FindMemberObject, profileIDs)
	case "study":
		message, err = slackmessages.ConstructStudyMessage(FindMemberObject, profileIDs)
	default:
		return errors.New("invalid type")
	}
	if err != nil {
		return err
	}
	_, _, err = s.client.PostMessage(s.findMemberChannelID, message)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", s.findMemberChannelID, err)
		return err
	}
	return nil
}

// SendAlertToUser 스터디 팀 공고 지원자/팀장 메시지 전송
func (s *alertService) SendAlertToUser(UserObject models.UserMessageSchema) error {
	leaderProfile, err := s.client.GetUserByEmail(UserObject.LeaderEmail)
	if err != nil {
		return err
	}
	applicantProfile, err := s.client.GetUserByEmail(UserObject.ApplicantEmail)
	if err != nil {
		return err
	}
	leaderMsg, applicantMsg, err := slackmessages.ConstructApplicantAndLeaderMessage(leaderProfile, applicantProfile, UserObject)
	if err != nil {
		return err
	}
	_, _, err = s.client.PostMessage(leaderProfile.ID, leaderMsg)
	if err != nil {
		return err
	}
	_, _, err = s.client.PostMessage(applicantProfile.ID, applicantMsg)
	if err != nil {
		return err
	}
	return nil
}
