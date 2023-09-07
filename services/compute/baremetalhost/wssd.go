// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalhost

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/diagnostics"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/marshal"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudcompute.BareMetalHostAgentClient
}

// newBareMetalHostClient - creates a client session with the backend wssdcloud agent
func newBareMetalHostClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetBareMetalHostClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]compute.BareMetalHost, error) {
	request, err := c.getBareMetalHostRequest(ctx, wssdcloudproto.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.BareMetalHostAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getBareMetalHostFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, sg *compute.BareMetalHost) (*compute.BareMetalHost, error) {
	request, err := c.getBareMetalHostRequest(ctx, wssdcloudproto.Operation_POST, location, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.BareMetalHostAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	bmhs := c.getBareMetalHostFromResponse(response, location)
	if len(*bmhs) == 0 {
		return nil, fmt.Errorf("Creation of Bare Metal Host failed to unknown reason.")
	}

	return &(*bmhs)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	bmhs, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*bmhs) == 0 {
		return fmt.Errorf("Bare Metal Host [%s] not found", name)
	}

	request, err := c.getBareMetalHostRequest(ctx, wssdcloudproto.Operation_DELETE, location, name, &(*bmhs)[0])
	if err != nil {
		return err
	}
	_, err = c.BareMetalHostAgentClient.Invoke(ctx, request)

	return err
}

// Query
func (c *client) Query(ctx context.Context, location, query string) (*[]compute.BareMetalHost, error) {
	bmhs, err := c.Get(ctx, location, "")
	if err != nil {
		return nil, err
	}

	filteredBytes, err := config.MarshalOutput(*bmhs, query, "json")
	if err != nil {
		return nil, err
	}

	err = marshal.FromJSONBytes(filteredBytes, bmhs)
	if err != nil {
		return nil, err
	}

	return bmhs, nil
}

// Private methods
func (c *client) getBareMetalHostFromResponse(response *wssdcloudcompute.BareMetalHostResponse, location string) *[]compute.BareMetalHost {
	bmhs := []compute.BareMetalHost{}
	for _, bmh := range response.GetBareMetalHosts() {
		bmhs = append(bmhs, *(c.getBareMetalHost(bmh, location)))
	}

	return &bmhs
}

func (c *client) getBareMetalHostRequest(ctx context.Context, opType wssdcloudproto.Operation, location, name string, bmh *compute.BareMetalHost) (*wssdcloudcompute.BareMetalHostRequest, error) {
	request := &wssdcloudcompute.BareMetalHostRequest{
		OperationType:  opType,
		BareMetalHosts: []*wssdcloudcompute.BareMetalHost{},
		Context: &wssdcloudproto.CallContext{
			CorrelationId: diagnostics.GetCorrelationId(ctx),
		},
	}
	var err error
	wssdbmh := &wssdcloudcompute.BareMetalHost{
		Name:         name,
		LocationName: location,
	}
	if bmh != nil {
		wssdbmh, err = c.getWssdBareMetalHost(bmh, location)
		if err != nil {
			return nil, err
		}
	}
	request.BareMetalHosts = append(request.BareMetalHosts, wssdbmh)
	return request, nil
}

func getComputeTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}
