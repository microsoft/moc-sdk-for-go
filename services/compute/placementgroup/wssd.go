// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package placementgroup

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
	wssdcloudcompute.PlacementGroupAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newPlacementGroupClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetPlacementGroupClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.PlacementGroup, error) {
	request, err := c.getPlacementGroupRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.PlacementGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getPlacementGroupFromResponse(response, group)
}

// Create
func (c *client) Create(ctx context.Context, group, name string, pgroup *compute.PlacementGroup) (*compute.PlacementGroup, error) {
	request, err := c.getPlacementGroupRequest(wssdcloudcommon.Operation_POST, group, name, pgroup)
	if err != nil {
		return nil, err
	}

	_, err = c.Get(ctx, group, name)
	if err == nil {
		// expect not found
		return nil, errors.Wrapf(errors.AlreadyExists,
			"Type[PlacementGroup] Group[%s] Name[%s]", group, name)
	} else if !errors.IsNotFound(err) {
		return nil, err
	}

	response, err := c.PlacementGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getPlacementGroupFromResponse(response, group)
	if err != nil {
		return nil, err
	}

	if len(*vmsss) == 0 {
		return nil, fmt.Errorf("creation of placement group failed to unknown reason")
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

	request, err := c.getPlacementGroupRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.PlacementGroupAgentClient.Invoke(ctx, request)
	return err
}

func (c *client) Precheck(ctx context.Context, group string, pgroups []*compute.PlacementGroup) (bool, error) {
	request, err := getPlacementGroupPrecheckRequest(group, pgroups)
	if err != nil {
		return false, err
	}
	response, err := c.PlacementGroupAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getPlacementGroupPrecheckResponse(response)
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getPlacementGroupFromResponse(response *wssdcloudcompute.PlacementGroupResponse, group string) (*[]compute.PlacementGroup, error) {
	pgroupsRet := []compute.PlacementGroup{}
	for _, pgroup := range response.GetPlacementGroups() {
		cpgroup, err := getWssdPlacementGroup(pgroup)
		if err != nil {
			return nil, err
		}
		pgroupsRet = append(pgroupsRet, *cpgroup)
	}

	return &pgroupsRet, nil

}

func (c *client) getPlacementGroupRequest(opType wssdcloudcommon.Operation, group, name string, pgroup *compute.PlacementGroup) (*wssdcloudcompute.PlacementGroupRequest, error) {
	request := &wssdcloudcompute.PlacementGroupRequest{
		OperationType:   opType,
		PlacementGroups: []*wssdcloudcompute.PlacementGroup{},
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	pgroupRet := &wssdcloudcompute.PlacementGroup{
		Name:      name,
		GroupName: group,
	}

	if pgroup != nil {
		var err error
		pgroupRet, err = getRpcPlacementGroup(pgroup, group)
		if err != nil {
			return nil, err
		}
	}

	request.PlacementGroups = append(request.PlacementGroups, pgroupRet)
	return request, nil

}

func getPlacementGroupPrecheckRequest(group string, pgroups []*compute.PlacementGroup) (*wssdcloudcompute.PlacementGroupPrecheckRequest, error) {
	request := &wssdcloudcompute.PlacementGroupPrecheckRequest{}

	protoPGroups := make([]*wssdcloudcompute.PlacementGroup, 0, len(pgroups))

	for _, pgroup := range pgroups {
		// can pgroup ever be nil here? what would be the meaning of that?
		if pgroup != nil {
			protoPGroup, err := getRpcPlacementGroup(pgroup, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert PlacementGroup to Protobuf representation")
			}
			protoPGroups = append(protoPGroups, protoPGroup)
		}
	}

	request.PlacementGroups = protoPGroups
	return request, nil
}

func getPlacementGroupPrecheckResponse(response *wssdcloudcompute.PlacementGroupPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}
