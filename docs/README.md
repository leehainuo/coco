# Coco Documentation

<div align="center">
  <img src="images/coco.png" alt="Coco Logo" width="160"/>
  <br/>
  <p>优雅 · 轻量 · 美观 | Elegant · Lightweight · Beautiful</p>
</div>

<br/>

Welcome to Coco documentation! Coco is an elegant, lightweight, and beautiful OpenAPI documentation renderer for Go.

欢迎使用 Coco 文档！Coco 是一个优雅、轻量、美观的 Go OpenAPI 文档渲染库。

Choose your preferred language / 选择你的语言：

## 📚 Documentation Languages

### 🇨🇳 中文文档 (Chinese)

完整的中文文档，包含详细的使用指南和示例。

**[进入中文文档 →](./zh/README.md)**

**文档列表**:
- [README](./zh/README.md) - 中文完整文档
- [贡献指南](./zh/CONTRIBUTING.md) - 如何为 Coco 做贡献
- [01. 快速入门](./zh/01-getting-started.md) - 5 分钟快速上手
- [02. 配置指南](./zh/02-configuration.md) - 详细配置说明
- [03. 框架集成](./zh/03-framework-integration.md) - 与各框架集成
- [04. OpenAPI 生成指南](./zh/04-openapi-generation.md) - 生成 OpenAPI 规范
- [05. API 参考](./zh/05-api-reference.md) - 完整 API 文档

---

### 🇬🇧 English Documentation

Complete English documentation with detailed guides and examples.

**[Go to English Docs →](../README.md)**

**Documentation List**:
- [README](../README.md) - Complete English documentation
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute to Coco
- [01. Quick Start](./en/01-getting-started.md) - Get started in 5 minutes
- [02. Configuration Guide](./en/02-configuration.md) - Detailed configuration
- [03. Framework Integration](./en/03-framework-integration.md) - Framework integration
- [04. OpenAPI Generation Guide](./en/04-openapi-generation.md) - Generate OpenAPI specs
- [05. API Reference](./en/05-api-reference.md) - Complete API docs

---

## 🎯 Quick Links

### For Beginners

- **中文用户**: 从 [快速入门](./zh/01-getting-started.md) 开始
- **English Users**: Start with [Quick Start](./en/01-getting-started.md)

### Common Tasks

| Task | 中文 | English |
|------|------|---------|
| Main README | [中文文档](./zh/README.md) | [English Docs](../README.md) |
| Installation & Setup | [快速入门](./zh/01-getting-started.md) | [Quick Start](./en/01-getting-started.md) |
| Configuration Options | [配置指南](./zh/02-configuration.md) | [Configuration Guide](./en/02-configuration.md) |
| Framework Integration | [框架集成](./zh/03-framework-integration.md) | [Framework Integration](./en/03-framework-integration.md) |
| Generate OpenAPI Specs | [OpenAPI 生成](./zh/04-openapi-generation.md) | [OpenAPI Generation](./en/04-openapi-generation.md) |
| API Reference | [API 参考](./zh/05-api-reference.md) | [API Reference](./en/05-api-reference.md) |

---

## 📖 Documentation Structure

```
coco/
├── README.md              # English main documentation
├── CONTRIBUTING.md        # English contributing guide
└── docs/
    ├── README.md          # This file (language selection)
    ├── images/
    ├── zh/                # Chinese documentation
    │   ├── README.md      # Chinese docs home
    │   ├── CONTRIBUTING.md
    │   ├── 01-getting-started.md
    │   ├── 02-configuration.md
    │   ├── 03-framework-integration.md
    │   ├── 04-openapi-generation.md
    │   └── 05-api-reference.md
    └── en/                # English detailed guides
        ├── CONTRIBUTING.md
        ├── 01-getting-started.md
        ├── 02-configuration.md
        ├── 03-framework-integration.md
        ├── 04-openapi-generation.md
        └── 05-api-reference.md
```

---

## 💡 Examples

Check the [example/framework](../example/framework) directory for complete runnable examples:

- **net/http** - Standard library examples
- **Gin** - Gin framework integration
- **Echo** - Echo framework integration
- **Fiber** - Fiber v3 framework integration
- **Chi** - Chi router integration

Each framework provides examples for both Huma and Swag OpenAPI generation methods.

---

## 🤝 Contributing

Documentation improvements are welcome in both languages! 

欢迎贡献文档改进，支持中英文！

---

**Made with ❤️ by [leehainuo](https://github.com/leehainuo)**
