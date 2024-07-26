// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string) (*[]cloud.AvailabilityZone, error)
	CreateOrUpdate(ctx context.Context, name string, avzone *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error)
	Delete(context.Context, string) error
}

type AvailabilityZoneClient struct {
	cloud.BaseClient
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
func (c *AvailabilityZoneClient) Get(ctx context.Context, name string) (*[]cloud.AvailabilityZone, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *AvailabilityZoneClient) CreateOrUpdate(ctx context.Context, name string, cloud *cloud.AvailabilityZone) (*cloud.AvailabilityZone, error) {
	return c.internal.CreateOrUpdate(ctx, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *AvailabilityZoneClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}