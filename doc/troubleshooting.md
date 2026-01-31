# Troubleshooting Guide

Common issues and solutions when using the MOC SDK for Go.

## Connection Issues

### "connection refused" Error

**Problem:** Cannot connect to MOC server.

**Possible Causes:**
1. MOC server not running
2. Incorrect server address
3. Firewall blocking connection
4. Wrong port number

**Solutions:**

```bash
# Test connectivity
ping moc-server.example.com

# Test port connectivity
telnet moc-server.example.com 55000

# On Windows PowerShell
Test-NetConnection -ComputerName moc-server.example.com -Port 55000
```

Check your server address:
```go
// Correct
cloudFQDN := "moc-server.example.com"  // Uses default port 55000

// With custom port
cloudFQDN := "moc-server.example.com:55001"
```

### "context deadline exceeded" Error

**Problem:** Operation timed out.

**Solutions:**

Increase timeout:
```go
// Increase from 30s to 2 minutes
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()
```

Check network latency:
```bash
ping -c 10 moc-server.example.com
```

## Authentication Issues

### "certificate verify failed" Error

**Problem:** Server certificate validation failed.

**Solutions:**

1. **Verify CA certificate is correct:**
```bash
openssl verify -CAfile ca-cert.pem client-cert.pem
```

2. **Check certificate expiration:**
```bash
openssl x509 -in client-cert.pem -noout -dates
openssl x509 -in ca-cert.pem -noout -dates
```

3. **Verify certificate chain:**
```bash
openssl verify -CAfile ca-cert.pem -untrusted intermediate.pem client-cert.pem
```

4. **For development, use debug mode (NOT for production):**
```go
import "os"
os.Setenv("WSSD_DEBUG_MODE", "on")
```

### "tls: bad certificate" Error

**Problem:** Server rejected client certificate.

**Solutions:**

1. **Ensure certificate is signed by trusted CA**
2. **Check private key matches certificate:**
```bash
# Get certificate modulus
openssl x509 -noout -modulus -in client-cert.pem | openssl md5

# Get private key modulus
openssl rsa -noout -modulus -in client-key.pem | openssl md5

# They should match
```

3. **Verify certificate format (PEM):**
```bash
openssl x509 -in client-cert.pem -text -noout
```

### "permission denied" Error

**Problem:** Authenticated but not authorized.

**Solutions:**

1. **Check RBAC role assignments:**
```go
roleClient, _ := roleassignment.NewRoleAssignmentClient(cloudFQDN, authorizer)
assignments, _ := roleClient.Get(ctx, group, "")

for _, assignment := range *assignments {
    fmt.Printf("Principal: %s, Role: %s\n", 
        *assignment.Properties.PrincipalID,
        *assignment.Properties.RoleDefinitionID)
}
```

2. **Contact MOC administrator to grant appropriate permissions**

## SDK Issues

### "module not found" Error

**Problem:** Cannot find SDK module.

**Solutions:**

1. **Add module to go.mod:**
```bash
go get github.com/microsoft/moc-sdk-for-go
```

2. **Set GOPRIVATE for private repos:**
```bash
export GOPRIVATE=github.com/microsoft
```

3. **Configure Git authentication:**
```bash
# Using personal access token
git config --global url."https://YOUR_TOKEN@github.com/".insteadOf "https://github.com/"

# Using SSH
git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
```

4. **Run go mod tidy:**
```bash
go mod tidy
```

### Build Errors with C++ Wrapper

**Problem:** Cannot build Windows DLL.

**Solutions:**

1. **Install mingw-w64:**
```bash
# WSL/Ubuntu
sudo apt-get install mingw-w64

# Windows with Chocolatey
choco install mingw

# macOS with Homebrew
brew install mingw-w64
```

2. **Verify compiler is in PATH:**
```bash
which x86_64-w64-mingw32-gcc
# Should show path to compiler
```

3. **Build only Go packages (skip wrapper):**
```bash
GOARCH=amd64 go build -v ./...
```

## Runtime Issues

### "resource not found" Error

**Problem:** Resource doesn't exist.

**Solutions:**

1. **Verify resource name and group:**
```go
// List all resources
resources, err := client.Get(ctx, groupName, "")
for _, r := range *resources {
    fmt.Printf("Found: %s\n", *r.Name)
}
```

2. **Check spelling and case sensitivity**

3. **Handle not found gracefully:**
```go
import "github.com/microsoft/moc/pkg/errors"

vm, err := vmClient.Get(ctx, group, name)
if err != nil {
    if errors.IsNotFound(err) {
        // Create resource
        vm, err = vmClient.CreateOrUpdate(ctx, group, name, vmSpec)
    } else {
        return err
    }
}
```

### "resource already exists" Error

**Problem:** Trying to create a resource that exists.

**Solutions:**

1. **Use CreateOrUpdate instead of Create:**
```go
// ✅ Good: Updates if exists
vm, err := vmClient.CreateOrUpdate(ctx, group, name, vmSpec)

// ❌ Bad: Fails if exists
vm, err := vmClient.Create(ctx, group, name, vmSpec)
```

2. **Handle already exists:**
```go
vm, err := vmClient.CreateOrUpdate(ctx, group, name, vmSpec)
if err != nil {
    if errors.IsAlreadyExists(err) {
        // Get existing resource
        vm, err = vmClient.Get(ctx, group, name)
    } else {
        return err
    }
}
```

### Memory Leaks or Connection Exhaustion

**Problem:** Too many open connections.

**Solutions:**

1. **Reuse clients:**
```go
// ✅ Good: Create once, reuse
vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
for _, name := range vmNames {
    vm, _ := vmClient.Get(ctx, group, name)
}

// ❌ Bad: Create in loop
for _, name := range vmNames {
    vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    vm, _ := vmClient.Get(ctx, group, name)
}
```

2. **Clear connection cache when done:**
```go
import "github.com/microsoft/moc-sdk-for-go/pkg/client"

defer client.ClearConnectionCache()
```

3. **Always use context with timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

## Performance Issues

### Slow Operations

**Problem:** Operations taking too long.

**Solutions:**

1. **Check network latency:**
```bash
ping -c 100 moc-server.example.com
```

2. **Use appropriate timeouts:**
```go
// For quick operations
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// For VM creation
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
```

3. **Batch operations when possible:**
```go
// Get all VMs at once
vms, _ := vmClient.Get(ctx, group, "")

// vs. Getting one by one
```

### High Memory Usage

**Problem:** Application using too much memory.

**Solutions:**

1. **Process resources in batches:**
```go
const batchSize = 100

vms, _ := vmClient.Get(ctx, group, "")
for i := 0; i < len(*vms); i += batchSize {
    end := i + batchSize
    if end > len(*vms) {
        end = len(*vms)
    }
    batch := (*vms)[i:end]
    // Process batch
}
```

2. **Clear unused resources:**
```go
vms = nil
runtime.GC()
```

## Debugging

### Enable Verbose Logging

```go
import "k8s.io/klog"

func init() {
    klog.InitFlags(nil)
    flag.Set("v", "4") // Set verbosity level
    flag.Parse()
}
```

### Inspect gRPC Communication

```go
import "google.golang.org/grpc/grpclog"

func init() {
    grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
}
```

### Debug Mode

Enable debug mode for development (disables TLS):

```go
os.Setenv("WSSD_DEBUG_MODE", "on")
```

Or using viper:
```go
viper.Set("Debug", true)
```

**Warning:** Never use debug mode in production!

## Getting Help

If you cannot resolve the issue:

1. **Check the FAQ:** [FAQ](faq.md)
2. **Search existing issues:** [GitHub Issues](https://github.com/microsoft/moc-sdk-for-go/issues)
3. **Open a new issue** with:
   - SDK version
   - Go version
   - Error messages
   - Code snippet
   - Steps to reproduce

## Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| `connection refused` | Server not reachable | Check network, firewall, server status |
| `certificate verify failed` | Invalid certificate | Verify certificate chain and expiration |
| `context deadline exceeded` | Operation timeout | Increase timeout or check network |
| `permission denied` | Insufficient permissions | Check RBAC role assignments |
| `resource not found` | Resource doesn't exist | Verify name and group |
| `resource already exists` | Duplicate create | Use CreateOrUpdate |
| `invalid input` | Bad parameters | Validate input parameters |

## Next Steps

- [FAQ](faq.md) - Frequently asked questions
- [Getting Started](getting-started.md) - Basic usage
- [GitHub Issues](https://github.com/microsoft/moc-sdk-for-go/issues) - Report issues
