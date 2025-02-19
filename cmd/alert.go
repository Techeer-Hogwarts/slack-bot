package cmd

import (
	"encoding/json"
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
	ID             int    `json:"id"`
	TeamName       string `json:"teamName"`
	Type           string `json:"type"`
	LeaderEmail    string `json:"leaderEmail"`
	ApplicantEmail string `json:"applicantEmail"`
	Result         string `json:"result"`
}

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
		message := sendStudyMessage(study, api, channelID)
		log.Println("Message:", message)

	case "project":
		var project projectSchema
		err := mapToStruct(temp, &project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error mapping to projectSchema:", err)
			return
		}
		message := sendProjectMessage(project, api, channelID)
		log.Println("Message:", message)
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

func sendProjectMessage(project projectSchema, api *slack.Client, channelID string) string {
	profile, err := api.GetUserByEmail(project.Email)
	if err != nil {
		log.Printf("Failed to get user by email %s: %v", project.Email, err)
		return ""
	}
	userCode := profile.ID
	testMessaege := "[" + emoji_people + " *새로운 프로젝트 팀 공고가 올라왔습니다* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + project.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <<@" + userCode + ">>\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명입니다*\n" + project.ProjectExplain + "\n\n\n\n" +
		"> " + emoji_notebook + " *이런 사람을 원합니다!*\n" + project.RecruitExplain + "\n\n\n\n" +
		"> " + emoji_stack + " *사용되는 기술입니다*\n" + convertStackToEmojiString(project.Stack) + "\n\n\n" +
		"> " + emoji_dart + " *모집하는 직군 & 인원*\n" + convertRecruitNumToEmojiString(project) + "\n\n\n\n" +
		"> " + ":notion:" + "*노션 링크* \n" + project.NotionLink + "\n\n자세한 문의사항은" + "<@" + userCode + ">" + "에게 DM으로 문의 주세요!"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", testMessaege, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section)
	_, _, err = api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return ""
	}
	return ""
}

func sendStudyMessage(study studySchema, api *slack.Client, channelID string) string {
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "스터디 테스트입니다", false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section)
	_, _, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
		return ""
	}
	return ""
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
