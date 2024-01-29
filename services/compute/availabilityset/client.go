// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.AvailabilitySet, error)
	CreateOrUpdate(ctx context.Context, group string, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error)
	Delete(context.Context, string, string) error
}

type AvailabilitySetClient struct {
	compute.BaseClient
	internal Service
}

func NewAvailabilitySetClient(cloudFQDN string, authorizer auth.Authorizer) (*AvailabilitySetClient, error) {
	c, err := newAvailabilitySetClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AvailabilitySetClient{internal: c}, nil
}

func NewAvailabilitySetMockClient(cloudFQDN string, authorizer auth.Authorizer) (*AvailabilitySetClient, error) {
	c, err := newAvailabilitySetMockClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AvailabilitySetClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *AvailabilitySetClient) Get(ctx context.Context, group, name string) (*[]compute.AvailabilitySet, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *AvailabilitySetClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *AvailabilitySetClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
