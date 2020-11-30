// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package controlplane

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloud.ControlPlaneAgentClient
}

// newControlPlaneClient - creates a client session with the backend wssd agent
func newControlPlaneClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetControlPlaneClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]cloud.ControlPlaneInfo, error) {
	request, err := c.getControlPlaneRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.ControlPlaneAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getControlPlaneFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, sg *cloud.ControlPlaneInfo) (controlPlane *cloud.ControlPlaneInfo, err error) {
	err = c.validate(ctx, sg, location)
	if err != nil {
		return
	}
	controlPlane = nil

	request, err := c.getControlPlaneRequest(wssdcloudcommon.Operation_POST, location, name, sg)
	if err != nil {
		return
	}

	response, err := c.ControlPlaneAgentClient.Invoke(ctx, request)
	if err != nil {
		return
	}
	gps := c.getControlPlaneFromResponse(response)
	if len(*gps) == 0 {
		return
	}

	controlPlane = &(*gps)[0]
	return
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	gp, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*gp) == 0 {
		return fmt.Errorf("ControlPlane [%s] not found", name)
	}

	request, err := c.getControlPlaneRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*gp)[0])
	if err != nil {
		return err
	}

	_, err = c.ControlPlaneAgentClient.Invoke(ctx, request)

	return err
}

///////////////////////////
// Private Methods
func (c *client) validate(ctx context.Context, sg *cloud.ControlPlaneInfo, location string) (err error) {
	if sg == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input is nil")
		return
	}
	if sg.Location == nil && len(location) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "Location is nil")
		return
	} else if sg.Location == nil {
		sg.Location = &location
	}

	if sg.ControlPlaneProperties == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing ControlPlaneProperties")
		return
	}
	if sg.ControlPlaneProperties.FQDN == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing ControlPlaneProperties.Fqdn")
		return
	}
	if sg.ControlPlaneProperties.Port == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing ControlPlaneProperties.Port")
		return
	}
	return

}
func (c *client) getControlPlaneFromResponse(response *wssdcloud.ControlPlaneResponse) *[]cloud.ControlPlaneInfo {
	gps := []cloud.ControlPlaneInfo{}
	for _, gp := range response.GetControlPlanes() {
		gps = append(gps, *(getControlPlane(gp)))
	}

	return &gps
}

func (c *client) getControlPlaneRequest(opType wssdcloudcommon.Operation, location, name string, cpRequest *cloud.ControlPlaneInfo) (*wssdcloud.ControlPlaneRequest, error) {
	request := &wssdcloud.ControlPlaneRequest{
		OperationType: opType,
		ControlPlanes: []*wssdcloud.ControlPlane{},
	}
	wssdControlPlane := &wssdcloud.ControlPlane{
		Name:         name,
		LocationName: location,
	}
	var err error
	if cpRequest != nil {
		wssdControlPlane, err = getWssdControlPlane(cpRequest, location)
		if err != nil {
			return nil, err
		}
	}
	request.ControlPlanes = append(request.ControlPlanes, wssdControlPlane)
	return request, nil
}
