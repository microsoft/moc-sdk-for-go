// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"testing"

	wssdcloud "github.com/microsoft/moc-proto/rpc/cloudagent/cloud"
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdNode(t *testing.T) {
	grp := &cloud.Node{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudNode := getWssdNode(grp)

	if *grp.ID != wssdcloudNode.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudNode.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getNode(t *testing.T) {
	wssdcloudNode := &wssdcloud.Node{
		Name: name,
		Id:   Id,
	}
	grp := getNode(wssdcloudNode)
	if *grp.ID != wssdcloudNode.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudNode.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
