# 框架集成指南

Coco 可以与任何 Go Web 框架无缝集成。本指南展示了如何在主流框架中使用 Coco。

## 核心概念

Coco 返回一个标准的 `http.Handler`，因此可以直接用于任何支持标准库接口的框架。

```go
handler := coco.New("./openapi.json")
// handler 实现了 http.Handler 接口
```

## 框架集成示例

### net/http (标准库)

最简单直接的方式，无需任何适配：

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

func main() {
    mux := http.NewServeMux()
    
    // 你的 API 路由
    mux.HandleFunc("/api/users", handleUsers)
    mux.HandleFunc("/api/posts", handlePosts)
    
    // 挂载文档
    mux.Handle("/docs/", coco.New("./openapi.json",
        coco.Title("My API Documentation"),
    ))
    
    http.ListenAndServe(":8080", mux)
}
```

**访问**: http://localhost:8080/docs/

---

### Gin

使用 `gin.WrapH` 包装 handler：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/leehainuo/coco"
)

func main() {
    r := gin.Default()
    
    // API 路由
    api := r.Group("/api")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
    }
    
    // 挂载文档 - 使用 gin.WrapH
    r.Any("/docs/*any", gin.WrapH(coco.New("./openapi.json",
        coco.Title("Gin API Documentation"),
    )))
    
    r.Run(":8080")
}
```

**要点**:
- 使用 `gin.WrapH()` 包装 handler
- 路由路径使用 `*any` 通配符
- 使用 `Any()` 方法处理所有 HTTP 方法

**访问**: http://localhost:8080/docs/

---

### Echo

使用 `echo.WrapHandler` 包装 handler：

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/leehainuo/coco"
)

func main() {
    e := echo.New()
    
    // API 路由
    e.GET("/api/users", getUsers)
    e.POST("/api/users", createUser)
    
    // 挂载文档 - 使用 echo.WrapHandler
    e.Any("/docs/*", echo.WrapHandler(coco.New("./openapi.json",
        coco.Title("Echo API Documentation"),
    )))
    
    e.Start(":8080")
}
```

**要点**:
- 使用 `echo.WrapHandler()` 包装 handler
- 路由路径使用 `*` 通配符
- 使用 `Any()` 方法处理所有 HTTP 方法

**访问**: http://localhost:8080/docs/

---

### Fiber v3

Fiber 需要使用适配器将 `http.Handler` 转换为 Fiber handler：

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
    
    // API 路由
    app.Get("/api/users", getUsers)
    app.Post("/api/users", createUser)
    
    // 挂载文档 - 使用 adaptor.HTTPHandler
    handler := coco.New("./openapi.json",
        coco.Title("Fiber API Documentation"),
    )
    app.All("/docs/*", adaptor.HTTPHandler(http.HandlerFunc(handler.ServeHTTP)))
    
    app.Listen(":8080")
}
```

**要点**:
- 需要导入 `github.com/gofiber/fiber/v3/middleware/adaptor`
- 使用 `adaptor.HTTPHandler()` 包装
- 需要将 handler 转换为 `http.HandlerFunc`

**访问**: http://localhost:8080/docs/

---

### Chi

Chi 直接支持标准 `http.Handler`：

```go
package main

import (
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/leehainuo/coco"
)

func main() {
    r := chi.NewRouter()
    
    // API 路由
    r.Get("/api/users", getUsers)
    r.Post("/api/users", createUser)
    
    // 挂载文档
    r.HandleFunc("/docs/*", func(w http.ResponseWriter, req *http.Request) {
        handler := coco.New("./openapi.json",
            coco.Title("Chi API Documentation"),
        )
        handler.ServeHTTP(w, req)
    })
    
    http.ListenAndServe(":8080", r)
}
```

**要点**:
- 使用 `HandleFunc` 并在内部调用 handler
- 路由路径使用 `*` 通配符

**访问**: http://localhost:8080/docs/

---

## 与 OpenAPI 生成工具集成

### Huma v2

Huma 可以自动生成 OpenAPI 规范，与 Coco 完美配合：

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
    
    // 创建 Huma API
    api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))
    
    // 注册操作
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
    
    // 获取 OpenAPI 规范（JSON 格式）
    spec, _ := api.OpenAPI().MarshalJSON()
    
    // 挂载 Coco 文档
    mux.Handle("/docs/", coco.New("",
        coco.Spec(spec),
        coco.Title("My API - Huma"),
    ))
    
    http.ListenAndServe(":8080", mux)
}
```

**要点**:
- 使用 `api.OpenAPI().MarshalJSON()` 获取 JSON 格式规范
- 使用 `coco.Spec()` 传递字节数组
- 第一个参数传空字符串

**Huma 适配器**:
- `humago` - net/http
- `humagin` - Gin
- `humaecho` - Echo
- `humafiber` - Fiber
- `humachi` - Chi

---

### Swag

Swag 通过注释生成 OpenAPI 规范：

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/leehainuo/coco"
)

// @title My API
// @version 1.0
// @description This is my API
// @host localhost:8080
// @BasePath /api

func main() {
    mux := http.NewServeMux()
    
    // API 路由
    mux.HandleFunc("/api/users", listUsers)
    
    // 挂载文档
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

**生成文档**:
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档（会创建 docs/swagger.json）
swag init

# 运行程序
go run main.go
```

---

## 路径配置对比

| 框架 | 路由定义 | 访问路径 |
|------|---------|---------|
| net/http | `/docs/` | http://localhost:8080/docs/ |
| Gin | `/docs/*any` | http://localhost:8080/docs/ |
| Echo | `/docs/*` | http://localhost:8080/docs/ |
| Fiber | `/docs/*` | http://localhost:8080/docs/ |
| Chi | `/docs/*` | http://localhost:8080/docs/ |

**注意**: 路径必须以 `/` 结尾才能正确匹配子路径。

---

## 最佳实践

### 1. 统一的文档路径

建议使用 `/docs/` 作为文档路径，这是业界常见的约定：

```go
// ✅ 推荐
mux.Handle("/docs/", handler)

// ❌ 不推荐
mux.Handle("/api-docs/", handler)
mux.Handle("/swagger/", handler)
```

### 2. 分离 API 和文档

使用路径前缀区分 API 和文档：

```go
// API 路由
r.Group("/api", func(r chi.Router) {
    r.Get("/users", getUsers)
    r.Post("/users", createUser)
})

// 文档路由
r.HandleFunc("/docs/*", docsHandler)
```

### 3. 环境隔离

开发和生产环境使用不同配置：

```go
func createDocsHandler() http.Handler {
    if os.Getenv("ENV") == "production" {
        return coco.New("./openapi.json",
            coco.EnableDebug(false),
        )
    }
    return coco.New("./openapi.json",
        coco.EnableDebug(true),
        coco.Lang("zh"),
    )
}
```

### 4. 中间件集成

可以为文档路由添加认证等中间件：

```go
// Gin 示例
authorized := r.Group("/docs")
authorized.Use(authMiddleware())
authorized.Any("/*any", gin.WrapH(handler))
```

---

## 完整示例

查看 [example/framework](../example/framework) 目录获取每个框架的完整可运行示例：

- `nethttp/` - net/http 标准库
- `gin/` - Gin 框架
- `echo/` - Echo 框架
- `fiber/` - Fiber v3 框架
- `chi/` - Chi 路由器

每个目录都包含 `huma/` 和 `swag/` 两种 OpenAPI 生成方式的示例。

---

## 相关文档

- [快速入门](./01-getting-started.md)
- [配置指南](./02-configuration.md)
- [OpenAPI 生成指南](./04-openapi-generation.md)
