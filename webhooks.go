package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dennisdijkstra/go/internal/database"
	"github.com/google/uuid"
)

type WebhookParams struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := WebhookParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateUserIsChirpyRed(req.Context(), database.UpdateUserIsChirpyRedParams{
		ID:          params.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while updating the user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
