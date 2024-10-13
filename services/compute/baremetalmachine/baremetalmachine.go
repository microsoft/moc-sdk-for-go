// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package baremetalmachine

import (
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/moc-sdk-for-go/services/compute"

	"github.com/microsoft/moc/pkg/status"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

func (c *client) getWssdBareMetalMachine(bmm *compute.BareMetalMachine, group string) (*wssdcloudcompute.BareMetalMachine, error) {
	if bmm.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Bare Metal Machine name is missing")
	}
	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Group not specified")
	}

	bmmOut := wssdcloudcompute.BareMetalMachine{
		Name:      *bmm.Name,
		GroupName: group,
		Tags:      getWssdTags(bmm.Tags),
	}

	if bmm.BareMetalMachineProperties != nil {
		storageConfig, err := c.getWssdBareMetalMachineStorageConfiguration(bmm.StorageProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Storage Configuration")
		}
		securityConfig, err := c.getWssdBareMetalMachineSecurityConfiguration(bmm)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get Security Configuration")
		}
		osConfig, err := c.getWssdBareMetalMachineOSConfiguration(bmm.OsProfile)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get OS Configuration")
		}

		bmmOut.Storage = storageConfig
		bmmOut.Security = securityConfig
		bmmOut.Os = osConfig

		if bmm.FQDN != nil {
			bmmOut.Fqdn = *bmm.FQDN
		}
	}

	if bmm.Version != nil {
		if bmmOut.Status == nil {
			bmmOut.Status = status.InitStatus()
		}
		bmmOut.Status.Version.Number = *bmm.Version
	}

	if bmm.Location != nil {
		bmmOut.LocationName = *bmm.Location
	}

	return &bmmOut, nil
}

func (c *client) getWssdBareMetalMachineStorageConfiguration(s *compute.BareMetalMachineStorageProfile) (*wssdcloudcompute.BareMetalMachineStorageConfiguration, error) {
	wssdstorage := &wssdcloudcompute.BareMetalMachineStorageConfiguration{}

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

	return wssdstorage, nil
}

func (c *client) getWssdBareMetalMachineStorageConfigurationImageReference(s *compute.BareMetalMachineImageReference) (string, error) {
	if s.Name == nil {
		return "", errors.Wrapf(errors.InvalidInput, "Invalid Image Reference Name")
	}
	return *s.Name, nil
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

func (c *client) getBareMetalMachine(bmm *wssdcloudcompute.BareMetalMachine, group string) *compute.BareMetalMachine {
	return &compute.BareMetalMachine{
		Name: &bmm.Name,
		ID:   &bmm.Id,
		Tags: getComputeTags(bmm.GetTags()),
		BareMetalMachineProperties: &compute.BareMetalMachineProperties{
			ProvisioningState: status.GetProvisioningState(bmm.GetStatus().GetProvisioningStatus()),
			Statuses:          status.GetStatuses(bmm.GetStatus()),
			StorageProfile:    c.getBareMetalMachineStorageProfile(bmm.Storage),
			SecurityProfile:   c.getBareMetalMachineSecurityProfile(bmm),
			OsProfile:         c.getBareMetalMachineOSProfile(bmm.Os),
			FQDN:              &bmm.Fqdn,
		},
		Version:  &bmm.Status.Version.Number,
		Location: &bmm.LocationName,
	}
}

func (c *client) getBareMetalMachineStorageProfile(s *wssdcloudcompute.BareMetalMachineStorageConfiguration) *compute.BareMetalMachineStorageProfile {
	if s == nil {
		return nil
	}
	return &compute.BareMetalMachineStorageProfile{
		ImageReference: c.getBareMetalMachineStorageProfileImageReference(s.ImageReference),
	}
}

func (c *client) getBareMetalMachineStorageProfileImageReference(imageReference string) *compute.BareMetalMachineImageReference {
	return &compute.BareMetalMachineImageReference{
		Name: &imageReference,
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

func (c *client) getBareMetalMachineLinuxConfiguration(linuxConfiguration *wssdcloudcompute.LinuxConfiguration) *compute.LinuxConfiguration {
	lc := &compute.LinuxConfiguration{}

	if linuxConfiguration != nil {
		lc.DisablePasswordAuthentication = &linuxConfiguration.DisablePasswordAuthentication
	}

	return lc
}

func (c *client) getBareMetalMachineOSProfile(osConfiguration *wssdcloudcompute.BareMetalMachineOperatingSystemConfiguration) *compute.BareMetalMachineOSProfile {
	if osConfiguration == nil {
		return nil
	}
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
