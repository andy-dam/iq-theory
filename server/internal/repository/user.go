package repository

import (
	"context"
	"database/sql"

	"github.com/andy-dam/iq-theory/server/internal/models"
	"github.com/andy-dam/iq-theory/server/pkg/database"
	"github.com/google/uuid"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, display_name, password_hash, avatar_url, created_at, updated_at, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.DisplayName, user.PasswordHash,
		user.AvatarURL, user.CreatedAt, user.UpdatedAt, user.IsActive, user.EmailVerified)

	return err
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, 
		       created_at, updated_at, is_active, email_verified
		FROM users 
		WHERE id = $1 AND is_active = true`

	user := &models.User{}
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.EmailVerified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, 
		       created_at, updated_at, is_active, email_verified
		FROM users 
		WHERE email = $1 AND is_active = true`

	user := &models.User{}
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.EmailVerified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by their username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, 
		       created_at, updated_at, is_active, email_verified
		FROM users 
		WHERE username = $1 AND is_active = true`

	user := &models.User{}
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.PasswordHash,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.EmailVerified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = $2, username = $3, display_name = $4, password_hash = $5, 
		    avatar_url = $6, updated_at = $7, is_active = $8, email_verified = $9
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.DisplayName, user.PasswordHash,
		user.AvatarURL, user.UpdatedAt, user.IsActive, user.EmailVerified)

	return err
}

// Delete soft deletes a user (sets is_active to false)
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// friendshipRepository implements the FriendshipRepository interface
type friendshipRepository struct {
	db *database.DB
}

// NewFriendshipRepository creates a new friendship repository instance
func NewFriendshipRepository(db *database.DB) FriendshipRepository {
	return &friendshipRepository{db: db}
}

// Create creates a new friendship request
func (r *friendshipRepository) Create(ctx context.Context, friendship *models.Friendship) error {
	query := `
		INSERT INTO friendships (id, requester_id, addressee_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		friendship.ID, friendship.RequesterID, friendship.AddresseeID,
		friendship.Status, friendship.CreatedAt, friendship.UpdatedAt)

	return err
}

// GetByID retrieves a friendship by ID
func (r *friendshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Friendship, error) {
	query := `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships 
		WHERE id = $1`

	friendship := &models.Friendship{}
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&friendship.ID, &friendship.RequesterID, &friendship.AddresseeID,
		&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return friendship, nil
}

// GetUserFriends retrieves all friendships for a user
func (r *friendshipRepository) GetUserFriends(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error) {
	query := `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships 
		WHERE (requester_id = $1 OR addressee_id = $1) AND status = 'accepted'
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*models.Friendship
	for rows.Next() {
		friendship := &models.Friendship{}
		err := rows.Scan(
			&friendship.ID, &friendship.RequesterID, &friendship.AddresseeID,
			&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, friendship)
	}

	return friendships, rows.Err()
}

// UpdateStatus updates the status of a friendship
func (r *friendshipRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE friendships SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

// Delete removes a friendship
func (r *friendshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM friendships WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
