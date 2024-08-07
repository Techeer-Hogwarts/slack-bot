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

type ExtraMessage struct {
	TeamID       int    `json:"team_id"`
	MessageTS    string `json:"message_ts"`
	UXWant       int    `json:"ux_want"`
	FrontendWant int    `json:"frontend_want"`
	BackendWant  int    `json:"backend_want"`
	DataWant     int    `json:"data_want"`
	DevopsWant   int    `json:"devops_want"`
	StudyWant    int    `json:"study_want"`
	EtcWant      int    `json:"etc_want"`
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

func AddExtraMessage(message ExtraMessage) error {
	_, err := DBMain.Exec("INSERT INTO messages (team_id, message_ts, ux_want, frontend_want, backend_want, data_want, devops_want, study_want, etc_want) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		message.TeamID, message.MessageTS, message.UXWant, message.FrontendWant, message.BackendWant, message.DataWant, message.DevopsWant, message.StudyWant, message.EtcWant)
	if err != nil {
		return fmt.Errorf("failed to insert new extra message: %s", err.Error())
	}
	log.Println("Extra message added to the database")
	return nil
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

func DeactivateRecruitTeam(ts string) error {
	// Deactivat a team from the database
	_, err := DBMain.Exec("UPDATE teams SET recruit_active = FALSE WHERE message_ts = $1", ts)
	if err != nil {
		return fmt.Errorf("failed to mark team as recruit deactive: %s", err.Error())
	}
	log.Printf("Team with message_ts %s marked as recruit deactive in the database", ts)
	return nil
}

func ActivateRecruitTeam(ts string) error {
	// Activate a team from the database
	_, err := DBMain.Exec("UPDATE teams SET recruit_active = TRUE WHERE message_ts = $1", ts)
	if err != nil {
		return fmt.Errorf("failed to mark team as recruit active: %s", err.Error())
	}
	log.Printf("Team with message_ts %s marked as recruit active in the database", ts)
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

func GetExtraMessage(ts string) (ExtraMessage, error) {
	// Get a team message from the database
	messageObj := ExtraMessage{}
	err := DBMain.QueryRow(
		"SELECT team_id, message_ts, ux_want, frontend_want, backend_want, data_want, devops_want, study_want, etc_want FROM messages WHERE message_ts = $1",
		ts,
	).Scan(
		&messageObj.TeamID,
		&messageObj.MessageTS,
		&messageObj.UXWant,
		&messageObj.FrontendWant,
		&messageObj.BackendWant,
		&messageObj.DataWant,
		&messageObj.DevopsWant,
		&messageObj.StudyWant,
		&messageObj.EtcWant,
	)
	if err == sql.ErrNoRows {
		return ExtraMessage{}, fmt.Errorf("message not found")
	}
	log.Printf("Extra message found in the database")
	return messageObj, nil
}

func UpdateExtraMessage(role string, ts string) error {
	if role == "backend" {
		_, err := DBMain.Exec("UPDATE messages SET backend_want = backend_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "frontend" {
		_, err := DBMain.Exec("UPDATE messages SET frontend_want = frontend_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "uxui" {
		_, err := DBMain.Exec("UPDATE messages SET ux_want = ux_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "devops" {
		_, err := DBMain.Exec("UPDATE messages SET devops_want = devops_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "data" {
		_, err := DBMain.Exec("UPDATE messages SET data_want = data_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "study" {
		_, err := DBMain.Exec("UPDATE messages SET study_want = study_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	} else if role == "etc" {
		_, err := DBMain.Exec("UPDATE messages SET etc_want = etc_want - 1 WHERE message_ts = $1", ts)
		if err != nil {
			return fmt.Errorf("failed to update extra message: %s", err.Error())
		}
	}
	return nil
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

// func GetUsersInTeam(teamId int) ([]int, error) {
// 	// Get all users in a team

// }

func GetUserInTeam(userID int, teamID int) (bool, error) {
	// Get a user in a team
	var rowID int
	err := DBMain.QueryRow("SELECT ut_id FROM user_teams WHERE user_id = $1 AND team_id = $2", userID, teamID).Scan(&rowID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return true, fmt.Errorf("failed to get user in team: %s", err)
	}
	return true, nil
}

func GetAllActiveTeams() ([]Team, error) {
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

func GetAllRecruitingTeams() ([]Team, error) {
	// Get all teams from the database
	rows, err := DBMain.Query("SELECT team_id, team_type, team_intro, team_name, team_leader, team_description, num_members, team_etc, message_ts FROM teams WHERE is_active = TRUE AND recruit_active = TRUE")
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

func GetTagsFromTeam(teamID int) ([]string, error) {
	// Get all tags from the database
	rows, err := DBMain.Query("SELECT t.tag_name FROM tags t JOIN team_tags tt ON t.tag_id = tt.tag_id WHERE tt.team_id = $1", teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %s", err.Error())
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tagObj string
		err := rows.Scan(
			&tagObj,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %s", err.Error())
		}
		tags = append(tags, tagObj)
	}
	return tags, nil
}
