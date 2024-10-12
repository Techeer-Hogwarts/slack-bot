package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var Srv *sheets.Service

type YourStruct struct {
	Key1 string
	Key2 string
	Key3 string
}

func ConnectGoogle(spreadsheetID string) {
	credentials := map[string]interface{}{
		"type":                        GetEnv("TYPE", "service_account"),
		"project_id":                  GetEnv("PROJECT_ID", ""),
		"private_key_id":              GetEnv("PRIVATE_KEY_ID", ""),
		"private_key":                 GetEnv("PRIVATE_KEY", ""),
		"client_email":                GetEnv("CLIENT_EMAIL", ""),
		"client_id":                   GetEnv("CLIENT_ID", ""),
		"auth_uri":                    GetEnv("AUTH_URI", "https://accounts.google.com/o/oauth2/auth"),
		"token_uri":                   GetEnv("TOKEN_URI", "https://oauth2.googleapis.com/token"),
		"auth_provider_x509_cert_url": GetEnv("AUTH_PROVIDER_X509_CERT_URL", "https://www.googleapis.com/oauth2/v1/certs"),
		"client_x509_cert_url":        GetEnv("CLIENT_X509_CERT_URL", ""),
		"universe_domain":             GetEnv("UNIVERSE_DOMAIN", "googleapis.com"),
	}

	// Convert the credentials map to JSON
	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("Unable to marshal credentials: %v", err)
	}

	// Create the Sheets service using the credentials JSON
	ctx := context.Background()

	// Initialize Google Sheets API client
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Step 1: Create new tab with today's date and time
	now := time.Now()
	tabName := now.Format("2006/01/02 - 15:04")

	// Create new sheet (tab)
	addSheetRequest := &sheets.AddSheetRequest{
		Properties: &sheets.SheetProperties{
			Title: tabName,
		},
	}

	// Batch update request to add the new sheet
	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: addSheetRequest,
			},
		},
	}

	resp, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, request).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to create new sheet: %v", err)
	}

	newSheetID := resp.Replies[0].AddSheet.Properties.SheetId
	log.Printf("Created new sheet: %s\n", tabName)

	// Step 2: Mock data from your database
	data := []YourStruct{
		{"Value 1A", "Value 1B", "Value 1C"},
		{"Value 2A", "Value 2B", "Value 2C"},
	}

	// Step 3: Convert struct to a 2D array
	var values [][]interface{}

	headers := []interface{}{"Key1", "Key2", "Key3"}
	values = append(values, headers)

	for _, item := range data {
		row := []interface{}{item.Key1, item.Key2, item.Key3}
		values = append(values, row)
	}

	// Populate the new sheet with data
	writeRange := fmt.Sprintf("%s!A1", tabName)
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	// Step 4: Format header row to be bold
	headerFormatRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Range: &sheets.GridRange{
						SheetId:          newSheetID,
						StartRowIndex:    0, // First row (headers)
						EndRowIndex:      1, // Apply to just the header row
						StartColumnIndex: 0,
						EndColumnIndex:   3, // Adjust based on the number of columns
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true, // Bold formatting
							},
							BackgroundColor: &sheets.Color{
								Red:   0.9,
								Green: 0.9,
								Blue:  0.9, // Light gray background for headers
							},
						},
					},
					Fields: "userEnteredFormat(textFormat,backgroundColor)", // Specify which fields to format
				},
			},
		},
	}

	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, headerFormatRequest).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to format header row: %v", err)
	}

	log.Println("Header row formatted successfully.")
}
