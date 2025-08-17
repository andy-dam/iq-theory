package handlers

// TODO: Implement your HTTP handlers (controllers)
// These handle incoming HTTP requests and return responses
// Keep them thin - delegate business logic to services

// Examples of what you might have:
// - AuthHandler (login, register, logout)
// - UserHandler (CRUD operations for users)
// - QuizHandler (CRUD operations for quizzes)
// - etc.

import (
	"github.com/andy-dam/iq-theory/server/internal/service"
)

type Handler struct {
	UserService service.UserService
	QuizService service.QuizService
}