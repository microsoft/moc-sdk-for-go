// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package galleryimage

import (
	"context"
	"encoding/json"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/rpc/common"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]compute.GalleryImage, error)
	CreateOrUpdate(context.Context, string, string, string, *compute.GalleryImage) (*compute.GalleryImage, error)
	Delete(context.Context, string, string) error
	Precheck(ctx context.Context, location, imagePath string, galleryImages []*compute.GalleryImage) (bool, error)
}

// Client structure
type GalleryImageClient struct {
	compute.BaseClient
	internal Service
}

// NewClient method returns new client
func NewGalleryImageClient(cloudFQDN string, authorizer auth.Authorizer) (*GalleryImageClient, error) {
	c, err := newGalleryImageClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &GalleryImageClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *GalleryImageClient) Get(ctx context.Context, location, name string) (*[]compute.GalleryImage, error) {
	return c.internal.Get(ctx, location, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *GalleryImageClient) CreateOrUpdate(ctx context.Context, location, imagePath, name string, compute *compute.GalleryImage) (*compute.GalleryImage, error) {
	if compute != nil && compute.GalleryImageProperties != nil {
		compute.SourceType = common.ImageSource_LOCAL_SOURCE
	}
	return c.internal.CreateOrUpdate(ctx, location, imagePath, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *GalleryImageClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// Prechecks whether the system is able to create specified resources.
// Returns true if it is possible; or false with reason in error message if not.
func (c *GalleryImageClient) Precheck(ctx context.Context, location, imagePath string, galleryImages []*compute.GalleryImage) (bool, error) {
	return c.internal.Precheck(ctx, location, imagePath, galleryImages)
}

// UploadImageFromLocal   methods invokes  UploadImageFromLocal  on the client
func (c *GalleryImageClient) UploadImageFromLocal(ctx context.Context, location, imagePath, name string, compute *compute.GalleryImage) (*compute.GalleryImage, error) {
	if compute != nil && compute.GalleryImageProperties != nil {
		compute.SourceType = common.ImageSource_LOCAL_SOURCE
	}
	return c.internal.CreateOrUpdate(ctx, location, imagePath, name, compute)
}

// UploadImageFromSFS   methods invokes  UploadImageFromSFS  on the client
func (c *GalleryImageClient) UploadImageFromSFS(ctx context.Context, location, name string, galImage *compute.GalleryImage, sfsImg *compute.SFSImageProperties) (*compute.GalleryImage, error) {
	// convert sfsImg struct to json string and use it as image-path
	data, err := json.Marshal(sfsImg)
	if err != nil {
		return nil, err
	}
	// update galImage with SourceType
	if galImage != nil && galImage.GalleryImageProperties != nil {
		galImage.SourceType = common.ImageSource_SFS_SOURCE
	}

	return c.internal.CreateOrUpdate(ctx, location, string(data), name, galImage)
}

func (c *GalleryImageClient) UploadImageFromHttp(ctx context.Context, location, name string, galImage *compute.GalleryImage, azHttpImg *compute.AzureGalleryImageProperties) (*compute.GalleryImage, error) {
	// convert httpImg struct to json string and use it as image-path
	data, err := json.Marshal(azHttpImg)
	if err != nil {
		return nil, err
	}
	if galImage != nil && galImage.GalleryImageProperties != nil {
		galImage.SourceType = common.ImageSource_HTTP_SOURCE
	}
	return c.internal.CreateOrUpdate(ctx, location, string(data), name, galImage)
}

// UploadImageFromLocal   methods invokes  UploadImageFromLocal  on the client
func (c *GalleryImageClient) UploadImageFromVMOsDisk(ctx context.Context, location, imagePath, name string, compute *compute.GalleryImage) (*compute.GalleryImage, error) {
	if compute != nil && compute.GalleryImageProperties != nil && compute.GalleryImageProperties.SourceVirtualMachine != nil {
		compute.SourceType = common.ImageSource_VMOSDISK_SOURCE
	}
	return c.internal.CreateOrUpdate(ctx, location, imagePath, name, compute)
}
