// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package renew

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc-sdk-for-go/services/security/renew/moc"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Renew(context.Context, string, *security.CertificateRequest) (*security.Certificate, string, error)
	RenewConfig(context.Context, *auth.WssdConfig) (*auth.WssdConfig, bool, error)
}

// Client structure
type RenewClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewRenewClient(cloudFQDN string, authorizer auth.Authorizer) (*RenewClient, error) {
	c, err := moc.NewRenewClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &RenewClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *RenewClient) Renew(ctx context.Context, group string, csr *security.CertificateRequest) (*security.Certificate, string, error) {
	return c.internal.Renew(ctx, group, csr)
}

func (c *RenewClient) RenewConfig(ctx context.Context, config *auth.WssdConfig) (*auth.WssdConfig, bool, error) {
	return c.internal.RenewConfig(ctx, config)
}
