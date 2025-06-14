package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Techeer-Hogwarts/slack-bot/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupDatabase() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.GetEnvVarAsString("DB_HOST", "localhost"),
		config.GetEnvVarAsString("DB_USER", "test"),
		config.GetEnvVarAsString("DB_PASSWORD", "test"),
		config.GetEnvVarAsString("DB_NAME", "bot"),
		config.GetEnvVarAsString("DB_PORT", "5432"),
		config.GetEnvVarAsString("DB_SSLMODE", "disable"),
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open connection: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	log.Println("Database connected")
	return db
}

func MigrateSQLFile(filePath string, db *sql.DB) error {
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}
	commands := strings.Split(string(sqlContent), ";")

	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		_, err := db.Exec(command)
		if err != nil {
			return fmt.Errorf("failed to execute SQL command: %v", err)
		}
	}
	log.Println("Database migrated")
	return nil
}
