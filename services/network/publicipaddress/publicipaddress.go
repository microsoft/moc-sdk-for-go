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

// GetWssdPublicIPAddress converts our internal representation of a PublicIPAddress (network.PublicIPAddress)
// to the cloud public IP address protobuf used by wssdcloudagent (wssdnetwork.PublicIPAddress)
// getWssdPublicIPAddress converts a network.PublicIPAddress object to a wssdcloudnetwork.PublicIPAddress object.
// It takes a network.PublicIPAddress and a group name as input parameters and returns a wssdcloudnetwork.PublicIPAddress object and an error if any.
func getWssdPublicIPAddress(networkPip *network.PublicIPAddress, group string) (wssdCloudPip *wssdcloudnetwork.PublicIPAddress, err error) {

	wssdCloudPip = &wssdcloudnetwork.PublicIPAddress{
		Name:      *networkPip.Name,
		GroupName: group,
	}

	if networkPip.Location != nil {
		wssdCloudPip.LocationName = *networkPip.Location
	}

	// Checks if the PublicIPAddressPropertiesFormat is nil and returns the initialized wssdCloudPip
	// because all other fields are under PublicIPAddressPropertiesFormat.
	if networkPip.PublicIPAddressPropertiesFormat == nil {
		return wssdCloudPip, nil
	}

	if networkPip.IdleTimeoutInMinutes != nil {
		wssdCloudPip.IdleTimeoutInMinutes = *networkPip.IdleTimeoutInMinutes
	}

	switch ipAllocationMethodSdkToProtobuf(networkPip.PublicIPAllocationMethod) {
	// Static allocation is not supported for this release, but we will need this check when it is supported
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

	if networkPip.Tags != nil {
		wssdCloudPip.Tags = tags.MapToProto(networkPip.Tags)
	}

	return wssdCloudPip, nil
}

// GetPublicIPAddress converts the cloud public IP address protobuf returned from wssdcloudagent (wssdcloudnetwork.PublicIPAddress)
// to our internal representation of a public IP address (network.PublicIPAddress)
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
		},
	}

	if wssdPip.Tags != nil {
		networkPip.Tags = tags.ProtoToMap(wssdPip.Tags)
	}

	return networkPip, nil
}

// ipAllocationMethodProtobufToSdk converts a protobuf IP allocation method to the SDK IP allocation method.
// It takes an IPAllocationMethod from the wssdcommonproto package and returns the corresponding IPAllocationMethod
// from the network package. If the input does not match any known allocation method, it defaults to returning network.Dynamic.
func ipAllocationMethodProtobufToSdk(allocation wssdcommonproto.IPAllocationMethod) network.IPAllocationMethod {
	switch allocation {
	case wssdcommonproto.IPAllocationMethod_Static:
		return network.Static
	case wssdcommonproto.IPAllocationMethod_Dynamic:
		return network.Dynamic
	}
	return network.Dynamic
}

// ipAllocationMethodSdkToProtobuf converts an IP allocation method from the SDK representation to the protobuf representation.
func ipAllocationMethodSdkToProtobuf(allocation network.IPAllocationMethod) wssdcommonproto.IPAllocationMethod {
	switch allocation {
	case network.Static:
		return wssdcommonproto.IPAllocationMethod_Static
	case network.Dynamic:
		return wssdcommonproto.IPAllocationMethod_Dynamic
	}
	return wssdcommonproto.IPAllocationMethod_Dynamic
}

// ipVersionProtobufToSdk converts a protobuf IPVersion to the corresponding SDK IPVersion. It takes an IPVersion
// from the wssdcommonproto package and returns the corresponding IPVersion from the network package.
// If the provided IPVersion does not match any known values, it defaults to returning network.IPv4.
func ipVersionProtobufToSdk(ipversion wssdcommonproto.IPVersion) network.IPVersion {
	switch ipversion {
	case wssdcommonproto.IPVersion_IPv4:
		return network.IPv4
	case wssdcommonproto.IPVersion_IPv6:
		return network.IPv6
	}
	return network.IPv4
}

// ipVersionSdkToProtobuf converts the given IP version from the SDK representation to the Protobuf representation.
func ipVersionSdkToProtobuf(ipversion network.IPVersion) wssdcommonproto.IPVersion {
	switch ipversion {
	case network.IPv4:
		return wssdcommonproto.IPVersion_IPv4
	case network.IPv6:
		return wssdcommonproto.IPVersion_IPv6
	}
	return wssdcommonproto.IPVersion_IPv4
}
