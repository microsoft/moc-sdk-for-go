// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package role

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string) (*[]security.Role, error)
	CreateOrUpdate(context.Context, string, *security.Role) (*security.Role, error)
	Delete(context.Context, string) error
}

// RoleClient structure
type RoleClient struct {
	security.BaseClient
	internal Service
}

// NewRoleClient method returns new client
func NewRoleClient(cloudFQDN string, authorizer auth.Authorizer) (*RoleClient, error) {
	c, err := newRoleClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &RoleClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *RoleClient) Get(ctx context.Context, name string) (*[]security.Role, error) {
	return c.internal.Get(ctx, name)
}

// Ensure methods invokes create or update on the client
func (c *RoleClient) CreateOrUpdate(ctx context.Context, name string, role *security.Role) (*security.Role, error) {
	return c.internal.CreateOrUpdate(ctx, name, role)
}

// Delete methods invokes delete of the security resource
func (c *RoleClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
