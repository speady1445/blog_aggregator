package main

import (
	"net/http"
	"strconv"

	"github.com/speady1445/blog_aggregator/internal/database"
)

func (c *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	dbPosts, err := c.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: dbUser.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get posts")
		return
	}

	respondWithJson(w, http.StatusOK, databasePostsToPosts(dbPosts))
}
