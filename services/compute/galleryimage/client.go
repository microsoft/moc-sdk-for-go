// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package galleryimage

import (
	"context"
	"encoding/json"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]compute.GalleryImage, error)
	CreateOrUpdate(context.Context, string, string, string, *compute.GalleryImage) (*compute.GalleryImage, error)
	Delete(context.Context, string, string) error
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
	return c.internal.CreateOrUpdate(ctx, location, imagePath, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *GalleryImageClient) Delete(ctx context.Context, location, name string) error {
	return c.internal.Delete(ctx, location, name)
}

// UploadImageFromLocal   methods invokes  UploadImageFromLocal  on the client
func (c *GalleryImageClient) UploadImageFromLocal(ctx context.Context, location, imagePath, name string, compute *compute.GalleryImage) (*compute.GalleryImage, error) {
	return c.internal.CreateOrUpdate(ctx, location, imagePath, name, compute)
}

// UploadImageFromSFS   methods invokes  UploadImageFromSFS  on the client
func (c *GalleryImageClient) UploadImageFromSFS(ctx context.Context, location, name string, galImage *compute.GalleryImage, sfsImg *compute.SFSImageProperties) (*compute.GalleryImage, error) {

	// convert sfsImg struct to json string and use it as image-path
	data, err := json.Marshal(sfsImg)
	// update galImage with SourceType
	galImage.SourceType = "sfs"

	if err != nil {
		return nil, err
	}
	return c.internal.CreateOrUpdate(ctx, location, string(data), name, galImage)
}
