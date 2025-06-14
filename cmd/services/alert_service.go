package services

import (
	"log"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/slack-go/slack"
)

type AlertService interface {
	SendAlert(channelID, message string) error
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
func (s *alertService) SendAlert(channelID, message string) error {
	log.Printf("channelID: %s, message: %s", channelID, message)
	_, _, err := s.client.PostMessage(channelID, slack.MsgOptionText(message, false))
	return err
}

// SendAlertToFindMember 스터디 팀 공고 메시지 전송 (/alert/find-member 와 /alert/channel 에서 사용)
func (s *alertService) SendAlertToFindMember(FindMemberObject models.FindMemberSchema) error {
	log.Printf("FindMemberObject: %+v", FindMemberObject)
	_, _, err := s.client.PostMessage(s.findMemberChannelID, slack.MsgOptionText("Some Message", false))
	return err
}

// SendAlertToUser 스터디 팀 공고 지원자/팀장 메시지 전송
func (s *alertService) SendAlertToUser(UserObject models.UserMessageSchema) error {
	log.Printf("UserObject: %+v", UserObject)
	_, _, err := s.client.PostMessage(s.findMemberChannelID, slack.MsgOptionText("Some Message", false))
	return err
}
