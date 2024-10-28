// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package placementgroup

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
)

// this is only used in create/update
func getRpcPlacementGroup(s *compute.PlacementGroup, group string) (*wssdcloudcompute.PlacementGroup, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Placement group object is nil")
	}

	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "PlacementGroup object Name is empty")
	}

	pgType := wssdcloudcompute.PlacementGroupType_Affinity
	if s.Type == compute.Affinity {
		pgType = wssdcloudcompute.PlacementGroupType_Affinity
	} else if s.Type == compute.AntiAffinity {
		pgType = wssdcloudcompute.PlacementGroupType_AntiAffinity
	} else if s.Type == compute.StrictAntiAffinity {
		pgType = wssdcloudcompute.PlacementGroupType_StrictAntiAffinity
	}

	pgScope := wssdcloudcompute.PlacementGroupScope_Server
	if s.Scope == compute.ZoneScope {
		pgScope = wssdcloudcompute.PlacementGroupScope_Zone
	}

	placementGroup := &wssdcloudcompute.PlacementGroup{
		Name:            *s.Name,
		GroupName:       group,
		Status:          status.GetFromStatuses(s.Statuses),
		VirtualMachines: getRpcVirtualMachineReferences(s.VirtualMachines),
		Type:            pgType,
		Scope:           pgScope,
	}

	if s.PlacementGroupProperties != nil {
		if s.PlacementGroupProperties.Zones != nil {
			placementGroup.Zones = &wssdcloudproto.ZoneConfiguration{
				Zones:           []*wssdcloudproto.ZoneReference{},
				StrictPlacement: s.PlacementGroupProperties.StrictPlacement,
			}

			for _, zn := range *s.PlacementGroupProperties.Zones {
				rpcZoneRef, err := getRpcZoneReference(&zn)
				if err != nil {
					return nil, err
				}
				placementGroup.Zones.Zones = append(placementGroup.Zones.Zones, rpcZoneRef)
			}
		}
	}

	return placementGroup, nil
}

func getRpcZoneReference(s *string) (*wssdcloudproto.ZoneReference, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Zone Name is missing")
	}

	return &wssdcloudproto.ZoneReference{
		Name: *s,
	}, nil
}

func getRpcWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func getRpcVirtualMachineReferences(resources []*compute.VirtualMachineReference) []*wssdcloudcompute.VirtualMachineReference {
	ret := []*wssdcloudcompute.VirtualMachineReference{}
	for _, res := range resources {
		ret = append(ret, getRpcVirtualMachineReference(res))
	}
	return ret
}

func getRpcVirtualMachineReference(s *compute.VirtualMachineReference) *wssdcloudcompute.VirtualMachineReference {
	if s == nil {
		return nil
	}

	vm := &wssdcloudcompute.VirtualMachineReference{
		Name:      *s.Name,
		GroupName: *s.GroupName,
	}
	return vm
}

// Convert from client model (rpc) to core model (compute)
func getWssdPlacementGroup(s *wssdcloudcompute.PlacementGroup) (*compute.PlacementGroup, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Placement group object is nil")
	}

	pgZone := []string{}

	if s.Zones != nil && s.Zones.Zones != nil {
		for _, zn := range s.Zones.Zones {
			pgZone = append(pgZone, zn.Name)
		}
    }

	pgScope := compute.ServerScope
	if s.Scope == wssdcloudcompute.PlacementGroupScope_Zone {
		pgScope = compute.ZoneScope
	}

	placementGroup := &compute.PlacementGroup{
		Name:     &s.Name,
		ID:       &s.Id,
		Location: &s.LocationName,
		Version:  &s.Status.Version.Number,
		PlacementGroupProperties: &compute.PlacementGroupProperties{
			VirtualMachines: getWssdVirtualMachineReferences(s.VirtualMachines),
			Statuses:        status.GetStatuses(s.Status),
			Zones:           &pgZone,
			Scope:           pgScope,
			StrictPlacement: s.Zones.StrictPlacement,
		},
	}

	return placementGroup, nil
}

func getWssdVirtualMachineReferences(cs []*wssdcloudcompute.VirtualMachineReference) []*compute.VirtualMachineReference {
	ret := []*compute.VirtualMachineReference{}
	for _, wssdvm := range cs {
		vm := getWssdVirtualMachineReference(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func getWssdVirtualMachineReference(s *wssdcloudcompute.VirtualMachineReference) *compute.VirtualMachineReference {
	if s == nil {
		return nil
	}

	vm := &compute.VirtualMachineReference{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return vm
}
