package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
// @description This is a sample API using Chi and Swag
// @host localhost:8000
// @BasePath /api

func main() {
	r := chi.NewRouter()

	r.Get("/api/users", listUsers)
	r.Post("/api/users", createUser)

	r.HandleFunc("/docs/*", func(w http.ResponseWriter, req *http.Request) {
		handler := coco.New("./docs/swagger.json",
			coco.Title("Coco Demo - Chi + Swag"),
		)
		handler.ServeHTTP(w, req)
	})

	println("Server starting on http://localhost:8000")
	println("API docs available at http://localhost:8000/docs/")
	println("Using Chi + Swag")
	println("Run 'swag init' to generate docs")
	http.ListenAndServe(":8000", r)
}

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Simplified JSON encoding
	w.Write([]byte(`[{"id":"1","name":"Alice","email":"alice@example.com"},{"id":"2","name":"Bob","email":"bob@example.com"}]`))
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id":"3","name":"Alice","email":"alice@example.com"}`))
}
