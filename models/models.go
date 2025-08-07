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

// AlternativeQuizSubmissionRequest represents an alternative format for quiz submission
type AlternativeQuizSubmissionRequest struct {
	QuizID  string            `json:"quizId,omitempty"`                  // Optional quiz identifier
	Answers map[string]string `json:"answers" validate:"required,min=1"` // questionId -> answer text
	UserID  string            `json:"userId,omitempty"`                  // Optional user identifier
}

// QuizSubmissionResponse represents the response for quiz submission
type QuizSubmissionResponse struct {
	Score          int                `json:"score"`
	Total          int                `json:"total"`
	CorrectAnswers int                `json:"correctAnswers"`
	Percentage     float64            `json:"percentage"`
	Passed         bool               `json:"passed"`
	SessionID      string             `json:"sessionId,omitempty"`
	Results        []QuizAnswerResult `json:"results,omitempty"` // Detailed per-question results
}

// QuizAnswerResult represents the result for a single question
type QuizAnswerResult struct {
	QuestionID   int    `json:"questionId"`
	Question     string `json:"question"`
	SelectedID   int    `json:"selectedId"`
	CorrectID    int    `json:"correctId"`
	IsCorrect    bool   `json:"isCorrect"`
	SelectedText string `json:"selectedText,omitempty"`
	CorrectText  string `json:"correctText,omitempty"`
}

// QuestionResponse represents a simplified question format for API responses
type QuestionResponse struct {
	QuestionID    string   `json:"questionId"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}
