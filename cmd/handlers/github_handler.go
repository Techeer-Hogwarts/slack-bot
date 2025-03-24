package handlers

import "github.com/Techeer-Hogwarts/slack-bot/cmd/services"

type GitHubHandler struct {
	service services.GitHubService
}

func NewGitHubHandler(service services.GitHubService) *GitHubHandler {
	return &GitHubHandler{service: service}
}
