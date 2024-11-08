package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

type deployRequest struct {
	ImageName  string `json:"imageName"`
	ImageTag   string `json:"imageTag"`
	CommitLink string `json:"commitLink"`
	Secret     string `json:"secret"`
}

type actionsRequestWrapper struct {
	Reference string         `json:"ref"`
	Inputs    actionsRequest `json:"inputs"`
}

type actionsRequest struct {
	ImageName     string `json:"imageName"`
	ImageTag      string `json:"imageTag"`
	ReplicaCouint string `json:"replicaCount"`
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
	imageNameWithTag := deployMessage.ImageName + ":" + deployMessage.ImageTag
	messageText := "이 메시지는 아래 커밋 메시지에 의해 트리거된 배포 파이프라인 입니다. \n 커밋 링크 & 메시지 \n" + deployMessage.CommitLink + "\n" + "아래 이미지를 배포할까요? \n 이미지 이름: " + deployMessage.ImageName + "\n 이미지 태그: " + deployMessage.ImageTag + "\n"
	deployButton := slack.NewButtonBlockElement("deploy_button", imageNameWithTag, slack.NewTextBlockObject("plain_text", ":white_check_mark: 네", false, false))
	noDeployButton := slack.NewButtonBlockElement("no_deploy_button", "delete", slack.NewTextBlockObject("plain_text", ":no_entry_sign: 아니요", false, false))
	inputElement := slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "기본값: 1", false, false), "replica_count")
	inputElement.MaxLength = 2
	inputBlock := slack.NewInputBlock("replica_action", slack.NewTextBlockObject("plain_text", "복제 컨테이너 개수", false, false), slack.NewTextBlockObject("plain_text", "Scale Out 갯수", false, false), inputElement)
	actionBlock := slack.NewActionBlock("deploy_action", deployButton, noDeployButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, inputBlock, actionBlock)
	_, _, err := api.PostMessage(channelID, messageBlocks)
	if err != nil {
		return err
	}
	return nil
}

func triggerDeployment(actionValue string, payload slack.InteractionCallback) error {
	api := slack.New(botToken)
	channelID := payload.Channel.ID
	imageNameWithTag := actionValue
	replicaCount := payload.BlockActionState.Values["replica_action"]["replica_count"].Value
	imageNameAndTag := strings.Split(imageNameWithTag, ":")
	imageName := imageNameAndTag[0]
	imageTag := imageNameAndTag[1]
	messageText := fmt.Sprintf("이미지 배포가 요청되었습니다.\n이미지 이름: %s\n이미지 태그: %s\n 복제 컨테이너 개수: %s", imageName, imageTag, replicaCount)
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		log.Println(err)
		return err
	}
	imageDeploy := actionsRequest{
		ImageName:     imageName,
		ImageTag:      imageTag,
		ReplicaCouint: replicaCount,
	}

	finalBody := actionsRequestWrapper{
		Reference: "main",
		Inputs:    imageDeploy,
	}

	err = sendDeploymentRequest(finalBody)
	return nil
}

func sendDeploymentRequest(deployBody actionsRequestWrapper) error {
	log.Println("Send Deployment Request")
	jsonBody, err := json.Marshal(deployBody)
	if err != nil {
		log.Println(err)
		return err
	}
	url := githubURL
	token := githubToken
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	log.Printf("Request URL: %s", req.URL)
	log.Printf("Request Header: %s", req.Header)
	log.Printf("Request Body: %s", req.Body)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}
	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Response Body: %s", body)
	if resp.StatusCode != http.StatusNoContent {
		log.Println("Failed to send deployment request")
		return fmt.Errorf("Failed to send deployment request")
	}
	return nil
}
