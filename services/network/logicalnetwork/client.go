// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package logicalnetwork

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.LogicalNetwork, error)
	CreateOrUpdate(context.Context, string, string, *network.LogicalNetwork) (*network.LogicalNetwork, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, location string, logicalNetworks []*network.LogicalNetwork) (bool, error)
	UpdateRegisteredIPs(ctx context.Context, locationName, name string, subnetRegisteredIPs []SubnetRegisteredIPs) (subnetPersistedIPs []SubnetRegisteredIPs, failures []IPAddressUpdateFailure, err error)
}

// Client structure
type LogicalNetworkClient struct {
	network.BaseClient
	internal Service
}

// NewClient method returns new client
func NewLogicalNetworkClient(cloudFQDN string, authorizer auth.Authorizer) (*LogicalNetworkClient, error) {
	c, err := newLogicalNetworkClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &LogicalNetworkClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LogicalNetworkClient) Get(ctx context.Context, location, name string) (*[]network.LogicalNetwork, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *LogicalNetworkClient) CreateOrUpdate(ctx context.Context, location, name string, network *network.LogicalNetwork) (*network.LogicalNetwork, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, network)
}

// Delete methods invokes delete of the logical network resource
func (c *LogicalNetworkClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// Prechecks whether the system is able to create specified logicalNetworks.
// Returns true if it is possible; or false with reason in error message if not.
func (c *LogicalNetworkClient) Precheck(ctx context.Context, location string, logicalNetworks []*network.LogicalNetwork) (bool, error) {
	return c.internal.Precheck(ctx, location, logicalNetworks)
}

// UpdateRegisteredIPs replaces the registeredIPAddresses list on each
// referenced subnet with the supplied list (subnet-scoped full-replace) and
// applies the corresponding bitmap delta in MOC.
//
// subnetPersistedIPs reflects what is now stored in MOC (only trustworthy when err
// is nil). failures contains IP-level rejections (may be non-empty even when err is
// nil - partial success). A non-nil err signals a whole-call failure
// (LNET-not-found, lookup error, lock-acquisition error, store save error,
// transport error).
func (c *LogicalNetworkClient) UpdateRegisteredIPs(ctx context.Context, locationName, name string, subnetRegisteredIPs []SubnetRegisteredIPs) (subnetPersistedIPs []SubnetRegisteredIPs, failures []IPAddressUpdateFailure, err error) {
	return c.internal.UpdateRegisteredIPs(ctx, locationName, name, subnetRegisteredIPs)
}
