// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package keyvault

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/security"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]security.KeyVault, error)
	CreateOrUpdate(context.Context, string, string, *security.KeyVault) (*security.KeyVault, error)
	Delete(context.Context, string, string) error
}

// Client structure
type KeyVaultClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewKeyVaultClient(cloudFQDN string, authorizer auth.Authorizer) (*KeyVaultClient, error) {
	c, err := newKeyVaultClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &KeyVaultClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *KeyVaultClient) Get(ctx context.Context, group, name string) (*[]security.KeyVault, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *KeyVaultClient) CreateOrUpdate(ctx context.Context, group, name string, keyvault *security.KeyVault) (*security.KeyVault, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, keyvault)
}

// Delete methods invokes delete of the keyvault resource
func (c *KeyVaultClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
