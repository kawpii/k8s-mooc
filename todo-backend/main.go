package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, done FROM todos")

	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo

	// 2. Iterate through the rows
	for rows.Next() {
		var t Todo
		// 3. Scan into a string variable, not the slice
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		todos = append(todos, t)
	}

	if todos == nil {
		todos = []Todo{}
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

// UPDATED PUT handler to use /todos/<id>
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// NEW: extract id from URL path
	// Expected format: /todos/123
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[2] == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// NEW: read JSON body instead of query param
	var payload struct {
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE todos SET done=$1 WHERE id=$2", payload.Done, id)
	if err != nil {
		log.Printf("Update error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Updated todo with ID: %s\n", idStr)

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

	http.HandleFunc("/todos/", updateTodoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
