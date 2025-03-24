package services

import (
	"github.com/Techeer-Hogwarts/slack-bot/cmd/repositories"
	"github.com/slack-go/slack"
)

type Service struct {
	UserService     UserService
	SlackService    SlackService
	GitHubService   GitHubService
	WorkflowService WorkflowService
}

// NewService creates a new instance of Service with all required services.
func NewService(repo *repositories.Repository, slackClient *slack.Client, githubURL, githubToken string) *Service {
	return &Service{
		UserService:     NewUserService(repo.UserRepository),
		SlackService:    NewSlackService(slackClient),
		GitHubService:   NewGitHubService(githubURL, githubToken),
		WorkflowService: NewWorkflowService(repo.WorkflowRepository),
	}
}
