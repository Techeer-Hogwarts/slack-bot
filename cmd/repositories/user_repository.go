package repositories

import (
	"database/sql"
	"errors"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserPasswordHash(userID int) (string, error)
	UpdateUserPassword(userID int, newPasswordHash string) error
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

func (r *userRepository) GetUserPasswordHash(userID int) (string, error) {
	var passwordHash string
	err := r.db.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (r *userRepository) UpdateUserPassword(userID int, newPasswordHash string) error {
	_, err := r.db.Exec("UPDATE users SET password = $1 WHERE id = $2", newPasswordHash, userID)
	return err
}
