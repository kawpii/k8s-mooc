package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT title FROM todos")

	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []string

	// 2. Iterate through the rows
	for rows.Next() {
		var title string
		// 3. Scan into a string variable, not the slice
		if err := rows.Scan(&title); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		todos = append(todos, title)
	}

	w.Header().Set("Content-Type", "application/json")
	// Handle empty case to return [] instead of null
	if todos == nil {
		todos = []string{}
	}
	json.NewEncoder(w).Encode(todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("todo")
	if text == "" || len(text) > 140 {
		log.Println("Todo must be 1–140 characters")
		http.Error(w, "Todo must be 1–140 characters", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
		INSERT INTO todos(title) VALUES ($1)`, text)

	if err != nil {
		log.Printf("Database error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Created todo: %s\n", text)
	// Redirect back to page (prevents duplicate submit on refresh)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	var err error
	// 1. Initialize the DB connection
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Verify connection
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

		if err = db.Ping(); err != nil {
			log.Fatal("Cannot connect to DB: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodosHandler(w, r)
		case http.MethodPost:
			createTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
