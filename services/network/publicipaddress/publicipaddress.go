// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package publicipaddress

import (
	"net"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
)

// GetWssdPublicIPAddress converts our internal representation of a PublicIPAddress (network.PublicIPAddress) to the cloud public IP address protobuf used by wssdcloudagent (wssdnetwork.PublicIPAddress)
func getWssdPublicIPAddress(networkPip *network.PublicIPAddress, group string) (wssdCloudPip *wssdcloudnetwork.PublicIPAddress, err error) {

	// // if networkPip == nil || networkPip.PublicIPAddressPropertiesFormat == nil {
	// // 	return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing public IP address Properties")
	// // }

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Group not specified")
	}

	if networkPip.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for public IP Address")
	}

	wssdCloudPip = &wssdcloudnetwork.PublicIPAddress{
		Name:      *networkPip.Name,
		GroupName: group,
		//IdleTimeoutInMinutes: uint32(*networkPip.IdleTimeoutInMinutes),
		//IpAddress:            *networkPip.IPAddress,
		//Allocation: ipAllocationMethodSdkToProtobuf(networkPip.PublicIPAllocationMethod),
		//IpVersion:  ipVersionSdkToProtobuf(networkPip.PublicIPAddressVersion),
		//DomainNameLabel:      *networkPip.DNSSettings.DomainNameLabel,
		//ReverseFqdn:          *networkPip.DNSSettings.ReverseFqdn,
	}

	if networkPip.Location != nil {
		wssdCloudPip.LocationName = *networkPip.Location
	}

	if networkPip.PublicIPAddressPropertiesFormat == nil {
		return wssdCloudPip, nil
	}

	if networkPip.IdleTimeoutInMinutes != nil {
		wssdCloudPip.IdleTimeoutInMinutes = uint32(*networkPip.IdleTimeoutInMinutes)
	}

	switch ipAllocationMethodSdkToProtobuf(networkPip.PublicIPAllocationMethod) {
	case wssdcommonproto.IPAllocationMethod_Static:
		if networkPip.IPAddress != nil {
			wssdCloudPip.IpAddress = *networkPip.IPAddress
			wssdCloudPip.Allocation = wssdcommonproto.IPAllocationMethod_Static
			parsedIP := net.ParseIP(*networkPip.IPAddress)
			switch ipVersionSdkToProtobuf(networkPip.PublicIPAddressVersion) {
			case wssdcommonproto.IPVersion_IPv4:
				if parsedIP != nil && parsedIP.To4() != nil {
					wssdCloudPip.IpVersion = wssdcommonproto.IPVersion_IPv4
				} else {
					return nil, errors.Wrapf(errors.InvalidInput, "Public IP address is not in IPv4 format")
				}
			case wssdcommonproto.IPVersion_IPv6:
				if parsedIP != nil && parsedIP.To16() != nil {
					wssdCloudPip.IpVersion = wssdcommonproto.IPVersion_IPv6
				} else {
					return nil, errors.Wrapf(errors.InvalidInput, "Public IP address is not in IPv6 format")
				}
			}

		} else {
			return nil, errors.Wrapf(errors.InvalidInput, "Missing public IP address with static allocation")
		}
	case wssdcommonproto.IPAllocationMethod_Dynamic:
		wssdCloudPip.Allocation = wssdcommonproto.IPAllocationMethod_Dynamic
		wssdCloudPip.IpVersion = ipVersionSdkToProtobuf(networkPip.PublicIPAddressVersion)
	}

	// No support for now
	// if networkPip.PublicIPAddressPropertiesFormat.DNSSettings == nil {
	// 	return wssdCloudPip, nil
	// }

	// if networkPip.DNSSettings.DomainNameLabel != nil {
	// 	wssdCloudPip.DomainNameLabel = *networkPip.DNSSettings.DomainNameLabel
	// }

	// if networkPip.DNSSettings.ReverseFqdn != nil {
	// 	wssdCloudPip.ReverseFqdn = *networkPip.DNSSettings.ReverseFqdn
	// }

	if networkPip.Tags != nil {
		wssdCloudPip.Tags = tags.MapToProto(networkPip.Tags)
	}

	return wssdCloudPip, nil
}

// GetPublicIPAddress converts the cloud public IP address protobuf returned from wssdcloudagent (wssdcloudnetwork.PublicIPAddress) to our internal representation of a public IP address (network.PublicIPAddress)
func getPublicIPAddress(wssdPip *wssdcloudnetwork.PublicIPAddress) (networkPip *network.PublicIPAddress, err error) {

	networkPip = &network.PublicIPAddress{
		Name:     &wssdPip.Name,
		Location: &wssdPip.LocationName,
		ID:       &wssdPip.Id,
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			Statuses:                 status.GetStatuses(wssdPip.GetStatus()),
			IdleTimeoutInMinutes:     &wssdPip.IdleTimeoutInMinutes,
			IPAddress:                &wssdPip.IpAddress,
			PublicIPAddressVersion:   ipVersionProtobufToSdk(wssdPip.IpVersion),
			PublicIPAllocationMethod: ipAllocationMethodProtobufToSdk(wssdPip.Allocation),
			// No support for now
			// DNSSettings: &network.PublicIPAddressDNSSettings{
			// 	DomainNameLabel: &wssdPip.DomainNameLabel,
			// 	Fqdn:            &wssdPip.Fqdn,
			// 	ReverseFqdn:     &wssdPip.ReverseFqdn,
			// },
		},
	}

	if wssdPip.Tags != nil {
		networkPip.Tags = tags.ProtoToMap(wssdPip.Tags)
	}

	return networkPip, nil
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

func ipVersionProtobufToSdk(ipversion wssdcommonproto.IPVersion) network.IPVersion {
	switch ipversion {
	case wssdcommonproto.IPVersion_IPv4:
		return network.IPv4
	case wssdcommonproto.IPVersion_IPv6:
		return network.IPv6
	}
	return network.IPv4
}

func ipVersionSdkToProtobuf(ipversion network.IPVersion) wssdcommonproto.IPVersion {
	switch ipversion {
	case network.IPv4:
		return wssdcommonproto.IPVersion_IPv4
	case network.IPv6:
		return wssdcommonproto.IPVersion_IPv6
	}
	return wssdcommonproto.IPVersion_IPv4
}
