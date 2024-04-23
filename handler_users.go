package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/speady1445/blog_aggregator/internal/auth"
	"github.com/speady1445/blog_aggregator/internal/database"
)

func (c *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if !validUserName(params.Name) {
		respondWithError(w, http.StatusBadRequest, "Invalid user name")
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	now := time.Now().UTC()
	dbUser, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJson(w, http.StatusCreated, databaseUserToUser(dbUser))
}

func validUserName(name string) bool {
	return name != ""
}

func (c *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request) {
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
