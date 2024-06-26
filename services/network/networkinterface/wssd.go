// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package networkinterface

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
	subID string
	wssdcloudnetwork.NetworkInterfaceAgentClient
}

// newInterfaceClient - creates a client session with the backend wssdcloud agent
func newInterfaceClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetNetworkInterfaceClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]network.Interface, error) {
	request, err := c.getNetworkInterfaceRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.NetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnetInt, err := c.getInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return vnetInt, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, vnetInterface *network.Interface) (*network.Interface, error) {
	request, err := c.getNetworkInterfaceRequest(wssdcloudcommon.Operation_POST, group, name, vnetInterface)
	if err != nil {
		return nil, err
	}
	response, err := c.NetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnets, err := c.getInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return &(*vnets)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vnetInterface, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vnetInterface) == 0 {
		return fmt.Errorf("Virtual Network Interface [%s] not found", name)
	}

	request, err := c.getNetworkInterfaceRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vnetInterface)[0])
	if err != nil {
		return err
	}
	_, err = c.NetworkInterfaceAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

func (c *client) Precheck(ctx context.Context, group string, resources []*network.Interface) (bool, error) {
	request, err := getNetworkInterfacePrecheckRequest(group, resources)
	if err != nil {
		return false, err
	}
	response, err := c.NetworkInterfaceAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getNetworkInterfacePrecheckResponse(response)
}

// ///////////// private methods  ///////////////
func (c *client) getNetworkInterfaceRequest(opType wssdcloudcommon.Operation, group, name string, networkInterface *network.Interface) (*wssdcloudnetwork.NetworkInterfaceRequest, error) {
	request := &wssdcloudnetwork.NetworkInterfaceRequest{
		OperationType:     opType,
		NetworkInterfaces: []*wssdcloudnetwork.NetworkInterface{},
	}
	var err error

	wssdCloudInterface := &wssdcloudnetwork.NetworkInterface{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if networkInterface != nil {
		wssdCloudInterface, err = getWssdNetworkInterface(networkInterface, group)
		if err != nil {
			return nil, err
		}
	}

	request.NetworkInterfaces = append(request.NetworkInterfaces, wssdCloudInterface)
	return request, nil
}

func (c *client) getInterfacesFromResponse(group string, response *wssdcloudnetwork.NetworkInterfaceResponse) (*[]network.Interface, error) {
	virtualNetworkInterfaces := []network.Interface{}

	for _, vnetInterface := range response.GetNetworkInterfaces() {
		vnetIntf, err := getNetworkInterface(c.subID, group, vnetInterface)
		if err != nil {
			return nil, err
		}

		virtualNetworkInterfaces = append(virtualNetworkInterfaces, *vnetIntf)
	}

	return &virtualNetworkInterfaces, nil
}

func getNetworkInterfacePrecheckRequest(group string, networkInterfaces []*network.Interface) (*wssdcloudnetwork.NetworkInterfacePrecheckRequest, error) {
	request := &wssdcloudnetwork.NetworkInterfacePrecheckRequest{}

	protoNetworkInterfaces := make([]*wssdcloudnetwork.NetworkInterface, 0, len(networkInterfaces))

	for _, networkInterface := range networkInterfaces {
		// can networkInterface ever be nil here? what would be the meaning of that?
		if networkInterface != nil {
			protoNetworkInterface, err := getWssdNetworkInterface(networkInterface, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert NetworkInterface to Protobuf representation")
			}
			protoNetworkInterfaces = append(protoNetworkInterfaces, protoNetworkInterface)
		}
	}

	request.NetworkInterfaces = protoNetworkInterfaces
	return request, nil
}

func getNetworkInterfacePrecheckResponse(response *wssdcloudnetwork.NetworkInterfacePrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}
