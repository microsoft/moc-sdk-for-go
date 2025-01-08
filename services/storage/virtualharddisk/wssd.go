// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualharddisk

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/storage"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudstorage "github.com/microsoft/moc/rpc/cloudagent/storage"
	"github.com/microsoft/moc/rpc/common"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
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
	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_GET, group, container, name, nil, "", common.ImageSource_LOCAL_SOURCE)
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
func (c *client) CreateOrUpdate(ctx context.Context, group, container, name string, vhd *storage.VirtualHardDisk, sourcePath string, sourceType common.ImageSource) (*storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_POST, group, container, name, vhd, sourcePath, sourceType)
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

// The hydrate call takes the group name, container name and the name of the disk file. The group is standard input for every call.
// Ultimately, we need the full path on disk to the disk file which we assemble from the path of the container plus the file name of the disk.
// (e.g. "C:\ClusterStorage\Userdata_1\abc123" for the container path and "my_disk.vhd" for the disk name)
func (c *client) Hydrate(ctx context.Context, group, container, name string, vhd *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_HYDRATE, group, container, name, vhd, "", common.ImageSource_LOCAL_SOURCE)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vhds := getVirtualHardDisksFromResponse(response, group)

	if len(*vhds) == 0 {
		return nil, fmt.Errorf("[VirtualHardDisk][Hydrate] Unexpected error: Hydrating a storage interface returned no result")
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
		return fmt.Errorf("[VirtualHardDisk][Delete] %s: not found", name)
	}

	request, err := getVirtualHardDiskRequest(wssdcloudcommon.Operation_DELETE, group, container, name, &(*vhd)[0], "", common.ImageSource_LOCAL_SOURCE)
	if err != nil {
		return err
	}
	_, err = c.VirtualHardDiskAgentClient.Invoke(ctx, request)

	return err

}

func (c *client) Precheck(ctx context.Context, group, container string, vhds []*storage.VirtualHardDisk) (bool, error) {
	request, err := getVirtualHardDiskPrecheckRequest(group, container, vhds)
	if err != nil {
		return false, err
	}

	response, err := c.VirtualHardDiskAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getVirtualHardDiskPrecheckResponse(response)
}

func (c *client) Upload(ctx context.Context, group, container string, vhd *storage.VirtualHardDisk, targetUrl string) error {
	request, err := getVirtualHardDiskOperationRequest(group, container, vhd, targetUrl, wssdcloudcommon.ProviderAccessOperation_VirtualHardDisk_Upload)
	if err != nil {
		return err
	}

	_, err = c.VirtualHardDiskAgentClient.Operate(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func getVirtualHardDiskPrecheckResponse(response *wssdcloudstorage.VirtualHardDiskPrecheckResponse) (bool, error) {
	var err error = nil
	result := response.GetResult().GetValue()
	if !result {
		err = errors.New(response.GetError())
	}
	return result, err
}

func getVirtualHardDiskPrecheckRequest(group, container string, vhds []*storage.VirtualHardDisk) (*wssdcloudstorage.VirtualHardDiskPrecheckRequest, error) {
	request := &wssdcloudstorage.VirtualHardDiskPrecheckRequest{
		VirtualHardDisks: []*wssdcloudstorage.VirtualHardDisk{},
	}
	for _, vhd := range vhds {
		wssdvhd, err := getWssdVirtualHardDisk(vhd, group, container, "", common.ImageSource_LOCAL_SOURCE)
		if err != nil {
			return nil, err
		}
		request.VirtualHardDisks = append(request.VirtualHardDisks, wssdvhd)
	}
	return request, nil
}

func getVirtualHardDiskOperationRequest(group, container string, vhd *storage.VirtualHardDisk, targetUrl string, opType wssdcloudcommon.ProviderAccessOperation) (*wssdcloudstorage.VirtualHardDiskOperationRequest, error) {
	request := &wssdcloudstorage.VirtualHardDiskOperationRequest{
		VirtualHardDisks: []*wssdcloudstorage.VirtualHardDisk{},
		OperationType:    opType,
	}

	var err error

	if vhd == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VirtualHardDisk object is nil")
	}

	wssdvhd, err := getWssdVirtualHardDisk(vhd, group, container, "", common.ImageSource_LOCAL_SOURCE) //sourcePath and SourceType are not used in this context
	if err != nil {
		return nil, err
	}
	wssdvhd.TargetUrl = targetUrl
	request.VirtualHardDisks = append(request.VirtualHardDisks, wssdvhd)

	return request, nil
}

func getVirtualHardDiskRequest(opType wssdcloudcommon.Operation, group, container, name string, storage *storage.VirtualHardDisk, sourcePath string, sourceType common.ImageSource) (*wssdcloudstorage.VirtualHardDiskRequest, error) {
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
		wssdvhd, err = getWssdVirtualHardDisk(storage, group, container, sourcePath, sourceType)
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
