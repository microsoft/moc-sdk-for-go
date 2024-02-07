// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

// Convert from core ("github.com/microsoft/moc-sdk-for-go/services/compute") model to client(rpc) model
func getWssdAvailabilitySet(s *compute.AvailabilitySet, group string) *wssdcloudcompute.AvailabilitySet {
	if s == nil {
		return nil
	}

	availabilitySet := &wssdcloudcompute.AvailabilitySet{
		Name:                     *s.Name,
		GroupName:                group,
		PlatformFaultDomainCount: *s.PlatformFaultDomainCount,
	}
	return availabilitySet
}

func getWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func getwssdCloudSubResources(resources []*compute.CloudSubResource) []*wssdcommon.CloudSubResource {
	ret := []*wssdcommon.CloudSubResource{}
	for _, res := range resources {
		ret = append(ret, getWssdCloudSubResource(res))
	}
	return ret
}

func getWssdCloudSubResource(s *compute.CloudSubResource) *wssdcommon.CloudSubResource {
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
func getComputeAvailabilitySet(s *wssdcloudcompute.AvailabilitySet) *compute.AvailabilitySet {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.AvailabilitySet{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		PlatformFaultDomainCount: &s.PlatformFaultDomainCount,
		Tags:                     getComputeTags(s.Tags),
		Version:                  &s.Status.Version.Number,
		VirtualMachines:          getCloudSubResources(s.VirtualMachines),
	}
	return availabilitySet
}

func getComputeTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getCloudSubResources(cs []*wssdcommon.CloudSubResource) []*compute.CloudSubResource {
	ret := []*compute.CloudSubResource{}
	for _, wssdvm := range cs {
		vm := getCloudSubResource(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func getCloudSubResource(s *wssdcommon.CloudSubResource) *compute.CloudSubResource {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.CloudSubResource{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return availabilitySet
}
