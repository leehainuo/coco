package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
// @description This is a sample API using Echo and Swag
// @host localhost:8000
// @BasePath /api

func main() {
	e := echo.New()

	e.GET("/api/users", listUsers)
	e.POST("/api/users", createUser)

	e.Any("/docs/*", echo.WrapHandler(coco.New("./docs/swagger.json",
		coco.Title("Coco Demo - Echo + Swag"),
	)))

	println("Server starting on http://localhost:8000")
	println("API docs available at http://localhost:8000/docs/")
	println("Using Echo + Swag")
	println("Run 'swag init' to generate docs")
	e.Start(":8000")
}

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func listUsers(c echo.Context) error {
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}
	return c.JSON(http.StatusOK, users)
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
func createUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := User{
		ID:    "3",
		Name:  req.Name,
		Email: req.Email,
	}
	return c.JSON(http.StatusCreated, user)
}
