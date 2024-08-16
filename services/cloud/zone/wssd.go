// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package zone

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
	wssdcloudcompute.ZoneAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newZoneClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetZoneClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location string, name string) (*[]cloud.Zone, error) {
	request, err := c.getZoneRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.ZoneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getZoneFromResponse(response)
}

// Create
func (c *client) CreateOrUpdate(ctx context.Context, location string, name string, avzone *cloud.Zone) (*cloud.Zone, error) {
	request, err := c.getZoneRequest(wssdcloudcommon.Operation_POST, location, name, avzone)
	if err != nil {
		return nil, err
	}

	response, err := c.ZoneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	avzones, err := c.getZoneFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*avzones) == 0 {
		return nil, fmt.Errorf("creation of zone failed to unknown reason")
	}

	return &(*avzones)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location string, name string) error {
	avzones, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*avzones) == 0 {
		return errors.NotFound
	}

	request, err := c.getZoneRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*avzones)[0])
	if err != nil {
		return err
	}
	_, err = c.ZoneAgentClient.Invoke(ctx, request)
	return err
}

func (c *client) Precheck(ctx context.Context, location string, zones []*cloud.Zone) (bool, error) {
	request, err := c.getZonePrecheckRequest(location, zones)
	if err != nil {
		return false, err
	}
	response, err := c.ZoneAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getZonePrecheckResponse(response)
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getZoneFromResponse(response *wssdcloudcompute.ZoneResponse) (*[]cloud.Zone, error) {
	avzonesRet := []cloud.Zone{}
	for _, avzone := range response.GetZones() {
		cavzone, err := getWssdZone(avzone)
		if err != nil {
			return nil, err
		}
		avzonesRet = append(avzonesRet, *cavzone)
	}
	return &avzonesRet, nil

}

func (c *client) getZoneRequest(opType wssdcloudcommon.Operation, location string, name string, avzone *cloud.Zone) (*wssdcloudcompute.ZoneRequest, error) {
	request := &wssdcloudcompute.ZoneRequest{
		OperationType: opType,
		Zones:         []*wssdcloudcompute.Zone{},
	}

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	avzoneRet := &wssdcloudcompute.Zone{
		Name:         name,
		LocationName: location,
	}

	if avzone != nil {
		var err error
		avzoneRet, err = getRpcZone(avzone)
		if err != nil {
			return nil, err
		}
	}

	request.Zones = append(request.Zones, avzoneRet)
	return request, nil

}

func getZonePrecheckResponse(response *wssdcloudcompute.ZonePrecheckResponse) (bool, error) {
	var err error = nil
	result := response.GetResult().GetValue()
	if !result {
		err = errors.New(response.GetError())
	}
	return result, err
}

func (c *client) getZonePrecheckRequest(location string, zones []*cloud.Zone) (*wssdcloudcompute.ZonePrecheckRequest, error) {
	request := &wssdcloudcompute.ZonePrecheckRequest{
		Zones: []*wssdcloudcompute.Zone{},
	}
	for _, zone := range zones {
		if zone != nil {
			moczone, err := getRpcZone(zone)
			if err != nil {
				return nil, err
			}
			request.Zones = append(request.Zones, moczone)
		}
	}
	return request, nil
}
