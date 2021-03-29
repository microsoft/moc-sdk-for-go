// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package roleassignment

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getScope(mocscope *wssdcloudcommon.Scope) *security.Scope {
	if mocscope == nil {
		return &security.Scope{}
	}
	return &security.Scope{
		Location: &mocscope.Location,
		Group:    &mocscope.ResourceGroup,
		Provider: security.GetProviderType(mocscope.ProviderType),
		Resource: &mocscope.Resource,
	}
}

func getRoleAssignment(mocra *wssdcloudsecurity.RoleAssignment) *security.RoleAssignment {
	return &security.RoleAssignment{
		ID:      &mocra.Id,
		Name:    &mocra.Name,
		Version: &mocra.Status.Version.Number,
		RoleAssignmentProperties: &security.RoleAssignmentProperties{
			RoleName:     &mocra.RoleName,
			IdentityName: &mocra.IdentityName,
			Scope:        getScope(mocra.Scope),
		},
	}
}

func getMocScope(scope *security.Scope) (*wssdcloudcommon.Scope, error) {
	mocscope := &wssdcloudcommon.Scope{}

	if scope == nil {
		return mocscope, nil
	}

	if scope.Location != nil {
		mocscope.Location = *scope.Location
	}

	if scope.Group != nil {
		mocscope.ResourceGroup = *scope.Group
	}

	providerType, err := security.GetMocProviderType(scope.Provider)
	if err != nil {
		return nil, err
	}
	mocscope.ProviderType = providerType

	if scope.Resource != nil {
		mocscope.Resource = *scope.Resource
	}

	return mocscope, nil
}

func getMocRoleAssignment(ra *security.RoleAssignment) (*wssdcloudsecurity.RoleAssignment, error) {
	mocra := &wssdcloudsecurity.RoleAssignment{}

	if ra.Name != nil {
		mocra.Name = *ra.Name
	}

	if ra.Version != nil {
		if mocra.Status == nil {
			mocra.Status = status.InitStatus()
		}
		mocra.Status.Version.Number = *ra.Version
	}

	if ra.RoleAssignmentProperties != nil {
		if ra.RoleAssignmentProperties.IdentityName != nil {
			mocra.IdentityName = *ra.IdentityName
		}

		if ra.RoleAssignmentProperties.RoleName != nil {
			mocra.RoleName = *ra.RoleName
		}

		if ra.RoleAssignmentProperties.Scope != nil {
			scope, err := getMocScope(ra.Scope)
			if err != nil {
				return nil, err
			}
			mocra.Scope = scope
		}
	}

	return mocra, nil
}
