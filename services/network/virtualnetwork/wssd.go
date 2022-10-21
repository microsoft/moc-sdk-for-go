// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualnetwork

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/errors"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudnetwork.VirtualNetworkAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVirtualNetworkClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVirtualNetworkClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]network.VirtualNetwork, error) {
	request, err := getVirtualNetworkRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualNetworksFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, vnet *network.VirtualNetwork) (*network.VirtualNetwork, error) {
	request, err := getVirtualNetworkRequest(wssdcloudcommon.Operation_POST, group, name, vnet)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnets := getVirtualNetworksFromResponse(response, group)

	if len(*vnets) == 0 {
		return nil, fmt.Errorf("[VirtualNetwork][Create] Unexpected error: Creating a network interface returned no result")
	}

	return &((*vnets)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vnet, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vnet) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request, err := getVirtualNetworkRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vnet)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualNetworkAgentClient.Invoke(ctx, request)

	return err
}

func getVirtualNetworkRequest(opType wssdcloudcommon.Operation, group, name string, network *network.VirtualNetwork) (*wssdcloudnetwork.VirtualNetworkRequest, error) {
	request := &wssdcloudnetwork.VirtualNetworkRequest{
		OperationType:   opType,
		VirtualNetworks: []*wssdcloudnetwork.VirtualNetwork{},
	}

	var err error

	wssdnetwork := &wssdcloudnetwork.VirtualNetwork{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if network != nil {
		wssdnetwork, err = getWssdVirtualNetwork(network, group)
		if err != nil {
			return nil, err
		}
	}
	request.VirtualNetworks = append(request.VirtualNetworks, wssdnetwork)

	return request, nil
}

func getVirtualNetworksFromResponse(response *wssdcloudnetwork.VirtualNetworkResponse, group string) *[]network.VirtualNetwork {
	virtualNetworks := []network.VirtualNetwork{}
	for _, vnet := range response.GetVirtualNetworks() {
		virtualNetworks = append(virtualNetworks, *(getVirtualNetwork(vnet, group)))
	}

	return &virtualNetworks
}
