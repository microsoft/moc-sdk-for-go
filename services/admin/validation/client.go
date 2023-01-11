// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package validation

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/admin/validation/internal"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interfacetype Service interface {
type Service interface {
	Validate(context.Context) error
}

// Client structure
type ValidationClient struct {
	internal Service
}

// NewClient method returns new client
func NewValidationClient(cloudFQDN string, authorizer auth.Authorizer) (*ValidationClient, error) {
	c, err := internal.NewValidationClient(cloudFQDN, authorizer)
	return &ValidationClient{c}, err
}

// gets a file from the corresponding node agent and writes it to filename
func (c *ValidationClient) Validate(ctx context.Context) error {
	return c.internal.Validate(ctx)
}
