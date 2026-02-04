package main

import (
	"net/http"
	"github.com/dennisdijkstra/go/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", server.HandleRoot)
	mux.HandleFunc("/get", server.HandleGet)
	mux.HandleFunc("/post", server.HandlePost)
	mux.HandleFunc("/healthz", server.HandleHealthz)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}