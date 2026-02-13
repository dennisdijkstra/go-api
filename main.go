package main

import (
	"net/http"
	"sync/atomic"
	"log"
	"database/sql"
	"os"
	"github.com/dennisdijkstra/go/server"
	"github.com/joho/godotenv"
	"github.com/dennisdijkstra/go/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	environment string
	jwtSecret string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	environment := os.Getenv("ENVIRONMENT")
	jwtSecret := os.Getenv("JWT_SECRET")

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
		environment: environment,
		jwtSecret: jwtSecret,
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

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByID)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)

	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.writeMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetAll)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}