// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	compute_pb "github.com/microsoft/moc/rpc/cloudagent/compute"
)

// GetGalleryImageClient returns the virtual machine client to communicate with the wssd agent
func GetGalleryImageClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.GalleryImageAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get GalleryImageClient. Failed to dial: %v", err)
	}

	return compute_pb.NewGalleryImageAgentClient(conn), nil
}

// GetVirtualMachineClient returns the virtual machine client to communicate with the wssd agent
func GetVirtualMachineClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineAgentClient(conn), nil
}

// GetAvailabilitySet returns the virtual machine client to communicate with the wssd agent
func GetAvailabilitySetClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.AvailabilitySetAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get AvailabilitySetClient. Failed to dial: %v", err)
	}

	return compute_pb.NewAvailabilitySetAgentClient(conn), nil
}

// GetPlacementGroup returns the virtual machine client to communicate with the wssd agent
func GetPlacementGroupClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.PlacementGroupAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get PlacementGroupClient. Failed to dial: %v", err)
	}

	return compute_pb.NewPlacementGroupAgentClient(conn), nil
}

// GetVirtualMachineScaleSetClient returns the virtual machine client to communicate with the wssd agent
func GetVirtualMachineScaleSetClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineScaleSetAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineScaleSetClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineScaleSetAgentClient(conn), nil
}

// GetBareMetalHostClient returns the bare metal machine client to communicate with the wssd agent
func GetBareMetalHostClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.BareMetalHostAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get BareMetalHostClient. Failed to dial: %v", err)
	}

	return compute_pb.NewBareMetalHostAgentClient(conn), nil
}

// GetBareMetalMachineClient returns the bare metal machine client to communicate with the wssd agent
func GetBareMetalMachineClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.BareMetalMachineAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get BareMetalMachineClient. Failed to dial: %v", err)
	}

	return compute_pb.NewBareMetalMachineAgentClient(conn), nil
}
