// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalhost

import (
	"github.com/microsoft/moc/pkg/convert"
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/compute"

	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

func (c *client) getWssdBareMetalHost(bmh *compute.BareMetalHost, location string) (*wssdcloudcompute.BareMetalHost, error) {
	if bmh.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Bare Metal Host name is missing")
	}
	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	bmhOut := wssdcloudcompute.BareMetalHost{
		Name:         *bmh.Name,
		LocationName: location,
		Tags:         getWssdTags(bmh.Tags),
	}

	if bmh.BareMetalHostProperties != nil {
		storageConfig, err := c.getWssdBareMetalHostStorageConfiguration(bmh.StorageProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Storage Configuration")
		}
		hardwareConfig, err := c.getWssdBareMetalHostHardwareConfiguration(bmh)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Hardware Configuration")
		}
		securityConfig, err := c.getWssdBareMetalHostSecurityConfiguration(bmh)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Security Configuration")
		}

		networkConfig, err := c.getWssdBareMetalHostNetworkConfiguration(bmh.NetworkProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Network Configuration")
		}

		bmhOut.Storage = storageConfig
		bmhOut.Hardware = hardwareConfig
		bmhOut.Security = securityConfig
		bmhOut.Network = networkConfig

		if bmh.BareMetalMachine != nil && bmh.BareMetalMachine.ID != nil {
			bmhOut.BareMetalMachineName = *bmh.BareMetalMachine.ID
		}

		if bmh.FQDN != nil {
			bmhOut.Fqdn = *bmh.FQDN
		}

		if bmh.Port != nil {
			bmhOut.Port = *bmh.Port
		}

		if bmh.AuthorizerPort != nil {
			bmhOut.AuthorizerPort = *bmh.AuthorizerPort
		}

		if bmh.Certificate != nil {
			bmhOut.Certificate = *bmh.Certificate
		}
	}

	if bmh.Version != nil {
		if bmhOut.Status == nil {
			bmhOut.Status = status.InitStatus()
		}
		bmhOut.Status.Version.Number = *bmh.Version
	}

	return &bmhOut, nil
}

func (c *client) getWssdBareMetalHostStorageConfiguration(s *compute.BareMetalHostStorageProfile) (*wssdcloudcompute.BareMetalHostStorageConfiguration, error) {
	wssdstorage := &wssdcloudcompute.BareMetalHostStorageConfiguration{
		Disks: []*wssdcloudcompute.BareMetalHostDisk{},
	}

	if s == nil {
		return wssdstorage, nil
	}

	if s.Disks == nil {
		return wssdstorage, nil
	}

	for _, disk := range *s.Disks {
		wssddisk, err := c.getWssdBareMetalHostStorageConfigurationDisk(&disk)
		if err != nil {
			return nil, err
		}
		wssdstorage.Disks = append(wssdstorage.Disks, wssddisk)
	}

	return wssdstorage, nil
}

func (c *client) getWssdBareMetalHostStorageConfigurationDisk(s *compute.BareMetalHostDisk) (*wssdcloudcompute.BareMetalHostDisk, error) {
	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Name is missing in BareMetalHostDisk")
	}
	if s.DiskSizeGB == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Disk Size is missing in BareMetalHostDisk")
	}
	return &wssdcloudcompute.BareMetalHostDisk{
		DiskName:   *s.Name,
		DiskSizeGB: *s.DiskSizeGB,
	}, nil
}

func (c *client) getWssdBareMetalHostHardwareConfiguration(bmh *compute.BareMetalHost) (*wssdcloudcompute.BareMetalHostHardwareConfiguration, error) {
	var machineSize *wssdcloudcompute.BareMetalHostSize
	if bmh.HardwareProfile != nil && bmh.HardwareProfile.MachineSize != nil {
		machineSize = &wssdcloudcompute.BareMetalHostSize{
			CpuCount: *bmh.HardwareProfile.MachineSize.CpuCount,
			GpuCount: *bmh.HardwareProfile.MachineSize.GpuCount,
			MemoryMB: *bmh.HardwareProfile.MachineSize.MemoryMB,
		}
	}
	wssdhardware := &wssdcloudcompute.BareMetalHostHardwareConfiguration{
		MachineSize: machineSize,
	}
	return wssdhardware, nil
}

func (c *client) getWssdBareMetalHostSecurityConfiguration(bmh *compute.BareMetalHost) (*wssdcloudcompute.SecurityConfiguration, error) {
	enableTPM := false
	if bmh.SecurityProfile != nil {
		enableTPM = *bmh.SecurityProfile.EnableTPM
	}
	wssdsecurity := &wssdcloudcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}
	return wssdsecurity, nil
}

func (c *client) getWssdBareMetalHostNetworkConfiguration(s *compute.BareMetalHostNetworkProfile) (*wssdcloudcompute.BareMetalHostNetworkConfiguration, error) {
	nc := &wssdcloudcompute.BareMetalHostNetworkConfiguration{
		Interfaces: []*wssdcloudcompute.BareMetalHostNetworkInterface{},
	}

	if s == nil || s.NetworkInterfaces == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.Name == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Network Interface Name is missing")
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcloudcompute.BareMetalHostNetworkInterface{NetworkInterfaceName: *nic.Name})
	}

	return nc, nil
}

// Conversion functions from wssdcloudcompute to compute

func (c *client) getBareMetalHost(bmh *wssdcloudcompute.BareMetalHost, location string) *compute.BareMetalHost {
	return &compute.BareMetalHost{
		Name: &bmh.Name,
		ID:   &bmh.Id,
		Tags: getComputeTags(bmh.GetTags()),
		BareMetalHostProperties: &compute.BareMetalHostProperties{
			ProvisioningState: status.GetProvisioningState(bmh.GetStatus().GetProvisioningStatus()),
			Statuses:          c.getBareMetalHostStatuses(bmh),
			StorageProfile:    c.getBareMetalHostStorageProfile(bmh.Storage),
			HardwareProfile:   c.getBareMetalHostHardwareProfile(bmh),
			SecurityProfile:   c.getBareMetalHostSecurityProfile(bmh),
			NetworkProfile:    c.getBareMetalHostNetworkProfile(bmh.Network),
			BareMetalMachine:  c.getBareMetalMachineDescription(bmh),
			FQDN:              &bmh.Fqdn,
			Port:              &bmh.Port,
			AuthorizerPort:    &bmh.AuthorizerPort,
			Certificate:       &bmh.Certificate,
		},
		Version:  &bmh.Status.Version.Number,
		Location: &bmh.LocationName,
	}
}

func (c *client) getBareMetalHostStatuses(bmh *wssdcloudcompute.BareMetalHost) map[string]*string {
	statuses := status.GetStatuses(bmh.GetStatus())
	statuses["PowerState"] = convert.ToStringPtr(bmh.GetPowerState().String())
	return statuses
}

func (c *client) getBareMetalHostStorageProfile(s *wssdcloudcompute.BareMetalHostStorageConfiguration) *compute.BareMetalHostStorageProfile {
	return &compute.BareMetalHostStorageProfile{
		Disks: c.getBareMetalHostStorageProfileDisks(s.Disks),
	}
}

func (c *client) getBareMetalHostStorageProfileDisks(d []*wssdcloudcompute.BareMetalHostDisk) *[]compute.BareMetalHostDisk {
	cd := []compute.BareMetalHostDisk{}

	for _, i := range d {
		cd = append(cd,
			compute.BareMetalHostDisk{
				Name:       &i.DiskName,
				DiskSizeGB: &i.DiskSizeGB,
			},
		)
	}

	return &cd
}

func (c *client) getBareMetalHostHardwareProfile(bmh *wssdcloudcompute.BareMetalHost) *compute.BareMetalHostHardwareProfile {
	var machineSize *compute.BareMetalHostSize
	if bmh.Hardware != nil && bmh.Hardware.MachineSize != nil {
		machineSize = &compute.BareMetalHostSize{
			CpuCount: &bmh.Hardware.MachineSize.CpuCount,
			GpuCount: &bmh.Hardware.MachineSize.GpuCount,
			MemoryMB: &bmh.Hardware.MachineSize.MemoryMB,
		}
	}
	return &compute.BareMetalHostHardwareProfile{
		MachineSize: machineSize,
	}
}

func (c *client) getBareMetalHostSecurityProfile(bmh *wssdcloudcompute.BareMetalHost) *compute.SecurityProfile {
	enableTPM := false
	if bmh.Security != nil {
		enableTPM = bmh.Security.EnableTPM
	}
	return &compute.SecurityProfile{
		EnableTPM: &enableTPM,
	}
}

func (c *client) getBareMetalMachineDescription(bmh *wssdcloudcompute.BareMetalHost) *compute.SubResource {
	return &compute.SubResource{
		ID: &bmh.BareMetalMachineName,
	}
}

func (c *client) getBareMetalHostNetworkProfile(n *wssdcloudcompute.BareMetalHostNetworkConfiguration) *compute.BareMetalHostNetworkProfile {
	np := &compute.BareMetalHostNetworkProfile{
		NetworkInterfaces: &[]compute.BareMetalHostNetworkInterface{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.BareMetalHostNetworkInterface{Name: &((*nic).NetworkInterfaceName)})
	}
	return np
}
