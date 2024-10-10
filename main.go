package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thomas-and-friends/slack-bot/cmd"
	"github.com/thomas-and-friends/slack-bot/config"
	"github.com/thomas-and-friends/slack-bot/db"
)

var err error

func main() {
	reload := config.GetEnv("RELOAD", "true")
	port := config.GetEnv("PORT", "")
	http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
	http.HandleFunc("/slack/interactions", cmd.HandleInteraction)
	http.HandleFunc("/", cmd.SendHelloWorld)
	http.HandleFunc("/api/v1/profile/picture", cmd.ZipPictureHandler)
	http.HandleFunc("/api/v1/profile/verify", cmd.ZipVerifyHandler)
	if port == "" {
		port = "8080"
	}

	db.DBMain, err = db.NewSQLDB("pgx")
	if err != nil {
		log.Fatal(err)
	}
	if reload == "true" {
		cmd.InitialDataUsers()
		cmd.InitialDataTags()
	}
	db.ExecuteSQLFile("slack.sql")
	// config.ConnectGoogle()
	defer func() {
		if err = db.DBMain.Close(); err != nil {
			panic(err)
		}
		log.Println("Disconnected from SQL Database")
	}()

	log.Printf("Server started on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
