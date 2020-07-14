// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachinescaleset

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"

	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"

	wssdcloudcompute "github.com/microsoft/moc-proto/rpc/cloudagent/compute"
	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
)

type client struct {
	subID    string
	vmclient *virtualmachine.VirtualMachineClient
	wssdcloudcompute.VirtualMachineScaleSetAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVirtualMachineScaleSetClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVirtualMachineScaleSetClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	vmc, err := virtualmachine.NewVirtualMachineClient(subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{subID, vmc, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineScaleSetFromResponse(response, group)
}

// GetVirtualMachines
func (c *client) GetVirtualMachines(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	vms := []compute.VirtualMachine{}
	for _, vmss := range response.GetVirtualMachineScaleSetSystems() {
		for _, vm := range vmss.GetVirtualMachineSystems() {
			tvms, err := c.vmclient.Get(ctx, group, vm.Name)
			if err != nil {
				return nil, err
			}
			if tvms == nil || len(*tvms) == 0 {
				return nil, fmt.Errorf("Vmss doesnt have any Vms")
			}
			// FIXME: Make sure Vms only on this scale set is returned.
			// If another Vm with the same name exists, that could also potentially be returned.
			vms = append(vms, (*tvms)[0])
		}
	}

	return &vms, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcloudcommon.Operation_POST, group, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getVirtualMachineScaleSetFromResponse(response, group)
	if err != nil {
		return nil, err
	}

	if len(*vmsss) == 0 {
		return &compute.VirtualMachineScaleSet{}, nil
	}

	return &(*vmsss)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vmss, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vmss) == 0 {
		return errors.NotFound
	}

	request, err := c.getVirtualMachineScaleSetRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	return err
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getVirtualMachineScaleSetFromResponse(response *wssdcloudcompute.VirtualMachineScaleSetResponse, group string) (*[]compute.VirtualMachineScaleSet, error) {
	vmsss := []compute.VirtualMachineScaleSet{}
	for _, vmss := range response.GetVirtualMachineScaleSetSystems() {
		cvmss, err := c.getVirtualMachineScaleSet(vmss, group)
		if err != nil {
			return nil, err
		}
		vmsss = append(vmsss, *cvmss)
	}

	return &vmsss, nil

}

func (c *client) getVirtualMachineScaleSetRequest(opType wssdcloudcommon.Operation, group, name string, vmss *compute.VirtualMachineScaleSet) (*wssdcloudcompute.VirtualMachineScaleSetRequest, error) {
	request := &wssdcloudcompute.VirtualMachineScaleSetRequest{
		OperationType:                 opType,
		VirtualMachineScaleSetSystems: []*wssdcloudcompute.VirtualMachineScaleSet{},
	}
	var err error
	wssdvmss := &wssdcloudcompute.VirtualMachineScaleSet{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if vmss != nil {
		wssdvmss, err = c.getWssdVirtualMachineScaleSet(vmss, group)
		if err != nil {
			return nil, err

		}
	}

	request.VirtualMachineScaleSetSystems = append(request.VirtualMachineScaleSetSystems, wssdvmss)
	return request, nil

}
