// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package key

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
)

// Service interface
type Service interface {
	Get(context.Context, string, string, string) (*[]keyvault.Key, error)
	CreateOrUpdate(context.Context, string, string, string, *keyvault.Key) (*keyvault.Key, error)
	Delete(context.Context, string, string, string) error
	Encrypt(context.Context, string, string, string, *keyvault.KeyOperationsParameters) (*keyvault.KeyOperationResult, error)
	Decrypt(context.Context, string, string, string, *keyvault.KeyOperationsParameters) (*keyvault.KeyOperationResult, error)
	WrapKey(context.Context, string, string, string, *keyvault.KeyOperationsParameters) (*keyvault.KeyOperationResult, error)
	UnwrapKey(context.Context, string, string, string, *keyvault.KeyOperationsParameters) (*keyvault.KeyOperationResult, error)
}

// Client structure
type KeyClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewKeyClient(cloudFQDN string, authorizer auth.Authorizer) (*KeyClient, error) {
	c, err := newKeyClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &KeyClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *KeyClient) Get(ctx context.Context, group, name, vaultName string) (*[]keyvault.Key, error) {
	return c.internal.Get(ctx, group, name, vaultName)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *KeyClient) CreateOrUpdate(ctx context.Context, group, vaultName, name string,
	param *keyvault.Key) (*keyvault.Key, error) {
	return c.internal.CreateOrUpdate(ctx, group, vaultName, name, param)
}

// Delete methods invokes delete of the keyvault resource
func (c *KeyClient) Delete(ctx context.Context, group, name, vaultName string) error {
	return c.internal.Delete(ctx, group, name, vaultName)
}

// Encrypt methods invokes encrypt of the keyvault resource
func (c *KeyClient) Encrypt(ctx context.Context, group, vaultName, name string, parameters *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	return c.internal.Encrypt(ctx, group, name, vaultName, parameters)
}

// Decrypt methods invokes encrypt of the keyvault resource
func (c *KeyClient) Decrypt(ctx context.Context, group, vaultName, name string, parameters *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	return c.internal.Decrypt(ctx, group, name, vaultName, parameters)
}

// WrapKey
func (c *KeyClient) WrapKey(ctx context.Context, group, vaultName, name string, parameters *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	return c.internal.WrapKey(ctx, group, name, vaultName, parameters)
}

// UnwrapKey
func (c *KeyClient) UnwrapKey(ctx context.Context, group, vaultName, name string, parameters *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	return c.internal.UnwrapKey(ctx, group, name, vaultName, parameters)
}
