package main

import (
	"net/http"
	"github.com/dennisdijkstra/go/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", server.HandleGet)
	mux.HandleFunc("/post", server.HandlePost)

	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", fs)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}