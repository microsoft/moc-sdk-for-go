// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"github.com/microsoft/moc-proto/pkg/convert"
	"github.com/microsoft/moc-proto/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/compute"

	"github.com/microsoft/moc-proto/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc-proto/rpc/cloudagent/compute"
	wssdcommon "github.com/microsoft/moc-proto/rpc/common"
)

func (c *client) getWssdVirtualMachine(vm *compute.VirtualMachine, group string) (*wssdcloudcompute.VirtualMachine, error) {
	if vm.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Machine name is missing")
	}
	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	storageConfig, err := c.getWssdVirtualMachineStorageConfiguration(vm.StorageProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Storage Configuration")
	}
	hardwareConfig, err := c.getWssdVirtualMachineHardwareConfiguration(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Hardware Configuration")
	}
	securityConfig, err := c.getWssdVirtualMachineSecurityConfiguration(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Security Configuration")
	}
	osconfig, err := c.getWssdVirtualMachineOSConfiguration(vm.OsProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get OS Configuration")
	}

	networkConfig, err := c.getWssdVirtualMachineNetworkConfiguration(vm.NetworkProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Network Configuration")
	}

	vmtype := wssdcloudcompute.VMType_TENANT
	if vm.VmType == compute.LoadBalancer {
		vmtype = wssdcloudcompute.VMType_LOADBALANCER
	} else if vm.VmType == compute.StackedControlPlane {
		vmtype = wssdcloudcompute.VMType_STACKEDCONTROLPLANE
	}

	vmOut := wssdcloudcompute.VirtualMachine{
		Name:      *vm.Name,
		Storage:   storageConfig,
		Hardware:  hardwareConfig,
		Security:  securityConfig,
		Os:        osconfig,
		Network:   networkConfig,
		GroupName: group,
		VmType:    vmtype,
	}

	if vm.DisableHighAvailability != nil {
		vmOut.DisableHighAvailability = *vm.DisableHighAvailability
	}

	if vm.Version != nil {
		if vmOut.Status == nil {
			vmOut.Status = status.InitStatus()
		}
		vmOut.Status.Version.Number = *vm.Version
	}

	if vm.Location != nil {
		vmOut.LocationName = *vm.Location
	}

	return &vmOut, nil
}

func (c *client) getWssdVirtualMachineStorageConfiguration(s *compute.StorageProfile) (*wssdcloudcompute.StorageConfiguration, error) {
	wssdstorage := &wssdcloudcompute.StorageConfiguration{
		Osdisk:    &wssdcloudcompute.Disk{},
		Datadisks: []*wssdcloudcompute.Disk{},
	}

	if s.ImageReference != nil {
		imageReference, err := c.getWssdVirtualMachineStorageConfigurationImageReference(s.ImageReference)
		if err != nil {
			return nil, err
		}
		wssdstorage.ImageReference = imageReference
	}

	if s.OsDisk != nil {
		osdisk, err := c.getWssdVirtualMachineStorageConfigurationOsDisk(s.OsDisk)
		if err != nil {
			return nil, errors.Wrapf(err, "Invalid Storage Configuration")
		}
		wssdstorage.Osdisk = osdisk
	}

	if s.DataDisks == nil {
		return wssdstorage, nil
	}

	for _, datadisk := range *s.DataDisks {
		wssddatadisk, err := c.getWssdVirtualMachineStorageConfigurationDataDisk(&datadisk)
		if err != nil {
			return nil, err
		}
		wssdstorage.Datadisks = append(wssdstorage.Datadisks, wssddatadisk)
	}

	if s.VmConfigContainerName != nil {
		wssdstorage.VmConfigContainerName = *s.VmConfigContainerName
	}

	return wssdstorage, nil
}

func (c *client) getWssdVirtualMachineStorageConfigurationImageReference(s *compute.ImageReference) (string, error) {
	if s.Name == nil {
		return "", errors.Wrapf(errors.InvalidInput, "Invalid Image Reference Name")
	}
	return *s.Name, nil
}
func (c *client) getWssdVirtualMachineStorageConfigurationOsDisk(s *compute.OSDisk) (*wssdcloudcompute.Disk, error) {
	if s.Vhd == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd Configuration is missing in OSDisk")
	}
	if s.Vhd.URI == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd URI Configuration is missing in OSDisk")
	}
	return &wssdcloudcompute.Disk{
		Diskname: *s.Vhd.URI,
	}, nil
}

func (c *client) getWssdVirtualMachineStorageConfigurationDataDisk(s *compute.DataDisk) (*wssdcloudcompute.Disk, error) {
	if s.Vhd == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd Configuration is missing in DataDisk")
	}
	if s.Vhd.URI == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd URI Configuration is missing in DataDisk ")
	}
	return &wssdcloudcompute.Disk{
		Diskname: *s.Vhd.URI,
	}, nil
}

func (c *client) getWssdVirtualMachineHardwareConfiguration(vm *compute.VirtualMachine) (*wssdcloudcompute.HardwareConfiguration, error) {
	sizeType := wssdcommon.VirtualMachineSizeType_Default
	if vm.HardwareProfile != nil {
		sizeType = compute.GetCloudVirtualMachineSizeFromCloudSdkVirtualMachineSize(vm.HardwareProfile.VMSize)
	}
	wssdhardware := &wssdcloudcompute.HardwareConfiguration{
		VMSize: sizeType,
	}
	return wssdhardware, nil
}

func (c *client) getWssdVirtualMachineSecurityConfiguration(vm *compute.VirtualMachine) (*wssdcloudcompute.SecurityConfiguration, error) {
	enableTPM := false
	if vm.SecurityProfile != nil {
		enableTPM = *vm.SecurityProfile.EnableTPM
	}
	wssdsecurity := &wssdcloudcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}
	return wssdsecurity, nil
}

func (c *client) getWssdVirtualMachineNetworkConfiguration(s *compute.NetworkProfile) (*wssdcloudcompute.NetworkConfiguration, error) {
	nc := &wssdcloudcompute.NetworkConfiguration{
		Interfaces: []*wssdcloudcompute.NetworkInterface{},
	}

	if s == nil || s.NetworkInterfaces == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.ID == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Network Interface ID/Name is missing")
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcloudcompute.NetworkInterface{NetworkInterfaceName: *nic.ID})
	}

	return nc, nil
}

func (c *client) getWssdVirtualMachineOSSSHPublicKeys(ssh *compute.SSHConfiguration) ([]*wssdcloudcompute.SSHPublicKey, error) {
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

func (c *client) getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) (*wssdcloudcompute.OperatingSystemConfiguration, error) {
	publickeys := []*wssdcloudcompute.SSHPublicKey{}
	osType := wssdcommon.OperatingSystemType_WINDOWS
	var err error

	if s.LinuxConfiguration != nil || s.WindowsConfiguration != nil {
		var sshConfiguration *compute.SSHConfiguration = nil

		if s.LinuxConfiguration != nil {
			sshConfiguration = s.LinuxConfiguration.SSH
		} else if s.WindowsConfiguration != nil {
			sshConfiguration = s.WindowsConfiguration.SSH
		}

		if sshConfiguration != nil {
			publickeys, err = c.getWssdVirtualMachineOSSSHPublicKeys(sshConfiguration)
			if err != nil {
				return nil, errors.Wrapf(err, "SSH Configuration Invalid")
			}
		}
	}

	switch s.OsType {
	case compute.Linux:
		osType = wssdcommon.OperatingSystemType_LINUX
	case compute.Windows:
		osType = wssdcommon.OperatingSystemType_WINDOWS
	default:
		if s.LinuxConfiguration != nil {
			osType = wssdcommon.OperatingSystemType_LINUX
		} else {
			osType = wssdcommon.OperatingSystemType_WINDOWS
		}
	}

	adminuser := &wssdcloudcompute.UserConfiguration{}
	if s.AdminUsername != nil {
		adminuser.Username = *s.AdminUsername
	}

	if s.AdminPassword != nil {
		adminuser.Password = *s.AdminPassword
	}

	if s.ComputerName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "ComputerName is missing")
	}

	osconfig := wssdcloudcompute.OperatingSystemConfiguration{
		ComputerName:  *s.ComputerName,
		Administrator: adminuser,
		Users:         []*wssdcloudcompute.UserConfiguration{},
		Publickeys:    publickeys,
		Ostype:        osType,
	}

	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}
	return &osconfig, nil
}

// Conversion functions from wssdcloudcompute to compute

func (c *client) getVirtualMachine(vm *wssdcloudcompute.VirtualMachine, group string) *compute.VirtualMachine {
	vmtype := compute.Tenant
	if vm.VmType == wssdcloudcompute.VMType_LOADBALANCER {
		vmtype = compute.LoadBalancer
	} else if vm.VmType == wssdcloudcompute.VMType_STACKEDCONTROLPLANE {
		vmtype = compute.StackedControlPlane
	}
	return &compute.VirtualMachine{
		Name: &vm.Name,
		ID:   &vm.Id,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			ProvisioningState:       status.GetProvisioningState(vm.GetStatus().GetProvisioningStatus()),
			Statuses:                c.getVirtualMachineStatuses(vm),
			StorageProfile:          c.getVirtualMachineStorageProfile(vm.Storage),
			HardwareProfile:         c.getVirtualMachineHardwareProfile(vm),
			SecurityProfile:         c.getVirtualMachineSecurityProfile(vm),
			OsProfile:               c.getVirtualMachineOSProfile(vm.Os),
			NetworkProfile:          c.getVirtualMachineNetworkProfile(vm.Network),
			VmType:                  vmtype,
			DisableHighAvailability: &vm.DisableHighAvailability,
			Host:                    c.getVirtualMachineHostDescription(vm),
		},
		Version:  &vm.Status.Version.Number,
		Location: &vm.LocationName,
	}
}

func (c *client) getVirtualMachineStatuses(vm *wssdcloudcompute.VirtualMachine) map[string]*string {
	statuses := status.GetStatuses(vm.GetStatus())
	statuses["PowerState"] = convert.ToStringPtr(vm.GetPowerState().String())
	return statuses
}

func (c *client) getVirtualMachineStorageProfile(s *wssdcloudcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		ImageReference:        c.getVirtualMachineStorageProfileImageReference(s.ImageReference),
		OsDisk:                c.getVirtualMachineStorageProfileOsDisk(s.Osdisk),
		DataDisks:             c.getVirtualMachineStorageProfileDataDisks(s.Datadisks),
		VmConfigContainerName: &s.VmConfigContainerName,
	}
}

func (c *client) getVirtualMachineStorageProfileImageReference(imageReference string) *compute.ImageReference {
	return &compute.ImageReference{
		Name: &imageReference,
	}
}
func (c *client) getVirtualMachineStorageProfileOsDisk(d *wssdcloudcompute.Disk) *compute.OSDisk {
	if d == nil {
		return &compute.OSDisk{}
	}
	return &compute.OSDisk{
		Vhd: &compute.VirtualHardDisk{URI: &d.Diskname},
	}
}

func (c *client) getVirtualMachineStorageProfileDataDisks(dd []*wssdcloudcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd,
			compute.DataDisk{
				Vhd: &compute.VirtualHardDisk{URI: &(i.Diskname)},
			},
		)
	}

	return &cdd

}

func (c *client) getVirtualMachineHardwareProfile(vm *wssdcloudcompute.VirtualMachine) *compute.HardwareProfile {
	sizeType := compute.VirtualMachineSizeTypesDefault
	if vm.Hardware != nil {
		sizeType = compute.GetCloudSdkVirtualMachineSizeFromCloudVirtualMachineSize(vm.Hardware.VMSize)
	}
	return &compute.HardwareProfile{
		VMSize: sizeType,
	}
}

func (c *client) getVirtualMachineSecurityProfile(vm *wssdcloudcompute.VirtualMachine) *compute.SecurityProfile {
	enableTPM := false
	if vm.Security != nil {
		enableTPM = vm.Security.EnableTPM
	}
	return &compute.SecurityProfile{
		EnableTPM: &enableTPM,
	}
}

func (c *client) getVirtualMachineHostDescription(vm *wssdcloudcompute.VirtualMachine) *compute.SubResource {
	return &compute.SubResource{
		ID: &vm.NodeName,
	}
}

func (c *client) getVirtualMachineNetworkProfile(n *wssdcloudcompute.NetworkConfiguration) *compute.NetworkProfile {
	np := &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.NetworkInterfaceReference{ID: &((*nic).NetworkInterfaceName)})
	}
	return np
}
func (c *client) getVirtualMachineOSProfile(o *wssdcloudcompute.OperatingSystemConfiguration) *compute.OSProfile {
	osType := compute.Windows
	switch o.Ostype {
	case wssdcommon.OperatingSystemType_LINUX:
		osType = compute.Linux
	case wssdcommon.OperatingSystemType_WINDOWS:
		osType = compute.Windows
	}

	return &compute.OSProfile{
		ComputerName: &o.ComputerName,
		OsType:       osType,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
	}
}
