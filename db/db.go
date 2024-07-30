package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/thomas-and-friends/slack-bot/config"
)

var DBMain *DB

type DB struct {
	*sql.DB
}

func connectSQLDB(dbDriver string) (*sql.DB, error) {
	urlDB := config.PostgresDBConfig()
	db, err := sql.Open(dbDriver, urlDB)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Database Connection Success!")
	return db, nil
}

func NewSQLDB(dbDriver string) (*DB, error) {
	db, err := connectSQLDB(dbDriver)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to SQL Database")
	return &DB{db}, nil
}
