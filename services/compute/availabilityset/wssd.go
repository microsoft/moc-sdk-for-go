// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	wssdcloudcommon "github.com/microsoft/moc/rpc/common"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

type client struct {
	subID string
	wssdcloudcompute.AvailabilitySetAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newAvailabilitySetClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetAvailabilitySetClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.AvailabilitySet, error) {
	request, err := c.getAvailabilitySetRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getAvailabilitySetFromResponse(response, group)
}

// Create
func (c *client) Create(ctx context.Context, group, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	request, err := c.getAvailabilitySetRequest(wssdcloudcommon.Operation_POST, group, name, avset)
	if err != nil {
		return nil, err
	}
	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getAvailabilitySetFromResponse(response, group)
	if err != nil {
		return nil, err
	}

	if len(*vmsss) == 0 {
		return &compute.AvailabilitySet{}, nil
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

	request, err := c.getAvailabilitySetRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.AvailabilitySetAgentClient.Invoke(ctx, request)
	return err
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getAvailabilitySetFromResponse(response *wssdcloudcompute.AvailabilitySetResponse, group string) (*[]compute.AvailabilitySet, error) {
	vmsss := []compute.AvailabilitySet{}
	for _, vmss := range response.GetAvailabilitySets() {
		cvmss := getComputeAvailabilitySet(vmss)
		vmsss = append(vmsss, *cvmss)
	}

	return &vmsss, nil

}

func (c *client) getAvailabilitySetRequest(opType wssdcloudcommon.Operation, group, name string, vmss *compute.AvailabilitySet) (*wssdcloudcompute.AvailabilitySetRequest, error) {
	request := &wssdcloudcompute.AvailabilitySetRequest{
		OperationType:    opType,
		AvailabilitySets: []*wssdcloudcompute.AvailabilitySet{},
	}
	wssdvmss := &wssdcloudcompute.AvailabilitySet{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if vmss != nil {
		wssdvmss = getWssdAvailabilitySet(vmss, group)
	}

	request.AvailabilitySets = append(request.AvailabilitySets, wssdvmss)
	return request, nil

}
