# Structure

> Last updated: 2026-01-11

## Directory Layout

```
moc-sdk-for-go/
├── pkg/                          # Shared packages
│   ├── client/                   # gRPC connection management
│   │   ├── client.go             # Connection caching, dial options
│   │   ├── admin.go              # Admin service clients
│   │   ├── cloud.go              # Cloud service clients
│   │   ├── compute.go            # Compute service clients
│   │   ├── network.go            # Network service clients
│   │   ├── security.go           # Security service clients
│   │   └── storage.go            # Storage service clients
│   └── constant/                 # SDK constants
│
├── services/                     # Service implementations
│   ├── admin/                    # Administrative services
│   │   ├── debug/                # Debug operations
│   │   ├── health/               # Health checks
│   │   ├── logging/              # Logging configuration
│   │   ├── recovery/             # Recovery operations
│   │   ├── validation/           # Validation checks
│   │   └── version/              # Version info
│   │
│   ├── cloud/                    # Cloud management
│   │   ├── cloud.go              # Type definitions
│   │   ├── group/                # Resource groups
│   │   ├── location/             # Locations
│   │   ├── node/                 # Cluster nodes
│   │   └── zone/                 # Availability zones
│   │
│   ├── compute/                  # Compute resources
│   │   ├── compute.go            # Type definitions (68K+ lines)
│   │   ├── vmSizes.go            # VM size constants
│   │   ├── availabilityset/      # Availability sets
│   │   ├── baremetalhost/        # Bare metal hosts
│   │   ├── baremetalmachine/     # Bare metal machines
│   │   ├── galleryimage/         # Gallery images
│   │   ├── placementgroup/       # Placement groups
│   │   ├── virtualmachine/       # Virtual machines
│   │   ├── virtualmachineimage/  # VM images
│   │   └── virtualmachinescaleset/ # VM scale sets
│   │
│   ├── network/                  # Network resources
│   │   ├── network.go            # Type definitions
│   │   ├── common.go             # Shared utilities
│   │   ├── loadbalancer/         # Load balancers
│   │   ├── logicalnetwork/       # Logical networks
│   │   ├── macpool/              # MAC address pools
│   │   ├── networkinterface/     # NICs
│   │   ├── networksecuritygroup/ # NSGs
│   │   ├── publicipaddress/      # Public IPs
│   │   ├── vippool/              # VIP pools
│   │   └── virtualnetwork/       # Virtual networks
│   │
│   ├── security/                 # Security resources
│   │   ├── security.go           # Type definitions
│   │   ├── providerTypes.go      # Provider types
│   │   ├── authentication/       # Auth operations
│   │   ├── certificate/          # Certificates
│   │   ├── identity/             # Identities
│   │   ├── keyvault/             # Key vault + key/
│   │   ├── role/                 # RBAC roles
│   │   └── roleassignment/       # Role assignments
│   │
│   └── storage/                  # Storage resources
│       ├── storage.go            # Type definitions
│       ├── container/            # Storage containers
│       └── virtualharddisk/      # VHDs
│
├── wrapper/                      # C++ interop wrapper
│   ├── telemetry.go              # Telemetry helpers
│   └── cpp/
│       └── main.go               # cgo exports (Windows DLL)
│
├── doc/                          # Documentation
│   └── README.md                 # (empty)
│
├── .github/                      # GitHub workflows
├── .pipelines/                   # Azure DevOps pipelines
├── .gdn/                         # Guardian config
│
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── Makefile                      # Build commands
├── README.md                     # Project readme
├── LICENSE                       # Apache 2.0
└── SECURITY.md                   # Security policy
```

## Key Files

| File | Purpose |
|------|---------|
| `pkg/client/client.go` | gRPC connection caching and management |
| `services/compute/compute.go` | All compute type definitions |
| `services/compute/virtualmachine/client.go` | VM client interface |
| `services/compute/virtualmachine/wssd.go` | VM gRPC implementation |
| `wrapper/cpp/main.go` | C++ DLL exports |
| `Makefile` | Build, test, format commands |

## Module Relationships

```
pkg/client         → provides gRPC clients to all services
    ↓
services/<domain>  → uses pkg/client for connections
    ↓
github.com/microsoft/moc → protobuf definitions (external)
```
