// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	mocclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	mocadmin "github.com/microsoft/moc/rpc/common/admin"
)

type client struct {
	mocadmin.VersionAgentClient
}

// NewVersionClient - creates a client session with the backend moc agent
func NewVersionClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mocclient.GetVersionClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// GetVersion
func (c *client) GetVersion(ctx context.Context) (string, string, error) {
	request := getVersionRequest(mocadmin.VersionOperation_VERSION)
	response, err := c.VersionAgentClient.Invoke(ctx, request)
	if err != nil {
		return "", "", err
	}
	return response.version, response.mocversion, nil
}

func getVersionRequest() *mocadmin.VersionRequest {
	return &mocadmin.VersionRequest{}
}
