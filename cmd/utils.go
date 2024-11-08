package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/thomas-and-friends/slack-bot/config"
	"github.com/thomas-and-friends/slack-bot/db"
)

var (
	signingKey string
	botToken   string
	channelID  string
	roleMap    map[string]string
	secret     string
	sheetsID   string
)

const (
	emoji_people   = ":people_holding_hands:"
	emoji_golf     = ":golf:"
	emoji_star     = ":star2:"
	emoji_notebook = ":notebook:"
	emoji_stack    = ":hammer_and_pick:"
	emoji_dart     = ":dart:"
)

func init() {
	config.LoadEnv()
	signingKey = config.GetEnv("SLACK_SIGNING_SECRET", "")
	botToken = config.GetEnv("SLACK_BOT_TOKEN", "")
	secret = config.GetEnv("SECRET_ALLOW", "")
	channelID = config.GetEnv("CHANNEL_ID2", "")
	sheetsID = config.GetEnv("GSHEETS_ID", "")
	roleMap = map[string]string{
		"frontend":  "Frontend Developer",
		"backend":   "Backend Developer",
		"fullstack": "Fullstack Developer",
		"uxui":      "UX/UI Designer",
		"devops":    "OPS/SRE",
		"data":      "Data Engineer",
		"study":     "스터디",
		"etc":       "기타",
	}
}

func VerifySlackRequest(req *http.Request) error {
	s, err := slack.NewSecretsVerifier(req.Header, signingKey)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body)) // Reassign body after reading it

	if _, err := s.Write(body); err != nil {
		return err
	}

	if err := s.Ensure(); err != nil {
		return err
	}

	return nil
}

// func getChannelMessages(api *slack.Client, channelID string) (*slack.GetConversationHistoryResponse, error) {
// 	historyParams := slack.GetConversationHistoryParameters{
// 		ChannelID: channelID,
// 		Limit:     100,
// 	}

// 	history, err := api.GetConversationHistory(&historyParams)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return history, nil
// }

func constructMessageText(message FormMessage) (string, error) {
	if message.Description == "" || message.TeamName == "" || message.TechStacks == nil {
		return "", errors.New("missing required fields except team leader")
	} else if message.TeamLeader == "" {
		return "", errors.New("missing team leader")
	}
	return "[" + emoji_people + message.TeamIntro + emoji_people + "]\n" +
		"> " + emoji_golf + " *팀 이름* \n " + message.TeamName + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <<@" + message.TeamLeader + ">>\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명*\n" + message.Description + "\n\n\n\n" +
		"> " + emoji_stack + " *사용되는 기술*\n" + formatListStacks(message.TechStacks) + "\n\n\n\n" +
		"> " + emoji_dart + " *모집하는 직군 & 인원*\n" + formatListRoles(message) + "\n\n\n\n" +
		"> " + "*그 외 추가적인 정보* \n" + message.Etc + "\n\n자세한 문의사항은" + "<@" + message.TeamLeader + ">" + "에게 DM으로 문의 주세요!", nil
}

func formatListRoles(message FormMessage) string {
	var roles []string
	uxNum, err := strconv.Atoi(message.UxMembers)
	if err != nil {
		uxNum = 0
	}
	if uxNum <= 0 {
		log.Println("No UX/UI members")
	} else {
		roles = append(roles, " • "+roleMap["uxui"]+" ("+message.UxMembers+"명)\n")
	}
	frontNum, err := strconv.Atoi(message.FrontMembers)
	if err != nil {
		frontNum = 0
	}
	if frontNum <= 0 {
		log.Println("No Frontend members")
	} else {
		roles = append(roles, " • "+roleMap["frontend"]+" ("+message.FrontMembers+"명)\n")
	}
	backNum, err := strconv.Atoi(message.BackMembers)
	if err != nil {
		backNum = 0
	}
	if backNum <= 0 {
		log.Println("No Backend members")
	} else {
		roles = append(roles, " • "+roleMap["backend"]+" ("+message.BackMembers+"명)\n")
	}
	dataNum, err := strconv.Atoi(message.DataMembers)
	if err != nil {
		dataNum = 0
	}
	if dataNum <= 0 {
		log.Println("No Data members")
	} else {
		roles = append(roles, " • "+roleMap["data"]+" ("+message.DataMembers+"명)\n")
	}
	opsNum, err := strconv.Atoi(message.OpsMembers)
	if err != nil {
		opsNum = 0
	}
	if opsNum <= 0 {
		log.Println("No OPS/SRE members")
	} else {
		roles = append(roles, " • "+roleMap["devops"]+" ("+message.OpsMembers+"명)\n")
	}
	studyNum, err := strconv.Atoi(message.StudyMembers)
	if err != nil {
		studyNum = 0
	}
	if studyNum <= 0 {
		log.Println("No Study members")
	} else {
		roles = append(roles, " • "+roleMap["study"]+" ("+message.StudyMembers+"명)\n")
	}
	etcNum, err := strconv.Atoi(message.EtcMembers)
	if err != nil {
		etcNum = 0
	}
	if etcNum <= 0 {
		log.Println("No Etc members")
	} else {
		roles = append(roles, " • "+roleMap["etc"]+" ("+message.EtcMembers+"명)\n")
	}
	return strings.Join(roles, "")
}

func formatListStacks(items []string) string {
	var backendStacks []string
	var frontendStacks []string
	var devopsStacks []string
	var otherStacks []string
	if len(items) == 0 {
		return "None"
	}

	for _, stack := range items {
		tagName, tagType, _, err := db.GetTag(stack)
		if err != nil {
			log.Printf("Failed to get tag: %v", err)
		}
		stack_text := "`" + tagName + "`"
		if tagType == "backend" {
			backendStacks = append(backendStacks, stack_text)
		} else if tagType == "frontend" {
			frontendStacks = append(frontendStacks, stack_text)
		} else if tagType == "devops" {
			devopsStacks = append(devopsStacks, stack_text)
		} else {
			otherStacks = append(otherStacks, stack_text)
		}
	}
	joinedBackStacks := "*백엔드 기술:* " + strings.Join(backendStacks, ", ") + "\n"
	joinedFrontStacks := "*프런트엔드 기술:* " + strings.Join(frontendStacks, ", ") + "\n"
	joinedDevopsStacks := "*데브옵스 기술:* " + strings.Join(devopsStacks, ", ") + "\n"
	joinedOtherStacks := "*그 외 기술:* " + strings.Join(otherStacks, ", ") + "\n"
	return joinedBackStacks + joinedFrontStacks + joinedDevopsStacks + joinedOtherStacks
}

func InitialDataUsers() {
	api := slack.New(botToken)
	err := addAllUsers(api)
	if err != nil {
		log.Printf("Failed to add all users: %v", err)
	}
}

func InitialDataTags() {
	stacks, err := loadStacksFromFile("stacks.json")
	if err != nil {
		log.Printf("Failed to load stacks from file: %v", err)
	}
	err = addAllTags(stacks)
	if err != nil {
		log.Printf("Failed to add all tags: %v", err)
	}
}

func addAllUsers(api *slack.Client) error {
	users, err := api.GetUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		username := user.Profile.DisplayNameNormalized
		email := user.Profile.Email
		if username == "" {
			username = user.Profile.RealNameNormalized
		}
		if !user.IsBot && !user.Deleted && !user.IsAppUser && !user.IsOwner && user.ID != "USLACKBOT" {
			ms, _, err := db.GetUser(user.ID)
			if ms == "na" {
				err = db.AddUser(user.ID, username, email)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return fmt.Errorf("failed to get user: %s", err.Error())
			}
		}
	}
	return nil
}

func getAllUsers(api *slack.Client) error {
	users, err := api.GetUsers()
	if err != nil {
		return err
	}
	counter := 0
	for _, user := range users {
		username := user.Profile.DisplayNameNormalized
		if username == "" {
			username = user.Profile.RealNameNormalized
		}
		if !user.IsBot && !user.Deleted && !user.IsAppUser && !user.IsOwner && user.ID != "USLACKBOT" {
			log.Printf("User ID: %v | User Name: %v ", user.ID, username)
			counter++
		}
	}
	log.Printf("Total number of users: %v", counter)

	return nil
}

func addAllTags(stacks []db.Stack) error {
	for _, value := range stacks {
		ms, _, _, err := db.GetTag(value.Key)
		if ms == "na" {
			err := db.AddTag(value.Key, value.Name, value.Type)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return fmt.Errorf("failed to get tag: %s", err.Error())
		}
	}
	return nil
}

func loadStacksFromFile(filename string) ([]db.Stack, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var stacks []db.Stack
	if err := json.Unmarshal(bytes, &stacks); err != nil {
		return nil, err
	}

	return stacks, nil
}
func deleteMessage(api *slack.Client, payload slack.InteractionCallback) error {
	actionUserID := payload.User.ID
	actionMessageTimestamp := payload.Message.Timestamp

	teamObj, err := db.GetTeam(actionMessageTimestamp)
	if err != nil {
		if actionUserID == "U02AES3BH17" || actionUserID == "U033UTX061X" {
			_, _, err = api.DeleteMessage(payload.Channel.ID, actionMessageTimestamp)
			if err != nil {
				return fmt.Errorf("failed to delete message from Slack: %s", err.Error())
			}
			err = sendSuccessMessage(api, payload.Channel.ID, payload.User.ID, "메시지가 삭제되었습니다.")
			if err != nil {
				return err
			}
		} else {
			_ = sendFailMessage(api, payload.Channel.ID, payload.User.ID, "오류가 발생했습니다. 다시 시도해주시거나 개발자를 연락 하세요")
			return err
		}
	}
	if teamObj.TeamLeader == actionUserID || actionUserID == "U02AES3BH17" || actionUserID == "U033UTX061X" {
		err = db.DeleteTeam(actionMessageTimestamp)
		if err != nil {
			return err
		}
		_, _, err = api.DeleteMessage(payload.Channel.ID, actionMessageTimestamp)
		if err != nil {
			return fmt.Errorf("failed to delete message from Slack: %s", err.Error())
		}
		err = sendSuccessMessage(api, payload.Channel.ID, payload.User.ID, "메시지가 삭제되었습니다.")
		if err != nil {
			return err
		}
	} else {
		err = sendFailMessage(api, payload.Channel.ID, payload.User.ID, "팀 리더가 아닙니다. 삭제 권한이 없습니다.")
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func addTeamToDB(message FormMessage, ts string) error {
	jsonBytes, _ := json.Marshal(message)
	log.Println(string(jsonBytes))
	teamObj := db.Team{
		TeamType:   message.TeamType,
		TeamIntro:  message.TeamIntro,
		TeamName:   message.TeamName,
		TeamLeader: message.TeamLeader,
		TeamDesc:   message.Description,
		NumMembers: message.NumCurrentMembers,
		TeamEtc:    message.Etc,
		TeamTs:     ts,
	}
	teamID, err := db.AddTeam(teamObj)
	if err != nil {
		return err
	}
	for _, stack := range message.TechStacks {
		_, _, stack_id, err := db.GetTag(stack)
		if err != nil {
			return err
		}
		err = db.AddTagsToTeam(teamID, stack_id)
		if err != nil {
			return err
		}
	}
	for _, user := range message.Members {
		_, user_id, err := db.GetUser(user)
		if err != nil {
			return err
		}
		err = db.AddUserToTeam(teamID, user_id)
		if err != nil {
			return err
		}
	}
	uxWant, err := strconv.Atoi(message.UxMembers)
	if err != nil {
		uxWant = 0
	}
	frontWant, err := strconv.Atoi(message.FrontMembers)
	if err != nil {
		frontWant = 0
	}
	backWant, err := strconv.Atoi(message.BackMembers)
	if err != nil {
		backWant = 0
	}
	dataWant, err := strconv.Atoi(message.DataMembers)
	if err != nil {
		dataWant = 0
	}
	devopsWant, err := strconv.Atoi(message.OpsMembers)
	if err != nil {
		devopsWant = 0
	}
	studyWant, err := strconv.Atoi(message.StudyMembers)
	if err != nil {
		studyWant = 0
	}
	etcWant, err := strconv.Atoi(message.EtcMembers)
	if err != nil {
		etcWant = 0
	}
	messageObj := db.ExtraMessage{
		TeamID:       teamID,
		MessageTS:    ts,
		UXWant:       uxWant,
		FrontendWant: frontWant,
		BackendWant:  backWant,
		DataWant:     dataWant,
		DevopsWant:   devopsWant,
		StudyWant:    studyWant,
		EtcWant:      etcWant,
	}
	err = db.AddExtraMessage(messageObj)
	if err != nil {
		return err
	}
	return nil
}

func sendSuccessMessage(api *slack.Client, channelID string, userID string, messageText string) error {
	_, err := api.PostEphemeral(channelID, userID, slack.MsgOptionText(messageText, false))
	return err
}

func sendFailMessage(api *slack.Client, channelID string, userID string, messageText string) error {
	_, err := api.PostEphemeral(channelID, userID, slack.MsgOptionText(messageText, false))
	return err
}

func enrollUser(api *slack.Client, value string, channelID string) (string, error) {
	values := strings.Split(value, "|")
	applicantID := values[0]
	role := values[2]
	teamID, err := strconv.Atoi(values[1])
	if err != nil {
		return "", err
	}
	teamObj, err := db.GetTeamByID(teamID)
	if err != nil {
		return "", err
	}
	_, applicantIDInt, err := db.GetUser(applicantID)
	if err != nil {
		return "", err
	}
	flag, _ := db.GetUserInTeam(applicantIDInt, teamID)
	if flag {
		log.Println("User already in team")
		err = sendFailMessage(api, channelID, teamObj.TeamLeader, "이미 팀에 속해있습니다.")
		return "", err
	}
	err = db.AddUserToTeam(teamID, applicantIDInt)
	if err != nil {
		return "", err
	}
	err = db.UpdateTeamMembers(teamID, teamObj.NumMembers+1)
	if err != nil {
		return "", err
	}
	err = db.UpdateExtraMessage(role, teamObj.TeamTs)
	if err != nil {
		return "", err
	}
	msgText := fmt.Sprintf("<@%s>님의 %s 팀 가입 신청을 수락하셨습니다.", applicantID, teamObj.TeamName)
	// _, err = api.PostEphemeral(channelID, teamObj.TeamLeader, slack.MsgOptionText(msgText, false))
	err = sendDMSuccessMessage(api, teamObj.TeamLeader, msgText)
	if err != nil {
		return "", err
	}
	msgText = fmt.Sprintf("%v 팀 가입 신청이 수락되었습니다.", teamObj.TeamName)
	err = sendDMSuccessMessage(api, applicantID, msgText)
	if err != nil {
		return "", err
	}
	return teamObj.TeamTs, nil
}
