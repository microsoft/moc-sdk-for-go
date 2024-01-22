// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineavailabilityset

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachineAvailabilitySet, error)
	GetVirtualMachines(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachineAvailabilitySet) (*compute.VirtualMachineAvailabilitySet, error)
	Delete(context.Context, string, string) error
}

type VirtualMachineAvailabilitySetClient struct {
	compute.BaseClient
	internal Service
}

func NewVirtualMachineAvailabilitySetClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineAvailabilitySetClient, error) {
	c, err := newAvailabilitySetClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineAvailabilitySetClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineAvailabilitySetClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineAvailabilitySet, error) {
	return c.internal.Get(ctx, group, name)
}

// Get methods invokes the client Get method
func (c *VirtualMachineAvailabilitySetClient) List(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.GetVirtualMachines(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineAvailabilitySetClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachineAvailabilitySet) (*compute.VirtualMachineAvailabilitySet, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineAvailabilitySetClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
