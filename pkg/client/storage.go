// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc-proto/pkg/auth"
	storage_pb "github.com/microsoft/moc-proto/rpc/cloudagent/storage"
)

// GetVirtualHardDiskClient returns the virtual network client to communicate with the wssdagent
func GetVirtualHardDiskClient(serverAddress *string, authorizer auth.Authorizer) (storage_pb.VirtualHardDiskAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualHardDiskClient. Failed to dial: %v", err)
	}

	return storage_pb.NewVirtualHardDiskAgentClient(conn), nil
}

// GetContainerClient returns the virtual network client to communicate with the wssdagent
func GetStorageContainerClient(serverAddress *string, authorizer auth.Authorizer) (storage_pb.ContainerAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get ContainerClient. Failed to dial: %v", err)
	}

	return storage_pb.NewContainerAgentClient(conn), nil
}
