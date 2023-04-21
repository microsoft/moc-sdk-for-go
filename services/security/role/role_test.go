// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package role

import (
	"fmt"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/security"
	"github.com/microsoft/moc/rpc/common"
)

var (
	name    = "test"
	Id      = "1234"
	version = "1"

	expectedMocRole = wssdcloud.Role{
		Name:             name,
		AssignableScopes: []*common.Scope{},
		Permissions: []*wssdcloud.Permission{
			{
				Actions: []*wssdcloud.Action{
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_All,
						ProviderType:     common.ProviderType_AnyProvider,
					},
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_Read,
						ProviderType:     common.ProviderType_VirtualMachine,
					},
				},
				NotActions: []*wssdcloud.Action{
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_Delete,
						ProviderType:     common.ProviderType_AnyProvider,
					},
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_Write,
						ProviderType:     common.ProviderType_Location,
					},
				},
			}, {
				Actions: []*wssdcloud.Action{
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_Read,
					},
					{
						GeneralOperation: wssdcloud.GeneralAccessOperation_Write,
					},
				},
			},
			{
				Actions: []*wssdcloud.Action{
					{
						ProviderType:      common.ProviderType_VirtualMachine,
						GeneralOperation:  wssdcloud.GeneralAccessOperation_ProviderAction,
						ProviderOperation: common.ProviderAccessOperation_VirtualMachine_Start,
					},
				},
			},
		},
		Status: &common.Status{
			Version: &common.Version{
				Number: version,
			},
		},
	}

	expectedRole = security.Role{
		Name:    &name,
		Version: &version,
		RoleProperties: &security.RoleProperties{
			Permissions: &[]security.RolePermission{
				{
					Actions: &[]security.Action{
						{
							Provider:         security.AnyProviderType,
							GeneralOperation: security.AllAccess,
						},
						{
							Provider:         security.VirtualMachineType,
							GeneralOperation: security.ReadAccess,
						},
					},
					NotActions: &[]security.Action{
						{
							Provider:         security.AnyProviderType,
							GeneralOperation: security.DeleteAccess,
						},
						{
							Provider:         security.LocationType,
							GeneralOperation: security.WriteAccess,
						},
					},
				},
				{
					Actions: &[]security.Action{
						{
							GeneralOperation: security.ReadAccess,
						},
						{
							GeneralOperation: security.WriteAccess,
						},
					},
				},
				{
					Actions: &[]security.Action{
						{
							Provider:          security.VirtualMachineType,
							GeneralOperation:  security.ProviderAction,
							ProviderOperation: security.VirtualMachine_StartAccess,
						},
					},
				},
			},
		},
	}
)

func Test_getMocRole(t *testing.T) {

	mocRole, err := getMocRole(&expectedRole)
	if err != nil {
		t.Errorf(err.Error())
	}

	if mocRole.Name != expectedMocRole.Name {
		t.Errorf("Role names don't match post conversion %v %v", mocRole.Name, expectedMocRole.Name)
	}
	if mocRole.Id != expectedMocRole.Id {
		t.Errorf("Role ids don't match post conversion %v %v", mocRole.Id, expectedMocRole.Id)
	}
	if mocRole.Status.Version.Number != expectedMocRole.Status.Version.Number {
		t.Errorf("Roles versions don't match post conversion %v %v", mocRole.Status.Version.Number, expectedMocRole.Status.Version.Number)
	}
	if len(mocRole.Permissions) != len(expectedMocRole.Permissions) {
		t.Errorf("Role permissions length don't match post conversion %v %v", len(mocRole.Permissions), len(expectedMocRole.Permissions))
	} else {
		for i, v := range expectedMocRole.Permissions {
			if len(mocRole.Permissions[i].Actions) != len(v.Actions) {
				t.Errorf("Role actions length don't match post conversion %v %v", len(mocRole.Permissions[i].Actions), len(v.Actions))
			} else {
				for j, w := range v.Actions {
					if mocRole.Permissions[i].Actions[j].GeneralOperation != w.GeneralOperation {
						t.Errorf("GeneralOperations don't match post conversion %v %v", mocRole.Permissions[i].Actions[j].GeneralOperation, w.GeneralOperation)
					}
					if mocRole.Permissions[i].Actions[j].ProviderType != w.ProviderType {
						t.Errorf("ProviderTypes don't match post conversion %v %v", mocRole.Permissions[i].Actions[j].ProviderType, w.ProviderType)
					}
					if mocRole.Permissions[i].Actions[j].ProviderOperation != w.ProviderOperation {
						t.Errorf("ProviderAccessOperations don't match post conversion %v %v", mocRole.Permissions[i].Actions[j].ProviderOperation, w.ProviderOperation)
					}
				}
			}
			if len(mocRole.Permissions[i].NotActions) != len(v.NotActions) {
				t.Errorf("Role not actions length don't match post conversion %v %v", len(mocRole.Permissions[i].NotActions), len(v.NotActions))
			} else {
				for j, w := range v.NotActions {
					if mocRole.Permissions[i].NotActions[j].GeneralOperation != w.GeneralOperation {
						t.Errorf("GeneralOperations don't match post conversion %v %v", mocRole.Permissions[i].NotActions[j].GeneralOperation, w.GeneralOperation)
					}
					if mocRole.Permissions[i].NotActions[j].ProviderType != w.ProviderType {
						t.Errorf("ProviderTypes don't match post conversion %v %v", mocRole.Permissions[i].NotActions[j].ProviderType, w.ProviderType)
					}
				}
			}
		}
	}
}

func Test_getRole(t *testing.T) {
	inputMocRole := expectedMocRole
	inputMocRole.Id = Id
	outRole := expectedRole
	outRole.ID = &Id

	role, err := getRole(&inputMocRole)
	if err != nil {
		t.Errorf(err.Error())
	}

	if *role.ID != *outRole.ID {
		t.Errorf("Role IDs don't match post conversion %v %v", *role.ID, *outRole.ID)
	}
	if *role.Name != *outRole.Name {
		t.Errorf("Role names don't match post conversion %v %v", *role.Name, *outRole.Name)
	}
	if *role.Version != *outRole.Version {
		t.Errorf("Roles versions don't match post conversion %v %v", *role.Version, *outRole.Version)
	}
	if len(*role.Permissions) != len(*outRole.Permissions) {
		t.Errorf("Role permissions length don't match post conversion %v %v", len(*role.Permissions), len(*outRole.Permissions))
	} else {
		perms := *role.Permissions
		for i, v := range perms {
			if perms[i].Actions == nil {
				if v.Actions != nil {
					t.Errorf("Role actions length don't match post conversion")
				}
			} else {
				actions := *perms[i].Actions
				if len(actions) != len(*v.Actions) {
					t.Errorf("Role actions length don't match post conversion %v %v", len(actions), len(*v.Actions))
				} else {
					for j, w := range *v.Actions {
						if actions[j].GeneralOperation != w.GeneralOperation {
							t.Errorf("GeneralOperations don't match post conversion %v %v", actions[j].GeneralOperation, w.GeneralOperation)
						}
						fmt.Println(actions[j].Provider)
						if actions[j].Provider != w.Provider {
							t.Errorf("ProviderTypes don't match post conversion %v %v", actions[j].Provider, w.Provider)
						}

						if actions[j].ProviderOperation != w.ProviderOperation {
							t.Errorf("ProviderAccessOperations don't match post conversion %v %v", actions[j].Provider, w.Provider)
						}
					}
				}
			}

			if perms[i].NotActions == nil {
				if v.NotActions != nil {
					t.Errorf("Role not actions length don't match post conversion")
				}
			} else {
				notactions := *perms[i].NotActions
				if *perms[i].NotActions != nil && len(notactions) != len(*v.NotActions) {
					t.Errorf("Role not actions length don't match post conversion %v %v", len(notactions), len(*v.NotActions))
				} else {
					for j, w := range *v.NotActions {

						if notactions[j].GeneralOperation != w.GeneralOperation {
							t.Errorf("GeneralOperations don't match post conversion %v %v", notactions[j].GeneralOperation, w.GeneralOperation)
						}
						if notactions[j].Provider != w.Provider {
							t.Errorf("ProviderTypes don't match post conversion %v %v", notactions[j].Provider, w.Provider)
						}
					}
				}
			}
		}
	}
}
