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
coco.Lang("zh")     // 中文
coco.Lang("en")     // 英文
coco.Lang("custom") // 自定义语言（见下文）
```

**可选值**:
- `"en"` - English
- `"zh"` - 中文
- `"custom"` - 通过 `coco.I18n` / `coco.I18nData` / `coco.I18nURL` 加载的自定义语言

**默认值**: `"en"`

**注意**: 用户可以在界面右上角随时切换语言。

### 自定义语言（i18n）

除内置的中文与英文外，你可以通过提供自己的翻译文件新增第三种语言。使用下列
任意一个选项注册翻译，再通过 `coco.Lang("custom")` 选中它：

```go
coco.I18n("./i18n.fr.json")                          // 从本地文件加载
coco.I18nData(frI18n)                                // 从内存字节加载
coco.I18nURL("https://cdn.example.com/i18n.fr.json") // 从远程 URL 加载
```

翻译文件是如下结构的 JSON 对象：

```json
{
  "name": "Français",
  "messages": {
    "search": "Recherche",
    "loading": "Chargement...",
    "response": "Réponse"
  }
}
```

- `name` - 语言切换器中显示的名称。
- `messages` - 翻译键到本地化文案的扁平映射表。

**说明**:
- 缺失的键会自动回退到英文。
- 未配置自定义语言时，切换器只显示中文与英文。
- 设置 `coco.Lang("custom")` 可将自定义语言设为初始语言；否则用户可从切换器中手动选择。

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

## 翻译配置

自定义语言从翻译文件加载。与规范一样，它可以来自三种来源，解析优先级为：
内存字节 → 远程 URL → 本地文件。

### 从文件加载

从本地文件系统加载翻译文件：

```go
handler := coco.New("./openapi.json",
    coco.I18n("./i18n.fr.json"),
    coco.Lang("custom"),
)
```

**使用场景**:
- 本地开发，修改翻译文件后无需重新编译
- 翻译文件随二进制一起部署

### 从字节数组加载

传入已在内存中的翻译字节：

```go
handler := coco.New("./openapi.json",
    coco.I18nData(frI18n),
    coco.Lang("custom"),
)
```

`I18nData` 接收任意 `[]byte`，来源完全由你决定：

**示例 - 使用 embed**（单文件自包含二进制）:
```go
import _ "embed"

//go:embed i18n.fr.json
var frI18n []byte

func main() {
    handler := coco.New("./openapi.json",
        coco.I18nData(frI18n),
        coco.Lang("custom"),
    )
    // ...
}
```

**示例 - 使用内联 Go 字符串**（完全不需要外部文件）:
```go
const frI18n = `{
    "name": "Français",
    "messages": {
        "search": "Recherche",
        "loading": "Chargement..."
    }
}`

handler := coco.New("./openapi.json",
    coco.I18nData([]byte(frI18n)),
    coco.Lang("custom"),
)
```

**示例 - 运行时构建**（例如从数据库读取）:
```go
payload, _ := json.Marshal(map[string]any{
    "name":     "Français",
    "messages": messagesFromDB, // map[string]string
})

handler := coco.New("./openapi.json",
    coco.I18nData(payload),
    coco.Lang("custom"),
)
```

**使用场景**:
- 使用 Go `embed` 嵌入翻译文件，做到零外部依赖
- 直接把 JSON 定义为 Go 字符串常量
- 动态生成翻译，或从数据库 / 配置中心读取

### 从远程 URL 加载

从远程服务器获取翻译文件：

```go
handler := coco.New("./openapi.json",
    coco.I18nURL("https://cdn.example.com/i18n.fr.json"),
    coco.Lang("custom"),
)
```

**使用场景**:
- 翻译托管在 CDN、对象存储或配置中心
- 多个服务共享同一份翻译并集中更新

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
