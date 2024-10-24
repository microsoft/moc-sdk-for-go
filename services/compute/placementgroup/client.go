// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package placementgroup

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.PlacementGroup, error)
	Create(ctx context.Context, group string, name string, pgroup *compute.PlacementGroup) (*compute.PlacementGroup, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, group string, pgroups []*compute.PlacementGroup) (bool, error)
}

type PlacementGroupClient struct {
	compute.BaseClient
	internal Service
}

func NewPlacementGroupClient(cloudFQDN string, authorizer auth.Authorizer) (*PlacementGroupClient, error) {
	c, err := newPlacementGroupClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &PlacementGroupClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *PlacementGroupClient) Get(ctx context.Context, group, name string) (*[]compute.PlacementGroup, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *PlacementGroupClient) Create(ctx context.Context, group, name string, compute *compute.PlacementGroup) (*compute.PlacementGroup, error) {
	return c.internal.Create(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *PlacementGroupClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Prechecks whether the system is able to create specified placement groups.
// Returns true if it is possible; or false with reason in error message if not.
func (c *PlacementGroupClient) Precheck(ctx context.Context, group string, pgroups []*compute.PlacementGroup) (bool, error) {
	return c.internal.Precheck(ctx, group, pgroups)
}
