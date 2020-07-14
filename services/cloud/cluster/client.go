// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package cluster

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

type Service interface {
	Get(context.Context, string, string) (*[]cloud.Cluster, error)
	GetNodes(context.Context, string, string) (*[]cloud.Node, error)
	Load(context.Context, string, string, *cloud.Cluster) (*cloud.Cluster, error)
	Unload(context.Context, string, string) error
}

type ClusterClient struct {
	internal Service
}

func NewClusterClient(cloudFQDN string, authorizer auth.Authorizer) (*ClusterClient, error) {
	c, err := newClusterClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &ClusterClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *ClusterClient) Get(ctx context.Context, location, name string) (*[]cloud.Cluster, error) {
	return c.internal.Get(ctx, location, name)
}

// GetNodes methods invokes the client GetNodes method
func (c *ClusterClient) GetNodes(ctx context.Context, location, name string) (*[]cloud.Node, error) {
	return c.internal.GetNodes(ctx, location, name)
}

// Load methods invokes create or update on the client
func (c *ClusterClient) Load(ctx context.Context, location, name string, cloud *cloud.Cluster) (*cloud.Cluster, error) {
	return c.internal.Load(ctx, location, name, cloud)
}

// Unload methods invokes delete of the cloud resource
func (c *ClusterClient) Unload(ctx context.Context, location, name string) error {
	return c.internal.Unload(ctx, location, name)
}
