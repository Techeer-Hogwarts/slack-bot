package handlers

import "github.com/Techeer-Hogwarts/slack-bot/cmd/services"

type WorkflowHanlder struct {
	service services.WorkflowService
}

func NewWorkflowHandler(service services.WorkflowService) *WorkflowHanlder {
	return &WorkflowHanlder{service: service}
}
