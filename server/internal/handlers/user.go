package handlers

import (
	"net/http"

	"github.com/andy-dam/iq-theory/server/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

// ExampleHandler handles the /example route
func (uh *UserHandler) ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Hello from user handler"}`))
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Handle user-related requests
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Get all users"}`))
}

// TODO: Add more user-related handlers here
// Examples:
// - GetUser (by ID)
// - CreateUser
// - UpdateUser
// - DeleteUser
