package main

import (
	"fmt"
	"net/http"
	"github.com/dennisdijkstra/go/server"
	"sync/atomic"
	"encoding/json"
	"log"
	"strings"
	"github.com/joho/godotenv"
	"database/sql"
	"os"
	"github.com/dennisdijkstra/go/internal/database"
	_ "github.com/lib/pq"
	"time"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

type ChirpParams struct {
	Body string `json:"body"`
}

type UserParams struct {
	Email string `json:"email"`
}

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseSuccess struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) writeMetrics(w http.ResponseWriter, req *http.Request) {
	numberOfRequests := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, numberOfRequests)
	w.Write([]byte(template))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

func marshalError(w http.ResponseWriter, err error) {
	log.Printf("Error marshalling JSON: %s", err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	resp := ResponseError{Error: msg}
	data, err := json.Marshal(resp)

	if err != nil {
		marshalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		marshalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func cleanMessage(message string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(message, " ")
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

func validateChirp(w http.ResponseWriter, req *http.Request) {
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

	cleanedMessage := cleanMessage(params.Body)
	body := ResponseSuccess{
		CleanedBody: cleanedMessage,
	}
	respondWithJSON(w, 200, body)
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := UserParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "Something went wrong while creating the user")
	}
	
	body := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 201, body)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	defer db.Close()

	dbQueries := database.New(db)


	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
	}
	mux := http.NewServeMux()

	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fs))

	mux.HandleFunc("/api/get", server.HandleGet)
	mux.HandleFunc("/api/post", server.HandlePost)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.writeMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}