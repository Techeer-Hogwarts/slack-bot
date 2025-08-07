package server

import (
	"fmt"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/handlers"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/Techeer-Hogwarts/slack-bot/config"
	"github.com/slack-go/slack"
)

var (
	slackToken             string
	slackClient            *slack.Client
	githubURL              string
	githubToken            string
	cicdChannelID          string
	findMemberChannelID    string
	findMemberChannelIDDev string
)

func StartServer(port string) {
	// db := database.SetupDatabase()
	// defer db.Close()

	// err := database.MigrateSQLFile("cmd/database/migration/init.sql", db)
	// if err != nil {
	// 	panic(err)
	// }

	// repo := repositories.NewRepository(db)

	slackToken = config.GetEnvVarAsString("SLACK_BOT_TOKEN", "")
	githubToken = config.GetEnvVarAsString("GITHUB_ACTIONS_TOKEN", "")
	githubURL = config.GetEnvVarAsString("GITHUB_URL", "")
	cicdChannelID = config.GetEnvVarAsString("CICD_CHANNEL_ID", "")
	findMemberChannelID = config.GetEnvVarAsString("FIND_MEMBER_CHANNEL_ID", "")
	findMemberChannelIDDev = config.GetEnvVarAsString("FIND_MEMBER_CHANNEL_ID_DEV", "")

	slackClient = slack.New(slackToken)

	service := services.NewService(slackClient, githubURL, githubToken, cicdChannelID, findMemberChannelID, findMemberChannelIDDev)

	handler := handlers.NewHandler(service)

	router := setupRouter(handler)
	router.Run(fmt.Sprintf(":%s", port))
}
