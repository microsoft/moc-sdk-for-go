// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package certificate

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]security.Certificate, error)
	CreateOrUpdate(context.Context, string, string, *security.Certificate) (*security.Certificate, error)
	Delete(context.Context, string, string) error
	Sign(context.Context, string, string, *security.CertificateRequest) (*security.Certificate, string, error)
	Renew(context.Context, string, string, *security.CertificateRequest) (*security.Certificate, string, error)
	Precheck(ctx context.Context, certificates []*security.Certificate) (bool, error)
}

// Client structure
type CertificateClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewCertificateClient(cloudFQDN string, authorizer auth.Authorizer) (*CertificateClient, error) {
	c, err := newCertificateClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &CertificateClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *CertificateClient) Get(ctx context.Context, group, name string) (*[]security.Certificate, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *CertificateClient) CreateOrUpdate(ctx context.Context, group, name string, Certificate *security.Certificate) (*security.Certificate, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, Certificate)
}

// Delete methods invokes delete of the Certificate resource
func (c *CertificateClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Sign methods invokes sign to create a CA-Signed Certificate
func (c *CertificateClient) Sign(ctx context.Context, group, name string, csr *security.CertificateRequest) (*security.Certificate, string, error) {
	return c.internal.Sign(ctx, group, name, csr)
}

// Renew methods invokes renew to renew signed-certificate
func (c *CertificateClient) Renew(ctx context.Context, group, name string, csr *security.CertificateRequest) (*security.Certificate, string, error) {
	return c.internal.Renew(ctx, group, name, csr)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *CertificateClient) Precheck(ctx context.Context, certificates []*security.Certificate) (bool, error) {
	return c.internal.Precheck(ctx, certificates)
}
