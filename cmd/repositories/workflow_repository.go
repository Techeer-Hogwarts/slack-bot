package repositories

import "database/sql"

type WorkflowRepository interface {
	// Define methods for user repository
}

type workflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) WorkflowRepository {
	return &workflowRepository{db: db}
}
