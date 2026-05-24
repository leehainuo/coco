# Coco

<p align="center">
<img align="center" width="160px" src="docs/images/coco.png">
</p>

<p align="center">An elegant, lightweight, and beautiful OpenAPI documentation renderer built for Go developers</p>

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/leehainuo/coco.svg)](https://pkg.go.dev/github.com/leehainuo/coco)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

## Preview

<p align="center">
  <img src="docs/images/coco-light.png" alt="Coco Light Theme" width="45%">
  <img src="docs/images/coco-dark.png" alt="Coco Dark Theme" width="45%">
</p>

## What is Coco?

English | [简体中文](docs/zh/README.md)

Coco is an elegant, lightweight, and beautiful Go library that renders OpenAPI/Swagger specifications into stunning interactive API documentation. Zero dependencies, single-file bundle, perfectly embedded into Go binaries. Supports OpenAPI 2.0 and 3.0+ specifications with seamless integration into any Go web framework.

#### Advantages of Coco:

* **Elegant DX** - Minimal API design, integrate in just 30 seconds
* **Beautiful UI** - Built with Vue 3 + TailwindCSS, supports automatic light/dark theme switching
* **Framework Agnostic** - Works with all Go web frameworks (Gin, Echo, Fiber, Chi, net/http, etc.)
* **Feature Complete** - Built-in API testing, spec export, request history, internationalization, and more

## Core Features

* **Internationalization** - Built-in English and Chinese support, extensible for more languages
* **Zero Dependencies** - Pure Go implementation, frontend assets fully embedded, no external tools required
* **Framework Agnostic** - Compatible with all Go web frameworks (Gin, Echo, Fiber, Chi, net/http, etc.)
* **Flexible Configuration** - Rich configuration options to meet various customization needs
* **API Testing** - Built-in interactive debug panel for instant API testing
* **Spec Export** - One-click export of OpenAPI/Swagger specification files
* **Request History** - Automatically saves debug history for easy review and reuse
* **High Performance** - Embedded static assets, single-file bundle, no additional HTTP requests

## Installation

```bash
go get github.com/leehainuo/coco
```

## Quick Start

### Using Swag (Recommended)

**Step 1: Install Swag**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Step 2: Add annotations to your code**

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

// @title           My API
// @version         1.0
// @description     This is a sample API
// @host            localhost:8000
// @BasePath        /api

func main() {
    mux := http.NewServeMux()
    
    // Your API routes
    mux.HandleFunc("/api/hello", handleHello)
    
    // Mount Coco docs (Swag generates docs/swagger.json)
    mux.Handle("/docs/", coco.New("./docs/swagger.json"))
    
    http.ListenAndServe(":8000", mux)
}

// @Summary      Hello endpoint
// @Description  Returns a greeting message
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "Hello, World!"
// @Router       /hello [get]
func handleHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}
```

**Step 3: Generate docs and run**

```bash
# Generate OpenAPI documentation
swag init

# Run your application
go run main.go
```

**Step 4: View your docs**

Open your browser and visit http://localhost:8000/docs/ to see your beautiful API documentation!

---

### Other Usage Methods

**Load from local file**
```go
handler := coco.New("./openapi.json")
```

**Load from URL**
```go
handler := coco.New("", coco.SpecURL("https://example.com/openapi.json"))
```

**Load from byte array**
```go
handler := coco.New("", coco.Spec(specBytes))
```

**For more configuration and integration options, see [Full Documentation](./docs/en/)**

## Integration Examples

### net/http

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
    
    // Mount documentation
    mux.Handle("/docs/", coco.New("./openapi.json",
        coco.Title("User API"),
    ))
    
    http.ListenAndServe(":8000", mux)
}
```

### Gin

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/leehainuo/coco"
)

func main() {
    r := gin.Default()
    
    // Your API routes
    r.GET("/api/users", getUsers)
    
    // Mount documentation
    r.Any("/docs/*any", gin.WrapH(coco.New("./openapi.json",
        coco.Title("User API"),
    )))
    
    r.Run(":8000")
}
```

### Echo

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/leehainuo/coco"
)

func main() {
    e := echo.New()
    
    // Your API routes
    e.GET("/api/users", getUsers)
    
    // Mount documentation
    e.Any("/docs/*", echo.WrapHandler(coco.New("./openapi.json",
        coco.Title("User API"),
    )))
    
    e.Start(":8000")
}
```

### Fiber

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
    
    // Your API routes
    app.Get("/api/users", getUsers)
    
    // Mount documentation
    handler := coco.New("./openapi.json",
        coco.Title("User API"),
    )
    app.All("/docs/*", adaptor.HTTPHandler(http.HandlerFunc(handler.ServeHTTP)))
    
    app.Listen(":8000")
}
```

### Chi

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/leehainuo/coco"
)

func main() {
    r := chi.NewRouter()
    
    // Your API routes
    r.Get("/api/users", getUsers)
    
    // Mount documentation
    r.HandleFunc("/docs/*", func(w http.ResponseWriter, req *http.Request) {
        handler := coco.New("./openapi.json",
            coco.Title("User API"),
        )
        handler.ServeHTTP(w, req)
    })
    
    http.ListenAndServe(":8000", r)
}
```

## Configuration Options

### Spec Sources

```go
// Load from file path
coco.New("./openapi.json")

// Load from byte array
coco.New("", coco.Spec(specBytes))

// Load from remote URL
coco.New("", coco.SpecURL("https://example.com/openapi.json"))
```

### UI Configuration

```go
// Set document title
coco.Title("My API Documentation")

// Set theme: "light", "dark", "auto"
coco.Theme("dark")

// Set language: "en", "zh"
coco.Lang("en")
```

### Feature Toggles

```go
// Enable/disable debug panel (enabled by default)
coco.EnableDebug(true)

// Enable/disable export feature (enabled by default)
coco.EnableExport(true)

// Enable/disable history (enabled by default)
coco.EnableHistory(true)
```

## OpenAPI Generation Tools Integration

### Using Huma v2

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/leehainuo/coco"
)

func main() {
    mux := http.NewServeMux()
    api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))
    
    // Register your API
    huma.Register(api, huma.Operation{
        OperationID: "get-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "Get users",
    }, func(ctx context.Context, input *struct{}) (*struct{}, error) {
        return &struct{}{}, nil
    })
    
    // Get OpenAPI spec and mount documentation
    spec, _ := api.OpenAPI().MarshalJSON()
    mux.Handle("/docs/", coco.New("",
        coco.Spec(spec),
        coco.Title("My API - Huma"),
    ))
    
    http.ListenAndServe(":8000", mux)
}
```

### Using Swag

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

// @title My API
// @version 1.0
// @description This is my API
// @host localhost:8000
// @BasePath /api

func main() {
    mux := http.NewServeMux()
    
    // Your API routes
    mux.HandleFunc("/api/users", getUsers)
    
    // Mount documentation (swag generates docs/swagger.json)
    mux.Handle("/docs/", coco.New("./docs/swagger.json",
        coco.Title("My API - Swag"),
    ))
    
    http.ListenAndServe(":8000", mux)
}
```

Generate documentation before running:
```bash
swag init
```

## Themes and Languages

### Theme Options

- `light` - Light theme
- `dark` - Dark theme
- `auto` - Follow system (default)

### Language Options

- `en` - English (default)
- `zh` - Chinese

Users can switch themes and languages anytime in the top-right corner of the interface.

## Documentation

- **English Docs**: [docs/en/](./docs/en/) - Complete English documentation
- **中文文档**: [docs/zh/](./docs/zh/) - 完整的中文使用指南
- **Documentation Home**: [docs/](./docs/) - Choose your language / 选择你的语言

## Complete Examples

Check the [example/framework](./example/framework) directory for complete examples:

- **net/http** - Standard library examples
- **Gin** - Gin framework integration
- **Echo** - Echo framework integration
- **Fiber** - Fiber v3 framework integration
- **Chi** - Chi router integration

Each framework provides examples for both Huma and Swag OpenAPI generation methods.

## 📚 Documentation

- **English**: [docs/en/](./docs/en/)
- **中文**: [docs/zh/](./docs/zh/)

## 🤝 Contributing

Contributions, issues, and suggestions are welcome!

If you're using Coco or find it helpful, please give us a Star ⭐

## License

MIT License - See [LICENSE](LICENSE) file for details

## Give a Star! ⭐

If you like this project or are using it to learn or build your solution, please give it a Star to get updates on new releases. Your support matters!

---

**Made with ❤️ by [leehainuo](https://github.com/leehainuo)**
