package services

import (
	"github.com/slack-go/slack"
)

type Service struct {
	SlackService   SlackService
	DeployService  DeployService
	ProfileService ProfileService
	AlertService   AlertService
}

// NewService creates a new instance of Service with all required services.
func NewService(slackClient *slack.Client, githubURL, githubToken, cicdChannelID, findMemberChannelID, findMemberChannelIDDev string) *Service {
	return &Service{
		AlertService:   NewAlertService(slackClient, findMemberChannelID, findMemberChannelIDDev),
		SlackService:   NewSlackService(slackClient),
		DeployService:  NewDeployService(slackClient, githubURL, githubToken, cicdChannelID),
		ProfileService: NewProfileService(slackClient),
	}
}
