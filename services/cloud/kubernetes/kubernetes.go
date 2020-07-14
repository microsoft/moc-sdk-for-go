// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package kubernetes

import (
	"github.com/microsoft/moc-proto/pkg/errors"
	"github.com/microsoft/moc-proto/pkg/status"
	wssdcloud "github.com/microsoft/moc-proto/rpc/cloudagent/cloud"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

// Conversion functions from cloud to wssdcloud
func (c *client) getWssdKubernetes(gp *cloud.Kubernetes, group string) (*wssdcloud.Kubernetes, error) {
	if gp.Network == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration")
	}
	wssdNetwork, err := c.getWssdKubernetesNetwork(gp.Network)
	if err != nil {
		return nil, err
	}
	if gp.Storage == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Storage Configuration")
	}
	wssdStorage, err := c.getWssdKubernetesStorage(gp.Storage)
	if err != nil {
		return nil, err
	}
	if gp.Cluster == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Cluster Configuration")
	}
	wssdCluster, err := c.getWssdKubernetesCluster(gp.Cluster)
	if err != nil {
		return nil, err
	}
	if gp.Compute == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute Configuration")
	}
	wssdCompute, err := c.getWssdKubernetesCompute(gp.Compute)
	if err != nil {
		return nil, err
	}

	if gp.ClusterAPI == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI Configuration")
	}
	wssdCapi, err := c.getWssdKubernetesCapi(gp.ClusterAPI)
	if err != nil {
		return nil, err
	}

	managementStrategyType := wssdcloud.ManagementStrategyType_Distinct
	switch gp.ManagementStrategy {
	case cloud.Pivoted:
		managementStrategyType = wssdcloud.ManagementStrategyType_Pivoted
	}

	kubernetes := &wssdcloud.Kubernetes{
		Name:               *gp.Name,
		GroupName:          group,
		Network:            wssdNetwork,
		Storage:            wssdStorage,
		Compute:            wssdCompute,
		Cluster:            wssdCluster,
		CapiConfig:         wssdCapi,
		ManagementStrategy: managementStrategyType,
		DeploymentManifest: gp.DeploymentManifest,
	}

	if gp.Version != nil {
		if kubernetes.Status == nil {
			kubernetes.Status = status.InitStatus()
		}
		kubernetes.Status.Version.Number = *gp.Version
	}

	if gp.ContainerRegistry != nil {
		wssdContainerRegistry, err := c.getWssdKubernetesContainerRegistry(gp.ContainerRegistry)
		if err != nil {
			return nil, err
		}
		kubernetes.ContainerRegistry = wssdContainerRegistry
	}

	return kubernetes, nil
}

func (c *client) getWssdKubernetesContainerRegistry(cfg *cloud.ContainerRegistryConfiguration) (*wssdcloud.ContainerRegistry, error) {
	if cfg.Name == nil || len(*cfg.Name) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Container Registry Name")
	}
	if cfg.Username == nil || len(*cfg.Username) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Container Registry Username")
	}
	if cfg.Password == nil || len(*cfg.Password) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Container Registry Password")
	}

	return &wssdcloud.ContainerRegistry{
		Name:     *cfg.Name,
		Username: *cfg.Username,
		Password: *cfg.Password,
	}, nil
}

func (c *client) getWssdKubernetesCapi(cfg *cloud.ClusterAPIConfiguration) (*wssdcloud.ClusterAPIConfiguration, error) {
	if cfg.ConfigurationEndpoint == nil || len(*cfg.ConfigurationEndpoint) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI ConfigurationEndpoint")
	}
	if cfg.InfrastructureProviderVersion == nil || len(*cfg.InfrastructureProviderVersion) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI InfrastructureProviderVersion")
	}
	if cfg.BootstrapProviderVersion == nil || len(*cfg.BootstrapProviderVersion) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI BootstrapProviderVersion")
	}
	if cfg.ControlPlaneProviderVersion == nil || len(*cfg.ControlPlaneProviderVersion) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI ControlPlaneProviderVersion")
	}
	if cfg.CoreProviderVersion == nil || len(*cfg.CoreProviderVersion) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing ClusterAPI CoreProviderVersion")
	}

	return &wssdcloud.ClusterAPIConfiguration{
		ConfigurationEndpoint:         *cfg.ConfigurationEndpoint,
		InfrastructureProviderVersion: *cfg.InfrastructureProviderVersion,
		BootstrapProviderVersion:      *cfg.BootstrapProviderVersion,
		ControlPlaneProviderVersion:   *cfg.ControlPlaneProviderVersion,
		CoreProviderVersion:           *cfg.CoreProviderVersion,
	}, nil
}
func (c *client) getWssdKubernetesNetwork(cfg *cloud.NetworkConfiguration) (*wssdcloud.NetworkConfiguration, error) {
	if cfg.CNI == nil || len(*cfg.CNI) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration.CNI")
	}
	if cfg.PodCIDR == nil || len(*cfg.PodCIDR) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration.PodCIDR")
	}
	if cfg.ClusterCIDR == nil || len(*cfg.ClusterCIDR) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration.ClusterCIDR")
	}
	if cfg.ControlPlaneCIDR == nil || len(*cfg.ControlPlaneCIDR) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration.ControlPlaneCIDR")
	}
	if cfg.VirtualNetwork == nil || len(*cfg.VirtualNetwork) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Configuration.Virtualnetwork")
	}

	return &wssdcloud.NetworkConfiguration{
		Cni:              *cfg.CNI,
		PodCidr:          *cfg.PodCIDR,
		ClusterCidr:      *cfg.ClusterCIDR,
		ControlPlaneCidr: *cfg.ControlPlaneCIDR,
		Virtualnetwork:   *cfg.VirtualNetwork,
	}, nil
}

func (c *client) getWssdKubernetesStorage(cfg *cloud.StorageConfiguration) (*wssdcloud.StorageConfiguration, error) {
	if cfg.CSI == nil || len(*cfg.CSI) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Storage Configuration.CSI")
	}

	return &wssdcloud.StorageConfiguration{
		Csi: *cfg.CSI,
	}, nil
}
func (c *client) getWssdKubernetesCluster(cfg *cloud.ClusterConfiguration) (*wssdcloud.ClusterConfiguration, error) {
	if cfg.Version == nil || len(*cfg.Version) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Storage Configuration.Version")
	}

	return &wssdcloud.ClusterConfiguration{
		Version: *cfg.Version,
	}, nil

}
func (c *client) getWssdKubernetesCompute(cfg *cloud.ComputeConfiguration) (*wssdcloud.ComputeConfiguration, error) {
	if cfg.CRI == nil || len(*cfg.CRI) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute Configuration.CRI")
	}

	if cfg.SSH == nil || cfg.SSH.PublicKey == nil || *cfg.SSH.PublicKey.KeyData == "" {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute Configuration.SSH.PublicKey")
	}

	if cfg.NodePools == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute NodePools Configuration")
	}

	publicKey := &wssdcloud.SSHPublicKey{
		KeyData: *cfg.SSH.PublicKey.KeyData,
	}

	nodePools := []*wssdcloud.NodePoolConfiguration{}
	for _, nodePool := range *cfg.NodePools {
		wssdnodepool, err := c.getWssdKubernetesComputeNodePool(&nodePool)
		if err != nil {
			return nil, err
		}
		nodePools = append(nodePools, wssdnodepool)
	}

	return &wssdcloud.ComputeConfiguration{
		Cri:       *cfg.CRI,
		PublicKey: publicKey,
		NodePools: nodePools,
	}, nil
}

func (c *client) getWssdKubernetesComputeNodePool(cfg *cloud.NodePoolConfiguration) (*wssdcloud.NodePoolConfiguration, error) {
	if cfg.Replicas == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute NodePoolConfiguration.Replicas")
	}
	if cfg.ImageReference == nil || len(*cfg.ImageReference) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute NodePoolConfiguration.ImageReference")
	}
	if cfg.VMSize == nil || len(*cfg.VMSize) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Compute NodePoolConfiguration.VMSize")
	}

	wssdcloudnodeType := wssdcloud.NodeType_ControlPlane

	switch cfg.NodeType {
	case cloud.ControlPlane:
	case cloud.LinuxWorker:
		wssdcloudnodeType = wssdcloud.NodeType_LinuxWorker
	case cloud.WindowsWorker:
		wssdcloudnodeType = wssdcloud.NodeType_WindowsWorker
	case cloud.LoadBalancer:
		wssdcloudnodeType = wssdcloud.NodeType_LoadBalancer
	default:
	}

	return &wssdcloud.NodePoolConfiguration{
		Replicas:       *cfg.Replicas,
		Imagereference: *cfg.ImageReference,
		NodeType:       wssdcloudnodeType,
		VMSize:         *cfg.VMSize,
	}, nil

}

// Conversion functions from wssdcloud to cloud
func (c *client) getKubernetes(gp *wssdcloud.Kubernetes) *cloud.Kubernetes {
	sshConfig := c.getKubernetesComputeSSHConfiguration(gp.Compute.PublicKey)

	nodepools := []cloud.NodePoolConfiguration{}
	for _, wssdnodepool := range gp.Compute.NodePools {
		nodepools = append(nodepools, *(c.getKubernetesComputeNodePool(wssdnodepool)))
	}

	containerRegistry := &cloud.ContainerRegistryConfiguration{}
	if gp.ContainerRegistry != nil {
		containerRegistry = &cloud.ContainerRegistryConfiguration{
			Name:     &gp.ContainerRegistry.Name,
			Username: &gp.ContainerRegistry.Username,
			Password: &gp.ContainerRegistry.Password,
		}
	}

	return &cloud.Kubernetes{
		Name:    &gp.Name,
		Version: &gp.Status.Version.Number,
		KubernetesProperties: &cloud.KubernetesProperties{
			Statuses: status.GetStatuses(gp.GetStatus()),
			Network: &cloud.NetworkConfiguration{
				CNI:              &gp.Network.Cni,
				PodCIDR:          &gp.Network.PodCidr,
				ClusterCIDR:      &gp.Network.ClusterCidr,
				ControlPlaneCIDR: &gp.Network.ControlPlaneCidr,
				VirtualNetwork:   &gp.Network.Virtualnetwork,
				LoadBalancerVip:  &gp.Network.LoadBalancerVip,
				LoadBalancerMac:  &gp.Network.LoadBalancerMac,
			},
			Storage: &cloud.StorageConfiguration{
				CSI: &gp.Storage.Csi,
			},
			Compute: &cloud.ComputeConfiguration{
				CRI:       &gp.Compute.Cri,
				SSH:       sshConfig,
				NodePools: &nodepools,
			},
			Cluster: &cloud.ClusterConfiguration{
				Version: &gp.Cluster.Version,
			},
			ClusterAPI: &cloud.ClusterAPIConfiguration{
				ConfigurationEndpoint:         &gp.CapiConfig.ConfigurationEndpoint,
				InfrastructureProviderVersion: &gp.CapiConfig.InfrastructureProviderVersion,
				BootstrapProviderVersion:      &gp.CapiConfig.BootstrapProviderVersion,
				ControlPlaneProviderVersion:   &gp.CapiConfig.ControlPlaneProviderVersion,
				CoreProviderVersion:           &gp.CapiConfig.CoreProviderVersion,
			},
			ContainerRegistry: containerRegistry,
			KubeConfig:        gp.KubeConfig,
		},
	}
}

func (c *client) getKubernetesComputeSSHConfiguration(pk *wssdcloud.SSHPublicKey) *cloud.SSHConfiguration {
	publicKey := cloud.SSHPublicKey{
		KeyData: &pk.KeyData,
	}

	return &cloud.SSHConfiguration{
		PublicKey: &publicKey,
	}
}

func (c *client) getKubernetesComputeNodePool(gp *wssdcloud.NodePoolConfiguration) *cloud.NodePoolConfiguration {
	nodeType := cloud.ControlPlane
	switch gp.NodeType {
	case wssdcloud.NodeType_ControlPlane:
	case wssdcloud.NodeType_LinuxWorker:
		nodeType = cloud.LinuxWorker
	case wssdcloud.NodeType_WindowsWorker:
		nodeType = cloud.WindowsWorker
	case wssdcloud.NodeType_LoadBalancer:
		nodeType = cloud.LoadBalancer
	}
	return &cloud.NodePoolConfiguration{
		NodeType:       nodeType,
		Replicas:       &gp.Replicas,
		ImageReference: &gp.Imagereference,
		VMSize:         &gp.VMSize,
	}

}
