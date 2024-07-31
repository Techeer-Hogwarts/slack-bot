package db

import (
	"database/sql"
	"fmt"
	"log"
)

func AddUser(userCode string, userName string) error {
	_, err := DBMain.Exec("INSERT INTO users (user_code, user_name) VALUES ($1, $2)", userCode, userName)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %s", err.Error())
	}
	log.Printf("User %s added to the database", userName)
	return nil
}

func GetUser(userCode string) (string, error) {
	var userName string
	err := DBMain.QueryRow("SELECT user_name FROM users WHERE user_code = $1", userCode).Scan(&userName)
	if err == sql.ErrNoRows {
		return "na", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("some other sql error: %s", err.Error())
	}
	log.Printf("User %s found in the database", userName)
	return userName, nil
}

func AddTeam() {
	// Add a team to the database
}

func GetTeams() {
	// Get all teams from the database
}

func AddUserToTeam() {
	// Add a user to a team
}

func GetUsersByTeam() {
	// Get all users in a team
}

func GetTeamPost() {
	// Get a team post
}
