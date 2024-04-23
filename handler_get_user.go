package main

import (
	"net/http"

	"github.com/speady1445/blog_aggregator/internal/auth"
)

func (c *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		switch err {
		case auth.ErrMissingAuthHeader:
			respondWithError(w, http.StatusBadRequest, "Missing Authorization header")
		case auth.ErrInvalidAuthHeader:
			respondWithError(w, http.StatusBadRequest, "Invalid Authorization header")
		}
		return
	}

	dbUser, err := c.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get user")
		return
	}

	respondWithJson(w, http.StatusOK, databaseUserToUser(dbUser))
}
