package db

import "fmt"

func AddUser(userID string, userName string) error {
	_, err := DBMain.Exec("INSERT INTO users (user_code, user_name) VALUES ($1, $2)", userID, userName)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %s", err.Error())
	}
	return nil

}

func GetUser(userID string) (string, error) {
	var userName string
	err := DBMain.QueryRow("SELECT user_name FROM users WHERE user_code = $1", userID).Scan(&userName)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %s", err.Error())
	}
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
