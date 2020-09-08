// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package cluster

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloud.ClusterAgentClient
}

// newClusterClient - creates a client session with the backend wssd agent
func newClusterClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetClusterClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]cloud.Cluster, error) {
	request, err := c.getClusterRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.ClusterAgentClient.GetCluster(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getClusterFromResponse(response), nil

}

// Conversion functions from wssdcloud to cloud
func getNode(gp *wssdcloud.Node) *cloud.Node {
	return &cloud.Node{
		Name:     &gp.Name,
		Location: &gp.LocationName,
		NodeProperties: &cloud.NodeProperties{
			FQDN: &gp.Fqdn,
		},
	}
}

// GetNodes
func (c *client) GetNodes(ctx context.Context, location, clusterName string) (*[]cloud.Node, error) {
	request, err := c.getClusterRequest(wssdcloudcommon.Operation_GET, location, clusterName, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.ClusterAgentClient.GetCluster(ctx, request)
	if err != nil {
		return nil, err
	}

	Nodes := []cloud.Node{}
	for _, cluster := range response.GetClusters() {
		pbNodeResponse, err := c.ClusterAgentClient.GetNodes(ctx, cluster)
		if err != nil {
			return nil, err
		}
		if pbNodeResponse.Nodes == nil || len(pbNodeResponse.Nodes) == 0 {
			return nil, fmt.Errorf("The cluster doesnt have any nodes")
		}
		for _, pbNode := range pbNodeResponse.Nodes {
			Nodes = append(Nodes, *getNode(pbNode))
		}
	}

	return &Nodes, nil
}

// Load
func (c *client) Load(ctx context.Context, location, name string, sg *cloud.Cluster) (cluster *cloud.Cluster, err error) {
	err = c.validate(ctx, sg, location)
	if err != nil {
		return
	}
	cluster = nil

	request, err := c.getClusterRequest(wssdcloudcommon.Operation_POST, location, name, sg)
	if err != nil {
		return
	}
	response, err := c.ClusterAgentClient.LoadCluster(ctx, request)
	if err != nil {
		return
	}
	gps := c.getClusterFromResponse(response)
	if len(*gps) == 0 {
		return
	}

	cluster = &(*gps)[0]
	return
}

// Unload methods invokes create or update on the client
func (c *client) Unload(ctx context.Context, location, name string) error {
	gp, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*gp) == 0 {
		return fmt.Errorf("Cluster [%s] not found", name)
	}

	request, err := c.getClusterRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*gp)[0])
	if err != nil {
		return err
	}
	_, err = c.ClusterAgentClient.UnloadCluster(ctx, request)

	return err
}

///////////////////////////
// Private Methods
func (c *client) validate(ctx context.Context, sg *cloud.Cluster, location string) (err error) {
	if sg == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input is nil")
		return
	}
	if sg.Location == nil && len(location) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "Location is nil")
		return
	} else if sg.Location == nil {
		sg.Location = &location
	}

	if sg.ClusterProperties == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing ClusterProperties")
		return
	}
	if sg.ClusterProperties.FQDN == nil {
		err = errors.Wrapf(errors.InvalidInput, "Missing ClusterProperties.FQDN")
		return
	}
	return

}
func (c *client) getClusterFromResponse(response *wssdcloud.ClusterResponse) *[]cloud.Cluster {
	gps := []cloud.Cluster{}
	for _, gp := range response.GetClusters() {
		gps = append(gps, *(getCluster(gp)))
	}

	return &gps
}

func (c *client) getClusterRequest(opType wssdcloudcommon.Operation, location, name string, gpss *cloud.Cluster) (*wssdcloud.Cluster, error) {
	wssdcluster := &wssdcloud.Cluster{
		Name:         name,
		LocationName: location,
	}
	var err error
	if gpss != nil {
		wssdcluster, err = getWssdCluster(gpss, location)
		if err != nil {
			return nil, err
		}
	}
	return wssdcluster, nil
}
