// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-proto/pkg/auth"
	wssdcloudsecurity "github.com/microsoft/moc-proto/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	log "k8s.io/klog"
)

type client struct {
	wssdcloudsecurity.IdentityAgentClient
}

// NewIdentityClientN- creates a client session with the backend wssdcloud agent
func newIdentityClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetIdentityClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.Identity, error) {

	request, err := getIdentityRequest(wssdcloudcommon.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.IdentityAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getIdentitysFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.Identity) (*security.Identity, error) {
	request, err := getIdentityRequest(wssdcloudcommon.Operation_POST, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.IdentityAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Identity] Create failed with error %v", err)
		return nil, err
	}

	cert := getIdentitysFromResponse(response)

	if len(*cert) == 0 {
		return nil, fmt.Errorf("[Identity][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*cert)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	id, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*id) == 0 {
		return fmt.Errorf("Identity [%s] not found", name)
	}

	request, err := getIdentityRequest(wssdcloudcommon.Operation_DELETE, name, &(*id)[0])
	if err != nil {
		return err
	}
	_, err = c.IdentityAgentClient.Invoke(ctx, request)
	return err
}

func getIdentitysFromResponse(response *wssdcloudsecurity.IdentityResponse) *[]security.Identity {
	certs := []security.Identity{}
	for _, identitys := range response.GetIdentitys() {
		certs = append(certs, *(getIdentity(identitys)))
	}

	return &certs
}

func getIdentityRequest(opType wssdcloudcommon.Operation, name string, cert *security.Identity) (*wssdcloudsecurity.IdentityRequest, error) {
	request := &wssdcloudsecurity.IdentityRequest{
		OperationType: opType,
		Identitys:     []*wssdcloudsecurity.Identity{},
	}
	wssdidentity := &wssdcloudsecurity.Identity{
		Name: name,
	}

	var err error
	if cert != nil {
		wssdidentity, err = getWssdIdentity(cert)
		if err != nil {
			return nil, err
		}
	}
	request.Identitys = append(request.Identitys, wssdidentity)
	return request, nil
}
