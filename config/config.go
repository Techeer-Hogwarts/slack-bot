package config

import (
	"github.com/thomas-and-friends/slack-bot/cmd" // Add this import statement
)

func PostgresDBConfig() string {
	cmd.LoadEnv()
	dbURL := cmd.GetEnv("DATABASE_URL", "")
	return dbURL
}
