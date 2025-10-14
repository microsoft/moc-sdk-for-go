// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package roleassignment

import (
	"fmt"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/security"
	"github.com/microsoft/moc/rpc/common"
)

var (
	name            = "test"
	Id              = "1234"
	roleName        = "testRole"
	identityName    = "testIdentity"
	locationName    = "testLocation"
	groupName       = "testGroup"
	providerName    = security.VirtualMachineType
	providerMocName = common.ProviderType_VirtualMachine
	resourceName    = "testResource"
	version         = "1"

	expectedMocRA = wssdcloud.RoleAssignment{
		RoleName:     roleName,
		IdentityName: identityName,
		Scope: &common.Scope{
			Location:      locationName,
			ResourceGroup: groupName,
			ProviderType:  providerMocName,
			Resource:      resourceName,
		},
		Status: &common.Status{
			Version: &common.Version{
				Number: version,
			},
		},
	}

	expectedRA = security.RoleAssignment{
		Version: &version,
		RoleAssignmentProperties: &security.RoleAssignmentProperties{
			RoleName:     &roleName,
			IdentityName: &identityName,
			Scope: &security.Scope{
				Location: &locationName,
				Group:    &groupName,
				Provider: providerName,
				Resource: &resourceName,
			},
		},
	}
)

func Test_getMocRoleAssignmentNoName(t *testing.T) {
	mocRA, err := getMocRoleAssignment(&expectedRA)
	if err != nil {
		t.Errorf(err.Error(), "error")
	}
	if err := compareMocRas(mocRA, &expectedMocRA); err != nil {
		t.Errorf(err.Error(), "error")
	}
}

func Test_getMocRoleAssignment(t *testing.T) {
	inputRA := expectedRA
	inputRA.Name = &name
	outMocRA := expectedMocRA
	outMocRA.Name = name

	mocRA, err := getMocRoleAssignment(&inputRA)
	if err != nil {
		t.Errorf(err.Error(), "error")
	}

	if err := compareMocRas(mocRA, &outMocRA); err != nil {
		t.Errorf(err.Error(), "error")
	}
}

func Test_getRoleAssignmentNoName(t *testing.T) {
	inputMocRA := expectedMocRA
	inputMocRA.Id = Id
	outRA := expectedRA
	outRA.ID = &Id
	emptyName := ""
	outRA.Name = &emptyName

	ra := getRoleAssignment(&inputMocRA)
	if err := compareRas(ra, &outRA); err != nil {
		t.Errorf(err.Error(), "error")
	}
}

func Test_getRoleAssignment(t *testing.T) {
	inputMocRA := expectedMocRA
	inputMocRA.Name = name
	inputMocRA.Id = Id
	outRA := expectedRA
	outRA.Name = &name
	outRA.ID = &Id

	ra := getRoleAssignment(&inputMocRA)
	if err := compareRas(ra, &outRA); err != nil {
		t.Errorf(err.Error(), "error")
	}
}

func compareMocRas(mocRA, mocRA2 *wssdcloud.RoleAssignment) error {
	if mocRA.Name != mocRA2.Name {
		return fmt.Errorf("Role assignment names don't match post conversion %v %v", mocRA.Name, mocRA2.Name)
	}
	if mocRA.Id != mocRA2.Id {
		return fmt.Errorf("Role assignment ids don't match post conversion %v %v", mocRA.Id, mocRA2.Id)
	}
	if mocRA.RoleName != mocRA2.RoleName {
		return fmt.Errorf("Role assignment role names don't match post conversion %v %v", mocRA.RoleName, mocRA2.RoleName)
	}
	if mocRA.IdentityName != mocRA2.IdentityName {
		return fmt.Errorf("Role assignment identity names don't match post conversion %v %v", mocRA.IdentityName, mocRA2.IdentityName)
	}
	if mocRA.Status.Version.Number != mocRA2.Status.Version.Number {
		return fmt.Errorf("Role assignment versions don't match post conversion %v %v", mocRA.Status.Version.Number, mocRA2.Status.Version.Number)
	}
	if mocRA.Scope == nil {
		if mocRA2.Scope != nil {
			return fmt.Errorf("Role assignment scopes don't match post conversion %v %v", mocRA.Scope, mocRA2.Scope)
		}
	} else {
		if mocRA.Scope.Location != mocRA2.Scope.Location {
			return fmt.Errorf("Role assignment locations don't match post conversion %v %v", mocRA.Scope.Location, mocRA2.Scope.Location)
		}
		if mocRA.Scope.ResourceGroup != mocRA2.Scope.ResourceGroup {
			return fmt.Errorf("Role assignment resource groups don't match post conversion %v %v", mocRA.Scope.ResourceGroup, mocRA2.Scope.ResourceGroup)
		}
		if mocRA.Scope.ProviderType != mocRA2.Scope.ProviderType {
			return fmt.Errorf("Role assignment  provider types don't match post conversion %v %v", mocRA.Scope.ProviderType, mocRA2.Scope.ProviderType)
		}
		if mocRA.Scope.Resource != mocRA2.Scope.Resource {
			return fmt.Errorf("Role assignment resource names don't match post conversion %v %v", mocRA.Scope.Resource, mocRA2.Scope.Resource)
		}
	}
	return nil
}

func compareRas(ra, ra2 *security.RoleAssignment) error {
	if ra.Name == nil {
		if ra2.Name != nil {
			return fmt.Errorf("Role assignment names don't match post conversion")
		}
	} else {
		if *ra.Name != *ra2.Name {
			return fmt.Errorf("Role assignment names don't match post conversion %v %v", *ra.Name, *ra2.Name)
		}
	}
	if *ra.ID != *ra2.ID {
		return fmt.Errorf("Role assignment IDs don't match post conversion %v %v", *ra.ID, *ra2.ID)
	}
	if *ra.Version != *ra2.Version {
		return fmt.Errorf("Role assignment versions don't match post conversion %v %v", *ra.Version, *ra2.Version)
	}

	if ra.RoleAssignmentProperties == nil {
		if ra2.RoleAssignmentProperties != nil {
			return fmt.Errorf("Role assignment properties don't match post conversion")
		}
	} else {
		if *ra.IdentityName != *ra2.IdentityName {
			return fmt.Errorf("Role assignment identity names don't match post conversion %v %v", *ra.IdentityName, *ra2.IdentityName)
		}
		if *ra.RoleName != *ra2.RoleName {
			return fmt.Errorf("Role assignment role names don't match post conversion %v %v", *ra.RoleName, *ra2.Name)
		}
		if ra.Scope == nil {
			if ra2.Scope != nil {
				return fmt.Errorf("Role assignment Scopes don't match post conversion")
			}
		} else {
			if *ra.Scope.Location != *ra2.Scope.Location {
				return fmt.Errorf("Role assignment Locations don't match post conversion %v %v", *ra.Scope.Location, *ra2.Scope.Location)
			}
			if *ra.Scope.Group != *ra2.Scope.Group {
				return fmt.Errorf("Role assignment Groups don't match post conversion %v %v", *ra.Scope.Group, *ra2.Scope.Group)
			}
			if ra.Scope.Provider != ra2.Scope.Provider {
				return fmt.Errorf("Role assignment Providers don't match post conversion %v %v,\n%v | %v", ra.Scope.Provider, ra2.Scope.Provider, *ra.Scope, *ra2.Scope)
			}
			if *ra.Scope.Resource != *ra2.Scope.Resource {
				return fmt.Errorf("Role assignment Resource names don't match post conversion %v %v", *ra.Scope.Resource, *ra2.Scope.Resource)
			}
		}
	}
	return nil
}
