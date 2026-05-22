# Framework Integration Guide

Coco integrates seamlessly with any Go web framework. This guide shows how to use Coco with popular frameworks.

## Core Concept

Coco returns a standard `http.Handler`, so it works directly with any framework that supports the standard library interface.

```go
handler := coco.New("./openapi.json")
// handler implements http.Handler interface
```

## Framework Integration Examples

### net/http (Standard Library)

The simplest and most direct way, no adaptation needed:

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

func main() {
    mux := http.NewServeMux()
    
    // Your API routes
    mux.HandleFunc("/api/users", handleUsers)
    mux.HandleFunc("/api/posts", handlePosts)
    
    // Mount documentation
    mux.Handle("/docs/", coco.New("./openapi.json",
        coco.Title("My API Documentation"),
    ))
    
    http.ListenAndServe(":8080", mux)
}
```

**Access**: http://localhost:8080/docs/

---

### Gin

Use `gin.WrapH` to wrap the handler:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/leehainuo/coco"
)

func main() {
    r := gin.Default()
    
    // API routes
    api := r.Group("/api")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
    }
    
    // Mount documentation - use gin.WrapH
    r.Any("/docs/*any", gin.WrapH(coco.New("./openapi.json",
        coco.Title("Gin API Documentation"),
    )))
    
    r.Run(":8080")
}
```

**Key points**:
- Use `gin.WrapH()` to wrap handler
- Route path uses `*any` wildcard
- Use `Any()` method to handle all HTTP methods

**Access**: http://localhost:8080/docs/

---

### Echo

Use `echo.WrapHandler` to wrap the handler:

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/leehainuo/coco"
)

func main() {
    e := echo.New()
    
    // API routes
    e.GET("/api/users", getUsers)
    e.POST("/api/users", createUser)
    
    // Mount documentation - use echo.WrapHandler
    e.Any("/docs/*", echo.WrapHandler(coco.New("./openapi.json",
        coco.Title("Echo API Documentation"),
    )))
    
    e.Start(":8080")
}
```

**Key points**:
- Use `echo.WrapHandler()` to wrap handler
- Route path uses `*` wildcard
- Use `Any()` method to handle all HTTP methods

**Access**: http://localhost:8080/docs/

---

### Fiber v3

Fiber requires an adapter to convert `http.Handler` to Fiber handler:

```go
package main

import (
    "net/http"
    
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/adaptor"
    "github.com/leehainuo/coco"
)

func main() {
    app := fiber.New()
    
    // API routes
    app.Get("/api/users", getUsers)
    app.Post("/api/users", createUser)
    
    // Mount documentation - use adaptor.HTTPHandler
    handler := coco.New("./openapi.json",
        coco.Title("Fiber API Documentation"),
    )
    app.All("/docs/*", adaptor.HTTPHandler(http.HandlerFunc(handler.ServeHTTP)))
    
    app.Listen(":8080")
}
```

**Key points**:
- Import `github.com/gofiber/fiber/v3/middleware/adaptor`
- Use `adaptor.HTTPHandler()` to wrap
- Convert handler to `http.HandlerFunc`

**Access**: http://localhost:8080/docs/

---

### Chi

Chi directly supports standard `http.Handler`:

```go
package main

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/leehainuo/coco"
)

func main() {
    r := chi.NewRouter()
    
    // API routes
    r.Get("/api/users", getUsers)
    r.Post("/api/users", createUser)
    
    // Mount documentation
    r.HandleFunc("/docs/*", func(w http.ResponseWriter, req *http.Request) {
        handler := coco.New("./openapi.json",
            coco.Title("Chi API Documentation"),
        )
        handler.ServeHTTP(w, req)
    })
    
    http.ListenAndServe(":8080", r)
}
```

**Key points**:
- Use `HandleFunc` and call handler internally
- Route path uses `*` wildcard

**Access**: http://localhost:8080/docs/

---

## Integration with OpenAPI Generation Tools

### Huma v2

Huma can automatically generate OpenAPI specs, working perfectly with Coco:

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/leehainuo/coco"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ListUsersOutput struct {
    Body []User
}

func main() {
    mux := http.NewServeMux()
    
    // Create Huma API
    api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))
    
    // Register operation
    huma.Register(api, huma.Operation{
        OperationID: "list-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "List all users",
    }, func(ctx context.Context, input *struct{}) (*ListUsersOutput, error) {
        users := []User{
            {ID: "1", Name: "Alice", Email: "alice@example.com"},
        }
        return &ListUsersOutput{Body: users}, nil
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

**Key points**:
- Use `api.OpenAPI().MarshalJSON()` to get JSON format spec
- Use `coco.Spec()` to pass byte array
- Pass empty string as first parameter

**Huma adapters**:
- `humago` - net/http
- `humagin` - Gin
- `humaecho` - Echo
- `humafiber` - Fiber
- `humachi` - Chi

---

### Swag

Swag generates OpenAPI specs through annotations:

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/leehainuo/coco"
)

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
    mux := http.NewServeMux()
    
    // API routes
    mux.HandleFunc("/api/users", listUsers)
    
    // Mount documentation
    mux.Handle("/docs/", coco.New("./docs/swagger.json",
        coco.Title("My API - Swag"),
    ))
    
    http.ListenAndServe(":8080", mux)
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
    users := []User{
        {ID: "1", Name: "Alice", Email: "alice@example.com"},
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**Generate documentation**:
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs (creates docs/swagger.json)
swag init

# Run the program
go run main.go
```

---

## Path Configuration Comparison

| Framework | Route Definition | Access Path |
|-----------|-----------------|-------------|
| net/http | `/docs/` | http://localhost:8080/docs/ |
| Gin | `/docs/*any` | http://localhost:8080/docs/ |
| Echo | `/docs/*` | http://localhost:8080/docs/ |
| Fiber | `/docs/*` | http://localhost:8080/docs/ |
| Chi | `/docs/*` | http://localhost:8080/docs/ |

**Note**: Path must end with `/` to correctly match sub-paths.

---

## Best Practices

### 1. Unified Documentation Path

Recommend using `/docs/` as documentation path:

```go
// ✅ Recommended
mux.Handle("/docs/", handler)

// ❌ Not recommended
mux.Handle("/api-docs/", handler)
```

### 2. Separate API and Documentation

Use path prefixes to distinguish API and documentation:

```go
// API routes
r.Group("/api", func(r chi.Router) {
    r.Get("/users", getUsers)
})

// Documentation route
r.HandleFunc("/docs/*", docsHandler)
```

### 3. Environment Isolation

Use different configurations for development and production:

```go
func createDocsHandler() http.Handler {
    if os.Getenv("ENV") == "production" {
        return coco.New("./openapi.json",
            coco.EnableDebug(false),
        )
    }
    return coco.New("./openapi.json",
        coco.EnableDebug(true),
    )
}
```

---

## Complete Examples

Check [example/framework](../../example/framework) directory for complete runnable examples.

---

## Related Documentation

- [Quick Start](./01-getting-started.md)
- [Configuration Guide](./02-configuration.md)
- [OpenAPI Generation Guide](./04-openapi-generation.md)
