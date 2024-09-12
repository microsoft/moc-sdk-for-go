// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineimage

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudcompute.VirtualMachineImageAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newVirtualMachineImageClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVirtualMachineImageClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
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

func (c *client) Precheck(ctx context.Context, group string, virtualMachineImages []*compute.VirtualMachineImage) (bool, error) {
	request, err := getVirtualMachineImagePrecheckRequest(group, virtualMachineImages)
	if err != nil {
		return false, err
	}
	response, err := c.VirtualMachineImageAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getVirtualMachineImagePrecheckResponse(response)
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

func getVirtualMachineImagePrecheckRequest(group string, vmImages []*compute.VirtualMachineImage) (*wssdcloudcompute.VirtualMachineImagePrecheckRequest, error) {
	request := &wssdcloudcompute.VirtualMachineImagePrecheckRequest{}

	protoVMImages := make([]*wssdcloudcompute.VirtualMachineImage, 0, len(vmImages))

	for _, vmImage := range vmImages {
		// can vm image ever be nil here? what would be the meaning of that?
		if vmImage != nil {
			protoVMImage, err := getWssdVirtualMachineImage(vmImage, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert VirtualMachineImage to Protobuf representation")
			}
			protoVMImages = append(protoVMImages, protoVMImage)
		}
	}

	request.VirtualMachineImages = protoVMImages
	return request, nil
}

func getVirtualMachineImagePrecheckResponse(response *wssdcloudcompute.VirtualMachineImagePrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}
