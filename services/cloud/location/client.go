// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package location

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

type Service interface {
	Get(context.Context, string) (*[]cloud.Location, error)
	CreateOrUpdate(context.Context, string, *cloud.Location) (*cloud.Location, error)
	Delete(context.Context, string) error
}

type LocationClient struct {
	internal Service
}

func NewLocationClient(cloudFQDN string, authorizer auth.Authorizer) (*LocationClient, error) {
	c, err := newLocationClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &LocationClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LocationClient) Get(ctx context.Context, name string) (*[]cloud.Location, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *LocationClient) CreateOrUpdate(ctx context.Context, name string, cloud *cloud.Location) (*cloud.Location, error) {
	return c.internal.CreateOrUpdate(ctx, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *LocationClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
