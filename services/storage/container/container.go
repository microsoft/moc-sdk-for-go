// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package container

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/microsoft/moc-sdk-for-go/services/storage"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudstorage "github.com/microsoft/moc/rpc/cloudagent/storage"
)

// Conversion functions from storage to wssdcloudstorage
func getWssdContainer(c *storage.Container, locationName string) (*wssdcloudstorage.Container, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Storage Container name is missing")
	}

	if len(locationName) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}
	wssdcontainer := &wssdcloudstorage.Container{
		Name:         *c.Name,
		LocationName: locationName,
		Tags:         tags.MapToProto(c.Tags),
	}

	if c.Version != nil {
		if wssdcontainer.Status == nil {
			wssdcontainer.Status = status.InitStatus()
		}
		wssdcontainer.Status.Version.Number = *c.Version
	}

	if c.ContainerProperties != nil {
		if c.Path != nil {
			wssdcontainer.Path = *c.Path
		}
		wssdcontainer.Isolated = c.Isolated
	}
	return wssdcontainer, nil
}

func getVirtualharddisktype(enum string) wssdcloudstorage.ContainerType {
	typevalue := wssdcloudstorage.ContainerType(0)
	typevTmp, ok := wssdcloudstorage.ContainerType_value[enum]
	if ok {
		typevalue = wssdcloudstorage.ContainerType(typevTmp)
	}
	return typevalue
}

// Conversion function from wssdcloudstorage to storage
func getContainer(c *wssdcloudstorage.Container, location string) *storage.Container {
	var totalSize string
	var availSize string
	var node string
	var zone string
	if c.Info != nil {
		totalSize = bytefmt.ByteSize(c.Info.Capacity.TotalBytes)
		availSize = bytefmt.ByteSize(c.Info.Capacity.AvailableBytes)
		node = c.Info.Node
		zone = c.Info.Zone
	}
	return &storage.Container{
		Name: &c.Name,
		ID:   &c.Id,
		ContainerProperties: &storage.ContainerProperties{
			Statuses: status.GetStatuses(c.GetStatus()),
			Path:     &c.Path,
			Isolated: c.Isolated,
			ContainerInfo: &storage.ContainerInfo{
				AvailableSize: availSize,
				TotalSize:     totalSize,
				Node:          node,
				Zone:          zone,
			},
		},
		Version: &c.Status.Version.Number,
		Tags:    tags.ProtoToMap(c.Tags),
	}
}
