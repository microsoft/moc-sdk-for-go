// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

// this is only used in create/update
func getRpcAvailabilityZone(s *compute.AvailabilityZone) (*wssdcloudcompute.AvailabilityZone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilityZone object Name is empty")
	}

	availabilityZone := &wssdcloudcompute.AvailabilityZone{
		Name:                     *s.Name,
		Status:                   status.GetFromStatuses(s.Statuses),
		VirtualMachines:          getRpcVirtualMachineReferences(s.VirtualMachines),
		Nodes:                    s.Nodes,
	}
	return availabilityZone, nil
}

func getRpcVirtualMachineReferences(resources []*compute.VirtualMachineReference) []*wssdcloudcompute.VirtualMachineRef {
	ret := []*wssdcloudcompute.VirtualMachineRef {}
	for _, res := range resources {
		ret = append(ret, getRpcVirtualMachineReference(res))
	}
	return ret
}

func getRpcVirtualMachineReference(s *compute.VirtualMachineReference) *wssdcloudcompute.VirtualMachineRef {
	if s == nil {
		return nil
	}

	vm := &wssdcloudcompute.VirtualMachineRef{
		Name:      *s.Name,
		GroupName: *s.GroupName,
	}
	return vm
}

// Convert from client model (rpc) to core model (compute)
func getWssdAvailabilityZone(s *wssdcloudcompute.AvailabilityZone) (*compute.AvailabilityZone, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Availability zone object is nil")
	}

	availabilityZone := &compute.AvailabilityZone{
		Name:                     &s.Name,
		ID:                       &s.Id,
		Location:                 &s.LocationName,
		Version:                  &s.Status.Version.Number,
		VirtualMachines:          getWssdVirtualMachineReferences(s.VirtualMachines),
		Statuses:                 status.GetStatuses(s.Status),
		Nodes:                    &s.Nodes,
	}
	return availabilityZone, nil
}

func getWssdVirtualMachineReferences(cs []*wssdcloudcompute.VirtualMachineRef) []*compute.VirtualMachineReference {
	ret := []*compute.VirtualMachineReference{}
	for _, wssdvm := range cs {
		vm := getWssdVirtualMachineReference(wssdvm)
		ret = append(ret, vm)
	}
	return ret
}

func getWssdVirtualMachineReference(s *wssdcloudcompute.VirtualMachineRef) *compute.VirtualMachineReference {
	if s == nil {
		return nil
	}

	vm := &compute.VirtualMachineReference{
		Name:      &s.Name,
		GroupName: &s.GroupName,
	}
	return vm
}