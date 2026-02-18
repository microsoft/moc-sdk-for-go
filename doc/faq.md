# Frequently Asked Questions (FAQ)

## General Questions

### What is the MOC SDK for Go?

The MOC SDK for Go is a client library for interacting with Microsoft's Management of Cloud (MOC) infrastructure services. It provides Go APIs for managing virtual machines, networks, storage, security, and other cloud resources.

### What version of Go is required?

Go 1.25.0 or later is required. Check your Go version:

```bash
go version
```

### Is the SDK production-ready?

Yes, the SDK is used in production environments. Follow best practices for authentication, error handling, and resource management.

### Where can I find the source code?

The SDK is open source on GitHub:
- SDK: https://github.com/microsoft/moc-sdk-for-go
- Protocol Definitions: https://github.com/microsoft/moc

## Installation & Setup

### How do I install the SDK?

```bash
go get github.com/microsoft/moc-sdk-for-go
```

See the [Installation Guide](installation.md) for details.

### Do I need access to private repositories?

Yes, the SDK depends on the private `github.com/microsoft/moc` repository. Configure access:

```bash
export GOPRIVATE=github.com/microsoft
```

### How do I build the C++ wrapper?

Install mingw-w64 and run `make`:

```bash
# Ubuntu/WSL
sudo apt-get install mingw-w64
make

# The DLL is created in bin/MocCppWrapper.dll
```

## Authentication

### What authentication methods are supported?

- Certificate-based authentication (recommended)
- Token-based authentication
- Environment-based authentication

See [Authentication Guide](authentication.md) for details.

### How do I create certificates?

See the [Authentication Guide](authentication.md#certificate-requirements) for certificate generation steps.

### Can I disable TLS for testing?

Yes, enable debug mode (development only):

```bash
export WSSD_DEBUG_MODE=on
```

**Never use debug mode in production!**

### How do I rotate credentials?

```go
// Create new authorizer
newAuth, _ := auth.NewAuthorizerFromCertificate(newCert, newKey, ca, "")

// Clear connection cache
client.ClearConnectionCache()

// Create new clients with new authorizer
vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, newAuth)
```

## Usage

### How do I create a virtual machine?

See the [Getting Started Guide](getting-started.md) and [Compute Services](services/compute.md) for examples.

### How do I list all resources?

Use an empty string for the name parameter:

```go
// List all VMs in a group
vms, err := vmClient.Get(ctx, "production", "")
```

### How do I handle errors?

```go
import "github.com/microsoft/moc/pkg/errors"

vm, err := vmClient.Get(ctx, group, name)
if err != nil {
    if errors.IsNotFound(err) {
        // Handle not found
    } else if errors.IsAlreadyExists(err) {
        // Handle already exists
    } else {
        // Handle other errors
        return err
    }
}
```

### Should I use individual clients or facade clients?

**Individual clients** - When you need specific service functionality:
```go
vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
```

**Facade clients** - For convenience when using multiple related services:
```go
computeClient, _ := client.NewComputeClient(cloudFQDN, authorizer)
// Access: computeClient.VirtualMachines, computeClient.GalleryImages, etc.
```

### How do I set timeouts?

Always use context with timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

vm, err := vmClient.Get(ctx, group, name)
```

## Architecture

### How does the SDK communicate with MOC?

The SDK uses gRPC over TLS for secure communication. See [Architecture](architecture.md) for details.

### What ports does the SDK use?

- **Port 55000**: Main service port (default)
- **Port 65000**: Authentication port

### Are connections pooled?

Yes, the SDK automatically caches and reuses connections. See [Connection Management](advanced/connection-management.md).

### Can I use multiple clients concurrently?

Yes, clients are safe for concurrent use from multiple goroutines.

## Resources

### How are resources organized?

Resources are organized by groups and names:

```go
// Pattern: /<group>/<resource-type>/<name>
vmClient.Get(ctx, "production", "web-server-01")
```

### Can I use tags?

Yes, most resources support tags:

```go
vmSpec.Tags = map[string]*string{
    "environment": stringPtr("production"),
    "owner":       stringPtr("team-platform"),
}
```

### How do I reference other resources?

Use resource IDs:

```go
vmSpec.Properties.NetworkProfile.NetworkInterfaces = &[]compute.NetworkInterfaceReference{
    {ID: stringPtr("/production/networkinterfaces/web-nic")},
}
```

## Performance

### How can I improve performance?

1. **Reuse clients** - Don't create new clients for each operation
2. **Batch operations** - List all resources at once instead of one by one
3. **Use appropriate timeouts** - Balance between responsiveness and reliability
4. **Clear cache when done** - Free connections when no longer needed

### Why are operations slow?

Check:
1. Network latency to MOC server
2. Server resource availability
3. Operation complexity (VM creation vs. simple query)
4. Timeout settings

## Troubleshooting

### Where can I find troubleshooting help?

See the [Troubleshooting Guide](troubleshooting.md) for common issues and solutions.

### How do I enable verbose logging?

```go
import "k8s.io/klog"

func init() {
    klog.InitFlags(nil)
    flag.Set("v", "4")
    flag.Parse()
}
```

### How do I debug gRPC issues?

Enable gRPC logging:

```go
import "google.golang.org/grpc/grpclog"

func init() {
    grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
}
```

## Contributing

### How can I contribute?

See the [Contributing Guide](development/contributing.md) for guidelines.

### How do I report bugs?

Open an issue on GitHub with:
- SDK version
- Go version
- Error messages
- Code snippet
- Steps to reproduce

### How do I build from source?

```bash
git clone https://github.com/microsoft/moc-sdk-for-go.git
cd moc-sdk-for-go
make
```

See [Building from Source](development/building.md) for details.

## Best Practices

### What are the recommended patterns?

1. **Always use context with timeout**
2. **Reuse clients** - Create once, use many times
3. **Handle errors appropriately** - Check error types
4. **Use least privilege** - Request minimum required permissions
5. **Tag resources** - For organization and cost tracking
6. **Store secrets in Key Vault** - Never hardcode credentials

### How should I organize my code?

```go
// Good structure
type CloudManager struct {
    vmClient      *virtualmachine.VirtualMachineClient
    networkClient *virtualnetwork.VirtualNetworkClient
    storageClient *virtualharddisk.VirtualHardDiskClient
}

func NewCloudManager(cloudFQDN string, authorizer auth.Authorizer) (*CloudManager, error) {
    vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    if err != nil {
        return nil, err
    }
    
    // ... create other clients
    
    return &CloudManager{
        vmClient: vmClient,
        // ...
    }, nil
}
```

### What should I avoid?

❌ **Don't:**
- Create clients in loops
- Hardcode credentials
- Ignore errors
- Use debug mode in production
- Skip context timeouts
- Forget to handle resource cleanup

## Support

### Where can I get help?

1. **Documentation**: https://github.com/microsoft/moc-sdk-for-go/tree/main/doc
2. **Issues**: https://github.com/microsoft/moc-sdk-for-go/issues
3. **Discussions**: Contact your MOC administrator

### How do I stay updated?

- Watch the GitHub repository for releases
- Check the changelog for updates
- Subscribe to GitHub notifications

## License

### What license is the SDK under?

Apache License 2.0. See the [LICENSE](../LICENSE) file.

### Can I use this in commercial projects?

Yes, the Apache License 2.0 allows commercial use.

## Next Steps

- [Getting Started](getting-started.md) - Build your first application
- [Architecture](architecture.md) - Understand the SDK design
- [Service Documentation](services/compute.md) - Explore available services
- [Examples](examples/vm-management.md) - Learn from examples
