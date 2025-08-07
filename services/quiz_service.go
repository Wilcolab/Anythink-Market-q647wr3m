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

// GetAllQuestionsForAPI returns questions in the simplified API format
func (s *QuizService) GetAllQuestionsForAPI() ([]models.QuestionResponse, error) {
	questions, err := s.questionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.convertToQuestionResponses(questions), nil
}

// GetQuestionForAPI returns a single question in the simplified API format
func (s *QuizService) GetQuestionForAPI(id int) (*models.QuestionResponse, error) {
	question, err := s.questionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, nil
	}

	response := s.convertToQuestionResponse(*question)
	return &response, nil
}

// convertToQuestionResponses converts internal Question models to API response format
func (s *QuizService) convertToQuestionResponses(questions []models.Question) []models.QuestionResponse {
	responses := make([]models.QuestionResponse, len(questions))
	for i, q := range questions {
		responses[i] = s.convertToQuestionResponse(q)
	}
	return responses
}

// convertToQuestionResponse converts a single Question to API response format
func (s *QuizService) convertToQuestionResponse(q models.Question) models.QuestionResponse {
	// Create questionId in the format "q{id}"
	questionID := fmt.Sprintf("q%d", q.ID)

	// Extract option texts and find correct answer
	options := make([]string, len(q.Options))
	var correctAnswer string

	for i, option := range q.Options {
		options[i] = option.Text
		if option.ID == q.Answer {
			correctAnswer = option.Text
		}
	}

	return models.QuestionResponse{
		QuestionID:    questionID,
		Text:          q.Question,
		Options:       options,
		CorrectAnswer: correctAnswer,
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

	// Calculate score and generate detailed results
	score, correctAnswers, results := s.calculateScoreWithDetails(submission.Answers, questionMap)

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
		Results:        results,
	}

	return response, nil
}

// SubmitQuizAlternativeFormat processes a quiz submission in the alternative format
func (s *QuizService) SubmitQuizAlternativeFormat(submission *models.AlternativeQuizSubmissionRequest) (*models.QuizSubmissionResponse, error) {
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

	// Create maps for conversion
	questionMap := make(map[int]*models.Question)
	questionIdMap := make(map[string]*models.Question) // "q1" -> Question

	for i := range allQuestions {
		q := &allQuestions[i]
		questionMap[q.ID] = q
		questionIdMap[fmt.Sprintf("q%d", q.ID)] = q
	}

	// Convert alternative format to standard format
	var standardAnswers []models.QuizAnswer
	for questionIdStr, answerText := range submission.Answers {
		// Find the question
		question, exists := questionIdMap[questionIdStr]
		if !exists {
			return nil, fmt.Errorf("invalid question ID: %s", questionIdStr)
		}

		// Find the option ID for the answer text
		var optionID int
		found := false
		for _, option := range question.Options {
			if option.Text == answerText {
				optionID = option.ID
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid answer '%s' for question %s", answerText, questionIdStr)
		}

		standardAnswers = append(standardAnswers, models.QuizAnswer{
			QuestionID: question.ID,
			AnswerID:   optionID,
		})
	}

	// Create standard submission request
	standardSubmission := &models.QuizSubmissionRequest{
		Answers: standardAnswers,
		UserID:  submission.UserID,
	}

	// Process using the standard method
	return s.SubmitQuiz(standardSubmission)
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

// calculateScore computes the score based on correct answers (kept for backward compatibility)
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

// calculateScoreWithDetails computes the score and returns detailed per-question results
func (s *QuizService) calculateScoreWithDetails(answers []models.QuizAnswer, questionMap map[int]*models.Question) (score, correctAnswers int, results []models.QuizAnswerResult) {
	results = make([]models.QuizAnswerResult, 0, len(answers))

	for _, answer := range answers {
		question := questionMap[answer.QuestionID]
		isCorrect := answer.AnswerID == question.Answer

		if isCorrect {
			score++
			correctAnswers++
		}

		// Find the selected option text and correct option text
		var selectedText, correctText string
		for _, option := range question.Options {
			if option.ID == answer.AnswerID {
				selectedText = option.Text
			}
			if option.ID == question.Answer {
				correctText = option.Text
			}
		}

		result := models.QuizAnswerResult{
			QuestionID:   answer.QuestionID,
			Question:     question.Question,
			SelectedID:   answer.AnswerID,
			CorrectID:    question.Answer,
			IsCorrect:    isCorrect,
			SelectedText: selectedText,
			CorrectText:  correctText,
		}

		results = append(results, result)
	}

	return score, correctAnswers, results
}
