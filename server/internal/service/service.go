package service

import (
	"github.com/andy-dam/iq-theory/server/internal/repository"
)

// Services aggregates all service implementations
type Services struct {
	User        UserService
	Friendship  FriendshipService
	Group       GroupService
	Quiz        QuizService
	Leaderboard LeaderboardService
}

// NewServices creates a new instance of all services
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User:        NewUserService(repos),
		Friendship:  NewFriendshipService(repos),
		Group:       NewGroupService(repos),
		Quiz:        NewQuizService(repos),
		Leaderboard: NewLeaderboardService(repos),
	}
}
