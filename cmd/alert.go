package cmd

import (
	"encoding/json"
	"log"
	"net/http"

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
	userCode := profile.ID
	testMessaege := "[" + emoji_people + project.ProjectExplain + emoji_people + "]\n" +
		"> " + emoji_golf + " *팀 이름* \n " + project.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* <<@" + userCode + ">>\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명*\n" + project.RecruitExplain + "\n\n\n\n" +
		"> " + emoji_stack + " *사용되는 기술*\n" + "테스트" + "\n\n\n\n" +
		"> " + emoji_dart + " *모집하는 직군 & 인원*\n" + "직군" + "\n\n\n\n" +
		"> " + "*그 외 추가적인 정보* \n" + project.NotionLink + "\n\n자세한 문의사항은" + "<@" + project.Leader + ">" + "에게 DM으로 문의 주세요!"
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
