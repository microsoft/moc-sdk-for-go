// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package publicipaddress

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service defines the interface for managing Public IP Addresses in the network service.
// It provides methods to get, create or update, delete, and precheck public IP addresses.
type Service interface {
	Get(context.Context, string, string) (*[]network.PublicIPAddress, error)
	CreateOrUpdate(context.Context, string, string, *network.PublicIPAddress) (*network.PublicIPAddress, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, group string, pips []*network.PublicIPAddress) (bool, error)
}

// PublicIPAddressAgentClient is a client for managing public IP addresses.
// It embeds the network.BaseClient and includes an internal Service for additional functionality.
type PublicIPAddressAgentClient struct {
	network.BaseClient
	internal Service
}

// NewPublicIPAddressClient creates a new instance of PublicIPAddressAgentClient.
// It takes a cloudFQDN string and an authorizer of type auth.Authorizer as parameters.
// Returns a pointer to PublicIPAddressAgentClient and an error if the client creation fails.
func NewPublicIPAddressClient(cloudFQDN string, authorizer auth.Authorizer) (*PublicIPAddressAgentClient, error) {
	c, err := newPublicIPAddressAgentClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &PublicIPAddressAgentClient{internal: c}, nil
}

// Get retrieves a list of PublicIPAddresses from the specified resource group and name.
func (c *PublicIPAddressAgentClient) Get(ctx context.Context, group, name string) (*[]network.PublicIPAddress, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate creates or updates a Public IP Address in the specified resource group.
func (c *PublicIPAddressAgentClient) CreateOrUpdate(ctx context.Context, group, name string, pip *network.PublicIPAddress) (*network.PublicIPAddress, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, pip)
}

// Delete removes a public IP address resource identified by the specified group and name.
func (c *PublicIPAddressAgentClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *PublicIPAddressAgentClient) Precheck(ctx context.Context, group string, pip []*network.PublicIPAddress) (bool, error) {
	return c.internal.Precheck(ctx, group, pip)
}
