package server

import (
	"fmt"
	"net/http"
)

func HandlerGet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get")
}

func HandlerPost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post")
}
