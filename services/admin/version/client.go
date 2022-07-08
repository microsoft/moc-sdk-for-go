// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package version

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/admin/version/internal"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interfacetype Service interface {
type Service interface {
	GetVersion(context.Context) (string, string, error)
}

// Client structure
type VersionClient struct {
	internal Service
}

// NewClient method returns new client
func NewVersionClient(cloudFQDN string, authorizer auth.Authorizer) (*VersionClient, error) {
	c, err := internal.NewVersionClient(cloudFQDN, authorizer)
	return &VersionClient{c}, err
}

// GetVersion
func (c *VersionClient) GetVersion(ctx context.Context) (string, string, error) {
	return c.internal.GetVersion(ctx)
}
