// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package group

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.Group, error)
	CreateOrUpdate(context.Context, string, string, *cloud.Group) (*cloud.Group, error)
	Delete(context.Context, string, string) error
}

type GroupClient struct {
	internal Service
}

func NewGroupClient(cloudFQDN string, authorizer auth.Authorizer) (*GroupClient, error) {
	c, err := newGroupClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &GroupClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *GroupClient) Get(ctx context.Context, location, name string) (*[]cloud.Group, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *GroupClient) CreateOrUpdate(ctx context.Context, location, name string, cloud *cloud.Group) (*cloud.Group, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *GroupClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
