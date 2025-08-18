package repository

import (
	"github.com/andy-dam/iq-theory/server/pkg/database"
)

// Repositories aggregates all repository implementations
type Repositories struct {
	User            UserRepository
	Friendship      FriendshipRepository
	Group           GroupRepository
	GroupMembership GroupMembershipRepository
	Quiz            QuizRepository
	QuizSession     QuizSessionRepository
	QuizAnswer      QuizAnswerRepository
	Leaderboard     LeaderboardRepository
}

// NewRepositories creates a new instance of all repositories
func NewRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:            NewUserRepository(db),
		Friendship:      NewFriendshipRepository(db),
		Group:           NewGroupRepository(db),
		GroupMembership: NewGroupMembershipRepository(db),
		Quiz:            NewQuizRepository(db),
		QuizSession:     NewQuizSessionRepository(db),
		QuizAnswer:      NewQuizAnswerRepository(db),
		Leaderboard:     NewLeaderboardRepository(db),
	}
}
