# Usage

> Last updated: 2026-01-11

## Prerequisites

- Go 1.24+ (per go.mod)
- For Windows DLL build: `mingw-w64` (install via `sudo apt-get install mingw-w64` on WSL)
- Access to Microsoft private Go modules: `GOPRIVATE=github.com/microsoft`

## Installation

```bash
go get github.com/microsoft/moc-sdk-for-go
```

## Building

```bash
# Full build (tidy, format, build, unittest)
make

# Just build
make build

# Format code
make format

# Run all tests
make test

# Run unit tests only
make unittest

# Run linter
make golangci-lint
```

## SDK Usage Examples

### Authentication

```go
import (
    "github.com/microsoft/moc/pkg/auth"
    "github.com/microsoft/moc-sdk-for-go/services/security/authentication"
)

// Create authorizer from environment
authorizer, err := auth.NewAuthorizerFromEnvironment(serverAddress)

// Or login with config
authClient, err := authentication.NewAuthenticationClient(cloudFQDN, authorizer)
result, err := authClient.LoginWithConfig(ctx, group, loginConfig, enableRenew)
```

### Virtual Machine Operations

```go
import (
    "github.com/microsoft/moc-sdk-for-go/services/compute"
    "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
)

// Create client
vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)

// Get VMs
vms, err := vmClient.Get(ctx, resourceGroup, vmName)

// Create/Update VM
vm := &compute.VirtualMachine{
    Name:     &vmName,
    Location: &location,
    // ... configure properties
}
result, err := vmClient.CreateOrUpdate(ctx, resourceGroup, vmName, vm)

// VM lifecycle
err = vmClient.Start(ctx, resourceGroup, vmName)
err = vmClient.Stop(ctx, resourceGroup, vmName)
err = vmClient.StopGraceful(ctx, resourceGroup, vmName)
err = vmClient.Restart(ctx, resourceGroup, vmName)

// Disk operations
err = vmClient.DiskAttach(ctx, resourceGroup, vmName, diskName)
err = vmClient.DiskDetach(ctx, resourceGroup, vmName, diskName)

// Network interface operations
err = vmClient.NetworkInterfaceAdd(ctx, resourceGroup, vmName, nicName)
err = vmClient.NetworkInterfaceRemove(ctx, resourceGroup, vmName, nicName)

// Query VMs
vms, err := vmClient.GetByComputerName(ctx, resourceGroup, computerName)
ips, err := vmClient.ListIPs(ctx, resourceGroup, vmName)
```

### Network Operations

```go
import (
    "github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"
    "github.com/microsoft/moc-sdk-for-go/services/network/virtualnetwork"
)

// Network interface client
nicClient, err := networkinterface.NewInterfaceClient(cloudFQDN, authorizer)
nics, err := nicClient.Get(ctx, group, nicName)

// Virtual network client
vnetClient, err := virtualnetwork.NewVirtualNetworkClient(cloudFQDN, authorizer)
```

### Storage Operations

```go
import "github.com/microsoft/moc-sdk-for-go/services/storage/virtualharddisk"

vhdClient, err := virtualharddisk.NewVirtualHardDiskClient(cloudFQDN, authorizer)
```

## Configuration

### Environment Variables

- `WSSD_DEBUG_MODE=on` - Enable insecure (non-TLS) connections for testing

### Viper Configuration

Debug mode can also be set via Viper config:
```go
viper.SetDefault("Debug", true)
```

## Common Commands

```bash
# Build everything
make

# Run tests
make test

# Run unit tests only (faster)
make unittest

# Fix module issues
go mod tidy

# Format all Go code
make format
```
