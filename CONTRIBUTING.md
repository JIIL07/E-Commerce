# Contributing to E-Commerce Skeleton

Thank you for your interest in contributing to the E-Commerce Skeleton project! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.23+
- Node.js 18+
- Make (optional)

### Development Setup
1. Fork the repository
2. Clone your fork: `git clone https://github.com/JIIL07/E-Commerce`
3. Setup environment: `make setup`
4. Start development: `make dev`

## 📝 How to Contribute

### Reporting Issues
- Use the GitHub issue tracker
- Provide clear description and steps to reproduce
- Include system information (OS, Docker version, etc.)

### Suggesting Features
- Open an issue with the "enhancement" label
- Describe the feature and its benefits
- Consider implementation complexity

### Code Contributions
1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes
3. Test thoroughly
4. Commit with clear messages
5. Push to your fork
6. Create a Pull Request

## 🧪 Testing

### Backend Testing
```bash
cd backend-go
go test ./...
```

### Frontend Testing
```bash
cd frontend
npm test
```

### Integration Testing
```bash
make test
```

## 📋 Code Standards

### Go Code
- Follow Go conventions
- Use meaningful variable names
- Add comments for public functions
- Run `gofmt` and `golint`

### TypeScript/React Code
- Use TypeScript strict mode
- Follow React best practices
- Use meaningful component names
- Add PropTypes or TypeScript interfaces

### Git Commits
- Use conventional commit format
- Keep commits atomic
- Write clear commit messages

## 🏗️ Architecture Guidelines

### Backend Structure
- Follow clean architecture principles
- Separate concerns (handlers, services, repositories)
- Use dependency injection
- Implement proper error handling

### Frontend Structure
- Use component composition
- Implement proper state management
- Follow accessibility guidelines
- Optimize for performance

## 📚 Documentation

### Code Documentation
- Document public APIs
- Add inline comments for complex logic
- Update README files when needed

### API Documentation
- Use Swagger/OpenAPI annotations
- Provide example requests/responses
- Document error codes

## 🔒 Security

### Security Guidelines
- Never commit secrets or API keys
- Use environment variables for configuration
- Implement proper input validation
- Follow OWASP guidelines

### Reporting Security Issues
- Email security issues to: security@example.com
- Do not open public issues for security vulnerabilities

## 🎯 Areas for Contribution

### High Priority
- Bug fixes
- Performance improvements
- Security enhancements
- Documentation improvements

### Medium Priority
- New features
- UI/UX improvements
- Test coverage
- Code refactoring

### Low Priority
- Code style improvements
- Minor optimizations
- Additional examples


Thank you for contributing to E-Commerce Skeleton! 🎉
