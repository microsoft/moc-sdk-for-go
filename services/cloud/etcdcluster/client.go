// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdcluster

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]cloud.EtcdCluster, error)
	CreateOrUpdate(context.Context, string, string, *cloud.EtcdCluster) (*cloud.EtcdCluster, error)
	Delete(context.Context, string, string) error
}

// Client structure
type EtcdClusterClient struct {
	//cloud.BaseClient
	internal Service
}

// NewClient method returns new client
func NewEtcdClusterClient(cloudFQDN string, authorizer auth.Authorizer) (*EtcdClusterClient, error) {
	c, err := newEtcdClusterClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &EtcdClusterClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *EtcdClusterClient) Get(ctx context.Context, group, name string) (*[]cloud.EtcdCluster, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *EtcdClusterClient) CreateOrUpdate(ctx context.Context, group, name string, etcdcluster *cloud.EtcdCluster) (*cloud.EtcdCluster, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, etcdcluster)
}

// Delete methods invokes delete of the etcdcluster resource
func (c *EtcdClusterClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
