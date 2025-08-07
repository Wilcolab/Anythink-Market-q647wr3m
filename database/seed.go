package database

import (
	"database/sql"
	"go-quiz-api/models"
	"go-quiz-api/repository"
	"log"
)

// SeedData populates the database with initial quiz questions
func SeedData(db *sql.DB) error {
	repo := repository.NewQuestionRepository(db)

	// Check if questions already exist
	questions, err := repo.GetAll()
	if err != nil {
		return err
	}

	if len(questions) > 0 {
		log.Println("📝 Database already contains questions, skipping seed")
		return nil
	}

	log.Println("🌱 Seeding database with initial questions...")

	sampleQuestions := []models.Question{
		{
			Question: "What is a goroutine in Go?",
			Options: []models.Option{
				{Text: "A lightweight thread managed by Go runtime", IsCorrect: true},
				{Text: "A Go function", IsCorrect: false},
				{Text: "A Go package", IsCorrect: false},
				{Text: "A Go variable", IsCorrect: false},
			},
		},
		{
			Question: "Which keyword is used to define a new type in Go?",
			Options: []models.Option{
				{Text: "struct", IsCorrect: false},
				{Text: "type", IsCorrect: true},
				{Text: "var", IsCorrect: false},
				{Text: "func", IsCorrect: false},
			},
		},
		{
			Question: "How do you declare a variable in Go?",
			Options: []models.Option{
				{Text: "let x int = 5", IsCorrect: false},
				{Text: "x := 5", IsCorrect: false},
				{Text: "var x = 5", IsCorrect: false},
				{Text: "Both 2 and 3", IsCorrect: true},
			},
		},
		{
			Question: "What is the zero value of an int in Go?",
			Options: []models.Option{
				{Text: "nil", IsCorrect: false},
				{Text: "0", IsCorrect: true},
				{Text: "undefined", IsCorrect: false},
				{Text: "false", IsCorrect: false},
			},
		},
		{
			Question: "Which package is used for formatted I/O in Go?",
			Options: []models.Option{
				{Text: "io", IsCorrect: false},
				{Text: "fmt", IsCorrect: true},
				{Text: "os", IsCorrect: false},
				{Text: "bufio", IsCorrect: false},
			},
		},
	}

	for _, question := range sampleQuestions {
		if err := repo.Create(&question); err != nil {
			return err
		}
	}

	log.Printf("✅ Successfully seeded %d questions", len(sampleQuestions))
	return nil
}
