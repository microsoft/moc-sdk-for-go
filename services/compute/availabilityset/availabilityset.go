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
)

// this is only used in create/update
func getRpcAvailabilitySet(s *compute.AvailabilitySet, group string) (*wssdcloudcompute.AvailabilitySet, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability set object is nil")
	}

	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilitySet object Name is empty")
	}

	if s.PlatformFaultDomainCount == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilitySet object PlatformFaultDomainCount is empty")
	}

	availabilitySet := &wssdcloudcompute.AvailabilitySet{
		Name:                     *s.Name,
		GroupName:                group,
		PlatformFaultDomainCount: *s.PlatformFaultDomainCount,
		Status:                   status.GetFromStatuses(s.Statuses),
		VirtualMachines:          getRpcVirtualMachineReferences(s.VirtualMachines),
		Tags:                     getRpcWssdTags(s.Tags),
	}
	return availabilitySet, nil
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
		VirtualMachines:          getWssdVirtualMachineReferences(s.VirtualMachines),
		Statuses:                 status.GetStatuses(s.Status),
	}
	return availabilitySet, nil
}

func getWssdTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
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
