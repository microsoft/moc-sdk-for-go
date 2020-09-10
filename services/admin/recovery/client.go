// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package recovery

import (
	"context"
	"github.com/microsoft/moc-sdk-for-go/services/admin/recovery/internal"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interfacetype Service interface {
type Service interface {
	Backup(context.Context, string) error
	Restore(context.Context, string) error
}

// Client structure
type RecoveryClient struct {
	internal Service
}

// NewClient method returns new client
func NewRecoveryClient(cloudFQDN string, authorizer auth.Authorizer) (*RecoveryClient, error) {
	c, err := internal.NewRecoveryClient(cloudFQDN, authorizer)
	return &RecoveryClient{c}, err
}

// Backupo
func (c *RecoveryClient) Backup(ctx context.Context, path string) error {
	return c.internal.Backup(ctx, path)
}

// Restore
func (c *RecoveryClient) Restore(ctx context.Context, path string) error {
	return c.internal.Restore(ctx, path)
}
