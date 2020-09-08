// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineimage

import (
	"context"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachineImage, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachineImage) (*compute.VirtualMachineImage, error)
	Delete(context.Context, string, string) error
}

// Client structure
type VirtualMachineImageClient struct {
	compute.BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualMachineImageClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineImageClient, error) {
	c, err := newVirtualMachineImageClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineImageClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineImageClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineImage, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineImageClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachineImage) (*compute.VirtualMachineImage, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineImageClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
