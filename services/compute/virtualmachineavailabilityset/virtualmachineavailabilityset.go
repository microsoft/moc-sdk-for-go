// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachineavailabilityset

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

func (c *client) getWssdAvailabilitySet(s *compute.VirtualMachineAvailabilitySet, group string) (*wssdcloudcompute.AvailabilitySet, error) {
	if s == nil {
		return nil, nil
	}

	availabilitySet := &wssdcloudcompute.AvailabilitySet{
		Name:                     *s.Name,
		Id:                       *s.ID,
		LocationName:             *s.Location,
		GroupName:                group,
		PlatformFaultDomainCount: *s.PlatformFaultDomainCount,
		Tags:                     getWssdTags(s.Tags),
	}
	return availabilitySet, nil
}

func (c *client) getComputeAvailabilitySet(s *wssdcloudcompute.AvailabilitySet) (*compute.VirtualMachineAvailabilitySet, error) {
	if s == nil {
		return nil, nil
	}

	availabilitySet := &compute.VirtualMachineAvailabilitySet{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		PlatformFaultDomainCount: &s.PlatformFaultDomainCount,
		Tags:                     getComputeTags(s.Tags),
		Version:                  &s.Status.Version.Number,
		VirtualMachines:          c.getCloudSubResources(s.VirtualMachines),
	}
	return availabilitySet, nil
}

func getComputeTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func (c *client) getWssdCloudSubResource(s *compute.CloudSubResource) (*wssdcommon.CloudSubResource, error) {
	if s == nil {
		return nil, nil
	}

	availabilitySet := &wssdcommon.CloudSubResource{
		Name:      *s.Name,
		GroupName: *s.GroupName,
	}
	return availabilitySet, nil
}

func (c *client) getCloudSubResources(cs []*wssdcommon.CloudSubResource) []*compute.CloudSubResource {
	ret := []*compute.CloudSubResource{}
	for _, wssdvm := range cs {
		vm := c.getCloudSubResource(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func (c *client) getCloudSubResource(s *wssdcommon.CloudSubResource) *compute.CloudSubResource {
	if s == nil {
		return nil
	}

	availabilitySet := &compute.CloudSubResource{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return availabilitySet
}
