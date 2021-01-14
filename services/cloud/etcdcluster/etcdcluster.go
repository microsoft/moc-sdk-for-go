// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdcluster

import (
	"github.com/Azure/go-autorest/autorest"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

type EtcdServerProperties struct {
	// ClusterName
	ClusterName *string `json:"clustername"`
	// FQDN address that can be used to talk to the server
	Fqdn *string `json:"fqdn"`
	// Client port that ETCD listens on for client communication
	ClientPort uint32 `json:"clientPort"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Server defines the structure of a server
type EtcdServer struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The Azure Resource Manager resource ID for the key cluster.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the key cluster.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type of the key cluster.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The supported Azure location where the key cluster should be created.
	Location *string `json:"location,omitempty"`
	// Tags - The tags that will be assigned to the key cluster.
	Tags map[string]*string `json:"tags"`
	// Properties - Properties of the cluster
	*EtcdServerProperties `json:"properties,omitempty"`
}

func getEtcdCluster(cluster *wssdcloudcloud.EtcdCluster, group string) *cloud.EtcdCluster {
	return &cloud.EtcdCluster{
		ID:       &cluster.Id,
		Name:     &cluster.Name,
		Version:  &cluster.Status.Version.Number,
		Location: &cluster.LocationName,
		EtcdClusterProperties: &cloud.EtcdClusterProperties{
			CaCertificate: &cluster.CaCertificate,
			CaKey:         &cluster.CaKey,
			Statuses:      status.GetStatuses(cluster.GetStatus()),
		},
	}
}

func getWssdEtcdCluster(cluster *cloud.EtcdCluster, group string) (*wssdcloudcloud.EtcdCluster, error) {
	if cluster.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdCluster name is missing")
	}
	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	if cluster.CaCertificate == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdCluster CaCertificate is missing")
	}
	if cluster.CaKey == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdCluster CaKey is missing")
	}
	etcdcluster := &wssdcloudcloud.EtcdCluster{
		GroupName:     group,
		Name:          *cluster.Name,
		CaCertificate: *cluster.CaCertificate,
		CaKey:         *cluster.CaKey,
	}

	if cluster.Version != nil {
		if etcdcluster.Status == nil {
			etcdcluster.Status = status.InitStatus()
		}
		etcdcluster.Status.Version.Number = *cluster.Version
	}

	return etcdcluster, nil
}
