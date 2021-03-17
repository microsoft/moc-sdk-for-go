// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package security

import (
	"github.com/microsoft/moc/pkg/errors"
	pbcom "github.com/microsoft/moc/rpc/common"
)

func GetProviderType(pbProvider pbcom.ProviderType) ProviderType {
	providerInt := int32(pbProvider)
	value, found := pbcom.ProviderType_name[providerInt]
	if !found {
		return AnyProviderType // Not found, return AnyProvider
	}
	return ProviderType(value)
}

func GetWssdProviderType(providerType ProviderType) (pbcom.ProviderType, error) {
	// Convert sdk enum to string representation
	providerString := string(providerType)

	var pbProvider pbcom.ProviderType
	if len(providerType) == 0 {
		pbProvider = pbcom.ProviderType_AnyProvider
	} else {
		// Find the corresponding string in provider map
		value, found := pbcom.ProviderType_value[providerString]
		if !found {
			// Not found, user supplied unsupported provider
			return pbcom.ProviderType_AnyProvider, errors.Wrapf(errors.NotSupported, "Provider Type [%+v] is not currently supported", providerType)
		}
		pbProvider = pbcom.ProviderType(value)
	}
	return pbProvider, nil
}

type ProviderType string

var (
	AnyProviderType            ProviderType = ""
	VirtualMachineType         ProviderType = "VirtualMachine"
	VirtualMachineScaleSetType ProviderType = "VirtualMachineScaleSet"
	LoadBalancerType           ProviderType = "LoadBalancer"
	VirtualNetworkType         ProviderType = "VirtualNetwork"
	VirtualHardDiskType        ProviderType = "VirtualHardDisk"
	GalleryImageType           ProviderType = "GalleryImage"
	VirtualMachineImageType    ProviderType = "VirtualMachineImage"
	NetworkInterfaceType       ProviderType = "NetworkInterface"
	KeyVaultType               ProviderType = "KeyVault"
	KubernetesType             ProviderType = "Kubernetes"
	ClusterType                ProviderType = "Cluster"
	ControlPlaneType           ProviderType = "ControlPlane"
	GroupType                  ProviderType = "Group"
	NodeType                   ProviderType = "Node"
	LocationType               ProviderType = "Location"
	StorageContainerType       ProviderType = "StorageContainer"
	SubscriptionType           ProviderType = "Subscription"
	VipPoolType                ProviderType = "VipPool"
	MacPoolType                ProviderType = "MacPool"
	EtcdClusterType            ProviderType = "EtcdCluster"
	BareMetalMachineType       ProviderType = "BareMetalMachine"
	RoleType                   ProviderType = "Role"
	RoleAssignmentType         ProviderType = "RoleAssignment"
)
