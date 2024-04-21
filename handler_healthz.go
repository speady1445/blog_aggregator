package main

import (
	"net/http"
)

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	respondWithJson(w, http.StatusOK, response{Status: "ok"})
}
