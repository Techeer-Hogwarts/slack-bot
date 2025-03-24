package repositories

import (
	"database/sql"
	"errors"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := "SELECT id, email, password FROM users WHERE email = $1"

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}
