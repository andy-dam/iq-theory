package main

import (
	"log"
	"net/http"

	"github.com/andy-dam/iq-theory/server/internal/config"
	"github.com/andy-dam/iq-theory/server/internal/handlers"
	"github.com/andy-dam/iq-theory/server/internal/repository"
	"github.com/andy-dam/iq-theory/server/internal/service"
	"github.com/andy-dam/iq-theory/server/pkg/database"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations (if implemented)
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup routes with all dependencies
	router := setupRoutes(db)

	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupRoutes initializes and configures all routes with their handlers
func setupRoutes(db *database.DB) *mux.Router {
	r := mux.NewRouter()

	// Initialize all repositories
	repos := repository.NewRepositories(db)

	// Initialize all services
	services := service.NewServices(repos)

	// Initialize handlers
	authHandler := &handlers.AuthHandler{
		UserService: services.User,
	}
	// userHandler := &handlers.UserHandler{
	// 	UserService: services.User,
	// }
	// quizHandler := &handlers.QuizHandler{
	// 	QuizService: services.Quiz,
	// }

	// Create API subrouter with /api prefix
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define API routes (these will be prefixed with /api)
	apiRouter.HandleFunc("/health", healthHandler(db)).Methods("GET")
	apiRouter.HandleFunc("/auth/register", authHandler.Register).Methods("POST")

	return r
}

// healthHandler provides a health check endpoint
func healthHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.HealthCheck(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Database connection failed"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
