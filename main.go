package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {

	// Sample quiz data
	type Option struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}

	type Question struct {
		ID       int      `json:"id"`
		Question string   `json:"question"`
		Options  []Option `json:"options"`
		Answer   int      `json:"answer"` // Option ID
	}

	questions := []Question{
		{ID: 1, Question: "What is a goroutine in Go?", Options: []Option{{ID: 1, Text: "A lightweight thread managed by Go runtime"}, {ID: 2, Text: "A Go function"}, {ID: 3, Text: "A Go package"}, {ID: 4, Text: "A Go variable"}}, Answer: 1},
		{ID: 2, Question: "Which keyword is used to define a new type in Go?", Options: []Option{{ID: 1, Text: "struct"}, {ID: 2, Text: "type"}, {ID: 3, Text: "var"}, {ID: 4, Text: "func"}}, Answer: 2},
		{ID: 3, Question: "How do you declare a variable in Go?", Options: []Option{{ID: 1, Text: "let x int = 5"}, {ID: 2, Text: "x := 5"}, {ID: 3, Text: "var x = 5"}, {ID: 4, Text: "Both 2 and 3"}}, Answer: 4},
		{ID: 4, Question: "What is the zero value of an int in Go?", Options: []Option{{ID: 1, Text: "nil"}, {ID: 2, Text: "0"}, {ID: 3, Text: "undefined"}, {ID: 4, Text: "false"}}, Answer: 2},
		{ID: 5, Question: "Which package is used for formatted I/O in Go?", Options: []Option{{ID: 1, Text: "io"}, {ID: 2, Text: "fmt"}, {ID: 3, Text: "os"}, {ID: 4, Text: "bufio"}}, Answer: 2},
	}

	// CORS middleware
	withCORS := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h(w, r)
		}
	}

	http.HandleFunc("/api/questions", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(questions)
	}))

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
		for _, q := range questions {
			if q.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(q)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Question not found"})
	}))
	// Connect to database
	db := connectDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go API! 🐹")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "ok", "message": "server is running"}
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("🚀 Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to database")
	return db
}
