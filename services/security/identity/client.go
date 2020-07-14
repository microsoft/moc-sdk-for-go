// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"context"

	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/security"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]security.Identity, error)
	CreateOrUpdate(context.Context, string, string, *security.Identity) (*security.Identity, error)
	Delete(context.Context, string, string) error
}

// Client structure
type IdentityClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewIdentityClient(cloudFQDN string, authorizer auth.Authorizer) (*IdentityClient, error) {
	c, err := newIdentityClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &IdentityClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *IdentityClient) Get(ctx context.Context, group, name string) (*[]security.Identity, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *IdentityClient) CreateOrUpdate(ctx context.Context, group, name string, identity *security.Identity) (*security.Identity, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, identity)
}

// Delete methods invokes delete of the Identity resource
func (c *IdentityClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
