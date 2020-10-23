// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package controlplane

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.ControlPlaneInfo, error)
	CreateOrUpdate(context.Context, string, string, *cloud.ControlPlaneInfo) (*cloud.ControlPlaneInfo, error)
	Delete(context.Context, string, string) error
}

type ControlPlaneClient struct {
	internal Service
}

func NewControlPlaneClient(cloudFQDN string, authorizer auth.Authorizer) (*ControlPlaneClient, error) {
	c, err := newControlPlaneClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &ControlPlaneClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *ControlPlaneClient) Get(ctx context.Context, location, name string) (*[]cloud.ControlPlaneInfo, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *ControlPlaneClient) CreateOrUpdate(ctx context.Context, location, name string, cloud *cloud.ControlPlaneInfo) (*cloud.ControlPlaneInfo, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *ControlPlaneClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
