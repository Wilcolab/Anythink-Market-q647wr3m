package main

import (
	"encoding/json"
	"fmt"
	"go-quiz-api/database"
	"go-quiz-api/models"
	"go-quiz-api/repository"
	"go-quiz-api/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Seed database with initial data
	if err := database.SeedData(db.DB); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	// Initialize repository
	questionRepo := repository.NewQuestionRepository(db.DB)

	// Initialize services
	quizService := services.NewQuizService(questionRepo)

	// CORS middleware
	withCORS := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h(w, r)
		}
	}

	// Get all questions endpoint
	http.HandleFunc("/api/questions", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		questions, err := quizService.GetAllQuestionsForAPI()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch questions"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(questions)
	}))

	// Get single question endpoint
	http.HandleFunc("/api/questions/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/questions/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid question ID"})
			return
		}

		question, err := quizService.GetQuestionForAPI(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch question"})
			return
		}

		if question == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Question not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(question)
	}))

	// Quiz submission endpoint
	http.HandleFunc("/api/quiz/submit", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		// Read the request body
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
			return
		}

		// Detect format based on the structure of "answers"
		var result *models.QuizSubmissionResponse
		var err error

		if answers, ok := body["answers"]; ok {
			switch answers.(type) {
			case []interface{}:
				// Standard format: answers is an array
				bodyBytes, _ := json.Marshal(body)
				var submission models.QuizSubmissionRequest
				if err := json.Unmarshal(bodyBytes, &submission); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid format for standard submission"})
					return
				}
				result, err = quizService.SubmitQuiz(&submission)

			case map[string]interface{}:
				// Alternative format: answers is an object/map
				bodyBytes, _ := json.Marshal(body)
				var submission models.AlternativeQuizSubmissionRequest
				if err := json.Unmarshal(bodyBytes, &submission); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid format for alternative submission"})
					return
				}
				result, err = quizService.SubmitQuizAlternativeFormat(&submission)

			default:
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid answers format"})
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing answers field"})
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}))

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "ok", "message": "server is running"}
		json.NewEncoder(w).Encode(response)
	})

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Quiz API with PostgreSQL! 🐹🐘")
	})

	fmt.Println("🚀 Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
