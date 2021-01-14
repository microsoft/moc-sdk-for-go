// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package etcdserver

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/cloud/etcdcluster"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string, string) (*[]etcdcluster.EtcdServer, error)
	CreateOrUpdate(context.Context, string, string, *etcdcluster.EtcdServer) (*etcdcluster.EtcdServer, error)
	Delete(context.Context, string, string, string) error
}

// Client structure
type EtcdServerClient struct {
	//cloud.BaseClient
	internal Service
}

// NewClient method returns new client
func NewEtcdServerClient(cloudFQDN string, authorizer auth.Authorizer) (*EtcdServerClient, error) {
	c, err := newEtcdServerClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &EtcdServerClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *EtcdServerClient) Get(ctx context.Context, group, name, clusterName string) (*[]etcdcluster.EtcdServer, error) {
	return c.internal.Get(ctx, group, name, clusterName)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *EtcdServerClient) CreateOrUpdate(ctx context.Context, group, name string, server *etcdcluster.EtcdServer) (*etcdcluster.EtcdServer, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, server)
}

// Delete methods invokes delete of the etcdcluster resource
func (c *EtcdServerClient) Delete(ctx context.Context, group, name, clusterName string) error {
	return c.internal.Delete(ctx, group, name, clusterName)
}
