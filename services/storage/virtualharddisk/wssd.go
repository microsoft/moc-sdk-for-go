// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualharddisk

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"
	wssdcloudstorage "github.com/microsoft/moc-proto/rpc/cloudagent/storage"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/storage"
)

type client struct {
	wssdcloudstorage.VirtualHardDiskAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVirtualHardDiskClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVirtualHardDiskClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, container, name string) (*[]storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_GET, group, container, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualHardDisksFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, container, name string, vhd *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_POST, group, container, name, vhd)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vhds := getVirtualHardDisksFromResponse(response, group)

	if len(*vhds) == 0 {
		return nil, fmt.Errorf("[VirtualHardDisk][Create] Unexpected error: Creating a storage interface returned no result")
	}

	return &((*vhds)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, container, name string) error {
	vhd, err := c.Get(ctx, group, container, name)
	if err != nil {
		return err
	}
	if len(*vhd) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_DELETE, group, container, name, &(*vhd)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualHardDiskAgentClient.Invoke(ctx, request)

	return err

}

func getVirtualHardDiskRequest(opType wssdcloudcommon.Operation, group, container, name string, storage *storage.VirtualHardDisk) (*wssdcloudstorage.VirtualHardDiskRequest, error) {
	request := &wssdcloudstorage.VirtualHardDiskRequest{
		OperationType:    opType,
		VirtualHardDisks: []*wssdcloudstorage.VirtualHardDisk{},
	}

	var err error

	wssdvhd := &wssdcloudstorage.VirtualHardDisk{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if storage != nil {
		wssdvhd, err = getWssdVirtualHardDisk(storage, group, container)
		if err != nil {
			return nil, err
		}
	}
	request.VirtualHardDisks = append(request.VirtualHardDisks, wssdvhd)

	return request, nil
}

func getVirtualHardDisksFromResponse(response *wssdcloudstorage.VirtualHardDiskResponse, group string) *[]storage.VirtualHardDisk {
	virtualHardDisks := []storage.VirtualHardDisk{}
	for _, vhd := range response.GetVirtualHardDisks() {
		virtualHardDisks = append(virtualHardDisks, *(getVirtualHardDisk(vhd, group)))
	}

	return &virtualHardDisks
}
