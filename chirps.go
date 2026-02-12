package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"database/sql"
	"time"
	"github.com/google/uuid"
	"github.com/dennisdijkstra/go/internal/database"
)

type ChirpParams struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := ChirpParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	maxLength := 140
	if len(params.Body) > maxLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := getCleanedBody(params.Body)
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, 500, "Something went wrong while creating the chirp")	
		return
	}

	body := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 201, body)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong while fetching chirps")
		return
	}

	response := make([]Chirp, 0, len(chirps))
	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, 400, "Chirp ID is required")
		return
	}

	chirp, err := cfg.db.GetChirpByID(req.Context(), uuid.MustParse(chirpID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Chirp not found")
			return
		}
		respondWithError(w, 500, "Something went wrong while fetching the chirp")
		return
	}

	body := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 200, body)
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