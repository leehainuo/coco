# API 参考

Coco 库的完整 API 参考文档。

## 核心函数

### New

创建一个新的文档处理器。

```go
func New(path string, opts ...Option) http.Handler
```

**参数**:
- `path` - OpenAPI 规范文件的路径（相对或绝对路径）
- `opts` - 可选的配置选项

**返回值**:
- `http.Handler` - 实现了标准 HTTP 处理器接口

**示例**:
```go
// 从文件加载
handler := coco.New("./openapi.json")

// 带配置选项
handler := coco.New("./openapi.json",
    coco.Title("My API"),
    coco.Theme("dark"),
)
```

---

## 配置选项

所有配置选项都是 `Option` 类型的函数，可以传递给 `New()` 函数。

### Spec

从字节数组加载 OpenAPI 规范。

```go
func Spec(data []byte) Option
```

**参数**:
- `data` - OpenAPI 规范的字节数组（JSON 格式）

**示例**:
```go
spec := []byte(`{"openapi": "3.0.0", ...}`)
handler := coco.New("", coco.Spec(spec))
```

**注意**: 使用此选项时，第一个参数应传空字符串。

---

### SpecURL

从远程 URL 加载 OpenAPI 规范。

```go
func SpecURL(url string) Option
```

**参数**:
- `url` - OpenAPI 规范的远程 URL

**示例**:
```go
handler := coco.New("", 
    coco.SpecURL("https://api.example.com/openapi.json"),
)
```

**注意**: 使用此选项时，第一个参数应传空字符串。

---

### Title

设置文档页面的标题。

```go
func Title(title string) Option
```

**参数**:
- `title` - 文档标题字符串

**默认值**: `"Coco API Docs"`

**示例**:
```go
handler := coco.New("./openapi.json",
    coco.Title("我的 API 文档"),
)
```

**效果**:
- 显示在浏览器标签页标题
- 显示在文档页面顶部

---

### Theme

设置默认主题。

```go
func Theme(theme string) Option
```

**参数**:
- `theme` - 主题名称

**可选值**:
- `"light"` - 亮色主题
- `"dark"` - 暗色主题
- `"auto"` - 自动跟随系统（默认）

**默认值**: `"auto"`

**示例**:
```go
handler := coco.New("./openapi.json",
    coco.Theme("dark"),
)
```

**注意**: 用户可以在界面右上角随时切换主题，选择会保存在浏览器本地存储中。

---

### Lang

设置界面语言。

```go
func Lang(lang string) Option
```

**参数**:
- `lang` - 语言代码

**可选值**:
- `"en"` - English
- `"zh"` - 中文

**默认值**: `"en"`

**示例**:
```go
handler := coco.New("./openapi.json",
    coco.Lang("zh"),
)
```

**注意**: 用户可以在界面右上角随时切换语言。

---

### EnableDebug

启用或禁用调试面板。

```go
func EnableDebug(debug bool) Option
```

**参数**:
- `debug` - `true` 启用，`false` 禁用

**默认值**: `true`

**示例**:
```go
// 启用调试面板（默认）
handler := coco.New("./openapi.json",
    coco.EnableDebug(true),
)

// 禁用调试面板
handler := coco.New("./openapi.json",
    coco.EnableDebug(false),
)
```

**功能**:
- 在文档界面测试 API
- 填写参数和请求体
- 查看响应结果
- 支持各种 HTTP 方法

---

### EnableExport

启用或禁用导出功能。

```go
func EnableExport(export bool) Option
```

**参数**:
- `export` - `true` 启用，`false` 禁用

**默认值**: `true`

**示例**:
```go
// 启用导出功能（默认）
handler := coco.New("./openapi.json",
    coco.EnableExport(true),
)

// 禁用导出功能
handler := coco.New("./openapi.json",
    coco.EnableExport(false),
)
```

**功能**:
- 允许用户下载 OpenAPI 规范文件（JSON 格式）

---

### EnableHistory

启用或禁用历史记录功能。

```go
func EnableHistory(history bool) Option
```

**参数**:
- `history` - `true` 启用，`false` 禁用

**默认值**: `true`

**示例**:
```go
// 启用历史记录（默认）
handler := coco.New("./openapi.json",
    coco.EnableHistory(true),
)

// 禁用历史记录
handler := coco.New("./openapi.json",
    coco.EnableHistory(false),
)
```

**功能**:
- 保存调试面板的请求历史
- 存储在浏览器本地存储中
- 方便重复测试

---

## 类型定义

### Option

配置选项函数类型。

```go
type Option func(*config.Config)
```

这是一个函数类型，用于修改配置。所有配置函数都返回此类型。

---

## 配置结构

虽然这些结构是内部使用的，但了解它们有助于理解配置的工作原理。

### Config

```go
type Config struct {
    Spec
    UI
    Feature
}
```

### Spec

规范配置。

```go
type Spec struct {
    Path string  // 文件路径
    Data []byte  // 字节数组
    URL  string  // 远程 URL
}
```

### UI

界面配置。

```go
type UI struct {
    Title string // 文档标题
    Theme string // 主题: "light", "dark", "auto"
    Lang  string // 语言: "en", "zh"
}
```

### Feature

功能开关。

```go
type Feature struct {
    Debug   bool // 调试面板
    Export  bool // 导出功能
    History bool // 历史记录
}
```

---

## 使用示例

### 最小配置

```go
handler := coco.New("./openapi.json")
http.Handle("/docs/", handler)
```

### 完整配置

```go
handler := coco.New("./openapi.json",
    coco.Title("我的 API 文档"),
    coco.Theme("dark"),
    coco.Lang("zh"),
    coco.EnableDebug(true),
    coco.EnableExport(true),
    coco.EnableHistory(true),
)
http.Handle("/docs/", handler)
```

### 使用字节数组

```go
import _ "embed"

//go:embed openapi.json
var spec []byte

handler := coco.New("",
    coco.Spec(spec),
    coco.Title("嵌入式 API"),
)
```

### 使用远程 URL

```go
handler := coco.New("",
    coco.SpecURL("https://api.example.com/openapi.json"),
    coco.Title("远程 API"),
)
```

### 环境相关配置

```go
import "os"

func createHandler() http.Handler {
    isProd := os.Getenv("ENV") == "production"
    
    opts := []coco.Option{
        coco.Title("My API"),
    }
    
    if !isProd {
        opts = append(opts,
            coco.EnableDebug(true),
            coco.Lang("zh"),
        )
    } else {
        opts = append(opts,
            coco.EnableDebug(false),
            coco.EnableExport(false),
        )
    }
    
    return coco.New("./openapi.json", opts...)
}
```

---

## 错误处理

Coco 会在以下情况返回错误响应：

1. **文件不存在** - 返回 404 Not Found
2. **文件格式错误** - 返回 500 Internal Server Error
3. **URL 无法访问** - 返回 500 Internal Server Error

建议在开发时检查日志输出以诊断问题。

---

## 性能考虑

1. **文件缓存** - OpenAPI 规范文件会被缓存在内存中
2. **静态资源** - 前端资源已嵌入到二进制文件中
3. **并发安全** - 所有操作都是并发安全的

---

## 兼容性

- **Go 版本**: >= 1.21
- **OpenAPI 版本**: 2.0, 3.0.x, 3.1.x
- **浏览器**: 现代浏览器（Chrome, Firefox, Safari, Edge）

---

## 相关文档

- [快速入门](./01-getting-started.md)
- [配置指南](./02-configuration.md)
- [框架集成](./03-framework-integration.md)
- [OpenAPI 生成](./04-openapi-generation.md)
