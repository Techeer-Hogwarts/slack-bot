package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"github.com/thomas-and-friends/slack-bot/config"
	"github.com/thomas-and-friends/slack-bot/db"
	"google.golang.org/api/sheets/v4"
)

type TeamWithUsers struct {
	Team  db.Team
	Users []db.UserObj
}

func writeGoogleSheet() string {
	now := time.Now()
	tabName := now.Format("2006/01/02 - 15:04")

	// Create new sheet (tab)
	addSheetRequest := &sheets.AddSheetRequest{
		Properties: &sheets.SheetProperties{
			Title: tabName,
		},
	}
	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: addSheetRequest,
			},
		},
	}

	resp, err := config.Srv.Spreadsheets.BatchUpdate(sheetsID, request).Context(config.SheetsCTX).Do()
	if err != nil {
		log.Fatalf("Unable to create new sheet: %v", err)
	}

	newSheetID := resp.Replies[0].AddSheet.Properties.SheetId
	log.Printf("Created new sheet: %s with id %v\n", tabName, newSheetID)

	allTeams, err := db.GetAllActiveTeams()
	if err != nil {
		log.Printf("Failed to get all teams: %v", err)
	}
	var data []TeamWithUsers

	for _, team := range allTeams {
		team_id := team.TeamID
		teamIDInt, _ := strconv.Atoi(team_id)
		allUsersInTeam, err := db.GetAllUsersInTeam(teamIDInt)
		if err != nil {
			log.Printf("Failed to get all users in team: %v", err)
		}
		data = append(data, TeamWithUsers{Team: team, Users: allUsersInTeam})
	}

	var values [][]interface{}

	headers := []interface{}{"팀 이름", "팀 종류", "팀 리더", "팀 소개", "팀 설명", "팀 멤버들", "기타"}
	values = append(values, headers)

	for _, team := range data {
		var users string
		for _, user := range team.Users {
			users += fmt.Sprintf("%s (%s)\n", user.UserName, user.UserEmail)
		}
		team.Team.TeamLeaderName, _, _ = db.GetUser(team.Team.TeamLeader)
		row := []interface{}{team.Team.TeamName, team.Team.TeamType, team.Team.TeamLeaderName, team.Team.TeamIntro, team.Team.TeamDesc, users, team.Team.TeamEtc}
		values = append(values, row)
	}

	// Populate the new sheet with data
	writeRange := fmt.Sprintf("%s!A1", tabName)
	_, err = config.Srv.Spreadsheets.Values.Update(sheetsID, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Unable to write data to sheet: %v", err)
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
						EndColumnIndex:   7, // Adjust based on the number of columns
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true, // Bold formatting
							},
							BackgroundColor: &sheets.Color{
								Red:   0.9,
								Green: 0.9,
								Blue:  0.9,
							},
						},
					},
					Fields: "userEnteredFormat(textFormat,backgroundColor)", // Specify which fields to format
				},
			},
		},
	}

	_, err = config.Srv.Spreadsheets.BatchUpdate(sheetsID, headerFormatRequest).Context(config.SheetsCTX).Do()
	if err != nil {
		log.Printf("Unable to format header row: %v", err)
	}

	log.Println("Header row formatted successfully.")
	resizeRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						SheetId:    newSheetID,
						Dimension:  "COLUMNS",
						StartIndex: 3, // First column (A)
						EndIndex:   7, // Up to third column (C)
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 300, // Adjust column width (in pixels)
					},
					Fields: "pixelSize",
				},
			},
		},
	}
	_, err = config.Srv.Spreadsheets.BatchUpdate(sheetsID, resizeRequest).Context(config.SheetsCTX).Do()
	if err != nil {
		log.Printf("Unable to resize columns: %v", err)
	}

	spreadsheetLink := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit#gid=%d", sheetsID, newSheetID)
	return spreadsheetLink
}

func ExportToGoogleSheet(w http.ResponseWriter, r *http.Request) {
	api := slack.New(botToken)
	userCode := r.FormValue("user_id")
	log.Println(userCode)
	if userCode == "U02AES3BH17" || userCode == "U033UTX061X" {
		log.Println("Exporting to Google Sheet")
		link := writeGoogleSheet()
		log.Printf("Exported to Google Sheet: %s", link)
		message := fmt.Sprintf("구글 시트 링크: %s", link)
		sendDMSuccessMessage(api, userCode, message)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Exported to Google Sheet"))
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

}
