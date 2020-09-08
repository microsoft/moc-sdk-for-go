// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package location

import (
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/status"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdLocation(lcn *cloud.Location) (*wssdcloud.Location, error) {

	if lcn.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	location := &wssdcloud.Location{
		Name: *lcn.Name,
	}

	if lcn.Version != nil {
		if location.Status == nil {
			location.Status = status.InitStatus()
		}
		location.Status.Version.Number = *lcn.Version
	}

	return location, nil
}

// Conversion functions from wssdcloud to cloud
func getLocation(lcn *wssdcloud.Location) *cloud.Location {
	return &cloud.Location{
		Name:    &lcn.Name,
		Version: &lcn.Status.Version.Number,
		LocationProperties: &cloud.LocationProperties{
			Statuses: status.GetStatuses(lcn.GetStatus()),
		},
	}
}
