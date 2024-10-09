package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/thomas-and-friends/slack-bot/db"
)

type zipRequest struct {
	Email  string `json:"email"`
	Secret string `json:"secret"`
}

func ZipPictureHandler(w http.ResponseWriter, r *http.Request) {
	api := slack.New(botToken)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req zipRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(req.Email)
	if req.Secret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	profile, err := api.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userImage := profile.Profile.ImageOriginal
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userImage)
}

func ZipVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req zipRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(req.Email)
	if req.Secret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, verify, err := db.CheckUserWithEmail(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if verify {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}
