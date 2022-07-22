// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdserver

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/cloud/etcdcluster"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudcloud.EtcdServerAgentClient
}

// NewEtcdServerClient - creates a client session with the backend wssdcloud agent
func newEtcdServerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetEtcdServerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name, clusterName string) (*[]etcdcluster.EtcdServer, error) {
	request, err := getEtcdServerRequest(wssdcloudcommon.Operation_GET, name, clusterName, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.EtcdServerAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	return getEtcdServersFromResponse(response, clusterName), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, server *etcdcluster.EtcdServer) (*etcdcluster.EtcdServer, error) {
	err := c.validate(ctx, group, name, server)
	if err != nil {
		return nil, err
	}
	request, err := getEtcdServerRequest(wssdcloudcommon.Operation_POST, name, *server.ClusterName, server)
	if err != nil {
		return nil, err
	}
	response, err := c.EtcdServerAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, errors.Wrapf(err, "EtcdServer Create failed")
	}

	servers := getEtcdServersFromResponse(response, *server.ClusterName)

	if len(*servers) == 0 {
		return nil, fmt.Errorf("[EtcdServer][Create] Unexpected error: Creating an etcdserver returned no result")
	}

	return &((*servers)[0]), err
}

func (c *client) validate(ctx context.Context, group, name string, server *etcdcluster.EtcdServer) (err error) {
	if server == nil || server.ClusterName == nil {
		return errors.Wrapf(errors.InvalidInput, "Invalid Configuration")
	}
	if len(*server.ClusterName) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Missing Cluster Name")
	}
	if server.ClientPort == 0 {
		return errors.Wrapf(errors.InvalidInput, "Missing ClientPort")
	}
	if len(*server.Fqdn) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Missing Fqdn")
	}

	if server.Name == nil {
		server.Name = &name
	}
	return nil
}

// Delete
func (c *client) Delete(ctx context.Context, group, name, clusterName string) error {
	etcdserver, err := c.Get(ctx, group, name, clusterName)
	if err != nil {
		return err
	}
	if len(*etcdserver) == 0 {
		return fmt.Errorf("etcdserver [%s] not found", name)
	}

	request, err := getEtcdServerRequest(wssdcloudcommon.Operation_DELETE, name, clusterName, &(*etcdserver)[0])
	if err != nil {
		return err
	}
	_, err = c.EtcdServerAgentClient.Invoke(ctx, request)
	services.HandleGRPCError(err)

	return err
}

func getEtcdServersFromResponse(response *wssdcloudcloud.EtcdServerResponse, clusterName string) *[]etcdcluster.EtcdServer {
	etcdServers := []etcdcluster.EtcdServer{}
	for _, etcdservers := range response.GetEtcdServers() {
		etcdServers = append(etcdServers, *(getEtcdServer(etcdservers, clusterName)))
	}

	return &etcdServers
}

func getEtcdServerRequest(opType wssdcloudcommon.Operation, name, clusterName string, server *etcdcluster.EtcdServer) (*wssdcloudcloud.EtcdServerRequest, error) {
	request := &wssdcloudcloud.EtcdServerRequest{
		OperationType: opType,
		EtcdServers:   []*wssdcloudcloud.EtcdServer{},
	}
	if server != nil {
		etcdserver, err := getWssdEtcdServer(server, opType)
		if err != nil {
			return nil, err
		}
		request.EtcdServers = append(request.EtcdServers, etcdserver)
	} else if len(name) > 0 { // TODO for some operations we will need more attributes
		request.EtcdServers = append(request.EtcdServers,
			&wssdcloudcloud.EtcdServer{
				Name:            name,
				EtcdClusterName: clusterName,
			})
	} else { // TODO for some operations we will need more attributes
		request.EtcdServers = append(request.EtcdServers,
			&wssdcloudcloud.EtcdServer{
				EtcdClusterName: clusterName,
			})
	}

	return request, nil
}
