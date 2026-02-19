package main

import (
	"net/http"

	"github.com/dennisdijkstra/go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) requireJWTUserID(r *http.Request) (uuid.UUID, int, string, bool) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, http.StatusUnauthorized, "Something went wrong while parsing the bearer token", false
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		return uuid.Nil, http.StatusUnauthorized, "Unauthorized", false
	}

	return userID, 0, "", true
}
