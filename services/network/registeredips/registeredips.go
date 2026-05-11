// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

// Package registeredips holds the SDK-side types and proto-conversion
// helpers shared by the virtualnetwork and logicalnetwork registered-IPs
// surfaces.
//
// VNET and LNET have parallel proto messages (VirtualSubnetIPUpdate and
// LogicalSubnetIPUpdate) with identical field shape, so the SDK code is
// otherwise byte-identical between the two networks. Each network package
// re-exports the public types here via Go type aliases.
package registeredips

import (
	"encoding/json"
	"fmt"

	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

// SubnetRegisteredIPs identifies a subnet on the target network and the
// desired (full-replace) registered IP list for that subnet.
type SubnetRegisteredIPs struct {
	// Subnet name on the target network (NOT the network name itself;
	// the network is identified by separate parameters on the
	// UpdateRegisteredIPs call).
	SubnetName            string
	RegisteredIPAddresses []string
}

// IPAddressUpdateFailure describes a single IP that MOC rejected during an
// UpdateRegisteredIPs call.
type IPAddressUpdateFailure struct {
	SubnetName string
	IPAddress  string
	Code       IPUpdateErrorCode
	Error      string
}

// IPUpdateErrorCode mirrors moc's wssdcommon IPUpdateErrorCode enum.
type IPUpdateErrorCode int32

const (
	IPUpdateUnknown          IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_UNKNOWN)
	IPUpdateInvalidFormat    IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_INVALID_FORMAT)
	IPUpdateOutOfRange       IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_OUT_OF_RANGE)
	IPUpdateSubnetNotFound   IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_SUBNET_NOT_FOUND)
	IPUpdateAlreadyAllocated IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_ALREADY_ALLOCATED)
	IPUpdateNoPoolsInSubnet  IPUpdateErrorCode = IPUpdateErrorCode(wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_NO_POOLS_IN_SUBNET)
)

// String returns a human-readable PascalCase name for the error code, e.g.
// "SubnetNotFound". Falls back to the underlying numeric value for unknown
// codes.
func (c IPUpdateErrorCode) String() string {
	switch c {
	case IPUpdateUnknown:
		return "Unknown"
	case IPUpdateInvalidFormat:
		return "InvalidFormat"
	case IPUpdateOutOfRange:
		return "OutOfRange"
	case IPUpdateSubnetNotFound:
		return "SubnetNotFound"
	case IPUpdateAlreadyAllocated:
		return "AlreadyAllocated"
	case IPUpdateNoPoolsInSubnet:
		return "NoPoolsInSubnet"
	default:
		return fmt.Sprintf("IPUpdateErrorCode(%d)", int32(c))
	}
}

// MarshalJSON renders the code as its readable string form.
func (c IPUpdateErrorCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// MarshalYAML renders the code as its readable string form.
func (c IPUpdateErrorCode) MarshalYAML() (interface{}, error) {
	return c.String(), nil
}

// SubnetIPProto is satisfied by *wssdcloudnetwork.VirtualSubnetIPUpdate and
// *wssdcloudnetwork.LogicalSubnetIPUpdate. Both are comparable proto pointer
// types and expose the same (SubnetName, RegisteredIPAddresses) getter pair.
type SubnetIPProto interface {
	comparable
	GetSubnetName() string
	GetRegisteredIPAddresses() []string
}

// SubnetsFromProto translates a slice of subnet IP-update proto pointers
// into the SDK's portable SubnetRegisteredIPs slice. nil entries are
// skipped defensively against malformed proto responses.
//
// The RegisteredIPAddresses slice on each result entry is always
// non-nil; a subnet whose list was just cleared is rendered as an empty
// slice rather than nil.
func SubnetsFromProto[T SubnetIPProto](in []T) []SubnetRegisteredIPs {
	if len(in) == 0 {
		return nil
	}
	out := make([]SubnetRegisteredIPs, 0, len(in))
	var zero T
	for _, s := range in {
		if s == zero {
			continue
		}
		ips := s.GetRegisteredIPAddresses()
		out = append(out, SubnetRegisteredIPs{
			SubnetName:            s.GetSubnetName(),
			RegisteredIPAddresses: append(make([]string, 0, len(ips)), ips...),
		})
	}
	return out
}

// FailuresFromProto translates the failure slice from either RPC response
// into the portable SDK failure slice.
func FailuresFromProto(in []*wssdcloudcommon.IPAddressUpdateFailure) []IPAddressUpdateFailure {
	if len(in) == 0 {
		return nil
	}
	out := make([]IPAddressUpdateFailure, 0, len(in))
	for _, f := range in {
		if f == nil {
			continue
		}
		out = append(out, IPAddressUpdateFailure{
			SubnetName: f.GetSubnetName(),
			IPAddress:  f.GetIPAddress(),
			Code:       IPUpdateErrorCode(f.GetCode()),
			Error:      f.GetError(),
		})
	}
	return out
}
