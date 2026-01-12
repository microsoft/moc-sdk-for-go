# Patterns

> Last updated: 2026-01-11

## Design Patterns

### 1. Two-Layer Client Pattern

Each service has a public client wrapping an internal WSSD client:

```go
// Public client (client.go)
type VirtualMachineClient struct {
    compute.BaseClient
    internal   Service      // Interface for testability
    cloudFQDN  string
    authorizer auth.Authorizer
}

func NewVirtualMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineClient, error) {
    c, err := newVirtualMachineClient(cloudFQDN, authorizer)  // internal wssd client
    return &VirtualMachineClient{internal: c, ...}, nil
}

// Public method delegates to internal
func (c *VirtualMachineClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
    return c.internal.Get(ctx, group, name)
}
```

### 2. Service Interface Pattern

Services define interfaces for testability and mocking:

```go
// services/compute/virtualmachine/client.go
type Service interface {
    Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
    CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
    Delete(context.Context, string, string) error
    // ... other operations
}
```

### 3. WSSD Client Pattern (gRPC Layer)

Internal clients in `wssd.go` handle protobuf translation:

```go
// wssd.go
type client struct {
    wssdcloudcompute.VirtualMachineAgentClient
}

func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
    request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_GET, group, name, nil)
    response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
    return c.getVirtualMachineFromResponse(response, group), nil
}
```

### 4. Retry with Version Check Pattern

Operations retry on version conflicts:

```go
func (c *VirtualMachineClient) Update(ctx context.Context, group, vmName string, updateFunctor UpdateFunctor) error {
    for {
        vms, err := c.Get(ctx, group, vmName)
        vm, err := updateFunctor.Update(ctx, &(*vms)[0])
        _, err = c.CreateOrUpdate(ctx, group, vmName, vm)
        if err != nil {
            if errors.IsInvalidVersion(err) {
                time.Sleep(100 * time.Millisecond)
                continue  // Retry on stale version
            }
            return err
        }
        break
    }
    return nil
}
```

### 5. Functor Pattern for Updates

Extensible update operations via interfaces:

```go
type UpdateFunctor interface {
    Update(context.Context, *compute.VirtualMachine) (*compute.VirtualMachine, error)
}
```

## Code Organization

### Service Structure

Each service follows this file layout:
```
services/<domain>/<resource>/
├── client.go           # Public client, Service interface
├── wssd.go             # gRPC client, protobuf translation
├── <resource>.go       # Type definitions (SDK ↔ protobuf helpers)
└── <resource>_test.go  # Unit tests
```

### Domain Type Files

Central type definitions per domain:
- `services/compute/compute.go` - VirtualMachine, StorageProfile, etc.
- `services/network/network.go` - VirtualNetwork, NetworkInterface, etc.
- `services/storage/storage.go` - VirtualHardDisk, Container, etc.
- `services/security/security.go` - Certificate, Identity, etc.

## Naming Conventions

- **Packages**: lowercase, single word (`virtualmachine`, `networkinterface`)
- **Clients**: `<Resource>Client` (e.g., `VirtualMachineClient`)
- **Constructors**: `New<Resource>Client(cloudFQDN, authorizer)`
- **Internal**: `new<Resource>Client()` (lowercase) for WSSD client
- **Methods**: Azure SDK style (`Get`, `CreateOrUpdate`, `Delete`)

## Error Handling

Uses wrapped errors from `github.com/microsoft/moc/pkg/errors`:

```go
if vms == nil || len(*vms) == 0 {
    return errors.Wrapf(errors.NotFound, "Virtual Machine [%s] not found", vmName)
}

if errors.IsInvalidVersion(err) {
    // Handle version conflict
}
```

## Async Patterns

No goroutines in core SDK - callers handle concurrency. Authentication renewal uses background goroutines with context cancellation.
