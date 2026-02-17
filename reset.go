package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetAll(w http.ResponseWriter, r *http.Request) {
	if cfg.environment != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while resetting the database")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics and database reset successfully"))
}
