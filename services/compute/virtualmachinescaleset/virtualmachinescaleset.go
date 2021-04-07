// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachinescaleset

import (
	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

func (c *client) getVirtualMachineScaleSet(vmss *wssdcloudcompute.VirtualMachineScaleSet, group string) (*compute.VirtualMachineScaleSet, error) {
	vmprofile, err := c.getVirtualMachineScaleSetVMProfile(vmss.Virtualmachineprofile)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set VM Profile is invalid")
	}
	return &compute.VirtualMachineScaleSet{
		Name:     &vmss.Name,
		ID:       &vmss.Id,
		Version:  &vmss.Status.Version.Number,
		Location: &vmss.LocationName,
		Sku: &compute.Sku{
			Name:     &vmss.Sku.Name,
			Capacity: &vmss.Sku.Capacity,
		},
		VirtualMachineScaleSetProperties: &compute.VirtualMachineScaleSetProperties{
			ProvisioningState:     status.GetProvisioningState(vmss.GetStatus().GetProvisioningStatus()),
			Statuses:              status.GetStatuses(vmss.GetStatus()),
			VirtualMachineProfile: vmprofile,
		},
	}, nil
}

func (c *client) getVirtualMachineScaleSetVMProfile(vm *wssdcloudcompute.VirtualMachineProfile) (*compute.VirtualMachineScaleSetVMProfile, error) {
	net, err := c.getVirtualMachineScaleSetNetworkProfile(vm.Network)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set VM Profile is invalid")
	}
	storage, err := c.getVirtualMachineScaleSetStorageProfile(vm.Storage)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set Storage Profile is invalid")
	}
	hardware, err := c.getVirtualMachineScaleSetHardwareProfile(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set Hardware Profile is invalid")
	}
	security, err := c.getVirtualMachineScaleSetSecurityProfile(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set Security Profile is invalid")
	}
	os, err := c.getVirtualMachineScaleSetOSProfile(vm.Os)
	if err != nil {
		return nil, errors.Wrapf(err, "Virtual Machine Scale Set OS Profile is invalid")
	}

	return &compute.VirtualMachineScaleSetVMProfile{
		StorageProfile:  storage,
		HardwareProfile: hardware,
		SecurityProfile: security,
		OsProfile:       os,
		NetworkProfile:  net,
	}, nil
}

func (c *client) getVirtualMachineScaleSetStorageProfile(s *wssdcloudcompute.StorageConfiguration) (*compute.VirtualMachineScaleSetStorageProfile, error) {
	osdisk, err := c.getVirtualMachineScaleSetStorageProfileOsDisk(s.Osdisk)
	if err != nil {
		return nil, errors.Wrapf(err, "VMSS Invalid Storage Profile")
	}
	dataDisks := []compute.VirtualMachineScaleSetDataDisk{}
	storageProfile := &compute.VirtualMachineScaleSetStorageProfile{
		ImageReference: &compute.ImageReference{Name: &s.ImageReference},
		OsDisk:         osdisk,
		DataDisks:      &dataDisks,
	}

	for _, dd := range s.Datadisks {
		datadisk, err := c.getVirtualMachineScaleSetStorageProfileDataDisk(dd)
		if err != nil {
			return nil, errors.Wrapf(err, "VMSS Invalid Storage Profile")
		}
		dataDisks = append(dataDisks, *datadisk)
	}

	return storageProfile, nil
}

func (c *client) getVirtualMachineScaleSetStorageProfileOsDisk(d *wssdcloudcompute.Disk) (*compute.VirtualMachineScaleSetOSDisk, error) {
	osdisk := &compute.VirtualMachineScaleSetOSDisk{
		Image: &compute.VirtualHardDisk{},
	}
	if d != nil {
		osdisk.Image.URI = &d.Diskname
	}
	return osdisk, nil
}

func (c *client) getVirtualMachineScaleSetStorageProfileDataDisk(dd *wssdcloudcompute.Disk) (*compute.VirtualMachineScaleSetDataDisk, error) {
	return &compute.VirtualMachineScaleSetDataDisk{Image: &compute.VirtualHardDisk{URI: &(dd.Diskname)}}, nil
}

func (c *client) getVirtualMachineScaleSetHardwareProfile(vm *wssdcloudcompute.VirtualMachineProfile) (*compute.VirtualMachineScaleSetHardwareProfile, error) {
	sizeType := compute.VirtualMachineSizeTypesDefault
	var customSize *compute.VirtualMachineCustomSize
	if vm.Hardware != nil {
		sizeType = compute.GetCloudSdkVirtualMachineSizeFromCloudVirtualMachineSize(vm.Hardware.VMSize)
		if vm.Hardware.CustomSize != nil {
			customSize = &compute.VirtualMachineCustomSize{
				CpuCount: &vm.Hardware.CustomSize.CpuCount,
				MemoryMB: &vm.Hardware.CustomSize.MemoryMB,
			}
		}
	}
	hardwareProfile := &compute.VirtualMachineScaleSetHardwareProfile{
		VMSize:     sizeType,
		CustomSize: customSize,
	}

	return hardwareProfile, nil
}

func (c *client) getVirtualMachineScaleSetSecurityProfile(vm *wssdcloudcompute.VirtualMachineProfile) (*compute.SecurityProfile, error) {
	enableTPM := false
	if vm.Security != nil {
		enableTPM = vm.Security.EnableTPM
	}
	securityProfile := &compute.SecurityProfile{
		EnableTPM: &enableTPM,
	}

	return securityProfile, nil
}

func (c *client) getVirtualMachineScaleSetNetworkProfile(n *wssdcloudcompute.NetworkConfigurationScaleSet) (*compute.VirtualMachineScaleSetNetworkProfile, error) {
	np := &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &[]compute.VirtualMachineScaleSetNetworkConfiguration{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		vnic, err := c.getVirtualMachineScaleSetNetworkConfiguration(nic)
		if err != nil {
			return nil, err
		}
		*np.NetworkInterfaceConfigurations = append(*np.NetworkInterfaceConfigurations, *vnic)
	}
	return np, nil
}

func (c *client) getVirtualMachineScaleSetNetworkConfiguration(nic *wssdcloudnetwork.NetworkInterface) (*compute.VirtualMachineScaleSetNetworkConfiguration, error) {
	ipConfigs := []compute.VirtualMachineScaleSetIPConfiguration{}
	for _, wssdipconfig := range nic.IpConfigurations {
		ipconfig, err := c.getVirtualMachineScaleSetNetworkConfigurationIpConfiguration(wssdipconfig)
		if err != nil {
			return nil, err
		}
		ipConfigs = append(ipConfigs, *ipconfig)
	}

	return &compute.VirtualMachineScaleSetNetworkConfiguration{
		VirtualMachineScaleSetNetworkConfigurationProperties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
			IPConfigurations: &ipConfigs,
		},
	}, nil
}

func (c *client) getVirtualMachineScaleSetNetworkConfigurationIpConfiguration(nic *wssdcloudnetwork.IpConfiguration) (*compute.VirtualMachineScaleSetIPConfiguration, error) {
	return &compute.VirtualMachineScaleSetIPConfiguration{
		VirtualMachineScaleSetIPConfigurationProperties: &compute.VirtualMachineScaleSetIPConfigurationProperties{
			Subnet: &compute.APIEntityReference{
				ID: &nic.Subnetid,
			},
			Primary: &nic.Primary,
		},
	}, nil
}

func (c *client) getVirtualMachineWindowsConfiguration(windowsConfiguration *wssdcloudcompute.WindowsConfiguration) *compute.WindowsConfiguration {
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

func (c *client) getVirtualMachineLinuxConfiguration(linuxConfiguration *wssdcloudcompute.LinuxConfiguration) *compute.LinuxConfiguration {
	lc := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &linuxConfiguration.DisablePasswordAuthentication,
	}

	return lc
}

func (c *client) getVirtualMachineScaleSetOSProfile(o *wssdcloudcompute.OperatingSystemConfiguration) (*compute.VirtualMachineScaleSetOSProfile, error) {
	publicKeys := []compute.SSHPublicKey{}
	for _, wssdpkey := range o.Publickeys {
		publicKey, err := c.getVirtualMachineScaleSetOSProfileSSHKeys(wssdpkey)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, *publicKey)
	}

	ssh := compute.SSHConfiguration{PublicKeys: &publicKeys}

	osBootstrapEngine := compute.CloudInit
	switch o.OsBootstrapEngine {
	case wssdcommon.OperatingSystemBootstrapEngine_WINDOWS_ANSWER_FILES:
		osBootstrapEngine = compute.WindowsAnswerFiles
	case wssdcommon.OperatingSystemBootstrapEngine_CLOUD_INIT:
		fallthrough
	default:
		osBootstrapEngine = compute.CloudInit
	}

	osprofile := &compute.VirtualMachineScaleSetOSProfile{
		ComputerNamePrefix: &o.ComputerName,
		CustomData:         &o.CustomData,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
		// Publickeys: &o.Publickeys,
		// Users : &o.Users,
		OsBootstrapEngine:    osBootstrapEngine,
		WindowsConfiguration: c.getVirtualMachineWindowsConfiguration(o.WindowsConfiguration),
		LinuxConfiguration:   c.getVirtualMachineLinuxConfiguration(o.LinuxConfiguration),
	}

	switch o.Ostype {
	case wssdcommon.OperatingSystemType_LINUX:
		osprofile.LinuxConfiguration = &compute.LinuxConfiguration{SSH: &ssh}
	case wssdcommon.OperatingSystemType_WINDOWS:
		osprofile.WindowsConfiguration = &compute.WindowsConfiguration{SSH: &ssh}
	}

	return osprofile, nil
}

func (c *client) getVirtualMachineScaleSetOSProfileSSHKeys(k *wssdcloudcompute.SSHPublicKey) (*compute.SSHPublicKey, error) {
	return &compute.SSHPublicKey{KeyData: &k.Keydata}, nil
}

// Conversion from sdk to protobuf
func (c *client) getWssdVirtualMachineScaleSet(vmss *compute.VirtualMachineScaleSet, group string) (*wssdcloudcompute.VirtualMachineScaleSet, error) {
	if vmss == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "VirtualMachineScaleSet Input is missing")
	}
	vm, err := c.getWssdVirtualMachineScaleSetVMProfile(vmss.VirtualMachineProfile)
	if err != nil {
		return nil, err
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	if vmss.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS name is missing")
	}
	if vmss.Sku == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Sku is missing")
	}

	if vmss.Sku.Name == nil || vmss.Sku.Capacity == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Sku [Capacity/Name] is missing")
	}

	vmScaleSet := wssdcloudcompute.VirtualMachineScaleSet{
		Name: *(vmss.Name),
		Sku: &wssdcloudcompute.Sku{
			Name:     *(vmss.Sku.Name),
			Capacity: *(vmss.Sku.Capacity),
		},
		Virtualmachineprofile: vm,
		GroupName:             group,
	}

	if vmss.Version != nil {
		if vmScaleSet.Status == nil {
			vmScaleSet.Status = status.InitStatus()
		}
		vmScaleSet.Status.Version.Number = *vmss.Version
	}

	if vmss.Location != nil {
		vmScaleSet.LocationName = *vmss.Location
	}

	return &vmScaleSet, nil
}

func (c *client) getWssdVirtualMachineScaleSetVMProfile(vmp *compute.VirtualMachineScaleSetVMProfile) (*wssdcloudcompute.VirtualMachineProfile, error) {
	if vmp == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "VirtualMachineScaleSetVMProfile Input is missing")
	}
	net, err := c.getWssdVirtualMachineScaleSetNetworkConfiguration(vmp.NetworkProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid VMSS VMProfile Network")
	}

	storage, err := c.getWssdVirtualMachineScaleSetStorageConfiguration(vmp.StorageProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid VMSS VMProfile Storage")
	}
	hardware, err := c.getWssdVirtualMachineScaleSetHardwareConfiguration(vmp)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid VMSS VMProfile Hardware")
	}
	security, err := c.getWssdVirtualMachineScaleSetSecurityConfiguration(vmp)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid VMSS VMProfile Security")
	}
	os, err := c.getWssdVirtualMachineScaleSetOSConfiguration(vmp.OsProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Invalid VMSS VMProfile OS")
	}

	return &wssdcloudcompute.VirtualMachineProfile{
		Storage:  storage,
		Hardware: hardware,
		Security: security,
		Os:       os,
		Network:  net,
	}, nil

}

func (c *client) getWssdVirtualMachineScaleSetStorageConfiguration(s *compute.VirtualMachineScaleSetStorageProfile) (*wssdcloudcompute.StorageConfiguration, error) {
	wssdstorage := &wssdcloudcompute.StorageConfiguration{
		Datadisks: []*wssdcloudcompute.Disk{},
	}

	if s.ImageReference != nil && s.ImageReference.Name != nil && len(*s.ImageReference.Name) > 0 {
		wssdstorage.ImageReference = *s.ImageReference.Name
	} else if s.OsDisk != nil {
		osdisk, err := c.getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s.OsDisk)
		if err != nil {
			return nil, errors.Wrapf(err, "Invalid VMSS Storage Configuration")
		}
		wssdstorage.Osdisk = osdisk
	} else {
		return nil, errors.Wrapf(errors.InvalidInput, "Either ImageReference or OsDisk is missing")
	}

	if s.DataDisks == nil {
		return wssdstorage, nil
	}
	for _, dd := range *s.DataDisks {
		wssddatadisk, err := c.getWssdVirtualMachineScaleSetStorageConfigurationDataDisk(&dd)
		if err != nil {
			return nil, errors.Wrapf(err, "Invalid VMSS Storage Configuration")
		}
		wssdstorage.Datadisks = append(wssdstorage.Datadisks, wssddatadisk)
	}
	return wssdstorage, nil

}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s *compute.VirtualMachineScaleSetOSDisk) (*wssdcloudcompute.Disk, error) {
	if s.Image == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Storage Configuration OSDisk is missing")
	}
	if s.Image.URI == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Storage Configuration OSDisk URI is missing")
	}
	return &wssdcloudcompute.Disk{
		Diskname: *s.Image.URI,
	}, nil
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationDataDisk(d *compute.VirtualMachineScaleSetDataDisk) (*wssdcloudcompute.Disk, error) {
	if d.Image == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Storage Configuration DataDisk is missing")
	}
	if d.Image.URI == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "VMSS Storage Configuration DataDisk URI is missing")
	}
	return &wssdcloudcompute.Disk{Diskname: *d.Image.URI}, nil
}

func (c *client) getWssdVirtualMachineScaleSetHardwareConfiguration(vmp *compute.VirtualMachineScaleSetVMProfile) (*wssdcloudcompute.HardwareConfiguration, error) {
	sizeType := wssdcommon.VirtualMachineSizeType_Default
	var customSize *wssdcommon.VirtualMachineCustomSize
	if vmp.HardwareProfile != nil {
		sizeType = compute.GetCloudVirtualMachineSizeFromCloudSdkVirtualMachineSize(vmp.HardwareProfile.VMSize)
		if vmp.HardwareProfile.CustomSize != nil {
			customSize = &wssdcommon.VirtualMachineCustomSize{
				CpuCount: *vmp.HardwareProfile.CustomSize.CpuCount,
				MemoryMB: *vmp.HardwareProfile.CustomSize.MemoryMB,
			}
		}
	}
	wssdhardware := &wssdcloudcompute.HardwareConfiguration{
		VMSize:     sizeType,
		CustomSize: customSize,
	}
	return wssdhardware, nil
}

func (c *client) getWssdVirtualMachineScaleSetSecurityConfiguration(vmp *compute.VirtualMachineScaleSetVMProfile) (*wssdcloudcompute.SecurityConfiguration, error) {
	enableTPM := false
	if vmp.SecurityProfile != nil {
		enableTPM = *vmp.SecurityProfile.EnableTPM
	}
	wssdsecurity := &wssdcloudcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}
	return wssdsecurity, nil
}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfiguration(s *compute.VirtualMachineScaleSetNetworkProfile) (*wssdcloudcompute.NetworkConfigurationScaleSet, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "VirtualMachineScaleSetNetworkProfile Input is missing")
	}
	nc := &wssdcloudcompute.NetworkConfigurationScaleSet{
		Interfaces: []*wssdcloudnetwork.NetworkInterface{},
	}
	if s == nil || s.NetworkInterfaceConfigurations == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaceConfigurations {
		vnic, err := c.getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface(&nic)
		if err != nil {
			return nil, err
		}
		nc.Interfaces = append(nc.Interfaces, vnic)
	}

	return nc, nil
}

// getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface gets
func (c *client) getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface(nic *compute.VirtualMachineScaleSetNetworkConfiguration) (*wssdcloudnetwork.NetworkInterface, error) {
	nicName := ""
	if nic.Name != nil {
		nicName = *nic.Name
	}
	if nic.IPConfigurations == nil {
		return nil, errors.InvalidConfiguration
	}

	wssdIpConfigs := []*wssdcloudnetwork.IpConfiguration{}
	for _, ipconfig := range *nic.IPConfigurations {
		wssdipconfig, err := c.getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterfaceIpConfiguration(&ipconfig)
		if err != nil {
			return nil, err
		}
		wssdIpConfigs = append(wssdIpConfigs, wssdipconfig)
	}

	return &wssdcloudnetwork.NetworkInterface{
		Name:             nicName,
		IpConfigurations: wssdIpConfigs,
		//Networkname: *nic.VirtualNetworkName,
	}, nil
}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterfaceIpConfiguration(ipconfig *compute.VirtualMachineScaleSetIPConfiguration) (*wssdcloudnetwork.IpConfiguration, error) {
	if ipconfig.VirtualMachineScaleSetIPConfigurationProperties == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing VirtualMachineScaleSetIPConfigurationProperties")
	}
	if ipconfig.Subnet == nil || ipconfig.Subnet.ID == nil || len(*ipconfig.Subnet.ID) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Subnet")
	}

	primary := false
	if ipconfig.Primary != nil {
		primary = *ipconfig.Primary
	}

	return &wssdcloudnetwork.IpConfiguration{
		Primary:  primary,
		Subnetid: *ipconfig.Subnet.ID,
	}, nil
}

func (c *client) getWssdVirtualMachineScaleSetOSSSHPublicKey(sshKey *compute.SSHPublicKey) (*wssdcloudcompute.SSHPublicKey, error) {
	if sshKey.KeyData == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing SSH.KeyData")
	}
	return &wssdcloudcompute.SSHPublicKey{Keydata: *sshKey.KeyData}, nil
}

func (c *client) getWssdVirtualMachineWindowsConfiguration(windowsConfiguration *compute.WindowsConfiguration) *wssdcloudcompute.WindowsConfiguration {
	wc := &wssdcloudcompute.WindowsConfiguration{
		RDPConfiguration: &wssdcloudcompute.RDPConfiguration{},
	}

	if windowsConfiguration == nil {
		return wc
	}

	if windowsConfiguration.RDP.DisableRDP != nil {
		wc.RDPConfiguration.DisableRDP = *windowsConfiguration.RDP.DisableRDP
	}

	if windowsConfiguration.EnableAutomaticUpdates != nil {
		wc.EnableAutomaticUpdates = *windowsConfiguration.EnableAutomaticUpdates
	}

	if windowsConfiguration.TimeZone != nil {
		wc.TimeZone = *windowsConfiguration.TimeZone
	}

	return wc
}

func (c *client) getWssdVirtualMachineLinuxConfiguration(linuxConfiguration *compute.LinuxConfiguration) *wssdcloudcompute.LinuxConfiguration {
	lc := &wssdcloudcompute.LinuxConfiguration{}

	if linuxConfiguration.DisablePasswordAuthentication != nil {
		lc.DisablePasswordAuthentication = *linuxConfiguration.DisablePasswordAuthentication
	}

	return lc

}

func (c *client) getWssdVirtualMachineScaleSetOSConfiguration(s *compute.VirtualMachineScaleSetOSProfile) (*wssdcloudcompute.OperatingSystemConfiguration, error) {
	sshConfig, err := c.getWssdVirtualMachineScaleSetOSConfigurationSSH(s)
	if err != nil {
		return nil, err
	}

	publickeys := []*wssdcloudcompute.SSHPublicKey{}
	if sshConfig.PublicKeys == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing SSH.PublicKeys")
	}
	for _, key := range *sshConfig.PublicKeys {
		publickey, err := c.getWssdVirtualMachineScaleSetOSSSHPublicKey(&key)
		if err != nil {
			return nil, err
		}
		publickeys = append(publickeys, publickey)
	}

	adminuser := &wssdcloudcompute.UserConfiguration{}
	if s.AdminUsername != nil {
		adminuser.Username = *s.AdminUsername
	}
	if s.AdminPassword != nil {
		adminuser.Password = *s.AdminPassword
	}

	osBootstrapEngine := wssdcommon.OperatingSystemBootstrapEngine_CLOUD_INIT
	switch s.OsBootstrapEngine {
	case compute.WindowsAnswerFiles:
		osBootstrapEngine = wssdcommon.OperatingSystemBootstrapEngine_WINDOWS_ANSWER_FILES
	case compute.CloudInit:
		fallthrough
	default:
		osBootstrapEngine = wssdcommon.OperatingSystemBootstrapEngine_CLOUD_INIT
	}

	var windowsConfiguration *wssdcloudcompute.WindowsConfiguration = nil
	if s.WindowsConfiguration != nil {
		windowsConfiguration = c.getWssdVirtualMachineWindowsConfiguration(s.WindowsConfiguration)
	}

	var linuxConfiguration *wssdcloudcompute.LinuxConfiguration = nil
	if s.LinuxConfiguration != nil {
		linuxConfiguration = c.getWssdVirtualMachineLinuxConfiguration(s.LinuxConfiguration)
	}

	osconfig := wssdcloudcompute.OperatingSystemConfiguration{
		ComputerName:         *s.ComputerNamePrefix,
		Administrator:        adminuser,
		Users:                []*wssdcloudcompute.UserConfiguration{},
		Publickeys:           publickeys,
		Ostype:               wssdcommon.OperatingSystemType_WINDOWS,
		OsBootstrapEngine:    osBootstrapEngine,
		WindowsConfiguration: windowsConfiguration,
		LinuxConfiguration:   linuxConfiguration,
	}

	if s.LinuxConfiguration != nil {
		osconfig.Ostype = wssdcommon.OperatingSystemType_LINUX
	}

	// Optional CustomData
	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}
	return &osconfig, nil
}

func (c *client) getWssdVirtualMachineScaleSetOSConfigurationSSH(s *compute.VirtualMachineScaleSetOSProfile) (*compute.SSHConfiguration, error) {
	if s.LinuxConfiguration != nil {
		if s.LinuxConfiguration.SSH == nil {
			return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing LinuxConfiguration.SSH")
		}
		return s.LinuxConfiguration.SSH, nil
	}
	if s.WindowsConfiguration != nil {
		if s.WindowsConfiguration.SSH == nil {
			// return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing WindowsConfiguration.SSH")
			// For now, this is not mandatory. Fill in a dummy value
			tmp := ""
			return &compute.SSHConfiguration{PublicKeys: &[]compute.SSHPublicKey{{KeyData: &tmp}}}, nil
		}
		return s.WindowsConfiguration.SSH, nil
	}
	return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing LinuxConfiguration or WindowsConfiguration")

}
