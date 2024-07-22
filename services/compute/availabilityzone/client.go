// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string) (*[]compute.AvailabilityZone, error)
	Create(ctx context.Context, name string, avzone *compute.AvailabilityZone) (*compute.AvailabilityZone, error)
	Delete(context.Context, string) error
}

type AvailabilityZoneClient struct {
	compute.BaseClient
	internal Service
}

func NewAvailabilityZoneClient(cloudFQDN string, authorizer auth.Authorizer) (*AvailabilityZoneClient, error) {
	c, err := newAvailabilityZoneClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AvailabilityZoneClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *AvailabilityZoneClient) Get(ctx context.Context, name string) (*[]compute.AvailabilityZone, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *AvailabilityZoneClient) Create(ctx context.Context, name string, compute *compute.AvailabilityZone) (*compute.AvailabilityZone, error) {
	return c.internal.Create(ctx, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *AvailabilityZoneClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}