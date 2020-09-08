// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	network_pb "github.com/microsoft/moc/rpc/cloudagent/network"
)

// GetVirtualNetworkClient returns the virtual network client to communicate with the wssdagent
func GetVirtualNetworkClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.VirtualNetworkAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualNetworkClient. Failed to dial: %v", err)
	}

	return network_pb.NewVirtualNetworkAgentClient(conn), nil
}

// GetNetworkInterfaceClient returns the virtual network interface client to communicate with the wssd agent
func GetNetworkInterfaceClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.NetworkInterfaceAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get NetworkInterfaceClient. Failed to dial: %v", err)
	}

	return network_pb.NewNetworkInterfaceAgentClient(conn), nil
}

// GetLoadBalancerClient returns the loadbalancer client to communicate with the wssd agent
func GetLoadBalancerClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.LoadBalancerAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get LoadBalancerClient. Failed to dial: %v", err)
	}

	return network_pb.NewLoadBalancerAgentClient(conn), nil
}

// GetVipPoolClient returns the vippool client to communicate with the wssd agent
func GetVipPoolClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.VipPoolAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VipPoolClient. Failed to dial: %v", err)
	}

	return network_pb.NewVipPoolAgentClient(conn), nil
}

// GetMacPoolClient returns the macpool client to communicate with the wssd agent
func GetMacPoolClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.MacPoolAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get MacPoolClient. Failed to dial: %v", err)
	}

	return network_pb.NewMacPoolAgentClient(conn), nil
}
