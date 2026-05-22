package main

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
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
// @description This is a sample API using Fiber and Swag
// @host localhost:8000
// @BasePath /api

func main() {
	app := fiber.New()

	app.Get("/api/users", listUsers)
	app.Post("/api/users", createUser)

	handler := coco.New("./docs/swagger.json",
		coco.Title("Coco Demo - Fiber + Swag"),
	)
	app.All("/docs/*", adaptor.HTTPHandler(http.HandlerFunc(handler.ServeHTTP)))

	println("Server starting on http://localhost:8000")
	println("API docs available at http://localhost:8000/docs/")
	println("Using Fiber + Swag")
	println("Run 'swag init' to generate docs")
	app.Listen(":8000")
}

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func listUsers(c *fiber.Ctx) error {
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}
	return (*c).JSON(users)
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
func createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := (*c).Bind().Body(&req); err != nil {
		return (*c).Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user := User{
		ID:    "3",
		Name:  req.Name,
		Email: req.Email,
	}
	return (*c).Status(201).JSON(user)
}
