// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package roleassignment

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudsecurity.RoleAssignmentAgentClient
}

// NewRoleAssignmentClient - creates a client session with the backend wssdcloud agent
func newRoleAssignmentClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetRoleAssignmentClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get - Retrieve roles assigned to named identity that match the role assignment definitions
func (c *client) Get(ctx context.Context, inputRa *security.RoleAssignment) (*[]security.RoleAssignment, error) {
	request, err := c.getRoleAssignmentRequest(wssdcloudcommon.Operation_GET, inputRa)
	if err != nil {
		return nil, err
	}

	response, err := c.RoleAssignmentAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}

	ras, err := c.getRoleAssignmentFromResponse(response)
	if err != nil {
		return nil, err
	}

	return ras, nil
}

// Delete - Remove role assigned to named identity that match the role assignment definitions
func (c *client) Delete(ctx context.Context, inputRa *security.RoleAssignment) error {
	err := c.validateWithName(ctx, inputRa)
	if err != nil {
		return err
	}

	request, err := c.getRoleAssignmentRequest(wssdcloudcommon.Operation_DELETE, inputRa)
	if err != nil {
		return err
	}
	_, err = c.RoleAssignmentAgentClient.Invoke(ctx, request)
	services.HandleGRPCError(err)
	return err
}

// CreateOrUpdate - Assign roles to identity
func (c *client) CreateOrUpdate(ctx context.Context, inputRa *security.RoleAssignment) (*security.RoleAssignment, error) {
	err := c.validate(ctx, inputRa)
	if err != nil {
		return nil, err
	}

	request, err := c.getRoleAssignmentRequest(wssdcloudcommon.Operation_POST, inputRa)
	if err != nil {
		return nil, err
	}

	response, err := c.RoleAssignmentAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}

	ras, err := c.getRoleAssignmentFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*ras) == 0 {
		return nil, fmt.Errorf("[RoleAssignment][Create] Unexpected error: Creating a role assignment returned no result")
	}

	return &((*ras)[0]), err
}

func (c *client) validateWithName(ctx context.Context, ra *security.RoleAssignment) (err error) {
	if ra == nil || ra.Name == nil {
		return c.validate(ctx, ra)
	}
	return
}

func (c *client) validate(ctx context.Context, ra *security.RoleAssignment) (err error) {
	if ra == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing input")
		return
	}

	if ra.RoleAssignmentProperties == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing RoleAssignmentProperties")
		return
	}

	if ra.RoleAssignmentProperties.IdentityName == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing Identity name for role assignment")
		return
	}

	if ra.RoleAssignmentProperties.RoleName == nil {
		err = errors.Wrapf(errors.InvalidConfiguration, "Missing Role name for role assignment")
		return
	}

	return
}

func (c *client) getRoleAssignmentFromResponse(response *wssdcloudsecurity.RoleAssignmentResponse) (*[]security.RoleAssignment, error) {
	ras := []security.RoleAssignment{}
	for _, ra := range response.GetRoleAssignments() {
		ras = append(ras, *getRoleAssignment(ra))
	}

	return &ras, nil
}

func (c *client) getRoleAssignmentRequest(opType wssdcloudcommon.Operation, ra *security.RoleAssignment) (*wssdcloudsecurity.RoleAssignmentRequest, error) {
	request := &wssdcloudsecurity.RoleAssignmentRequest{
		OperationType:   opType,
		RoleAssignments: []*wssdcloudsecurity.RoleAssignment{},
	}

	wssdra, err := getMocRoleAssignment(ra)
	if err != nil {
		return nil, err
	}

	request.RoleAssignments = append(request.RoleAssignments, wssdra)
	return request, nil
}
