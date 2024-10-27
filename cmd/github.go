package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type deployRequest struct {
	ImageName string `json:"imageName"`
	ImageTag  string `json:"imageTag"`
	Secret    string `json:"secret"`
}

func DeployImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("Deploy Image Handler")
	var temp deployRequest
	requestBody := r.Body
	err := json.NewDecoder(requestBody).Decode(&temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(temp)
	if temp.Secret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	defer requestBody.Close()
	log.Println(temp.ImageName)
	log.Println(temp.ImageTag)
	err = sendDeploymentMessageToChannel(temp.ImageName, temp.ImageTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func sendDeploymentMessageToChannel(imageName string, imageTag string) error {
	api := slack.New(botToken)
	channelID := "C07H5TFEKBM"
	messageText := "아래 이미지를 배포할까요? \n 이미지 이름: " + imageName + "\n 이미지 태그: " + imageTag + "\n"
	deployButton := slack.NewButtonBlockElement("deploy_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 네", false, false))
	noDeployButton := slack.NewButtonBlockElement("no_deploy_button", "delete", slack.NewTextBlockObject("plain_text", ":no_entry_sign: 아니요", false, false))
	actionBlock := slack.NewActionBlock("apply_action", deployButton, noDeployButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	_, _, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		return err
	}
	return nil
}
