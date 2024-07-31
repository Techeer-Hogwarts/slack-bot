package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("unable to load .env file")
	}
}

func GetEnv(key string, defaultVal string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultVal
	}
	log.Printf("env %s: found \n", key)
	return value
}
