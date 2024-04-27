package main

import (
	"net/http"

	"github.com/speady1445/blog_aggregator/internal/auth"
	"github.com/speady1445/blog_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (c *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, dbUser)
	}
}
