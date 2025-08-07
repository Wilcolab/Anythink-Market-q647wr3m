package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	_ "github.com/lib/pq"
)

func main() {
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
