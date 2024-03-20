// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string) error
	Query(context.Context, string, string) (*[]compute.VirtualMachine, error)
	Start(context.Context, string, string) error
	Stop(context.Context, string, string) error
	Pause(context.Context, string, string) error
	Save(context.Context, string, string) error
	RepairGuestAgent(context.Context, string, string) error
	CreateCheckpoint(context.Context, string, string, string) error
	RunCommand(context.Context, string, string, *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error)
	Validate(context.Context, string, string) error
}

type VirtualMachineClient struct {
	compute.BaseClient
	internal   Service
	cloudFQDN  string
	authorizer auth.Authorizer
}

func NewVirtualMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineClient, error) {
	c, err := newVirtualMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineClient{internal: c,
		cloudFQDN:  cloudFQDN,
		authorizer: authorizer,
	}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineClient) Delete(ctx context.Context, group string, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Query method invokes the client Get method and uses the provided query to filter the returned results
func (c *VirtualMachineClient) Query(ctx context.Context, group, query string) (*[]compute.VirtualMachine, error) {
	return c.internal.Query(ctx, group, query)
}

// Start the Virtual Machine
func (c *VirtualMachineClient) Start(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Start(ctx, group, name)
	return
}

// Stop the Virtual Machine
func (c *VirtualMachineClient) Stop(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Stop(ctx, group, name)
	return
}

// Restart the Virtual Machine
func (c *VirtualMachineClient) Restart(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Stop(ctx, group, name)
	if err != nil {
		return
	}
	err = c.internal.Start(ctx, group, name)
	return
}

// Pause the Virtual Machine
func (c *VirtualMachineClient) Pause(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Pause(ctx, group, name)
	return
}

// Save the Virtual Machine
func (c *VirtualMachineClient) Save(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Save(ctx, group, name)
	return
}

type UpdateFunctor interface {
	Update(context.Context, *compute.VirtualMachine) (*compute.VirtualMachine, error)
}

// Update the VM with a retry
func (c *VirtualMachineClient) Update(ctx context.Context, group string, vmName string, updateFunctor UpdateFunctor) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Virtual Machine [%s] not found", vmName)
		}

		vm, err := updateFunctor.Update(ctx, &(*vms)[0])
		if err != nil {
			return err
		}

		_, err = c.CreateOrUpdate(ctx, group, vmName, vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}

// Resize the Virtual Machine
func (c *VirtualMachineClient) Resize(ctx context.Context, group string, vmName string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Virtual Machine [%s] not found", vmName)
		}

		vm := (*vms)[0]
		vm.HardwareProfile.VMSize = newSize
		vm.HardwareProfile.CustomSize = newCustomSize

		_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}

func (c *VirtualMachineClient) DiskAttach(ctx context.Context, group string, vmName, diskName string) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
		}

		vm := (*vms)[0]

		for _, disk := range *vm.StorageProfile.DataDisks {
			if *disk.Vhd.URI == diskName {
				return errors.Wrapf(errors.AlreadyExists, "DataDisk [%s] is already attached to the VM [%s]", diskName, vmName)
			}
		}

		*vm.StorageProfile.DataDisks = append(*vm.StorageProfile.DataDisks, compute.DataDisk{Vhd: &compute.VirtualHardDisk{URI: &diskName}})

		_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}
func (c *VirtualMachineClient) DiskDetach(ctx context.Context, group string, vmName, diskName string) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
		}

		vm := (*vms)[0]

		for i, element := range *vm.StorageProfile.DataDisks {
			if *element.Vhd.URI == diskName {
				*vm.StorageProfile.DataDisks = append((*vm.StorageProfile.DataDisks)[:i], (*vm.StorageProfile.DataDisks)[i+1:]...)
				break
			}
		}

		_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				log.Printf("Retrying because of stale version\n")
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}

func (c *VirtualMachineClient) NetworkInterfaceAdd(ctx context.Context, group string, vmName, nicName string) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
		}

		vm := (*vms)[0]

		for _, nic := range *vm.NetworkProfile.NetworkInterfaces {
			if *nic.ID == nicName {
				return errors.Wrapf(errors.AlreadyExists, "NetworkInterface [%s] is already attached to the VM [%s]", nicName, vmName)
			}
		}

		*vm.NetworkProfile.NetworkInterfaces = append(*vm.NetworkProfile.NetworkInterfaces,
			compute.NetworkInterfaceReference{
				ID: &nicName,
			},
		)

		_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}

func (c *VirtualMachineClient) NetworkInterfaceRemove(ctx context.Context, group string, vmName, nicName string) (err error) {
	for {
		vms, err := c.Get(ctx, group, vmName)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
		}

		vm := (*vms)[0]

		for i, element := range *vm.NetworkProfile.NetworkInterfaces {
			if *element.ID == nicName {
				*vm.NetworkProfile.NetworkInterfaces = append((*vm.NetworkProfile.NetworkInterfaces)[:i], (*vm.NetworkProfile.NetworkInterfaces)[i+1:]...)
				break
			}
		}

		_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
		if err != nil {
			if errors.IsInvalidVersion(err) {
				log.Printf("Retrying because of stale version\n")
				// Retry only on invalid version
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}
		break
	}
	return
}

// Get the Virtual Machine by querying for the specified computer name
func (c *VirtualMachineClient) GetByComputerName(ctx context.Context, group string, computerName string) (*[]compute.VirtualMachine, error) {
	query := fmt.Sprintf("[?virtualmachineproperties.osprofile.computername=='%s']", computerName)

	vms, err := c.Query(ctx, group, query)
	if err != nil {
		return nil, err
	}

	return vms, nil
}

func (c *VirtualMachineClient) RunCommand(ctx context.Context, group, vmName string, request *compute.VirtualMachineRunCommandRequest) (response *compute.VirtualMachineRunCommandResponse, err error) {
	return c.internal.RunCommand(ctx, group, vmName, request)
}

func (c *VirtualMachineClient) RepairGuestAgent(ctx context.Context, group, vmName string) (err error) {
	return c.internal.RepairGuestAgent(ctx, group, vmName)
}

func (c *VirtualMachineClient) CreateCheckpoint(ctx context.Context, group, vmName, checkpointName string) (err error) {
	return c.internal.CreateCheckpoint(ctx, group, vmName, checkpointName)
}

// ListIPs for specified VM
func (c *VirtualMachineClient) ListIPs(ctx context.Context, group, name string) ([]string, error) {
	if len(name) == 0 {
		return nil, errors.Wrap(errors.InvalidInput, "ListIPs requires VM name input")
	}
	vms, err := c.Get(ctx, group, name)
	if err != nil {
		return nil, err
	}
	if len(*vms) == 0 {
		return nil, errors.NotFound
	}

	ips := []string{}
	if (*vms)[0].NetworkProfile == nil {
		return ips, nil
	}

	for _, vmnic := range *(*vms)[0].NetworkProfile.NetworkInterfaces {
		nicCli, err := networkinterface.NewInterfaceClient(c.cloudFQDN, c.authorizer)
		if err != nil {
			return nil, err
		}

		nics, err := nicCli.Get(ctx, group, *vmnic.ID)
		if err != nil {
			return nil, err
		}

		if len(*nics) == 0 || (*nics)[0].IPConfigurations == nil {
			break
		}

		for _, ipConfig := range *(*nics)[0].IPConfigurations {
			if ipConfig.PrivateIPAddress != nil && len(*ipConfig.PrivateIPAddress) > 0 {
				ips = append(ips, *ipConfig.PrivateIPAddress)
			}
		}
	}
	return ips, nil
}

// Validate methods invokes the validate Get method
func (c *VirtualMachineClient) Validate(ctx context.Context, group, name string) error {
	return c.internal.Validate(ctx, group, name)
}
