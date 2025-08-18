package repository

import (
	"context"

	"github.com/andy-dam/iq-theory/server/internal/models"
	"github.com/google/uuid"
)

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// FriendshipRepository defines methods for friendship data access
type FriendshipRepository interface {
	Create(ctx context.Context, friendship *models.Friendship) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Friendship, error)
	GetUserFriends(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// GroupRepository defines methods for group data access
type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	GetByJoinCode(ctx context.Context, joinCode string) (*models.Group, error)
	GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// GroupMembershipRepository defines methods for group membership data access
type GroupMembershipRepository interface {
	Create(ctx context.Context, membership *models.GroupMembership) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.GroupMembership, error)
	GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*models.GroupMembership, error)
	GetUserMemberships(ctx context.Context, userID uuid.UUID) ([]*models.GroupMembership, error)
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// QuizRepository defines methods for quiz configuration data access
type QuizRepository interface {
	GetClefTypes(ctx context.Context) ([]*models.ClefType, error)
	GetDurationOptions(ctx context.Context) ([]*models.DurationOption, error)
	GetLedgerLineOptions(ctx context.Context) ([]*models.LedgerLineOption, error)
	GetAvailableConfigurations(ctx context.Context) ([]*models.AvailableQuizConfiguration, error)
}

// QuizSessionRepository defines methods for quiz session data access
type QuizSessionRepository interface {
	Create(ctx context.Context, session *models.QuizSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.QuizSession, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*models.QuizSession, error)
	Update(ctx context.Context, session *models.QuizSession) error
	Complete(ctx context.Context, id uuid.UUID, score int, timeTaken int) error
}

// QuizAnswerRepository defines methods for quiz answer data access
type QuizAnswerRepository interface {
	Create(ctx context.Context, answer *models.QuizAnswer) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.QuizAnswer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.QuizAnswer, error)
	Update(ctx context.Context, answer *models.QuizAnswer) error
}

// LeaderboardRepository defines methods for leaderboard data access
type LeaderboardRepository interface {
	GetGlobalLeaderboard(ctx context.Context, clef string, duration int, maxLedgerLines int, limit int) ([]*models.LeaderboardEntry, error)
	GetUserRanking(ctx context.Context, userID uuid.UUID, clef string, duration int, maxLedgerLines int) (*models.LeaderboardEntry, error)
	RefreshLeaderboard(ctx context.Context) error
}
