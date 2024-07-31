package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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
	stackMap   map[string]string
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
	channelID = config.GetEnv("CHANNEL_ID", "")
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
	stackMap = map[string]string{
		"none":          "없음",
		"react":         "React.js",
		"vue":           "Vue.js",
		"next":          "Next.js",
		"svelte":        "SvelteKit",
		"angular":       "Angular",
		"django":        "Django",
		"flask":         "Flask",
		"rails":         "Ruby on Rails",
		"spring":        "Spring Boot",
		"express":       "Express.js",
		"laravel":       "Laravel",
		"s3":            "S3/Cloud Storage",
		"go":            "Go Lang",
		"ai":            "AI/ML (Tensorflow, PyTorch)",
		"kube":          "Kubernetes",
		"jenkins":       "Jenkins CI",
		"actions":       "Github Actions",
		"spin":          "Spinnaker",
		"graphite":      "Graphite",
		"kafka":         "Kafka",
		"docker":        "Docker",
		"ansible":       "Ansible",
		"terraform":     "Terraform",
		"fastapi":       "FastAPI",
		"redis":         "Redis",
		"msa":           "MSA",
		"java":          "Java",
		"python":        "Python",
		"jsts":          "JavaScript/TypeScript",
		"cpp":           "C/C++",
		"csharp":        "C#",
		"ruby":          "Ruby",
		"aws":           "AWS",
		"gcp":           "GCP",
		"ELK":           "ELK Stack",
		"elasticsearch": "Elasticsearch",
		"prom":          "Prometheus",
		"grafana":       "Grafana",
		"celery":        "Celery",
		"nginx":         "Nginx",
		"cdn":           "CDN (CloudFront)",
		"nestjs":        "Nest.JS",
		"zustand":       "Zustand",
		"tailwind":      "Tailwind CSS",
		"bootstrap":     "Bootstrap",
		"postgre":       "PostgreSQL",
		"mysql":         "MySQL",
		"mongo":         "MongoDB",
		"node":          "Node.js",
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

func getChannelMessages(api *slack.Client, channelID string) (*slack.GetConversationHistoryResponse, error) {
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
	}

	history, err := api.GetConversationHistory(&historyParams)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func TriggerEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a trigger event request")

	api := slack.New(botToken)
	history, err := getChannelMessages(api, channelID)
	if err != nil {
		log.Printf("Failed to retrieve messages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("channelID: %v", channelID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Trigger event processed successfully")
}

// func getUsernameAndEmail(api *slack.Client, userID string) (string, error) {
// 	user, err := api.GetUserInfo(userID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return user.RealName, nil
// }

func constructMessageText(message FormMessage) (string, error) {
	if len(message.Members) == 0 || message.Description == "" || message.TeamName == "" || message.TechStacks == nil {
		return "", errors.New("missing required fields except team leader")
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
	if uxNum == 0 {
		log.Println("No UX/UI members")
	} else {
		roles = append(roles, " • "+roleMap["uxui"]+" ("+message.UxMembers+"명)\n")
	}
	frontNum, err := strconv.Atoi(message.FrontMembers)
	if err != nil {
		frontNum = 0
	}
	if frontNum == 0 {
		log.Println("No Frontend members")
	} else {
		roles = append(roles, " • "+roleMap["frontend"]+" ("+message.FrontMembers+"명)\n")
	}
	backNum, err := strconv.Atoi(message.BackMembers)
	if err != nil {
		backNum = 0
	}
	if backNum == 0 {
		log.Println("No Backend members")
	} else {
		roles = append(roles, " • "+roleMap["backend"]+" ("+message.BackMembers+"명)\n")
	}
	dataNum, err := strconv.Atoi(message.DataMembers)
	if err != nil {
		dataNum = 0
	}
	if dataNum == 0 {
		log.Println("No Data members")
	} else {
		roles = append(roles, " • "+roleMap["data"]+" ("+message.DataMembers+"명)\n")
	}
	opsNum, err := strconv.Atoi(message.OpsMembers)
	if err != nil {
		opsNum = 0
	}
	if opsNum == 0 {
		log.Println("No OPS/SRE members")
	} else {
		roles = append(roles, " • "+roleMap["devops"]+" ("+message.OpsMembers+"명)\n")
	}
	studyNum, err := strconv.Atoi(message.StudyMembers)
	if err != nil {
		studyNum = 0
	}
	if studyNum == 0 {
		log.Println("No Study members")
	} else {
		roles = append(roles, " • "+roleMap["study"]+" ("+message.StudyMembers+"명)\n")
	}
	etcNum, err := strconv.Atoi(message.EtcMembers)
	if err != nil {
		etcNum = 0
	}
	if etcNum == 0 {
		log.Println("No Etc members")
	} else {
		roles = append(roles, " • "+roleMap["etc"]+" ("+message.EtcMembers+"명)\n")
	}
	return strings.Join(roles, "")
}

func formatListStacks(items []string) string {
	if len(items) == 0 {
		return "None"
	}
	var stacks []string
	for _, stack := range items {
		stack_text := "`" + stackMap[stack] + "`"
		stacks = append(stacks, stack_text)
	}
	return strings.Join(stacks, ", ")
}

func InitialDataUsers() {
	api := slack.New(botToken)
	err := addAllUsers(api)
	if err != nil {
		log.Printf("Failed to add all users: %v", err)
	}
}

func InitialDataTags() {
	api := slack.New(botToken)
	err := addAllTags(api)
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
		if username == "" {
			username = user.Profile.RealNameNormalized
		}
		if !user.IsBot && !user.Deleted && !user.IsAppUser && !user.IsOwner && user.ID != "USLACKBOT" {
			ms, _, err := db.GetUser(user.ID)
			if ms == "na" {
				err = db.AddUser(user.ID, username)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return fmt.Errorf("failed to get user: %s", err.Error())
			}
		}
	}
	if err != nil {
		return err
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
	if err != nil {
		return err
	}

	return nil
}

func addAllTags(api *slack.Client) error {
	for key, value := range stackMap {
		ms, err := db.GetTag(key)
		if ms == "na" {
			err := db.AddTag(key, value)
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err != nil {
				return fmt.Errorf("failed to get tag: %s", err.Error())
			}
		}
	}
	return nil
}
func deleteMessage(payload slack.InteractionCallback) error {
	api := slack.New(botToken)
	actionUserID := payload.User.ID
	actionMessageTimestamp := payload.Message.Timestamp
	actionContainerTimestamp := payload.Container.MessageTs
	log.Printf("User ID: %v | Message Timestamp: %v | Container Timestamp: %v", actionUserID, actionMessageTimestamp, actionContainerTimestamp)
	teamObj, _ := db.GetTeam(actionMessageTimestamp)
	// if err != nil {
	// 	return err
	// }
	teamTs := teamObj.TeamTs
	log.Println("Team Timestamp: ", teamTs)
	err := sendSuccessMessage(api, payload.Channel.ID, payload.User.ID, "Message deleted successfully")
	if err != nil {
		return err
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
	err := db.AddTeam(teamObj)
	if err != nil {
		return err
	}
	return nil
}

func openApplyModal(triggerID string) error {
	api := slack.New(botToken)
	modalRequest := slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: "apply_form",
		Title:      slack.NewTextBlockObject("plain_text", "Apply to Team", false, false),
		Close:      slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"team_select",
					slack.NewTextBlockObject("plain_text", "Select a Team", false, false),
					slack.NewTextBlockObject("plain_text", "Select a team", false, false),
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeStatic,
						slack.NewTextBlockObject("plain_text", "내부", false, false),
						"selected_team",
						slack.NewOptionBlockObject("내부 키", slack.NewTextBlockObject("plain_text", "Team 1", false, false), nil),
						slack.NewOptionBlockObject("team2", slack.NewTextBlockObject("plain_text", "Team 2", false, false), nil),
					),
				),
				slack.NewInputBlock(
					"resume_input",
					slack.NewTextBlockObject("plain_text", "Upload Resume", false, false),
					slack.NewTextBlockObject("plain_text", "Paste your resume link", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "내부", false, false), "resume_link"),
				),
			},
		},
	}
	_, err := api.OpenView(triggerID, modalRequest)
	return err
}

func sendSuccessMessage(api *slack.Client, channelID string, userID string, messageText string) error {
	_, err := api.PostEphemeral(channelID, userID, slack.MsgOptionText(messageText, false))
	return err
}

func sendFailMessage(api *slack.Client, channelID string, userID string, messageText string) error {
	log.Println(userID)
	_, err := api.PostEphemeral(channelID, userID, slack.MsgOptionText(messageText, false))
	return err
}
