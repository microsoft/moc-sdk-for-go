// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package cloud

import (
	"github.com/Azure/go-autorest/autorest"
)

// LocationProperties the resource group properties.
type LocationProperties struct {
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Location resource group information.
type Location struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	// Properties
	*LocationProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

// GroupProperties the resource group properties.
type GroupProperties struct {
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Group resource group information.
type Group struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	// Properties
	*GroupProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// ManagedBy - The ID of the resource that manages this resource group.
	ManagedBy *string `json:"managedBy,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

// NodeProperties the resource group properties.
type NodeProperties struct {
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// FQDN
	FQDN *string `json:"fqdn,omitempty"`

	Port *int32 `json:"port,omitempty"`

	AuthorizerPort *int32 `json:"authorizerPort,omitempty"`

	Certificate *string `json:"certificate,omitempty"`
}

// Node resource group information.
type Node struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	//Properties
	*NodeProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

type ClusterConfiguration struct {
	// Version
	Version *string `json:"version,omitempty"`
}

type ManagementStrategyType string

const (
	Pivoted  ManagementStrategyType = "Pivoted"
	Distinct ManagementStrategyType = "Distinct"
)

type NetworkConfiguration struct {
	// CNI
	CNI *string `json:"cni,omitempty"`
	// PodCIDR
	PodCIDR *string `json:"pidCidr,omitempty"`
	// ClusterCIDR
	ClusterCIDR *string `json:"clusterCidr,omitempty"`
	// ControlPlaneCIDR
	ControlPlaneCIDR *string `json:"controlPlaneCidr,omitempty"`
	// VirtualNetwork
	VirtualNetwork *string `json:"virtualNetwork,omitempty"`
	// LoadBalancerVip
	LoadBalancerVip *string `json:"loadBalancerVip,omitempty"`
	// LoadBalancerMac
	LoadBalancerMac *string `json:"loadBalancerMac,omitempty"`
}

type NodeType string

const (
	ControlPlane  NodeType = "ControlPlane"
	LinuxWorker   NodeType = "LinuxWorker"
	WindowsWorker NodeType = "WindowsWorker"
	LoadBalancer  NodeType = "LoadBalancer"
)

type NodePoolConfiguration struct {
	// NodeType
	NodeType NodeType `json:"nodeType,omitempty"`
	// Replicas
	Replicas *int32 `json:"replicas,omitempty"`
	// ImageReference
	ImageReference *string `json:"imageReference,omitempty"`
	// VMSize
	VMSize *string `json:"vmSize,omitempty"`
}

type SSHPublicKey struct {
	// KeyData - SSH public key certificate used to authenticate with the VM through ssh. The key needs to be at least 2048-bit and in ssh-rsa format.
	KeyData *string `json:"keyData,omitempty"`
}

type SSHConfiguration struct {
	// PublicKeys - The SSH public key used to authenticate with linux based VMs.
	PublicKey *SSHPublicKey `json:"publicKey,omitempty"`
}

type ComputeConfiguration struct {
	// CRI
	CRI *string `json:"cri,omitempty"`
	// SSH
	SSH *SSHConfiguration `json:"ssh,omitempty"`
	// NodePools
	NodePools *[]NodePoolConfiguration `json:"nodePools,omitempty"`
}

type StorageConfiguration struct {
	// Version
	CSI *string `json:"csi,omitempty"`
}

// ClusterAPIConfiguration is the configuration needed for setting up Cluster API
type ClusterAPIConfiguration struct {
	// ConfigurationEndpoint
	ConfigurationEndpoint *string `json:"configurationEndpoint,omitempty"`
	// InfrastructureProviderVersion
	InfrastructureProviderVersion *string `json:"infrastructureProviderVersion,omitempty"`
	// BootstrapProviderVersion
	BootstrapProviderVersion *string `json:"bootstrapProviderVersion,omitempty"`
	// ControlPlaneProviderVersion
	ControlPlaneProviderVersion *string `json:"controlPlaneProviderVersion,omitempty"`
	// CoreProviderVersion
	CoreProviderVersion *string `json:"coreProviderVersion,omitempty"`
}

// ContainerRegistryConfiguration is the configuration needed for a container registry
type ContainerRegistryConfiguration struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Username
	Username *string `json:"username,omitempty"`
	// Password
	Password *string `json:"password,omitempty"`
}

// KubernetesProperties the resource group properties.
type KubernetesProperties struct {
	// Cluster
	Cluster *ClusterConfiguration `json:"cluster,omitempty"`
	// Network
	Network *NetworkConfiguration `json:"network,omitempty"`
	// Storage
	Storage *StorageConfiguration `json:"storage,omitempty"`
	// Compute
	Compute *ComputeConfiguration `json:"compute,omitempty"`
	// ClusterAPI
	ClusterAPI *ClusterAPIConfiguration `json:"clusterapi,omitempty"`
	// ContainerRegistry
	ContainerRegistry *ContainerRegistryConfiguration `json:"containerregistry,omitempty"`
	// ManagementStrategy
	ManagementStrategy ManagementStrategyType `json:"managementstrategy,omitempty"`
	// KubeConfig
	KubeConfig []byte `json:"kubeconfig,omitempty"`
	// DeploymentManifest
	DeploymentManifest []byte `json:"deploymentManifest,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Kubernetes resource group information.
type Kubernetes struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	//Properties
	*KubernetesProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

// ClusterProperties the resource group properties.
type ClusterProperties struct {
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// FQDN
	FQDN *string `json:"fqdn,omitempty"`
}

// Cluster resource group information.
type Cluster struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	//Properties
	*ClusterProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
	// Nodes
	Nodes *[]Node `json:"nodes,omitempty"`
}

// ControlPlaneProperties the resource group properties.
type ControlPlaneProperties struct {
	// Statuses - provides state of the ControlPlane like denoting whether
	// each ControlPlane is the Leader or Active
	Statuses map[string]*string `json:"statuses"`
	// FQDN - provides the ControlPlane FQDN (or IP) used for the leadership
	// election.
	FQDN *string `json:"fqdn,omitempty"`
	// Port - provides the ControlPlane Port (or IP) used for the leadership
	// election.
	Port *int32 `json:"port,omitempty"`
}

// ControlPlane resource group information.
type ControlPlaneInfo struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	// Properties
	*ControlPlaneProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

// EtcdClusterProperties the resource group properties.
type EtcdClusterProperties struct {
	// CaCertificate used as root certificate for communication among ETCD nodes
	// and to the ETCD cluster
	CaCertificate *string `json:"cacertificate,omitempty"`
	// CaKey is the private key corresponding to the CaCertificate
	CaKey *string `json:"cakey,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// EtcdCluster resource group information.
type EtcdCluster struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the resource group.
	ID *string `json:"id,omitempty"`
	// Name - The name of the resource group.
	Name *string `json:"name,omitempty"`
	// Properties
	*EtcdClusterProperties `json:"properties,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The location of the resource group. It cannot be changed after the resource group has been created. It must be one of the supported Azure locations.
	Location *string `json:"location,omitempty"`
	// Tags - The tags attached to the resource group.
	Tags map[string]*string `json:"tags"`
}

// AvailabilityZone describes the availabilityZone setting for a virtual machine
type AvailabilityZone struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`

	*AvailabilityZoneProperties `json:"availabilityzoneproperties,omitempty"`
}

type AvailabilityZoneProperties struct {
	// Statuses - Statuses
	Statuses map[string]*string `json:"statuses"`
	// Nodes
	Nodes *[]string `json:"nodes,omitempty"`
}
