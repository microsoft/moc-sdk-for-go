// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

type avsetStore struct {
	avsets map[string]*wssdcloudcompute.AvailabilitySet
}

type mockClient struct {
	avsetStore
}

// newClient - creates a mockClient session with the backend wssdcloud agent
func newAvailabilitySetMockClient(subID string, authorizer auth.Authorizer) (*mockClient, error) {
	store := avsetStore{
		avsets: make(map[string]*wssdcloudcompute.AvailabilitySet),
	}
	return &mockClient{store}, nil
}

// Get
func (c *mockClient) Get(ctx context.Context, group, name string) (*[]compute.AvailabilitySet, error) {
	// if name == nil, return all avsets
	if name == "" {
		ret := []compute.AvailabilitySet{}
		for _, wssdavset := range c.avsets {
			avset := getComputeAvailabilitySet(wssdavset)
			ret = append(ret, *avset)
		}
		return &ret, nil
	}

	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		wssdavset := c.avsets[name]
		avset := getComputeAvailabilitySet(wssdavset)
		// if it does, return the value
		return &[]compute.AvailabilitySet{*avset}, nil
	}

	return nil, errors.NotFound
}

// CreateOrUpdate
func (c *mockClient) CreateOrUpdate(ctx context.Context, group, name string, as *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	wssdavset := getWssdAvailabilitySet(as, group)

	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, check that the platform fault domain count is the same
		if c.avsets[name].PlatformFaultDomainCount != wssdavset.PlatformFaultDomainCount {
			return nil, errors.Wrapf(errors.InvalidInput, "PlatformFaultDomainCount cannot be changed")
		}

		// if it does, update the value
		c.avsets[name] = wssdavset
		return as, nil
	}

	// if it doesn't, create it
	c.avsets[name] = wssdavset
	return as, nil
}

// Delete methods invokes create or update on the mockClient
func (c *mockClient) Delete(ctx context.Context, group, name string) error {
	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, check if it has any VM members
		if len(c.avsets[name].VirtualMachines) > 0 {
			return errors.Wrapf(errors.InUse, "AvailabilitySet %s has VM members, cannot delete an availability set with VM members", name)
		}

		// if it doesn't, delete it
		delete(c.avsets, name)
		return nil
	}

	return errors.NotFound
}
