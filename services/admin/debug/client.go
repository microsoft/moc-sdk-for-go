// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package debug

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/admin/debug/internal"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interfacetype Service interface {
type Service interface {
	Stacktrace(context.Context) (string, error)
}

// Client structure
type DebugClient struct {
	internal Service
}

// NewClient method returns new client
func NewDebugClient(cloudFQDN string, authorizer auth.Authorizer) (*DebugClient, error) {
	c, err := internal.NewDebugClient(cloudFQDN, authorizer)
	return &DebugClient{c}, err
}

// Stacktrace
func (c *DebugClient) Stacktrace(ctx context.Context) (string, error) {
	return c.internal.Stacktrace(ctx)
}
