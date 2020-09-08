// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package networkinterface

import (
	"context"
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.Interface, error)
	CreateOrUpdate(context.Context, string, string, *network.Interface) (*network.Interface, error)
	Delete(context.Context, string, string) error
}

// InterfaceClient structure
type InterfaceClient struct {
	network.BaseClient
	internal Service
}

// NewInterfaceClient method returns new client
func NewInterfaceClient(cloudFQDN string, authorizer auth.Authorizer) (*InterfaceClient, error) {
	c, err := newInterfaceClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &InterfaceClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *InterfaceClient) Get(ctx context.Context, group, name string) (*[]network.Interface, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *InterfaceClient) CreateOrUpdate(ctx context.Context, group, name string, networkInterface *network.Interface) (*network.Interface, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, networkInterface)
}

// Delete methods invokes delete of the network interface resource
func (c *InterfaceClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
