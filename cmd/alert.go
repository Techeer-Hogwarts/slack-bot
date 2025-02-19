package cmd

import (
	"encoding/json"
	"log"
	"net/http"
)

type alertRequest struct {
	Secret string
	Alert  string
}

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
	var temp interface{}
	requestBody := r.Body

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
