// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package logging

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-sdk-for-go/services/admin/logging/internal"
)

// Service interfacetype Service interface {
type Service interface {
	GetLogFile(context.Context, string, string) error
}

// Client structure
type LoggingClient struct {
	internal Service
}

// NewClient method returns new client
func NewLoggingClient(cloudFQDN string, authorizer auth.Authorizer) (*LoggingClient, error) {
	c, err := internal.NewLoggingClient(cloudFQDN, authorizer)
	return &LoggingClient{c}, err
}

// gets a file from the corresponding node agent and writes it to filename
func (c *LoggingClient) GetLogFile(ctx context.Context, location string, filename string) error {
	return c.internal.GetLogFile(ctx, location, filename)
}
