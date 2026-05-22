# 配置指南

Coco 提供了丰富的配置选项，让你可以自定义文档的外观和行为。

## 配置方式

所有配置都通过 `Option` 函数传递给 `coco.New()`：

```go
handler := coco.New("./openapi.json",
    coco.Title("我的 API"),
    coco.Theme("dark"),
    coco.Lang("zh"),
    // ... 更多配置
)
```

## UI 配置

### 文档标题

设置文档页面的标题（显示在浏览器标签和页面顶部）：

```go
coco.Title("我的 API 文档")
```

**默认值**: `"Coco API Docs"`

### 主题设置

设置默认主题：

```go
coco.Theme("dark")  // 暗色主题
coco.Theme("light") // 亮色主题
coco.Theme("auto")  // 自动跟随系统（默认）
```

**可选值**:
- `"light"` - 亮色主题
- `"dark"` - 暗色主题
- `"auto"` - 自动跟随系统设置

**默认值**: `"auto"`

**注意**: 用户可以在界面右上角随时切换主题，选择会保存在浏览器本地存储中。

### 语言设置

设置界面语言：

```go
coco.Lang("zh") // 中文
coco.Lang("en") // 英文
```

**可选值**:
- `"en"` - English
- `"zh"` - 中文

**默认值**: `"en"`

**注意**: 用户可以在界面右上角随时切换语言。

## 规范配置

### 从文件加载

最常用的方式，从本地文件加载 OpenAPI 规范：

```go
coco.New("./openapi.json")
coco.New("./swagger.json")
coco.New("./docs/api-spec.json")
```

支持相对路径和绝对路径，仅支持 JSON 格式。

### 从字节数组加载

适合动态生成或嵌入的规范：

```go
spec := []byte(`{
    "openapi": "3.0.0",
    "info": {
        "title": "My API",
        "version": "1.0.0"
    }
}`)

handler := coco.New("", coco.Spec(spec))
```

**使用场景**:
- 与 Huma 等代码生成工具集成
- 从数据库或配置中心加载规范
- 使用 Go embed 嵌入规范文件

**示例 - 使用 embed**:
```go
import _ "embed"

//go:embed openapi.json
var spec []byte

func main() {
    handler := coco.New("", coco.Spec(spec))
    // ...
}
```

### 从远程 URL 加载

从远程服务器加载规范：

```go
handler := coco.New("", 
    coco.SpecURL("https://api.example.com/openapi.json"),
)
```

**注意**: 
- 确保 URL 可访问
- 考虑网络延迟和可用性
- 生产环境建议使用本地文件或嵌入方式

## 功能开关

### 调试面板

启用或禁用 API 调试面板：

```go
coco.EnableDebug(true)  // 启用（默认）
coco.EnableDebug(false) // 禁用
```

**默认值**: `true`

**功能**:
- 直接在文档界面测试 API
- 填写参数、请求体
- 查看响应结果
- 支持各种 HTTP 方法

**建议**:
- 开发环境：启用
- 生产环境：根据需要决定（如果是内部文档可以启用）

### 导出功能

启用或禁用 OpenAPI 规范导出功能：

```go
coco.EnableExport(true)  // 启用（默认）
coco.EnableExport(false) // 禁用
```

**默认值**: `true`

**功能**:
- 允许用户下载 OpenAPI 规范文件（JSON 格式）
- 方便与其他工具集成

### 历史记录

启用或禁用请求历史记录功能：

```go
coco.EnableHistory(true)  // 启用（默认）
coco.EnableHistory(false) // 禁用
```

**默认值**: `true`

**功能**:
- 保存调试面板的请求历史
- 存储在浏览器本地存储中
- 方便重复测试

## 完整配置示例

### 开发环境配置

```go
handler := coco.New("./openapi.json",
    coco.Title("开发环境 - API 文档"),
    coco.Theme("dark"),
    coco.Lang("zh"),
    coco.EnableDebug(true),
    coco.EnableExport(true),
    coco.EnableHistory(true),
)
```

### 生产环境配置

```go
handler := coco.New("./openapi.json",
    coco.Title("生产环境 - API 文档"),
    coco.Theme("auto"),
    coco.Lang("en"),
    coco.EnableDebug(false),  // 禁用调试
    coco.EnableExport(false), // 禁用导出
    coco.EnableHistory(false),
)
```

### 内部文档配置

```go
handler := coco.New("./openapi.json",
    coco.Title("内部 API 文档"),
    coco.Theme("auto"),
    coco.Lang("zh"),
    coco.EnableDebug(true),   // 内部可以调试
    coco.EnableExport(true),  // 允许导出
    coco.EnableHistory(true),
)
```

### 使用嵌入规范

```go
import _ "embed"

//go:embed openapi.json
var apiSpec []byte

func main() {
    handler := coco.New("",
        coco.Spec(apiSpec),
        coco.Title("嵌入式 API 文档"),
        coco.Lang("zh"),
    )
    // ...
}
```

## 高级用法

### 动态配置

根据环境变量动态配置：

```go
import "os"

func createHandler() http.Handler {
    isDev := os.Getenv("ENV") == "development"
    
    return coco.New("./openapi.json",
        coco.Title(getTitle()),
        coco.EnableDebug(isDev),
        coco.EnableExport(isDev),
    )
}

func getTitle() string {
    if os.Getenv("ENV") == "production" {
        return "生产环境 API"
    }
    return "开发环境 API"
}
```

### 多个文档实例

为不同的 API 版本创建不同的文档：

```go
// API v1 文档
v1Handler := coco.New("./openapi-v1.json",
    coco.Title("API v1 文档"),
)
http.Handle("/docs/v1/", v1Handler)

// API v2 文档
v2Handler := coco.New("./openapi-v2.json",
    coco.Title("API v2 文档"),
)
http.Handle("/docs/v2/", v2Handler)
```

## 最佳实践

1. **使用环境变量** - 根据环境动态配置功能开关
2. **嵌入规范文件** - 生产环境使用 `embed` 避免文件路径问题
3. **合理的默认值** - 开发环境启用所有功能，生产环境谨慎选择
4. **清晰的标题** - 使用描述性的标题，方便用户识别
5. **考虑用户习惯** - 中文用户使用中文界面，国际化产品使用英文

## 相关文档

- [快速入门](./01-getting-started.md)
- [框架集成](./03-framework-integration.md)
- [API 参考](./05-api-reference.md)
