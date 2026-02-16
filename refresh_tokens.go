package main

import (
	"net/http"
	"time"

	"github.com/dennisdijkstra/go/internal/auth"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong while parsing the bearer token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating the access token")
		return
	}

	body := RefreshToken{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, body)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong while parsing the bearer token")
		return
	}
	err = cfg.db.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while revoking the refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
