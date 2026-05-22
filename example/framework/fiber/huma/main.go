package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/leehainuo/coco"
)

type User struct {
	ID    string `json:"id" example:"1"`
	Name  string `json:"name" example:"Alice"`
	Email string `json:"email" example:"alice@example.com"`
}

type CreateUserInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxLength:"50"`
		Email string `json:"email" format:"email"`
	}
}

type CreateUserOutput struct {
	Body User
}

type ListUsersOutput struct {
	Body []User
}

func main() {
	app := fiber.New()

	api := humafiber.New(app, huma.DefaultConfig("Coco Demo API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      "GET",
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
		Method:      "POST",
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

	spec, _ := api.OpenAPI().MarshalJSON()

	handler := coco.New("",
		coco.Spec(spec),
		coco.Title("Coco Demo - Fiber + Huma"),
	)

	app.All("/docs/*", adaptor.HTTPHandler(http.HandlerFunc(handler.ServeHTTP)))

	fmt.Println("Server starting on http://localhost:8000")
	fmt.Println("API docs available at http://localhost:8000/docs/")
	fmt.Println("Using Fiber + Huma")
	app.Listen(":8000")
}
