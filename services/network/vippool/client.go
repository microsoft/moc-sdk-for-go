// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package vippool

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.VipPool, error)
	CreateOrUpdate(context.Context, string, string, *network.VipPool) (*network.VipPool, error)
	Delete(context.Context, string, string) error
}

// VipPoolClient structure
type VipPoolClient struct {
	network.BaseClient
	internal Service
}

// NewVipPoolClient method returns new client
func NewVipPoolClient(cloudFQDN string, authorizer auth.Authorizer) (*VipPoolClient, error) {
	c, err := newVipPoolClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VipPoolClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VipPoolClient) Get(ctx context.Context, location, name string) (*[]network.VipPool, error) {
	return c.internal.Get(ctx, location, name)
}

// Ensure methods invokes create or update on the client
func (c *VipPoolClient) CreateOrUpdate(ctx context.Context, location, name string, vp *network.VipPool) (*network.VipPool, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, vp)
}

// Delete methods invokes delete of the network resource
func (c *VipPoolClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
