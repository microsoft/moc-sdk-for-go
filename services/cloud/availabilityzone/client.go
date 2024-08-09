// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.AvailabilityZone, error)
	CreateOrUpdate(ctx context.Context, location string, name string, avzone *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, location string, avzones []*cloud.AvailabilityZone) (bool, error)
}

type AvailabilityZoneClient struct {
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
func (c *AvailabilityZoneClient) Get(ctx context.Context, location string, name string) (*[]cloud.AvailabilityZone, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *AvailabilityZoneClient) CreateOrUpdate(ctx context.Context, location string, name string, cloud *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *AvailabilityZoneClient) Delete(ctx context.Context, location string, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// Prechecks whether the system is able to create specified availability zones.
// Returns true if it is possible; or false with reason in error message if not.
func (c *AvailabilityZoneClient) Precheck(ctx context.Context, location string, avzones []*cloud.AvailabilityZone) (bool, error) {
	return c.internal.Precheck(ctx, location, avzones)
}
