package slackmessages

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/slack-go/slack"
)

func ConstructApplicantAndLeaderMessage(leaderProfile, applicantProfile *slack.User, userMessage models.UserMessageSchema) (slack.MsgOption, slack.MsgOption, error) {
	log.Printf("UserObject: %+v", userMessage)
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

func ConstructProjectMessage(project models.FindMemberSchema) string {
	log.Printf("UserObject: %+v", project)
	return "Some Message"
}

func ConstructStudyMessage(study models.FindMemberSchema) string {
	log.Printf("UserObject: %+v", study)
	return "Some Message"
}

func convertRecruitNumToEmojiString(project models.FindMemberSchema) string {
	var recruitString string
	if project.FrontNum > 0 {
		recruitString += ":frontend:" + " " + strconv.Itoa(project.FrontNum) + "명\n"
	}
	return recruitString
}

func convertStackToEmojiString(stack []string) string {
	var stackString string
	for _, s := range stack {
		stackString += s + " "
	}
	return stackString
}
