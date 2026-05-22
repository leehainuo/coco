# 快速入门指南

本指南将帮助你在 5 分钟内开始使用 Coco。

## 前置要求

- Go 1.21 或更高版本
- 一个 OpenAPI/Swagger 规范文件（JSON 格式， 可以通过 Swag 和 Huma 等工具生成）

## 第一步：安装

使用 `go get` 安装 Coco：

```bash
go get github.com/leehainuo/coco
```

## 第二步：准备 OpenAPI 规范

你需要一个 OpenAPI 规范文件。如果还没有，可以创建一个简单的示例：

**openapi.json**
```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "Hello API",
    "version": "1.0.0",
    "description": "一个简单的示例 API"
  },
  "paths": {
    "/hello": {
      "get": {
        "summary": "打招呼",
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "message": {
                      "type": "string",
                      "example": "Hello, World!"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

## 第三步：编写代码

创建 `main.go` 文件：

```go
package main

import (
    "net/http"
    "github.com/leehainuo/coco"
)

func main() {
    // 创建文档处理器
    handler := coco.New("./openapi.json",
        coco.Title("Hello API 文档"),
        coco.Lang("zh"),
    )
    
    // 挂载到 /docs/ 路径
    http.Handle("/docs/", handler)
    
    // 启动服务器
    println("文档服务器启动在 http://localhost:8080/docs/")
    http.ListenAndServe(":8080", nil)
}
```

## 第四步：运行

```bash
go run main.go
```

打开浏览器访问 http://localhost:8080/docs/，你将看到美观的 API 文档界面！

## 下一步

现在你已经成功运行了 Coco，可以：

1. **自定义配置** - 查看 [配置指南](./02-configuration.md) 了解所有配置选项
2. **框架集成** - 查看 [框架集成指南](./03-framework-integration.md) 学习如何与你喜欢的框架集成
3. **OpenAPI 生成** - 查看 [OpenAPI 生成指南](./04-openapi-generation.md) 学习如何自动生成 OpenAPI 规范

## 常见问题

### 文档页面显示空白？

检查以下几点：
1. OpenAPI 文件路径是否正确
2. OpenAPI 文件格式是否有效（可以使用在线验证器验证）
3. 浏览器控制台是否有错误信息

### 如何自定义文档路径？

```go
// 路径由你在挂载时决定，没有默认值
http.Handle("/docs/", handler)      // 使用 /docs/
http.Handle("/api-docs/", handler)  // 使用 /api-docs/
http.Handle("/swagger/", handler)   // 使用 /swagger/
```

### 支持哪些 OpenAPI 版本？

Coco 支持：
- OpenAPI 3.0.x
- OpenAPI 3.1.x
- Swagger 2.0

### 如何在生产环境使用？

建议禁用调试功能：

```go
handler := coco.New("./openapi.json",
    coco.EnableDebug(false),
)
```

## 提示

- 使用 `coco.Theme("dark")` 设置默认暗色主题
- 使用 `coco.Lang("zh")` 设置中文界面
- 文档界面右上角可以切换主题和语言
- 调试面板可以直接测试 API 接口

## 相关资源

- [完整 API 参考](./05-api-reference.md)
- [示例代码](../../example/framework/)
- [配置选项详解](./02-configuration.md)
