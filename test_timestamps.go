package main

import (
	"fmt"
	"go-quiz-api/database"
	"go-quiz-api/repository"
	"log"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	questionRepo := repository.NewQuestionRepository(db.DB)

	// Get all questions to check timestamps
	questions, err := questionRepo.GetAll()
	if err != nil {
		log.Fatal("Failed to get questions:", err)
	}

	fmt.Printf("Found %d questions:\n", len(questions))
	for i, q := range questions {
		fmt.Printf("Question %d:\n", i+1)
		fmt.Printf("  ID: %d\n", q.ID)
		fmt.Printf("  Text: %s\n", q.Question)
		fmt.Printf("  Created At: %s\n", q.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Updated At: %s\n", q.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Options: %d\n", len(q.Options))
		fmt.Println()
	}
}
