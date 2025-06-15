package slackmessages

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/slack-go/slack"
)

func ConstructApplicantAndLeaderMessage(leaderProfile, applicantProfile *slack.User, userMessage models.UserMessageSchema) (slack.MsgOption, slack.MsgOption, error) {
	var leaderStatus string
	var applicantStatus string
	switch userMessage.Result {
	case "PENDING":
		leaderStatus = "지원자가 있습니다."
		applicantStatus = "지원이 완료됐습니다."
	case "CANCELLED":
		leaderStatus = "지원자께서 지원을 취소 하셨습니다."
		applicantStatus = "지원자께서 지원을 취소 하셨습니다."
	case "REJECT":
		leaderStatus = "팀원중 한명이 지원자를 거절 하셨습니다."
		applicantStatus = "지원이 거절되었습니다. 다음 기회에 함께할 수 있길 바랍니다! :pray:"
	case "APPROVED":
		leaderStatus = "팀원중 한명이 지원자를 수락 하셨습니다."
		applicantStatus = "지원이 승인되어 팀에 합류하셨습니다! 함께 열심히 해봐요! :rocket:"
	default:
		log.Println("Invalid result:", userMessage.Result)
		return nil, nil, fmt.Errorf("invalid result: %s", userMessage.Result)
	}

	leaderMsg := "[" + emoji_people + " *지원 결과 알림* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + userMessage.TeamName + "\n\n\n\n" +
		"> " + emoji_star + " *지원자* <@" + applicantProfile.ID + ">\n\n\n\n" +
		"> " + emoji_notebook + " *상태:* " + leaderStatus + "\n\n\n\n" +
		"> " + emoji_dart + " *링크* \n" + fmt.Sprintf(redirectURL, userMessage.Type, userMessage.TeamID) + "\n\n\n\n"
	leaderSection := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", leaderMsg, false, false), nil, nil)
	leaderMessageBlock := slack.MsgOptionBlocks(leaderSection)

	applicantMsg := "[" + emoji_people + " *지원 결과 알림* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + userMessage.TeamName + "\n\n\n\n" +
		"> " + emoji_star + " *팀장:* <@" + leaderProfile.ID + ">\n\n\n\n" +
		"> " + emoji_notebook + " *지원 결과:* " + applicantStatus + "\n\n\n\n" +
		"> " + emoji_dart + " *링크* \n" + fmt.Sprintf(redirectURL, userMessage.Type, userMessage.TeamID) + "\n\n\n\n"
	applicantsSction := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", applicantMsg, false, false), nil, nil)
	applicantMessageBlock := slack.MsgOptionBlocks(applicantsSction)
	return leaderMessageBlock, applicantMessageBlock, nil
}

func ConstructProjectMessage(project models.FindMemberSchema, profileIDs []string) (slack.MsgOption, error) {
	teamLeaderString := strings.Join(profileIDs, ">, <@")
	teamLeaderString = "<@" + teamLeaderString + ">"
	projectMessage := "[" + emoji_people + " *새로운 프로젝트 팀 공고가 올라왔습니다* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + project.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* " + teamLeaderString + "\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명입니다*\n" + project.ProjectExplain + "\n\n\n\n" +
		"> " + ":woman-raising-hand:" + " *이런 사람을 원합니다!*\n" + project.RecruitExplain + "\n\n\n\n" +
		"> " + emoji_stack + " *사용되는 기술입니다*\n" + convertStackToEmojiString(project.Stack) + "\n\n\n" +
		"> " + emoji_dart + " *모집하는 직군 & 인원*\n" + convertRecruitNumToEmojiString(project) + "\n\n\n\n" +
		"> " + ":notion:" + " *노션 링크* \n" + project.NotionLink + "\n\n자세한 문의사항은 " + teamLeaderString + " 에게 DM으로 문의 주세요!"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", projectMessage, false, false), nil, nil)
	applyButton := slack.NewButtonBlockElement("", "apply", slack.NewTextBlockObject("plain_text", ":link: 팀 지원하러 가기!", false, false))
	applyButton.URL = fmt.Sprintf(redirectURL, project.Type, project.ID)
	profileIDsJSON, _ := json.Marshal(profileIDs)
	deleteButton := slack.NewButtonBlockElement("delete_button", string(profileIDsJSON), slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	return messageBlocks, nil
}

func ConstructStudyMessage(study models.FindMemberSchema, profileIDs []string) (slack.MsgOption, error) {
	teamLeaderString := strings.Join(profileIDs, ">, <@")
	teamLeaderString = "<@" + teamLeaderString + ">"
	studyMessage := "[" + emoji_people + " *새로운 스터디 팀 공고가 올라왔습니다* " + emoji_people + "]\n" +
		"> " + ":name_badge:" + " *팀 이름* \n " + study.Name + "\n\n\n\n" +
		"> " + emoji_star + " *팀장* " + teamLeaderString + "\n\n\n\n" +
		"> " + emoji_notebook + " *팀/프로젝트 설명입니다*\n" + study.StudyExplain + "\n\n\n\n" +
		"> " + ":man-raising-hand:" + " *이런 사람을 원합니다!*\n" + study.RecruitExplain + "\n\n\n\n" +
		"> " + ":pencil:" + " *지켜야 하는 규칙입니다!*\n" + study.Rule + "\n\n\n" +
		"> " + emoji_dart + " *모집하는 스터디 인원*\n" + strconv.Itoa(study.RecruitNum) + "명\n\n\n\n" +
		"> " + ":notion:" + " *노션 링크* \n" + study.NotionLink + "\n\n자세한 문의사항은 " + teamLeaderString + " 에게 DM으로 문의 주세요!"
	section := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", studyMessage, false, false), nil, nil)
	applyButton := slack.NewButtonBlockElement("", "apply", slack.NewTextBlockObject("plain_text", ":link: 팀 지원하러 가기!", false, false))
	applyButton.URL = fmt.Sprintf(redirectURL, study.Type, study.ID)
	profileIDsJSON, _ := json.Marshal(profileIDs)
	deleteButton := slack.NewButtonBlockElement("delete_button", string(profileIDsJSON), slack.NewTextBlockObject("plain_text", ":warning: 삭제하기!", false, false))
	actionBlock := slack.NewActionBlock("apply_action", applyButton, deleteButton)
	messageBlocks := slack.MsgOptionBlocks(section, actionBlock)
	return messageBlocks, nil
}

func convertRecruitNumToEmojiString(project models.FindMemberSchema) string {
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
	if project.FullStack > 0 {
		recruitString += ":fullstack:" + " " + strconv.Itoa(project.FullStack) + "명\n"
	}
	return recruitString
}

func convertStackToEmojiString(stack []string) string {
	var backArray []string
	var frontArray []string
	var devOpsArray []string
	var otherArray []string
	var databaseArray []string
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
		case "DATABASE":
			databaseArray = append(databaseArray, emoji)
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
	if len(databaseArray) > 0 {
		stackString += ":database:" + " : " + strings.Join(databaseArray, " ") + "\n"
	}
	return stackString
}
