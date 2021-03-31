// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]security.Identity, error)
	CreateOrUpdate(context.Context, string, string, *security.Identity) (*security.Identity, error)
	Delete(context.Context, string, string) error
	Revoke(context.Context, string, string) (*security.Identity, error)
	Rotate(context.Context, string, string) (*security.Identity, error)
	CreateCertificate(context.Context, string, string, []*security.CertificateRequest) ([]*security.Certificate, string, error)
	RenewCertificate(context.Context, string, string, []*security.CertificateRequest) ([]*security.Certificate, string, error)
}

// Client structure
type IdentityClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewIdentityClient(cloudFQDN string, authorizer auth.Authorizer) (*IdentityClient, error) {
	c, err := newIdentityClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &IdentityClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *IdentityClient) Get(ctx context.Context, group, name string) (*[]security.Identity, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *IdentityClient) CreateOrUpdate(ctx context.Context, group, name string, identity *security.Identity) (*security.Identity, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, identity)
}

// Delete methods invokes delete of the Identity resource
func (c *IdentityClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Revoke methods invokes revokes an identity
func (c *IdentityClient) Revoke(ctx context.Context, group, name string) (*security.Identity, error) {
	return c.internal.Revoke(ctx, group, name)
}

// Rotate methods rotates identity token
func (c *IdentityClient) Rotate(ctx context.Context, group, name string) (*security.Identity, error) {
	return c.internal.Rotate(ctx, group, name)
}

// CreateCertificate methods invokes creates client certificate for the identity
func (c *IdentityClient) CreateCertificate(ctx context.Context, group, name string, csr []*security.CertificateRequest) ([]*security.Certificate, string, error) {
	return c.internal.CreateCertificate(ctx, group, name, csr)
}

// RenewCertificate methods invokes renew client certificate for the identity
func (c *IdentityClient) RenewCertificate(ctx context.Context, group, name string, csr []*security.CertificateRequest) ([]*security.Certificate, string, error) {
	return c.internal.RenewCertificate(ctx, group, name, csr)
}
