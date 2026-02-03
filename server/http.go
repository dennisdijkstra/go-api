package server

import (
	"fmt"
	"net/http"
)

func HandleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get")
}

func HandlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post")
}