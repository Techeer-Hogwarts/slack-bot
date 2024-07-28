package main

import (
	"log"
	"net/http"
	"os"

	"github.com/thomas-and-friends/slack-bot/cmd"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
	http.HandleFunc("/trigger_event", cmd.TriggerEvent)
	http.HandleFunc("/", cmd.SendHelloWorld)

	log.Printf("Server started on port: %v", port)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
