package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/andy-dam/iq-theory/server/internal/models"
	"github.com/andy-dam/iq-theory/server/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// userService implements the UserServiceInterface
type userService struct {
	userRepo       repository.UserRepository
	friendshipRepo repository.FriendshipRepository
}

// NewUserService creates a new user service instance
func NewUserService(repos *repository.Repositories) UserService {
	return &userService{
		userRepo:       repos.User,
		friendshipRepo: repos.Friendship,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	now := time.Now()
	user := &models.User{
		ID:            uuid.New(),
		Email:         req.Email,
		Username:      req.Username,
		DisplayName:   req.DisplayName,
		PasswordHash:  string(hashedPassword),
		CreatedAt:     now,
		UpdatedAt:     now,
		IsActive:      true,
		EmailVerified: false,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AuthenticateUser validates user credentials
func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateProfile updates user profile information
func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName string, avatarURL *string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.DisplayName = displayName
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// VerifyEmail marks a user's email as verified
func (s *userService) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// friendshipService implements the FriendshipServiceInterface
type friendshipService struct {
	friendshipRepo repository.FriendshipRepository
	userRepo       repository.UserRepository
}

// NewFriendshipService creates a new friendship service instance
func NewFriendshipService(repos *repository.Repositories) FriendshipService {
	return &friendshipService{
		friendshipRepo: repos.Friendship,
		userRepo:       repos.User,
	}
}

// SendFriendRequest sends a friend request
func (s *friendshipService) SendFriendRequest(ctx context.Context, requesterID, addresseeID uuid.UUID) error {
	if requesterID == addresseeID {
		return fmt.Errorf("cannot send friend request to yourself")
	}

	// Check if users exist
	if _, err := s.userRepo.GetByID(ctx, requesterID); err != nil {
		return fmt.Errorf("requester not found: %w", err)
	}
	if _, err := s.userRepo.GetByID(ctx, addresseeID); err != nil {
		return fmt.Errorf("addressee not found: %w", err)
	}

	// Create friendship request
	now := time.Now()
	friendship := &models.Friendship{
		ID:          uuid.New(),
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.friendshipRepo.Create(ctx, friendship); err != nil {
		return fmt.Errorf("failed to create friend request: %w", err)
	}

	return nil
}

// AcceptFriendRequest accepts a friend request
func (s *friendshipService) AcceptFriendRequest(ctx context.Context, friendshipID uuid.UUID) error {
	return s.friendshipRepo.UpdateStatus(ctx, friendshipID, "accepted")
}

// DeclineFriendRequest declines a friend request
func (s *friendshipService) DeclineFriendRequest(ctx context.Context, friendshipID uuid.UUID) error {
	return s.friendshipRepo.UpdateStatus(ctx, friendshipID, "declined")
}

// RemoveFriend removes a friendship
func (s *friendshipService) RemoveFriend(ctx context.Context, friendshipID uuid.UUID) error {
	return s.friendshipRepo.Delete(ctx, friendshipID)
}

// GetUserFriends gets all friends for a user
func (s *friendshipService) GetUserFriends(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error) {
	return s.friendshipRepo.GetUserFriends(ctx, userID)
}

// GetPendingRequests gets pending friend requests for a user
func (s *friendshipService) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error) {
	// TODO: Implement method in repository to get pending requests
	// For now, return empty slice
	return []*models.Friendship{}, nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
