// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package authentication

import (
	"context"

	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	//log "k8s.io/klog"
)

type client struct {
	wssdsecurity.AuthenticationAgentClient
}

// NewAuthenticationClient creates a client session with the backend wssd agent
func newAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetAuthenticationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Login
func (c *client) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	request := getAuthenticationRequest(identity)
	response, err := c.AuthenticationAgentClient.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return &response.Token, nil
}

func getAuthenticationRequest(identity *security.Identity) *wssdsecurity.AuthenticationRequest {
	certs := map[string]string{"": *identity.Certificate}
	request := &wssdsecurity.AuthenticationRequest{
		Identity: &wssdsecurity.Identity{
			Certificates: certs,
		},
	}
	if identity.Name != nil {
		request.Identity.Name = *identity.Name
	}

	return request
}
