// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package macpool

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/network"

	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"
	"github.com/microsoft/moc-proto/pkg/status"
	wssdcloudnetwork "github.com/microsoft/moc-proto/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
)

type client struct {
	wssdcloudnetwork.MacPoolAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newMacPoolClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetMacPoolClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get MAC pools by name.  If name is nil, get all MAC pools
func (c *client) Get(ctx context.Context, location, name string) (*[]network.MACPool, error) {

	request, err := c.getMacPoolRequestByName(wssdcloudcommon.Operation_GET, location, name)
	if err != nil {
		return nil, err
	}

	response, err := c.MacPoolAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	macpools, err := c.getMacPoolsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return macpools, nil

}

// CreateOrUpdate creates a MAC pool if it does not exist, or updates an existing MAC pool
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, inputMacPool *network.MACPool) (*network.MACPool, error) {

	if inputMacPool == nil || inputMacPool.MACPoolPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing MAC pool Properties")
	}

	request, err := c.getMacPoolRequest(wssdcloudcommon.Operation_POST, location, name, inputMacPool)
	if err != nil {
		return nil, err
	}
	response, err := c.MacPoolAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	macpools, err := c.getMacPoolsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return &(*macpools)[0], nil
}

// Delete a MAC pool
func (c *client) Delete(ctx context.Context, location, name string) error {
	macpools, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*macpools) == 0 {
		return fmt.Errorf("MAC pool [%s] not found", name)
	}

	request, err := c.getMacPoolRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*macpools)[0])
	if err != nil {
		return err
	}
	_, err = c.MacPoolAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

func (c *client) getMacPoolRequestByName(opType wssdcloudcommon.Operation, location, name string) (*wssdcloudnetwork.MacPoolRequest, error) {
	networkMacPool := network.MACPool{
		Name: &name,
	}
	return c.getMacPoolRequest(opType, location, name, &networkMacPool)
}

// getMacPoolRequest converts our internal representation of a MAC pool (network.MACPool) into a protobuf request (wssdcloudnetwork.MacPoolRequest) that can be sent to wssdcloudagent
func (c *client) getMacPoolRequest(opType wssdcloudcommon.Operation, location, name string, networkMacPool *network.MACPool) (*wssdcloudnetwork.MacPoolRequest, error) {

	if networkMacPool == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.MacPoolRequest{
		OperationType: opType,
		MacPools:      []*wssdcloudnetwork.MacPool{},
	}
	var err error

	wssdCloudMacPool := &wssdcloudnetwork.MacPool{
		Name:         name,
		LocationName: location,
	}

	if networkMacPool != nil {
		wssdCloudMacPool, err = getWssdMacPool(networkMacPool, location)
		if err != nil {
			return nil, err
		}
	}

	request.MacPools = append(request.MacPools, wssdCloudMacPool)
	return request, nil
}

// getMacPoolsFromResponse converts a protobuf response from wssdcloudagent (wssdcloudnetwork.MacPoolResponse) to out internal representation of a MAC pool (network.MACPool)
func (c *client) getMacPoolsFromResponse(response *wssdcloudnetwork.MacPoolResponse) (*[]network.MACPool, error) {
	networkMacPools := []network.MACPool{}

	for _, wssdCloudMacPool := range response.GetMacPools() {
		networkMacPool, err := getMacPool(wssdCloudMacPool)
		if err != nil {
			return nil, err
		}

		networkMacPools = append(networkMacPools, *networkMacPool)
	}

	return &networkMacPools, nil
}

// getWssdMacPool convert our internal representation of a macpool (network.MACPool) to the cloud MAC pool protobuf used by wssdcloudagent (wssdnetwork.MacPool)
func getWssdMacPool(networkMacPool *network.MACPool, location string) (wssdCloudMacPool *wssdcloudnetwork.MacPool, err error) {

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if networkMacPool.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for MAC pool")
	}

	wssdCloudMacPool = &wssdcloudnetwork.MacPool{
		Name:         *networkMacPool.Name,
		LocationName: location,
	}

	if networkMacPool.Version != nil {
		if wssdCloudMacPool.Status == nil {
			wssdCloudMacPool.Status = status.InitStatus()
		}
		wssdCloudMacPool.Status.Version.Number = *networkMacPool.Version
	}

	if networkMacPool.MACPoolPropertiesFormat != nil && networkMacPool.MACPoolPropertiesFormat.Range != nil &&
		networkMacPool.MACPoolPropertiesFormat.Range.StartMACAddress != nil && networkMacPool.MACPoolPropertiesFormat.Range.EndMACAddress != nil {
		wssdCloudMacPool.Range = &wssdcloudnetwork.MacRange{
			StartMacAddress: *networkMacPool.MACPoolPropertiesFormat.Range.StartMACAddress,
			EndMacAddress:   *networkMacPool.MACPoolPropertiesFormat.Range.EndMACAddress,
		}
	}

	return wssdCloudMacPool, nil
}

// getMacPool converts the cloud MAC pool protobuf returned from wssdcloudagent (wssdcloudnetwork.MacPool) to our internal representation of a macpool (network.MACPool)
func getMacPool(wssdMacPool *wssdcloudnetwork.MacPool) (networkMacPool *network.MACPool, err error) {
	networkMacPool = &network.MACPool{
		Name:     &wssdMacPool.Name,
		Location: &wssdMacPool.LocationName,
		ID:       &wssdMacPool.Id,
		Version:  &wssdMacPool.Status.Version.Number,
		MACPoolPropertiesFormat: &network.MACPoolPropertiesFormat{
			Range: &network.MACRange{
				StartMACAddress: &wssdMacPool.Range.StartMacAddress,
				EndMACAddress:   &wssdMacPool.Range.EndMacAddress,
			},
			Statuses: status.GetStatuses(wssdMacPool.GetStatus()),
		},
	}

	return networkMacPool, nil
}
