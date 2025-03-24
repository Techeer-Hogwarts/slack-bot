package services

import "github.com/Techeer-Hogwarts/slack-bot/cmd/repositories"

type WorkflowService interface {
	// Define methods for workflow service
}

type workflowService struct {
	workflowRepo repositories.WorkflowRepository
}

func NewWorkflowService(workflowRepo repositories.WorkflowRepository) *workflowService {
	return &workflowService{workflowRepo: workflowRepo}
}
