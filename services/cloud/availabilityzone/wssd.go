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
func (c *client) Get(ctx context.Context, location string, name string) (*[]cloud.AvailabilityZone, error) {
	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_GET, location, name, nil)
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
func (c *client) CreateOrUpdate(ctx context.Context, location string, name string, avzone *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error) {
	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_POST, location, name, avzone)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilityZoneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	avzones, err := c.getAvailabilityZoneFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*avzones) == 0 {
		return nil, fmt.Errorf("creation of availability zone failed to unknown reason")
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

	request, err := c.getAvailabilityZoneRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*avzones)[0])
	if err != nil {
		return err
	}
	_, err = c.AvailabilityZoneAgentClient.Invoke(ctx, request)
	return err
}

func (c *client) Precheck(ctx context.Context, group string, avzones []*compute.AvailabilityZone) (bool, error) {
	request, err := getAvailabilityZonePrecheckRequest(group, avzones)
	if err != nil {
		return false, err
	}
	response, err := c.AvailabilityZoneAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getAvailabilityZonePrecheckResponse(response)
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

func (c *client) getAvailabilityZoneRequest(opType wssdcloudcommon.Operation, location string, name string, avzone *cloud.AvailabilityZone) (*wssdcloudcompute.AvailabilityZoneRequest, error) {
	request := &wssdcloudcompute.AvailabilityZoneRequest{
		OperationType:    opType,
		AvailabilityZones: []*wssdcloudcompute.AvailabilityZone{},
	}

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	avzoneRet := &wssdcloudcompute.AvailabilityZone{
		Name:      name,
		LocationName: location,
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

func getAvailabilityZonePrecheckRequest(group string, avzones []*compute.AvailabilityZone) (*wssdcloudcompute.AvailabilityZonePrecheckRequest, error) {
	request := &wssdcloudcompute.AvailabilityZonePrecheckRequest{}

	protoAvZones := make([]*wssdcloudcompute.AvailabilityZone, 0, len(avzones))

	for _, avzone := range avzones {
		// can avzone ever be nil here? what would be the meaning of that?
		if avset != nil {
			protoAvZone, err := getRpcAvailabilityZone(avzone, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert AvailabilityZone to Protobuf representation")
			}
			protoAvZones = append(protoAvZones, protoAvZone)
		}
	}

	request.AvailabilityZones = protoAvZones
	return request, nil
}

func getAvailabilityZonePrecheckResponse(response *wssdcloudcompute.AvailabilityZonePrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}