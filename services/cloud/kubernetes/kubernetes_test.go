// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package kubernetes

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

var (
	name = "test"
	Id   = "1234"
)

func Test_getWssdKubernetes(t *testing.T) {
	grp := &cloud.Kubernetes{
		Name: &name,
		ID:   &Id,
	}
	wssdcloudKubernetes := getWssdKubernetes(grp)

	if *grp.ID != wssdcloudKubernetes.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudKubernetes.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
func Test_getKubernetes(t *testing.T) {
	wssdcloudKubernetes := &wssdcloud.Kubernetes{
		Name: name,
		Id:   Id,
	}
	grp := getKubernetes(wssdcloudKubernetes)
	if *grp.ID != wssdcloudKubernetes.Id {
		t.Errorf("ID doesnt match post conversion")
	}
	if *grp.Name != wssdcloudKubernetes.Name {
		t.Errorf("Name doesnt match post conversion")
	}
}
