// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package authentication

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Login(context.Context, string, *security.Identity) (*string, error)
	LoginWithConfig(context.Context, string, auth.LoginConfig, bool) (*auth.WssdConfig, error)
}

// Client structure
type AuthenticationClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewAuthenticationClient(cloudFQDN string, authorizer auth.Authorizer) (*AuthenticationClient, error) {
	c, err := newAuthenticationClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AuthenticationClient{internal: c}, nil
}

// NewClient method returns new client based on the authentication mode
func NewAuthenticationClientAuthMode(cloudFQDN string, loginconfig auth.LoginConfig) (*AuthenticationClient, error) {
	authorizer, err := auth.NewAuthorizerForAuth(loginconfig.Token, loginconfig.Certificate, cloudFQDN)
	if err != nil {
		return nil, err
	}

	c, err := newAuthenticationClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}
	return &AuthenticationClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *AuthenticationClient) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	return c.internal.Login(ctx, group, identity)
}

// Get methods invokes the client Get method
func (c *AuthenticationClient) LoginWithConfig(ctx context.Context, group string, loginconfig auth.LoginConfig, enableRenewRoutine bool) (*auth.WssdConfig, error) {
	return c.internal.LoginWithConfig(ctx, group, loginconfig, enableRenewRoutine)
}
