// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package group

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
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
		Tags:         tags.MapToProto(gp.Tags),
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
	var version *string = nil
	var locationName *string = nil
	if gp.LocationName != "" {
		locationName = &gp.LocationName
	}
	if gp.Status.Version.Number != "" {
		version = &gp.Status.Version.Number
	}
	return &cloud.Group{
		Name:     &gp.Name,
		Location: locationName,
		Version:  version,
		GroupProperties: &cloud.GroupProperties{
			Statuses: status.GetStatuses(gp.GetStatus()),
		},
		Tags: tags.ProtoToMap(gp.Tags),
	}
}
