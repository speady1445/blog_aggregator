package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/speady1445/blog_aggregator/internal/database"
)

func (c *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if !validateFeed(params.Name, params.URL) {
		respondWithError(w, http.StatusBadRequest, "Invalid feed name or url")
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	now := time.Now().UTC()
	feed, err := c.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create feed")
		return
	}

	respondWithJson(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func validateFeed(name, url_ string) bool {
	_, err := url.ParseRequestURI(url_)
	return name != "" && err == nil
}

func (c *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := c.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get feeds")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToFeeds((dbFeeds)))
}
