// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package location

import (
	"testing"

	wssdcloud "github.com/microsoft/moc-proto/rpc/cloudagent/cloud"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdLocation(t *testing.T) {
	lcn := &cloud.Location{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudLocation := getWssdLocation(lcn)

	if *lcn.ID != wssdcloudLocation.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *lcn.Name != wssdcloudLocation.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getLocation(t *testing.T) {
	wssdcloudLocation := &wssdcloud.Location{
		Name: name,
		Id:   Id,
	}
	lcn := getLocation(wssdcloudLocation)
	if *lcn.ID != wssdcloudLocation.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *lcn.Name != wssdcloudLocation.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
