// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package galleryimage

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudcompute.GalleryImageAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newGalleryImageClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetGalleryImageClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]compute.GalleryImage, error) {
	request, err := getGalleryImageRequest(wssdcloudcommon.Operation_GET, location, "", name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.GalleryImageAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getGalleryImagesFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, imagePath, name string, galleryimage *compute.GalleryImage) (*compute.GalleryImage, error) {
	request, err := getGalleryImageRequest(wssdcloudcommon.Operation_POST, location, imagePath, name, galleryimage)
	if err != nil {
		return nil, err
	}
	response, err := c.GalleryImageAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	galleryimages := getGalleryImagesFromResponse(response, location)

	if len(*galleryimages) == 0 {
		return nil, fmt.Errorf("[GalleryImage][Create] Unexpected error: Creating a compute interface returned no result")
	}

	return &((*galleryimages)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	galleryimage, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*galleryimage) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request, err := getGalleryImageRequest(wssdcloudcommon.Operation_DELETE, location, "", name, &(*galleryimage)[0])
	if err != nil {
		return err
	}
	_, err = c.GalleryImageAgentClient.Invoke(ctx, request)

	return err

}

func (c *client) Precheck(ctx context.Context, location, imagePath string, galleryImages []*compute.GalleryImage) (bool, error) {
	request, err := getGalleryImagePrecheckRequest(location, imagePath, galleryImages)
	if err != nil {
		return false, err
	}
	response, err := c.GalleryImageAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getGalleryImagePrecheckResponse(response)
}

func getGalleryImageRequest(opType wssdcloudcommon.Operation, location, imagePath, name string, compute *compute.GalleryImage) (*wssdcloudcompute.GalleryImageRequest, error) {
	request := &wssdcloudcompute.GalleryImageRequest{
		OperationType: opType,
		GalleryImages: []*wssdcloudcompute.GalleryImage{},
	}

	var err error

	wssdgalleryimage := &wssdcloudcompute.GalleryImage{
		Name:         name,
		LocationName: location,
		SourcePath:   imagePath,
	}

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if compute != nil {
		wssdgalleryimage, err = getWssdGalleryImage(compute, location, imagePath)
		if err != nil {
			return nil, err
		}
		if compute.GalleryImageProperties != nil {
			wssdgalleryimage.SourceType = compute.SourceType
		}
	}
	request.GalleryImages = append(request.GalleryImages, wssdgalleryimage)

	return request, nil
}

func getGalleryImagesFromResponse(response *wssdcloudcompute.GalleryImageResponse, location string) *[]compute.GalleryImage {
	virtualHardDisks := []compute.GalleryImage{}
	for _, galleryimage := range response.GetGalleryImages() {
		virtualHardDisks = append(virtualHardDisks, *(getGalleryImage(galleryimage, location)))
	}

	return &virtualHardDisks
}

func getGalleryImagePrecheckRequest(location, imagePath string, galleryImages []*compute.GalleryImage) (*wssdcloudcompute.GalleryImagePrecheckRequest, error) {
	request := &wssdcloudcompute.GalleryImagePrecheckRequest{}

	protoGalleryImages := make([]*wssdcloudcompute.GalleryImage, 0, len(galleryImages))

	for _, galImage := range galleryImages {
		// can an entry in the galleryImages slice ever be nil here? what would be the meaning of that?
		if galImage != nil {
			protoGalImage, err := getWssdGalleryImage(galImage, location, imagePath)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert GalleryImage to Protobuf representation")
			}
			protoGalleryImages = append(protoGalleryImages, protoGalImage)
		}
	}

	request.GalleryImages = protoGalleryImages
	return request, nil
}

func getGalleryImagePrecheckResponse(response *wssdcloudcompute.GalleryImagePrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}
