// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalmachine

import (
	"github.com/microsoft/moc/pkg/convert"
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/compute"

	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

func (c *client) getWssdBareMetalMachine(bmm *compute.BareMetalMachine, location string) (*wssdcloudcompute.BareMetalMachine, error) {
	if bmm.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Bare Metal Machine name is missing")
	}
	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	bmmOut := wssdcloudcompute.BareMetalMachine{
		Name:         *bmm.Name,
		LocationName: location,
		Tags:         getWssdTags(bmm.Tags),
	}

	if bmm.BareMetalMachineProperties != nil {
		storageConfig, err := c.getWssdBareMetalMachineStorageConfiguration(bmm.StorageProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Storage Configuration")
		}
		hardwareConfig, err := c.getWssdBareMetalMachineHardwareConfiguration(bmm)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Hardware Configuration")
		}
		securityConfig, err := c.getWssdBareMetalMachineSecurityConfiguration(bmm)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Security Configuration")
		}
		osConfig, err := c.getWssdBareMetalMachineOSConfiguration(bmm.OsProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get OS Configuration")
		}

		networkConfig, err := c.getWssdBareMetalMachineNetworkConfiguration(bmm.NetworkProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Network Configuration")
		}

		bmmOut.Storage = storageConfig
		bmmOut.Hardware = hardwareConfig
		bmmOut.Security = securityConfig
		bmmOut.Os = osConfig
		bmmOut.Network = networkConfig

		if bmm.Host != nil && bmm.Host.ID != nil {
			bmmOut.NodeName = *bmm.Host.ID
		}

		if bmm.FQDN != nil {
			bmmOut.Fqdn = *bmm.FQDN
		}

		if bmm.Port != nil {
			bmmOut.Port = *bmm.Port
		}

		if bmm.AuthorizerPort != nil {
			bmmOut.AuthorizerPort = *bmm.AuthorizerPort
		}

		if bmm.Certificate != nil {
			bmmOut.Certificate = *bmm.Certificate
		}
	}

	if bmm.Version != nil {
		if bmmOut.Status == nil {
			bmmOut.Status = status.InitStatus()
		}
		bmmOut.Status.Version.Number = *bmm.Version
	}

	return &bmmOut, nil
}

func (c *client) getWssdBareMetalMachineStorageConfiguration(s *compute.BareMetalMachineStorageProfile) (*wssdcloudcompute.BareMetalMachineStorageConfiguration, error) {
	wssdstorage := &wssdcloudcompute.BareMetalMachineStorageConfiguration{
		Disks: []*wssdcloudcompute.BareMetalMachineDisk{},
	}

	if s == nil {
		return wssdstorage, nil
	}

	if s.ImageReference != nil {
		imageReference, err := c.getWssdBareMetalMachineStorageConfigurationImageReference(s.ImageReference)
		if err != nil {
			return nil, err
		}
		wssdstorage.ImageReference = imageReference
	}

	if s.Disks == nil {
		return wssdstorage, nil
	}

	for _, disk := range *s.Disks {
		wssddisk, err := c.getWssdBareMetalMachineStorageConfigurationDisk(&disk)
		if err != nil {
			return nil, err
		}
		wssdstorage.Disks = append(wssdstorage.Disks, wssddisk)
	}

	return wssdstorage, nil
}

func (c *client) getWssdBareMetalMachineStorageConfigurationImageReference(s *compute.BareMetalMachineImageReference) (string, error) {
	if s.Name == nil {
		return "", errors.Wrapf(errors.InvalidInput, "Invalid Image Reference Name")
	}
	return *s.Name, nil
}

func (c *client) getWssdBareMetalMachineStorageConfigurationDisk(s *compute.BareMetalMachineDisk) (*wssdcloudcompute.BareMetalMachineDisk, error) {
	if s.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Name is missing in BareMetalMachineDisk")
	}
	if s.DiskSizeGB == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Disk Size is missing in BareMetalMachineDisk")
	}
	return &wssdcloudcompute.BareMetalMachineDisk{
		DiskName:   *s.Name,
		DiskSizeGB: *s.DiskSizeGB,
	}, nil
}

func (c *client) getWssdBareMetalMachineHardwareConfiguration(bmm *compute.BareMetalMachine) (*wssdcloudcompute.BareMetalMachineHardwareConfiguration, error) {
	var machineSize *wssdcloudcompute.BareMetalMachineSize
	if bmm.HardwareProfile != nil && bmm.HardwareProfile.MachineSize != nil {
		machineSize = &wssdcloudcompute.BareMetalMachineSize{
			CpuCount: *bmm.HardwareProfile.MachineSize.CpuCount,
			GpuCount: *bmm.HardwareProfile.MachineSize.GpuCount,
			MemoryMB: *bmm.HardwareProfile.MachineSize.MemoryMB,
		}
	}
	wssdhardware := &wssdcloudcompute.BareMetalMachineHardwareConfiguration{
		MachineSize: machineSize,
	}
	return wssdhardware, nil
}

func (c *client) getWssdBareMetalMachineSecurityConfiguration(bmm *compute.BareMetalMachine) (*wssdcloudcompute.SecurityConfiguration, error) {
	enableTPM := false
	if bmm.SecurityProfile != nil {
		enableTPM = *bmm.SecurityProfile.EnableTPM
	}
	wssdsecurity := &wssdcloudcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}
	return wssdsecurity, nil
}

func (c *client) getWssdBareMetalMachineNetworkConfiguration(s *compute.BareMetalMachineNetworkProfile) (*wssdcloudcompute.BareMetalMachineNetworkConfiguration, error) {
	nc := &wssdcloudcompute.BareMetalMachineNetworkConfiguration{
		Interfaces: []*wssdcloudcompute.BareMetalMachineNetworkInterface{},
	}

	if s == nil || s.NetworkInterfaces == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.Name == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Network Interface Name is missing")
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcloudcompute.BareMetalMachineNetworkInterface{NetworkInterfaceName: *nic.Name})
	}

	return nc, nil
}

func (c *client) getWssdBareMetalMachineOSSSHPublicKeys(ssh *compute.SSHConfiguration) ([]*wssdcloudcompute.SSHPublicKey, error) {
	keys := []*wssdcloudcompute.SSHPublicKey{}
	if ssh == nil {
		return keys, nil
	}
	for _, key := range *ssh.PublicKeys {
		if key.KeyData == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "SSH KeyData is missing")
		}
		keys = append(keys, &wssdcloudcompute.SSHPublicKey{Keydata: *key.KeyData})
	}
	return keys, nil

}

func (c *client) getWssdBareMetalMachineLinuxConfiguration(linuxConfiguration *compute.LinuxConfiguration) *wssdcloudcompute.LinuxConfiguration {
	lc := &wssdcloudcompute.LinuxConfiguration{}

	if linuxConfiguration.DisablePasswordAuthentication != nil {
		lc.DisablePasswordAuthentication = *linuxConfiguration.DisablePasswordAuthentication
	}

	return lc

}

func (c *client) getWssdBareMetalMachineOSConfiguration(s *compute.BareMetalMachineOSProfile) (*wssdcloudcompute.BareMetalMachineOperatingSystemConfiguration, error) {
	publicKeys := []*wssdcloudcompute.SSHPublicKey{}
	var err error

	osConfig := wssdcloudcompute.BareMetalMachineOperatingSystemConfiguration{
		Users: []*wssdcloudcompute.UserConfiguration{},
	}

	if s == nil {
		return &osConfig, nil
	}

	if s.LinuxConfiguration != nil {
		var sshConfiguration *compute.SSHConfiguration = s.LinuxConfiguration.SSH

		if sshConfiguration != nil {
			publicKeys, err = c.getWssdBareMetalMachineOSSSHPublicKeys(sshConfiguration)
			if err != nil {
				return nil, errors.Wrapf(err, "SSH Configuration Invalid")
			}
		}
	}

	adminUser := &wssdcloudcompute.UserConfiguration{}
	if s.AdminUsername != nil {
		adminUser.Username = *s.AdminUsername
	}

	if s.AdminPassword != nil {
		adminUser.Password = *s.AdminPassword
	}

	if s.ComputerName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "ComputerName is missing")
	}

	osConfig.ComputerName = *s.ComputerName
	osConfig.Administrator = adminUser
	osConfig.PublicKeys = publicKeys

	if s.LinuxConfiguration != nil {
		osConfig.LinuxConfiguration = c.getWssdBareMetalMachineLinuxConfiguration(s.LinuxConfiguration)
	}

	if s.CustomData != nil {
		osConfig.CustomData = *s.CustomData
	}
	return &osConfig, nil
}

// Conversion functions from wssdcloudcompute to compute

func (c *client) getBareMetalMachine(bmm *wssdcloudcompute.BareMetalMachine, location string) *compute.BareMetalMachine {
	return &compute.BareMetalMachine{
		Name: &bmm.Name,
		ID:   &bmm.Id,
		Tags: getComputeTags(bmm.GetTags()),
		BareMetalMachineProperties: &compute.BareMetalMachineProperties{
			ProvisioningState: status.GetProvisioningState(bmm.GetStatus().GetProvisioningStatus()),
			Statuses:          c.getBareMetalMachineStatuses(bmm),
			StorageProfile:    c.getBareMetalMachineStorageProfile(bmm.Storage),
			HardwareProfile:   c.getBareMetalMachineHardwareProfile(bmm),
			SecurityProfile:   c.getBareMetalMachineSecurityProfile(bmm),
			OsProfile:         c.getBareMetalMachineOSProfile(bmm.Os),
			NetworkProfile:    c.getBareMetalMachineNetworkProfile(bmm.Network),
			Host:              c.getBareMetalMachineHostDescription(bmm),
			FQDN:              &bmm.Fqdn,
			Port:              &bmm.Port,
			AuthorizerPort:    &bmm.AuthorizerPort,
			Certificate:       &bmm.Certificate,
		},
		Version:  &bmm.Status.Version.Number,
		Location: &bmm.LocationName,
	}
}

func (c *client) getBareMetalMachineStatuses(bmm *wssdcloudcompute.BareMetalMachine) map[string]*string {
	statuses := status.GetStatuses(bmm.GetStatus())
	statuses["PowerState"] = convert.ToStringPtr(bmm.GetPowerState().String())
	return statuses
}

func (c *client) getBareMetalMachineStorageProfile(s *wssdcloudcompute.BareMetalMachineStorageConfiguration) *compute.BareMetalMachineStorageProfile {
	return &compute.BareMetalMachineStorageProfile{
		ImageReference: c.getBareMetalMachineStorageProfileImageReference(s.ImageReference),
		Disks:          c.getBareMetalMachineStorageProfileDisks(s.Disks),
	}
}

func (c *client) getBareMetalMachineStorageProfileImageReference(imageReference string) *compute.BareMetalMachineImageReference {
	return &compute.BareMetalMachineImageReference{
		Name: &imageReference,
	}
}

func (c *client) getBareMetalMachineStorageProfileDisks(d []*wssdcloudcompute.BareMetalMachineDisk) *[]compute.BareMetalMachineDisk {
	cd := []compute.BareMetalMachineDisk{}

	for _, i := range d {
		cd = append(cd,
			compute.BareMetalMachineDisk{
				Name:       &i.DiskName,
				DiskSizeGB: &i.DiskSizeGB,
			},
		)
	}

	return &cd
}

func (c *client) getBareMetalMachineHardwareProfile(bmm *wssdcloudcompute.BareMetalMachine) *compute.BareMetalMachineHardwareProfile {
	var machineSize *compute.BareMetalMachineSize
	if bmm.Hardware != nil && bmm.Hardware.MachineSize != nil {
		machineSize = &compute.BareMetalMachineSize{
			CpuCount: &bmm.Hardware.MachineSize.CpuCount,
			GpuCount: &bmm.Hardware.MachineSize.GpuCount,
			MemoryMB: &bmm.Hardware.MachineSize.MemoryMB,
		}
	}
	return &compute.BareMetalMachineHardwareProfile{
		MachineSize: machineSize,
	}
}

func (c *client) getBareMetalMachineSecurityProfile(bmm *wssdcloudcompute.BareMetalMachine) *compute.SecurityProfile {
	enableTPM := false
	if bmm.Security != nil {
		enableTPM = bmm.Security.EnableTPM
	}
	return &compute.SecurityProfile{
		EnableTPM: &enableTPM,
	}
}

func (c *client) getBareMetalMachineHostDescription(bmm *wssdcloudcompute.BareMetalMachine) *compute.SubResource {
	return &compute.SubResource{
		ID: &bmm.NodeName,
	}
}

func (c *client) getBareMetalMachineNetworkProfile(n *wssdcloudcompute.BareMetalMachineNetworkConfiguration) *compute.BareMetalMachineNetworkProfile {
	np := &compute.BareMetalMachineNetworkProfile{
		NetworkInterfaces: &[]compute.BareMetalMachineNetworkInterface{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.BareMetalMachineNetworkInterface{Name: &((*nic).NetworkInterfaceName)})
	}
	return np
}

func (c *client) getBareMetalMachineWindowsConfiguration(windowsConfiguration *wssdcloudcompute.WindowsConfiguration) *compute.WindowsConfiguration {
	wc := &compute.WindowsConfiguration{
		RDP: &compute.RDPConfiguration{},
	}

	if windowsConfiguration == nil {
		return wc
	}

	if windowsConfiguration.RDPConfiguration != nil {
		wc.RDP.DisableRDP = &windowsConfiguration.RDPConfiguration.DisableRDP
	}

	wc.EnableAutomaticUpdates = &windowsConfiguration.EnableAutomaticUpdates
	wc.TimeZone = &windowsConfiguration.TimeZone

	return wc
}

func (c *client) getBareMetalMachineLinuxConfiguration(linuxConfiguration *wssdcloudcompute.LinuxConfiguration) *compute.LinuxConfiguration {
	lc := &compute.LinuxConfiguration{}

	if linuxConfiguration != nil {
		lc.DisablePasswordAuthentication = &linuxConfiguration.DisablePasswordAuthentication
	}

	return lc
}

func (c *client) getBareMetalMachineOSProfile(osConfiguration *wssdcloudcompute.BareMetalMachineOperatingSystemConfiguration) *compute.BareMetalMachineOSProfile {
	op := &compute.BareMetalMachineOSProfile{
		ComputerName:       &osConfiguration.ComputerName,
		CustomData:         &osConfiguration.CustomData,
		LinuxConfiguration: c.getBareMetalMachineLinuxConfiguration(osConfiguration.LinuxConfiguration),
	}

	if osConfiguration.Administrator != nil {
		op.AdminUsername = &osConfiguration.Administrator.Username
	}

	if osConfiguration.Administrator != nil {
		op.AdminPassword = &osConfiguration.Administrator.Password
	}

	return op
}
