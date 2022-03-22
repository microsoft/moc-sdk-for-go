// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package storage

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/microsoft/moc/rpc/common"
)

// VirtualHardDiskProperties defines the structure of a Load Balancer
type VirtualHardDiskProperties struct {
	// DiskSizeBytes
	DiskSizeBytes *int64 `json:"diskSizeGB,omitempty"`
	// Dynamic
	Dynamic *bool `json:"dynamic,omitempty"`
	// Blocksizebytes - TODO: Revisit exposing this
	Blocksizebytes *int32 `json:"blocksizebytes,omitempty"`
	//Logicalsectorbytes - TODO: Revisit exposing this
	Logicalsectorbytes *int32 `json:"logicalsectorbytes,omitempty"`
	//Physicalsectorbytes - TODO: Revisit exposing this
	Physicalsectorbytes *int32 `json:"physicalsectorbytes,omitempty"`
	//Controllernumber - TODO: Revisit exposing this
	Controllernumber *int64 `json:"controllernumber,omitempty"`
	//Controllerlocation - TODO: Revisit exposing this
	Controllerlocation *int64 `json:"controllerlocation,omitempty"`
	//Disknumber - TODO: Revisit exposing this
	Disknumber *int64 `json:"disknumber,omitempty"`
	// READONLY - VirtualMachineName to which this disk is attached to
	VirtualMachineName *string `json:"virtualmachinename,omitempty"`
	//Scsipath - TODO: Revisit exposing this
	Scsipath *string `json:"scsipath,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	//HyperVGeneration - Gets the HyperVGenerationType of the VirtualMachine created from the image.
	HyperVGeneration common.HyperVGeneration `json:"hyperVGeneration,omitempty"`
	//DiskFileFormat - File format of the disk
	DiskFileFormat common.DiskFileFormat `json:"diskFileFormat,omitempty"`
}

// VirtualHardDisk defines the structure of a VHD
type VirtualHardDisk struct {
	autorest.Response `json:"-"`
	// Properties
	*VirtualHardDiskProperties `json:"properties,omitempty"`
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}

type ContainerInfo struct {
	AvailableSize string `json:"AvailableSize,omitempty"`
	TotalSize     string `json:"TotalSize,omitempty"`
}

// ContainerProperties defines the structure of a Load Balancer
type ContainerProperties struct {
	// Path
	Path *string `json:"path,omitempty"`
	// State - State
	Statuses       map[string]*string `json:"statuses"`
	*ContainerInfo `json:"info"`
}

// Container defines the structure of a VHD
type Container struct {
	autorest.Response `json:"-"`
	// Properties
	*ContainerProperties `json:"properties,omitempty"`
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}
