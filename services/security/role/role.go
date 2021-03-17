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

func getActions(wssdactions []*wssdcloudsecurity.Action) ([]security.Action, error) {
	var actions []security.Action
	for _, wssdaction := range wssdactions {
		action := security.Action{}
		switch wssdaction.Operation {
		case wssdcloudsecurity.AccessOperation_Read:
			action.Operation = security.ReadAccess
		case wssdcloudsecurity.AccessOperation_Write:
			action.Operation = security.WriteAccess
		case wssdcloudsecurity.AccessOperation_Delete:
			action.Operation = security.DeleteAccess
		case wssdcloudsecurity.AccessOperation_All:
			action.Operation = security.AllAccess
		default:
			return nil, errors.Wrapf(errors.InvalidInput, "Access: [%v]", wssdaction.Operation)
		}

		action.Provider = security.GetProviderType(wssdaction.ProviderType)
		actions = append(actions, action)
	}

	return actions, nil
}

func getPermissions(wssdperms []*wssdcloudsecurity.Permission) (*[]security.RolePermission, error) {
	permissions := []security.RolePermission{}

	for _, perm := range wssdperms {
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

func getAssignableScopes(wssdscopes []*wssdcloudcommon.Scope) (*[]security.Scope, error) {
	scopes := []security.Scope{}

	for _, wssdscope := range wssdscopes {
		scopes = append(scopes, security.Scope{
			Location: &wssdscope.Location,
			Group:    &wssdscope.ResourceGroup,
			Provider: security.GetProviderType(wssdscope.ProviderType),
			Resource: &wssdscope.Resource,
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

func getWssdAction(action *security.Action) (*wssdcloudsecurity.Action, error) {
	wssdaction := &wssdcloudsecurity.Action{}

	if action == nil {
		return wssdaction, nil
	}

	switch action.Operation {
	case security.ReadAccess:
		wssdaction.Operation = wssdcloudsecurity.AccessOperation_Read
	case security.WriteAccess:
		wssdaction.Operation = wssdcloudsecurity.AccessOperation_Write
	case security.DeleteAccess:
		wssdaction.Operation = wssdcloudsecurity.AccessOperation_Delete
	case security.AllAccess:
		wssdaction.Operation = wssdcloudsecurity.AccessOperation_All
	default:
		return nil, errors.Wrapf(errors.InvalidInput, "Access: [%v]", action.Operation)
	}

	providerType, err := security.GetWssdProviderType(action.Provider)
	if err != nil {
		return nil, err
	}
	wssdaction.ProviderType = providerType

	return wssdaction, nil
}

func getWssdPermission(permission *security.RolePermission) (*wssdcloudsecurity.Permission, error) {
	wssdperm := &wssdcloudsecurity.Permission{}

	if permission == nil {
		return wssdperm, nil
	}

	if permission.Actions != nil {
		for _, action := range *permission.Actions {
			wssdaction, err := getWssdAction(&action)
			if err != nil {
				return nil, err
			}
			wssdperm.Actions = append(wssdperm.Actions, wssdaction)
		}
	}

	if permission.NotActions != nil {
		for _, action := range *permission.NotActions {
			wssdaction, err := getWssdAction(&action)
			if err != nil {
				return nil, err
			}
			wssdperm.NotActions = append(wssdperm.NotActions, wssdaction)
		}
	}

	return wssdperm, nil
}

func getWssdAssignableScope(scope *security.Scope) (*wssdcloudcommon.Scope, error) {
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

func getWssdRole(role *security.Role) (*wssdcloudsecurity.Role, error) {
	if role.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Role name is missing")
	}

	wssdrole := &wssdcloudsecurity.Role{
		Name: *role.Name,
	}

	if role.Version != nil {
		if wssdrole.Status == nil {
			wssdrole.Status = status.InitStatus()
		}
		wssdrole.Status.Version.Number = *role.Version
	}

	if role.RoleProperties != nil {
		if role.RoleProperties.Permissions != nil {
			for _, permission := range *role.Permissions {
				wssdperm, err := getWssdPermission(&permission)
				if err != nil {
					return nil, err
				}
				wssdrole.Permissions = append(wssdrole.Permissions, wssdperm)
			}
		}

		if role.RoleProperties.AssignableScopes != nil {
			for _, scope := range *role.AssignableScopes {
				wssdscope, err := getWssdAssignableScope(&scope)
				if err != nil {
					return nil, err
				}
				wssdrole.AssignableScopes = append(wssdrole.AssignableScopes, wssdscope)
			}
		}
	}

	return wssdrole, nil
}
