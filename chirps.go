package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dennisdijkstra/go/internal/database"
	"github.com/google/uuid"
)

type ChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	userID, code, msg, ok := cfg.requireJWTUserID(r)
	if !ok {
		respondWithError(w, code, msg)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := ChirpParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong while decoding the request body")
		return
	}

	maxLength := 140
	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := getCleanedBody(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating the chirp")
		return
	}

	body := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, body)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var authorID uuid.UUID
	var chirps []database.Chirp
	var err error

	authorIDQuery := r.URL.Query().Get("author_id")
	sortQuery := r.URL.Query().Get("sort")

	if authorIDQuery != "" {
		authorID, err = uuid.Parse(authorIDQuery)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
	}

	if authorID != uuid.Nil {
		chirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching chirps")
			return
		}
	}

	response := make([]Chirp, 0, len(chirps))
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	sort.Slice(response, func(i, j int) bool {
		if sortQuery == "desc" {
			return response[i].CreatedAt.After(response[j].CreatedAt)
		}
		return response[i].CreatedAt.Before(response[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching the chirp")
		return
	}

	body := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, body)
}

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}

	userID, code, msg, ok := cfg.requireJWTUserID(r)
	if !ok {
		respondWithError(w, code, msg)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while fetching the chirp")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not allowed to delete this chirp")
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while deleting the chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getCleanedBody(body string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	cleanedMessage := make([]string, 0, len(words))

	for _, word := range words {
		isProfane := false

		for _, profanity := range profanities {
			if strings.ToLower(word) == profanity {
				isProfane = true
				break
			}
		}

		if isProfane {
			cleanedMessage = append(cleanedMessage, "****")
		} else {
			cleanedMessage = append(cleanedMessage, word)
		}
	}

	return strings.Join(cleanedMessage, " ")
}
