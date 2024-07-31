package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Team struct {
	TeamID     string
	TeamType   string
	TeamIntro  string
	TeamName   string
	TeamLeader string
	TeamDesc   string
	NumMembers int
	TeamEtc    string
	TeamTs     string
}

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

func AddTeam(teamobj Team) error {
	_, err := DBMain.Exec("INSERT INTO teams (team_type, team_intro, team_name, team_leader, team_desc, num_members, team_etc, team_ts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", teamobj.TeamType, teamobj.TeamIntro, teamobj.TeamName, teamobj.TeamLeader, teamobj.TeamDesc, teamobj.NumMembers, teamobj.TeamEtc, teamobj.TeamTs)
	if err != nil {
		return fmt.Errorf("failed to insert new team: %s", err.Error())
	}
	log.Printf("Team %s added to the database", teamobj.TeamName)
	return nil
}

func DeleteTeam() {
	// Delete a team from the database
}

func GetTeam(ts string) (Team, error) {
	// Get a team from the database
	teamObj := Team{}
	err := DBMain.QueryRow("SELECT * FROM teams WHERE team_ts = $1", ts).Scan(&teamObj.TeamID, &teamObj.TeamType, &teamObj.TeamIntro, &teamObj.TeamName, &teamObj.TeamLeader, &teamObj.TeamDesc, &teamObj.NumMembers, &teamObj.TeamEtc, &teamObj.TeamTs)
	if err == sql.ErrNoRows {
		return Team{}, fmt.Errorf("team not found")
	}
	return teamObj, nil
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

func GetTag(key string) (string, error) {
	var tag string
	err := DBMain.QueryRow("SELECT tag_long_name FROM tags WHERE tag_name = $1", key).Scan(&tag)
	if err == sql.ErrNoRows {
		return "na", fmt.Errorf("tag not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get tag: %s", err.Error())
	}
	log.Printf("Tag %s found in the database", tag)
	return tag, nil
}

func AddTagsToTeam() {
	// Add tags to a team post
}

func AddTag(key string, value string) error {
	_, err := DBMain.Exec("INSERT INTO tags (tag_name, tag_long_name) VALUES ($1, $2)", key, value)
	if err != nil {
		return fmt.Errorf("failed to insert new tag: %s", err.Error())
	}
	log.Printf("Tag %s added to the database", key)
	return nil
}
