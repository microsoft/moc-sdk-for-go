// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package virtualharddisk

import (
	"github.com/microsoft/moc-sdk-for-go/services/storage"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudstorage "github.com/microsoft/moc/rpc/cloudagent/storage"
	"github.com/microsoft/moc/rpc/common"
)

// Conversion functions from storage to wssdcloudstorage
func getWssdVirtualHardDisk(c *storage.VirtualHardDisk, groupName, containerName string) (*wssdcloudstorage.VirtualHardDisk, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Hard Disk name is missing")
	}

	if len(groupName) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	wssdvhd := &wssdcloudstorage.VirtualHardDisk{
		Name:          *c.Name,
		GroupName:     groupName,
		ContainerName: containerName,
	}

	if c.Version != nil {
		if wssdvhd.Status == nil {
			wssdvhd.Status = status.InitStatus()
		}
		wssdvhd.Status.Version.Number = *c.Version
	}

	if c.VirtualHardDiskProperties != nil {
		if c.Blocksizebytes != nil {
			wssdvhd.Blocksizebytes = *c.Blocksizebytes
		}
		if c.Dynamic != nil {
			wssdvhd.Dynamic = *c.Dynamic
		}
		if c.Physicalsectorbytes != nil {
			wssdvhd.Physicalsectorbytes = *c.Physicalsectorbytes
		}
		if c.DiskSizeBytes != nil {
			wssdvhd.Size = *c.DiskSizeBytes
		}
		if c.Logicalsectorbytes != nil {
			wssdvhd.Logicalsectorbytes = *c.Logicalsectorbytes
		}
		if c.VirtualMachineName != nil {
			wssdvhd.VirtualmachineName = *c.VirtualMachineName
		}
		if &c.HyperVGeneration != nil {
			wssdvhd.HyperVGeneration = c.HyperVGeneration
		} else {
			wssdvhd.HyperVGeneration = common.HyperVGeneration_HyperVGenerationV2
		}
		if &c.DiskFileFormat != nil {
			wssdvhd.DiskFileFormat = c.DiskFileFormat
		} else {
			wssdvhd.DiskFileFormat = common.DiskFileFormat_DiskFileFormatVHDX
		}
		if &c.CloudInitDataSource != nil {
			wssdvhd.CloudInitDataSource = c.CloudInitDataSource
		} else {
			wssdvhd.CloudInitDataSource = common.CloudInitDataSource_NoCloud
		}
	}
	return wssdvhd, nil
}

// Conversion function from wssdcloudstorage to storage
func getVirtualHardDisk(c *wssdcloudstorage.VirtualHardDisk, group string) *storage.VirtualHardDisk {
	return &storage.VirtualHardDisk{
		Name:    &c.Name,
		ID:      &c.Id,
		Version: &c.Status.Version.Number,
		VirtualHardDiskProperties: &storage.VirtualHardDiskProperties{
			Statuses:            status.GetStatuses(c.GetStatus()),
			DiskSizeBytes:       &c.Size,
			Dynamic:             &c.Dynamic,
			Blocksizebytes:      &c.Blocksizebytes,
			Logicalsectorbytes:  &c.Logicalsectorbytes,
			Physicalsectorbytes: &c.Physicalsectorbytes,
			Controllernumber:    &c.Controllernumber,
			Controllerlocation:  &c.Controllerlocation,
			Disknumber:          &c.Disknumber,
			VirtualMachineName:  &c.VirtualmachineName,
			Scsipath:            &c.Scsipath,
			HyperVGeneration:    c.HyperVGeneration,
			DiskFileFormat:      c.DiskFileFormat,
		},
	}
}
