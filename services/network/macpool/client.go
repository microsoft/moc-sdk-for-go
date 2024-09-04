// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package macpool

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.MACPool, error)
	CreateOrUpdate(context.Context, string, string, *network.MACPool) (*network.MACPool, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, location string, macPools []*network.MACPool) (bool, error)
}

// MacPoolClient structure
type MacPoolClient struct {
	network.BaseClient
	internal Service
}

// NewMacPoolClient method returns new client
func NewMacPoolClient(cloudFQDN string, authorizer auth.Authorizer) (*MacPoolClient, error) {
	c, err := newMacPoolClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &MacPoolClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *MacPoolClient) Get(ctx context.Context, location, name string) (*[]network.MACPool, error) {
	return c.internal.Get(ctx, location, name)
}

// Ensure methods invokes create or update on the client
func (c *MacPoolClient) CreateOrUpdate(ctx context.Context, location, name string, macpool *network.MACPool) (*network.MACPool, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, macpool)
}

// Delete methods invokes delete of the network resource
func (c *MacPoolClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *MacPoolClient) Precheck(ctx context.Context, location string, macPools []*network.MACPool) (bool, error) {
	return c.internal.Precheck(ctx, location, macPools)
}
