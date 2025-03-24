package repositories

import "database/sql"

type Repository struct {
	UserRepository     UserRepository
	WorkflowRepository WorkflowRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:     NewUserRepository(db),
		WorkflowRepository: NewWorkflowRepository(db),
	}
}
