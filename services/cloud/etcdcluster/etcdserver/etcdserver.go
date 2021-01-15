// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdserver

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud/etcdcluster"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getEtcdServer(sec *wssdcloudcloud.EtcdServer, clusterName string) *etcdcluster.EtcdServer {
	return &etcdcluster.EtcdServer{
		ID:      &sec.Id,
		Name:    &sec.Name,
		Version: &sec.Status.Version.Number,
		EtcdServerProperties: &etcdcluster.EtcdServerProperties{
			ClusterName: &clusterName,
			Statuses:    status.GetStatuses(sec.GetStatus()),
			Fqdn:        &sec.Fqdn,
			ClientPort:  sec.ClientPort,
		},
	}
}

func getWssdEtcdServer(server *etcdcluster.EtcdServer, opType wssdcloudcommon.Operation) (*wssdcloudcloud.EtcdServer, error) {
	if server.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdServer name is missing")
	}
	if server.ClusterName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdCluster name is missing")
	}
	if server.Fqdn == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdServer Fqdn is missing")
	}
	if server.ClientPort == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "EtcdServer ClientPort is missing")
	}
	etcdserver := &wssdcloudcloud.EtcdServer{
		Name:            *server.Name,
		EtcdClusterName: *server.ClusterName,
		Fqdn:            *server.Fqdn,
		ClientPort:      server.ClientPort,
	}

	if server.Version != nil {
		if etcdserver.Status == nil {
			etcdserver.Status = status.InitStatus()
		}
		etcdserver.Status.Version.Number = *server.Version
	}

	return etcdserver, nil
}
