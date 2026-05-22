# OpenAPI Specification Generation Guide

This guide introduces how to generate OpenAPI specification files for use with Coco.

## Two Main Approaches

### 1. Code-First - Huma

Define APIs through code, automatically generate OpenAPI specs.

**Pros**:
- ✅ Type-safe
- ✅ Automatic validation
- ✅ Specs stay in sync with code
- ✅ No manual documentation maintenance

**Cons**:
- ❌ Steeper learning curve
- ❌ Requires framework adapters

### 2. Annotation-First - Swag

Generate OpenAPI specs through code annotations.

**Pros**:
- ✅ Simple and easy to use
- ✅ Easy integration with existing code
- ✅ Flexible documentation control

**Cons**:
- ❌ Need to manually keep annotations in sync with code
- ❌ Lacks compile-time checking
- ❌ Requires additional generation step

---

## Using Huma v2

### Installation

```bash
go get github.com/danielgtaylor/huma/v2
```

Choose adapter based on your framework:

```bash
# net/http
go get github.com/danielgtaylor/huma/v2/adapters/humago

# Gin
go get github.com/danielgtaylor/huma/v2/adapters/humagin

# Echo
go get github.com/danielgtaylor/huma/v2/adapters/humaecho

# Fiber
go get github.com/danielgtaylor/huma/v2/adapters/humafiber

# Chi
go get github.com/danielgtaylor/huma/v2/adapters/humachi
```

### Basic Usage

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/leehainuo/coco"
)

// Define data models
type User struct {
    ID    string `json:"id" example:"1" doc:"User ID"`
    Name  string `json:"name" example:"Alice" doc:"User name"`
    Email string `json:"email" example:"alice@example.com" doc:"Email"`
}

// Define input
type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" minLength:"1" maxLength:"50" doc:"User name"`
        Email string `json:"email" format:"email" doc:"Email address"`
    }
}

// Define output
type CreateUserOutput struct {
    Body User
}

type ListUsersOutput struct {
    Body []User
}

func main() {
    mux := http.NewServeMux()
    
    // Create Huma API
    api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))
    
    // Register GET operation
    huma.Register(api, huma.Operation{
        OperationID: "list-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "Get user list",
        Description: "Return all users",
        Tags:        []string{"Users"},
    }, func(ctx context.Context, input *struct{}) (*ListUsersOutput, error) {
        users := []User{
            {ID: "1", Name: "Alice", Email: "alice@example.com"},
            {ID: "2", Name: "Bob", Email: "bob@example.com"},
        }
        return &ListUsersOutput{Body: users}, nil
    })
    
    // Register POST operation
    huma.Register(api, huma.Operation{
        OperationID: "create-user",
        Method:      http.MethodPost,
        Path:        "/api/users",
        Summary:     "Create user",
        Description: "Create a new user",
        Tags:        []string{"Users"},
    }, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
        user := User{
            ID:    "3",
            Name:  input.Body.Name,
            Email: input.Body.Email,
        }
        return &CreateUserOutput{Body: user}, nil
    })
    
    // Get OpenAPI spec (JSON format)
    spec, _ := api.OpenAPI().MarshalJSON()
    
    // Mount Coco documentation
    mux.Handle("/docs/", coco.New("",
        coco.Spec(spec),
        coco.Title("My API - Huma"),
    ))
    
    http.ListenAndServe(":8080", mux)
}
```

### Advanced Features

#### Path Parameters

```go
type GetUserInput struct {
    UserID string `path:"id" doc:"User ID"`
}

type GetUserOutput struct {
    Body User
}

huma.Register(api, huma.Operation{
    OperationID: "get-user",
    Method:      http.MethodGet,
    Path:        "/api/users/{id}",
    Summary:     "Get single user",
}, func(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
    user := User{ID: input.UserID, Name: "Alice", Email: "alice@example.com"}
    return &GetUserOutput{Body: user}, nil
})
```

#### Query Parameters

```go
type ListUsersInput struct {
    Page  int    `query:"page" default:"1" doc:"Page number"`
    Limit int    `query:"limit" default:"10" doc:"Items per page"`
    Sort  string `query:"sort" enum:"name,email,created" doc:"Sort field"`
}
```

#### Request Headers

```go
type AuthInput struct {
    Authorization string `header:"Authorization" doc:"Auth token"`
}
```

#### Validation Rules

```go
type CreateUserInput struct {
    Body struct {
        Name     string `json:"name" minLength:"2" maxLength:"50"`
        Email    string `json:"email" format:"email"`
        Age      int    `json:"age" minimum:"18" maximum:"120"`
        Website  string `json:"website" format:"uri"`
        Password string `json:"password" pattern:"^[a-zA-Z0-9]{8,}$"`
    }
}
```

---

## Using Swag

### Installation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Basic Usage

#### 1. Add General Annotations

Add at the top of `main.go`:

```go
package main

// @title My API
// @version 1.0
// @description This is my API server
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
    // ...
}
```

#### 2. Add Annotations for Each API

```go
// listUsers godoc
// @Summary Get user list
// @Description Return all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func listUsers(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// createUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User info"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func createUser(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// getUser godoc
// @Summary Get single user
// @Description Get user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func getUser(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

#### 3. Generate Documentation

```bash
swag init
```

This generates in `docs/` directory:
- `docs.go`
- `swagger.json`
- `swagger.yaml`

#### 4. Integrate with Coco

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

func main() {
    mux := http.NewServeMux()
    
    // API routes
    mux.HandleFunc("/api/users", listUsers)
    
    // Mount documentation
    mux.Handle("/docs/", coco.New("./docs/swagger.json",
        coco.Title("My API - Swag"),
    ))
    
    http.ListenAndServe(":8080", mux)
}
```

### Swag Annotation Reference

#### General Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@title` | API title | `@title My API` |
| `@version` | API version | `@version 1.0` |
| `@description` | API description | `@description This is my API` |
| `@host` | Host address | `@host localhost:8080` |
| `@BasePath` | Base path | `@BasePath /api` |

#### Operation Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description | `@Summary Get users` |
| `@Description` | Detailed description | `@Description Get all users` |
| `@Tags` | Tag grouping | `@Tags Users` |
| `@Accept` | Accept content type | `@Accept json` |
| `@Produce` | Produce content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path string true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} User` |
| `@Failure` | Failure response | `@Failure 404 {object} Error` |
| `@Router` | Route path | `@Router /users [get]` |
| `@Security` | Security auth | `@Security ApiKeyAuth` |

#### Parameter Types

```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter
// @Param page query int false "Page number"

// Header
// @Param Authorization header string true "Bearer token"

// Request body
// @Param user body CreateUserRequest true "User info"

// Form parameter
// @Param name formData string true "User name"
```

---

## Comparison Summary

| Feature | Huma | Swag |
|---------|------|------|
| **Learning Curve** | Medium | Easy |
| **Type Safety** | ✅ Strong typing | ❌ Depends on annotations |
| **Auto Validation** | ✅ Built-in | ❌ Manual |
| **Code Generation** | ❌ Not needed | ✅ Need `swag init` |
| **Maintenance Cost** | Low | Medium (keep annotations in sync) |
| **Flexibility** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Framework Support** | Needs adapters | Framework agnostic |
| **Use Case** | New projects, type safety | Existing projects, quick integration |

---

## Best Practices

### Huma Best Practices

1. **Use meaningful OperationID**
   ```go
   OperationID: "list-users"  // ✅ Good
   OperationID: "op1"         // ❌ Bad
   ```

2. **Add detailed documentation**
   ```go
   type User struct {
       ID   string `json:"id" doc:"User unique identifier"`
       Name string `json:"name" doc:"User display name"`
   }
   ```

3. **Use tag grouping**
   ```go
   Tags: []string{"Users", "Management"}
   ```

4. **Define error responses**
   ```go
   type ErrorOutput struct {
       Body struct {
           Error string `json:"error" doc:"Error message"`
       }
   }
   ```

### Swag Best Practices

1. **Keep annotations in sync with code**
   - Update annotations immediately after modifying API
   - Run `swag init` regularly to validate

2. **Use consistent naming**
   ```go
   // ✅ Good
   // @Summary Get user
   // @Router /users/{id} [get]
   
   // ❌ Bad
   // @Summary getUserById
   // @Router /user/{userId} [get]
   ```

3. **Add example values**
   ```go
   type User struct {
       ID   string `json:"id" example:"123"`
       Name string `json:"name" example:"Alice"`
   }
   ```

4. **Use CI/CD for automatic generation**
   ```yaml
   # .github/workflows/docs.yml
   - name: Generate docs
     run: swag init
   ```

---

## Related Resources

- [Huma Official Documentation](https://huma.rocks/)
- [Swag Official Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Complete Examples](../../example/framework/)

---

## Related Documentation

- [Quick Start](./01-getting-started.md)
- [Framework Integration](./03-framework-integration.md)
- [Configuration Guide](./02-configuration.md)
