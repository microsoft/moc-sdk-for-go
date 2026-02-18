# Client Usage Guide

This guide covers how to initialize, configure, and use MOC SDK clients effectively.

## Client Types

The SDK provides two approaches to client usage:

1. **Individual Service Clients** - Direct access to specific services
2. **Facade Clients** - Aggregated clients for convenience

## Individual Service Clients

### Creating Service Clients

Each service has its own client that can be created independently:

```go
import (
    "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
    "github.com/microsoft/moc/pkg/auth"
)

func main() {
    // Create authorizer
    authorizer, err := auth.NewAuthorizerFromCertificate(
        "certs/client.pem",
        "certs/client-key.pem",
        "certs/ca.pem",
        "",
    )
    if err != nil {
        panic(err)
    }
    
    // Create VM client
    cloudFQDN := "moc-server.example.com"
    vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    if err != nil {
        panic(err)
    }
    
    // Use the client
    vms, err := vmClient.Get(ctx, "default", "")
}
```

### Available Service Clients

**Compute Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachinescaleset"
import "github.com/microsoft/moc-sdk-for-go/services/compute/galleryimage"
import "github.com/microsoft/moc-sdk-for-go/services/compute/baremetalhost"
import "github.com/microsoft/moc-sdk-for-go/services/compute/availabilityset"
```

**Network Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/network/virtualnetwork"
import "github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"
import "github.com/microsoft/moc-sdk-for-go/services/network/loadbalancer"
import "github.com/microsoft/moc-sdk-for-go/services/network/publicipaddress"
import "github.com/microsoft/moc-sdk-for-go/services/network/networksecuritygroup"
```

**Storage Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/storage/virtualharddisk"
import "github.com/microsoft/moc-sdk-for-go/services/storage/container"
```

**Security Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/security/identity"
import "github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
import "github.com/microsoft/moc-sdk-for-go/services/security/certificate"
import "github.com/microsoft/moc-sdk-for-go/services/security/role"
import "github.com/microsoft/moc-sdk-for-go/services/security/roleassignment"
```

**Cloud Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/location"
import "github.com/microsoft/moc-sdk-for-go/services/cloud/zone"
import "github.com/microsoft/moc-sdk-for-go/services/cloud/node"
import "github.com/microsoft/moc-sdk-for-go/services/cloud/group"
```

**Admin Services:**
```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/version"
import "github.com/microsoft/moc-sdk-for-go/services/admin/health"
import "github.com/microsoft/moc-sdk-for-go/services/admin/logging"
import "github.com/microsoft/moc-sdk-for-go/services/admin/recovery"
import "github.com/microsoft/moc-sdk-for-go/services/admin/validation"
```

## Facade Clients

Facade clients aggregate related services for convenience:

### Using Facade Clients

```go
import "github.com/microsoft/moc-sdk-for-go/pkg/client"

func main() {
    authorizer, _ := createAuthorizer()
    cloudFQDN := "moc-server.example.com"
    
    // Create compute facade client
    computeClient, err := client.NewComputeClient(cloudFQDN, authorizer)
    if err != nil {
        panic(err)
    }
    
    // Access individual services through facade
    vms, err := computeClient.VirtualMachines.Get(ctx, "default", "")
    images, err := computeClient.GalleryImages.Get(ctx, "default", "")
    vmss, err := computeClient.VirtualMachineScaleSets.Get(ctx, "default", "")
}
```

### Available Facade Clients

```go
// Compute facade
computeClient, err := client.NewComputeClient(cloudFQDN, authorizer)

// Network facade
networkClient, err := client.NewNetworkClient(cloudFQDN, authorizer)

// Storage facade
storageClient, err := client.NewStorageClient(cloudFQDN, authorizer)

// Security facade
securityClient, err := client.NewSecurityClient(cloudFQDN, authorizer)

// Cloud facade
cloudClient, err := client.NewCloudClient(cloudFQDN, authorizer)

// Admin facade
adminClient, err := client.NewAdminClient(cloudFQDN, authorizer)
```

## Context Usage

All client operations accept a `context.Context` for timeout and cancellation control.

### With Timeout

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

vm, err := vmClient.Get(ctx, "default", "my-vm")
if err != nil {
    // Handle error (including timeout)
}
```

### With Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

// Cancel from another goroutine
go func() {
    time.Sleep(5 * time.Second)
    cancel() // Cancel the operation
}()

vm, err := vmClient.CreateOrUpdate(ctx, "default", "my-vm", vmSpec)
```

### Best Practices

```go
// ✅ Good: Always use timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// ❌ Bad: Using background context without timeout
ctx := context.Background() // No timeout!
```

## Resource Operations

### Get Resource

Get a specific resource by name:

```go
ctx := context.Background()
groupName := "production"
vmName := "web-server-01"

vm, err := vmClient.Get(ctx, groupName, vmName)
if err != nil {
    log.Fatalf("Failed to get VM: %v", err)
}

fmt.Printf("VM: %s, Status: %s\n", *vm.Name, vm.Statuses)
```

### List Resources

List all resources in a group (empty name):

```go
groupName := "production"
emptyName := "" // Empty string lists all

vms, err := vmClient.Get(ctx, groupName, emptyName)
if err != nil {
    log.Fatalf("Failed to list VMs: %v", err)
}

for _, vm := range *vms {
    fmt.Printf("VM: %s\n", *vm.Name)
}
```

### Create or Update Resource

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute"

vmSpec := &compute.VirtualMachine{
    Name:     stringPtr("my-vm"),
    Location: stringPtr("default"),
    Properties: &compute.VirtualMachineProperties{
        HardwareProfile: &compute.HardwareProfile{
            VMSize: &compute.VMSize{
                VCPUs:    int32Ptr(2),
                MemoryMB: int32Ptr(4096),
            },
        },
        StorageProfile: &compute.StorageProfile{
            ImageReference: &compute.ImageReference{
                Name: stringPtr("ubuntu-20.04"),
            },
        },
        OsProfile: &compute.OSProfile{
            ComputerName:  stringPtr("my-vm"),
            AdminUsername: stringPtr("azureuser"),
        },
        NetworkProfile: &compute.NetworkProfile{
            NetworkInterfaces: &[]compute.NetworkInterfaceReference{
                {ID: stringPtr("/production/networkinterfaces/my-nic")},
            },
        },
    },
}

vm, err := vmClient.CreateOrUpdate(ctx, "production", "my-vm", vmSpec)
if err != nil {
    log.Fatalf("Failed to create VM: %v", err)
}

fmt.Printf("Created VM: %s\n", *vm.Name)

// Helper functions
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }
```

### Delete Resource

```go
err := vmClient.Delete(ctx, "production", "my-vm")
if err != nil {
    log.Fatalf("Failed to delete VM: %v", err)
}

fmt.Println("VM deleted successfully")
```

### Query Resources

Query with filters:

```go
query := "status eq 'Running'"
vms, err := vmClient.Query(ctx, "production", query)
if err != nil {
    log.Fatalf("Failed to query VMs: %v", err)
}

fmt.Printf("Found %d running VMs\n", len(*vms))
```

## Resource Lifecycle Operations

Some resources support lifecycle operations:

### Virtual Machine Lifecycle

```go
// Start VM
err := vmClient.Start(ctx, "production", "my-vm")

// Stop VM (forced)
err := vmClient.Stop(ctx, "production", "my-vm")

// Stop VM (graceful shutdown)
err := vmClient.StopGraceful(ctx, "production", "my-vm")

// Pause VM
err := vmClient.Pause(ctx, "production", "my-vm")

// Save VM state
err := vmClient.Save(ctx, "production", "my-vm")
```

### Checking Resource Status

```go
vms, err := vmClient.Get(ctx, "production", "my-vm")
if err != nil {
    log.Fatalf("Failed to get VM: %v", err)
}

vm := (*vms)[0]
if vm.Statuses != nil {
    fmt.Printf("VM Status: %s\n", *vm.Statuses)
}
```

## Error Handling

### Checking Error Types

```go
import "github.com/microsoft/moc/pkg/errors"

vm, err := vmClient.Get(ctx, "production", "my-vm")
if err != nil {
    if errors.IsNotFound(err) {
        fmt.Println("VM not found")
        // Handle not found case
    } else if errors.IsAlreadyExists(err) {
        fmt.Println("VM already exists")
        // Handle already exists case
    } else if errors.IsInvalidInput(err) {
        fmt.Println("Invalid input parameters")
        // Handle invalid input
    } else {
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

### Retry Logic

```go
import "time"

func retryOperation(maxRetries int, operation func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = operation()
        if err == nil {
            return nil
        }
        
        // Don't retry certain errors
        if errors.IsNotFound(err) || errors.IsInvalidInput(err) {
            return err
        }
        
        // Wait before retry
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return err
}

// Usage
err := retryOperation(3, func() error {
    _, err := vmClient.Get(ctx, "production", "my-vm")
    return err
})
```

## Connection Management

### Connection Caching

The SDK automatically caches connections:

```go
// First call creates connection
vm1Client, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)

// Second call reuses cached connection
vm2Client, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
```

### Clearing Connection Cache

```go
import "github.com/microsoft/moc-sdk-for-go/pkg/client"

// Clear all cached connections
client.ClearConnectionCache()
```

### Connection Validation

Connections are validated before reuse:

```go
// SDK automatically checks connection state
// - TransientFailure: marked invalid
// - Shutdown: marked invalid
// - Other states: valid
```

## Configuration

### Debug Mode

Enable debug mode for development:

```go
import "os"

// Enable debug mode (disables TLS)
os.Setenv("WSSD_DEBUG_MODE", "on")

// Or using viper
import "github.com/spf13/viper"
viper.Set("Debug", true)
```

### Custom Ports

Specify custom server ports:

```go
// With port
cloudFQDN := "moc-server.example.com:55001"

// Default port (55000)
cloudFQDN := "moc-server.example.com"
```

## Best Practices

### 1. Reuse Clients

```go
// ✅ Good: Create once, reuse
vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
for _, vmName := range vmNames {
    vm, _ := vmClient.Get(ctx, group, vmName)
}

// ❌ Bad: Create client in loop
for _, vmName := range vmNames {
    vmClient, _ := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    vm, _ := vmClient.Get(ctx, group, vmName)
}
```

### 2. Always Use Context

```go
// ✅ Good: Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// ❌ Bad: No timeout
ctx := context.Background()
```

### 3. Handle Errors Appropriately

```go
// ✅ Good: Check error types
if errors.IsNotFound(err) {
    // Create resource
} else if err != nil {
    return err
}

// ❌ Bad: Ignore errors
vm, _ := vmClient.Get(ctx, group, name)
```

### 4. Use Defer for Cleanup

```go
// ✅ Good: Ensure cancel is called
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Continue with operation
```

## Next Steps

- [Service Documentation](services/compute.md) - Explore specific services
- [Code Examples](examples/vm-management.md) - Detailed examples
- [Error Handling](advanced/error-handling.md) - Advanced error handling
