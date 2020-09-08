// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package kubernetes

import (
	"context"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]cloud.Kubernetes, error)
	CreateOrUpdate(context.Context, string, string, *cloud.Kubernetes) (*cloud.Kubernetes, error)
	Delete(context.Context, string, string) error
}

// Client structure
type KubernetesClient struct {
	internal Service
}

// NewClient method returns new client
func NewKubernetesClient(cloudFQDN string, authorizer auth.Authorizer) (*KubernetesClient, error) {
	c, err := newKubernetesClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &KubernetesClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *KubernetesClient) Get(ctx context.Context, group, name string) (*[]cloud.Kubernetes, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *KubernetesClient) CreateOrUpdate(ctx context.Context, group, name string, cloud *cloud.Kubernetes) (*cloud.Kubernetes, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, cloud)
}

// Delete methods invokes delete of the cloud resource
func (c *KubernetesClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
