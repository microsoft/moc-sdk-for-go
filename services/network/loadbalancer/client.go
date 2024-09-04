// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package loadbalancer

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.LoadBalancer, error)
	CreateOrUpdate(context.Context, string, string, *network.LoadBalancer) (*network.LoadBalancer, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, group string, loadBalancers []*network.LoadBalancer) (bool, error)
}

// LoadBalancerClient structure
type LoadBalancerClient struct {
	network.BaseClient
	internal Service
}

// NewLoadBalancerClient method returns new client
func NewLoadBalancerClient(cloudFQDN string, authorizer auth.Authorizer) (*LoadBalancerClient, error) {
	c, err := newLoadBalancerClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LoadBalancerClient) Get(ctx context.Context, group, name string) (*[]network.LoadBalancer, error) {
	return c.internal.Get(ctx, group, name)
}

// Ensure methods invokes create or update on the client
func (c *LoadBalancerClient) CreateOrUpdate(ctx context.Context, group, name string, lb *network.LoadBalancer) (*network.LoadBalancer, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, lb)
}

// Delete methods invokes delete of the network resource
func (c *LoadBalancerClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Prechecks whether the system is able to create specified loadBalancers.
// Returns true if it is possible; or false with reason in error message if not.
func (c *LoadBalancerClient) Precheck(ctx context.Context, group string, loadBalancers []*network.LoadBalancer) (bool, error) {
	return c.internal.Precheck(ctx, group, loadBalancers)
}
