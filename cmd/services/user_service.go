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
