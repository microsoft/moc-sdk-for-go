# SDK Architecture

This document describes the architecture and design principles of the MOC SDK for Go.

## Overview

The MOC SDK for Go is designed as a layered architecture that provides type-safe, idiomatic Go interfaces for interacting with Microsoft's cloud infrastructure management services.

```
┌─────────────────────────────────────────────────────────┐
│                   Your Application                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              Client Layer (pkg/client/)                  │
│  High-level facades: ComputeClient, NetworkClient, etc. │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│           Service Layer (services/*/client.go)           │
│  VirtualMachine, VirtualNetwork, LoadBalancer, etc.     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│         WSSD Adapters (services/*/wssd.go)              │
│  Convert SDK types ↔ Protocol buffer types              │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│          gRPC Client Stubs (from moc package)           │
│  Generated from protobuf definitions                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              MOC Backend Service (WSSD)                  │
│  Cloud infrastructure management service                 │
└─────────────────────────────────────────────────────────┘
```

## Layer Breakdown

### 1. Client Layer (`pkg/client/`)

The client layer provides high-level facades that aggregate multiple services. This layer is optional but convenient for common scenarios.

**Purpose:**
- Simplified API for common multi-service operations
- Connection management and caching
- Default configuration

**Example:**
```go
type ComputeClient struct {
    VirtualMachines      *virtualmachine.VirtualMachineClient
    GalleryImages        *galleryimage.GalleryImageClient
    AvailabilitySets     *availabilityset.AvailabilitySetClient
    // ... more compute services
}
```

**Files:**
- `client.go` - Connection management, gRPC dial options
- `compute.go` - Compute service facade
- `network.go` - Network service facade
- `storage.go` - Storage service facade
- `security.go` - Security service facade
- `admin.go` - Admin service facade
- `cloud.go` - Cloud service facade

### 2. Service Layer (`services/*/client.go`)

The service layer provides individual service clients. Each client corresponds to a specific resource type.

**Purpose:**
- Type-safe Go interfaces for resources
- Context-aware operations
- Resource lifecycle management
- Error handling and conversion

**Structure:**
```go
type Service interface {
    Get(context.Context, string, string) (*[]Resource, error)
    CreateOrUpdate(context.Context, string, string, *Resource) (*Resource, error)
    Delete(context.Context, string, string) error
    // Resource-specific operations
}

type ResourceClient struct {
    BaseClient
    internal   Service
    cloudFQDN  string
    authorizer auth.Authorizer
}
```

**Example Services:**
- `virtualmachine` - VM management
- `virtualnetwork` - Network management
- `loadbalancer` - Load balancer management
- `identity` - Identity management

### 3. WSSD Adapter Layer (`services/*/wssd.go`)

The WSSD (Windows Software-Defined) adapter layer converts between SDK types and protocol buffer types.

**Purpose:**
- Type conversion (SDK ↔ Protobuf)
- gRPC client initialization
- Connection establishment
- Request/response marshaling

**Naming Convention:**
```go
// Convert SDK type to protobuf
func getWssdVirtualMachine(c *compute.VirtualMachine) *wssdcloudcompute.VirtualMachine

// Convert protobuf to SDK type
func getVirtualMachine(c *wssdcloudcompute.VirtualMachine) *compute.VirtualMachine
```

### 4. gRPC Layer

The gRPC layer is provided by the `github.com/microsoft/moc` package and generated from protobuf definitions.

**Features:**
- Strongly-typed RPC stubs
- Automatic serialization
- Connection multiplexing
- Keep-alive support

## Core Components

### Connection Management

Connection caching for efficiency:

```go
var (
    mux             sync.Mutex
    connectionCache map[string]*grpc.ClientConn
)

func getClientConnection(serverAddress *string, authorizer auth.Authorizer) (*grpc.ClientConn, error) {
    endpoint := getServerEndpoint(serverAddress)
    
    mux.Lock()
    defer mux.Unlock()
    
    // Check cache
    if conn, ok := connectionCache[endpoint]; ok {
        if isValidConnection(conn) {
            return conn, nil
        }
    }
    
    // Create new connection
    conn, err := grpc.Dial(endpoint, getDefaultDialOption(authorizer)...)
    if err != nil {
        return nil, err
    }
    
    connectionCache[endpoint] = conn
    return conn, nil
}
```

### Resource Naming Convention

All resources follow a consistent naming pattern:

```
/<group>/<resource-type>/<name>
```

**Example:**
```go
// Resource group: "production"
// Resource type: "virtualmachines"
// Resource name: "web-server-01"
// Full path: "/production/virtualmachines/web-server-01"

vm, err := client.Get(ctx, "production", "web-server-01")
```

### Context Propagation

All operations accept a `context.Context` for:
- Timeout control
- Cancellation
- Request tracing

```go
func (c *VirtualMachineClient) Get(
    ctx context.Context,
    group string,
    name string,
) (*[]compute.VirtualMachine, error)
```

### Error Handling

Errors are wrapped with context:

```go
import "github.com/microsoft/moc/pkg/errors"

// Check error types
if errors.IsNotFound(err) {
    // Resource doesn't exist
}

if errors.IsAlreadyExists(err) {
    // Resource already exists
}

if errors.IsInvalidInput(err) {
    // Invalid parameters
}
```

## Data Models

### Resource Structure

All resources follow a common pattern:

```go
type VirtualMachine struct {
    // ARM-style resource properties
    ID       *string
    Name     *string
    Type     *string
    Location *string
    Tags     map[string]*string
    Version  *string
    
    // Resource-specific properties
    Properties *VirtualMachineProperties
    
    // System metadata
    Statuses *string
}
```

### Properties Pattern

Complex resource settings are in a `Properties` struct:

```go
type VirtualMachineProperties struct {
    HardwareProfile *HardwareProfile
    StorageProfile  *StorageProfile
    OsProfile       *OSProfile
    NetworkProfile  *NetworkProfile
    SecurityProfile *SecurityProfile
    // ... more profiles
}
```

### Pointer vs Value Types

The SDK uses pointers for optional fields:

```go
type VMSize struct {
    VCPUs    *int32  // Optional - pointer
    MemoryMB *int32  // Optional - pointer
}

// Check if field is set
if vm.Properties.HardwareProfile.VMSize.VCPUs != nil {
    cpuCount := *vm.Properties.HardwareProfile.VMSize.VCPUs
}
```

## Design Patterns

### 1. Client Construction

Consistent client initialization:

```go
func NewResourceClient(cloudFQDN string, authorizer auth.Authorizer) (*ResourceClient, error) {
    c, err := newInternalClient(cloudFQDN, authorizer)
    if err != nil {
        return nil, err
    }
    
    return &ResourceClient{
        internal:   c,
        cloudFQDN:  cloudFQDN,
        authorizer: authorizer,
    }, nil
}
```

### 2. Service Interface

Service interface pattern for testability:

```go
type Service interface {
    Get(context.Context, string, string) (*[]Resource, error)
    CreateOrUpdate(context.Context, string, string, *Resource) (*Resource, error)
    Delete(context.Context, string, string) error
}

type client struct {
    cloudFQDN  string
    authorizer auth.Authorizer
}

func (c *client) Get(ctx context.Context, group, name string) (*[]Resource, error) {
    // Implementation
}
```

### 3. Factory Functions

Helper functions for creating resources:

```go
func NewVirtualMachine(name, location string) *compute.VirtualMachine {
    return &compute.VirtualMachine{
        Name:     &name,
        Location: &location,
        Properties: &compute.VirtualMachineProperties{
            HardwareProfile: &compute.HardwareProfile{},
            StorageProfile:  &compute.StorageProfile{},
            NetworkProfile:  &compute.NetworkProfile{},
        },
    }
}
```

## Directory Structure

```
moc-sdk-for-go/
├── pkg/
│   ├── client/              # High-level client facades
│   │   ├── client.go        # Connection management
│   │   ├── compute.go       # Compute client facade
│   │   ├── network.go       # Network client facade
│   │   └── ...
│   └── constant/            # SDK constants
│
├── services/
│   ├── compute/             # Compute services
│   │   ├── compute.go       # Common types
│   │   ├── virtualmachine/  # VM service
│   │   │   ├── client.go    # VM client
│   │   │   ├── wssd.go      # Protobuf conversion
│   │   │   └── virtualmachine.go  # VM operations
│   │   └── .../             # Other compute services
│   │
│   ├── network/             # Network services
│   ├── storage/             # Storage services
│   ├── security/            # Security services
│   ├── cloud/               # Cloud services
│   └── admin/               # Admin services
│
└── wrapper/                 # C++ wrapper
    └── cpp/                 # C-shared library
```

## Extensibility

### Adding a New Resource Type

1. **Define types** in `services/<category>/<resource>.go`
2. **Implement Service interface** in `services/<category>/<resource>/client.go`
3. **Add WSSD converters** in `services/<category>/<resource>/wssd.go`
4. **Add to facade** (optional) in `pkg/client/<category>.go`

### Custom Interceptors

Add custom gRPC interceptors:

```go
import "google.golang.org/grpc"

opts := []grpc.DialOption{
    grpc.WithUnaryInterceptor(myInterceptor),
}
```

## Performance Considerations

### Connection Pooling

- Connections are cached per endpoint
- Reused across multiple operations
- Validated before reuse

### Keep-Alive

- Prevents connection drops
- Reduces reconnection overhead
- Configured automatically

### Lazy Initialization

- Clients created on-demand
- Connections established when needed
- Resources cleaned up when unused

## Security Architecture

### Authentication Flow

```
1. Application creates Authorizer
2. Client receives Authorizer
3. Connection established with TLS
4. Each request includes auth metadata
5. Server validates credentials
6. Response returned
```

### TLS Configuration

- Mutual TLS (mTLS) by default
- Client and server authentication
- Certificate validation
- Debug mode for testing only

## C++ Wrapper

The SDK includes a C++ wrapper for cross-language support:

- **Built as**: Windows DLL (`MocCppWrapper.dll`)
- **Exports**: C-compatible functions
- **Uses**: cgo for Go-C bridging

See [C++ Wrapper Documentation](advanced/cpp-wrapper.md) for details.

## Next Steps

- [Client Usage](client-usage.md) - Using the client layer
- [gRPC Communication](advanced/grpc-communication.md) - Deep dive into gRPC
- [Error Handling](advanced/error-handling.md) - Error patterns
