package main

import (
	"github.com/Techeer-Hogwarts/slack-bot/cmd/server"
	"github.com/Techeer-Hogwarts/slack-bot/config"
)

// var err error

// @securityDefinitions.apikey SlackSigningSecret
// @in header
// @name X-Slack-Signature

// @securityDefinitions.apikey APIKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey JwtAuth
// @in cookie
// @name access_token
func main() {
	// reload := config.GetEnv("RELOAD", "true")
	// port := config.GetEnv("PORT", "")
	// http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
	// http.HandleFunc("/slack/interactions", cmd.HandleInteraction)
	// http.HandleFunc("/", cmd.SendHelloWorld)
	// http.HandleFunc("/api/v1/profile/picture", cmd.ZipPictureHandler)
	// http.HandleFunc("/api/v1/deploy/image", cmd.DeployImageHandler)
	// http.HandleFunc("/api/v1/deploy/status", cmd.DeployStatusHandler)
	// // 알림 기능
	// http.HandleFunc("/api/v1/alert/user", cmd.AlertUserHandler)
	// http.HandleFunc("/api/v1/alert/channel", cmd.AlertChannelHandler)
	// if port == "" {
	// 	port = "8080"
	// }

	// db.DBMain, err = db.NewSQLDB("pgx")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if reload == "true" {
	// 	cmd.InitialDataUsers()
	// 	cmd.InitialDataTags()
	// }
	// db.ExecuteSQLFile("slack.sql")
	// config.ConnectGoogle()
	// defer func() {
	// 	if err = db.DBMain.Close(); err != nil {
	// 		panic(err)
	// 	}
	// 	log.Println("Disconnected from SQL Database")
	// }()

	// log.Printf("Server started on port: %v", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	config.LoadEnvFile(".env")

	port := config.GetEnvVarAsString("PORT", "8080")
	server.StartServer(port)
}
