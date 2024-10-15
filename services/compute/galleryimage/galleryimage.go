// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.
package galleryimage

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

// Conversion functions from compute to wssdcloudcompute
func getWssdGalleryImage(c *compute.GalleryImage, locationName, imagePath string) (*wssdcloudcompute.GalleryImage, error) {
	if c.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Hard Disk name is missing")
	}

	if len(locationName) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}
	wssdgalleryimage := &wssdcloudcompute.GalleryImage{
		Name:         *c.Name,
		LocationName: locationName,
		SourcePath:   imagePath,
		Tags:         tags.MapToProto(c.Tags),
	}

	if c.GalleryImageProperties != nil && c.GalleryImageProperties.ContainerName != nil {
		wssdgalleryimage.ContainerName = *c.GalleryImageProperties.ContainerName
	}

	if c.GalleryImageProperties != nil {
		wssdgalleryimage.SourceType = c.SourceType
		wssdgalleryimage.CloudInitDataSource = c.GalleryImageProperties.CloudInitDataSource
		wssdgalleryimage.HyperVGeneration = c.HyperVGeneration
		switch c.OsType {
		case compute.Windows:
			wssdgalleryimage.ImageOSType = wssdcloudcompute.GalleryImageOSType_WINDOWS
		case compute.Linux:
			wssdgalleryimage.ImageOSType = wssdcloudcompute.GalleryImageOSType_LINUX
		}
	}

	if c.Version != nil {
		if wssdgalleryimage.Status == nil {
			wssdgalleryimage.Status = status.InitStatus()
		}
		wssdgalleryimage.Status.Version.Number = *c.Version
	}

	return wssdgalleryimage, nil
}

// Conversion function from wssdcloudcompute to compute
func getGalleryImage(c *wssdcloudcompute.GalleryImage, location string) *compute.GalleryImage {
	galleryImg := &compute.GalleryImage{
		Name:    &c.Name,
		ID:      &c.Id,
		Version: &c.Status.Version.Number,
		GalleryImageProperties: &compute.GalleryImageProperties{
			Statuses:         status.GetStatuses(c.GetStatus()),
			ContainerName:    &c.ContainerName,
			HyperVGeneration: c.HyperVGeneration,
		},
		Tags: tags.ProtoToMap(c.Tags),
	}
	switch c.ImageOSType {
	case wssdcloudcompute.GalleryImageOSType_WINDOWS:
		galleryImg.OsType = compute.Windows
	case wssdcloudcompute.GalleryImageOSType_LINUX:
		galleryImg.OsType = compute.Linux
	}
	return galleryImg
}
