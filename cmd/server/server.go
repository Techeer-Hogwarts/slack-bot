package server

import (
	"fmt"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/database"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/handlers"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/repositories"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/Techeer-Hogwarts/slack-bot/config"
	"github.com/slack-go/slack"
)

var (
	slackClient *slack.Client
	githubURL   string
	githubToken string
)

func StartServer(port string) {
	db := database.SetupDatabase()
	defer db.Close()

	err := database.MigrateSQLFile("cmd/database/migration/init.sql", db)
	if err != nil {
		panic(err)
	}

	repo := repositories.NewRepository(db)

	slackToken := config.GetEnvVarAsString("SLACK_BOT_TOKEN", "")
	githubToken = config.GetEnvVarAsString("GITHUB_ACTIONS_TOKEN", "")
	githubURL = config.GetEnvVarAsString("GITHUB_URL", "")

	slackClient = slack.New(slackToken)

	service := services.NewService(repo, slackClient, githubURL, githubToken)

	handler := handlers.NewHandler(service)

	router := setupRouter(handler)
	router.Run(fmt.Sprintf(":%s", port))
}
