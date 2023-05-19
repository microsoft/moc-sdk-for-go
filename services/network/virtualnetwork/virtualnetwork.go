// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package virtualnetwork

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
func getWssdVirtualNetwork(c *network.VirtualNetwork, groupName string) (*wssdcloudnetwork.VirtualNetwork, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Network name is missing")
	}
	if len(groupName) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	wssdnetwork := &wssdcloudnetwork.VirtualNetwork{
		Name:      *c.Name,
		GroupName: groupName,
		Tags:      tags.MapToProto(c.Tags),
	}

	if c.Version != nil {
		if wssdnetwork.Status == nil {
			wssdnetwork.Status = status.InitStatus()
		}
		wssdnetwork.Status.Version.Number = *c.Version
	}

	if c.Location != nil {
		wssdnetwork.LocationName = *c.Location
	}

	if c.VirtualNetworkPropertiesFormat != nil {
		subnets, err := getWssdNetworkSubnets(c.Subnets)
		if err != nil {
			return nil, err
		}
		wssdnetwork.Subnets = subnets

		if c.VirtualNetworkPropertiesFormat.MacPoolName != nil {
			wssdnetwork.MacPoolName = *c.VirtualNetworkPropertiesFormat.MacPoolName
		}

		if c.DhcpOptions != nil && c.DhcpOptions.DNSServers != nil {
			wssdnetwork.Dns = &wssdcommonproto.Dns{
				Servers: *c.DhcpOptions.DNSServers,
			}
		}
	}

	if c.Type == nil {
		emptyString := ""
		c.Type = &emptyString
	}

	networkType, err := virtualNetworkTypeFromString(*c.Type)
	if err != nil {
		return nil, err
	}

	wssdnetwork.Type = networkType

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

func getWssdNetworkIPPoolInfo(ippoolinfo *network.IPPoolInfo) *wssdcloudnetwork.IPPoolInfo {
	if ippoolinfo != nil {
		return &wssdcloudnetwork.IPPoolInfo{
			Used:      ippoolinfo.Used,
			Available: ippoolinfo.Available,
		}
	}
	return nil
}
func getWssdNetworkSubnets(subnets *[]network.Subnet) (wssdsubnets []*wssdcloudnetwork.Subnet, err error) {
	if subnets == nil {
		return
	}

	for _, subnet := range *subnets {
		wssdsubnet := &wssdcloudnetwork.Subnet{}
		if subnet.Name == nil {
			err = errors.Wrapf(errors.InvalidInput, "Name is missing in subnet")
			return
		}
		wssdsubnet.Name = *subnet.Name

		if subnet.SubnetPropertiesFormat == nil {
			continue
		}

		if subnet.Vlan == nil {
			wssdsubnet.Vlan = 0
		} else {
			wssdsubnet.Vlan = uint32(*subnet.Vlan)
		}

		if subnet.TrunkVlan == nil {
			wssdsubnet.Trunkvlan = nil
		} else {
			trunkVlan := *subnet.TrunkVlan
			allowedVlanIdList := []uint32{}
			nativeVlanId := uint32(0)
			if trunkVlan.AllowedVlanIdList != nil {
				allowedVlanIdList = *trunkVlan.AllowedVlanIdList
			}
			if subnet.TrunkVlan.NativeVlanId != nil {
				nativeVlanId = *trunkVlan.NativeVlanId
			}
			wssdsubnet.Trunkvlan = &wssdcloudnetwork.TrunkVlan{
				Allowedvlanidlist: allowedVlanIdList,
				Nativevlanid:      nativeVlanId,
			}

		}

		wssdsubnetRoutes, err1 := getWssdNetworkRoutes(subnet.RouteTable)
		if err1 != nil {
			err = err1
			return
		}
		wssdsubnet.Routes = wssdsubnetRoutes
		wssdsubnet.Allocation = ipAllocationMethodSdkToProtobuf(subnet.IPAllocationMethod)

		if subnet.AddressPrefix != nil {
			wssdsubnet.Cidr = *subnet.AddressPrefix
		}

		//An address prefix is required if using ippools
		if len(subnet.IPPools) > 0 && subnet.AddressPrefix == nil {
			err = errors.Wrapf(errors.InvalidInput, "AddressPrefix is missing")
			return
		}

		for _, ippool := range subnet.IPPools {
			ippoolType := wssdcloudnetwork.IPPoolType_VM
			if strings.EqualFold(string(ippool.Type), string(network.VIPPOOL)) {
				ippoolType = wssdcloudnetwork.IPPoolType_VIPPool
			}
			wssdsubnet.Ippools = append(wssdsubnet.Ippools, &wssdcloudnetwork.IPPool{
				Name:  ippool.Name,
				Type:  ippoolType,
				Start: ippool.Start,
				End:   ippool.End,
				Info:  getWssdNetworkIPPoolInfo(ippool.Info),
			})
		}

		wssdsubnets = append(wssdsubnets, wssdsubnet)
	}

	return
}

func getWssdNetworkIpams(subnets *[]network.Subnet) []*wssdcloudnetwork.Ipam {
	ipam := wssdcloudnetwork.Ipam{}
	if subnets == nil {
		return []*wssdcloudnetwork.Ipam{}
	}

	for _, subnet := range *subnets {
		wssdsubnet := &wssdcloudnetwork.Subnet{
			Name: *subnet.Name,
			// TODO: implement something for IPConfigurationReferences
		}

		if subnet.AddressPrefix != nil {
			wssdsubnet.Cidr = *subnet.AddressPrefix
		}
		routes, err := getWssdNetworkRoutes(subnet.RouteTable)
		if err != nil {
			routes = []*wssdcloudnetwork.Route{}
		}
		wssdsubnet.Routes = routes

		ipam.Subnets = append(ipam.Subnets, wssdsubnet)
	}

	return []*wssdcloudnetwork.Ipam{&ipam}
}

func getWssdNetworkRoutes(routetable *network.RouteTable) (wssdcloudroutes []*wssdcloudnetwork.Route, err error) {
	if routetable == nil {
		return
	}

	for _, route := range *routetable.Routes {
		// RouteTable is optional
		if route.RoutePropertiesFormat == nil {
			continue
		}
		if route.NextHopIPAddress == nil || route.AddressPrefix == nil {
			err = errors.Wrapf(errors.InvalidInput, "NextHopIpAddress or AddressPrefix is missing")
			return
		}

		wssdcloudroutes = append(wssdcloudroutes, &wssdcloudnetwork.Route{
			Nexthop:           *route.NextHopIPAddress,
			Destinationprefix: *route.AddressPrefix,
		})
	}

	return
}

// Conversion function from wssdcloudnetwork to network
func getVirtualNetwork(c *wssdcloudnetwork.VirtualNetwork, group string) *network.VirtualNetwork {
	stringType := virtualNetworkTypeToString(c.Type)
	dnsservers := []string{}
	if c.Dns != nil {
		dnsservers = c.Dns.Servers
	}
	return &network.VirtualNetwork{
		Name:     &c.Name,
		Location: &c.LocationName,
		ID:       &c.Id,
		Type:     &stringType,
		Version:  &c.Status.Version.Number,
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
			Subnets:     getNetworkSubnets(c.Subnets),
			Statuses:    status.GetStatuses(c.GetStatus()),
			MacPoolName: &c.MacPoolName,
			DhcpOptions: &network.DhcpOptions{
				DNSServers: &dnsservers,
			},
		},
		Tags: tags.ProtoToMap(c.Tags),
	}
}

func getNetworkSubnets(wssdsubnets []*wssdcloudnetwork.Subnet) *[]network.Subnet {
	subnets := []network.Subnet{}

	for _, subnet := range wssdsubnets {
		subnets = append(subnets, network.Subnet{
			Name: &subnet.Name,
			ID:   &subnet.Id,
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: &subnet.Cidr,
				RouteTable:    getNetworkRoutetable(subnet.Routes),
				// TODO: implement something for IPConfigurationReferences
				IPAllocationMethod: ipAllocationMethodProtobufToSdk(subnet.Allocation),
				Vlan:               getVlan(subnet.Vlan),
				IPPools:            getIPPools(subnet.Ippools),
				TrunkVlan:          getTrunkVlan(subnet.Trunkvlan),
			},
		})
	}

	return &subnets
}

func getNetworkIPPoolInfo(wssdcloudippool *wssdcloudnetwork.IPPool) *network.IPPoolInfo {
	if wssdcloudippool.Info != nil {
		return &network.IPPoolInfo{
			Used:      wssdcloudippool.Info.Used,
			Available: wssdcloudippool.Info.Available,
		}
	}
	return nil
}

func getIPPools(wssdcloudippools []*wssdcloudnetwork.IPPool) []network.IPPool {
	ippool := []network.IPPool{}
	for _, wssdcloudippool := range wssdcloudippools {
		ippoolType := network.VM
		if wssdcloudippool.Type == wssdcloudnetwork.IPPoolType_VIPPool {
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

func getNetworkRoutetable(wssdcloudroutes []*wssdcloudnetwork.Route) *network.RouteTable {
	routes := []network.Route{}

	for _, route := range wssdcloudroutes {
		routes = append(routes, network.Route{
			RoutePropertiesFormat: &network.RoutePropertiesFormat{
				NextHopIPAddress: &route.Nexthop,
				AddressPrefix:    &route.Destinationprefix,
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

func getTrunkVlan(wssdtrunkvlan *wssdcloudnetwork.TrunkVlan) *network.TrunkVlan {
	if wssdtrunkvlan == nil {
		return nil
	}
	vlan := network.TrunkVlan{
		AllowedVlanIdList: &wssdtrunkvlan.Allowedvlanidlist,
		NativeVlanId:      &wssdtrunkvlan.Nativevlanid,
	}
	return &vlan
}

func virtualNetworkTypeToString(vnetType wssdcloudnetwork.VirtualNetworkType) string {
	typename, ok := wssdcloudnetwork.VirtualNetworkType_name[int32(vnetType)]
	if !ok {
		return "Unknown"
	}
	return typename

}

func virtualNetworkTypeFromString(vnNetworkString string) (wssdcloudnetwork.VirtualNetworkType, error) {
	typevalue := wssdcloudnetwork.VirtualNetworkType_ICS
	if len(vnNetworkString) > 0 {
		typevTmp, ok := wssdcloudnetwork.VirtualNetworkType_value[vnNetworkString]
		if ok {
			typevalue = wssdcloudnetwork.VirtualNetworkType(typevTmp)
		}
	}
	return typevalue, nil
}
