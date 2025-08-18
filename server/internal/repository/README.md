# Repository Layer

This directory contains the data access layer (repository pattern) for the IQ Theory application. The repository layer provides an abstraction between the business logic and data storage, making the code more testable and maintainable.

## Architecture Overview

The repository layer follows these principles:

- **Interface-driven design**: All repositories implement interfaces for easy testing and mocking
- **Single Responsibility**: Each repository handles one domain/entity
- **Dependency Injection**: Repositories are injected into services
- **Context-aware**: All operations accept a context for cancellation and timeouts

## File Structure

```
repository/
├── README.md           # This file
├── interfaces.go       # All repository interface definitions
├── repository.go       # Repository aggregator and constructor
├── user.go            # User & Friendship repository implementations
├── group.go           # Group & GroupMembership repository implementations
└── quiz.go            # Quiz, QuizSession, QuizAnswer, Leaderboard implementations
```

## Repository Interfaces

### Core Entities

- **UserRepository**: User management (CRUD operations)
- **FriendshipRepository**: Friend relationships between users
- **GroupRepository**: Study groups/classrooms
- **GroupMembershipRepository**: User membership in groups

### Quiz System

- **QuizRepository**: Quiz configurations and options
- **QuizSessionRepository**: Individual quiz attempts
- **QuizAnswerRepository**: Individual question answers
- **LeaderboardRepository**: Leaderboard and ranking data

## Usage Examples

### In Services

```go
type UserService struct {
    userRepo       repository.UserRepository
    friendshipRepo repository.FriendshipRepository
}

func NewUserService(repos *repository.Repositories) *UserService {
    return &UserService{
        userRepo:       repos.User,
        friendshipRepo: repos.Friendship,
    }
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
    return s.userRepo.GetByID(ctx, id)
}
```

### In Main Application

```go
// Initialize repositories
repos := repository.NewRepositories(db)

// Pass to services
userService := service.NewUserService(repos)
```

### Testing with Mocks

```go
type mockUserRepo struct{}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    return &models.User{ID: id, Username: "testuser"}, nil
}

// Use in tests
service := &UserService{userRepo: &mockUserRepo{}}
```

## Implementation Status

| Repository                | Status         | Notes                     |
| ------------------------- | -------------- | ------------------------- |
| UserRepository            | ✅ Implemented | Full CRUD operations      |
| FriendshipRepository      | ✅ Implemented | Friend request management |
| GroupRepository           | 🚧 Placeholder | TODO: Implement methods   |
| GroupMembershipRepository | 🚧 Placeholder | TODO: Implement methods   |
| QuizRepository            | 🚧 Placeholder | TODO: Implement methods   |
| QuizSessionRepository     | 🚧 Placeholder | TODO: Implement methods   |
| QuizAnswerRepository      | 🚧 Placeholder | TODO: Implement methods   |
| LeaderboardRepository     | 🚧 Placeholder | TODO: Implement methods   |

## Database Conventions

### Query Patterns

- Use `QueryRowContext` for single record queries
- Use `QueryContext` + `rows.Next()` for multiple records
- Use `ExecContext` for INSERT/UPDATE/DELETE operations
- Always handle `sql.ErrNoRows` appropriately

### Error Handling

```go
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    // ... query setup
    err := row.Scan(/* fields */)
    if err == sql.ErrNoRows {
        return nil, nil  // Not found, not an error
    }
    if err != nil {
        return nil, err  // Actual error
    }
    return user, nil
}
```

### Soft Deletes

- Use `is_active` field for soft deletes where applicable
- Filter by `is_active = true` in queries by default
- Set `is_active = false` and update `updated_at` for deletions

## Adding New Repositories

1. **Add interface** to `interfaces.go`
2. **Add to Repositories struct** in `repository.go`
3. **Add constructor call** in `NewRepositories()`
4. **Create implementation file** (e.g., `new_domain.go`)
5. **Implement all interface methods**
6. **Update this README**

## Best Practices

### Context Usage

- Always pass `context.Context` as the first parameter
- Use `ctx` for cancellation and timeouts
- Don't store context in structs

### SQL Queries

- Use multi-line strings with proper indentation
- Use parameterized queries ($1, $2, etc.) to prevent SQL injection
- Include field names explicitly in SELECT statements

### Error Handling

- Return errors from the database layer up to the service layer
- Don't log errors in repositories (let services decide)
- Use meaningful error messages

### Testing

- Mock repository interfaces in service tests
- Test repository implementations against a test database
- Use table-driven tests for multiple scenarios

## Dependencies

- `github.com/google/uuid` - UUID handling
- `database/sql` - Standard database operations
- `context` - Request context management
- `github.com/lib/pq` - PostgreSQL driver (via database package)

## Related Documentation

- [Database Schema](../../docs/database_schema.md)
- [Service Layer](../service/README.md)
- [Models](../models/README.md)
