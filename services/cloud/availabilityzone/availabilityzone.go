// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// this is only used in create/update
func getRpcAvailabilityZone(s *cloud.AvailabilityZone) (*wssdcloudcompute.AvailabilityZone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilityZone object Name is empty")
	}

	availabilityZone := &wssdcloudcompute.AvailabilityZone{
		Name:                     *s.Name,
		LocationName:             *s.Location,
	}

	if s.Version != nil {  
		if availabilityZone.Status == nil {
			availabilityZone.Status = status.InitStatus()
		}
		availabilityZone.Status.Version.Number = *s.Version
	}

	if s.AvailabilityZoneProperties.Nodes != nil {
		availabilityZone.Nodes = *s.AvailabilityZoneProperties.Nodes
	}

	return availabilityZone, nil
}

// Convert from client model (rpc) to core model (compute)
func getWssdAvailabilityZone(s *wssdcloudcompute.AvailabilityZone) (*cloud.AvailabilityZone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	availabilityZone := &cloud.AvailabilityZone{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		Version:                  &s.Status.Version.Number,
		AvailabilityZoneProperties: &cloud.AvailabilityZoneProperties{
			Statuses:                 status.GetStatuses(s.Status),
			Nodes:                    &s.Nodes,
		},
	}

	return availabilityZone, nil
}

