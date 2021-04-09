// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalhost

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.BareMetalHost, error)
	CreateOrUpdate(context.Context, string, string, *compute.BareMetalHost) (*compute.BareMetalHost, error)
	Delete(context.Context, string, string) error
	Query(context.Context, string, string) (*[]compute.BareMetalHost, error)
}

type BareMetalHostClient struct {
	compute.BaseClient
	internal Service
}

func NewBareMetalHostClient(cloudFQDN string, authorizer auth.Authorizer) (*BareMetalHostClient, error) {
	c, err := newBareMetalHostClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &BareMetalHostClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *BareMetalHostClient) Get(ctx context.Context, location, name string) (*[]compute.BareMetalHost, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *BareMetalHostClient) CreateOrUpdate(ctx context.Context, location, name string, compute *compute.BareMetalHost) (*compute.BareMetalHost, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *BareMetalHostClient) Delete(ctx context.Context, location string, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// Query method invokes the client Get method and uses the provided query to filter the returned results
func (c *BareMetalHostClient) Query(ctx context.Context, location, query string) (*[]compute.BareMetalHost, error) {
	return c.internal.Query(ctx, location, query)
}
