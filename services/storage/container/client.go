// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package container

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/storage"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]storage.Container, error)
	CreateOrUpdate(context.Context, string, string, *storage.Container) (*storage.Container, error)
	Delete(context.Context, string, string) error
}

// Client structure
type ContainerClient struct {
	storage.BaseClient
	internal Service
}

// NewClient method returns new client
func NewContainerClient(cloudFQDN string, authorizer auth.Authorizer) (*ContainerClient, error) {
	c, err := newContainerClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &ContainerClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *ContainerClient) Get(ctx context.Context, location, name string) (*[]storage.Container, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *ContainerClient) CreateOrUpdate(ctx context.Context, location, name string, storage *storage.Container) (*storage.Container, error) {
	return c.internal.CreateOrUpdate(ctx, location, name, storage)
}

// Delete methods invokes delete of the storage resource
func (c *ContainerClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}
