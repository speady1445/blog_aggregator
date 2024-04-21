package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, errMsg string) {
	if status <= 500 {
		log.Printf("Status %d response: %s", status, errMsg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, status, errorResponse{Error: errMsg})
}

func respondWithJson(w http.ResponseWriter, status int, content interface{}) {
	w.Header().Set("Content-Type", "application/json")

	returnData, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(status)
	w.Write(returnData)
}
