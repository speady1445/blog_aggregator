package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/speady1445/blog_aggregator/internal/database"
)

func (c *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	feedFollow, err := c.feedFollowsCreate(r, user.ID, params.FeedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create feed follow")
	}

	respondWithJson(w, http.StatusCreated, feedFollow)
}

func (c *apiConfig) feedFollowsCreate(r *http.Request, userID uuid.UUID, feedID uuid.UUID) (FeedFollow, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return FeedFollow{}, err
	}

	now := time.Now().UTC()
	dbFeedFollow, err := c.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
	})
	if err != nil {
		return FeedFollow{}, err
	}

	return databaseFeedFollowToFeedFollow(dbFeedFollow), nil
}

func (c *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := r.PathValue("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID.")
		return
	}

	err = c.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete feed follow")
		return
	}

	respondWithJson(w, http.StatusOK, nil)
}

func (c *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	dbFeedFollows, err := c.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get feed follows")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedFollowsToFeedFollows((dbFeedFollows)))
}
