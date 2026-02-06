package main

import (
	"fmt"
	"net/http"
	"github.com/dennisdijkstra/go/server"
	"sync/atomic"
	"encoding/json"
	"log"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type parameters struct {
	Body string `json:"body"`
}

type responseError struct {
	Error string `json:"error"`
}

type responseSuccess struct {
	Valid bool `json:"valid"`
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
	resp := responseError{Error: msg}
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

func ValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
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

	body := responseSuccess{
		Valid: true,
	}
	respondWithJSON(w, 200, body)
}

func main() {
	apiCfg := &apiConfig{}
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

	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.writeMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetrics)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}