package cmd

import (
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func exportToGoogleSheet(w http.ResponseWriter, r *http.Request, api *slack.Client) {
	userCode := r.FormValue("user_id")
	if userCode == "U02AES3BH17" || userCode == "U033UTX061X" {
		log.Println("Exporting to Google Sheet")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

}
