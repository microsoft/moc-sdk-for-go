// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

// this is only used in create/update
func getRpcAvailabilitySet(s *compute.AvailabilitySet, group string) (*wssdcloudcompute.AvailabilitySet, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability set object is nil")
	}

	if s.Name == nil || s.ID == nil || s.Location == nil || s.PlatformFaultDomainCount == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilitySet object contains nil fields on Name/Id/Location/PlatformdomainCount")
	}

	availabilitySet := &wssdcloudcompute.AvailabilitySet{
		Name:                     *s.Name,
		GroupName:                group,
		PlatformFaultDomainCount: *s.PlatformFaultDomainCount,
		Status:                   status.GetFromStatuses(s.Statuses),
		VirtualMachines:          getRpcSubResources(s.VirtualMachines),
		Tags:                     getRpcWssdTags(s.Tags),
		LocationName:             *s.Location,
		Id:                       *s.ID,
	}
	return availabilitySet, nil
}

func getRpcWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func getRpcSubResources(resources []*compute.CloudSubResource) []*wssdcommon.CloudSubResource {
	ret := []*wssdcommon.CloudSubResource{}
	for _, res := range resources {
		ret = append(ret, getRpcSubResource(res))
	}
	return ret
}

func getRpcSubResource(s *compute.CloudSubResource) *wssdcommon.CloudSubResource {
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
func getWssdAvailabilitySet(s *wssdcloudcompute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability set object is nil")
	}

	availabilitySet := &compute.AvailabilitySet{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		PlatformFaultDomainCount: &s.PlatformFaultDomainCount,
		Tags:                     getWssdTags(s.Tags),
		Version:                  &s.Status.Version.Number,
		VirtualMachines:          getWssdSubResources(s.VirtualMachines),
		Statuses:                 status.GetStatuses(s.Status),
	}
	return availabilitySet, nil
}

func getWssdTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdSubResources(cs []*wssdcommon.CloudSubResource) []*compute.CloudSubResource {
	ret := []*compute.CloudSubResource{}
	for _, wssdvm := range cs {
		vm := getWssdSubResource(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func getWssdSubResource(s *wssdcommon.CloudSubResource) *compute.CloudSubResource {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.CloudSubResource{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return availabilitySet
}
