// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package role

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getActions(mocactions []*wssdcloudsecurity.Action) ([]security.Action, error) {
	var actions []security.Action
	for _, mocaction := range mocactions {
		action := security.Action{}
		switch mocaction.Operation {
		case wssdcloudsecurity.GeneralAccessOperation_Read:
			action.Operation = security.ReadAccess
		case wssdcloudsecurity.GeneralAccessOperation_Write:
			action.Operation = security.WriteAccess
		case wssdcloudsecurity.GeneralAccessOperation_Delete:
			action.Operation = security.DeleteAccess
		case wssdcloudsecurity.GeneralAccessOperation_All:
			action.Operation = security.AllAccess
		default:
			return nil, errors.Wrapf(errors.InvalidInput, "Access: [%v]", mocaction.Operation)
		}

		action.Provider = security.GetProviderType(mocaction.ProviderType)
		actions = append(actions, action)
	}

	return actions, nil
}

func getPermissions(mocperms []*wssdcloudsecurity.Permission) (*[]security.RolePermission, error) {
	permissions := []security.RolePermission{}

	for _, perm := range mocperms {
		actions, err := getActions(perm.Actions)
		if err != nil {
			return nil, err
		}

		notActions, err := getActions(perm.NotActions)
		if err != nil {
			return nil, err
		}

		permission := security.RolePermission{}
		if actions != nil {
			permission.Actions = &actions
		}
		if notActions != nil {
			permission.NotActions = &notActions
		}
		permissions = append(permissions, permission)
	}

	return &permissions, nil
}

func getAssignableScopes(mocscopes []*wssdcloudcommon.Scope) (*[]security.Scope, error) {
	scopes := []security.Scope{}

	for _, mocscope := range mocscopes {
		scopes = append(scopes, security.Scope{
			Location: &mocscope.Location,
			Group:    &mocscope.ResourceGroup,
			Provider: security.GetProviderType(mocscope.ProviderType),
			Resource: &mocscope.Resource,
		})
	}

	return &scopes, nil
}

func getRole(role *wssdcloudsecurity.Role) (*security.Role, error) {
	permissions, err := getPermissions(role.Permissions)
	if err != nil {
		return nil, err
	}

	scopes, err := getAssignableScopes(role.AssignableScopes)
	if err != nil {
		return nil, err
	}

	return &security.Role{
		ID:      &role.Id,
		Name:    &role.Name,
		Version: &role.Status.Version.Number,
		RoleProperties: &security.RoleProperties{
			Statuses:         status.GetStatuses(role.GetStatus()),
			Permissions:      permissions,
			AssignableScopes: scopes,
		},
	}, nil
}

func getMocProviderAction(action *security.Action) (wssdcloudcommon.ProviderAccessOperation, error) {
	mocaction := &wssdcloudsecurity.Action{}

	if action == nil {
		return wssdcloudcommon.ProviderAccessOperation_Unspecified, nil
	}

	switch action.ProviderOperation {
	case security.VirtualMachine_StartAccess:
		mocaction.ProviderOperation = wssdcloudcommon.ProviderAccessOperation_VirtualMachine_Start
	case security.VirtualMachine_StopAccess:
		mocaction.ProviderOperation = wssdcloudcommon.ProviderAccessOperation_VirtualMachine_Stop
	case security.VirtualMachine_ResetAccess:
		mocaction.ProviderOperation = wssdcloudcommon.ProviderAccessOperation_VirtualMachine_Reset
	default:
		return wssdcloudcommon.ProviderAccessOperation_Unspecified, errors.Wrapf(errors.InvalidInput, "Access: [%v]", action.ProviderOperation)
	}

	return wssdcloudcommon.ProviderAccessOperation_Unspecified, nil
}

func getMocAction(action *security.Action) (*wssdcloudsecurity.Action, error) {
	mocaction := &wssdcloudsecurity.Action{}

	if action == nil {
		return mocaction, nil
	}

	switch action.Operation {
	case security.ReadAccess:
		mocaction.Operation = wssdcloudsecurity.GeneralAccessOperation_Read
	case security.WriteAccess:
		mocaction.Operation = wssdcloudsecurity.GeneralAccessOperation_Write
	case security.DeleteAccess:
		mocaction.Operation = wssdcloudsecurity.GeneralAccessOperation_Delete
	case security.AllAccess:
		mocaction.Operation = wssdcloudsecurity.GeneralAccessOperation_All
	case security.ProviderAction:
		mocaction.Operation = wssdcloudsecurity.GeneralAccessOperation_ProviderAction
		mocaction.ProviderOperation, _ = getMocProviderAction(action)

	default:
		return nil, errors.Wrapf(errors.InvalidInput, "Access: [%v]", action.Operation)
	}

	providerType, err := security.GetMocProviderType(action.Provider)
	if err != nil {
		return nil, err
	}
	mocaction.ProviderType = providerType

	return mocaction, nil
}

func getMocPermission(permission *security.RolePermission) (*wssdcloudsecurity.Permission, error) {
	mocperm := &wssdcloudsecurity.Permission{}

	if permission == nil {
		return mocperm, nil
	}

	if permission.Actions != nil {
		for _, action := range *permission.Actions {
			mocaction, err := getMocAction(&action)
			if err != nil {
				return nil, err
			}
			mocperm.Actions = append(mocperm.Actions, mocaction)
		}
	}

	if permission.NotActions != nil {
		for _, action := range *permission.NotActions {
			mocaction, err := getMocAction(&action)
			if err != nil {
				return nil, err
			}
			mocperm.NotActions = append(mocperm.NotActions, mocaction)
		}
	}

	return mocperm, nil
}

func getMocAssignableScope(scope *security.Scope) (*wssdcloudcommon.Scope, error) {
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

func getMocRole(role *security.Role) (*wssdcloudsecurity.Role, error) {
	if role.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Role name is missing")
	}

	mocrole := &wssdcloudsecurity.Role{
		Name: *role.Name,
	}

	if role.Version != nil {
		if mocrole.Status == nil {
			mocrole.Status = status.InitStatus()
		}
		mocrole.Status.Version.Number = *role.Version
	}

	if role.RoleProperties != nil {
		if role.RoleProperties.Permissions != nil {
			for _, permission := range *role.Permissions {
				mocperm, err := getMocPermission(&permission)
				if err != nil {
					return nil, err
				}
				mocrole.Permissions = append(mocrole.Permissions, mocperm)
			}
		}

		if role.RoleProperties.AssignableScopes != nil {
			for _, scope := range *role.AssignableScopes {
				mocscope, err := getMocAssignableScope(&scope)
				if err != nil {
					return nil, err
				}
				mocrole.AssignableScopes = append(mocrole.AssignableScopes, mocscope)
			}
		}
	}

	return mocrole, nil
}
