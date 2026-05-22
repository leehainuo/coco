# Coco Framework Examples

这个目录包含了 Coco 与各种流行 Go Web 框架集成的示例，每个框架都提供了两种 OpenAPI 规范生成方式：**Huma** 和 **Swag**。

## 📁 目录结构

```
framework/
├── nethttp/
│   ├── huma/     # net/http + Huma v2 (原生标准库)
│   └── swag/     # net/http + Swag (原生标准库)
├── gin/
│   ├── huma/     # Gin + Huma v2
│   └── swag/     # Gin + Swag
├── echo/
│   ├── huma/     # Echo + Huma v2
│   └── swag/     # Echo + Swag
├── fiber/
│   ├── huma/     # Fiber + Huma v2
│   └── swag/     # Fiber + Swag
└── chi/
    ├── huma/     # Chi + Huma v2
    └── swag/     # Chi + Swag
```

## 🎯 两种方式对比

### Huma v2

**特点：**
- ✅ 代码优先（Code-First）
- ✅ 类型安全，自动生成 OpenAPI 规范
- ✅ 内置验证和序列化
- ✅ 支持多种框架适配器
- ✅ 无需额外的代码生成步骤

**适用场景：**
- 新项目
- 追求类型安全
- 喜欢声明式 API 定义

### Swag

**特点：**
- ✅ 注释优先（Annotation-First）
- ✅ 通过注释生成 OpenAPI 规范
- ✅ 与现有代码集成简单
- ✅ 需要运行 `swag init` 生成文档

**适用场景：**
- 已有项目
- 喜欢注释风格
- 需要更灵活的文档控制

## 🚀 快速开始

### 0. net/http 原生标准库

#### net/http + Huma
```bash
cd nethttp/huma
go mod tidy
go run main.go
```

#### net/http + Swag
```bash
cd nethttp/swag
~/go/bin/swag init
go run main.go
```

访问：http://localhost:8000/docs/

---

### 1. Gin 框架

#### Gin + Huma
```bash
cd gin/huma
go mod tidy
go run main.go
```

#### Gin + Swag
```bash
cd gin/swag
go install github.com/swaggo/swag/cmd/swag@latest
swag init
go run main.go
```

访问：http://localhost:8000/docs/

---

### 2. Echo 框架

#### Echo + Huma
```bash
cd echo/huma
go mod tidy
go run main.go
```

#### Echo + Swag
```bash
cd echo/swag
swag init
go run main.go
```

访问：http://localhost:8000/docs/

---

### 3. Fiber 框架

#### Fiber + Huma
```bash
cd fiber/huma
go mod tidy
go run main.go
```

#### Fiber + Swag
```bash
cd fiber/swag
swag init
go run main.go
```

访问：http://localhost:8000/docs/

---

### 4. Chi 框架

#### Chi + Huma
```bash
cd chi/huma
go mod tidy
go run main.go
```

#### Chi + Swag
```bash
cd chi/swag
swag init
go run main.go
```

访问：http://localhost:8000/docs/

## 📦 依赖安装

### Huma 示例依赖

```bash
# net/http (原生标准库)
go get github.com/danielgtaylor/huma/v2
go get github.com/danielgtaylor/huma/v2/adapters/humago

# Gin
go get github.com/danielgtaylor/huma/v2/adapters/humagin
go get github.com/gin-gonic/gin

# Echo
go get github.com/danielgtaylor/huma/v2/adapters/humaecho
go get github.com/labstack/echo/v4

# Fiber
go get github.com/danielgtaylor/huma/v2/adapters/humafiber
go get github.com/gofiber/fiber/v3

# Chi
go get github.com/danielgtaylor/huma/v2/adapters/humachi
go get github.com/go-chi/chi/v5
```

### Swag 示例依赖

```bash
# 安装 swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# 各框架依赖（net/http 是标准库，无需额外安装）
go get github.com/gin-gonic/gin
go get github.com/labstack/echo/v4
go get github.com/gofiber/fiber/v3
go get github.com/go-chi/chi/v5
```

## 🔧 Swag 使用说明

对于所有 Swag 示例，需要先生成文档：

```bash
# 在对应的 swag 目录下运行
swag init

# 这会生成 docs/ 目录，包含：
# - docs.go
# - swagger.json
# - swagger.yaml
```

## 💡 代码示例

### Huma 风格

```go
// 定义输入输出类型
type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" minLength:"1" maxLength:"50"`
        Email string `json:"email" format:"email"`
    }
}

type CreateUserOutput struct {
    Body User
}

// 注册操作
huma.Register(api, huma.Operation{
    OperationID: "create-user",
    Method:      http.MethodPost,
    Path:        "/api/users",
    Summary:     "Create a new user",
    Tags:        []string{"Users"},
}, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
    // 实现逻辑
    return &CreateUserOutput{Body: user}, nil
})

// 挂载 Coco 文档
spec, _ := api.OpenAPI().YAML()
router.Handle("/docs/*", coco.New("",
    coco.Spec([]byte(spec)),
    coco.Title("My API"),
))
```

### Swag 风格

```go
// @title My API
// @version 1.0
// @description This is my API
// @host localhost:8000
// @BasePath /api

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
    // 实现逻辑
}

// 挂载 Coco 文档
router.Handle("/docs/*", coco.New("./docs/swagger.json",
    coco.Title("My API"),
))
```

## 🎨 特点对比

| 特性 | Huma | Swag |
|------|------|------|
| **类型安全** | ✅ 强类型 | ⚠️ 依赖注释 |
| **自动验证** | ✅ 内置 | ❌ 需手动 |
| **代码生成** | ❌ 不需要 | ✅ 需要 swag init |
| **学习曲线** | 📈 中等 | 📉 简单 |
| **灵活性** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **维护成本** | 低 | 中（需保持注释同步） |

## 📝 注意事项

1. **端口冲突**：所有示例默认使用 8000 端口，不能同时运行
2. **Swag 文档生成**：Swag 示例需要先运行 `swag init` 生成文档
3. **依赖安装**：首次运行需要 `go mod tidy` 安装依赖
4. **Lint 错误**：未安装依赖前会有 lint 错误，这是正常的

## 🌟 推荐选择

- **新项目 + 追求类型安全** → 选择 **Huma**
- **已有项目 + 快速集成** → 选择 **Swag**
- **团队熟悉注释风格** → 选择 **Swag**
- **需要强验证和序列化** → 选择 **Huma**

## 🔗 相关链接

- [Huma 官方文档](https://huma.rocks/)
- [Swag 官方文档](https://github.com/swaggo/swag)
- [Coco 项目主页](https://github.com/leehainuo/coco)
