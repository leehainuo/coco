# Coco Basic Examples

这个目录包含了 Coco 的基础使用示例，展示了不同的 OpenAPI 规范加载方式。

## 示例列表

### 1. Remote - 远程 URL 加载

**位置：** `cmd/remote/main.go`

**端口：** 8000

**特点：**
- 从远程 URL 加载 OpenAPI 规范
- 使用 Petstore API 作为示例
- 适合展示第三方 API 文档

**运行：**
```bash
cd cmd/remote
go run main.go
```

**访问：** http://localhost:8000/docs/

---

### 2. Embedded - 嵌入式加载

**位置：** `cmd/embedded/main.go`

**端口：** 8000

**特点：**
- 使用 `//go:embed` 将 OpenAPI JSON 嵌入到二进制文件中
- 无需外部文件依赖
- 适合生产环境部署

**运行：**
```bash
cd cmd/embedded
go run main.go
```

**访问：** http://localhost:8000/docs/

---

### 3. Local - 本地文件加载

**位置：** `cmd/local/main.go`

**端口：** 8000

**特点：**
- 从本地文件系统加载 OpenAPI 规范
- 支持开发时动态修改规范文件
- 适合开发和测试环境

**运行：**
```bash
cd cmd/local
go run main.go
```

**访问：** http://localhost:8000/docs/

---

### 4. Net/HTTP - 纯标准库示例

**位置：** `cmd/nethttp/main.go`

**端口：** 8082

**特点：**
- 仅使用 Go 标准库 `net/http`
- 包含简单的 API 实现
- 展示如何在纯标准库项目中使用 Coco

**运行：**
```bash
cd cmd/nethttp
go run main.go
```

**访问：**
- 文档：http://localhost:8000/docs/
- API：
  - http://localhost:8000/api/health
  - http://localhost:8000/api/users

---

## 配置选项示例

所有示例都支持以下配置选项：

```go
coco.New(specPath,
    coco.Title("自定义标题"),           // 设置文档标题
    coco.Theme("dark"),               // 主题：auto/light/dark
    coco.Lang("zh"),                  // 语言：en/zh
    coco.EnableDebug(true),           // 启用调试模式
    coco.EnableExport(true),          // 启用导出功能
    coco.EnableHistory(true),         // 启用历史记录
)
```

## 项目结构

```
example/basic/
├── cmd/
│   ├── remote/          # 远程 URL 示例
│   │   └── main.go
│   ├── embedded/        # 嵌入式示例
│   │   ├── main.go
│   │   └── openapi.json
│   ├── local/           # 本地文件示例
│   │   └── main.go
│   └── nethttp/         # 纯标准库示例
│       └── main.go
├── pkg/                 # 共享的 API 实现
│   └── api.go
└── README.md
```

## 快速开始

1. **安装依赖：**
```bash
go mod tidy
```

2. **选择一个示例运行：**
```bash
# 运行远程示例
cd cmd/remote && go run main.go

# 或运行嵌入式示例
cd cmd/embedded && go run main.go

# 或运行本地文件示例
cd cmd/local && go run main.go

# 或运行纯标准库示例
cd cmd/nethttp && go run main.go
```

3. **在浏览器中访问对应的文档地址**

## 注意事项

- 每个示例使用不同的端口，可以同时运行多个示例
- Remote 示例需要网络连接来获取 Petstore API 规范
- Embedded 和 Local 示例使用相同的 `openapi.json` 文件
- 所有示例都展示了 Coco 的零配置、自动路径检测特性
