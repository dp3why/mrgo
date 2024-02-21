package handler

import (
	"net/http"
)

func CheckRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Server is running successfully!"))
}