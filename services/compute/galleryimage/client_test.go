// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package galleryimage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeInternal captures the args passed to CreateOrUpdate so the wrapper logic
// in UploadImageFromAzureStorageBlob can be asserted without a real backend.
type fakeInternal struct {
	gotLocation  string
	gotImagePath string
	gotName      string
	gotImage     *compute.GalleryImage
	callCount    int

	createOrUpdateResp *compute.GalleryImage
	createOrUpdateErr  error
}

var _ Service = (*fakeInternal)(nil)

func (f *fakeInternal) Get(ctx context.Context, location, name string) (*[]compute.GalleryImage, error) {
	return nil, nil
}

func (f *fakeInternal) CreateOrUpdate(ctx context.Context, location, imagePath, name string, img *compute.GalleryImage) (*compute.GalleryImage, error) {
	f.callCount++
	f.gotLocation = location
	f.gotImagePath = imagePath
	f.gotName = name
	f.gotImage = img
	return f.createOrUpdateResp, f.createOrUpdateErr
}

func (f *fakeInternal) Delete(ctx context.Context, location, name string) error {
	return nil
}

func (f *fakeInternal) Precheck(ctx context.Context, location, imagePath string, galleryImages []*compute.GalleryImage) (bool, error) {
	return true, nil
}

func TestUploadImageFromAzureStorageBlob_SetsSourceTypeAndEncodesJSON(t *testing.T) {
	fake := &fakeInternal{}
	c := &GalleryImageClient{internal: fake}

	blobImg := &compute.AzureBlobImageProperties{
		CatalogName: "cat",
		Audience:    "aud",
		Version:     "1.2.3",
		ReleaseName: "rel",
		Parts:       4,
		Cloud:       "AzurePublicCloud",
		Endpoint:    "https://account.blob.core.windows.net/",
	}
	galImage := &compute.GalleryImage{
		GalleryImageProperties: &compute.GalleryImageProperties{},
	}

	_, err := c.UploadImageFromAzureStorageBlob(context.Background(), "loc-1", "img-1", galImage, blobImg)
	require.NoError(t, err)
	require.Equal(t, 1, fake.callCount)

	assert.Equal(t, "loc-1", fake.gotLocation)
	assert.Equal(t, "img-1", fake.gotName)
	require.NotNil(t, fake.gotImage)
	assert.Equal(t, common.ImageSource_AZURESTORAGEBLOB_SOURCE, fake.gotImage.SourceType)

	// imagePath must be the JSON-encoded AzureBlobImageProperties.
	var roundTrip compute.AzureBlobImageProperties
	require.NoError(t, json.Unmarshal([]byte(fake.gotImagePath), &roundTrip))
	assert.Equal(t, *blobImg, roundTrip)
}

func TestUploadImageFromAzureStorageBlob_NilPropertiesReturnsError(t *testing.T) {
	fake := &fakeInternal{}
	c := &GalleryImageClient{internal: fake}

	// galImage without GalleryImageProperties: the wrapper must return an error
	// because SourceType cannot be set, which would cause the download to be
	// routed incorrectly downstream.
	galImage := &compute.GalleryImage{}
	blobImg := &compute.AzureBlobImageProperties{CatalogName: "cat"}

	_, err := c.UploadImageFromAzureStorageBlob(context.Background(), "loc-1", "img-1", galImage, blobImg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "GalleryImageProperties are required")
	assert.Equal(t, 0, fake.callCount, "CreateOrUpdate must not be called when properties are nil")
}
