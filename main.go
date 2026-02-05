package main

import (
	"fmt"
	"net/http"
	"github.com/dennisdijkstra/go/server"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) writeMetrics(w http.ResponseWriter, req *http.Request) {
	numberOfRequests := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits: " + fmt.Sprint(numberOfRequests)))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
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

	mux.HandleFunc("GET /api/metrics", apiCfg.writeMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.resetMetrics)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}