// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

// this is only used in create/update
func convertFromWssdAvailabilitySet(s *compute.AvailabilitySet, group string) *wssdcloudcompute.AvailabilitySet {
	if s == nil {
		return nil
	}

	availabilitySet := &wssdcloudcompute.AvailabilitySet{
		Name:                     *s.Name,
		GroupName:                group,
		PlatformFaultDomainCount: *s.PlatformFaultDomainCount,
		Status:                   status.GetFromStatuses(s.Statuses),
		VirtualMachines:          convertFromWssdSubResources(s.VirtualMachines),
		Tags:                     convertFromWssdTags(s.Tags),
		LocationName:             *s.Location,
		Id:                       *s.ID,
	}
	return availabilitySet
}

func convertFromWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func convertFromWssdSubResources(resources []*compute.CloudSubResource) []*wssdcommon.CloudSubResource {
	ret := []*wssdcommon.CloudSubResource{}
	for _, res := range resources {
		ret = append(ret, convertFromWssdSubResource(res))
	}
	return ret
}

func convertFromWssdSubResource(s *compute.CloudSubResource) *wssdcommon.CloudSubResource {
	if s == nil {
		return nil
	}

	availabilitySet := &wssdcommon.CloudSubResource{
		Name:      *s.Name,
		GroupName: *s.GroupName,
	}
	return availabilitySet
}

// Convert from client model (rpc) to core model (compute)
func convertToWssdAvailabilitySet(s *wssdcloudcompute.AvailabilitySet) *compute.AvailabilitySet {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.AvailabilitySet{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		PlatformFaultDomainCount: &s.PlatformFaultDomainCount,
		Tags:                     convertToWssdTags(s.Tags),
		Version:                  &s.Status.Version.Number,
		VirtualMachines:          convertToWssdSubResources(s.VirtualMachines),
		Statuses:                 status.GetStatuses(s.Status),
	}
	return availabilitySet
}

func convertToWssdTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func convertToWssdSubResources(cs []*wssdcommon.CloudSubResource) []*compute.CloudSubResource {
	ret := []*compute.CloudSubResource{}
	for _, wssdvm := range cs {
		vm := convertToWssdSubResource(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func convertToWssdSubResource(s *wssdcommon.CloudSubResource) *compute.CloudSubResource {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.CloudSubResource{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return availabilitySet
}
