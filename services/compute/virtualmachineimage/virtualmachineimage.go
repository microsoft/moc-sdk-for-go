// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package virtualmachineimage

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"

	"github.com/microsoft/moc/pkg/errors"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	
)

// Conversion functions from compute to wssdcloudcompute
func getWssdVirtualMachineImage(c *compute.VirtualMachineImage, groupName string) (*wssdcloudcompute.VirtualMachineImage, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Machine Image name is missing")
	}

	if len(groupName) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	wssdvhd := &wssdcloudcompute.VirtualMachineImage{
		Name:      *c.Name,
		GroupName: groupName,
	}

	if c.VirtualMachineImageProperties != nil {
		if &c.VirtualMachineImageProperties.HyperVGeneration != nil {
			wssdvhd.HyperVGeneration = c.VirtualMachineImageProperties.HyperVGeneration
		} else {
			wssdvhd.HyperVGeneration = common.HyperVGeneration_HyperVGenerationV2
		}
	}
	return wssdvhd, nil
}

// Conversion function from wssdcloudcompute to compute
func getVirtualMachineImage(c *wssdcloudcompute.VirtualMachineImage, group string) *compute.VirtualMachineImage {
	return &compute.VirtualMachineImage{
		Name:                          &c.Name,
		ID:                            &c.Id,
		VirtualMachineImageProperties: &compute.VirtualMachineImageProperties{
			HyperVGeneration: c.HyperVGeneration,
		},
	}
}
