package cmd

import (
	"encoding/json"
	"log"
	"net/http"
)

func DeployImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deploy Image Handler")
	var temp interface{}
	requestBody := r.Body
	err := json.NewDecoder(requestBody).Decode(&temp)
	if err != nil {
		panic(err)
	}
	defer requestBody.Close()
	log.Println(temp)
}
