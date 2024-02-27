// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	computeavset "github.com/microsoft/moc-sdk-for-go/services/compute/availabilityset"
	"github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
)

type VirtualMachineMockClient struct {
	vmClient      VirtualMachineClient
	avsetInternal computeavset.Service
}

func NewVirtualMachineMockClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineMockClient, error) {
	c, err := NewVirtualMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	vmClient := &VirtualMachineClient{internal: c,
		cloudFQDN:  cloudFQDN,
		authorizer: authorizer,
	}
	avsetMockService, err := computeavset.NewAvailabilitySetMockClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}
	mockClient := &VirtualMachineMockClient{
		avsetInternal: avsetMockService,
		vmClient:      *vmClient,
	}

	return mockClient, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineMockClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.vmClient.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineMockClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {

	vm, err := c.vmClient.internal.CreateOrUpdate(ctx, group, name, compute)
	if err == nil && compute.AvailabilitySetProfile != nil {
		err = c.attachVmToAvset(ctx, group, compute)
		return vm, err
	}
	return vm, err
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineMockClient) Delete(ctx context.Context, group string, name string) error {
	err := c.vmClient.internal.Delete(ctx, group, name)
	if err != nil {
		return err
	}
	return c.detachVmFromAvset(ctx, group, name)
}

func (c *VirtualMachineMockClient) validateVmAvst(ctx context.Context, group string, vm *compute.VirtualMachine) (err error) {
	if vm.AvailabilitySetProfile != nil {
		_, err := c.avsetInternal.Get(ctx, *vm.AvailabilitySetProfile.GroupName, *vm.AvailabilitySetProfile.Name)
		if err != nil {
			// avset needs to be created first
			return err
		}
		attachedAvset, err := c.getVirtualMachineAvailabilitySet(ctx, group, *vm.Name)
		if errors.IsNotFound(err) {
			return nil
		}
		if attachedAvset.Name != vm.AvailabilitySetProfile.Name {
			return errors.Wrapf(errors.InvalidConfiguration, "updating avset for vm is not allowed")
		}
		return nil
	}
	return nil
}

func (c *VirtualMachineMockClient) attachVmToAvset(ctx context.Context, group string, vm *compute.VirtualMachine) (err error) {
	if vm.AvailabilitySetProfile == nil {
		return
	}
	avset, err := c.avsetInternal.Get(ctx, *vm.AvailabilitySetProfile.GroupName, *vm.AvailabilitySetProfile.Name)
	if err != nil {
		return err
	}
	vmAvset := &(*avset)[0]
	newVmReference := &compute.VirtualMachineReference{
		GroupName: &group,
		Name:      vm.Name,
	}
	vmAvset.VirtualMachines = append(vmAvset.VirtualMachines, newVmReference)
	_, err = c.avsetInternal.Create(ctx, group, *vmAvset.Name, vmAvset)
	if err != nil {
		return err
	}
	return nil
}

func (c *VirtualMachineMockClient) detachVmFromAvset(ctx context.Context, group string, name string) (err error) {
	vmAvset, err := c.getVirtualMachineAvailabilitySet(ctx, group, name)
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}
	newVmSets := []*compute.VirtualMachineReference{}
	for _, vm := range vmAvset.VirtualMachines {
		if *vm.Name != name {
			newVmSets = append(newVmSets, vm)
		}
	}
	vmAvset.VirtualMachines = newVmSets
	_, err = c.avsetInternal.Create(ctx, group, *vmAvset.Name, vmAvset)
	return
}

func (c *VirtualMachineMockClient) getVirtualMachineAvailabilitySet(ctx context.Context, group, vmName string) (*compute.AvailabilitySet, error) {
	avsets, err := c.avsetInternal.Get(ctx, group, "")
	if err != nil {
		return nil, err
	}
	for _, avset := range *avsets {
		for _, vmSubResource := range avset.VirtualMachines {
			if group == *vmSubResource.GroupName && vmName == *vmSubResource.Name {
				return &avset, nil
			}
		}
	}
	return nil, errors.NotFound
}

// Query method invokes the client Get method and uses the provided query to filter the returned results
func (c *VirtualMachineMockClient) Query(ctx context.Context, group, query string) (*[]compute.VirtualMachine, error) {
	return c.vmClient.internal.Query(ctx, group, query)
}

// Start the Virtual Machine
func (c *VirtualMachineMockClient) Start(ctx context.Context, group string, name string) (err error) {
	err = c.vmClient.internal.Start(ctx, group, name)
	return
}

// Stop the Virtual Machine
func (c *VirtualMachineMockClient) Stop(ctx context.Context, group string, name string) (err error) {
	err = c.vmClient.internal.Stop(ctx, group, name)
	return
}

// Restart the Virtual Machine
func (c *VirtualMachineMockClient) Restart(ctx context.Context, group string, name string) (err error) {
	err = c.vmClient.internal.Stop(ctx, group, name)
	if err != nil {
		return
	}
	err = c.vmClient.internal.Start(ctx, group, name)
	return
}

// Update the VM with a retry
func (c *VirtualMachineMockClient) Update(ctx context.Context, group string, vmName string, updateFunctor UpdateFunctor) (err error) {
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
func (c *VirtualMachineMockClient) Resize(ctx context.Context, group string, vmName string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize) (err error) {
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

func (c *VirtualMachineMockClient) DiskAttach(ctx context.Context, group string, vmName, diskName string) (err error) {
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
func (c *VirtualMachineMockClient) DiskDetach(ctx context.Context, group string, vmName, diskName string) (err error) {
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

func (c *VirtualMachineMockClient) NetworkInterfaceAdd(ctx context.Context, group string, vmName, nicName string) (err error) {
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

func (c *VirtualMachineMockClient) NetworkInterfaceRemove(ctx context.Context, group string, vmName, nicName string) (err error) {
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
func (c *VirtualMachineMockClient) GetByComputerName(ctx context.Context, group string, computerName string) (*[]compute.VirtualMachine, error) {
	query := fmt.Sprintf("[?virtualmachineproperties.osprofile.computername=='%s']", computerName)

	vms, err := c.Query(ctx, group, query)
	if err != nil {
		return nil, err
	}

	return vms, nil
}

func (c *VirtualMachineMockClient) RunCommand(ctx context.Context, group, vmName string, request *compute.VirtualMachineRunCommandRequest) (response *compute.VirtualMachineRunCommandResponse, err error) {
	return c.vmClient.internal.RunCommand(ctx, group, vmName, request)
}

func (c *VirtualMachineMockClient) RepairGuestAgent(ctx context.Context, group, vmName string) (err error) {
	return c.vmClient.internal.RepairGuestAgent(ctx, group, vmName)
}

// ListIPs for specified VM
func (c *VirtualMachineMockClient) ListIPs(ctx context.Context, group, name string) ([]string, error) {
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
		nicCli, err := networkinterface.NewInterfaceClient(c.vmClient.cloudFQDN, c.vmClient.authorizer)
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
func (c *VirtualMachineMockClient) Validate(ctx context.Context, group, name string) error {
	return c.vmClient.internal.Validate(ctx, group, name)
}
