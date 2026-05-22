package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leehainuo/coco"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id" example:"1"`
	Name  string `json:"name" example:"Alice"`
	Email string `json:"email" example:"alice@example.com"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"Alice"`
	Email string `json:"email" binding:"required,email" example:"alice@example.com"`
}

// @title Coco Demo API
// @version 1.0
// @description This is a sample API using Gin and Swag
// @host localhost:8000
// @BasePath /api

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/users", listUsers)
		api.POST("/users", createUser)
	}

	// Mount Coco docs - swag will generate openapi.json
	r.Any("/docs/*any", gin.WrapH(coco.New("./docs/swagger.json",
		coco.Title("Coco Demo - Gin + Swag"),
	)))

	println("Server starting on http://localhost:8000")
	println("API docs available at http://localhost:8000/docs/")
	println("Using Gin + Swag")
	println("Run 'swag init' to generate docs")
	r.Run(":8000")
}

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func listUsers(c *gin.Context) {
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}
	c.JSON(http.StatusOK, users)
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
func createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := User{
		ID:    "3",
		Name:  req.Name,
		Email: req.Email,
	}
	c.JSON(http.StatusCreated, user)
}
