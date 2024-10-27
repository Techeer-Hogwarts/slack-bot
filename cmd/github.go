package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type deployRequest struct {
	ImageName  string `json:"imageName"`
	ImageTag   string `json:"imageTag"`
	Message    string `json:"message"`
	CommitLink string `json:"commitLink"`
	Secret     string `json:"secret"`
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
		log.Println(err)
		return
	}
	if temp.Secret != secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauthorized")
		return
	}
	defer requestBody.Close()
	log.Println(temp.ImageName)
	log.Println(temp.ImageTag)
	err = sendDeploymentMessageToChannel(temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func sendDeploymentMessageToChannel(deployMessage deployRequest) error {
	api := slack.New(botToken)
	channelID := "C07H5TFEKBM"
	messageText := "이 메시지는 아래 커밋 메시지에 의해 트리거된 배포 파이프라인 입니다. \n 커밋 링크 & 메시지" + deployMessage.CommitLink + "\n" + deployMessage.Message + "아래 이미지를 배포할까요? \n 이미지 이름: " + deployMessage.ImageName + "\n 이미지 태그: " + deployMessage.ImageTag + "\n"
	deployButton := slack.NewButtonBlockElement("deploy_button", "apply", slack.NewTextBlockObject("plain_text", ":white_check_mark: 네", false, false))
	noDeployButton := slack.NewButtonBlockElement("no_deploy_button", "delete", slack.NewTextBlockObject("plain_text", ":no_entry_sign: 아니요", false, false))
	actionBlock := slack.NewActionBlock("deploy_action", deployButton, noDeployButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	_, _, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		return err
	}
	return nil
}
