// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	mocclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/rpc/common"
	mocadmin "github.com/microsoft/moc/rpc/common/admin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	mocadmin.HealthAgentClient
}

// NewHealthClient - creates a client session with the backend moc agent
func NewHealthClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mocclient.GetHealthClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) CheckHealth(ctx context.Context, timeoutSeconds uint32) error {
	request := mocadmin.HealthRequest{TimeoutSeconds: timeoutSeconds}
	_, err := c.HealthAgentClient.CheckHealth(ctx, &request)
	return err
}

// GetAgentInfo
func (c *client) GetAgentInfo(ctx context.Context) (*common.NodeInfo, error) {
	response, err := c.HealthAgentClient.GetAgentInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return &common.NodeInfo{}, err
	}
	return response.Node, nil
}

// GetDeploymentId
func (c *client) GetDeploymentId(ctx context.Context) (string, error) {
	response, err := c.HealthAgentClient.GetAgentInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return response.DeploymentId, nil
}
