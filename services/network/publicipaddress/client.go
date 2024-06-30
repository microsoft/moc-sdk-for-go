// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package publicipaddress

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.PublicIPAddress, error)
	CreateOrUpdate(context.Context, string, string, *network.PublicIPAddress) (*network.PublicIPAddress, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, group string, pips []*network.PublicIPAddress) (bool, error)
}

// PublicIPAddressAgentClient structure
type PublicIPAddressAgentClient struct {
	network.BaseClient
	internal Service
}

// PublicIPAddressAgentClient method returns new client
func NewPublicIPAddressClient(cloudFQDN string, authorizer auth.Authorizer) (*PublicIPAddressAgentClient, error) {
	c, err := newPublicIPAddressAgentClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &PublicIPAddressAgentClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *PublicIPAddressAgentClient) Get(ctx context.Context, group, name string) (*[]network.PublicIPAddress, error) {
	return c.internal.Get(ctx, group, name)
}

// Ensure methods invokes create or update on the client
func (c *PublicIPAddressAgentClient) CreateOrUpdate(ctx context.Context, group, name string, pip *network.PublicIPAddress) (*network.PublicIPAddress, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, pip)
}

// Delete methods invokes delete of the network resource
func (c *PublicIPAddressAgentClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *PublicIPAddressAgentClient) Precheck(ctx context.Context, group string, pip []*network.PublicIPAddress) (bool, error) {
	return c.internal.Precheck(ctx, group, pip)
}
