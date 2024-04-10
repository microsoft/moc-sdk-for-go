// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package networksecuritygroup

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.SecurityGroup, error)
	CreateOrUpdate(context.Context, string, string, *network.SecurityGroup) (*network.SecurityGroup, error)
	Delete(context.Context, string, string) error
}

// NetworkSecurityGroupAgentClient structure
type NetworkSecurityGroupAgentClient struct {
	network.BaseClient
	internal Service
}

// NeNetworkSecurityGroupClient method returns new client
func NewSecurityGroupClient(cloudFQDN string, authorizer auth.Authorizer) (*NetworkSecurityGroupAgentClient, error) {
	c, err := newNetworkSecurityGroupClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &NetworkSecurityGroupAgentClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *NetworkSecurityGroupAgentClient) Get(ctx context.Context, location, name string) (*[]network.SecurityGroup, error) {
	return c.internal.Get(ctx, location, name)
}

// Ensure methods invokes create or update on the client
func (c *NetworkSecurityGroupAgentClient) CreateOrUpdate(ctx context.Context, location, name string, nsg *network.SecurityGroup) (*network.SecurityGroup, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, nsg)
}

// Delete methods invokes delete of the network resource
func (c *NetworkSecurityGroupAgentClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
