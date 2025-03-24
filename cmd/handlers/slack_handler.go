package handlers

import "github.com/Techeer-Hogwarts/slack-bot/cmd/services"

type SlackHandler struct {
	service services.SlackService
}

func NewSlackHandler(service services.SlackService) *SlackHandler {
	return &SlackHandler{service: service}
}
