// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	wssdcloudcommon "github.com/microsoft/moc/rpc/common"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

type client struct {
	subID string
	wssdcloudcompute.AvailabilityZoneAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newAvailabilityZoneClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetAvailabilityZoneClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]cloud.AvailabilityZone, error) {
	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilityZoneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getAvailabilityZoneFromResponse(response)
}

// Create
func (c *client) CreateOrUpdate(ctx context.Context, name string, avzone *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error) {
	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_POST, name, avzone)
	if err != nil {
		return nil, err
	}

	_, err = c.Get(ctx, name)
	if err == nil {
		// expect not found
		return nil, errors.Wrapf(errors.AlreadyExists,
			"Type[AvailabilityZone] Name[%s]", name)
	} else if !errors.IsNotFound(err) {
		return nil, err
	}

	response, err := c.AvailabilityZoneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getAvailabilityZoneFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*vmsss) == 0 {
		return nil, fmt.Errorf("creation of availability zone failed to unknown reason")
	}

	return &(*vmsss)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string) error {
	vmss, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*vmss) == 0 {
		return errors.NotFound
	}

	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_DELETE, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.AvailabilityZoneAgentClient.Invoke(ctx, request)
	return err
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getAvailabilityZoneFromResponse(response *wssdcloudcompute.AvailabilityZoneResponse) (*[]cloud.AvailabilityZone, error) {
	avzonesRet := []cloud.AvailabilityZone{}
	for _, avzone := range response.GetAvailabilityZones() {
		cavzone, err := getWssdAvailabilityZone(avzone)
		if err != nil {
			return nil, err
		}
		avzonesRet = append(avzonesRet, *cavzone)
	}

	return &avzonesRet, nil

}

func (c *client) getAvailabilityZoneRequest(opType wssdcloudcommon.Operation, name string, avzone *cloud.AvailabilityZone) (*wssdcloudcompute.AvailabilityZoneRequest, error) {
	request := &wssdcloudcompute.AvailabilityZoneRequest{
		OperationType:    opType,
		AvailabilityZones: []*wssdcloudcompute.AvailabilityZone{},
	}

	if len(name) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Name not specified")
	}

	avzoneRet := &wssdcloudcompute.AvailabilityZone{
		Name:      name,
	}

	if avzone != nil {
		var err error
		avzoneRet, err = getRpcAvailabilityZone(avzone)
		if err != nil {
			return nil, err
		}
	}

	request.AvailabilityZones = append(request.AvailabilityZones, avzoneRet)
	return request, nil

}