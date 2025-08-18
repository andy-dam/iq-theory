# Service Layer

This directory contains the business logic layer for the IQ Theory application. The service layer acts as an intermediary between the handlers (presentation layer) and repositories (data access layer), implementing the core business rules and workflows.

## Architecture Overview

The service layer follows these principles:

- **Interface-driven design**: All services implement interfaces for easy testing and mocking
- **Single Responsibility**: Each service handles one domain/business area
- **Dependency Injection**: Services receive repository interfaces, not concrete implementations
- **Business Logic Isolation**: Complex business rules are encapsulated in services
- **Context-aware**: All operations accept a context for cancellation and timeouts

## File Structure

```
service/
├── README.md           # This file
├── interfaces.go       # All service interface definitions
├── service.go          # Service aggregator and constructor
├── user.go            # User & Friendship service implementations
├── group.go           # Group service implementation
└── quiz.go            # Quiz & Leaderboard service implementations
```

## Service Interfaces

### User Management

- **UserService**: User CRUD, authentication, profile management
- **FriendshipService**: Friend requests and relationships

### Group Management

- **GroupService**: Study groups/classrooms and membership management

### Quiz System

- **QuizService**: Quiz sessions, questions, answers, and scoring
- **LeaderboardService**: Rankings and leaderboard management

## Usage Examples

### In Handlers

```go
type UserHandler struct {
    UserService service.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    // ... parse request

    user, err := h.UserService.CreateUser(r.Context(), &req)
    if err != nil {
        // ... handle error
    }
    // ... return response
}
```

### In Main Application

```go
// Initialize services
repos := repository.NewRepositories(db)
services := service.NewServices(repos)

// Pass to handlers
userHandler := &handlers.UserHandler{
    UserService: services.User,
}
```

### Testing with Mocks

```go
type mockUserService struct{}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    return &models.User{ID: uuid.New(), Username: req.Username}, nil
}

// Use in tests
handler := &UserHandler{UserService: &mockUserService{}}
```

## Implementation Status

| Service            | Status         | Notes                                           |
| ------------------ | -------------- | ----------------------------------------------- |
| UserService        | ✅ Implemented | Full user management with bcrypt authentication |
| FriendshipService  | ✅ Implemented | Friend request management                       |
| GroupService       | 🚧 Placeholder | TODO: Implement business logic                  |
| QuizService        | 🚧 Placeholder | TODO: Implement quiz logic                      |
| LeaderboardService | 🚧 Placeholder | TODO: Implement leaderboard logic               |

## Business Logic Examples

### User Registration

```go
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // 1. Validate input
    // 2. Check for existing email/username
    // 3. Hash password
    // 4. Create user record
    // 5. Send verification email (future)
    // 6. Return user (without sensitive data)
}
```

### Friend Request Flow

```go
func (s *friendshipService) SendFriendRequest(ctx context.Context, requesterID, addresseeID uuid.UUID) error {
    // 1. Validate users exist
    // 2. Check for existing friendship
    // 3. Create pending friendship record
    // 4. Send notification (future)
}
```

### Quiz Session Management

```go
func (s *quizService) CreateQuizSession(ctx context.Context, userID uuid.UUID, clef string, duration int, maxLedgerLines int) (*models.QuizSession, error) {
    // 1. Validate configuration exists
    // 2. Create session record
    // 3. Generate question pool
    // 4. Set session state to "ready"
    // 5. Return session details
}
```

## Security Considerations

### Password Management

- **Hashing**: Uses bcrypt with default cost for secure password storage
- **Validation**: Enforce strong password policies
- **Change Flow**: Verify current password before allowing changes

### Authentication

- **Session Management**: Implement JWT or session tokens
- **Rate Limiting**: Prevent brute force attacks
- **Email Verification**: Verify email addresses before activation

### Authorization

- **Group Permissions**: Only admins can remove members/change roles
- **Quiz Access**: Users can only access their own quiz sessions
- **Data Privacy**: Users can only see their own profile data

## Error Handling Patterns

### Business Logic Errors

```go
// Return descriptive errors for business rule violations
if existingUser != nil {
    return nil, fmt.Errorf("email already exists")
}

// Wrap repository errors with context
if err := s.userRepo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

### Validation

```go
// Validate at service layer before repository calls
if req.Password == "" {
    return nil, fmt.Errorf("password is required")
}
```

## Performance Considerations

### Database Optimization

- Use transactions for multi-step operations
- Implement proper indexing strategies
- Consider caching for frequently accessed data

### Quiz Performance

- Pre-generate question pools
- Cache leaderboard data
- Implement pagination for large result sets

## Testing Strategy

### Unit Tests

- Mock repository interfaces
- Test business logic in isolation
- Verify error handling paths

### Integration Tests

- Test with real database
- Verify end-to-end workflows
- Test transaction rollbacks

## Dependencies

- `context` - Request context management
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto/bcrypt` - Password hashing
- Repository interfaces - Data access abstraction

## Adding New Services

1. **Add interface** to `interfaces.go`
2. **Add to Services struct** in `service.go`
3. **Add constructor call** in `NewServices()`
4. **Create implementation file** (e.g., `new_domain.go`)
5. **Implement all interface methods**
6. **Add business logic and validation**
7. **Update handlers** to use new service
8. **Update this README**

## Best Practices

### Context Usage

- Always accept `context.Context` as first parameter
- Pass context to repository calls
- Use context for request scoping and cancellation

### Error Messages

- Use descriptive error messages for user-facing errors
- Wrap repository errors with business context
- Don't expose internal implementation details

### Validation

- Validate inputs at service layer
- Return early for invalid data
- Use consistent error formats

### Testing

- Mock all dependencies (repositories)
- Test happy path and error scenarios
- Use table-driven tests for multiple cases

## Related Documentation

- [Repository Layer](../repository/README.md)
- [Handlers Layer](../handlers/README.md)
- [Models](../models/README.md)
- [API Documentation](../../docs/api.md)
