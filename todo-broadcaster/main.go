package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Todo struct matches the publisher
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Message represents the full payload from publisher
type Message struct {
	Operation string `json:"operation"`
	Todo      Todo   `json:"todo"`
}

// Store for received messages
var (
	messages   []Message
	messagesMu sync.Mutex
)

func main() {
	// Connect to NATS
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to "todos" subject
	_, err = nc.Subscribe("todos", func(msg *nats.Msg) {
		var m Message
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}

		messagesMu.Lock()
		messages = append(messages, m)
		messagesMu.Unlock()

		log.Printf("Received message: %+v\n", m)
	})
	if err != nil {
		log.Fatal(err)
	}

	// HTTP endpoints
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if nc == nil {
			http.Error(w, "NATS not initialized", http.StatusServiceUnavailable)
			return
		}

		if err := nc.FlushTimeout(2 * time.Second); err != nil {
			http.Error(w, "NATS not reachable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		messagesMu.Lock()
		defer messagesMu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(messages); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
