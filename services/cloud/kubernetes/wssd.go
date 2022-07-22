// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package kubernetes

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	wssdcloudk8s "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudk8s.KubernetesAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newKubernetesClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetKubernetesClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]cloud.Kubernetes, error) {
	request, err := c.getKubernetesRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.KubernetesAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	return c.getKubernetessFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, k8s *cloud.Kubernetes) (*cloud.Kubernetes, error) {
	request, err := c.getKubernetesRequest(wssdcloudcommon.Operation_POST, group, name, k8s)
	if err != nil {
		return nil, err
	}
	response, err := c.KubernetesAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	k8ss := c.getKubernetessFromResponse(response, group)

	if len(*k8ss) == 0 {
		return nil, fmt.Errorf("[Kubernetes][Create] Unexpected error: Creating a cloud interface returned no result")
	}

	return &((*k8ss)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	k8s, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*k8s) == 0 {
		return fmt.Errorf("Kubernetes Cluster [%s] not found", name)
	}

	request, err := c.getKubernetesRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*k8s)[0])
	if err != nil {
		return err
	}
	_, err = c.KubernetesAgentClient.Invoke(ctx, request)
	services.HandleGRPCError(err)
	return err
}

func (c *client) getKubernetesRequest(opType wssdcloudcommon.Operation, group, name string, cloud *cloud.Kubernetes) (*wssdcloudk8s.KubernetesRequest, error) {
	request := &wssdcloudk8s.KubernetesRequest{
		OperationType: opType,
		Kubernetess:   []*wssdcloudk8s.Kubernetes{},
	}

	var err error

	wssdcloud := &wssdcloudk8s.Kubernetes{
		GroupName: group,
		Name:      name,
	}

	if cloud != nil {
		wssdcloud, err = c.getWssdKubernetes(cloud, group)
		if err != nil {
			return nil, err
		}
	}
	request.Kubernetess = append(request.Kubernetess, wssdcloud)

	return request, nil
}

func (c *client) getKubernetessFromResponse(response *wssdcloudk8s.KubernetesResponse, group string) *[]cloud.Kubernetes {
	kubes := []cloud.Kubernetes{}
	for _, k8s := range response.GetKubernetess() {
		kubes = append(kubes, *(c.getKubernetes(k8s)))
	}

	return &kubes
}
