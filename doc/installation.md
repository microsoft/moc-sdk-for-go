# Installation Guide

This guide walks you through installing and setting up the MOC SDK for Go in your project.

## Prerequisites

Before installing the MOC SDK, ensure you have:

- **Go 1.25.0 or later** installed
- **Git** for version control
- **Access to GitHub** (for private repository access)
- **MOC backend service** endpoint available
- **Authentication credentials** for the MOC service

## Installation Methods

### Method 1: Go Modules (Recommended)

The MOC SDK uses Go modules for dependency management. Add it to your project using `go get`:

```bash
go get github.com/microsoft/moc-sdk-for-go
```

Or add it directly to your `go.mod` file:

```go
module your-project

go 1.25.0

require (
    github.com/microsoft/moc-sdk-for-go v0.39.0
)
```

Then run:

```bash
go mod download
go mod tidy
```

### Method 2: Clone Repository

For development or contributing to the SDK:

```bash
# Clone the repository
git clone https://github.com/microsoft/moc-sdk-for-go.git
cd moc-sdk-for-go

# Install dependencies
make vendor

# Build the SDK
make build
```

## Private Repository Access

The MOC SDK depends on private Microsoft repositories. Configure Go to access them:

```bash
# Set GOPRIVATE environment variable
export GOPRIVATE=github.com/microsoft

# Configure Git credentials for private repos
git config --global url."https://YOUR_TOKEN@github.com/".insteadOf "https://github.com/"
```

Or use SSH:

```bash
git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
```

## Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/microsoft/moc-sdk-for-go/pkg/client"
)

func main() {
    fmt.Println("MOC SDK for Go installed successfully!")
    
    // Check if client package is accessible
    _ = client.ServerPort
    fmt.Printf("Server port: %d\n", client.ServerPort)
}
```

Run the test:

```bash
go run main.go
```

Expected output:
```
MOC SDK for Go installed successfully!
Server port: 55000
```

## Building the C++ Wrapper (Optional)

If you need the C++ wrapper for cross-language interoperability:

### Prerequisites for Windows

Install mingw-w64 compiler:

**On Windows (with Chocolatey):**
```powershell
choco install mingw
```

**On WSL/Ubuntu:**
```bash
sudo apt-get install mingw-w64
```

### Build the Wrapper

```bash
# Build everything including the C++ wrapper
make

# Or build only the wrapper
make build
```

This generates:
- `bin/MocCppWrapper.dll` - Windows DLL
- `bin/MocCppWrapper.lib` - Import library

## Project Structure

After installation, your project structure should look like:

```
your-project/
├── go.mod                  # Go module file
├── go.sum                  # Dependency checksums
├── main.go                 # Your application code
└── vendor/                 # Vendored dependencies (optional)
    └── github.com/
        └── microsoft/
            ├── moc-sdk-for-go/
            └── moc/
```

## Dependencies

The MOC SDK requires these key dependencies (automatically installed):

- **github.com/microsoft/moc** - MOC protocol definitions and common libraries
- **google.golang.org/grpc** - gRPC framework for communication
- **github.com/spf13/viper** - Configuration management
- **k8s.io/klog** - Logging framework
- **github.com/Azure/go-autorest** - Azure REST client infrastructure

## Environment Variables

Configure these optional environment variables:

```bash
# Enable debug mode (disable TLS for testing)
export WSSD_DEBUG_MODE=on

# Set private repository access
export GOPRIVATE=github.com/microsoft

# Enable Go modules
export GO111MODULE=on
```

## Configuration File

Create a configuration file for SDK settings:

```yaml
# config.yaml
debug: false
server:
  address: "your-moc-server.example.com"
  port: 55000
auth:
  port: 65000
```

Load the configuration in your application:

```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()
```

## Troubleshooting

### Issue: "module github.com/microsoft/moc-sdk-for-go: reading... 404 Not Found"

**Solution:** Configure private repository access (see above) and ensure you have proper GitHub authentication.

### Issue: "could not import github.com/microsoft/moc"

**Solution:** The MOC dependency is private. Set `GOPRIVATE=github.com/microsoft` and configure Git credentials.

### Issue: "mingw-w64 not found" (Windows DLL build)

**Solution:** Install mingw-w64 compiler:
```bash
# WSL/Ubuntu
sudo apt-get install mingw-w64

# Windows with Chocolatey
choco install mingw
```

### Issue: Build errors with Go version mismatch

**Solution:** Ensure you're using Go 1.25.0 or later:
```bash
go version
```

## Next Steps

Now that the SDK is installed:

1. Read the [Getting Started Guide](getting-started.md) for your first application
2. Learn about [Authentication](authentication.md) to connect to MOC services
3. Explore [Client Usage](client-usage.md) for detailed client configuration

## Additional Resources

- [Building from Source](development/building.md) - Detailed build instructions
- [Development Setup](development/contributing.md) - Contributing to the SDK
- [Makefile Reference](../Makefile) - Available make targets
