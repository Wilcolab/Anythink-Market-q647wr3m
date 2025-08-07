package services

import (
	"fmt"
	"go-quiz-api/models"
	"go-quiz-api/repository"
	"strings"
	"time"
)

// QuizService handles quiz-related business logic
type QuizService struct {
	questionRepo *repository.QuestionRepository
}

// NewQuizService creates a new quiz service
func NewQuizService(questionRepo *repository.QuestionRepository) *QuizService {
	return &QuizService{
		questionRepo: questionRepo,
	}
}

// SubmitQuiz processes a quiz submission and returns the score
func (s *QuizService) SubmitQuiz(submission *models.QuizSubmissionRequest) (*models.QuizSubmissionResponse, error) {
	// Validate submission
	if len(submission.Answers) == 0 {
		return nil, fmt.Errorf("no answers provided")
	}

	// Get all questions to validate against
	allQuestions, err := s.questionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	if len(allQuestions) == 0 {
		return nil, fmt.Errorf("no questions available")
	}

	// Create a map for quick question lookup
	questionMap := make(map[int]*models.Question)
	for i := range allQuestions {
		questionMap[allQuestions[i].ID] = &allQuestions[i]
	}

	// Validate that all required questions are answered
	if err := s.validateAnswers(submission.Answers, questionMap); err != nil {
		return nil, err
	}

	// Calculate score
	score, correctAnswers := s.calculateScore(submission.Answers, questionMap)

	total := len(allQuestions)
	percentage := float64(score) / float64(total) * 100
	passed := percentage >= 60.0 // 60% passing threshold

	// Generate session ID (simple timestamp-based)
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())

	response := &models.QuizSubmissionResponse{
		Score:          score,
		Total:          total,
		CorrectAnswers: correctAnswers,
		Percentage:     percentage,
		Passed:         passed,
		SessionID:      sessionID,
	}

	return response, nil
}

// validateAnswers ensures all questions are answered and answers are valid
func (s *QuizService) validateAnswers(answers []models.QuizAnswer, questionMap map[int]*models.Question) error {
	answeredQuestions := make(map[int]bool)

	for _, answer := range answers {
		// Check if question exists
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			return fmt.Errorf("invalid question ID: %d", answer.QuestionID)
		}

		// Check if answer option is valid for this question
		validOption := false
		for _, option := range question.Options {
			if option.ID == answer.AnswerID {
				validOption = true
				break
			}
		}

		if !validOption {
			return fmt.Errorf("invalid answer ID %d for question %d", answer.AnswerID, answer.QuestionID)
		}

		// Check for duplicate answers
		if answeredQuestions[answer.QuestionID] {
			return fmt.Errorf("duplicate answer for question %d", answer.QuestionID)
		}

		answeredQuestions[answer.QuestionID] = true
	}

	// Check if all questions are answered
	var missingQuestions []string
	for questionID := range questionMap {
		if !answeredQuestions[questionID] {
			missingQuestions = append(missingQuestions, fmt.Sprintf("%d", questionID))
		}
	}

	if len(missingQuestions) > 0 {
		return fmt.Errorf("missing answers for questions: %s", strings.Join(missingQuestions, ", "))
	}

	return nil
}

// calculateScore computes the score based on correct answers
func (s *QuizService) calculateScore(answers []models.QuizAnswer, questionMap map[int]*models.Question) (score, correctAnswers int) {
	for _, answer := range answers {
		question := questionMap[answer.QuestionID]

		// Check if the selected answer is correct
		if answer.AnswerID == question.Answer {
			score++
			correctAnswers++
		}
	}

	return score, correctAnswers
}
