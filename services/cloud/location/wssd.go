// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package location

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
	wssdcloud.LocationAgentClient
}

// newLocationClient - creates a client session with the backend wssdcloud agent
func newLocationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetLocationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]cloud.Location, error) {
	request, err := c.getLocationRequest(wssdcloudcommon.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.LocationAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getLocationFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, lcn *cloud.Location) (*cloud.Location, error) {
	request, err := c.getLocationRequest(wssdcloudcommon.Operation_POST, name, lcn)
	if err != nil {
		return nil, err
	}
	response, err := c.LocationAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	lcns := c.getLocationFromResponse(response)
	if len(*lcns) == 0 {
		return nil, fmt.Errorf("Creation of Location failed to unknown reason.")
	}

	return &(*lcns)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string) error {
	lcn, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*lcn) == 0 {
		return fmt.Errorf("Location [%s] not found", name)
	}

	request, err := c.getLocationRequest(wssdcloudcommon.Operation_DELETE, name, &(*lcn)[0])
	if err != nil {
		return err
	}

	_, err = c.LocationAgentClient.Invoke(ctx, request)

	return err
}

///////////////////////////
// Private Methods
func (c *client) getLocationFromResponse(response *wssdcloud.LocationResponse) *[]cloud.Location {
	lcns := []cloud.Location{}
	for _, lcn := range response.GetLocations() {
		lcns = append(lcns, *(getLocation(lcn)))
	}

	return &lcns
}

func (c *client) getLocationRequest(opType wssdcloudcommon.Operation, name string, lcnss *cloud.Location) (*wssdcloud.LocationRequest, error) {
	request := &wssdcloud.LocationRequest{
		OperationType: opType,
		Locations:     []*wssdcloud.Location{},
	}
	if lcnss != nil {
		wssdLocation, err := getWssdLocation(lcnss)
		if err != nil {
			return nil, err
		}
		request.Locations = append(request.Locations, wssdLocation)
	} else if len(name) > 0 {
		request.Locations = append(request.Locations,
			&wssdcloud.Location{
				Name: name,
			})
	}
	return request, nil
}
