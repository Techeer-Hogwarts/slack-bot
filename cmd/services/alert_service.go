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
	client                 *slack.Client
	findMemberChannelID    string
	findMemberChannelIDDev string
}

func NewAlertService(client *slack.Client, findMemberChannelID, findMemberChannelIDDev string) *alertService {
	return &alertService{
		client:                 client,
		findMemberChannelID:    findMemberChannelID,
		findMemberChannelIDDev: findMemberChannelIDDev,
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
	var slackChannel string
	switch FindMemberObject.Environment {
	case "staging":
		slackChannel = s.findMemberChannelIDDev
	case "production":
		slackChannel = s.findMemberChannelID
	default:
		slackChannel = s.findMemberChannelIDDev
	}
	_, _, err = s.client.PostMessage(slackChannel, message)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", slackChannel, err)
		return err
	}
	return nil
}

// SendAlertToUser 스터디 팀 공고 지원자/팀장 메시지 전송
func (s *alertService) SendAlertToUser(UserObject models.UserMessageSchema) error {
	applicantProfile, err := s.client.GetUserByEmail(UserObject.ApplicantEmail)
	if err != nil {
		log.Printf("Applicant Email: %s", UserObject.ApplicantEmail)
		return err
	}
	leaderMsg, applicantMsg, err := slackmessages.ConstructApplicantAndLeaderMessage(applicantProfile, UserObject)
	if err != nil {
		return err
	}
	emails := len(UserObject.LeaderEmails)
	log.Printf("Leader Emails: %v", UserObject.LeaderEmails)
	for i := range emails {
		email := UserObject.LeaderEmails[i]
		leaderProfile, err := s.client.GetUserByEmail(email)
		if err != nil {
			log.Printf("Leader Email: %s", email)
			return err
		}
		_, _, err = s.client.PostMessage(leaderProfile.ID, leaderMsg)
		if err != nil {
			return err
		}
	}
	_, _, err = s.client.PostMessage(applicantProfile.ID, applicantMsg)
	if err != nil {
		return err
	}
	return nil
}
