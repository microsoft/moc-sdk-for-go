// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdcluster

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudcloud.EtcdClusterAgentClient
}

// NewEtcdClusterClient creates a client session with the backend wssdcloud agent
func newEtcdClusterClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetEtcdClusterClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]cloud.EtcdCluster, error) {
	request, err := getEtcdClusterRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.EtcdClusterAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getEtcdClustersFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, cluster *cloud.EtcdCluster) (*cloud.EtcdCluster, error) {
	request, err := getEtcdClusterRequest(wssdcloudcommon.Operation_POST, group, name, cluster)
	if err != nil {
		return nil, err
	}
	response, err := c.EtcdClusterAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	vault := getEtcdClustersFromResponse(response, group)

	if len(*vault) == 0 {
		return nil, fmt.Errorf("[EtcdCluster][Create] Unexpected error: Creating an etcdcluster returned no result")
	}

	return &((*vault)[0]), err
}

// Delete
func (c *client) Delete(ctx context.Context, group, name string) error {
	vault, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vault) == 0 {
		return fmt.Errorf("EtcdCluster [%s] not found", name)
	}

	request, err := getEtcdClusterRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vault)[0])
	if err != nil {
		return err
	}
	_, err = c.EtcdClusterAgentClient.Invoke(ctx, request)
	return err
}

func getEtcdClustersFromResponse(response *wssdcloudcloud.EtcdClusterResponse, group string) *[]cloud.EtcdCluster {
	vaults := []cloud.EtcdCluster{}
	for _, etcdclusters := range response.GetEtcdClusters() {
		vaults = append(vaults, *(getEtcdCluster(etcdclusters, group)))
	}

	return &vaults
}

func getEtcdClusterRequest(opType wssdcloudcommon.Operation, group, name string, vault *cloud.EtcdCluster) (*wssdcloudcloud.EtcdClusterRequest, error) {
	request := &wssdcloudcloud.EtcdClusterRequest{
		OperationType: opType,
		EtcdClusters:  []*wssdcloudcloud.EtcdCluster{},
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	wssdetcdcluster := &wssdcloudcloud.EtcdCluster{
		Name:      name,
		GroupName: group,
	}

	var err error
	if vault != nil {
		wssdetcdcluster, err = getWssdEtcdCluster(vault, group)
		if err != nil {
			return nil, err
		}
	}
	request.EtcdClusters = append(request.EtcdClusters, wssdetcdcluster)
	return request, nil
}
