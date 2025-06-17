// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package container

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/storage"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudstorage "github.com/microsoft/moc/rpc/cloudagent/storage"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudstorage.ContainerAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newContainerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetStorageContainerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]storage.Container, error) {
	request, err := getContainerRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getContainersFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, container *storage.Container) (*storage.Container, error) {
	request, err := getContainerRequest(wssdcloudcommon.Operation_POST, location, name, container)
	if err != nil {
		return nil, err
	}
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	containers := getContainersFromResponse(response, location)

	if len(*containers) == 0 {
		return nil, fmt.Errorf("[Container][Create] Unexpected error: Creating a storage interface returned no result")
	}

	return &((*containers)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	container, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*container) == 0 {
		return fmt.Errorf("Container [%s] not found", name)
	}

	request, err := getContainerRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*container)[0])
	if err != nil {
		return err
	}
	_, err = c.ContainerAgentClient.Invoke(ctx, request)

	return err

}

func (c *client) Precheck(ctx context.Context, location string, containers []*storage.Container) (bool, error) {
	request, err := getContainerPrecheckRequest(location, containers)
	if err != nil {
		return false, err
	}
	response, err := c.ContainerAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getContainerPrecheckResponse(response)
}

func getContainerPrecheckRequest(location string, containers []*storage.Container) (*wssdcloudstorage.ContainerPrecheckRequest, error) {
	request := &wssdcloudstorage.ContainerPrecheckRequest{}

	protoContainers := make([]*wssdcloudstorage.Container, 0, len(containers))

	for _, container := range containers {
		// can container ever be nil here? what would be the meaning of that?
		if container != nil {
			protoContainer, err := getWssdContainer(container, location)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert Container to Protobuf representation")
			}
			protoContainers = append(protoContainers, protoContainer)
		}
	}

	request.Containers = protoContainers
	return request, nil
}

func getContainerPrecheckResponse(response *wssdcloudstorage.ContainerPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}

func getContainerRequest(opType wssdcloudcommon.Operation, location, name string, storage *storage.Container) (*wssdcloudstorage.ContainerRequest, error) {
	request := &wssdcloudstorage.ContainerRequest{
		OperationType: opType,
		Containers:    []*wssdcloudstorage.Container{},
	}

	var err error

	wssdcontainer := &wssdcloudstorage.Container{
		Name:         name,
		LocationName: location,
	}

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if storage != nil {
		wssdcontainer, err = getWssdContainer(storage, location)
		if err != nil {
			return nil, err
		}
	}
	request.Containers = append(request.Containers, wssdcontainer)

	return request, nil
}

func getContainersFromResponse(response *wssdcloudstorage.ContainerResponse, location string) *[]storage.Container {
	virtualHardDisks := []storage.Container{}
	for _, container := range response.GetContainers() {
		virtualHardDisks = append(virtualHardDisks, *(getContainer(container, location)))
	}

	return &virtualHardDisks
}
