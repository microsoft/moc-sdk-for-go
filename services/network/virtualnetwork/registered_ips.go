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

// UpdateRegisteredIPs wraps the gRPC UpdateRegisteredIPs RPC. See
// VirtualNetworkClient.UpdateRegisteredIPs in interfaces.go for the full
// contract (subnet-scoped full-replace, IP-level best effort, partial-success
// semantics).
func (c *client) UpdateRegisteredIPs(ctx context.Context, groupName, name string, subnetRegisteredIPs []SubnetRegisteredIPs) (subnetPersistedIPs []SubnetRegisteredIPs, failures []IPAddressUpdateFailure, err error) {
	return c.UpdateRegisteredIPsWithVersion(ctx, groupName, name, subnetRegisteredIPs, Version_Default)
}

// UpdateRegisteredIPsWithVersion is the API-version-aware variant of
// UpdateRegisteredIPs.apiVersion=="" or "1.0" preserves the default v1
// behavior; apiVersion=="2.0" opts the caller out of the cluster-wide
// VNET->LNET migration shim on the cloudagent side, so the registered-IP
// write lands on the VNET provider's IPAM regardless of
// disableLogicalNetworkMigration. Mirrors the Get/CreateOrUpdate/Delete
// *WithVersion shape used elsewhere in the VNET SDK.
func (c *client) UpdateRegisteredIPsWithVersion(ctx context.Context, groupName, name string, subnetRegisteredIPs []SubnetRegisteredIPs, apiVersion string) (subnetPersistedIPs []SubnetRegisteredIPs, failures []IPAddressUpdateFailure, err error) {
	if len(groupName) == 0 {
		return nil, nil, errors.Wrapf(errors.InvalidInput, "GroupName is not specified")
	}
	if len(name) == 0 {
		return nil, nil, errors.Wrapf(errors.InvalidInput, "Name is not specified")
	}

	version, err := getApiVersion(apiVersion)
	if err != nil {
		return nil, nil, err
	}

	req := &wssdcloudnetwork.VirtualNetworkIPUpdateRequest{
		GroupName: groupName,
		Name:      name,
		IPUpdates: subnetsToProto(subnetRegisteredIPs),
		Version:   version,
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
