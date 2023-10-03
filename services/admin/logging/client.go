// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package logging

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/admin/logging/internal"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interfacetype Service interface {
type Service interface {
	GetLogFile(context.Context, string, string) error
	SetVerbosityLevel(context.Context, int32, bool) error
	GetVerbosityLevel(context.Context) (string, error)
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

func (c *LoggingClient) SetVerbosityLevel(ctx context.Context, verbositylevel int32, include_nodeagents bool) error {
	return c.internal.SetVerbosityLevel(ctx, verbositylevel, include_nodeagents)
}

func (c *LoggingClient) GetVerbosityLevel(ctx context.Context) (string, error) {
	return c.internal.GetVerbosityLevel(ctx)
}
