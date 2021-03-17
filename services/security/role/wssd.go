// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package role

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
	log "k8s.io/klog"
)

type client struct {
	wssdcloudsecurity.RoleAgentClient
}

// NewRoleClient - creates a client session with the backend wssdcloud agent
func newRoleClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetRoleClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]security.Role, error) {
	request, err := c.getRoleRequestByName(wssdcloudcommon.Operation_GET, name)
	if err != nil {
		return nil, err
	}

	response, err := c.RoleAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	roles, err := c.getRolesFromResponse(response)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, role *security.Role) (*security.Role, error) {
	err := c.validate(ctx, role)
	if err != nil {
		return nil, err
	}

	request, err := c.getRoleRequest(wssdcloudcommon.Operation_POST, name, role)
	if err != nil {
		return nil, err
	}

	response, err := c.RoleAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Role] Create failed with error %v", err)
		return nil, err
	}

	roles, err := c.getRolesFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*roles) == 0 {
		return nil, fmt.Errorf("[Role][Create] Unexpected error: Creating a role returned no result")
	}

	return &((*roles)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string) error {
	role, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*role) == 0 {
		return fmt.Errorf("Role [%s] not found", name)
	}

	request, err := c.getRoleRequest(wssdcloudcommon.Operation_DELETE, name, &(*role)[0])
	if err != nil {
		return err
	}
	_, err = c.RoleAgentClient.Invoke(ctx, request)
	return err
}

func (c *client) validate(ctx context.Context, role *security.Role) (err error) {
	if role == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing input")
		return
	}

	if role.Name == nil {
		return errors.Wrapf(errors.InvalidConfiguration, "Missing Name for Role")
	}

	if role.RoleProperties == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing RoleProperties")
		return
	}

	if role.RoleProperties.Permissions == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing RoleProperties.Permissions")
		return
	}

	return
}

func (c *client) getRolesFromResponse(response *wssdcloudsecurity.RoleResponse) (*[]security.Role, error) {
	roles := []security.Role{}
	for _, wssdrole := range response.GetRoles() {
		role, err := getRole(wssdrole)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return &roles, nil
}

func (c *client) getRoleRequestByName(opType wssdcloudcommon.Operation, name string) (*wssdcloudsecurity.RoleRequest, error) {
	role := security.Role{
		Name: &name,
	}
	return c.getRoleRequest(opType, name, &role)
}

func (c *client) getRoleRequest(opType wssdcloudcommon.Operation, name string, role *security.Role) (*wssdcloudsecurity.RoleRequest, error) {
	request := &wssdcloudsecurity.RoleRequest{
		OperationType: opType,
		Roles:         []*wssdcloudsecurity.Role{},
	}

	wssdrole := &wssdcloudsecurity.Role{
		Name: name,
	}

	var err error
	if role != nil {
		wssdrole, err = getWssdRole(role)
		if err != nil {
			return nil, err
		}
	}

	request.Roles = append(request.Roles, wssdrole)
	return request, nil
}
