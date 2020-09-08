// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package group

import (
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/status"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdGroup(gp *cloud.Group, location string) (*wssdcloud.Group, error) {

	if gp.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	group := &wssdcloud.Group{
		Name:         *gp.Name,
		LocationName: location,
	}

	if gp.Version != nil {
		if group.Status == nil {
			group.Status = status.InitStatus()
		}
		group.Status.Version.Number = *gp.Version
	}

	return group, nil
}

// Conversion functions from wssdcloud to cloud
func getGroup(gp *wssdcloud.Group) *cloud.Group {
	return &cloud.Group{
		Name:     &gp.Name,
		Location: &gp.LocationName,
		Version:  &gp.Status.Version.Number,
		GroupProperties: &cloud.GroupProperties{
			Statuses: status.GetStatuses(gp.GetStatus()),
		},
	}
}
