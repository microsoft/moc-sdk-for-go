// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalmachine

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/marshal"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

type client struct {
	wssdcloudcompute.BareMetalMachineAgentClient
}

// newBareMetalMachineClient - creates a client session with the backend wssdcloud agent
func newBareMetalMachineClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetBareMetalMachineClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]compute.BareMetalMachine, error) {
	request, err := c.getBareMetalMachineRequest(wssdcloudproto.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.BareMetalMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getBareMetalMachineFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, sg *compute.BareMetalMachine) (*compute.BareMetalMachine, error) {
	request, err := c.getBareMetalMachineRequest(wssdcloudproto.Operation_POST, location, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.BareMetalMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	bmms := c.getBareMetalMachineFromResponse(response, location)
	if len(*bmms) == 0 {
		return nil, fmt.Errorf("Creation of Bare Metal Machine failed to unknown reason.")
	}

	return &(*bmms)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	bmms, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*bmms) == 0 {
		return fmt.Errorf("Bare Metal Machine [%s] not found", name)
	}

	request, err := c.getBareMetalMachineRequest(wssdcloudproto.Operation_DELETE, location, name, &(*bmms)[0])
	if err != nil {
		return err
	}
	_, err = c.BareMetalMachineAgentClient.Invoke(ctx, request)

	return err
}

// Query
func (c *client) Query(ctx context.Context, location, query string) (*[]compute.BareMetalMachine, error) {
	bmms, err := c.Get(ctx, location, "")
	if err != nil {
		return nil, err
	}

	filteredBytes, err := config.MarshalOutput(*bmms, query, "json")
	if err != nil {
		return nil, err
	}

	err = marshal.FromJSONBytes(filteredBytes, bmms)
	if err != nil {
		return nil, err
	}

	return bmms, nil
}

// Private methods
func (c *client) getBareMetalMachineFromResponse(response *wssdcloudcompute.BareMetalMachineResponse, location string) *[]compute.BareMetalMachine {
	bmms := []compute.BareMetalMachine{}
	for _, bmm := range response.GetBareMetalMachines() {
		bmms = append(bmms, *(c.getBareMetalMachine(bmm, location)))
	}

	return &bmms
}

func (c *client) getBareMetalMachineRequest(opType wssdcloudproto.Operation, location, name string, bmm *compute.BareMetalMachine) (*wssdcloudcompute.BareMetalMachineRequest, error) {
	request := &wssdcloudcompute.BareMetalMachineRequest{
		OperationType:     opType,
		BareMetalMachines: []*wssdcloudcompute.BareMetalMachine{},
	}
	var err error
	wssdbmm := &wssdcloudcompute.BareMetalMachine{
		Name:         name,
		LocationName: location,
	}
	if bmm != nil {
		wssdbmm, err = c.getWssdBareMetalMachine(bmm, location)
		if err != nil {
			return nil, err
		}
	}
	request.BareMetalMachines = append(request.BareMetalMachines, wssdbmm)
	return request, nil
}

func getComputeTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}
