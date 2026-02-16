package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseSuccess struct {
	CleanedBody string `json:"cleaned_body"`
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
