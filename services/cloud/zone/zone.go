// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package zone

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// this is only used in create/update
func getRpcZone(s *cloud.Zone) (*wssdcloudcompute.Zone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Zone object Name is empty")
	}

	zone := &wssdcloudcompute.Zone{
		Name:                     *s.Name,
		LocationName:             *s.Location,
	}

	if s.Version != nil {  
		if zone.Status == nil {
			zone.Status = status.InitStatus()
		}
		zone.Status.Version.Number = *s.Version
	}

	if s.ZoneProperties.Nodes != nil {
		zone.Nodes = *s.ZoneProperties.Nodes
	}

	return zone, nil
}

// Convert from client model (rpc) to core model (compute)
func getWssdZone(s *wssdcloudcompute.Zone) (*cloud.Zone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	zone := &cloud.Zone{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		Version:                  &s.Status.Version.Number,
		ZoneProperties: &cloud.ZoneProperties{
			Statuses:                 status.GetStatuses(s.Status),
			Nodes:                    &s.Nodes,
		},
	}

	return zone, nil
}

