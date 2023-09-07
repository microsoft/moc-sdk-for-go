// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package vippool

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/network"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/diagnostics"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudnetwork.VipPoolAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVipPoolClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVipPoolClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get vip pools by name.  If name is nil, get all vip pools
func (c *client) Get(ctx context.Context, location, name string) (*[]network.VipPool, error) {

	request, err := c.getVipPoolRequestByName(ctx, wssdcloudcommon.Operation_GET, location, name)
	if err != nil {
		return nil, err
	}

	response, err := c.VipPoolAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vps, err := c.getVipPoolsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return vps, nil

}

// CreateOrUpdate creates a vip pool if it does not exist, or updates an existing vip pool
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, inputVP *network.VipPool) (*network.VipPool, error) {

	if inputVP == nil || inputVP.VipPoolPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing vip pool Properties")
	}

	request, err := c.getVipPoolRequest(ctx, wssdcloudcommon.Operation_POST, location, name, inputVP)
	if err != nil {
		return nil, err
	}
	response, err := c.VipPoolAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vps, err := c.getVipPoolsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return &(*vps)[0], nil
}

// Delete a vip pool
func (c *client) Delete(ctx context.Context, location, name string) error {
	vps, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*vps) == 0 {
		return fmt.Errorf("vip pool [%s] not found", name)
	}

	request, err := c.getVipPoolRequest(ctx, wssdcloudcommon.Operation_DELETE, location, name, &(*vps)[0])
	if err != nil {
		return err
	}
	_, err = c.VipPoolAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

func (c *client) getVipPoolRequestByName(ctx context.Context, opType wssdcloudcommon.Operation, location, name string) (*wssdcloudnetwork.VipPoolRequest, error) {
	networkVP := network.VipPool{
		Name: &name,
	}
	return c.getVipPoolRequest(ctx, opType, location, name, &networkVP)
}

// getVipPoolRequest converts our internal representation of a vip pool (network.VipPool) into a protobuf request (wssdcloudnetwork.VipPoolRequest) that can be sent to wssdcloudagent
func (c *client) getVipPoolRequest(ctx context.Context, opType wssdcloudcommon.Operation, location, name string, networkVP *network.VipPool) (*wssdcloudnetwork.VipPoolRequest, error) {

	if networkVP == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.VipPoolRequest{
		OperationType: opType,
		VipPools:      []*wssdcloudnetwork.VipPool{},
		Context: &wssdcloudcommon.CallContext{
			CorrelationId: diagnostics.GetCorrelationId(ctx),
		},
	}
	var err error

	wssdCloudVP := &wssdcloudnetwork.VipPool{
		Name:         name,
		LocationName: location,
	}

	if networkVP != nil {
		wssdCloudVP, err = getWssdVipPool(networkVP, location)
		if err != nil {
			return nil, err
		}
	}

	request.VipPools = append(request.VipPools, wssdCloudVP)
	return request, nil
}

// getVipPoolsFromResponse converts a protobuf response from wssdcloudagent (wssdcloudnetwork.VipPoolResponse) to out internal representation of a vip pool (network.VipPool)
func (c *client) getVipPoolsFromResponse(response *wssdcloudnetwork.VipPoolResponse) (*[]network.VipPool, error) {
	networkVPs := []network.VipPool{}

	for _, wssdCloudVP := range response.GetVipPools() {
		networkVP, err := getVipPool(wssdCloudVP)
		if err != nil {
			return nil, err
		}

		networkVPs = append(networkVPs, *networkVP)
	}

	return &networkVPs, nil
}

// getWssdVipPool convert our internal representation of a vippool (network.VipPool) to the cloud vip pool protobuf used by wssdcloudagent (wssdnetwork.VipPool)
func getWssdVipPool(networkVP *network.VipPool, location string) (wssdCloudVP *wssdcloudnetwork.VipPool, err error) {

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if networkVP.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for vip pool")
	}

	wssdCloudVP = &wssdcloudnetwork.VipPool{
		Name:         *networkVP.Name,
		LocationName: location,
	}

	if networkVP.VipPoolPropertiesFormat != nil {

		if networkVP.VipPoolPropertiesFormat.IPPrefix != nil {
			wssdCloudVP.Cidr = *networkVP.VipPoolPropertiesFormat.IPPrefix
		} else {
			if networkVP.VipPoolPropertiesFormat.StartIP == nil || *networkVP.VipPoolPropertiesFormat.StartIP == "" {
				return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing StartIP for vip pool")
			}
			wssdCloudVP.Startip = *networkVP.VipPoolPropertiesFormat.StartIP

			if networkVP.VipPoolPropertiesFormat.EndIP == nil || *networkVP.VipPoolPropertiesFormat.EndIP == "" {
				return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing EndIP for vip pool")
			}
			wssdCloudVP.Endip = *networkVP.VipPoolPropertiesFormat.EndIP
		}
	}

	if networkVP.Version != nil {
		if wssdCloudVP.Status == nil {
			wssdCloudVP.Status = status.InitStatus()
		}
		wssdCloudVP.Status.Version.Number = *networkVP.Version
	}

	return wssdCloudVP, nil
}

// getVipPool converts the cloud vip pool protobuf returned from wssdcloudagent (wssdcloudnetwork.VipPool) to our internal representation of a vippool (network.VipPool)
func getVipPool(wssdVP *wssdcloudnetwork.VipPool) (networkVP *network.VipPool, err error) {
	networkVP = &network.VipPool{
		Name:     &wssdVP.Name,
		Location: &wssdVP.LocationName,
		ID:       &wssdVP.Id,
		Version:  &wssdVP.Status.Version.Number,
		VipPoolPropertiesFormat: &network.VipPoolPropertiesFormat{
			IPPrefix: &wssdVP.Cidr,
			StartIP:  &wssdVP.Startip,
			EndIP:    &wssdVP.Endip,
			Statuses: status.GetStatuses(wssdVP.GetStatus()),
		},
	}

	return networkVP, nil
}
