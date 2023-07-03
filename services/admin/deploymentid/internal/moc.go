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
	mocadmin.DeploymentIdAgentClient
}

// NewDeploymentIdClient - creates a client session with the backend moc agent
func NewDeploymentIdClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mocclient.GetDeploymentIdClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// GetDeploymentId
func (c *client) GetDeploymentId(ctx context.Context) (string, error) {
	request := getDeploymentIdRequest()
	response, err := c.DeploymentIdAgentClient.Get(ctx, request)
	if err != nil {
		return "", err
	}
	return response.DeploymentId, nil
}

func getDeploymentIdRequest() *mocadmin.DeploymentIdRequest {
	return &mocadmin.DeploymentIdRequest{}
}
