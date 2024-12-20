// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package networkinterface

import (
	"strings"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
)

const (
	LNET_PREFIX        = "logicalnetworks"
	LNET_PREFIX_LEGACY = "logicalnetwork"
	VNET_PREFIX        = "virtualnetworks"
	SUBNET_PREFIX      = "subnets"
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
		wssdipconfig, err := getWssdNetworkInterfaceIPConfig(&ipconfig, c.Location)
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
		Dns:              getDns(c.DNSSettings),
		Tags:             tags.MapToProto(c.Tags),
		AdvancedPolicies: network.GetWssdAdvancedNetworkPolicies(c.AdvancedNetworkPolicies),
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

	if c.EnableAcceleratedNetworking != nil {
		if *c.EnableAcceleratedNetworking {
			vnic.IovWeight = uint32(100)
		} else {
			vnic.IovWeight = uint32(0)
		}
	}

	return vnic, nil
}

func getWssdDNSSettings(dnssetting *wssdcommonproto.Dns) *network.InterfaceDNSSettings {
	if dnssetting == nil {
		return nil
	}
	return &network.InterfaceDNSSettings{
		DNSServers:               &dnssetting.Servers,
		InternalDomainNameSuffix: &dnssetting.Domain,
	}
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

func getWssdNetworkInterfaceIPConfig(ipConfig *network.InterfaceIPConfiguration, location *string) (*wssdcloudnetwork.IpConfiguration, error) {
	if ipConfig.InterfaceIPConfigurationPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Interface IPConfiguration Properties")
	}

	if ipConfig.Subnet == nil ||
		ipConfig.Subnet.ID == nil ||
		len(*ipConfig.Subnet.ID) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Subnet Reference")
	}

	wssdipconfig := &wssdcloudnetwork.IpConfiguration{}

	subnet, err := getWssdSubnetReference(ipConfig.Subnet.ID)

	if err != nil {
		return nil, err
	}

	if subnet != nil {
		wssdipconfig.Subnet = subnet
	} else {
		// for backward compat
		wssdipconfig.Subnetid = *ipConfig.Subnet.ID
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
	if ipConfig.Primary != nil {
		wssdipconfig.Primary = *ipConfig.Primary
	}
	if ipConfig.NetworkSecurityGroup != nil && ipConfig.NetworkSecurityGroup.ID != nil {
		wssdipconfig.NetworkSecurityGroupRef = &wssdcommonproto.NetworkSecurityGroupReference{
			ResourceRef: &wssdcommonproto.ResourceReference{
				Name: *ipConfig.NetworkSecurityGroup.ID,
			},
		}
	}
	if ipConfig.PublicIPAddress != nil && ipConfig.PublicIPAddress.ID != nil {
		wssdipconfig.PublicIPAddressRef = &wssdcommonproto.PublicIPAddressReference{
			ResourceRef: &wssdcommonproto.ResourceReference{
				Name: *ipConfig.PublicIPAddress.ID,
			},
		}
	}

	ipAllocationMethodSdkToProtobuf(ipConfig, wssdipconfig)

	if ipConfig.LoadBalancerBackendAddressPools != nil {
		for _, addresspool := range *ipConfig.LoadBalancerBackendAddressPools {
			wssdipconfig.Loadbalanceraddresspool = append(wssdipconfig.Loadbalanceraddresspool, *addresspool.Name)
			// SLB V2 BackendAddressPoolRef
			if len(strings.Split(*addresspool.Name, "/")) == 5 {
				wssdipconfig.LoadBalancerAddressPoolsRef = append(
					wssdipconfig.LoadBalancerAddressPoolsRef,
					&wssdcommonproto.BackendAddressPoolReference{
						ResourceRef: &wssdcommonproto.ResourceReference{
							Name: *addresspool.Name,
						},
					})
			}
		}
	}
	return wssdipconfig, nil
}

// Conversion function from wssdcloud network interface to network interface
func getNetworkInterface(server, group string, c *wssdcloudnetwork.NetworkInterface) (*network.Interface, error) {
	if c == nil {
		return &network.Interface{}, nil
	}

	ipConfigs := []network.InterfaceIPConfiguration{}
	for _, wssdipconfig := range c.IpConfigurations {
		ipConfigs = append(ipConfigs, *(getNetworkIpConfig(wssdipconfig)))
	}

	version := ""
	if c.Status != nil && c.Status.Version != nil {
		version = c.Status.Version.Number
	}

	advancedPolicies := network.GetNetworkAdvancedNetworkPolicies(c.AdvancedPolicies)

	vnetIntf := &network.Interface{
		Name:    &c.Name,
		ID:      &c.Id,
		Version: &version,
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			MacAddress: &c.Macaddress,
			// TODO: Type
			IPConfigurations:            &ipConfigs,
			Statuses:                    status.GetStatuses(c.GetStatus()),
			EnableAcceleratedNetworking: getIovSetting(c),
			DNSSettings:                 getWssdDNSSettings(c.Dns),
			AdvancedNetworkPolicies:     &advancedPolicies,
		},
		Tags: tags.ProtoToMap(c.Tags),
	}

	return vnetIntf, nil
}

func getDns(dnssetting *network.InterfaceDNSSettings) *wssdcommonproto.Dns {
	if dnssetting == nil {
		return nil
	}
	var dns wssdcommonproto.Dns
	if dnssetting.DNSServers != nil {
		dns.Servers = *dnssetting.DNSServers
	}
	if dnssetting.InternalDomainNameSuffix != nil {
		dns.Domain = *dnssetting.InternalDomainNameSuffix
	}
	return &dns
}

func getNetworkIpConfig(wssdcloudipconfig *wssdcloudnetwork.IpConfiguration) *network.InterfaceIPConfiguration {
	ipconfig := &network.InterfaceIPConfiguration{
		InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAddress: &wssdcloudipconfig.Ipaddress,
			Gateway:          &wssdcloudipconfig.Gateway,
			PrefixLength:     &wssdcloudipconfig.Prefixlength,
			Primary:          &wssdcloudipconfig.Primary,
		},
	}

	if wssdcloudipconfig.Subnet != nil {
		ipconfig.Subnet = getQualifiedSubnetReference(wssdcloudipconfig.Subnet)
	} else {
		ipconfig.Subnet = &network.APIEntityReference{ID: &wssdcloudipconfig.Subnetid}
	}

	if wssdcloudipconfig.NetworkSecurityGroupRef != nil && wssdcloudipconfig.NetworkSecurityGroupRef.ResourceRef != nil {
		ipconfig.InterfaceIPConfigurationPropertiesFormat.NetworkSecurityGroup = &network.SubResource{
			ID: &wssdcloudipconfig.NetworkSecurityGroupRef.ResourceRef.Name,
		}
	}

	// Set public IP Address reference
	if wssdcloudipconfig.PublicIPAddressRef != nil && wssdcloudipconfig.PublicIPAddressRef.ResourceRef != nil {
		ipconfig.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress = &network.PublicIPAddress{
			ID: &wssdcloudipconfig.PublicIPAddressRef.ResourceRef.Name,
		}
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

	var inboundNatRules []network.InboundNatRule
	for _, inboundnat := range wssdcloudipconfig.InboundNatRules {
		natrule := network.InboundNatRule{
			Name: &inboundnat.Name,
		}
		inboundNatRules = append(inboundNatRules, natrule)
	}

	ipconfig.LoadBalancerInboundNatRules = &inboundNatRules
	return ipconfig
}

func getIovSetting(vnic *wssdcloudnetwork.NetworkInterface) *bool {
	isAcceleratedNetworkingEnabled := false
	if vnic.IovWeight > 0 {
		isAcceleratedNetworkingEnabled = true
	}
	return &isAcceleratedNetworkingEnabled
}

func getWssdSubnetReference(subnetId *string) (*wssdcommonproto.SubnetReference, error) {

	// /virtualnetworks/<networkname>/subnets/<subnetname>
	subnetComponents := strings.Split(*subnetId, "/")

	if len(subnetComponents) != 5 {
		return nil, nil
	}

	networkRef := &wssdcommonproto.NetworkReference{
		ResourceRef: &wssdcommonproto.ResourceReference{
			Name: subnetComponents[2],
		},
	}

	if strings.EqualFold(subnetComponents[1], VNET_PREFIX) {
		networkRef.NetworkType = wssdcommonproto.NetworkType_VIRTUAL_NETWORK
	} else if strings.EqualFold(subnetComponents[1], LNET_PREFIX) || strings.EqualFold(subnetComponents[1], LNET_PREFIX_LEGACY) {
		networkRef.NetworkType = wssdcommonproto.NetworkType_LOGICAL_NETWORK
	} else {
		return nil, errors.Wrapf(errors.InvalidInput, "Cannot parse network type for the vnic")
	}

	return &wssdcommonproto.SubnetReference{
		Network: networkRef,
		ResourceRef: &wssdcommonproto.ResourceReference{
			Name: subnetComponents[4],
		},
	}, nil
}

func getQualifiedSubnetReference(subnet *wssdcommonproto.SubnetReference) *network.APIEntityReference {

	var networkPrefix string

	if subnet.GetNetwork().NetworkType == wssdcommonproto.NetworkType_VIRTUAL_NETWORK {
		networkPrefix = VNET_PREFIX
	} else {
		networkPrefix = LNET_PREFIX

	}

	id := "/" + networkPrefix + "/" + subnet.GetNetwork().GetResourceRef().GetName() + "/subnets/" + subnet.GetResourceRef().GetName()

	return &network.APIEntityReference{
		ID: &id,
	}
}
