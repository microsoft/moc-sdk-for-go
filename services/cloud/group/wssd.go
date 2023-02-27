// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package group

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloud.GroupAgentClient
}

// newGroupClient - creates a client session with the backend wssdcloud agent
func newGroupClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetGroupClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]cloud.Group, error) {
	request, err := c.getGroupRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.GroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getGroupFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, sg *cloud.Group) (*cloud.Group, error) {
	request, err := c.getGroupRequest(wssdcloudcommon.Operation_POST, location, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.GroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	gps := c.getGroupFromResponse(response)
	if len(*gps) == 0 {
		return nil, fmt.Errorf("Creation of Group failed to unknown reason.")
	}

	return &(*gps)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	gp, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*gp) == 0 {
		return fmt.Errorf("Group [%s] not found", name)
	}

	request, err := c.getGroupRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*gp)[0])
	if err != nil {
		return err
	}

	_, err = c.GroupAgentClient.Invoke(ctx, request)

	return err
}

// /////////////////////////
// Private Methods
func (c *client) getGroupFromResponse(response *wssdcloud.GroupResponse) *[]cloud.Group {
	gps := []cloud.Group{}
	for _, gp := range response.GetGroups() {
		gps = append(gps, *(getGroup(gp)))
	}

	return &gps
}

func (c *client) getGroupRequest(opType wssdcloudcommon.Operation, location, name string, gpss *cloud.Group) (*wssdcloud.GroupRequest, error) {
	request := &wssdcloud.GroupRequest{
		OperationType: opType,
		Groups:        []*wssdcloud.Group{},
	}

	wssdGroup := &wssdcloud.Group{
		Name:         name,
		LocationName: location,
	}

	var err error
	if gpss != nil {
		wssdGroup, err = getWssdGroup(gpss, location)
		if err != nil {
			return nil, err
		}
	}

	request.Groups = append(request.Groups, wssdGroup)
	return request, nil
}
