// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package zone

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.Zone, error)
	CreateOrUpdate(ctx context.Context, location string, name string, avzone *cloud.Zone) (*cloud.Zone, error)
	Delete(context.Context, string, string) error
}

type ZoneClient struct {
	internal Service
}

func NewZoneClient(cloudFQDN string, authorizer auth.Authorizer) (*ZoneClient, error) {
	c, err := newZoneClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &ZoneClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *ZoneClient) Get(ctx context.Context, location string, name string) (*[]cloud.Zone, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *ZoneClient) CreateOrUpdate(ctx context.Context, location string, name string, cloud *cloud.Zone) (*cloud.Zone, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *ZoneClient) Delete(ctx context.Context, location string, name string) error {
	return c.internal.Delete(ctx, location, name)
}
