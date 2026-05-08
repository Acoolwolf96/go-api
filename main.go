package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
)

//go:embed static/index.html
var staticFiles embed.FS

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var items = []Item{
	{ID: 1, Name: "apple"},
	{ID: 2, Name: "banana"},
	{ID: 3, Name: "cherry"},
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok", Version: "1.0.0"})
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content, _ := staticFiles.ReadFile("static/index.html")
	w.Write(content)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/items", itemsHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}