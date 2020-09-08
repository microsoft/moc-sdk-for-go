// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package networkinterface

import (
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
)

// Conversion functions from network interface to wssdcloud network interface
func getWssdNetworkInterface(c *network.Interface, group string) (*wssdcloudnetwork.NetworkInterface, error) {
	if c == nil || c.InterfacePropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Interface Properties")
	}
	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if c.IPConfigurations == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing IPConfigurations")
	}

	wssdipconfigs := []*wssdcloudnetwork.IpConfiguration{}
	for _, ipconfig := range *c.IPConfigurations {
		wssdipconfig, err := getWssdNetworkInterfaceIPConfig(&ipconfig)
		if err != nil {
			return nil, err
		}
		wssdipconfigs = append(wssdipconfigs, wssdipconfig)
	}

	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for Network Interface")
	}

	vnic := &wssdcloudnetwork.NetworkInterface{
		Name:             *c.Name,
		IpConfigurations: wssdipconfigs,
		GroupName:        group,
	}

	if c.Version != nil {
		if vnic.Status == nil {
			vnic.Status = status.InitStatus()
		}
		vnic.Status.Version.Number = *c.Version
	}

	if c.MacAddress != nil {
		vnic.Macaddress = *c.MacAddress
	}
	return vnic, nil
}

func ipAllocationMethodProtobufToSdk(wssdIpconfig *wssdcloudnetwork.IpConfiguration, ipConfig *network.InterfaceIPConfiguration) {
	if wssdIpconfig.Allocation == wssdcommonproto.IPAllocationMethod_Invalid {
		return
	}
	var val network.IPAllocationMethod
	switch wssdIpconfig.Allocation {
	case wssdcommonproto.IPAllocationMethod_Static:
		val = network.Static
	case wssdcommonproto.IPAllocationMethod_Dynamic:
		val = network.Dynamic
	}
	ipConfig.PrivateIPAllocationMethod = &val
}

func ipAllocationMethodSdkToProtobuf(ipConfig *network.InterfaceIPConfiguration, wssdIpConfig *wssdcloudnetwork.IpConfiguration) {
	if ipConfig.PrivateIPAllocationMethod == nil {
		return
	}
	switch *ipConfig.PrivateIPAllocationMethod {
	case network.Static:
		wssdIpConfig.Allocation = wssdcommonproto.IPAllocationMethod_Static
	case network.Dynamic:
		wssdIpConfig.Allocation = wssdcommonproto.IPAllocationMethod_Dynamic
	}
}

func getWssdNetworkInterfaceIPConfig(ipConfig *network.InterfaceIPConfiguration) (*wssdcloudnetwork.IpConfiguration, error) {
	if ipConfig.InterfaceIPConfigurationPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Interface IPConfiguration Properties")
	}

	if ipConfig.Subnet == nil ||
		ipConfig.Subnet.ID == nil ||
		len(*ipConfig.Subnet.ID) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Subnet Reference")
	}

	wssdipconfig := &wssdcloudnetwork.IpConfiguration{
		Subnetid: *ipConfig.Subnet.ID,
	}
	if ipConfig.PrivateIPAddress != nil {
		wssdipconfig.Ipaddress = *ipConfig.PrivateIPAddress
	}
	if ipConfig.PrefixLength != nil {
		wssdipconfig.Prefixlength = *ipConfig.PrefixLength
	}
	if ipConfig.Gateway != nil {
		wssdipconfig.Gateway = *ipConfig.Gateway
	}
	ipAllocationMethodSdkToProtobuf(ipConfig, wssdipconfig)

	if ipConfig.LoadBalancerBackendAddressPools != nil {
		for _, addresspool := range *ipConfig.LoadBalancerBackendAddressPools {
			wssdipconfig.Loadbalanceraddresspool = append(wssdipconfig.Loadbalanceraddresspool, *addresspool.Name)
		}
	}
	return wssdipconfig, nil
}

// Conversion function from wssdcloud network interface to network interface
func getNetworkInterface(server, group string, c *wssdcloudnetwork.NetworkInterface) (*network.Interface, error) {
	ipConfigs := []network.InterfaceIPConfiguration{}
	for _, wssdipconfig := range c.IpConfigurations {
		ipConfigs = append(ipConfigs, *(getNetworkIpConfig(wssdipconfig)))
	}

	vnetIntf := &network.Interface{
		Name:    &c.Name,
		ID:      &c.Id,
		Version: &c.Status.Version.Number,
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			MacAddress: &c.Macaddress,
			// TODO: Type
			IPConfigurations: &ipConfigs,
			Statuses:         status.GetStatuses(c.GetStatus()),
		},
	}

	return vnetIntf, nil
}

func getNetworkIpConfig(wssdcloudipconfig *wssdcloudnetwork.IpConfiguration) *network.InterfaceIPConfiguration {
	ipconfig := &network.InterfaceIPConfiguration{
		InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAddress: &wssdcloudipconfig.Ipaddress,
			Subnet:           &network.APIEntityReference{ID: &wssdcloudipconfig.Subnetid},
			Gateway:          &wssdcloudipconfig.Gateway,
			PrefixLength:     &wssdcloudipconfig.Prefixlength,
		},
	}

	ipAllocationMethodProtobufToSdk(wssdcloudipconfig, ipconfig)

	var addresspools []network.BackendAddressPool
	for _, addresspool := range wssdcloudipconfig.Loadbalanceraddresspool {
		bap := network.BackendAddressPool{
			Name: &addresspool,
		}
		addresspools = append(addresspools, bap)
	}
	ipconfig.LoadBalancerBackendAddressPools = &addresspools
	return ipconfig
}
