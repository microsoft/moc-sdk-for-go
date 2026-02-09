# API Reference

Complete API reference for the MOC SDK for Go.

## Client Initialization

All service clients follow the same initialization pattern:

```go
func New<Service>Client(cloudFQDN string, authorizer auth.Authorizer) (*<Service>Client, error)
```

**Parameters:**
- `cloudFQDN` - MOC server fully qualified domain name
- `authorizer` - Authentication authorizer

**Returns:**
- Service client instance or error

## Common Method Signatures

### Get

Retrieve one or all resources:

```go
Get(ctx context.Context, group string, name string) (*[]Resource, error)
```

- `name` empty = list all resources in group
- `name` specified = get specific resource

### CreateOrUpdate

Create or update a resource:

```go
CreateOrUpdate(ctx context.Context, group string, name string, resource *Resource) (*Resource, error)
```

### Delete

Delete a resource:

```go
Delete(ctx context.Context, group string, name string) error
```

### Query

Query resources with filter:

```go
Query(ctx context.Context, group string, query string) (*[]Resource, error)
```

## Compute API

### Virtual Machine Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"

type VirtualMachineClient struct { ... }

// CRUD Operations
Get(ctx, group, name string) (*[]compute.VirtualMachine, error)
CreateOrUpdate(ctx, group, name string, vm *compute.VirtualMachine) (*compute.VirtualMachine, error)
Delete(ctx, group, name string) error
Query(ctx, group, query string) (*[]compute.VirtualMachine, error)

// Lifecycle
Start(ctx, group, name string) error
Stop(ctx, group, name string) error
StopGraceful(ctx, group, name string) error
Pause(ctx, group, name string) error
Save(ctx, group, name string) error

// Disk Management
DiskAttach(ctx, group, vmName, diskName string) error
DiskDetach(ctx, group, vmName, diskName string) error
RemoveIsoDisk(ctx, group, name string) error

// Network Management
NetworkInterfaceAdd(ctx, group, vmName, nicName string) error
NetworkInterfaceRemove(ctx, group, vmName, nicName string) error
ListIPs(ctx, group, name string) ([]string, error)

// VM Management
Resize(ctx, group, vmName string, newSize VirtualMachineSizeTypes, customSize *VirtualMachineCustomSize) error
RunCommand(ctx, group, vmName string, request *VirtualMachineRunCommandRequest) (*VirtualMachineRunCommandResponse, error)
RepairGuestAgent(ctx, group, vmName string) error

// Diagnostics
GetHyperVVmId(ctx, group, name string) (*VirtualMachineHyperVVmId, error)
GetHostNodeName(ctx, group, name string) (*VirtualMachineHostNodeName, error)
GetHostNodeIpAddress(ctx, group, name string) (*VirtualMachineHostNodeIpAddress, error)

// Validation
Validate(ctx, group, name string) error
Precheck(ctx, group string, vms []*compute.VirtualMachine) (bool, error)
```

## Network API

### Virtual Network Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/virtualnetwork"

type VirtualNetworkClient struct { ... }

Get(ctx, group, name string) (*[]network.VirtualNetwork, error)
CreateOrUpdate(ctx, group, name string, vnet *network.VirtualNetwork) (*network.VirtualNetwork, error)
Delete(ctx, group, name string) error
```

### Network Interface Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"

type NetworkInterfaceClient struct { ... }

Get(ctx, group, name string) (*[]network.NetworkInterface, error)
CreateOrUpdate(ctx, group, name string, nic *network.NetworkInterface) (*network.NetworkInterface, error)
Delete(ctx, group, name string) error
```

### Load Balancer Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/loadbalancer"

type LoadBalancerClient struct { ... }

Get(ctx, group, name string) (*[]network.LoadBalancer, error)
CreateOrUpdate(ctx, group, name string, lb *network.LoadBalancer) (*network.LoadBalancer, error)
Delete(ctx, group, name string) error
```

## Storage API

### Virtual Hard Disk Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/storage/virtualharddisk"

type VirtualHardDiskClient struct { ... }

Get(ctx, group, name string) (*[]storage.VirtualHardDisk, error)
CreateOrUpdate(ctx, group, name string, vhd *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error)
Delete(ctx, group, name string) error
```

## Security API

### Identity Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/identity"

type IdentityClient struct { ... }

Get(ctx, group, name string) (*[]security.Identity, error)
CreateOrUpdate(ctx, group, name string, identity *security.Identity) (*security.Identity, error)
Delete(ctx, group, name string) error
```

### Key Vault Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

type KeyVaultClient struct { ... }

Get(ctx, group, name string) (*[]security.KeyVault, error)
CreateOrUpdate(ctx, group, name string, kv *security.KeyVault) (*security.KeyVault, error)
Delete(ctx, group, name string) error
```

## Cloud API

### Location Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/location"

type LocationClient struct { ... }

Get(ctx, location, name string) (*[]cloud.Location, error)
CreateOrUpdate(ctx, location, name string, loc *cloud.Location) (*cloud.Location, error)
Delete(ctx, location, name string) error
```

### Node Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/node"

type NodeClient struct { ... }

Get(ctx, location, name string) (*[]cloud.Node, error)
CreateOrUpdate(ctx, location, name string, node *cloud.Node) (*cloud.Node, error)
Delete(ctx, location, name string) error
Query(ctx, location, query string) (*[]cloud.Node, error)
```

## Admin API

### Health Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/health"

type HealthClient struct { ... }

Get(ctx, location, name string) (*[]admin.Health, error)
```

### Version Client

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/version"

type VersionClient struct { ... }

Get(ctx context.Context) (*admin.Version, error)
```

## Data Types

### Resource Common Fields

All resources share these common fields:

```go
type Resource struct {
    ID       *string
    Name     *string
    Type     *string
    Location *string
    Tags     map[string]*string
    Version  *string
    Statuses *string
}
```

### Context Usage

All operations require `context.Context`:

```go
import "context"

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## Error Handling

```go
import "github.com/microsoft/moc/pkg/errors"

// Check error types
if errors.IsNotFound(err) { }
if errors.IsAlreadyExists(err) { }
if errors.IsInvalidInput(err) { }
```

## Constants

### Server Ports

```go
const (
    ServerPort int = 55000  // Main service port
    AuthPort   int = 65000  // Authentication port
)
```

### Debug Mode

```go
const debugModeTLS = "WSSD_DEBUG_MODE"
```

## Next Steps

- [Getting Started](getting-started.md) - Quick start guide
- [Client Usage](client-usage.md) - Client configuration
- [Service Documentation](services/compute.md) - Detailed service docs
