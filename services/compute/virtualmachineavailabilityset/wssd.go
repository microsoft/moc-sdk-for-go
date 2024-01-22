// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineavailabilityset

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	wssdcloudcommon "github.com/microsoft/moc/rpc/common"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc-sdk-for-go/services/compute/virtualmachine"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

type client struct {
	subID    string
	vmclient *virtualmachine.VirtualMachineClient
	wssdcloudcompute.AvailabilitySetAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newAvailabilitySetClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetAvailabilitySetClient(&subID, authorizer)
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
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineAvailabilitySet, error) {
	request, err := c.getVirtualMachineAvailabilitySetRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineAvailabilitySetFromResponse(response, group)
}

// GetVirtualMachines
func (c *client) GetVirtualMachines(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineAvailabilitySetRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	vms := []compute.VirtualMachine{}
	for _, vmss := range response.GetAvailabilitySets() {
		for _, vm := range vmss.VirtualMachines {
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
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachineAvailabilitySet) (*compute.VirtualMachineAvailabilitySet, error) {
	request, err := c.getVirtualMachineAvailabilitySetRequest(wssdcloudcommon.Operation_POST, group, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getVirtualMachineAvailabilitySetFromResponse(response, group)
	if err != nil {
		return nil, err
	}

	if len(*vmsss) == 0 {
		return &compute.VirtualMachineAvailabilitySet{}, nil
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

	request, err := c.getVirtualMachineAvailabilitySetRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.AvailabilitySetAgentClient.Invoke(ctx, request)
	return err
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getVirtualMachineAvailabilitySetFromResponse(response *wssdcloudcompute.AvailabilitySetResponse, group string) (*[]compute.VirtualMachineAvailabilitySet, error) {
	vmsss := []compute.VirtualMachineAvailabilitySet{}
	for _, vmss := range response.GetAvailabilitySets() {
		cvmss, err := c.getComputeAvailabilitySet(vmss)
		if err != nil {
			return nil, err
		}
		vmsss = append(vmsss, *cvmss)
	}

	return &vmsss, nil

}

func (c *client) getVirtualMachineAvailabilitySetRequest(opType wssdcloudcommon.Operation, group, name string, vmss *compute.VirtualMachineAvailabilitySet) (*wssdcloudcompute.AvailabilitySetRequest, error) {
	request := &wssdcloudcompute.AvailabilitySetRequest{
		OperationType:    opType,
		AvailabilitySets: []*wssdcloudcompute.AvailabilitySet{},
	}
	var err error
	wssdvmss := &wssdcloudcompute.AvailabilitySet{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if vmss != nil {
		wssdvmss, err = c.getWssdAvailabilitySet(vmss, group)
		if err != nil {
			return nil, err

		}
	}

	request.AvailabilitySets = append(request.AvailabilitySets, wssdvmss)
	return request, nil

}
