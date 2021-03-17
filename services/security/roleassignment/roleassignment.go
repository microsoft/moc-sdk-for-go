// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package roleassignment

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getScope(wssdscope *wssdcloudcommon.Scope) *security.Scope {
	if wssdscope == nil {
		return &security.Scope{}
	}
	return &security.Scope{
		Location: &wssdscope.Location,
		Group:    &wssdscope.ResourceGroup,
		Provider: security.GetProviderType(wssdscope.ProviderType),
		Resource: &wssdscope.Resource,
	}
}

func getRoleAssignment(wssdra *wssdcloudsecurity.RoleAssignment) *security.RoleAssignment {
	return &security.RoleAssignment{
		ID:      &wssdra.Id,
		Name:    &wssdra.Name,
		Version: &wssdra.Status.Version.Number,
		RoleAssignmentProperties: &security.RoleAssignmentProperties{
			RoleName:     &wssdra.RoleName,
			IdentityName: &wssdra.IdentityName,
			Scope:        getScope(wssdra.Scope),
		},
	}
}

func getWssdScope(scope *security.Scope) (*wssdcloudcommon.Scope, error) {
	wssdscope := &wssdcloudcommon.Scope{}

	if scope == nil {
		return wssdscope, nil
	}

	if scope.Location != nil {
		wssdscope.Location = *scope.Location
	}

	if scope.Group != nil {
		wssdscope.ResourceGroup = *scope.Group
	}

	providerType, err := security.GetWssdProviderType(scope.Provider)
	if err != nil {
		return nil, err
	}
	wssdscope.ProviderType = providerType

	if scope.Resource != nil {
		wssdscope.Resource = *scope.Resource
	}

	return wssdscope, nil
}

func getWssdRoleAssignment(ra *security.RoleAssignment) (*wssdcloudsecurity.RoleAssignment, error) {
	wssdra := &wssdcloudsecurity.RoleAssignment{}

	if ra.Version != nil {
		if wssdra.Status == nil {
			wssdra.Status = status.InitStatus()
		}
		wssdra.Status.Version.Number = *ra.Version
	}

	if ra.RoleAssignmentProperties != nil {
		if ra.RoleAssignmentProperties.IdentityName != nil {
			wssdra.IdentityName = *ra.IdentityName
		}

		if ra.RoleAssignmentProperties.RoleName != nil {
			wssdra.RoleName = *ra.RoleName
		}

		if ra.RoleAssignmentProperties.Scope != nil {
			scope, err := getWssdScope(ra.Scope)
			if err != nil {
				return nil, err
			}
			wssdra.Scope = scope
		}
	}

	return wssdra, nil
}
