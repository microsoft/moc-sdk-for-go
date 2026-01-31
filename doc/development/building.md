# Building from Source

This guide explains how to build the MOC SDK for Go from source.

## Prerequisites

- Go 1.25.0 or later
- Git
- mingw-w64 (for Windows DLL wrapper)
- Make

## Clone Repository

```bash
git clone https://github.com/microsoft/moc-sdk-for-go.git
cd moc-sdk-for-go
```

## Build Steps

### 1. Install Dependencies

```bash
go mod download
go mod tidy
```

### 2. Build All Packages

```bash
make
```

This runs:
- `make vendor` - Updates dependencies
- `make format` - Formats code
- `make build` - Builds packages and wrapper
- `make unittest` - Runs unit tests

### Individual Build Targets

```bash
# Update dependencies
make vendor

# Format code
make format

# Build packages only
make build

# Build without wrapper
GOARCH=amd64 go build -v ./...

# Run all tests
make test

# Run unit tests only
make unittest

# Run linter
make golangci-lint

# Clean build artifacts
make clean
```

## Building the C++ Wrapper

### Install mingw-w64

**Ubuntu/WSL:**
```bash
sudo apt-get update
sudo apt-get install mingw-w64
```

**macOS:**
```bash
brew install mingw-w64
```

**Windows (Chocolatey):**
```powershell
choco install mingw
```

### Build Wrapper

```bash
make build
```

Output:
- `bin/MocCppWrapper.dll` - Windows DLL
- `bin/MocCppWrapper.lib` - Import library

## Development Build

For development with hot reload:

```bash
# Watch and rebuild on changes (requires entr or similar)
find . -name "*.go" | entr -r make build
```

## Build Configuration

### Environment Variables

```bash
# Private repository access
export GOPRIVATE=github.com/microsoft

# Enable Go modules
export GO111MODULE=on

# Build tags
go build -tags debug ./...
```

### Build Flags

```bash
# Verbose output
go build -v ./...

# With race detector
go build -race ./...

# Optimized build
go build -ldflags="-s -w" ./...
```

## Cross-Compilation

### For Windows

```bash
GOOS=windows GOARCH=amd64 go build ./...
```

### For Linux

```bash
GOOS=linux GOARCH=amd64 go build ./...
```

## Troubleshooting

### "cannot find package"

```bash
go mod download
go mod tidy
```

### "mingw-w64 not found"

Install mingw-w64 and ensure it's in PATH:

```bash
which x86_64-w64-mingw32-gcc
```

### Build fails on Windows

Use WSL or ensure mingw-w64 is properly installed.

## Next Steps

- [Testing](testing.md) - Run tests
- [Contributing](contributing.md) - Contribute code
- [CI/CD](ci-cd.md) - Understand pipelines
