# Contributing to MOC SDK for Go

Thank you for your interest in contributing! This guide will help you get started.

## Code of Conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).

## Contributor License Agreement (CLA)

You must sign a [Contributor License Agreement](https://cla.opensource.microsoft.com) before your PR will be merged. This is a one-time requirement.

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR-USERNAME/moc-sdk-for-go.git
   cd moc-sdk-for-go
   ```

3. **Set up upstream remote**
   ```bash
   git remote add upstream https://github.com/microsoft/moc-sdk-for-go.git
   ```

4. **Install dependencies**
   ```bash
   make vendor
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-new-feature
```

### 2. Make Changes

Follow the coding standards:
- Use `gofmt` for formatting
- Follow Go best practices
- Add tests for new code
- Update documentation

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run unit tests
make unittest

# Run linter
make golangci-lint

# Format code
make format
```

### 4. Commit Changes

```bash
git add .
git commit -m "Add feature: description of changes"
```

**Commit Message Guidelines:**
- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and PRs in the body

### 5. Push Changes

```bash
git push origin feature/my-new-feature
```

### 6. Create Pull Request

1. Go to your fork on GitHub
2. Click "Pull Request"
3. Select your branch
4. Fill out the PR template
5. Submit

## Pull Request Guidelines

- **Title**: Clear and descriptive
- **Description**: Explain what and why
- **Tests**: Include tests for new functionality
- **Documentation**: Update relevant docs
- **Commits**: Keep commits focused and logical

### PR Checklist

- [ ] Code builds successfully
- [ ] All tests pass
- [ ] Code is formatted (`make format`)
- [ ] Linter passes (`make golangci-lint`)
- [ ] Documentation updated
- [ ] CHANGELOG updated (if applicable)
- [ ] CLA signed

## Coding Standards

### Go Style

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

### Naming Conventions

```go
// Exported functions and types: PascalCase
func CreateVirtualMachine() {}
type VirtualMachine struct {}

// Unexported: camelCase
func internalHelper() {}
type internalState struct {}

// Constants: PascalCase
const MaxRetries = 3
```

### Error Handling

```go
// Good: Return errors, don't panic
func doSomething() error {
    if err := operation(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}

// Bad: Panic
func doSomething() {
    if err := operation(); err != nil {
        panic(err)
    }
}
```

### Comments

```go
// Package comment at top of file
// Package virtualmachine provides...
package virtualmachine

// Exported function must have comment
// CreateVirtualMachine creates a new virtual machine.
func CreateVirtualMachine() {}
```

## Testing Guidelines

### Unit Tests

```go
func TestCreateVirtualMachine(t *testing.T) {
    // Arrange
    client := setupTestClient()
    
    // Act
    vm, err := client.CreateOrUpdate(ctx, group, name, vmSpec)
    
    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if vm == nil {
        t.Fatal("expected VM, got nil")
    }
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid input", "valid", true, false},
        {"invalid input", "", false, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use complete sentences
- Provide examples for complex functionality

### README and Guides

- Keep documentation up to date
- Add examples for new features
- Update API reference

## Review Process

1. **Automated Checks**: CI/CD runs automatically
2. **Code Review**: Maintainers review your code
3. **Feedback**: Address review comments
4. **Approval**: Two approvals required
5. **Merge**: Maintainer merges PR

## Communication

- **Issues**: Report bugs and request features
- **Discussions**: Ask questions and discuss ideas
- **Pull Requests**: Submit code changes

## Resources

- [GitHub Repository](https://github.com/microsoft/moc-sdk-for-go)
- [Building Guide](building.md)
- [Testing Guide](testing.md)
- [CI/CD Documentation](ci-cd.md)

## Questions?

If you have questions, please:
1. Check existing documentation
2. Search existing issues
3. Open a new issue

Thank you for contributing! 🎉
