// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualnetwork

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.VirtualNetwork, error)
	CreateOrUpdate(context.Context, string, string, *network.VirtualNetwork) (*network.VirtualNetwork, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, group string, virtualNetworks []*network.VirtualNetwork) (bool, error)
}

// Client structure
type VirtualNetworkClient struct {
	network.BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualNetworkClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualNetworkClient, error) {
	c, err := newVirtualNetworkClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualNetworkClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualNetworkClient) Get(ctx context.Context, group, name string) (*[]network.VirtualNetwork, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualNetworkClient) CreateOrUpdate(ctx context.Context, group, name string, network *network.VirtualNetwork) (*network.VirtualNetwork, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, network)
}

// Delete methods invokes delete of the network resource
func (c *VirtualNetworkClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *VirtualNetworkClient) Precheck(ctx context.Context, group string, virtualNetworks []*network.VirtualNetwork) (bool, error) {
	return c.internal.Precheck(ctx, group, virtualNetworks)
}
