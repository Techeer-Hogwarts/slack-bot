package services

import (
	"errors"
	"log"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/auth"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(email, password string) (string, error)
	ResetPassword(userID int, currentPass, newPass string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) ResetPassword(userID int, currentPass, newPass string) error {
	// Retrieve user by ID
	storedHash, err := s.userRepo.GetUserPasswordHash(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(currentPass)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password in database
	return s.userRepo.UpdateUserPassword(userID, string(newHash))
}
