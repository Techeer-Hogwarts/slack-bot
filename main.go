// package main

// import (
// 	"log"
// 	"net/http"

// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	"github.com/thomas-and-friends/slack-bot/cmd"
// 	"github.com/thomas-and-friends/slack-bot/config"
// 	"github.com/thomas-and-friends/slack-bot/db"
// )

// var err error

// func main() {
// 	port := config.GetEnv("PORT", "")
// 	http.HandleFunc("/slack/commands", cmd.HandleSlashCommand)
// 	// http.HandleFunc("/trigger_event", cmd.TriggerEvent)
// 	http.HandleFunc("/slack/interactions", cmd.HandleInteraction)
// 	http.HandleFunc("/", cmd.SendHelloWorld)
// 	// http.HandleFunc("/test", cmd.TestEvent)
// 	if port == "" {
// 		port = "8080"
// 	}

// 	db.DBMain, err = db.NewSQLDB("pgx")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db.ExecuteSQLFile("slack.sql")
// 	cmd.InitialDataUsers()
// 	cmd.InitialDataTags()
// 	defer func() {
// 		if err = db.DBMain.Close(); err != nil {
// 			panic(err)
// 		}
// 		log.Println("Disconnected from SQL Database")
// 	}()

//		log.Printf("Server started on port: %v", port)
//		log.Fatal(http.ListenAndServe(":"+port, nil))
//	}
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Specify the file path
	filePath := "ex.text"

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert content to a string
	text := string(content)

	// Replace all occurrences of 'oldChar' with 'newChar'
	oldChar := "\\" // Character to replace
	newChar := "\"" // Replacement character
	modifiedText := strings.ReplaceAll(text, oldChar, newChar)

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(modifiedText), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File updated successfully!")
}
