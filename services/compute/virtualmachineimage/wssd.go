// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineimage

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"
	wssdcloudcompute "github.com/microsoft/moc-proto/rpc/cloudagent/compute"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
)

type client struct {
	wssdcloudcompute.VirtualMachineImageAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVirtualMachineImageClient(subID string, authorizer auth.Authorizer) (*client, error) {
	return nil, errors.NotImplemented
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineImage, error) {
	request, err := getVirtualMachineImageRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineImageAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualMachineImagesFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, vhd *compute.VirtualMachineImage) (*compute.VirtualMachineImage, error) {
	request, err := getVirtualMachineImageRequest(wssdcloudcommon.Operation_POST, group, name, vhd)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineImageAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vhds := getVirtualMachineImagesFromResponse(response, group)

	if len(*vhds) == 0 {
		return nil, fmt.Errorf("[VirtualMachineImage][Create] Unexpected error: Creating a compute interface returned no result")
	}

	return &((*vhds)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vhd, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vhd) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request, err := getVirtualMachineImageRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vhd)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineImageAgentClient.Invoke(ctx, request)

	return err

}

func getVirtualMachineImageRequest(opType wssdcloudcommon.Operation, group, name string, compute *compute.VirtualMachineImage) (*wssdcloudcompute.VirtualMachineImageRequest, error) {
	request := &wssdcloudcompute.VirtualMachineImageRequest{
		OperationType:        opType,
		VirtualMachineImages: []*wssdcloudcompute.VirtualMachineImage{},
	}

	var err error

	wssdvhd := &wssdcloudcompute.VirtualMachineImage{
		Name:      name,
		GroupName: group,
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if compute != nil {
		wssdvhd, err = getWssdVirtualMachineImage(compute, group)
		if err != nil {
			return nil, err
		}
	}
	request.VirtualMachineImages = append(request.VirtualMachineImages, wssdvhd)

	return request, nil
}

func getVirtualMachineImagesFromResponse(response *wssdcloudcompute.VirtualMachineImageResponse, group string) *[]compute.VirtualMachineImage {
	virtualHardDisks := []compute.VirtualMachineImage{}
	for _, vhd := range response.GetVirtualMachineImages() {
		virtualHardDisks = append(virtualHardDisks, *(getVirtualMachineImage(vhd, group)))
	}

	return &virtualHardDisks
}
