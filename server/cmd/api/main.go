package main

// TODO: Implement your main application entry point
// This is where you'll:
// - Load configuration
// - Initialize database connections
// - Set up routing
// - Start the HTTP server

import (
	"net/http"

	"github.com/andy-dam/iq-theory/server/internal/handlers"
	"github.com/andy-dam/iq-theory/server/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initialize services (you'll need to implement these)
	userService := &service.UserService{}
	// quizService := &service.QuizService{}

	// Initialize handlers
	userHandler := &handlers.UserHandler{
		UserService: userService,
	}
	// quizHandler := &handlers.QuizHandler{
	// 	QuizService: quizService,
	// }

	// Define routes
	r.HandleFunc("/example", userHandler.ExampleHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
