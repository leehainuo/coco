# 🚀 Contributing to Coco

English | [简体中文](docs/zh/CONTRIBUTING.md)

Welcome to Coco! I'm glad you're here. Contributing to Coco is a great way to help improve the OpenAPI documentation experience for Go developers. Let's get started!

## 📜 Before You Start

As the author, I'm still a student with limited experience. This library may have some shortcomings, and I really appreciate your help and support!

### 🤝 Code of Conduct

Please be respectful and considerate in all interactions. Everyone can contribute.

### 🌟 Community Expectations

At Coco, we value:
- **Respect**: Treat everyone with kindness and professionalism
- **Collaboration**: Share knowledge and help each other grow
- **Quality**: Write clean, well-tested code
- **Communication**: Be clear and constructive in discussions

## 🚀 Getting Started

Here's how to begin your contribution journey:

1. 🍴 **Fork the Repository**: Fork Coco to your GitHub account.

2. 🛠️ **Clone Your Fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/coco.git
   cd coco
   ```

3. 🔧 **Set Up Development Environment**:
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Set up frontend (if working on UI)
   cd frontend
   npm i
   ```

4. 🌱 **Create a Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

5. 🚀 **Submit a Pull Request**: When ready, submit a PR to the main repository!

## 🌟 Your First Contribution

I welcome all kinds of contributions:

- 🐛 **Bug Fixes**: Found a bug? Fix it!
- ✨ **New Features**: Have an idea? Implement it! (must be reasonable)
- 📝 **Documentation**: Improve docs, add examples, fix typos
- 🌍 **Translations**: Add support for more languages
- 🎨 **UI Improvements**: Enhance the frontend experience
- 🧪 **Tests**: Add or improve test coverage

If you have questions or need guidance, feel free to submit an **issue**

## 🔍 Find Something to Work On

### 💼 Good First Issues

 - Improve Coco's documentation and examples

### 🪄 Work on an Issue

1. Comment on the **issue** to let me know you're working on it
2. Wait for a maintainer to assign it to you
3. Start coding!

### 📢 Report a Bug

Found a bug? Please submit an **issue** with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (Go version, OS, etc.)

## 🎯 Development Workflow

### Backend (Go)

1. **Write Code**: Follow Go best practices and conventions
2. **Format Code**: Run `gofmt -s -w .` before committing
3. **Run Tests**: Ensure all tests pass
   ```bash
   go test ./...
   go test -race ./...  # Check for race conditions
   ```
4. **Build**: Verify the build works
   ```bash
   go build
   ```

### Frontend (Vue 3)

1. **Development Server**: Run the dev server
   ```bash
   cd frontend
   npm run dev
   ```
2. **Format Code**: Format with Prettier
   ```bash
   npm run format
   ```
3. **Type Check**: Ensure TypeScript types are correct
   ```bash
   npm run type-check
   ```
4. **Build**: Test the production build
   ```bash
   npm run build
   ```

## 🌠 Creating Pull Requests

When submitting a PR:

### ✅ Checklist

- [ ] Code is formatted (`gofmt` for Go, `prettier` for frontend)
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation is updated (if needed)
- [ ] Commit messages are clear and descriptive
- [ ] PR description explains what and why

### 📝 PR Guidelines

- **Title**: Use a clear, descriptive title
  - ✅ `feat: add support for OpenAPI 3.1`
  - ✅ `fix: resolve theme switching bug`
  - ✅ `docs: update installation guide`
  
- **Description**: Explain:
  - What changes you made
  - Why you made them
  - How to test them
  - Any breaking changes

- **Size**: Keep PRs focused and reasonably sized
  - Large changes? Consider breaking into smaller PRs
  - Each PR should address one concern

## 👁️ Code Review

All PRs go through code review. We need to focus on:

### 🧙‍♀️ Code Quality

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write clear, self-documenting code
- Add comments for complex logic
- Handle errors properly
- Avoid unnecessary complexity

### 📝 Commit Messages

Write meaningful commit messages:

```
feat: add dark mode support for API explorer

- Implement theme toggle in header
- Add CSS variables for theme colors
- Persist theme preference in localStorage

Closes #123
```

Follow this format:
- **Type**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- **Subject**: Brief description (50 chars or less)
- **Body**: Detailed explanation (if needed)
- **Footer**: Reference issues/PRs

### 🔍 Review Process

1. A maintainer will review your PR
2. Address any feedback or requested changes
3. Once approved, your PR will be merged!

## 🎨 Project Structure

```
coco/
├── coco.go              # Main API entry point
├── internal/            # Internal packages
│   ├── handler.go       # HTTP handler logic
│   ├── config/          # Configuration
│   └── assets/          # Embedded frontend assets
├── frontend/            # Vue 3 frontend
│   ├── src/             # Source code
│   ├── vite.config.ts   # Build configuration
│   └── package.json     # Dependencies
├── example/             # Example integrations
│   └── framework/       # Framework examples
├── docs/                # Documentation
│   ├── en/              # English docs
│   └── zh/              # Chinese docs
└── README.md            # Project readme
```

## 🌍 Adding Translations

Want to add support for a new language?

1. Add translations in `frontend/src/i18n/`
2. Update language selector in UI
3. Update documentation to mention the new language
4. Test thoroughly

## 🧪 Testing

We value tests! When adding features:

- Write unit tests for new functions
- Add integration tests for new features
- Ensure existing tests still pass
- Aim for good test coverage

## 📚 Documentation

Good documentation is crucial:

- Update README if adding features
- Add code comments for complex logic
- Update API documentation
- Add examples for new features
- Keep docs in sync with code

## 💡 Tips for Success

- **Start Small**: Begin with small contributions to get familiar
- **Ask Questions**: Don't hesitate to ask for help
- **Be Patient**: The author still has studies, reviews take time, I will handle your PR
- **Stay Updated**: Pull latest changes regularly
- **Have Fun**: Enjoy the process of contributing!

## 🙏 Thank You!

Every contribution, no matter how small, makes Coco better. Thank you for being part of Coco!

If you have ideas to improve this guide, please let me know. Happy contributing! 🌟

---

**Questions?** Submit an **issue** or start a discussion. I'm here to help!
