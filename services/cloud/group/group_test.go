// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package group

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdGroup(t *testing.T) {
	grp := &cloud.Group{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudGroup := getWssdGroup(grp)

	if *grp.ID != wssdcloudGroup.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudGroup.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getGroup(t *testing.T) {
	wssdcloudGroup := &wssdcloud.Group{
		Name: name,
		Id:   Id,
	}
	grp := getGroup(wssdcloudGroup)
	if *grp.ID != wssdcloudGroup.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudGroup.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
