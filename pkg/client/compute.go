// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	compute_pb "github.com/microsoft/moc/rpc/cloudagent/compute"
)

// GetGalleryImageClient returns the virtual machine client to comminicate with the wssd agent
func GetGalleryImageClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.GalleryImageAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get GalleryImageClient. Failed to dial: %v", err)
	}

	return compute_pb.NewGalleryImageAgentClient(conn), nil
}

// GetVirtualMachineClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineAgentClient(conn), nil
}

// GetVirtualMachineScaleSetClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineScaleSetClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineScaleSetAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineScaleSetClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineScaleSetAgentClient(conn), nil
}
