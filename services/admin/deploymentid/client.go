// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package deploymentid

import (
	"context"

	"github.com/microsoft/moc-sdk-for-go/services/admin/deploymentid/internal"
	"github.com/microsoft/moc/pkg/auth"
)

var (
	deploymentId = ""
)

// Service interfacetype Service interface {
type Service interface {
	GetDeploymentId(context.Context) (string, error)
}

// Client structure
type DeploymentIdClient struct {
	internal Service
}

// NewClient method returns new client
func NewDeploymentIdClient(cloudFQDN string, authorizer auth.Authorizer) (*DeploymentIdClient, error) {
	c, err := internal.NewDeploymentIdClient(cloudFQDN, authorizer)
	return &DeploymentIdClient{c}, err
}

// GetDeploymentId
func (c *DeploymentIdClient) GetDeploymentId(ctx context.Context) (string, error) {
	//if deploymentId is cached, directly return it
	if len(deploymentId) != 0 {
		return deploymentId, nil
	}
	id, err := c.internal.GetDeploymentId(ctx)
	//if met error, return empty string
	if err != nil {
		return "", err
	}
	deploymentId = id
	return deploymentId, nil
}
