package config

import (
	"os"

	"github.com/thomas-and-friends/slack-bot/cmd"
)

func PostgresDBConfig() string {
	cmd.LoadEnv()
	dbURL := os.Getenv("DATABASE_URL")
	return dbURL
}
