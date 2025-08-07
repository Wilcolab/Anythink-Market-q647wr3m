package repository

import (
	"database/sql"
	"fmt"
	"go-quiz-api/models"
)

// QuestionRepository handles database operations for questions
type QuestionRepository struct {
	db *sql.DB
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// GetAll retrieves all questions with their options
func (r *QuestionRepository) GetAll() ([]models.Question, error) {
	questions, err := r.getQuestions()
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	for i := range questions {
		options, err := r.getOptionsByQuestionID(questions[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get options for question %d: %w", questions[i].ID, err)
		}
		questions[i].Options = options

		// Set the correct answer ID
		for _, option := range options {
			if option.IsCorrect {
				questions[i].Answer = option.ID
				break
			}
		}
	}

	return questions, nil
}

// GetByID retrieves a single question by ID with its options
func (r *QuestionRepository) GetByID(id int) (*models.Question, error) {
	question := &models.Question{}

	query := `
		SELECT id, question, created_at, updated_at 
		FROM questions 
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&question.ID,
		&question.Question,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	options, err := r.getOptionsByQuestionID(question.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	question.Options = options

	// Set the correct answer ID
	for _, option := range options {
		if option.IsCorrect {
			question.Answer = option.ID
			break
		}
	}

	return question, nil
}

// Create creates a new question with options
func (r *QuestionRepository) Create(question *models.Question) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert question
	query := `
		INSERT INTO questions (question) 
		VALUES ($1) 
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query, question.Question).Scan(
		&question.ID,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	// Insert options
	for i := range question.Options {
		optionQuery := `
			INSERT INTO options (question_id, text, is_correct) 
			VALUES ($1, $2, $3) 
			RETURNING id
		`

		err = tx.QueryRow(optionQuery, question.ID, question.Options[i].Text, question.Options[i].IsCorrect).Scan(&question.Options[i].ID)
		if err != nil {
			return fmt.Errorf("failed to create option: %w", err)
		}
		question.Options[i].QuestionID = question.ID
	}

	return tx.Commit()
}

func (r *QuestionRepository) getQuestions() ([]models.Question, error) {
	query := `
		SELECT id, question, created_at, updated_at 
		FROM questions 
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.ID, &q.Question, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func (r *QuestionRepository) getOptionsByQuestionID(questionID int) ([]models.Option, error) {
	query := `
		SELECT id, question_id, text, is_correct 
		FROM options 
		WHERE question_id = $1 
		ORDER BY id
	`

	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.Option
	for rows.Next() {
		var o models.Option
		err := rows.Scan(&o.ID, &o.QuestionID, &o.Text, &o.IsCorrect)
		if err != nil {
			return nil, err
		}
		options = append(options, o)
	}

	return options, rows.Err()
}
