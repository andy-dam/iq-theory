package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Username      string    `json:"username" db:"username"`
	DisplayName   string    `json:"display_name" db:"display_name"`
	PasswordHash  string    `json:"-" db:"password_hash"` // Hidden from JSON
	AvatarURL     *string   `json:"avatar_url" db:"avatar_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
}

// Group represents a classroom or study group
type Group struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	JoinCode    string    `json:"join_code" db:"join_code"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	MaxMembers  int       `json:"max_members" db:"max_members"`
}

// GroupMembership represents user membership in a group
type GroupMembership struct {
	ID       uuid.UUID `json:"id" db:"id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	GroupID  uuid.UUID `json:"group_id" db:"group_id"`
	Role     string    `json:"role" db:"role"` // admin, member
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

// Friendship represents friend relationships between users
type Friendship struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RequesterID uuid.UUID `json:"requester_id" db:"requester_id"`
	AddresseeID uuid.UUID `json:"addressee_id" db:"addressee_id"`
	Status      string    `json:"status" db:"status"` // pending, accepted, declined, blocked
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ClefType represents a musical clef type
type ClefType struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// DurationOption represents a quiz duration option
type DurationOption struct {
	ID              int    `json:"id" db:"id"`
	DurationSeconds int    `json:"duration_seconds" db:"duration_seconds"`
	DisplayName     string `json:"display_name" db:"display_name"`
	IsActive        bool   `json:"is_active" db:"is_active"`
}

// LedgerLineOption represents a ledger line limit option
type LedgerLineOption struct {
	ID          int    `json:"id" db:"id"`
	MaxLines    int    `json:"max_lines" db:"max_lines"`
	DisplayName string `json:"display_name" db:"display_name"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// AvailableQuizConfiguration represents a dynamically generated quiz configuration
type AvailableQuizConfiguration struct {
	ConfigurationName string `json:"configuration_name" db:"configuration_name"`
	Clef              string `json:"clef" db:"clef"`
	ClefDisplay       string `json:"clef_display" db:"clef_display"`
	DurationSeconds   int    `json:"duration_seconds" db:"duration_seconds"`
	DurationDisplay   string `json:"duration_display" db:"duration_display"`
	MaxLedgerLines    int    `json:"max_ledger_lines" db:"max_ledger_lines"`
	LedgerDisplay     string `json:"ledger_display" db:"ledger_display"`
	IsAvailable       bool   `json:"is_available" db:"is_available"`
}

// QuizSession represents an individual quiz attempt
type QuizSession struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	// Store parameters directly instead of referencing a configuration
	Clef            string `json:"clef" db:"clef"`
	DurationSeconds int    `json:"duration_seconds" db:"duration_seconds"`
	MaxLedgerLines  int    `json:"max_ledger_lines" db:"max_ledger_lines"`

	Score              int        `json:"score" db:"score"`
	TotalQuestions     int        `json:"total_questions" db:"total_questions"`
	CorrectAnswers     int        `json:"correct_answers" db:"correct_answers"`
	TimeTakenSeconds   int        `json:"time_taken_seconds" db:"time_taken_seconds"`
	StartedAt          time.Time  `json:"started_at" db:"started_at"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
	Status             string     `json:"status" db:"status"` // in_progress, completed, abandoned
	AccuracyPercentage float64    `json:"accuracy_percentage" db:"accuracy_percentage"`
}

// QuizAnswer represents an individual question answer within a quiz session
type QuizAnswer struct {
	ID             uuid.UUID `json:"id" db:"id"`
	QuizSessionID  uuid.UUID `json:"quiz_session_id" db:"quiz_session_id"`
	QuestionNumber int       `json:"question_number" db:"question_number"`
	NoteImage      string    `json:"note_image" db:"note_image"`     // e.g., "A4.png"
	CorrectNote    string    `json:"correct_note" db:"correct_note"` // e.g., "A4"
	UserAnswer     *string   `json:"user_answer" db:"user_answer"`
	IsCorrect      bool      `json:"is_correct" db:"is_correct"`
	TimeTakenMs    int       `json:"time_taken_ms" db:"time_taken_ms"`
	AnsweredAt     time.Time `json:"answered_at" db:"answered_at"`
}

// LeaderboardEntry represents a leaderboard entry (from materialized view)
type LeaderboardEntry struct {
	Clef            string     `json:"clef" db:"clef"`
	DurationSeconds int        `json:"duration_seconds" db:"duration_seconds"`
	MaxLedgerLines  int        `json:"max_ledger_lines" db:"max_ledger_lines"`
	QuizName        string     `json:"quiz_name" db:"quiz_name"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Username        string     `json:"username" db:"username"`
	DisplayName     string     `json:"display_name" db:"display_name"`
	BestScore       int        `json:"best_score" db:"best_score"`
	BestAccuracy    float64    `json:"best_accuracy" db:"best_accuracy"`
	FastestTime     int        `json:"fastest_time" db:"fastest_time"`
	TotalAttempts   int        `json:"total_attempts" db:"total_attempts"`
	AverageScore    float64    `json:"average_score" db:"average_score"`
	LastAttempt     *time.Time `json:"last_attempt" db:"last_attempt"`
	GlobalRank      int        `json:"global_rank" db:"global_rank"`
}

// DTOs for API requests/responses

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
	Password    string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description"`
	MaxMembers  int     `json:"max_members" validate:"min=1,max=1000"`
}

// JoinGroupRequest represents the request to join a group
type JoinGroupRequest struct {
	JoinCode string `json:"join_code" validate:"required"`
}

// StartQuizRequest represents the request to start a new quiz
type StartQuizRequest struct {
	Clef            string `json:"clef" validate:"required,oneof=treble bass alto tenor"`
	DurationSeconds int    `json:"duration_seconds" validate:"required,oneof=30 60 120"`
	MaxLedgerLines  int    `json:"max_ledger_lines" validate:"required,min=0,max=3"`
}

// SubmitAnswerRequest represents the request to submit an answer
type SubmitAnswerRequest struct {
	QuestionNumber int    `json:"question_number" validate:"required,min=1"`
	UserAnswer     string `json:"user_answer" validate:"required"`
	TimeTakenMs    int    `json:"time_taken_ms" validate:"required,min=0"`
}

// LeaderboardRequest represents the request for leaderboard data
type LeaderboardRequest struct {
	Clef            string     `json:"clef" validate:"required,oneof=treble bass alto tenor"`
	DurationSeconds int        `json:"duration_seconds" validate:"required,oneof=30 60 120"`
	MaxLedgerLines  int        `json:"max_ledger_lines" validate:"required,min=0,max=3"`
	Scope           string     `json:"scope" validate:"required,oneof=global group friends"` // global, group, friends
	GroupID         *uuid.UUID `json:"group_id,omitempty"`
	Limit           int        `json:"limit" validate:"min=1,max=100"`
	Offset          int        `json:"offset" validate:"min=0"`
}

// QuizOptionsResponse represents available quiz configuration options
type QuizOptionsResponse struct {
	Clefs       []ClefType         `json:"clefs"`
	Durations   []DurationOption   `json:"durations"`
	LedgerLines []LedgerLineOption `json:"ledger_lines"`
}

// Note represents a musical note with its images
type Note struct {
	Note string   `json:"note"`
	Imgs []string `json:"imgs"`
}

// Question represents a quiz question
type Question struct {
	Image         string   `json:"img"`
	Note          string   `json:"note"`
	CorrectAnswer string   `json:"correctAnswer"`
	Answers       []string `json:"answers"`
}
