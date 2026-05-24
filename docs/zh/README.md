# Coco

<p align="center">
<img align="center" width="160px" src="../images/coco.png">
</p>

<p align="center">一个优雅、轻量、美观的 OpenAPI 文档渲染库，专为 Go 开发者打造</p>

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/leehainuo/coco.svg)](https://pkg.go.dev/github.com/leehainuo/coco)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

## 效果展示

<p align="center">
  <img src="../images/coco-light.png" alt="Coco 亮色主题" width="45%">
  <img src="../images/coco-dark.png" alt="Coco 暗色主题" width="45%">
</p>

## 什么是 Coco？

[English](../../README.md) | 简体中文

Coco 是一个**优雅、轻量、美观**的 Go 库，用于将 OpenAPI/Swagger 规范渲染为精美的交互式 API 文档。零依赖设计，单文件打包，完美嵌入 Go 二进制文件。支持 OpenAPI 2.0 和 3.0+ 规范，可无缝集成到任何 Go Web 框架。


#### Coco 的优势：

* **优雅的DX** - 极简的 API 设计，最快 30 秒完成集成
* **精美的UI** - 基于 Vue 3 + TailwindCSS 构建，支持亮/暗主题自动切换
* **框架无关** - 所有 Go Web 框架（Gin、Echo、Fiber、Chi、net/http 等）可用
* **功能齐全** - 内置 API 调试、规范导出、请求历史、国际化等功能

## 核心特性

* **国际化支持** - 内置中英文，可扩展更多语言
* **零依赖** - 纯 Go 实现，前端资源完全嵌入，无需任何外部工具
* **框架无关** - 兼容所有 Go Web 框架（Gin、Echo、Fiber、Chi、net/http 等）
* **灵活配置** - 丰富的配置选项，满足各种定制需求
* **API 调试** - 内置交互式调试面板，即时测试 API 接口
* **规范导出** - 一键导出 OpenAPI/Swagger 规范文件
* **请求历史** - 自动保存调试历史，方便回顾和重用
* **高性能** - 嵌入式静态资源，单文件打包，无额外 HTTP 请求

## 安装

```bash
go get github.com/leehainuo/coco
```

## 快速开始

### 使用 Swag 生成文档（推荐）

**第一步：安装 Swag**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**第二步：在代码中添加注释**

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

// @title           我的 API
// @version         1.0
// @description     这是一个示例 API
// @host            localhost:8000
// @BasePath        /api

func main() {
    mux := http.NewServeMux()
    
    // 你的 API 路由
    mux.HandleFunc("/api/hello", handleHello)
    
    // 挂载 Coco 文档（Swag 会生成 docs/swagger.json）
    mux.Handle("/docs/", coco.New("./docs/swagger.json"))
    
    http.ListenAndServe(":8000", mux)
}

// @Summary      问候接口
// @Description  返回问候消息
// @Tags         示例
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "Hello, World!"
// @Router       /hello [get]
func handleHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}
```

**第三步：生成文档并运行**

```bash
# 生成 OpenAPI 文档
swag init

# 运行程序
go run main.go
```

**第四步：访问文档**

打开浏览器访问 http://localhost:8000/docs/ 即可看到精美的 API 文档！

---

### 其他使用方式

**从本地文件加载**
```go
handler := coco.New("./openapi.json")
```

**从 URL 加载**
```go
handler := coco.New("", coco.SpecURL("https://example.com/openapi.json"))
```

**从字节数组加载**
```go
handler := coco.New("", coco.Spec(specBytes))
```

**更多配置和集成方式请查看下方的完整文档**

## 集成示例

### net/http

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
    
    // 挂载文档
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
    
    // 你的 API 路由
    r.GET("/api/users", getUsers)
    
    // 挂载文档
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
    
    // 你的 API 路由
    e.GET("/api/users", getUsers)
    
    // 挂载文档
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
    
    // 你的 API 路由
    app.Get("/api/users", getUsers)
    
    // 挂载文档
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
    
    // 你的 API 路由
    r.Get("/api/users", getUsers)
    
    // 挂载文档
    r.HandleFunc("/docs/*", func(w http.ResponseWriter, req *http.Request) {
        handler := coco.New("./openapi.json",
            coco.Title("User API"),
        )
        handler.ServeHTTP(w, req)
    })
    
    http.ListenAndServe(":8000", r)
}
```

## 配置选项

### 规范来源

```go
// 从文件路径加载
coco.New("./openapi.json")

// 从字节数组加载
coco.New("", coco.Spec(specBytes))

// 从远程 URL 加载
coco.New("", coco.SpecURL("https://example.com/openapi.json"))
```

### UI 配置

```go
// 设置文档标题
coco.Title("我的 API 文档")

// 设置主题: "light", "dark", "auto"
coco.Theme("dark")

// 设置语言: "en", "zh"
coco.Lang("zh")
```

### 功能开关

```go
// 启用/禁用调试面板（默认启用）
coco.EnableDebug(true)

// 启用/禁用导出功能（默认启用）
coco.EnableExport(true)

// 启用/禁用历史记录（默认启用）
coco.EnableHistory(true)
```

## OpenAPI 生成工具集成

### 使用 Huma v2

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
    
    // 注册你的 API
    huma.Register(api, huma.Operation{
        OperationID: "get-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "Get users",
    }, func(ctx context.Context, input *struct{}) (*struct{}, error) {
        return &struct{}{}, nil
    })
    
    // 获取 OpenAPI 规范并挂载文档
    spec, _ := api.OpenAPI().MarshalJSON()
    mux.Handle("/docs/", coco.New("",
        coco.Spec(spec),
        coco.Title("My API - Huma"),
    ))
    
    http.ListenAndServe(":8000", mux)
}
```

### 使用 Swag

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
    
    // 你的 API 路由
    mux.HandleFunc("/api/users", getUsers)
    
    // 挂载文档（swag 会生成 docs/swagger.json）
    mux.Handle("/docs/", coco.New("./docs/swagger.json",
        coco.Title("My API - Swag"),
    ))
    
    http.ListenAndServe(":8000", mux)
}
```

运行前需要生成文档：
```bash
swag init
```

## 主题和语言

### 主题选项

- `light` - 亮色主题
- `dark` - 暗色主题
- `auto` - 自动跟随系统（默认）

### 语言选项

- `en` - English（默认）
- `zh` - 中文

用户可以在界面右上角随时切换主题和语言。

## 完整示例

查看 [example/framework](../example/framework) 目录获取更多完整示例：

- **net/http** - 原生标准库示例
- **Gin** - Gin 框架集成
- **Echo** - Echo 框架集成
- **Fiber** - Fiber v3 框架集成
- **Chi** - Chi 路由器集成

每个框架都提供了 Huma 和 Swag 两种 OpenAPI 生成方式的示例。

## 文档

* [快速入门指南](./01-getting-started.md)
* [配置指南](./02-configuration.md)
* [框架集成示例](./03-framework-integration.md)
* [OpenAPI 生成指南](./04-openapi-generation.md)
* [API 参考](./05-api-reference.md)
* [完整示例代码](../example/framework)

## 贡献

欢迎贡献代码、报告问题或提出建议！

如果你在使用 Coco 或觉得它对你有帮助，请给我们一个 Star ⭐

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 给个 Star！⭐

如果你喜欢这个项目，或者正在使用它来学习或构建你的解决方案，请给它一个 Star 以获取新版本的更新通知。你的支持很重要！

---

**Made with ❤️ by [leehainuo](https://github.com/leehainuo)**
