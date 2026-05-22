package main

import (
	"encoding/json"
	"net/http"

	"github.com/leehainuo/coco"
)

type User struct {
	ID    string `json:"id" example:"1"`
	Name  string `json:"name" example:"Alice"`
	Email string `json:"email" example:"alice@example.com"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required" example:"Alice"`
	Email string `json:"email" validate:"required,email" example:"alice@example.com"`
}

// @title Coco Demo API
// @version 1.0
// @description This is a sample API using net/http and Swag
// @host localhost:8000
// @BasePath /api

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listUsers(w, r)
		case http.MethodPost:
			createUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		handler := coco.New("./docs/swagger.json",
			coco.Title("Coco Demo - net/http + Swag"),
		)
		handler.ServeHTTP(w, r)
	})

	println("Server starting on http://localhost:8000")
	println("API docs available at http://localhost:8000/docs/")
	println("Using net/http + Swag")
	println("Run 'swag init' to generate docs")
	http.ListenAndServe(":8000", mux)
}

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func listUsers(w http.ResponseWriter, _ *http.Request) {
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Create user"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	user := User{
		ID:    "3",
		Name:  req.Name,
		Email: req.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
