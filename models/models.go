package models

import "time"

// Option represents a quiz question option
type Option struct {
	ID         int    `json:"id" db:"id"`
	QuestionID int    `json:"question_id" db:"question_id"`
	Text       string `json:"text" db:"text"`
	IsCorrect  bool   `json:"-" db:"is_correct"` // Hidden from JSON response
}

// Question represents a quiz question
type Question struct {
	ID        int       `json:"id" db:"id"`
	Question  string    `json:"question" db:"question"`
	Options   []Option  `json:"options"`
	Answer    int       `json:"answer"` // The ID of the correct option
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// QuizResult represents a user's quiz result
type QuizResult struct {
	ID          int       `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	QuestionID  int       `json:"question_id" db:"question_id"`
	SelectedID  int       `json:"selected_id" db:"selected_id"`
	IsCorrect   bool      `json:"is_correct" db:"is_correct"`
	CompletedAt time.Time `json:"completed_at" db:"completed_at"`
}

// QuizSession represents a complete quiz session
type QuizSession struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	Score       int        `json:"score" db:"score"`
	TotalCount  int        `json:"total_count" db:"total_count"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
}

// QuizAnswer represents a single answer in a quiz submission
type QuizAnswer struct {
	QuestionID int `json:"questionId" validate:"required"`
	AnswerID   int `json:"answerId" validate:"required"`
}

// QuizSubmissionRequest represents the request body for quiz submission
type QuizSubmissionRequest struct {
	Answers []QuizAnswer `json:"answers" validate:"required,min=1"`
	UserID  string       `json:"userId,omitempty"` // Optional user identifier
}

// QuizSubmissionResponse represents the response for quiz submission
type QuizSubmissionResponse struct {
	Score          int     `json:"score"`
	Total          int     `json:"total"`
	CorrectAnswers int     `json:"correctAnswers"`
	Percentage     float64 `json:"percentage"`
	Passed         bool    `json:"passed"`
	SessionID      string  `json:"sessionId,omitempty"`
}
