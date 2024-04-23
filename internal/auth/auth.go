package auth

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "ApiKey" {
		return "", ErrInvalidAuthHeader
	}

	return authParts[1], nil
}
