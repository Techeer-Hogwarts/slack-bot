package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/slack-go/slack"
)

type DeployService interface {
	SendDeploymentMessage(req models.ImageDeployRequest) error
	TriggerDeployment(actionValue string, payload slack.InteractionCallback) error
	SendDeploymentStatus(status models.StatusRequest) error
}

type deployService struct {
	slackClient *slack.Client
	githubURL   string
	githubToken string
	channelID   string
}

func NewDeployService(slackClient *slack.Client, githubURL, githubToken, channelID string) DeployService {
	return &deployService{
		slackClient: slackClient,
		githubURL:   githubURL,
		githubToken: githubToken,
		channelID:   channelID,
	}
}

func (s *deployService) SendDeploymentMessage(req models.ImageDeployRequest) error {
	imageNameWithTag := req.ImageName + ":" + req.ImageTag + ":" + req.Environment
	messageText := fmt.Sprintf("새로운 이미지가 빌드 되었습니다.\n>커밋 링크 & 메시지: \n%s\n*아래 이미지를 배포할까요?*\n이미지 이름: `%s`\n이미지 태그: `%s`\n배포 환경: `%s`",
		req.CommitLink, req.ImageName, req.ImageTag, req.Environment)

	deployButton := slack.NewButtonBlockElement("deploy_button", imageNameWithTag,
		slack.NewTextBlockObject("plain_text", ":white_check_mark: 배포하기", false, false))
	noDeployButton := slack.NewButtonBlockElement("no_deploy_button", "delete",
		slack.NewTextBlockObject("plain_text", ":no_entry_sign: 삭제하기", false, false))

	inputElement := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "기본값: 1", false, false), "replica_count")
	inputElement.MaxLength = 2

	inputBlock := slack.NewInputBlock("replica_action",
		slack.NewTextBlockObject("plain_text", "복제 컨테이너 개수", false, false),
		slack.NewTextBlockObject("plain_text", "Scale Out 갯수", false, false),
		inputElement)

	actionBlock := slack.NewActionBlock("deploy_action", deployButton, noDeployButton)
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section, inputBlock, actionBlock)

	_, _, err := s.slackClient.PostMessage(s.channelID, messageBlocks)
	return err
}

func (s *deployService) SendDeploymentStatus(status models.StatusRequest) error {

	var messageText string
	if status.Status == "success" {
		messageText = fmt.Sprintf(":approved::approved: *새로운 이미지 배포를 성공하였습니다.* :approved::approved:\n이미지 이름: `%s`\n이미지 태그: `%s`\n링크: %s",
			status.ImageName, status.ImageTag, status.JobURL)
	} else {
		messageText = fmt.Sprintf(":exclamation::exclamation: *새로운 이미지 배포를 실패하였습니다.* :exclamation::exclamation:\n이미지 이름: `%s`\n이미지 태그: `%s`\n링크: %s\n실패한 단계: %s\n로그: ```%s```",
			status.ImageName, status.ImageTag, status.JobURL, status.FailedStep, status.Logs)
	}

	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", messageText, false, false), nil, nil)
	messageBlocks := slack.MsgOptionBlocks(section)

	_, _, err := s.slackClient.PostMessage(s.channelID, messageBlocks)
	return err
}

func (s *deployService) TriggerDeployment(actionValue string, payload slack.InteractionCallback) error {
	channelID := payload.Channel.ID

	replicaCount := payload.BlockActionState.Values["replica_action"]["replica_count"].Value
	if replicaCount == "" || replicaCount == "0" {
		replicaCount = "1"
	}

	imageNameAndTag := strings.Split(actionValue, ":")
	if len(imageNameAndTag) != 3 {
		return fmt.Errorf("invalid image name and tag format")
	}

	imageName := imageNameAndTag[0]
	imageTag := imageNameAndTag[1]
	environment := imageNameAndTag[2]
	messageText := fmt.Sprintf("*이미지 배포가 요청되었습니다.*\n이미지 이름: `%s`\n이미지 태그: `%s`\n복제 컨테이너 개수: `%s`\n배포 환경: `%s`\n요청 처리중......",
		imageName, imageTag, replicaCount, environment)

	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		return fmt.Errorf("failed to send deployment message: %w", err)
	}

	deployBody := models.ActionsRequestWrapper{
		Reference: "main",
		Inputs: models.DeployRequest{
			ImageName:   imageName,
			ImageTag:    imageTag,
			Replicas:    replicaCount,
			Environment: environment,
		},
	}
	log.Printf("deployBody: %+v", deployBody)
	return s.sendDeploymentRequest(deployBody)
}

func (s *deployService) sendDeploymentRequest(deployBody models.ActionsRequestWrapper) error {
	jsonBody, err := json.Marshal(deployBody)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment request: %w", err)
	}

	req, err := http.NewRequest("POST", s.githubURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+s.githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("deployment request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
