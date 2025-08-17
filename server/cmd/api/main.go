package main

// TODO: Implement your main application entry point
// This is where you'll:
// - Load configuration
// - Initialize database connections
// - Set up routing
// - Start the HTTP server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	
	
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
