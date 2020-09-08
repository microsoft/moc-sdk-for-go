// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package cluster

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdCluster(t *testing.T) {
	grp := &cloud.Cluster{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudCluster := getWssdCluster(grp)

	if *grp.ID != wssdcloudCluster.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudCluster.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getCluster(t *testing.T) {
	wssdcloudCluster := &wssdcloud.Cluster{
		Name: name,
		Id:   Id,
	}
	grp := getCluster(wssdcloudCluster)
	if *grp.ID != wssdcloudCluster.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudCluster.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
