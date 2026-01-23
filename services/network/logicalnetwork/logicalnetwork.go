// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package logicalnetwork

import (
	"strings"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
)

// Conversion functions from network to wssdcloudnetwork
func getWssdLogicalNetwork(c *network.LogicalNetwork) (*wssdcloudnetwork.LogicalNetwork, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Logical Network name is missing")
	}
	if c.Location == nil || len(*c.Location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location is not specified")
	}

	wssdnetwork := &wssdcloudnetwork.LogicalNetwork{
		Name:         *c.Name,
		LocationName: *c.Location,
		Tags:         tags.MapToProto(c.Tags),
	}

	if c.Version != nil {
		if wssdnetwork.Status == nil {
			wssdnetwork.Status = status.InitStatus()
		}
		wssdnetwork.Status.Version.Number = *c.Version
	}

	if c.LogicalNetworkPropertiesFormat != nil {
		subnets, err := getWssdNetworkSubnets(c.Subnets, *c.Location)
		if err != nil {
			return nil, err
		}
		wssdnetwork.Subnets = subnets

		if c.LogicalNetworkPropertiesFormat.MacPoolName != nil {
			wssdnetwork.MacPoolName = *c.LogicalNetworkPropertiesFormat.MacPoolName
		}

		if c.LogicalNetworkPropertiesFormat.NetworkVirtualizationEnabled != nil {
			wssdnetwork.NetworkVirtualizationEnabled = *c.LogicalNetworkPropertiesFormat.NetworkVirtualizationEnabled
		}

		if c.LogicalNetworkPropertiesFormat.AdvancedNetworkPolicies != nil {
			wssdnetwork.AdvancedPolicies = network.GetWssdAdvancedNetworkPolicies(c.LogicalNetworkPropertiesFormat.AdvancedNetworkPolicies)
		}
	}

	return wssdnetwork, nil
}

func ipAllocationMethodProtobufToSdk(allocation wssdcommonproto.IPAllocationMethod) network.IPAllocationMethod {
	switch allocation {
	case wssdcommonproto.IPAllocationMethod_Static:
		return network.Static
	case wssdcommonproto.IPAllocationMethod_Dynamic:
		return network.Dynamic
	}
	return network.Dynamic
}

func ipAllocationMethodSdkToProtobuf(allocation network.IPAllocationMethod) wssdcommonproto.IPAllocationMethod {
	switch allocation {
	case network.Static:
		return wssdcommonproto.IPAllocationMethod_Static
	case network.Dynamic:
		return wssdcommonproto.IPAllocationMethod_Dynamic
	}
	return wssdcommonproto.IPAllocationMethod_Dynamic
}

func getWssdNetworkIPPoolInfo(ippoolinfo *network.IPPoolInfo) *wssdcommonproto.IPPoolInfo {
	if ippoolinfo != nil {
		return &wssdcommonproto.IPPoolInfo{
			Used:      ippoolinfo.Used,
			Available: ippoolinfo.Available,
		}
	}
	return nil
}
func getWssdNetworkSubnets(subnets *[]network.LogicalSubnet, location string) (wssdsubnets []*wssdcloudnetwork.LogicalSubnet, err error) {
	if subnets == nil {
		return
	}

	for _, subnet := range *subnets {
		wssdsubnet := &wssdcloudnetwork.LogicalSubnet{}
		if subnet.Name == nil {
			err = errors.Wrapf(errors.InvalidInput, "Subnet name is missing")
			return
		}
		wssdsubnet.Name = *subnet.Name

		if subnet.LogicalSubnetPropertiesFormat == nil {
			continue
		}

		if subnet.Vlan == nil {
			wssdsubnet.Vlan = 0
		} else {
			wssdsubnet.Vlan = uint32(*subnet.Vlan)
		}

		wssdsubnetRoutes, err1 := getWssdNetworkRoutes(subnet.RouteTable)
		if err1 != nil {
			err = err1
			return
		}
		wssdsubnet.Routes = wssdsubnetRoutes
		wssdsubnet.Allocation = ipAllocationMethodSdkToProtobuf(subnet.IPAllocationMethod)

		if subnet.AddressPrefix != nil {
			wssdsubnet.AddressPrefix = *subnet.AddressPrefix
		}

		// An address prefix is required if using ippools
		if len(subnet.IPPools) > 0 && subnet.AddressPrefix == nil {
			err = errors.Wrapf(errors.InvalidInput, "AddressPrefix is missing")
			return
		}

		if subnet.DhcpOptions != nil && subnet.DhcpOptions.DNSServers != nil {
			wssdsubnet.Dns = &wssdcommonproto.Dns{
				Servers: *subnet.DhcpOptions.DNSServers,
			}
		}

		for _, ippool := range subnet.IPPools {
			ippoolType := wssdcommonproto.IPPoolType_VM
			if strings.EqualFold(string(ippool.Type), string(network.VIPPOOL)) {
				ippoolType = wssdcommonproto.IPPoolType_VIPPool
			}
			wssdsubnet.IpPools = append(wssdsubnet.IpPools, &wssdcommonproto.IPPool{
				Name:  ippool.Name,
				Type:  ippoolType,
				Start: ippool.Start,
				End:   ippool.End,
				Info:  getWssdNetworkIPPoolInfo(ippool.Info),
			})
		}

		if subnet.Public != nil {
			wssdsubnet.IsPublic = *subnet.Public
		}

		if subnet.NetworkSecurityGroup != nil {
			wssdsubnet.NetworkSecurityGroupRef = &wssdcommonproto.NetworkSecurityGroupReference{
				ResourceRef: &wssdcommonproto.ResourceReference{
					Name: *subnet.NetworkSecurityGroup.ID,
				},
			}
		}

		wssdsubnets = append(wssdsubnets, wssdsubnet)
	}

	return
}

func getWssdNetworkRoutes(routetable *network.RouteTable) (wssdcloudroutes []*wssdcommonproto.Route, err error) {
	if routetable == nil {
		return
	}

	for _, route := range *routetable.Routes {
		// RouteTable is optional
		if route.RoutePropertiesFormat == nil {
			continue
		}
		if route.NextHopIPAddress == nil || route.AddressPrefix == nil {
			err = errors.Wrapf(errors.InvalidInput, "NextHopIpAddress or AddressPrefix is missing in Route")
			return
		}

		wssdcloudroutes = append(wssdcloudroutes, &wssdcommonproto.Route{
			NextHop:           *route.NextHopIPAddress,
			DestinationPrefix: *route.AddressPrefix,
		})
	}

	return
}

// Conversion function from wssdcloudnetwork to network
func getLogicalNetwork(c *wssdcloudnetwork.LogicalNetwork) *network.LogicalNetwork {

	advancedPolicies := network.GetNetworkAdvancedNetworkPolicies(c.AdvancedPolicies)

	lnet := &network.LogicalNetwork{
		Name:     &c.Name,
		Location: &c.LocationName,
		ID:       &c.Id,
		Version:  &c.Status.Version.Number,
		LogicalNetworkPropertiesFormat: &network.LogicalNetworkPropertiesFormat{
			Subnets:                      getNetworkSubnets(c.Subnets),
			Statuses:                     status.GetStatuses(c.GetStatus()),
			MacPoolName:                  &c.MacPoolName,
			NetworkVirtualizationEnabled: &c.NetworkVirtualizationEnabled,
			AdvancedNetworkPolicies:      &advancedPolicies,
		},
		Tags:                    tags.ProtoToMap(c.Tags),
		NetworkControllerConfig: getNetworkControllerConfig(c.NetworkControllerConfig),
	}

	return lnet

}

// getNetworkControllerConfig converts proto NetworkControllerConfig to SDK NetworkControllerConfig
func getNetworkControllerConfig(c *wssdcommonproto.NetworkControllerConfig) *network.NetworkControllerConfig {
	if c == nil {
		return nil
	}
	return &network.NetworkControllerConfig{
		IsSdnEnabled:       c.IsSdnEnabled,
		IsSdnVnetEnabled:   c.IsSdnVnetEnabled,
		IsSdnVnetV2Enabled: c.IsSdnVnetV2Enabled,
		IsSdnLBV2Enabled:   c.IsSdnLBV2Enabled,
		IsLegacySdnEnabled: c.IsLegacySdnEnabled,
	}
}

func getNetworkSubnets(wssdsubnets []*wssdcloudnetwork.LogicalSubnet) *[]network.LogicalSubnet {
	subnets := []network.LogicalSubnet{}

	for _, subnet := range wssdsubnets {
		dnsservers := []string{}
		if subnet.Dns != nil {
			dnsservers = subnet.Dns.Servers
		}
		subnets = append(subnets, network.LogicalSubnet{
			Name: &subnet.Name,
			ID:   &subnet.Id,
			LogicalSubnetPropertiesFormat: &network.LogicalSubnetPropertiesFormat{
				AddressPrefix: &subnet.AddressPrefix,
				RouteTable:    getNetworkRoutetable(subnet.Routes),
				// TODO: implement something for IPConfigurationReferences
				IPAllocationMethod: ipAllocationMethodProtobufToSdk(subnet.Allocation),
				Vlan:               getVlan(subnet.Vlan),
				IPPools:            getIPPools(subnet.IpPools),
				DhcpOptions: &network.DhcpOptions{
					DNSServers: &dnsservers,
				},
				NetworkSecurityGroup: getNetworkSecurityGroup(subnet.NetworkSecurityGroupRef),
				Public:               &subnet.IsPublic,
			},
		})
	}

	return &subnets
}

func getNetworkIPPoolInfo(wssdcloudippool *wssdcommonproto.IPPool) *network.IPPoolInfo {
	if wssdcloudippool.Info != nil {
		return &network.IPPoolInfo{
			Used:      wssdcloudippool.Info.Used,
			Available: wssdcloudippool.Info.Available,
		}
	}
	return nil
}

func getIPPools(wssdcloudippools []*wssdcommonproto.IPPool) []network.IPPool {
	ippool := []network.IPPool{}
	for _, wssdcloudippool := range wssdcloudippools {
		ippoolType := network.VM
		if wssdcloudippool.Type == wssdcommonproto.IPPoolType_VIPPool {
			ippoolType = network.VIPPOOL
		}
		ippool = append(ippool, network.IPPool{
			Name:  wssdcloudippool.Name,
			Type:  ippoolType,
			Start: wssdcloudippool.Start,
			End:   wssdcloudippool.End,
			Info:  getNetworkIPPoolInfo(wssdcloudippool),
		})
	}
	return ippool
}

func getNetworkRoutetable(wssdcloudroutes []*wssdcommonproto.Route) *network.RouteTable {
	routes := []network.Route{}

	for _, route := range wssdcloudroutes {
		routes = append(routes, network.Route{
			RoutePropertiesFormat: &network.RoutePropertiesFormat{
				NextHopIPAddress: &route.NextHop,
				AddressPrefix:    &route.DestinationPrefix,
			},
		})
	}

	return &network.RouteTable{
		RouteTablePropertiesFormat: &network.RouteTablePropertiesFormat{
			Routes: &routes,
		},
	}
}

func getVlan(wssdvlan uint32) *uint16 {
	vlan := uint16(wssdvlan)
	return &vlan
}

func getNetworkSecurityGroup(wssdNsg *wssdcommonproto.NetworkSecurityGroupReference) *network.SubResource {
	if wssdNsg == nil || wssdNsg.ResourceRef == nil {
		return nil
	}

	return &network.SubResource{
		ID: &wssdNsg.ResourceRef.Name,
	}
}
