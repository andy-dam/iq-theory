# Go REST API Structure

This is a well-organized folder structure for building REST APIs in Go following best practices.

## Project Structure

```
server/
├── cmd/
│   └── api/                    # Application entry points
│       └── main.go            # Main application file
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   ├── handlers/              # HTTP handlers (controllers)
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Data models/structs
│   ├── repository/            # Data access layer
│   ├── service/               # Business logic layer
│   └── utils/                 # Internal utility functions
├── pkg/                       # Public/reusable packages
│   ├── auth/                  # Authentication utilities
│   ├── database/              # Database connection/utilities
│   └── logger/                # Logging utilities
├── migrations/                # Database migrations
├── scripts/                   # Build and deployment scripts
├── tests/                     # Test files
├── docs/                      # Documentation
├── go.mod                     # Go module file
└── go.sum                     # Go dependencies checksum
```

## Directory Explanations

- **cmd/api/**: Contains the main application entry point
- **internal/**: Private code that shouldn't be imported by other projects
- **pkg/**: Public packages that could be reused by other projects
- **migrations/**: Database schema changes
- **scripts/**: Automation scripts for building, testing, deployment
- **tests/**: Integration and end-to-end tests
- **docs/**: API documentation, swagger files, etc.

## Common Patterns

- Use dependency injection
- Implement interfaces for testability
- Separate concerns (handlers, services, repositories)
- Use middleware for cross-cutting concerns
- Keep business logic in service layer
