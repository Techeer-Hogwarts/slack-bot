package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Team struct {
	TeamID         string
	TeamType       string
	TeamIntro      string
	TeamName       string
	TeamLeader     string
	TeamLeaderName string
	TeamDesc       string
	NumMembers     int
	TeamEtc        string
	TeamTs         string
}

type Stack struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func AddUser(userCode string, userName string, email string) error {
	_, err := DBMain.Exec("INSERT INTO users (user_code, user_name, user_email) VALUES ($1, $2, $3)", userCode, userName, email)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %s", err.Error())
	}
	log.Printf("User %s added to the database", userName)
	return nil
}

func GetUser(userCode string) (string, int, error) {
	var userName string
	var userID int
	err := DBMain.QueryRow("SELECT user_name, user_id FROM users WHERE user_code = $1", userCode).Scan(&userName, &userID)
	if err == sql.ErrNoRows {
		return "na", 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return "", 0, fmt.Errorf("some other sql error: %s", err.Error())
	}
	log.Printf("User %s found in the database", userName)
	return userName, userID, nil
}

func GetUserWithID(userID int) (string, string, error) {
	var userName string
	var userCode string
	err := DBMain.QueryRow("SELECT user_name, user_code FROM users WHERE user_id = $1", userID).Scan(&userName, &userCode)
	if err == sql.ErrNoRows {
		return "na", "", fmt.Errorf("user not found. Error content: %v", err)
	}
	if err != nil {
		return "", "", fmt.Errorf("some other sql error: %s", err.Error())
	}
	log.Printf("User %s found in the database", userName)
	return userName, userCode, nil
}

func AddTeam(teamobj Team) (int, error) {
	var teamID int
	_, leader_id, _ := GetUser(teamobj.TeamLeader)
	err := DBMain.QueryRow(
		"INSERT INTO teams (team_type, team_intro, team_name, team_leader, team_description, num_members, team_etc, message_ts) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING team_id",
		teamobj.TeamType, teamobj.TeamIntro, teamobj.TeamName, leader_id, teamobj.TeamDesc, teamobj.NumMembers, teamobj.TeamEtc, teamobj.TeamTs,
	).Scan(&teamID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new team: %s", err.Error())
	}
	log.Printf("Team %s added to the database", teamobj.TeamName)
	return teamID, nil
}

func DeleteTeam(ts string) error {
	// Delete a team from the database
	_, err := DBMain.Exec("UPDATE teams SET is_active = FALSE WHERE message_ts = $1", ts)
	if err != nil {
		return fmt.Errorf("failed to mark team as inactive: %s", err.Error())
	}
	log.Printf("Team with message_ts %s marked as inactive in the database", ts)
	return nil
}

func GetTeam(ts string) (Team, error) {
	// Get a team from the database
	teamObj := Team{}
	var teamLeaderID int
	err := DBMain.QueryRow(
		"SELECT team_id, team_type, team_intro, team_name, team_leader, team_description, num_members, team_etc, message_ts FROM teams WHERE message_ts = $1 AND is_active = TRUE",
		ts,
	).Scan(
		&teamObj.TeamID,
		&teamObj.TeamType,
		&teamObj.TeamIntro,
		&teamObj.TeamName,
		&teamLeaderID,
		&teamObj.TeamDesc,
		&teamObj.NumMembers,
		&teamObj.TeamEtc,
		&teamObj.TeamTs,
	)
	if err == sql.ErrNoRows {
		return Team{}, fmt.Errorf("team not found")
	}
	_, teamLeaderCode, err := GetUserWithID(teamLeaderID)
	if err == sql.ErrNoRows {
		return Team{}, fmt.Errorf("leader not found")
	}
	if err != nil {
		return Team{}, fmt.Errorf("failed to get leader for team: %s", err.Error())
	}
	log.Printf("Team %s found in the database", teamObj.TeamName)
	teamObj.TeamLeader = teamLeaderCode
	return teamObj, nil
}

func GetTeamByID(teamID int) (Team, error) {
	// Get a team from the database
	teamObj := Team{}
	var teamLeaderID int
	err := DBMain.QueryRow(
		"SELECT team_id, team_type, team_intro, team_name, team_leader, team_description, num_members, team_etc, message_ts FROM teams WHERE team_id = $1 AND is_active = TRUE",
		teamID,
	).Scan(
		&teamObj.TeamID,
		&teamObj.TeamType,
		&teamObj.TeamIntro,
		&teamObj.TeamName,
		&teamLeaderID,
		&teamObj.TeamDesc,
		&teamObj.NumMembers,
		&teamObj.TeamEtc,
		&teamObj.TeamTs,
	)
	if err == sql.ErrNoRows {
		return Team{}, fmt.Errorf("team not found")
	}
	_, teamLeaderCode, err := GetUserWithID(teamLeaderID)
	if err == sql.ErrNoRows {
		return Team{}, fmt.Errorf("leader not found")
	}
	if err != nil {
		return Team{}, fmt.Errorf("failed to get leader for team: %s", err.Error())
	}
	log.Printf("Team %s found in the database", teamObj.TeamName)
	teamObj.TeamLeader = teamLeaderCode
	return teamObj, nil
}

func UpdateTeamMembers(teamID int, numMembers int) error {
	_, err := DBMain.Exec("UPDATE teams SET num_members = $1 WHERE team_id = $2", numMembers, teamID)
	if err != nil {
		return fmt.Errorf("failed to update team members: %s", err.Error())
	}
	log.Printf("Team %d now has %d members", teamID, numMembers)
	return nil
}

func AddUserToTeam(teamID int, userID int) error {
	_, err := DBMain.Exec("INSERT INTO user_teams (team_id, user_id) VALUES ($1, $2)", teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to insert new user to team: %s", err.Error())
	}
	log.Printf("User %d added to team %d", userID, teamID)
	return nil
}

func GetUsersInTeam() {
	// Get all users in a team
}

func GetUserInTeam(userID int, teamID int) (bool, error) {
	// Get a user in a team
	rows, err := DBMain.Query("SELECT * FROM user_teams WHERE user_id = $1 AND team_id = $2", userID, teamID)
	if err != nil {
		return true, fmt.Errorf("failed to get user in team: %s", err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		return false, fmt.Errorf("user %d is not in team %d", userID, teamID)
	}
	log.Println("User is in team")
	return true, nil
}

func GetAllTeams() ([]Team, error) {
	// Get all teams from the database
	rows, err := DBMain.Query("SELECT team_id, team_type, team_intro, team_name, team_leader, team_description, num_members, team_etc, message_ts FROM teams WHERE is_active = TRUE")
	if err != nil {
		return []Team{}, fmt.Errorf("failed to get all teams: %s", err.Error())
	}
	defer rows.Close()
	teams := []Team{}
	for rows.Next() {
		teamObj := Team{}
		var teamLeaderID int
		err := rows.Scan(
			&teamObj.TeamID,
			&teamObj.TeamType,
			&teamObj.TeamIntro,
			&teamObj.TeamName,
			&teamLeaderID,
			&teamObj.TeamDesc,
			&teamObj.NumMembers,
			&teamObj.TeamEtc,
			&teamObj.TeamTs,
		)
		if err != nil {
			return []Team{}, fmt.Errorf("failed to scan team: %s", err.Error())
		}
		_, teamLeaderCode, err := GetUserWithID(teamLeaderID)
		if err == sql.ErrNoRows {
			return []Team{}, fmt.Errorf("leader not found")
		}
		if err != nil {
			return []Team{}, fmt.Errorf("failed to get leader for team: %s", err.Error())
		}
		teamObj.TeamLeader = teamLeaderCode
		teams = append(teams, teamObj)
	}
	return teams, nil
}

func GetTag(key string) (string, string, int, error) {
	var tagName string
	var tagID int
	var tagType string
	err := DBMain.QueryRow("SELECT tag_long_name, tag_id, tag_type FROM tags WHERE tag_name = $1", key).Scan(&tagName, &tagID, &tagType)
	if err == sql.ErrNoRows {
		return "na", "", 0, nil
	} else if err != nil {
		return "", "", 0, fmt.Errorf("failed to get tag: %s", err.Error())
	}
	log.Printf("Tag %s found in the database", tagName)
	return tagName, tagType, tagID, nil
}

func AddTag(key string, value string, tagType string) error {
	_, err := DBMain.Exec("INSERT INTO tags (tag_name, tag_long_name, tag_type) VALUES ($1, $2, $3)", key, value, tagType)
	if err != nil {
		return fmt.Errorf("failed to insert new tag: %s", err.Error())
	}
	log.Printf("Tag %s added to the database", key)
	return nil
}

func AddTagsToTeam(teamID int, tag int) error {
	_, err := DBMain.Exec("INSERT INTO team_tags (team_id, tag_id) VALUES ($1, $2)", teamID, tag)
	if err != nil {
		return fmt.Errorf("failed to insert new tag to team: %s", err.Error())
	}
	return nil
}
