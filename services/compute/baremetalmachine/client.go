// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalmachine

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.BareMetalMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.BareMetalMachine) (*compute.BareMetalMachine, error)
	Delete(context.Context, string, string) error
	Query(context.Context, string, string) (*[]compute.BareMetalMachine, error)
}

type BareMetalMachineClient struct {
	compute.BaseClient
	internal Service
}

func NewBareMetalMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*BareMetalMachineClient, error) {
	c, err := newBareMetalMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &BareMetalMachineClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *BareMetalMachineClient) Get(ctx context.Context, group, name string) (*[]compute.BareMetalMachine, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *BareMetalMachineClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.BareMetalMachine) (*compute.BareMetalMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *BareMetalMachineClient) Delete(ctx context.Context, group string, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Query method invokes the client Get method and uses the provided query to filter the returned results
func (c *BareMetalMachineClient) Query(ctx context.Context, group, query string) (*[]compute.BareMetalMachine, error) {
	return c.internal.Query(ctx, group, query)
}

// Get the Virtual Machine by querying for the specified computer name
func (c *BareMetalMachineClient) GetByComputerName(ctx context.Context, group string, computerName string) (*[]compute.BareMetalMachine, error) {
	query := fmt.Sprintf("[?baremetalmachineproperties.osprofile.computername=='%s']", computerName)

	bmms, err := c.Query(ctx, group, query)
	if err != nil {
		return nil, err
	}

	return bmms, nil
}
