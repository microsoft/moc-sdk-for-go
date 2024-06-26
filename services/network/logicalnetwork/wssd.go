// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package logicalnetwork

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudnetwork.LogicalNetworkAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newLogicalNetworkClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetLogicalNetworkClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]network.LogicalNetwork, error) {
	request, err := getLogicalNetworkRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.LogicalNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getLogicalNetworksFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, lnet *network.LogicalNetwork) (*network.LogicalNetwork, error) {
	request, err := getLogicalNetworkRequest(wssdcloudcommon.Operation_POST, location, name, lnet)
	if err != nil {
		return nil, err
	}
	response, err := c.LogicalNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	lnets := getLogicalNetworksFromResponse(response, location)

	if len(*lnets) == 0 {
		return nil, fmt.Errorf("[LogicalNetwork][Create] Unexpected error: Creating a Logical Network returned no result")
	}

	return &((*lnets)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	lnet, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*lnet) == 0 {
		return fmt.Errorf("[LogicalNetwork][Delete] Logical Network [%s] not found", name)
	}

	request, err := getLogicalNetworkRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*lnet)[0])
	if err != nil {
		return err
	}
	_, err = c.LogicalNetworkAgentClient.Invoke(ctx, request)

	return err
}

func (c *client) Precheck(ctx context.Context, location string, logicalNetworks []*network.LogicalNetwork) (bool, error) {
	request, err := getLogicalNetworkPrecheckRequest(location, logicalNetworks)
	if err != nil {
		return false, err
	}
	response, err := c.LogicalNetworkAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getLogicalNetworkPrecheckResponse(response)
}

func getLogicalNetworkPrecheckRequest(location string, logicalNetworks []*network.LogicalNetwork) (*wssdcloudnetwork.LogicalNetworkPrecheckRequest, error) {
	request := &wssdcloudnetwork.LogicalNetworkPrecheckRequest{}

	protoLogicalNetworks := make([]*wssdcloudnetwork.LogicalNetwork, 0, len(logicalNetworks))

	for _, logicalNetwork := range logicalNetworks {
		// can logical network ever be nil here? what would be the meaning of that?
		if logicalNetwork != nil {

			// TODO (aweston): double check this
			if logicalNetwork.Location == nil {
				logicalNetwork.Location = &location
			}

			protoLogicalNetwork, err := getWssdLogicalNetwork(logicalNetwork)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert LogicalNetwork to Protobuf representation")
			}
			protoLogicalNetworks = append(protoLogicalNetworks, protoLogicalNetwork)
		}
	}

	request.LogicalNetworks = protoLogicalNetworks
	return request, nil
}

func getLogicalNetworkPrecheckResponse(response *wssdcloudnetwork.LogicalNetworkPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}

func getLogicalNetworkRequest(opType wssdcloudcommon.Operation, location, name string, network *network.LogicalNetwork) (*wssdcloudnetwork.LogicalNetworkRequest, error) {
	request := &wssdcloudnetwork.LogicalNetworkRequest{
		OperationType:   opType,
		LogicalNetworks: []*wssdcloudnetwork.LogicalNetwork{},
	}

	var err error

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location is not specified")
	}

	wssdnetwork := &wssdcloudnetwork.LogicalNetwork{
		Name:         name,
		LocationName: location,
	}

	if network != nil {
		wssdnetwork, err = getWssdLogicalNetwork(network)
		if err != nil {
			return nil, err
		}
	}
	request.LogicalNetworks = append(request.LogicalNetworks, wssdnetwork)

	return request, nil
}

func getLogicalNetworksFromResponse(response *wssdcloudnetwork.LogicalNetworkResponse, location string) *[]network.LogicalNetwork {
	logicalNetworks := []network.LogicalNetwork{}
	for _, lnet := range response.GetLogicalNetworks() {
		logicalNetworks = append(logicalNetworks, *(getLogicalNetwork(lnet)))
	}

	return &logicalNetworks
}
