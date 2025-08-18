package service

import (
	"context"

	"github.com/andy-dam/iq-theory/server/internal/models"
	"github.com/google/uuid"
)

// UserService defines methods for user-related business logic
type UserService interface {
	// User management
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Authentication
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

	// Profile management
	UpdateProfile(ctx context.Context, userID uuid.UUID, displayName string, avatarURL *string) error
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
}

// FriendshipService defines methods for friendship-related business logic
type FriendshipService interface {
	SendFriendRequest(ctx context.Context, requesterID, addresseeID uuid.UUID) error
	AcceptFriendRequest(ctx context.Context, friendshipID uuid.UUID) error
	DeclineFriendRequest(ctx context.Context, friendshipID uuid.UUID) error
	RemoveFriend(ctx context.Context, friendshipID uuid.UUID) error
	GetUserFriends(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error)
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*models.Friendship, error)
}

// GroupService defines methods for group-related business logic
type GroupService interface {
	CreateGroup(ctx context.Context, creatorID uuid.UUID, name, description string, maxMembers int) (*models.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetGroupByJoinCode(ctx context.Context, joinCode string) (*models.Group, error)
	GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error

	// Group membership management
	JoinGroup(ctx context.Context, userID uuid.UUID, joinCode string) error
	JoinGroupByID(ctx context.Context, userID, groupID uuid.UUID) error
	LeaveGroup(ctx context.Context, userID, groupID uuid.UUID) error
	RemoveMember(ctx context.Context, adminID, memberID, groupID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, adminID, memberID, groupID uuid.UUID, role string) error
	GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*models.GroupMembership, error)
}

// QuizService defines methods for quiz-related business logic
type QuizService interface {
	// Quiz configuration
	GetAvailableConfigurations(ctx context.Context) ([]*models.AvailableQuizConfiguration, error)
	GetClefTypes(ctx context.Context) ([]*models.ClefType, error)
	GetDurationOptions(ctx context.Context) ([]*models.DurationOption, error)
	GetLedgerLineOptions(ctx context.Context) ([]*models.LedgerLineOption, error)

	// Quiz session management
	CreateQuizSession(ctx context.Context, userID uuid.UUID, clef string, duration int, maxLedgerLines int) (*models.QuizSession, error)
	GetQuizSession(ctx context.Context, sessionID uuid.UUID) (*models.QuizSession, error)
	GetUserQuizSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*models.QuizSession, error)
	StartQuizSession(ctx context.Context, sessionID uuid.UUID) error
	CompleteQuizSession(ctx context.Context, sessionID uuid.UUID) (*models.QuizSession, error)
	AbandonQuizSession(ctx context.Context, sessionID uuid.UUID) error

	// Question and answer management
	GetNextQuestion(ctx context.Context, sessionID uuid.UUID) (string, error) // Returns note to identify
	SubmitAnswer(ctx context.Context, sessionID uuid.UUID, questionNumber int, userAnswer string, timeTakenMs int) (*models.QuizAnswer, error)
	GetSessionAnswers(ctx context.Context, sessionID uuid.UUID) ([]*models.QuizAnswer, error)

	// Results and scoring
	CalculateSessionScore(ctx context.Context, sessionID uuid.UUID) (int, float64, error) // score, accuracy
	GetSessionResults(ctx context.Context, sessionID uuid.UUID) (*models.QuizSession, error)
}

// LeaderboardService defines methods for leaderboard-related business logic
type LeaderboardService interface {
	GetGlobalLeaderboard(ctx context.Context, clef string, duration int, maxLedgerLines int, limit int) ([]*models.LeaderboardEntry, error)
	GetUserRanking(ctx context.Context, userID uuid.UUID, clef string, duration int, maxLedgerLines int) (*models.LeaderboardEntry, error)
	GetGroupLeaderboard(ctx context.Context, groupID uuid.UUID, clef string, duration int, maxLedgerLines int, limit int) ([]*models.LeaderboardEntry, error)
	RefreshLeaderboards(ctx context.Context) error
}
