// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package internal

import (
	"context"

	mocclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	mocadmin "github.com/microsoft/moc/rpc/common/admin"
)

type client struct {
	mocadmin.ValidationAgentClient
}

// NewValidationgingClient - creates a client session with the backend wssd agent
func NewValidationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mocclient.GetValidationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Validate(ctx context.Context) error {
	request := getValidationRequest()
	_, err := c.ValidationAgentClient.Invoke(ctx, request)
	return err
}

func getValidationRequest() *mocadmin.ValidationRequest {
	return &mocadmin.ValidationRequest{}
}
