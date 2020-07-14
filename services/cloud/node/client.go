// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.Node, error)
	CreateOrUpdate(context.Context, string, string, *cloud.Node) (*cloud.Node, error)
	Delete(context.Context, string, string) error
}

type NodeClient struct {
	internal Service
}

func NewNodeClient(cloudFQDN string, authorizer auth.Authorizer) (*NodeClient, error) {
	c, err := newNodeClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &NodeClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *NodeClient) Get(ctx context.Context, location, name string) (*[]cloud.Node, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *NodeClient) CreateOrUpdate(ctx context.Context, location, name string, cloud *cloud.Node) (*cloud.Node, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *NodeClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
