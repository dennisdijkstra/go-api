package server

import (
	"fmt"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, req *http.Request) {
	fs := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	fs.ServeHTTP(w, req)
}

func HandleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get")
}

func HandlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post")
}

func HandleHealthz(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("charset", "utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}