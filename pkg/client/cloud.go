// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	cloud_pb "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// GetLocationClient returns the virtual machine client to comminicate with the wssd agent
func GetLocationClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.LocationAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get LocationClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewLocationAgentClient(conn), nil
}

// GetGroupClient returns the virtual machine client to comminicate with the wssd agent
func GetGroupClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.GroupAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get GroupClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewGroupAgentClient(conn), nil
}

// GetNodeClient returns the virtual machine client to comminicate with the wssd agent
func GetNodeClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.NodeAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get NodeClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewNodeAgentClient(conn), nil
}

// GetKubernetesClient returns the virtual machine client to comminicate with the wssd agent
func GetKubernetesClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.KubernetesAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get KubernetesClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewKubernetesAgentClient(conn), nil
}

// GetClusterClient returns the cluster client to communicate with the wssd agent
func GetClusterClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.ClusterAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get ClusterClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewClusterAgentClient(conn), nil
}

// GetControlPlaneClient returns the cluster client to communicate with the wssd agent
func GetControlPlaneClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.ControlPlaneAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get ControlPlaneClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewControlPlaneAgentClient(conn), nil
}

// GetEtcdClusterClient returns the cluster client to communicate with the wssd agent
func GetEtcdClusterClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.EtcdClusterAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get EtcdClusterClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewEtcdClusterAgentClient(conn), nil
}

// GetEtcdServerClient returns the server client to communicate with the wssd agent
func GetEtcdServerClient(serverAddress *string, authorizer auth.Authorizer) (cloud_pb.EtcdServerAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get EtcdServerClient. Failed to dial: %v", err)
	}

	return cloud_pb.NewEtcdServerAgentClient(conn), nil
}
