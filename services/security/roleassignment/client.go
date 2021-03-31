// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package roleassignment

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, *security.RoleAssignment) (*[]security.RoleAssignment, error)
	CreateOrUpdate(context.Context, *security.RoleAssignment) (*security.RoleAssignment, error)
	Delete(context.Context, *security.RoleAssignment) error
}

// Client structure
type RoleAssignmentClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewRoleAssignmentClient(cloudFQDN string, authorizer auth.Authorizer) (*RoleAssignmentClient, error) {
	c, err := newRoleAssignmentClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &RoleAssignmentClient{internal: c}, nil
}

// Get invokes a role assignment retrieval
func (c *RoleAssignmentClient) Get(ctx context.Context, ra *security.RoleAssignment) (*[]security.RoleAssignment, error) {
	return c.internal.Get(ctx, ra)
}

// CreateOrUpdate method invokes a role assignment
func (c *RoleAssignmentClient) CreateOrUpdate(ctx context.Context, ra *security.RoleAssignment) (*security.RoleAssignment, error) {
	return c.internal.CreateOrUpdate(ctx, ra)
}

// Delete invokes a role assignment removal
func (c *RoleAssignmentClient) Delete(ctx context.Context, ra *security.RoleAssignment) error {
	return c.internal.Delete(ctx, ra)
}
