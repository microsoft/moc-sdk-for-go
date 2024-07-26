// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"context"
	"fmt"

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

	_, err = c.Get(ctx, group, name)
	if err == nil {
		// expect not found
		return nil, errors.Wrapf(errors.AlreadyExists,
			"Type[AvailabilitySet] Group[%s] Name[%s]", group, name)
	} else if !errors.IsNotFound(err) {
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
		return nil, fmt.Errorf("creation of availability set failed to unknown reason")
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

func (c *client) Precheck(ctx context.Context, group string, avsets []*compute.AvailabilitySet) (bool, error) {
	request, err := getAvailabilitySetPrecheckRequest(group, avsets)
	if err != nil {
		return false, err
	}
	response, err := c.AvailabilitySetAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getAvailabilitySetPrecheckResponse(response)
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getAvailabilitySetFromResponse(response *wssdcloudcompute.AvailabilitySetResponse, group string) (*[]compute.AvailabilitySet, error) {
	avsetsRet := []compute.AvailabilitySet{}
	for _, avset := range response.GetAvailabilitySets() {
		cavset, err := getWssdAvailabilitySet(avset)
		if err != nil {
			return nil, err
		}
		avsetsRet = append(avsetsRet, *cavset)
	}

	return &avsetsRet, nil

}

func (c *client) getAvailabilitySetRequest(opType wssdcloudcommon.Operation, group, name string, avset *compute.AvailabilitySet) (*wssdcloudcompute.AvailabilitySetRequest, error) {
	request := &wssdcloudcompute.AvailabilitySetRequest{
		OperationType:    opType,
		AvailabilitySets: []*wssdcloudcompute.AvailabilitySet{},
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	avsetRet := &wssdcloudcompute.AvailabilitySet{
		Name:      name,
		GroupName: group,
	}

	if avset != nil {
		var err error
		avsetRet, err = getRpcAvailabilitySet(avset, group)
		if err != nil {
			return nil, err
		}
	}

	request.AvailabilitySets = append(request.AvailabilitySets, avsetRet)
	return request, nil

}

func getAvailabilitySetPrecheckRequest(group string, avsets []*compute.AvailabilitySet) (*wssdcloudcompute.AvailabilitySetPrecheckRequest, error) {
	request := &wssdcloudcompute.AvailabilitySetPrecheckRequest{}

	protoAvSets := make([]*wssdcloudcompute.AvailabilitySet, 0, len(avsets))

	for _, avset := range avsets {
		// can avset ever be nil here? what would be the meaning of that?
		if avset != nil {
			protoAvSet, err := getRpcAvailabilitySet(avset, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert AvailabilitySet to Protobuf representation")
			}
			protoAvSets = append(protoAvSets, protoAvSet)
		}
	}

	request.AvailabilitySets = protoAvSets
	return request, nil
}

func getAvailabilitySetPrecheckResponse(response *wssdcloudcompute.AvailabilitySetPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}
