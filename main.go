package main

import (
	"log"
	"net/http"

	"github.com/thomas-and-friends/slack-bot/cmd"
)

func main() {
	http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
	http.HandleFunc("/trigger_event", cmd.TriggerEvent)
	http.HandleFunc("/", cmd.SendHelloWorld)

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
