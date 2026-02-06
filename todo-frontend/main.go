package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

const (
	imageURL  = "https://picsum.photos/400/400"
	cacheDir  = "cache"
	imageName = "picsum.jpg"
	interval  = 10 * time.Minute
)

var (
	imageBytes []byte
	mu         sync.RWMutex
)

// download and replace cache
func refreshImage() {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Println("Image download failed:", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Image read failed:", err)
		return
	}

	// Save to disk
	os.MkdirAll(cacheDir, 0755)
	path := filepath.Join(cacheDir, imageName)
	os.WriteFile(path, data, 0644)

	// Save to memory
	mu.Lock()
	imageBytes = data
	mu.Unlock()

	log.Println("Image cache refreshed")
}

// background scheduler
func startImageRefresher() {
	// Initial download
	refreshImage()

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			refreshImage()
		}
	}()
}

// serve cached image
func imageHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	data := imageBytes
	mu.RUnlock()

	if data == nil {
		http.Error(w, "Image not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=600") // 10 minutes
	w.Write(data)
}

// serve HTML page
func pageHandler(w http.ResponseWriter, r *http.Request) {

	url := os.Getenv("BACKEND_URL")
	resp, err := http.Get(url + "/todos")
	if err != nil {
		http.Error(w, "Failed to fetch data:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var todos []string
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		http.Error(w, "Invalid response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	var pageTmpl = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Image + Text Submit</title>
	<style>
		body {
			font-family: sans-serif;
			max-width: 600px;
			margin: 40px auto;
		}
		img {
			max-width: 100%;
			border-radius: 8px;
			margin-bottom: 20px;
		}
		.counter {
			font-size: 0.9em;
			color: #555;
		}
	</style>
</head>
<body>
<h2>Cached Image (updates every 10 minutes)</h2>
<img src="/image" alt="Cached Picsum Image">
<h3>Enter text (max 140 characters)</h3>
<form method="POST" action="/todos">
	<textarea
		name="todo"
		rows="4"
		style="width: 100%;"
		maxlength="140"
		oninput="updateCounter(this)"
		required
	></textarea>

	<div class="counter" id="counter">
		140 characters remaining
	</div>

	<br>
	<button type="submit">Submit</button>
</form>

<h2>Todo List</h2>

<ul>
{{range .}}
	<li>{{.}}</li>
{{else}}
	<li>No todos found</li>
{{end}}
</ul>

</body>
</html>
`))

	w.Header().Set("Content-Type", "text/html")
	if err := pageTmpl.Execute(w, todos); err != nil {
		http.Error(w, "Template error", 500)
	}
}

// func submitHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	msg := r.FormValue("message")

// 	if utf8.RuneCountInString(msg) > 140 {
// 		http.Error(w, "Message exceeds 140 characters", http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// 	w.Write([]byte("<h3>Message received:</h3><p>" +
// 		template.HTMLEscapeString(msg) + "</p><a href='/'>Back</a>"))
// }

func main() {
	startImageRefresher()

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/image", imageHandler)
	// http.HandleFunc("/submit", submitHandler)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
