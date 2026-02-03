package main

import (
	"net/http"
	"github.com/dennisdijkstra/go/server"
)

func main() {
	http.HandleFunc("/get", server.HandleGet)
	http.HandleFunc("/post", server.HandlePost)

	httpServer := &http.Server{
		Addr: ":8080",
	}
	httpServer.ListenAndServe()
}