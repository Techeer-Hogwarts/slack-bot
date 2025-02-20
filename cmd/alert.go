package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

type studySchema struct {
	ID             int    `json:"id"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	StudyExplain   string `json:"studyExplain"`
	RecruitNum     int    `json:"recruitNum"`
	Leader         string `json:"leader"`
	Email          string `json:"email"`
	RecruitExplain string `json:"recruitExplain"`
	NotionLink     string `json:"notionLink"`
	Goal           string `json:"goal"`
	Rule           string `json:"rule"`
}

type projectSchema struct {
	ID             int      `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	ProjectExplain string   `json:"projectExplain"`
	FrontNum       int      `json:"frontNum"`
	BackNum        int      `json:"backNum"`
	DataEngNum     int      `json:"dataEngNum"`
	DevOpsNum      int      `json:"devOpsNum"`
	uiUxNum        int      `json:"uiUxNum"`
	Leader         string   `json:"leader"`
	Email          string   `json:"email"`
	RecruitExplain string   `json:"recruitExplain"`
	NotionLink     string   `json:"notionLink"`
	Stack          []string `json:"stack"`
}

type userMessageSchema struct {
	TeamID         int    `json:"teamId"`
	TeamName       string `json:"teamName"`
	Type           string `json:"type"`
	LeaderEmail    string `json:"leaderEmail"`
	ApplicantEmail string `json:"applicantEmail"`
	Result         string `json:"result"`
}

const (
	redirectURL = "https://www.techeerzip.cloud/project/detail/%s/%d"
)

func AlertChannelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Alert Channel Handler")
	var temp map[string]interface{}
	requestBody := r.Body

	err := json.NewDecoder(requestBody).Decode(&temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	jsonVal, _ := json.MarshalIndent(temp, "", "  ")
	log.Println(string(jsonVal))
	typeStr, ok := temp["type"].(string)
	if !ok {
		http.Error(w, "Missing or invalid 'type' field", http.StatusBadRequest)
		log.Println("Missing 'type' field")
		return
	}
	api := slack.New(botToken)
	switch typeStr {
	case "study":
		var study studySchema
		err := mapToStruct(temp, &study)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error mapping to studySchema:", err)
			return
		}
		err = sendStudyMessage(study, api, channelID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error sending study message:", err)
			return
		}

	case "project":
		var project projectSchema
		err := mapToStruct(temp, &project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error mapping to projectSchema:", err)
			return
		}
		err = sendProjectMessage(project, api, channelID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error sending project message:", err)
			return
		}
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		log.Println("Invalid type:", typeStr)
		return
	}
	// if temp.Secret != secret {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	log.Println("Unauthorized")
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert Channel Handler"))
}

func AlertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Alert User Handler")
	var temp userMessageSchema
	requestBody := r.Body
	log.Println(requestBody)
	err := json.NewDecoder(requestBody).Decode(&temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	jsonVal, _ := json.MarshalIndent(temp, "", "  ")
	log.Println(string(jsonVal))

	result := temp.Result
	var statusMsg string
	api := slack.New(botToken)

	switch result {
	case "PENDING":
		statusMsg = "지원이 완료됐습니다."
	case "CANCELLED":
		statusMsg = "지원자께서 취소 하셨습니다."
	case "REJECT":
		statusMsg = "거절 돼서 팀에 합류하지 못하셨습니다."
	case "APPROVED":
		statusMsg = "수락 돼서 팀에 합류하셨습니다!"
	default:
		http.Error(w, "Invalid result", http.StatusBadRequest)
		log.Println("Invalid result:", result)
		return
	}
	err = sendUserStatusMessage(statusMsg, temp, api)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error sending user message:", err)
		return
	}
	err = sendLeaderStatusMessage(statusMsg, temp, api)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error sending leader message:", err)
		return
	}
	// if temp.Secret != secret {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	log.Println("Unauthorized")
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert User Handler"))
}

func mapToStruct(m map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func sendProjectMessage(project projectSchema, api *slack.Client, channelID string) error {
	profile, err := api.GetUserByEmail(project.Email)
	if err != nil {
		log.Printf("Failed to get user by email %s: %v", project.Email, err)
		return err
	}
	userCode := profile.ID
	projectMessage := "[" + emoji_people + " *새로운 프로젝트 팀 공고가 올라왔습니다* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + project.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <@" + userCode + ">\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명입니다*\n" + project.ProjectExplain + "\n\n\n\n" +
		"> " + ":woman-raising-hand:" + " *이런 사람을 원합니다!*\n" + project.RecruitExplain + "\n\n\n\n" +
		"> " + emoji_stack + " *사용되는 기술입니다*\n" + convertStackToEmojiString(project.Stack) + "\n\n\n" +
		"> " + emoji_dart + " *모집하는 직군 & 인원*\n" + convertRecruitNumToEmojiString(project) + "\n\n\n\n" +
		"> " + ":notion:" + "*노션 링크* \n" + project.NotionLink + "\n\n자세한 문의사항은" + "<@" + userCode + ">" + "에게 DM으로 문의 주세요!"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", projectMessage, false, false), nil, nil)
	applyButton := slack.NewButtonBlockElement("", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false))
	applyButton.URL = fmt.Sprintf(redirectURL, project.Type, project.ID)
	deleteButton := slack.NewButtonBlockElement("delete_button2", project.Email, slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	_, _, err = api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	return nil
}

func sendStudyMessage(study studySchema, api *slack.Client, channelID string) error {
	profile, err := api.GetUserByEmail(study.Email)
	if err != nil {
		log.Printf("Failed to get user by email %s: %v", study.Email, err)
		return err
	}
	userCode := profile.ID
	studyMessage := "[" + emoji_people + " *새로운 스터디 팀 공고가 올라왔습니다* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + study.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <@" + userCode + ">\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명입니다*\n" + study.StudyExplain + "\n\n\n\n" +
		"> " + ":man-raising-hand:" + " *이런 사람을 원합니다!*\n" + study.RecruitExplain + "\n\n\n\n" +
		"> " + ":pencil:" + " *지켜야 하는 규칙입니다!*\n" + study.Rule + "\n\n\n" +
		"> " + emoji_dart + " *모집하는 스터디 인원*\n" + strconv.Itoa(study.RecruitNum) + "명\n\n\n\n" +
		"> " + ":notion:" + " *노션 링크* \n" + study.NotionLink + "\n\n자세한 문의사항은" + "<@" + userCode + ">" + "에게 DM으로 문의 주세요!"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", studyMessage, false, false), nil, nil)
	applyButton := slack.NewButtonBlockElement("", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 팀 지원하기!", false, false))
	applyButton.URL = fmt.Sprintf(redirectURL, study.Type, study.ID)
	deleteButton := slack.NewButtonBlockElement("delete_button2", study.Email, slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	_, _, err = api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return err
	}
	return nil
}

func sendUserStatusMessage(status string, userMessage userMessageSchema, api *slack.Client) error {
	profile, err := api.GetUserByEmail(userMessage.ApplicantEmail)
	if err != nil {
		log.Printf("Failed to get user by email %s: %v", userMessage.ApplicantEmail, err)
		return err
	}
	msg := "[" + emoji_people + " *지원 결과 알림* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + userMessage.TeamName + "\n\n\n\n" +
		"> " + emoji_star + " *지원자:* <@" + profile.ID + ">\n\n\n\n" +
		"> " + emoji_notebook + " *지원 결과:* " + status + "\n\n\n\n" +
		"> " + emoji_dart + " *링크* \n" + fmt.Sprintf(redirectURL, userMessage.Type, userMessage.TeamID) + "\n\n\n\n"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", msg, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section)
	_, _, err = api.PostMessage(profile.ID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", profile.ID, err)
		return err
	}
	return nil
}

func sendLeaderStatusMessage(status string, userMessage userMessageSchema, api *slack.Client) error {
	profile, err := api.GetUserByEmail(userMessage.LeaderEmail)
	if err != nil {
		log.Printf("Failed to get user by email %s: %v", userMessage.LeaderEmail, err)
		return err
	}
	msg := "[" + emoji_people + " *지원 결과 알림* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + userMessage.TeamName + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <@" + profile.ID + ">\n\n\n\n" +
		"> " + emoji_notebook + " *지원 결과입니다:* " + status + "\n\n\n\n" +
		"> " + emoji_dart + " *링크* \n" + fmt.Sprintf(redirectURL, userMessage.Type, userMessage.TeamID) + "\n\n\n\n"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", msg, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section)
	_, _, err = api.PostMessage(profile.ID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", profile.ID, err)
		return err
	}
	return nil
}

func convertStackToEmojiString(stack []string) string {
	var backArray []string
	var frontArray []string
	var devOpsArray []string
	var otherArray []string
	var stackString string

	for _, s := range stack {
		category := categoryMap[s]
		emoji := stackMap[s]
		switch category {
		case "BACKEND":
			backArray = append(backArray, emoji)
		case "FRONTEND":
			frontArray = append(frontArray, emoji)
		case "DEVOPS":
			devOpsArray = append(devOpsArray, emoji)
		case "OTHER":
			otherArray = append(otherArray, emoji)
		default:
			log.Printf("Unknown category: %s", category)
			otherArray = append(otherArray, s)
		}
	}
	if len(backArray) > 0 {
		stackString += ":backend:" + " : " + strings.Join(backArray, " ") + "\n"
	}
	if len(frontArray) > 0 {
		stackString += ":frontend:" + " : " + strings.Join(frontArray, " ") + "\n"
	}
	if len(devOpsArray) > 0 {
		stackString += ":devops:" + " : " + strings.Join(devOpsArray, " ") + "\n"
	}
	if len(otherArray) > 0 {
		stackString += ":other:" + " : " + strings.Join(otherArray, " ") + "\n"
	}
	return stackString
}

func convertRecruitNumToEmojiString(project projectSchema) string {
	var recruitString string
	if project.FrontNum > 0 {
		recruitString += ":frontend:" + " " + strconv.Itoa(project.FrontNum) + "명\n"
	}
	if project.BackNum > 0 {
		recruitString += ":backend:" + " " + strconv.Itoa(project.BackNum) + "명\n"
	}
	if project.DataEngNum > 0 {
		recruitString += ":data_engineer:" + " " + strconv.Itoa(project.DataEngNum) + "명\n"
	}
	if project.DevOpsNum > 0 {
		recruitString += ":devops:" + " " + strconv.Itoa(project.DevOpsNum) + "명\n"
	}
	if project.uiUxNum > 0 {
		recruitString += ":figma:" + " " + strconv.Itoa(project.uiUxNum) + "명\n"
	}
	return recruitString
}
