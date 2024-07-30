package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thomas-and-friends/slack-bot/cmd"
)

var err error

func main() {
	// port := os.Getenv("PORT")
	http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
	http.HandleFunc("/trigger_event", cmd.TriggerEvent)
	http.HandleFunc("/slack/interactions", cmd.HandleInteraction)
	http.HandleFunc("/", cmd.SendHelloWorld)
	http.HandleFunc("/test", cmd.TestEvent)
	port := "8080"

	// db.DBMain, err = db.NewSQLDB("pgx")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() {
	// 	if err = db.DBMain.Close(); err != nil {
	// 		panic(err)
	// 	}
	// 	log.Println("Disconnected from SQL Database")
	// }()

	log.Printf("Server started on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
