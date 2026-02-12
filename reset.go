package main

import (
	"net/http"
)

func (cfg *apiConfig) resetAll(w http.ResponseWriter, req *http.Request) {
	if cfg.environment != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong while resetting the database")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics and database reset successfully"))
}
