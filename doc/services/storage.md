# Storage Services

The Storage Services provide APIs for managing virtual hard disks and storage containers.

## Overview

The MOC SDK Storage Services include:

- **Virtual Hard Disks** - Create and manage virtual hard disks (VHDs/VHDXs)
- **Containers** - Manage storage containers

## Virtual Hard Disks

Virtual hard disks provide persistent storage for VMs.

### Creating a Virtual Hard Disk

```go
import (
    "context"
    "github.com/microsoft/moc-sdk-for-go/services/storage"
    "github.com/microsoft/moc-sdk-for-go/services/storage/virtualharddisk"
)

func createVHD(vhdClient *virtualharddisk.VirtualHardDiskClient) error {
    ctx := context.Background()
    
    vhdSpec := &storage.VirtualHardDisk{
        Name:     stringPtr("data-disk-01"),
        Location: stringPtr("default"),
        Properties: &storage.VirtualHardDiskProperties{
            DiskSizeGB: int64Ptr(100), // 100 GB disk
        },
    }
    
    vhd, err := vhdClient.CreateOrUpdate(ctx, "production", "data-disk-01", vhdSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created VHD: %s, Size: %d GB\n", *vhd.Name, *vhd.Properties.DiskSizeGB)
    return nil
}
```

### VHD Operations

```go
// Get a specific VHD
vhd, err := vhdClient.Get(ctx, "production", "data-disk-01")

// List all VHDs
vhds, err := vhdClient.Get(ctx, "production", "")

// Update VHD (resize)
vhd.Properties.DiskSizeGB = int64Ptr(200) // Resize to 200 GB
updatedVHD, err := vhdClient.CreateOrUpdate(ctx, "production", "data-disk-01", vhd)

// Delete VHD
err := vhdClient.Delete(ctx, "production", "data-disk-01")
```

### VHD Types

#### Dynamic Disk

```go
vhdSpec := &storage.VirtualHardDisk{
    Name:     stringPtr("dynamic-disk"),
    Location: stringPtr("default"),
    Properties: &storage.VirtualHardDiskProperties{
        DiskSizeGB: int64Ptr(100),
        Dynamic:    boolPtr(true), // Dynamic disk
    },
}
```

#### Fixed Disk

```go
vhdSpec := &storage.VirtualHardDisk{
    Name:     stringPtr("fixed-disk"),
    Location: stringPtr("default"),
    Properties: &storage.VirtualHardDiskProperties{
        DiskSizeGB: int64Ptr(100),
        Dynamic:    boolPtr(false), // Fixed disk
    },
}
```

### Attaching VHD to VM

```go
// In VM specification
vmSpec.Properties.StorageProfile.DataDisks = &[]compute.DataDisk{
    {
        Lun:  int32Ptr(0),
        Name: stringPtr("data-disk-01"),
        Vhd: &compute.VirtualHardDisk{
            URI: stringPtr("/production/virtualharddisks/data-disk-01"),
        },
        CreateOption: compute.DiskCreateOptionTypesAttach,
    },
}
```

### OS Disk

```go
vmSpec.Properties.StorageProfile.OsDisk = &compute.OSDisk{
    Name: stringPtr("os-disk"),
    Vhd: &compute.VirtualHardDisk{
        URI: stringPtr("/production/virtualharddisks/os-disk"),
    },
    CreateOption: compute.DiskCreateOptionTypesFromImage,
}
```

## Storage Containers

Manage storage containers for organizing resources.

### Creating a Storage Container

```go
import "github.com/microsoft/moc-sdk-for-go/services/storage/container"

func createContainer(containerClient *container.ContainerClient) error {
    ctx := context.Background()
    
    containerSpec := &storage.Container{
        Name:     stringPtr("vm-storage"),
        Location: stringPtr("default"),
        Properties: &storage.ContainerProperties{
            Path: stringPtr("/mnt/storage/vms"),
        },
    }
    
    container, err := containerClient.CreateOrUpdate(ctx, "production", "vm-storage", containerSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created container: %s\n", *container.Name)
    return nil
}
```

### Container Operations

```go
// Get container
container, err := containerClient.Get(ctx, "production", "vm-storage")

// List all containers
containers, err := containerClient.Get(ctx, "production", "")

// Delete container
err := containerClient.Delete(ctx, "production", "vm-storage")
```

## Storage Patterns

### VM with Multiple Data Disks

```go
vmSpec := &compute.VirtualMachine{
    Name:     stringPtr("database-vm"),
    Location: stringPtr("default"),
    Properties: &compute.VirtualMachineProperties{
        StorageProfile: &compute.StorageProfile{
            OsDisk: &compute.OSDisk{
                Name: stringPtr("os-disk"),
                Vhd: &compute.VirtualHardDisk{
                    URI: stringPtr("/production/virtualharddisks/db-os-disk"),
                },
                CreateOption: compute.DiskCreateOptionTypesFromImage,
            },
            DataDisks: &[]compute.DataDisk{
                {
                    Lun:  int32Ptr(0),
                    Name: stringPtr("data-disk-01"),
                    Vhd: &compute.VirtualHardDisk{
                        URI: stringPtr("/production/virtualharddisks/db-data-01"),
                    },
                    DiskSizeGB:   int64Ptr(500),
                    CreateOption: compute.DiskCreateOptionTypesEmpty,
                },
                {
                    Lun:  int32Ptr(1),
                    Name: stringPtr("data-disk-02"),
                    Vhd: &compute.VirtualHardDisk{
                        URI: stringPtr("/production/virtualharddisks/db-data-02"),
                    },
                    DiskSizeGB:   int64Ptr(500),
                    CreateOption: compute.DiskCreateOptionTypesEmpty,
                },
            },
        },
    },
}
```

### Disk Caching

```go
dataDisk := compute.DataDisk{
    Lun:  int32Ptr(0),
    Name: stringPtr("cache-disk"),
    Vhd: &compute.VirtualHardDisk{
        URI: stringPtr("/production/virtualharddisks/cache-disk"),
    },
    Caching:      compute.CachingTypesReadWrite,
    CreateOption: compute.DiskCreateOptionTypesEmpty,
}
```

### Disk Create Options

```go
// Create empty disk
CreateOption: compute.DiskCreateOptionTypesEmpty

// Create from image
CreateOption: compute.DiskCreateOptionTypesFromImage

// Attach existing disk
CreateOption: compute.DiskCreateOptionTypesAttach

// Copy from source
CreateOption: compute.DiskCreateOptionTypesCopy
```

## Best Practices

### 1. Disk Sizing

```go
// ✅ Good: Appropriate size for workload
DiskSizeGB: int64Ptr(100) // 100 GB for data

// ❌ Bad: Oversized disk
DiskSizeGB: int64Ptr(10000) // 10 TB when not needed
```

### 2. Dynamic vs Fixed Disks

```go
// Use dynamic disks for development/test
Dynamic: boolPtr(true)

// Use fixed disks for production (better performance)
Dynamic: boolPtr(false)
```

### 3. Multiple Data Disks

```go
// ✅ Good: Multiple smaller disks for RAID/striping
DataDisks: &[]compute.DataDisk{
    {Lun: int32Ptr(0), DiskSizeGB: int64Ptr(500)},
    {Lun: int32Ptr(1), DiskSizeGB: int64Ptr(500)},
    {Lun: int32Ptr(2), DiskSizeGB: int64Ptr(500)},
}

// vs. Single large disk
```

### 4. Caching Strategy

```go
// OS Disk: ReadWrite caching
OsDisk.Caching = compute.CachingTypesReadWrite

// Data Disk (read-heavy): ReadOnly caching
DataDisk.Caching = compute.CachingTypesReadOnly

// Data Disk (write-heavy): None caching
DataDisk.Caching = compute.CachingTypesNone
```

## Disk Management Operations

### Hot Add/Remove Disks

```go
// Add disk to running VM
err := vmClient.DiskAttach(ctx, "production", "database-vm", "new-data-disk")

// Remove disk from running VM
err := vmClient.DiskDetach(ctx, "production", "database-vm", "old-data-disk")
```

### Resize Disk

```go
// Get current disk
vhds, err := vhdClient.Get(ctx, "production", "data-disk-01")
if err != nil {
    return err
}

vhd := &(*vhds)[0]

// Increase size
currentSize := *vhd.Properties.DiskSizeGB
newSize := currentSize + 100 // Add 100 GB
vhd.Properties.DiskSizeGB = &newSize

// Update disk
updatedVHD, err := vhdClient.CreateOrUpdate(ctx, "production", "data-disk-01", vhd)
```

**Note:** You can only increase disk size, not decrease it.

## Storage Monitoring

### Check Disk Usage

```go
vhds, err := vhdClient.Get(ctx, "production", "")
if err != nil {
    return err
}

for _, vhd := range *vhds {
    if vhd.Properties != nil {
        fmt.Printf("Disk: %s\n", *vhd.Name)
        fmt.Printf("  Size: %d GB\n", *vhd.Properties.DiskSizeGB)
        if vhd.Properties.DiskFileFormat != nil {
            fmt.Printf("  Format: %s\n", *vhd.Properties.DiskFileFormat)
        }
        if vhd.Statuses != nil {
            fmt.Printf("  Status: %s\n", *vhd.Statuses)
        }
    }
}
```

## Error Handling

### Common Storage Errors

```go
import "github.com/microsoft/moc/pkg/errors"

vhd, err := vhdClient.CreateOrUpdate(ctx, group, name, vhdSpec)
if err != nil {
    if errors.IsAlreadyExists(err) {
        fmt.Println("Disk already exists")
    } else if errors.IsOutOfSpace(err) {
        fmt.Println("Insufficient storage space")
    } else if errors.IsInvalidInput(err) {
        fmt.Println("Invalid disk configuration")
    } else {
        log.Fatalf("Failed to create disk: %v", err)
    }
}
```

## Next Steps

- [Compute Services](compute.md) - Use disks with VMs
- [Storage Examples](../examples/storage-operations.md) - Detailed examples
- [Architecture](../architecture.md) - Storage architecture
