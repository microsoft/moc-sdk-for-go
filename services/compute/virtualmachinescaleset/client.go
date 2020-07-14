// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachinescaleset

import (
	"context"

	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachineScaleSet, error)
	GetVirtualMachines(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error)
	Delete(context.Context, string, string) error
}

type VirtualMachineScaleSetClient struct {
	compute.BaseClient
	internal Service
}

func NewVirtualMachineScaleSetClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineScaleSetClient, error) {
	c, err := newVirtualMachineScaleSetClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineScaleSetClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineScaleSetClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineScaleSet, error) {
	return c.internal.Get(ctx, group, name)
}

// Get methods invokes the client Get method
func (c *VirtualMachineScaleSetClient) List(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.GetVirtualMachines(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineScaleSetClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineScaleSetClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
