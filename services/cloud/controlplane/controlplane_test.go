// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package controlplane

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdControlPlane(t *testing.T) {
	grp := &cloud.ControlPlaneInfo{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudControlPlane := getWssdControlPlane(grp)

	if *grp.ID != wssdcloudControlPlane.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudControlPlane.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getControlPlane(t *testing.T) {
	wssdcloudControlPlane := &wssdcloud.ControlPlane{
		Name: name,
		Id:   Id,
	}
	grp := getControlPlane(wssdcloudControlPlane)
	if *grp.ID != wssdcloudControlPlane.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudControlPlane.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
