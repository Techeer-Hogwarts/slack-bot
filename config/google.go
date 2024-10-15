package config

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var Srv *sheets.Service
var SheetsCTX context.Context

type YourStruct struct {
	Key1 string
	Key2 string
	Key3 string
}

func ConnectGoogle() {
	credentials := map[string]interface{}{
		"type":                        GetEnv("GTYPE", "service_account"),
		"project_id":                  GetEnv("GPROJECT_ID", ""),
		"private_key_id":              GetEnv("GPRIVATE_KEY_ID", ""),
		"private_key":                 GetEnv("GPRIVATE_KEY", ""),
		"client_email":                GetEnv("GCLIENT_EMAIL", ""),
		"client_id":                   GetEnv("GCLIENT_ID", ""),
		"auth_uri":                    GetEnv("GAUTH_URI", "https://accounts.google.com/o/oauth2/auth"),
		"token_uri":                   GetEnv("GTOKEN_URI", "https://oauth2.googleapis.com/token"),
		"auth_provider_x509_cert_url": GetEnv("GAUTH_PROVIDER_X509_CERT_URL", "https://www.googleapis.com/oauth2/v1/certs"),
		"client_x509_cert_url":        GetEnv("GCLIENT_X509_CERT_URL", ""),
		"universe_domain":             GetEnv("GUNIVERSE_DOMAIN", "googleapis.com"),
	}
	log.Println(credentials["private_key"])

	// Convert the credentials map to JSON
	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		log.Printf("Unable to marshal credentials: %v", err)
	}

	// Create the Sheets service using the credentials JSON
	SheetsCTX = context.Background()

	// Initialize Google Sheets API client
	Srv, err = sheets.NewService(SheetsCTX, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	log.Println("Connected to Google Sheets API")
}
