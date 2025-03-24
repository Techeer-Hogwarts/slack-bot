package handlers

import "github.com/Techeer-Hogwarts/slack-bot/cmd/services"

type Hanlder struct {
	GitHubHandler   *GitHubHandler
	UserHandler     *UserHandler
	SlackHandler    *SlackHandler
	WorkflowHanlder *WorkflowHanlder
}

func NewHandler(service *services.Service) *Hanlder {
	return &Hanlder{
		GitHubHandler:   NewGitHubHandler(service.GitHubService),
		UserHandler:     NewUserHandler(service.UserService),
		SlackHandler:    NewSlackHandler(service.SlackService),
		WorkflowHanlder: NewWorkflowHandler(service.WorkflowService),
	}
}
