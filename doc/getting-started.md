# Getting Started with MOC SDK for Go

This guide will help you build your first application using the MOC SDK for Go. By the end of this guide, you'll be able to connect to a MOC service and perform basic operations.

## Prerequisites

Before you begin, ensure you have:

- [Installed the MOC SDK](installation.md)
- Go 1.25.0 or later
- Access to a MOC backend service
- Authentication credentials (certificate or token)

## Your First MOC Application

### Step 1: Import Required Packages

Create a new Go file `main.go` and import the necessary packages:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/microsoft/moc-sdk-for-go/services/compute"
    "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
    "github.com/microsoft/moc/pkg/auth"
)
```

### Step 2: Initialize the Authorizer

The SDK requires an authorizer for authentication. Here's how to create one:

```go
func createAuthorizer() (auth.Authorizer, error) {
    // Option 1: Certificate-based authentication
    authorizer, err := auth.NewAuthorizerFromCertificate(
        "path/to/certificate.pem",  // Client certificate
        "path/to/key.pem",           // Private key
        "path/to/ca.pem",            // CA certificate
        "", // optional certificate password
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create authorizer: %w", err)
    }
    
    return authorizer, nil
}
```

### Step 3: Create a Service Client

Now create a client for the service you want to use. Let's start with virtual machines:

```go
func main() {
    // Create authorizer
    authorizer, err := createAuthorizer()
    if err != nil {
        log.Fatalf("Failed to create authorizer: %v", err)
    }
    
    // MOC server FQDN
    cloudFQDN := "your-moc-server.example.com"
    
    // Create VM client
    vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    if err != nil {
        log.Fatalf("Failed to create VM client: %v", err)
    }
    
    fmt.Println("Successfully connected to MOC service!")
}
```

### Step 4: Perform Operations

Now let's list existing virtual machines:

```go
func listVirtualMachines(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Query VMs in a resource group
    groupName := "default"
    vmName := "" // Empty string lists all VMs
    
    vms, err := vmClient.Get(ctx, groupName, vmName)
    if err != nil {
        return fmt.Errorf("failed to list VMs: %w", err)
    }
    
    if vms == nil || len(*vms) == 0 {
        fmt.Println("No virtual machines found")
        return nil
    }
    
    fmt.Printf("Found %d virtual machine(s):\n", len(*vms))
    for _, vm := range *vms {
        fmt.Printf("  - Name: %s, Status: %s\n", *vm.Name, vm.Statuses)
    }
    
    return nil
}
```

### Complete Example

Here's a complete working example:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

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
        log.Fatalf("Failed to create authorizer: %v", err)
    }

    // MOC server FQDN
    cloudFQDN := "your-moc-server.example.com"

    // Create VM client
    vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    if err != nil {
        log.Fatalf("Failed to create VM client: %v", err)
    }

    fmt.Println("Successfully connected to MOC service!")

    // List virtual machines
    if err := listVirtualMachines(vmClient); err != nil {
        log.Fatalf("Failed to list VMs: %v", err)
    }
}

func listVirtualMachines(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    groupName := "default"
    vmName := ""

    vms, err := vmClient.Get(ctx, groupName, vmName)
    if err != nil {
        return fmt.Errorf("failed to list VMs: %w", err)
    }

    if vms == nil || len(*vms) == 0 {
        fmt.Println("No virtual machines found")
        return nil
    }

    fmt.Printf("Found %d virtual machine(s):\n", len(*vms))
    for _, vm := range *vms {
        fmt.Printf("  - Name: %s\n", *vm.Name)
        if vm.Properties != nil && vm.Properties.HardwareProfile != nil {
            fmt.Printf("    CPU: %d cores, Memory: %dMB\n",
                *vm.Properties.HardwareProfile.VMSize.VCPUs,
                *vm.Properties.HardwareProfile.VMSize.MemoryMB)
        }
    }

    return nil
}
```

### Run the Application

```bash
go run main.go
```

## Common Patterns

### Using Context with Timeout

Always use context with timeout for operations:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.SomeOperation(ctx, params...)
```

### Error Handling

The SDK returns detailed errors. Always check and handle them:

```go
vm, err := vmClient.Get(ctx, group, name)
if err != nil {
    if errors.IsNotFound(err) {
        fmt.Println("VM not found")
    } else {
        log.Fatalf("Failed to get VM: %v", err)
    }
}
```

### Resource Naming

Resources are organized by groups and names:

```go
groupName := "my-resource-group"
resourceName := "my-virtual-machine"

// Get specific resource
resource, err := client.Get(ctx, groupName, resourceName)

// List all resources in group (empty name)
resources, err := client.Get(ctx, groupName, "")
```

## Debug Mode

For development and testing, you can enable debug mode to disable TLS:

```bash
export WSSD_DEBUG_MODE=on
```

Or in your configuration:

```go
import "github.com/spf13/viper"

viper.Set("Debug", true)
```

**Warning:** Never use debug mode in production!

## Next Steps

Now that you have a working application, explore more:

1. **[Authentication](authentication.md)** - Learn about different authentication methods
2. **[Client Usage](client-usage.md)** - Advanced client configuration
3. **[Service Documentation](services/compute.md)** - Explore all available services
4. **[Code Examples](examples/vm-management.md)** - More detailed examples

## Quick Reference

### Available Service Clients

```go
// Compute services
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachinescaleset"
import "github.com/microsoft/moc-sdk-for-go/services/compute/galleryimage"

// Network services
import "github.com/microsoft/moc-sdk-for-go/services/network/virtualnetwork"
import "github.com/microsoft/moc-sdk-for-go/services/network/loadbalancer"
import "github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"

// Storage services
import "github.com/microsoft/moc-sdk-for-go/services/storage/virtualharddisk"

// Security services
import "github.com/microsoft/moc-sdk-for-go/services/security/identity"
import "github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
```

### Common Operations

```go
// Create or update resource
resource, err := client.CreateOrUpdate(ctx, group, name, &resourceSpec)

// Get resource
resource, err := client.Get(ctx, group, name)

// Delete resource
err := client.Delete(ctx, group, name)

// List resources in group
resources, err := client.Get(ctx, group, "")
```

## Troubleshooting

### "connection refused" Error

Ensure the MOC server is running and accessible:

```bash
# Test connectivity
ping your-moc-server.example.com

# Check port is open
telnet your-moc-server.example.com 55000
```

### "certificate verify failed" Error

Check your certificate paths and ensure they're valid:

```bash
# Verify certificate
openssl x509 -in certs/client.pem -text -noout
```

### "permission denied" Error

Ensure your credentials have proper permissions for the operation. Check with your MOC administrator.

For more troubleshooting, see the [Troubleshooting Guide](troubleshooting.md).
