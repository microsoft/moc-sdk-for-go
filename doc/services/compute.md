# Compute Services

The Compute Services provide APIs for managing virtual machines, scale sets, images, and compute infrastructure.

## Overview

The MOC SDK Compute Services include:

- **Virtual Machines** - Create and manage VMs with full lifecycle control
- **Virtual Machine Scale Sets** - Manage groups of identical VMs
- **Gallery Images** - Manage OS and application images
- **Availability Sets** - Configure high availability for VMs
- **Bare Metal Hosts** - Manage physical compute hosts
- **Placement Groups** - Control VM placement for performance

## Virtual Machines

Virtual machines are the core compute resource.

### Creating a Virtual Machine

```go
import (
    "context"
    "os"
    "github.com/microsoft/moc-sdk-for-go/services/compute"
    "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
)

func createVM(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    
    vmSpec := &compute.VirtualMachine{
        Name:     stringPtr("web-server-01"),
        Location: stringPtr("default"),
        Properties: &compute.VirtualMachineProperties{
            HardwareProfile: &compute.HardwareProfile{
                VMSize: &compute.VMSize{
                    VCPUs:    int32Ptr(4),
                    MemoryMB: int32Ptr(8192),
                },
            },
            StorageProfile: &compute.StorageProfile{
                ImageReference: &compute.ImageReference{
                    Name: stringPtr("ubuntu-20.04"),
                },
                DataDisks: &[]compute.DataDisk{},
            },
            OsProfile: &compute.OSProfile{
                ComputerName:  stringPtr("web-server-01"),
                AdminUsername: stringPtr("azureuser"),
                AdminPassword: stringPtr(os.Getenv("VM_ADMIN_PASSWORD")), // Use environment variable
                LinuxConfiguration: &compute.LinuxConfiguration{
                    DisablePasswordAuthentication: boolPtr(false),
                },
            },
            NetworkProfile: &compute.NetworkProfile{
                NetworkInterfaces: &[]compute.NetworkInterfaceReference{
                    {
                        ID: stringPtr("/production/networkinterfaces/web-nic"),
                    },
                },
            },
        },
    }
    
    vm, err := vmClient.CreateOrUpdate(ctx, "production", "web-server-01", vmSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created VM: %s\n", *vm.ID)
    return nil
}
```

**Security Note:** The example above uses `os.Getenv("VM_ADMIN_PASSWORD")` to retrieve the password from an environment variable. Never hardcode passwords in your code. Set the environment variable before running:

```bash
export VM_ADMIN_PASSWORD="your-secure-password"
```

### VM Operations

#### Core CRUD Operations

```go
// Get a specific VM
vm, err := vmClient.Get(ctx, "production", "web-server-01")

// List all VMs in a group
vms, err := vmClient.Get(ctx, "production", "")

// Update VM
updatedVM, err := vmClient.CreateOrUpdate(ctx, "production", "web-server-01", vmSpec)

// Delete VM
err := vmClient.Delete(ctx, "production", "web-server-01")

// Query VMs with filter
vms, err := vmClient.Query(ctx, "production", "status eq 'Running'")
```

#### Lifecycle Operations

```go
// Start VM
err := vmClient.Start(ctx, "production", "web-server-01")

// Stop VM (forced)
err := vmClient.Stop(ctx, "production", "web-server-01")

// Graceful shutdown
err := vmClient.StopGraceful(ctx, "production", "web-server-01")

// Pause VM (suspend to memory)
err := vmClient.Pause(ctx, "production", "web-server-01")

// Save VM state to disk
err := vmClient.Save(ctx, "production", "web-server-01")
```

#### Disk Management

```go
// Attach a disk
err := vmClient.DiskAttach(ctx, "production", "web-server-01", "data-disk-01")

// Detach a disk
err := vmClient.DiskDetach(ctx, "production", "web-server-01", "data-disk-01")

// Remove ISO disk
err := vmClient.RemoveIsoDisk(ctx, "production", "web-server-01")
```

#### Network Management

```go
// Add network interface
err := vmClient.NetworkInterfaceAdd(ctx, "production", "web-server-01", "new-nic")

// Remove network interface
err := vmClient.NetworkInterfaceRemove(ctx, "production", "web-server-01", "old-nic")

// List VM IP addresses
ips, err := vmClient.ListIPs(ctx, "production", "web-server-01")
for _, ip := range ips {
    fmt.Printf("IP: %s\n", ip)
}
```

#### VM Resizing

```go
// Resize VM with predefined size
err := vmClient.Resize(ctx, "production", "web-server-01", 
    compute.VirtualMachineSizeTypesStandardA4, nil)

// Resize with custom size
customSize := &compute.VirtualMachineCustomSize{
    CpuCount: int32Ptr(8),
    MemoryMB: int32Ptr(16384),
}
err := vmClient.Resize(ctx, "production", "web-server-01", 
    compute.VirtualMachineSizeTypesCustom, customSize)

// Resize with GPU support
err := vmClient.ResizeEx(ctx, "production", "web-server-01",
    compute.VirtualMachineSizeTypesCustom, customSize, gpuCount)
```

#### Run Commands on VM

```go
// Execute command on VM
request := &compute.VirtualMachineRunCommandRequest{
    Command: stringPtr("ls -la /home"),
}

response, err := vmClient.RunCommand(ctx, "production", "web-server-01", request)
if err != nil {
    return err
}

fmt.Printf("Exit Code: %d\n", *response.ExitCode)
fmt.Printf("Output: %s\n", *response.Output)
```

#### Diagnostics

```go
// Repair guest agent
err := vmClient.RepairGuestAgent(ctx, "production", "web-server-01")

// Get Hyper-V VM ID
hvID, err := vmClient.GetHyperVVmId(ctx, "production", "web-server-01")
fmt.Printf("Hyper-V ID: %s\n", *hvID.ID)

// Get host node name
hostNode, err := vmClient.GetHostNodeName(ctx, "production", "web-server-01")
fmt.Printf("Host: %s\n", *hostNode.Name)

// Get host node IP
hostIP, err := vmClient.GetHostNodeIpAddress(ctx, "production", "web-server-01")
fmt.Printf("Host IP: %s\n", *hostIP.IpAddress)
```

#### Validation

```go
// Validate VM configuration
err := vmClient.Validate(ctx, "production", "web-server-01")

// Precheck VM placement
vmsToCreate := []*compute.VirtualMachine{vmSpec1, vmSpec2}
canPlace, err := vmClient.Precheck(ctx, "production", vmsToCreate)
if !canPlace {
    fmt.Println("Cannot place VMs with current resources")
}
```

## Virtual Machine Scale Sets

Manage groups of identical VMs for scalability.

### Creating a Scale Set

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachinescaleset"

func createVMSS(vmssClient *virtualmachinescaleset.VirtualMachineScaleSetClient) error {
    ctx := context.Background()
    
    vmssSpec := &compute.VirtualMachineScaleSet{
        Name:     stringPtr("web-vmss"),
        Location: stringPtr("default"),
        Sku: &compute.Sku{
            Name:     stringPtr("Standard_A2"),
            Capacity: int64Ptr(3), // 3 instances
        },
        Properties: &compute.VirtualMachineScaleSetProperties{
            VirtualMachineProfile: &compute.VirtualMachineScaleSetVMProfile{
                StorageProfile: &compute.VirtualMachineScaleSetStorageProfile{
                    ImageReference: &compute.ImageReference{
                        Name: stringPtr("ubuntu-20.04"),
                    },
                },
                OsProfile: &compute.VirtualMachineScaleSetOSProfile{
                    ComputerNamePrefix: stringPtr("web-"),
                    AdminUsername:      stringPtr("azureuser"),
                },
                NetworkProfile: &compute.VirtualMachineScaleSetNetworkProfile{
                    NetworkInterfaceConfigurations: &[]compute.VirtualMachineScaleSetNetworkConfiguration{
                        {
                            Name: stringPtr("default-nic"),
                            Properties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
                                Primary: boolPtr(true),
                            },
                        },
                    },
                },
            },
        },
    }
    
    vmss, err := vmssClient.CreateOrUpdate(ctx, "production", "web-vmss", vmssSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created scale set: %s\n", *vmss.Name)
    return nil
}
```

### Scale Set Operations

```go
// Get scale set
vmss, err := vmssClient.Get(ctx, "production", "web-vmss")

// List VMs in scale set
vms, err := vmssClient.List(ctx, "production", "web-vmss")
fmt.Printf("Scale set has %d VMs\n", len(*vms))

// Delete scale set
err := vmssClient.Delete(ctx, "production", "web-vmss")
```

## Gallery Images

Manage OS and application images for VM deployment.

### Uploading an Image

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/galleryimage"

func uploadImage(imageClient *galleryimage.GalleryImageClient) error {
    ctx := context.Background()
    
    imageSpec := &compute.GalleryImage{
        Name:     stringPtr("my-custom-image"),
        Location: stringPtr("default"),
        Properties: &compute.GalleryImageProperties{
            Identifier: &compute.GalleryImageIdentifier{
                Publisher: stringPtr("MyCompany"),
                Offer:     stringPtr("MyApp"),
                Sku:       stringPtr("1.0"),
            },
            OsType:  compute.OperatingSystemTypesLinux,
            OsState: compute.OperatingSystemStateTypesGeneralized,
        },
    }
    
    // Upload from local file
    image, err := imageClient.UploadImageFromLocal(
        ctx,
        "default",              // location
        "/path/to/image.vhdx",  // local image path
        "my-custom-image",       // image name
        imageSpec,
    )
    if err != nil {
        return err
    }
    
    fmt.Printf("Uploaded image: %s\n", *image.Name)
    return nil
}
```

### Image Upload Methods

```go
// Upload from local file
image, err := imageClient.UploadImageFromLocal(ctx, location, imagePath, name, imageSpec)

// Upload from SFS (Scale-out File Server)
sfsProps := &compute.SFSImageProperties{
    Path: stringPtr("\\\\sfs\\share\\image.vhdx"),
}
image, err := imageClient.UploadImageFromSFS(ctx, location, name, imageSpec, sfsProps)

// Upload from HTTP/Azure
azureProps := &compute.AzureGalleryImageProperties{
    ImageUrl: stringPtr("https://storage.blob.core.windows.net/images/image.vhdx"),
}
image, err := imageClient.UploadImageFromHttp(ctx, location, name, imageSpec, azureProps)
```

### Image Operations

```go
// Get image
images, err := imageClient.Get(ctx, "default", "my-custom-image")

// List all images
images, err := imageClient.Get(ctx, "default", "")

// Delete image
err := imageClient.Delete(ctx, "default", "my-custom-image")

// Validate image placement
canPlace, err := imageClient.Precheck(ctx, "default", imagePath, []*compute.GalleryImage{imageSpec})
```

## Availability Sets

Configure high availability for VMs.

### Creating an Availability Set

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/availabilityset"

func createAvailabilitySet(avsetClient *availabilityset.AvailabilitySetClient) error {
    ctx := context.Background()
    
    avsetSpec := &compute.AvailabilitySet{
        Name:     stringPtr("web-avset"),
        Location: stringPtr("default"),
        Properties: &compute.AvailabilitySetProperties{
            PlatformFaultDomainCount:  int32Ptr(2),
            PlatformUpdateDomainCount: int32Ptr(5),
        },
    }
    
    avset, err := avsetClient.Create(ctx, "production", "web-avset", avsetSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created availability set: %s\n", *avset.Name)
    return nil
}
```

### Availability Set Operations

```go
// Get availability set
avset, err := avsetClient.Get(ctx, "production", "web-avset")

// List all availability sets
avsets, err := avsetClient.Get(ctx, "production", "")

// Delete availability set
err := avsetClient.Delete(ctx, "production", "web-avset")

// Precheck placement
canPlace, err := avsetClient.Precheck(ctx, "production", []*compute.AvailabilitySet{avsetSpec})
```

### Using Availability Sets with VMs

```go
vmSpec := &compute.VirtualMachine{
    Name:     stringPtr("web-server-01"),
    Location: stringPtr("default"),
    Properties: &compute.VirtualMachineProperties{
        AvailabilitySet: &compute.SubResource{
            ID: stringPtr("/production/availabilitysets/web-avset"),
        },
        // ... other properties
    },
}
```

## Bare Metal Hosts

Manage physical compute hosts.

### Querying Bare Metal Hosts

```go
import "github.com/microsoft/moc-sdk-for-go/services/compute/baremetalhost"

func listHosts(hostClient *baremetalhost.BareMetalHostClient) error {
    ctx := context.Background()
    
    // Get all hosts
    hosts, err := hostClient.Get(ctx, "default", "")
    if err != nil {
        return err
    }
    
    for _, host := range *hosts {
        fmt.Printf("Host: %s\n", *host.Name)
        if host.Properties != nil {
            fmt.Printf("  FQDN: %s\n", *host.Properties.Fqdn)
            fmt.Printf("  Status: %s\n", host.Statuses)
        }
    }
    
    return nil
}
```

### Bare Metal Host Operations

```go
// Get specific host
host, err := hostClient.Get(ctx, "default", "host-01")

// Query hosts
hosts, err := hostClient.Query(ctx, "default", "status eq 'Ready'")

// Create or update host
hostSpec := &compute.BareMetalHost{
    Name:     stringPtr("host-01"),
    Location: stringPtr("default"),
    Properties: &compute.BareMetalHostProperties{
        Fqdn: stringPtr("host01.example.com"),
    },
}
host, err := hostClient.CreateOrUpdate(ctx, "default", "host-01", hostSpec)

// Delete host
err := hostClient.Delete(ctx, "default", "host-01")
```

## Common Patterns

### Resource References

```go
// Reference resources by ID
vmSpec.Properties.AvailabilitySet = &compute.SubResource{
    ID: stringPtr("/production/availabilitysets/web-avset"),
}

vmSpec.Properties.NetworkProfile.NetworkInterfaces = &[]compute.NetworkInterfaceReference{
    {ID: stringPtr("/production/networkinterfaces/nic-01")},
}
```

### Helper Functions

```go
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }
func int64Ptr(i int64) *int64    { return &i }
func boolPtr(b bool) *bool       { return &b }
```

## Next Steps

- [Network Services](network.md) - Configure networking
- [Storage Services](storage.md) - Manage storage
- [Code Examples](../examples/vm-management.md) - Detailed VM examples
