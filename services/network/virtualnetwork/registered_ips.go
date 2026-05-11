// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualnetwork

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network/registeredips"
	"github.com/microsoft/moc/pkg/errors"

	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
)

// SubnetRegisteredIPs is the desired registered-IP list for one subnet.
// It is re-exported from the shared registeredips package.
type SubnetRegisteredIPs = registeredips.SubnetRegisteredIPs

// IPAddressUpdateFailure is the single-IP failure record. It is re-exported
// from the shared registeredips package.
type IPAddressUpdateFailure = registeredips.IPAddressUpdateFailure

// IPUpdateErrorCode is the failure-code enum. It is re-exported from the
// shared registeredips package; the String, MarshalJSON, and MarshalYAML
// methods defined on the underlying type are accessible through this alias.
type IPUpdateErrorCode = registeredips.IPUpdateErrorCode

const (
	IPUpdateUnknown          = registeredips.IPUpdateUnknown
	IPUpdateInvalidFormat    = registeredips.IPUpdateInvalidFormat
	IPUpdateOutOfRange       = registeredips.IPUpdateOutOfRange
	IPUpdateSubnetNotFound   = registeredips.IPUpdateSubnetNotFound
	IPUpdateAlreadyAllocated = registeredips.IPUpdateAlreadyAllocated
	IPUpdateNoPoolsInSubnet  = registeredips.IPUpdateNoPoolsInSubnet
)

// UpdateRegisteredIPs wraps the gRPC UpdateRegisteredIPs RPC. See
// VirtualNetworkClient.UpdateRegisteredIPs in interfaces.go for the full
// contract (subnet-scoped full-replace, IP-level best effort, partial-success
// semantics).
func (c *client) UpdateRegisteredIPs(ctx context.Context, groupName, name string, subnetRegisteredIPs []SubnetRegisteredIPs) (subnetPersistedIPs []SubnetRegisteredIPs, failures []IPAddressUpdateFailure, err error) {
	if len(groupName) == 0 {
		return nil, nil, errors.Wrapf(errors.InvalidInput, "GroupName is not specified")
	}
	if len(name) == 0 {
		return nil, nil, errors.Wrapf(errors.InvalidInput, "Name is not specified")
	}

	req := &wssdcloudnetwork.VirtualNetworkIPUpdateRequest{
		GroupName: groupName,
		Name:      name,
		IPUpdates: subnetsToProto(subnetRegisteredIPs),
	}

	resp, err := c.VirtualNetworkAgentClient.UpdateRegisteredIPs(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return registeredips.SubnetsFromProto(resp.GetPersistedIPs()),
		registeredips.FailuresFromProto(resp.GetFailures()),
		nil
}

func subnetsToProto(in []SubnetRegisteredIPs) []*wssdcloudnetwork.VirtualSubnetIPUpdate {
	out := make([]*wssdcloudnetwork.VirtualSubnetIPUpdate, 0, len(in))
	for _, s := range in {
		out = append(out, &wssdcloudnetwork.VirtualSubnetIPUpdate{
			SubnetName:            s.SubnetName,
			RegisteredIPAddresses: append([]string(nil), s.RegisteredIPAddresses...),
		})
	}
	return out
}
