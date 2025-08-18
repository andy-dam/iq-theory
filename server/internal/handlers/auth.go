package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andy-dam/iq-theory/server/internal/models"
	"github.com/andy-dam/iq-theory/server/internal/service"
)

type AuthHandler struct {
	UserService service.UserService
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Handle user registration
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Call the service to register the user
	usr, err := ah.UserService.CreateUser(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "Failed to register user"}`, http.StatusInternalServerError)
		return
	}

	// Return the created user (without sensitive data)
	response, err := json.Marshal(usr)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
