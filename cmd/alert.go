package cmd

import (
	"net/http"
)

type alertRequest struct {
	Secret string
	Alert  string
}

func AlertChannelHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// log.Println("Alert Channel Handler")
	// var temp alertRequest
	// requestBody := r.Body
	// err := json.NewDecoder(requestBody).Decode(&temp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	log.Println(err)
	// 	return
	// }
	// if temp.Secret != secret {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	log.Println("Unauthorized")
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert Channel Handler"))
}

func AlertUserHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	// log.Println("Alert User Handler")
	// var temp alertRequest
	// requestBody := r.Body
	// err := json.NewDecoder(requestBody).Decode(&temp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	log.Println(err)
	// 	return
	// }
	// if temp.Secret != secret {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	log.Println("Unauthorized")
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert User Handler"))
}
