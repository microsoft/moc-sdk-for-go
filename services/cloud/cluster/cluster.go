// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package cluster

import (
	"github.com/microsoft/moc-proto/pkg/errors"
	"github.com/microsoft/moc-proto/pkg/status"
	wssdcloud "github.com/microsoft/moc-proto/rpc/cloudagent/cloud"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdCluster(gp *cloud.Cluster, location string) (*wssdcloud.Cluster, error) {
	if gp.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	if gp.FQDN == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing FQDN in Configuration")
	}

	cluster := &wssdcloud.Cluster{
		Name:         *gp.Name,
		Fqdn:         *gp.FQDN,
		LocationName: location,
	}

	if gp.Version != nil {
		if cluster.Status == nil {
			cluster.Status = status.InitStatus()
		}
		cluster.Status.Version.Number = *gp.Version
	}

	return cluster, nil
}

// Conversion functions from wssdcloud to cloud
func getCluster(gp *wssdcloud.Cluster) *cloud.Cluster {
	nodes := []cloud.Node{}

	for _, pbNode := range gp.Nodes {
		node := cloud.Node{}
		node.Name = &pbNode.Name
		node.FQDN = &pbNode.Fqdn

		nodes = append(nodes, node)
	}

	return &cloud.Cluster{
		Name: &gp.Name,
		ClusterProperties: &cloud.ClusterProperties{
			FQDN:     &gp.Fqdn,
			Statuses: status.GetStatuses(gp.GetStatus()),
		},
		Nodes:    &nodes,
		Location: &gp.LocationName,
		Version:  &gp.Status.Version.Number,
	}
}
