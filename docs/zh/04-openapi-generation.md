# OpenAPI 规范生成指南

本指南介绍如何生成 OpenAPI 规范文件，以便与 Coco 配合使用。

## 两种主流方式

### 1. 代码优先 (Code-First) - Huma

通过代码定义 API，自动生成 OpenAPI 规范。

**优点**:
- ✅ 类型安全
- ✅ 自动验证
- ✅ 规范与代码同步
- ✅ 无需手动维护文档

**缺点**:
- ❌ 学习曲线较陡
- ❌ 需要适配不同框架

### 2. 注释优先 (Annotation-First) - Swag

通过代码注释生成 OpenAPI 规范。

**优点**:
- ✅ 简单易用
- ✅ 与现有代码集成方便
- ✅ 灵活控制文档

**缺点**:
- ❌ 需要手动保持注释与代码同步
- ❌ 缺少编译时检查
- ❌ 需要额外的生成步骤

---

## 使用 Huma v2

### 安装

```bash
go get github.com/danielgtaylor/huma/v2
```

根据你的框架选择适配器：

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

### 基本用法

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/leehainuo/coco"
)

// 定义数据模型
type User struct {
    ID    string `json:"id" example:"1" doc:"用户 ID"`
    Name  string `json:"name" example:"Alice" doc:"用户名"`
    Email string `json:"email" example:"alice@example.com" doc:"邮箱"`
}

// 定义输入
type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" minLength:"1" maxLength:"50" doc:"用户名"`
        Email string `json:"email" format:"email" doc:"邮箱地址"`
    }
}

// 定义输出
type CreateUserOutput struct {
    Body User
}

type ListUsersOutput struct {
    Body []User
}

func main() {
    mux := http.NewServeMux()
    
    // 创建 Huma API
    api := humago.New(mux, huma.DefaultConfig("My API", "1.0.0"))
    
    // 注册 GET 操作
    huma.Register(api, huma.Operation{
        OperationID: "list-users",
        Method:      http.MethodGet,
        Path:        "/api/users",
        Summary:     "获取用户列表",
        Description: "返回所有用户",
        Tags:        []string{"Users"},
    }, func(ctx context.Context, input *struct{}) (*ListUsersOutput, error) {
        users := []User{
            {ID: "1", Name: "Alice", Email: "alice@example.com"},
            {ID: "2", Name: "Bob", Email: "bob@example.com"},
        }
        return &ListUsersOutput{Body: users}, nil
    })
    
    // 注册 POST 操作
    huma.Register(api, huma.Operation{
        OperationID: "create-user",
        Method:      http.MethodPost,
        Path:        "/api/users",
        Summary:     "创建用户",
        Description: "创建一个新用户",
        Tags:        []string{"Users"},
    }, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
        user := User{
            ID:    "3",
            Name:  input.Body.Name,
            Email: input.Body.Email,
        }
        return &CreateUserOutput{Body: user}, nil
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

### 高级特性

#### 路径参数

```go
type GetUserInput struct {
    UserID string `path:"id" doc:"用户 ID"`
}

type GetUserOutput struct {
    Body User
}

huma.Register(api, huma.Operation{
    OperationID: "get-user",
    Method:      http.MethodGet,
    Path:        "/api/users/{id}",
    Summary:     "获取单个用户",
}, func(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
    user := User{ID: input.UserID, Name: "Alice", Email: "alice@example.com"}
    return &GetUserOutput{Body: user}, nil
})
```

#### 查询参数

```go
type ListUsersInput struct {
    Page  int    `query:"page" default:"1" doc:"页码"`
    Limit int    `query:"limit" default:"10" doc:"每页数量"`
    Sort  string `query:"sort" enum:"name,email,created" doc:"排序字段"`
}
```

#### 请求头

```go
type AuthInput struct {
    Authorization string `header:"Authorization" doc:"认证令牌"`
}
```

#### 验证规则

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

## 使用 Swag

### 安装

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 基本用法

#### 1. 添加通用注释

在 `main.go` 文件顶部添加：

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

#### 2. 为每个 API 添加注释

```go
// listUsers godoc
// @Summary 获取用户列表
// @Description 返回所有用户
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func listUsers(w http.ResponseWriter, r *http.Request) {
    // 实现
}

// createUser godoc
// @Summary 创建用户
// @Description 创建一个新用户
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func createUser(w http.ResponseWriter, r *http.Request) {
    // 实现
}

// getUser godoc
// @Summary 获取单个用户
// @Description 根据 ID 获取用户
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func getUser(w http.ResponseWriter, r *http.Request) {
    // 实现
}
```

#### 3. 生成文档

```bash
swag init
```

这会在 `docs/` 目录下生成：
- `docs.go`
- `swagger.json`
- `swagger.yaml`

#### 4. 集成 Coco

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

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
```

### Swag 注释参考

#### 通用注释

| 注释 | 说明 | 示例 |
|------|------|------|
| `@title` | API 标题 | `@title My API` |
| `@version` | API 版本 | `@version 1.0` |
| `@description` | API 描述 | `@description This is my API` |
| `@host` | 主机地址 | `@host localhost:8080` |
| `@BasePath` | 基础路径 | `@BasePath /api` |

#### 操作注释

| 注释 | 说明 | 示例 |
|------|------|------|
| `@Summary` | 简短描述 | `@Summary Get users` |
| `@Description` | 详细描述 | `@Description Get all users` |
| `@Tags` | 标签分组 | `@Tags Users` |
| `@Accept` | 接受的内容类型 | `@Accept json` |
| `@Produce` | 返回的内容类型 | `@Produce json` |
| `@Param` | 参数定义 | `@Param id path string true "User ID"` |
| `@Success` | 成功响应 | `@Success 200 {object} User` |
| `@Failure` | 失败响应 | `@Failure 404 {object} Error` |
| `@Router` | 路由路径 | `@Router /users [get]` |
| `@Security` | 安全认证 | `@Security ApiKeyAuth` |

#### 参数类型

```go
// 路径参数
// @Param id path string true "User ID"

// 查询参数
// @Param page query int false "Page number"

// 请求头
// @Param Authorization header string true "Bearer token"

// 请求体
// @Param user body CreateUserRequest true "User info"

// 表单参数
// @Param name formData string true "User name"
```

---

## 对比总结

| 特性 | Huma | Swag |
|------|------|------|
| **学习曲线** | 中等 | 简单 |
| **类型安全** | ✅ 强类型 | ❌ 依赖注释 |
| **自动验证** | ✅ 内置 | ❌ 需手动 |
| **代码生成** | ❌ 不需要 | ✅ 需要 `swag init` |
| **维护成本** | 低 | 中（需保持注释同步） |
| **灵活性** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **框架支持** | 需要适配器 | 框架无关 |
| **适用场景** | 新项目、追求类型安全 | 已有项目、快速集成 |

---

## 最佳实践

### Huma 最佳实践

1. **使用有意义的 OperationID**
   ```go
   OperationID: "list-users"  // ✅ 好
   OperationID: "op1"         // ❌ 差
   ```

2. **添加详细的文档**
   ```go
   type User struct {
       ID   string `json:"id" doc:"用户唯一标识符"`
       Name string `json:"name" doc:"用户显示名称"`
   }
   ```

3. **使用标签分组**
   ```go
   Tags: []string{"Users", "Management"}
   ```

4. **定义错误响应**
   ```go
   type ErrorOutput struct {
       Body struct {
           Error string `json:"error" doc:"错误信息"`
       }
   }
   ```

### Swag 最佳实践

1. **保持注释与代码同步**
   - 修改 API 后立即更新注释
   - 定期运行 `swag init` 验证

2. **使用一致的命名**
   ```go
   // ✅ 好
   // @Summary Get user
   // @Router /users/{id} [get]
   
   // ❌ 差
   // @Summary getUserById
   // @Router /user/{userId} [get]
   ```

3. **添加示例值**
   ```go
   type User struct {
       ID   string `json:"id" example:"123"`
       Name string `json:"name" example:"Alice"`
   }
   ```

4. **使用 CI/CD 自动生成**
   ```yaml
   # .github/workflows/docs.yml
   - name: Generate docs
     run: swag init
   ```

---

## 相关资源

- [Huma 官方文档](https://huma.rocks/)
- [Swag 官方文档](https://github.com/swaggo/swag)
- [OpenAPI 规范](https://swagger.io/specification/)
- [完整示例](../../example/framework/)

---

## 相关文档

- [快速入门](./01-getting-started.md)
- [框架集成](./03-framework-integration.md)
- [配置指南](./02-configuration.md)
