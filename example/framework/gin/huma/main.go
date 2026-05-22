package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/leehainuo/coco"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id" example:"1" doc:"User ID"`
	Name  string `json:"name" example:"Alice" doc:"User name"`
	Email string `json:"email" example:"alice@example.com" doc:"User email"`
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxLength:"50" doc:"User name"`
		Email string `json:"email" format:"email" doc:"User email"`
	}
}

// CreateUserOutput represents the output for creating a user
type CreateUserOutput struct {
	Body User
}

// ListUsersOutput represents the output for listing users
type ListUsersOutput struct {
	Body []User
}

func main() {
	// Create Gin router
	r := gin.Default()

	// Create Huma API
	api := humagin.New(r, huma.DefaultConfig("Coco Demo API", "1.0.0"))

	// Register operations
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/api/users",
		Summary:     "List all users",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct{}) (*ListUsersOutput, error) {
		users := []User{
			{ID: "1", Name: "Alice", Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Email: "bob@example.com"},
		}
		return &ListUsersOutput{Body: users}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/api/users",
		Summary:     "Create a new user",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
		user := User{
			ID:    "3",
			Name:  input.Body.Name,
			Email: input.Body.Email,
		}
		return &CreateUserOutput{Body: user}, nil
	})

	// Get OpenAPI spec from Huma
	spec, _ := api.OpenAPI().MarshalJSON()

	// Mount Coco docs
	r.Any("/docs/*any", gin.WrapH(coco.New("",
		coco.Spec(spec),
		coco.Title("Coco Demo - Gin + Huma"),
	)))

	fmt.Println("Server starting on http://localhost:8000")
	fmt.Println("API docs available at http://localhost:8000/docs/")
	fmt.Println("Using Gin + Huma")
	r.Run(":8000")
}
