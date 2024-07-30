package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

var (
	signingKey string
	botToken   string
	channelID  string
	roleMap    map[string]string
	stackMap   map[string]string
)

const (
	emoji_people   = ":people_holding_hands::skin-tone-2-3:"
	emoji_golf     = ":golf:"
	emoji_star     = ":star2:"
	emoji_notebook = ":notebook:"
	emoji_stack    = ":hammer_and_pick:"
	emoji_dart     = ":dart:"
)

func init() {
	LoadEnv()
	signingKey = GetEnv("SLACK_SIGNING_SECRET", "")
	botToken = GetEnv("SLACK_BOT_TOKEN", "")
	channelID = GetEnv("CHANNEL_ID", "")
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

// func TriggerEvent(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Received a trigger event request")

// 	api := slack.New(botToken)
// 	messages, err := getChannelMessages(api, channelID)
// 	log.Printf("channelID: %v", channelID)
// 	if err != nil {
// 		log.Printf("Failed to retrieve messages: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(messages); err != nil {
// 		log.Printf("Failed to encode messages: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Println("Trigger event processed successfully")
// }

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

func getUsernameAndEmail(api *slack.Client, userID string) (string, error) {
	user, err := api.GetUserInfo(userID)
	if err != nil {
		return "", err
	}

	return user.RealName, nil
}

func constructMessageText(message FormMessage) (string, error) {
	if len(message.TeamRoles) == 0 || message.NumNewMembers == "" {
		return "", errors.New("TeamRoles is nil")
	}
	return "[" + emoji_people + message.TeamIntro + emoji_people + "]\n" +
		"> " + emoji_golf + "* 팀 이름 * \n " + message.TeamName + "\n\n\n\n" +
		"> " + emoji_star + "* 팀장 * <<@" + message.TeamLeader + ">>\n\n\n\n" +
		"> " + emoji_notebook + "* 팀/프로젝트 설명 *\n" + message.Description + "\n\n\n\n" +
		"> " + emoji_stack + "* 사용되는 기술 *\n" + formatListStacks(message.TechStacks) + "\n\n\n\n" +
		"> " + emoji_dart + "* 모집하는 직군 & 인원 *\n" + formatListRoles(message.TeamRoles, message.NumNewMembers) + "\n\n\n\n" +
		"> " + "* 그 외 추가적인 정보 * \n" + message.Etc + "자세한 문의사항은" + "<@" + message.TeamLeader + ">" + "에게 DM으로 문의 주세요!", nil
}

func formatListRoles(items []string, numPeople string) string {
	if len(items) == 0 {
		return "None"
	}
	var roles []string
	for _, role := range items {
		role_text := " • " + roleMap[role] + " (" + numPeople + "명)\n"
		roles = append(roles, role_text)
	}
	return strings.Join(roles, "")
}

func formatListStacks(items []string) string {
	if len(items) == 0 {
		return "None"
	}
	var stacks []string
	for _, stack := range items {
		stack_text := "`" + stackMap[stack] + "`, "
		stacks = append(stacks, stack_text)
	}
	return strings.Join(stacks, "")
}

// func formatListMembers(items []string) string {
// 	if len(items) == 0 {
// 		return "None"
// 	}
// 	return "<@" + strings.Join(items, "> <@")
// }

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

func openApplyModal(triggerID string) error {
	api := slack.New(botToken)
	modalRequest := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject("plain_text", "Apply to Team", false, false),
		Close:  slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"team_select",
					slack.NewTextBlockObject("plain_text", "Select a Team", false, false),
					slack.NewTextBlockObject("plain_text", "Select a team", false, false),
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeStatic,
						slack.NewTextBlockObject("plain_text", "Select a team", false, false),
						"selected_team",
						slack.NewOptionBlockObject("team1", slack.NewTextBlockObject("plain_text", "Team 1", false, false), nil),
						slack.NewOptionBlockObject("team2", slack.NewTextBlockObject("plain_text", "Team 2", false, false), nil),
					),
				),
				slack.NewInputBlock(
					"resume_input",
					slack.NewTextBlockObject("plain_text", "Upload Resume", false, false),
					slack.NewTextBlockObject("plain_text", "Paste your resume link", false, false),
					slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "Paste your resume link", false, false), "resume_link"),
				),
			},
		},
	}
	_, err := api.OpenView(triggerID, modalRequest)
	return err
}
