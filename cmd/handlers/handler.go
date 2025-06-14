package handlers

import "github.com/Techeer-Hogwarts/slack-bot/cmd/services"

type Handler struct {
	AlertHandler   *AlertHandler
	ProfileHandler *ProfileHandler
	SlackHandler   *SlackHandler
	DeployHandler  *DeployHandler
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		AlertHandler:   NewAlertHandler(service.AlertService),
		ProfileHandler: NewProfileHandler(service.ProfileService),
		SlackHandler:   NewSlackHandler(service.SlackService, service.DeployService),
		DeployHandler:  NewDeployHandler(service.DeployService),
	}
}
