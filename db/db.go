package db

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/thomas-and-friends/slack-bot/config"
// )

// var DBMain *DB

// type DB struct {
// 	*sql.DB
// }

// func connectSQLDB(dbDriver string) (*sql.DB, error) {
// 	urlDB := config.PostgresDBConfig()
// 	db, err := sql.Open(dbDriver, urlDB)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database: %v\n", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println("Database Connection Success!")
// 	return db, nil
// }

// func NewSQLDB(dbDriver string) (*DB, error) {
// 	db, err := connectSQLDB(dbDriver)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Println("Connected to SQL Database")
// 	return &DB{db}, nil
// }

// func ExecuteSQLFile(filePath string) error {
// 	sqlContent, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read SQL file: %v", err)
// 	}
// 	commands := strings.Split(string(sqlContent), ";")

// 	for _, command := range commands {
// 		command = strings.TrimSpace(command)
// 		if command == "" {
// 			continue
// 		}

// 		_, err := DBMain.Exec(command)
// 		if err != nil {
// 			return fmt.Errorf("failed to execute SQL command: %v", err)
// 		}
// 	}

// 	return nil
// }
