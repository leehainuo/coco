package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leehainuo/coco"
)

func main() {
	mux := http.NewServeMux()

	// Simple API endpoints
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API is healthy",
		})
	})

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users := []map[string]any{
			{"id": "1", "name": "Alice", "email": "alice@example.com"},
			{"id": "2", "name": "Bob", "email": "bob@example.com"},
		}
		json.NewEncoder(w).Encode(users)
	})

	// API documentation
	mux.Handle("/docs/", coco.New("../embedded/openapi.json",
		coco.Title("Coco Demo - Pure net/http"),
		coco.Theme("dark"),
	))

	log.Println("Server starting on http://localhost:8082")
	log.Println("API docs available at http://localhost:8082/docs/")
	log.Println("API endpoints:")
	log.Println("  - GET http://localhost:8082/api/health")
	log.Println("  - GET http://localhost:8082/api/users")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
