// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"

	wssdcloud "github.com/microsoft/moc-proto/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
)

type client struct {
	wssdcloud.NodeAgentClient
}

// newNodeClient - creates a client session with the backend wssd agent
func newNodeClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetNodeClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]cloud.Node, error) {
	request, err := c.getNodeRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.NodeAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getNodeFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, sg *cloud.Node) (node *cloud.Node, err error) {
	err = c.validate(ctx, sg, location)
	if err != nil {
		return
	}
	node = nil

	request, err := c.getNodeRequest(wssdcloudcommon.Operation_POST, location, name, sg)
	if err != nil {
		return
	}

	response, err := c.NodeAgentClient.Invoke(ctx, request)
	if err != nil {
		return
	}
	gps := c.getNodeFromResponse(response)
	if len(*gps) == 0 {
		return
	}

	node = &(*gps)[0]
	return
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	gp, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*gp) == 0 {
		return fmt.Errorf("Node [%s] not found", name)
	}

	request, err := c.getNodeRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*gp)[0])
	if err != nil {
		return err
	}

	_, err = c.NodeAgentClient.Invoke(ctx, request)

	return err
}

///////////////////////////
// Private Methods
func (c *client) validate(ctx context.Context, sg *cloud.Node, location string) (err error) {
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

	if sg.NodeProperties == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing NodeProperties")
		return
	}
	if sg.NodeProperties.FQDN == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing NodeProperties.FQDN")
		return
	}
	if sg.NodeProperties.Port == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing NodeProperties.Port")
		return
	}
	if sg.NodeProperties.AuthorizerPort == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing NodeProperties.AuthorizerPort")
		return
	}
	return

}
func (c *client) getNodeFromResponse(response *wssdcloud.NodeResponse) *[]cloud.Node {
	gps := []cloud.Node{}
	for _, gp := range response.GetNodes() {
		gps = append(gps, *(getNode(gp)))
	}

	return &gps
}

func (c *client) getNodeRequest(opType wssdcloudcommon.Operation, location, name string, gpss *cloud.Node) (*wssdcloud.NodeRequest, error) {
	request := &wssdcloud.NodeRequest{
		OperationType: opType,
		Nodes:         []*wssdcloud.Node{},
	}
	wssdNode := &wssdcloud.Node{
		Name:         name,
		LocationName: location,
	}
	var err error
	if gpss != nil {
		wssdNode, err = getWssdNode(gpss, location)
		if err != nil {
			return nil, err
		}
	}
	request.Nodes = append(request.Nodes, wssdNode)
	return request, nil
}
