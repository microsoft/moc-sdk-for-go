# Virtual Machine Management Examples

Complete examples for managing virtual machines with the MOC SDK.

## Basic VM Creation

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

    // Create VM client
    cloudFQDN := "moc-server.example.com"
    vmClient, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create a VM
    if err := createVM(vmClient); err != nil {
        log.Fatalf("Failed to create VM: %v", err)
    }

    // List VMs
    if err := listVMs(vmClient); err != nil {
        log.Fatalf("Failed to list VMs: %v", err)
    }

    // Manage VM lifecycle
    if err := manageVMLifecycle(vmClient); err != nil {
        log.Fatalf("Failed to manage VM: %v", err)
    }
}

func createVM(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    vmSpec := &compute.VirtualMachine{
        Name:     stringPtr("web-server-01"),
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
                OsDisk: &compute.OSDisk{
                    Name: stringPtr("os-disk"),
                    Vhd: &compute.VirtualHardDisk{
                        URI: stringPtr("/production/virtualharddisks/web-os-disk"),
                    },
                    CreateOption: compute.DiskCreateOptionTypesFromImage,
                },
                DataDisks: &[]compute.DataDisk{
                    {
                        Lun:  int32Ptr(0),
                        Name: stringPtr("data-disk-01"),
                        Vhd: &compute.VirtualHardDisk{
                            URI: stringPtr("/production/virtualharddisks/web-data-01"),
                        },
                        DiskSizeGB:   int64Ptr(100),
                        CreateOption: compute.DiskCreateOptionTypesEmpty,
                    },
                },
            },
            OsProfile: &compute.OSProfile{
                ComputerName:  stringPtr("web-server-01"),
                AdminUsername: stringPtr("azureuser"),
                AdminPassword: stringPtr("SecurePassword123!"),
                LinuxConfiguration: &compute.LinuxConfiguration{
                    DisablePasswordAuthentication: boolPtr(false),
                },
            },
            NetworkProfile: &compute.NetworkProfile{
                NetworkInterfaces: &[]compute.NetworkInterfaceReference{
                    {
                        ID: stringPtr("/production/networkinterfaces/web-nic-01"),
                    },
                },
            },
        },
        Tags: map[string]*string{
            "environment": stringPtr("production"),
            "role":        stringPtr("webserver"),
        },
    }

    vm, err := vmClient.CreateOrUpdate(ctx, "production", "web-server-01", vmSpec)
    if err != nil {
        return fmt.Errorf("failed to create VM: %w", err)
    }

    fmt.Printf("Created VM: %s\n", *vm.Name)
    fmt.Printf("  ID: %s\n", *vm.ID)
    fmt.Printf("  Location: %s\n", *vm.Location)

    return nil
}

func listVMs(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    vms, err := vmClient.Get(ctx, "production", "")
    if err != nil {
        return fmt.Errorf("failed to list VMs: %w", err)
    }

    fmt.Printf("\nFound %d VMs:\n", len(*vms))
    for _, vm := range *vms {
        fmt.Printf("  - %s\n", *vm.Name)
        if vm.Properties != nil && vm.Properties.HardwareProfile != nil {
            profile := vm.Properties.HardwareProfile
            if profile.VMSize != nil {
                fmt.Printf("    CPU: %d cores, Memory: %d MB\n",
                    *profile.VMSize.VCPUs,
                    *profile.VMSize.MemoryMB)
            }
        }
    }

    return nil
}

func manageVMLifecycle(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    vmName := "web-server-01"
    group := "production"

    // Start VM
    fmt.Println("\nStarting VM...")
    if err := vmClient.Start(ctx, group, vmName); err != nil {
        return fmt.Errorf("failed to start VM: %w", err)
    }
    fmt.Println("VM started")

    // Wait for VM to be ready
    time.Sleep(30 * time.Second)

    // Get VM IP addresses
    ips, err := vmClient.ListIPs(ctx, group, vmName)
    if err != nil {
        return fmt.Errorf("failed to get IPs: %w", err)
    }
    fmt.Printf("VM IP addresses: %v\n", ips)

    // Run command on VM
    cmdRequest := &compute.VirtualMachineRunCommandRequest{
        Command: stringPtr("hostname && uptime"),
    }
    response, err := vmClient.RunCommand(ctx, group, vmName, cmdRequest)
    if err != nil {
        return fmt.Errorf("failed to run command: %w", err)
    }
    fmt.Printf("Command output:\n%s\n", *response.Output)

    // Stop VM gracefully
    fmt.Println("Stopping VM gracefully...")
    if err := vmClient.StopGraceful(ctx, group, vmName); err != nil {
        return fmt.Errorf("failed to stop VM: %w", err)
    }
    fmt.Println("VM stopped")

    return nil
}

// Helper functions
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }
func int64Ptr(i int64) *int64    { return &i }
func boolPtr(b bool) *bool       { return &b }
```

## Advanced VM Management

### Resize VM

```go
func resizeVM(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    
    // Resize to predefined size
    err := vmClient.Resize(ctx, "production", "web-server-01",
        compute.VirtualMachineSizeTypesStandardA4, nil)
    if err != nil {
        return err
    }
    
    // Or resize with custom size
    customSize := &compute.VirtualMachineCustomSize{
        CpuCount: int32Ptr(8),
        MemoryMB: int32Ptr(16384),
    }
    err = vmClient.Resize(ctx, "production", "web-server-01",
        compute.VirtualMachineSizeTypesCustom, customSize)
    
    return err
}
```

### Attach/Detach Disks

```go
func manageDisk(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    
    // Attach disk
    err := vmClient.DiskAttach(ctx, "production", "web-server-01", "new-data-disk")
    if err != nil {
        return fmt.Errorf("failed to attach disk: %w", err)
    }
    
    // Detach disk
    err = vmClient.DiskDetach(ctx, "production", "web-server-01", "old-data-disk")
    if err != nil {
        return fmt.Errorf("failed to detach disk: %w", err)
    }
    
    return nil
}
```

### VM with High Availability

```go
func createHAVM(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    
    vmSpec := &compute.VirtualMachine{
        Name:     stringPtr("ha-vm-01"),
        Location: stringPtr("eastus"),
        Zones:    &[]string{"zone-1"}, // Place in availability zone
        Properties: &compute.VirtualMachineProperties{
            AvailabilitySet: &compute.SubResource{
                ID: stringPtr("/production/availabilitysets/web-avset"),
            },
            // ... other properties
        },
    }
    
    vm, err := vmClient.CreateOrUpdate(ctx, "production", "ha-vm-01", vmSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created HA VM: %s in zone %s\n", *vm.Name, (*vm.Zones)[0])
    return nil
}
```

## Error Handling

```go
func robustVMOperation(vmClient *virtualmachine.VirtualMachineClient) error {
    ctx := context.Background()
    
    vm, err := vmClient.Get(ctx, "production", "web-server-01")
    if err != nil {
        if errors.IsNotFound(err) {
            // VM doesn't exist, create it
            return createVM(vmClient)
        }
        return fmt.Errorf("unexpected error: %w", err)
    }
    
    // VM exists, update it
    vmToUpdate := &(*vm)[0]
    vmToUpdate.Tags["last-updated"] = stringPtr(time.Now().Format(time.RFC3339))
    
    _, err = vmClient.CreateOrUpdate(ctx, "production", "web-server-01", vmToUpdate)
    return err
}
```

## See Also

- [Compute Services](../services/compute.md) - Full API reference
- [Network Setup](network-setup.md) - Network configuration examples
- [Storage Operations](storage-operations.md) - Storage management examples
