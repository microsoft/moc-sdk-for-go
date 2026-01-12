# Architecture

> Last updated: 2026-01-11

## Overview

The MOC SDK for Go (`moc-sdk-for-go`) is a Go SDK providing Azure-compatible API abstractions for Microsoft's MOC (Microsoft On-premises Cloud) infrastructure. It enables programmatic management of compute, storage, network, and security resources through gRPC-based communication with MOC cloud agents.

## System Type

**SDK/Client Library** - Provides client interfaces for MOC cloud agent services, with an optional C++ wrapper DLL for Windows interoperability.

## Core Components

### 1. Service Clients (`services/`)

Organized by resource domain, each providing Azure-style client interfaces:

- **compute/** - VirtualMachine, VirtualMachineScaleSet, AvailabilitySet, GalleryImage, PlacementGroup, BareMetalHost/Machine
- **network/** - VirtualNetwork, LogicalNetwork, NetworkInterface, NetworkSecurityGroup, LoadBalancer, PublicIPAddress, VIPPool, MACPool
- **storage/** - VirtualHardDisk, Container
- **security/** - Authentication, Certificate, Identity, KeyVault, Role, RoleAssignment
- **cloud/** - Location, Group, Node, Zone
- **admin/** - Health, Logging, Debug, Recovery, Validation, Version

### 2. Client Connection Layer (`pkg/client/`)

Central gRPC connection management:
- Connection caching with state validation
- TLS/mTLS authentication via `auth.Authorizer`
- Debug mode support for insecure connections
- Keepalive parameters and error interceptors

### 3. C++ Wrapper (`wrapper/cpp/`)

Windows DLL (`MocCppWrapper.dll`) exposing SDK functions to C++ consumers:
- Security login/logout operations
- KeyVault key operations (wrap/unwrap/encrypt/decrypt)
- Telemetry instrumentation

## Data Flow

```
Client Code
    ↓
[Service Client] (e.g., VirtualMachineClient)
    ↓
[WSSD Client] (internal gRPC client)
    ↓
[gRPC Connection] (cached, authenticated)
    ↓
[MOC Cloud Agent] (wssdagent)
```

## External Dependencies

### Core Dependencies
- `github.com/microsoft/moc` - MOC protobuf definitions and core utilities
- `google.golang.org/grpc` - gRPC communication
- `github.com/Azure/go-autorest` - Azure-compatible REST patterns

### Testing
- `github.com/stretchr/testify` - Assertions and mocks
- `sigs.k8s.io/controller-runtime` - Kubernetes controller utilities

## Key Design Decisions

1. **Azure API Compatibility**: Type definitions and client patterns mirror Azure SDK for Go (`github.com/Azure/azure-sdk-for-go`) for familiar developer experience

2. **Two-Layer Client Pattern**: Public `*Client` wraps internal WSSD `client` to separate Azure-style interface from gRPC implementation

3. **Connection Caching**: Single shared connection pool with automatic reconnection on failure states

4. **Protobuf Translation**: SDK types ↔ protobuf types conversion in `wssd.go` files using helper functions in `virtualmachine.go` (or equivalent)

5. **Cross-Platform Support**: Primary Go SDK with Windows C-shared library for C++ integration
