// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

// Based on https://godoc.org/github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2015-06-15/compute

package compute

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/microsoft/moc/rpc/common"
)

// SubResource ...
type SubResource struct {
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
}

type OperatingSystemTypes string

const (
	// Linux
	Linux OperatingSystemTypes = "Linux"
	// Windows
	Windows OperatingSystemTypes = "Windows"
)

type OperatingSystemBootstrapEngine string

const (
	CloudInit          OperatingSystemBootstrapEngine = "CloudInit"
	WindowsAnswerFiles OperatingSystemBootstrapEngine = "WindowsAnswerFiles"
)

type VMType string

const (
	Tenant              VMType = "Tenant"
	LoadBalancer        VMType = "LoadBalancer"
	StackedControlPlane VMType = "StackedControlPlane"
)

// IPVersion enumerates the values for ip version.
type IPVersion string

const (
	// IPv4 ...
	IPv4 IPVersion = "IPv4"
	// IPv6 ...
	IPv6 IPVersion = "IPv6"
)

// ImageReference specifies information about the image to use. You can specify information about platform
// images, marketplace images, or virtual machine images. This element is required when you want to use a
// platform image, marketplace image, or virtual machine image, but is not used in other creation
// operations.
type ImageReference struct {
	// Publisher - The image publisher.
	Publisher *string `json:"publisher,omitempty"`
	// Offer - Specifies the offer of the platform image or marketplace image used to create the virtual machine.
	Offer *string `json:"offer,omitempty"`
	// Sku - The image SKU.
	Sku *string `json:"sku,omitempty"`
	// Version - Specifies the version of the platform image or marketplace image used to create the virtual machine. The allowed formats are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers. Specify 'latest' to use the latest version of an image available at deploy time. Even if you use 'latest', the VM image will not automatically update after deploy time even if a new version becomes available.
	Version *string `json:"version,omitempty"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
	// Name - Name of the image
	Name *string `json:"name,omitempty"`
}

// VirtualHardDisk describes the uri of a disk.
type VirtualHardDisk struct {
	// URI - Specifies the virtual hard disk's uri.
	URI *string `json:"uri,omitempty"`
}

type OSDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Vhd - The virtual hard disk.
	Vhd *VirtualHardDisk `json:"vhd,omitempty"`
	// Image - The source user image virtual hard disk. The virtual hard disk will be copied before being attached to the virtual machine. If SourceImage is provided, the destination virtual hard drive must not exist.
	Image *VirtualHardDisk `json:"image,omitempty"`
}

type DataDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Vhd - The virtual hard disk.
	Vhd *VirtualHardDisk `json:"vhd,omitempty"`
	// Image - The source user image virtual hard disk. The virtual hard disk will be copied before being attached to the virtual machine. If SourceImage is provided, the destination virtual hard drive must not exist.

	// ImageReference
	ImageReference *ImageReference `json:"imageReference,omitempty"`
}

type StorageProfile struct {
	// ImageReference - Specifies information about the image to use. You can specify information about platform images, marketplace images, or virtual machine images. This element is required when you want to use a platform image, marketplace image, or virtual machine image, but is not used in other creation operations.
	ImageReference *ImageReference `json:"imagereference,omitempty"`
	// OSDisk
	OsDisk *OSDisk `json:"osdisk,omitempty"`
	// DataDisks
	DataDisks *[]DataDisk `json:"datadisks,omitempty"`
	// VMConfigContainerName - Name of the storage container that hosts the VM configuration file
	VmConfigContainerName *string `json:"vmConfigContainerName,omitempty"`
}
type SSHPublicKey struct {
	// Path - Specifies the full path on the created VM where ssh public key is stored. If the file already exists, the specified key is appended to the file. Example: /home/user/.ssh/authorized_keys
	Path *string `json:"path,omitempty"`
	// KeyData - SSH public key certificate used to authenticate with the VM through ssh. The key needs to be at least 2048-bit and in ssh-rsa format. <br><br> For creating ssh keys, see [Create SSH keys on Linux and Mac for Li      nux VMs in Azure](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-mac-create-ssh-keys?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json).
	KeyData *string `json:"keyData,omitempty"`
}

type SSHConfiguration struct {
	// PublicKeys - The list of SSH public keys used to authenticate with linux based VMs.
	PublicKeys *[]SSHPublicKey `json:"publicKeys,omitempty"`
}

type RDPConfiguration struct {
	// Set to 'true' to disable Remote Desktop
	DisableRDP *bool
	// Specify custom port for Remote Desktop
	Port *uint16
}

// ProtocolTypes enumerates the values for protocol types.
type ProtocolTypes string

const (
	// HTTP ...
	HTTP ProtocolTypes = "Http"
	// HTTPS ...
	HTTPS ProtocolTypes = "Https"
)

// WinRMConfiguration describes Windows Remote Management configuration of the VM
type WinRMConfiguration struct {
	// Listeners - The list of Windows Remote Management listeners
	Listeners *[]WinRMListener `json:"listeners,omitempty"`
}

// WinRMListener describes Protocol and thumbprint of Windows Remote Management listener
type WinRMListener struct {
	// Protocol - Specifies the protocol of WinRM listener. Possible values include: 'HTTP', 'HTTPS'
	Protocol ProtocolTypes `json:"protocol,omitempty"`
}

// Based on https://godoc.org/github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2015-06-15/compute
type WindowsConfiguration struct {
	// EnableAutomaticUpdates
	EnableAutomaticUpdates *bool `json:"enableAutomaticUpdates,omitempty"`
	// TimeZone
	TimeZone *string `json:"timeZone,omitempty"`
	// AdditionalUnattendContent
	// AdditionalUnattendContent *[]AdditionalUnattendContent `json:"additionalUnattendContent,omitempty"`
	// SSH
	SSH *SSHConfiguration `json:"ssh,omitempty"`
	// RDP
	RDP *RDPConfiguration `json:"rdp,omitempty"`
	// WinRM - Specifies the Windows Remote Management listeners. This enables remote Windows PowerShell.
	WinRM *WinRMConfiguration `json:"winRM,omitempty"`
}

// Based on https://godoc.org/github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2015-06-15/compute#LinuxConfiguration
type LinuxConfiguration struct {
	// SSH
	SSH *SSHConfiguration `json:"ssh,omitempty"`
	// DisablePasswordAuthentication
	DisablePasswordAuthentication *bool `json:"disablePasswordAuthentication,omitempty"`
}

type OSProfile struct {
	// ComputerName
	ComputerName *string `json:"computername,omitempty"`
	// AdminUsername
	AdminUsername *string `json:"adminusername,omitempty"`
	// AdminPassword
	AdminPassword *string `json:"adminpassword,omitempty"`
	// CustomData Specifies a base-64 encoded string of custom data. The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes. <br><br> For using cloud-init for your VM, see [Using cloud-init to customize a Linux VM during creation](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-using-cloud-init?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json)
	CustomData *string `json:"customdata,omitempty"`
	// OsType
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// WindowsConfiguration
	WindowsConfiguration *WindowsConfiguration `json:"windowsconfiguration,omitempty"`
	// LinuxConfiguration
	LinuxConfiguration *LinuxConfiguration `json:"linuxconfiguration,omitempty"`
	// Bootstrap engine
	OsBootstrapEngine OperatingSystemBootstrapEngine `json:"osbootstrapengine,omitempty"`
}

// VirtualMachineCustomSize Specifies cpu/memory information for custom VMSize types.
type VirtualMachineCustomSize struct {
	CpuCount *int32 `json:"cpucount,omitempty"`
	MemoryMB *int32 `json:"memorymb,omitempty"`
}

// DynamicMemoryConfiguration Specifies the dynamic memory configuration for a VM.
type DynamicMemoryConfiguration struct {
	// MaximumMemoryMB - Specifies the maximum amount of memory the VM is allowed to use.
	MaximumMemoryMB *uint64 `json:"maximummemorymb,omitempty"`
	// MinimumMemoryMB - Specifies the minimum amount of memory the VM is allocated.
	MinimumMemoryMB *uint64 `json:"minimummemorymb,omitempty"`
	// TargetMemoryBuffer - Specifies the size of the VMs memory buffer as a percentage of the current memory usage.
	TargetMemoryBuffer *uint32 `json:"targetmemorybuffer,omitempty"`
}

type HardwareProfile struct {
	VMSize     VirtualMachineSizeTypes   `json:"vmsize,omitempty"`
	CustomSize *VirtualMachineCustomSize `json:"customsize,omitempty"`
	// DynamicMemoryConfig - Specifies the dynamic memory configuration for a VM, dynamic memory will be enabled if this field is present.
	DynamicMemoryConfig *DynamicMemoryConfiguration `json:"dynamicmemoryconfig,omitempty"`
}

// NetworkInterfaceReferenceProperties describes a network interface reference properties.
type NetworkInterfaceReferenceProperties struct {
	// Primary - Specifies the primary network interface in case the virtual machine has more than 1 network interface.
	Primary *bool `json:"primary,omitempty"`
}

type NetworkInterfaceReference struct {
	*NetworkInterfaceReferenceProperties `json:"properties,omitempty"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
}

type NetworkProfile struct {
	// NetworkInterfaces
	NetworkInterfaces *[]NetworkInterfaceReference `json:"networkinterfaces,omitempty"`
}

type GuestAgentProfile struct {
	// Enabled - Specifies whether guest agent should be enabled on the virtual machine.
	Enabled *bool `json:"enabled,omitempty"`
}

type InstanceViewStatus struct {
	// Code - READ-ONLY; The status code, which only appears in the response.
	Code string `json:"code,omitempty"`
	// Level - READ-ONLY; The level code, which only appears in the response.
	Level string `json:"level,omitempty"`
	// DisplayStatus - READ-ONLY; The short localizable label for the status, which only appears in the response.
	DisplayStatus string `json:"displayStatus,omitempty"`
	// Message - READ-ONLY; The detailed status message, including for alerts and error messages, which only appears in the response.
	Message string `json:"message,omitempty"`
	// Time - READ-ONLY; The time of the status, which only appears in the response.
	Time string `json:"time,omitempty"`
}

type GuestAgentInstanceView struct {
	// AgentVersion - READ-ONLY; The Guest Agent full version, which only appears in the response.
	AgentVersion string `json:"agentVersion,omitempty"`
	// Statuses - READ-ONLY; The resource status information, which only appears in the response.
	Statuses []*InstanceViewStatus `json:"statuses,omitempty"`
}

type UefiSettings struct {
	// SecureBootEnabled - Specifies whether secure boot should be enabled on the virtual machine.
	SecureBootEnabled *bool `json:"secureBootEnabled,omitempty"`
}

type SecurityProfile struct {
	EnableTPM *bool `json:"enableTPM,omitempty"`
	//Security related configuration used while creating the virtual machine.
	UefiSettings *UefiSettings `json:"uefiSettings,omitempty"`
}

// Plan specifies information about the marketplace image used to create the virtual machine. This element
// is only used for marketplace images. Before you can use a marketplace image from an API, you must enable
// the image for programmatic use.  In the Azure portal, find the marketplace image that you want to use
// and then click **Want to deploy programmatically, Get Started ->**. Enter any required information and
// then click **Save**.
type Plan struct {
	// Name - The plan ID.
	Name *string `json:"name,omitempty"`
	// Publisher - The publisher ID.
	Publisher *string `json:"publisher,omitempty"`
	// Product - Specifies the product of the image from the marketplace. This is the same value as Offer under the imageReference element.
	Product *string `json:"product,omitempty"`
	// PromotionCode - The promotion code.
	PromotionCode *string `json:"promotionCode,omitempty"`
}

// VirtualMachineProperties describes the properties of a Virtual Machine.
type VirtualMachineProperties struct {
	// StorageProfile
	StorageProfile *StorageProfile `json:"storageprofile,omitempty"`
	// OsProfile
	OsProfile *OSProfile `json:"osprofile,omitempty"`
	// NetworkProfile
	NetworkProfile *NetworkProfile `json:"networkprofile,omitempty"`
	// HardwareProfile - Specifies the hardware settings for the virtual machine.
	HardwareProfile *HardwareProfile `json:"hardwareprofile,omitempty"`
	// GuestAgentProfile - Specifies the guest agent settings for the virtual machine.
	GuestAgentProfile *GuestAgentProfile `json:"guestAgentProfile,omitempty"`
	// SecurityProfile - Specifies the security settings for the virtual machine.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// Host - Specifies information about the dedicated host that the virtual machine resides in. <br><br>Minimum api-version: 2018-10-01.
	Host *SubResource `json:"host,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// ValidationState - READ-ONLY; The validation status, which only appears in the response.
	ValidationStatus []*common.ValidationState `json:"validationStatus"`
	// GuestAgentInstanceView - READ-ONLY; The info of the Agent running on the virtual machine, which only appears in the response.
	GuestAgentInstanceView *GuestAgentInstanceView `json:"guestAgentInstanceView,omitempty"`
	// VMID - READ-ONLY; Specifies the VM unique ID which is a 128-bits identifier that is encoded and stored in all Azure IaaS VMs SMBIOS and can be read using platform BIOS commands.
	VMID *string `json:"vmId,omitempty"`
	// VmType - The type of the VM.  Can be either tenant or loadbalancer vm
	VmType VMType `json:"vmType,omitempty"`
	// Disable High Availability
	DisableHighAvailability *bool `json:"disableHighAvailability,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

type VirtualMachine struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Zones - The virtual machine scale set zones.
	Zones *[]string `json:"zones,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Plan
	Plan *Plan `json:"plan,omitempty"`
	// Properties
	*VirtualMachineProperties `json:"virtualmachineproperties,omitempty"`
}

type Sku struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Capacity
	Capacity *int64 `json:"capacity,omitempty"`
}

// APIEntityReference the API entity reference.
type APIEntityReference struct {
	// ID - The ARM resource id in the form of /subscriptions/{SubscriptionId}/resourceGroups/{ResourceGroupName}/...
	ID *string `json:"id,omitempty"`
}

// VirtualMachineScaleSetNetworkConfigurationDNSSettings describes a virtual machines scale sets network
// configuration's DNS settings.
type VirtualMachineScaleSetNetworkConfigurationDNSSettings struct {
	// DNSServers - List of DNS servers IP addresses
	DNSServers *[]string `json:"dnsServers,omitempty"`
}

// VirtualMachineScaleSetIPTag contains the IP tag associated with the public IP address.
type VirtualMachineScaleSetIPTag struct {
	// IPTagType - IP tag type. Example: FirstPartyUsage.
	IPTagType *string `json:"ipTagType,omitempty"`
	// Tag - IP tag associated with the public IP. Example: SQL, Storage etc.
	Tag *string `json:"tag,omitempty"`
}

// VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings describes a virtual machines scale sets
// network configuration's DNS settings.
type VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings struct {
	// DomainNameLabel - The Domain name label.The concatenation of the domain name label and vm index will be the domain name labels of the PublicIPAddress resources that will be created
	DomainNameLabel *string `json:"domainNameLabel,omitempty"`
}

// VirtualMachineScaleSetPublicIPAddressConfigurationProperties describes a virtual machines scale set IP
// Configuration's PublicIPAddress configuration
type VirtualMachineScaleSetPublicIPAddressConfigurationProperties struct {
	// IdleTimeoutInMinutes - The idle timeout of the public IP address.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// DNSSettings - The dns settings to be applied on the publicIP addresses .
	DNSSettings *VirtualMachineScaleSetPublicIPAddressConfigurationDNSSettings `json:"dnsSettings,omitempty"`
	// IPTags - The list of IP tags associated with the public IP address.
	IPTags *[]VirtualMachineScaleSetIPTag `json:"ipTags,omitempty"`
	// PublicIPPrefix - The PublicIPPrefix from which to allocate publicIP addresses.
	PublicIPPrefix *SubResource `json:"publicIPPrefix,omitempty"`
}

// VirtualMachineScaleSetPublicIPAddressConfiguration describes a virtual machines scale set IP
// Configuration's PublicIPAddress configuration
type VirtualMachineScaleSetPublicIPAddressConfiguration struct {
	// Name - The publicIP address configuration name.
	Name                                                          *string `json:"name,omitempty"`
	*VirtualMachineScaleSetPublicIPAddressConfigurationProperties `json:"properties,omitempty"`
}

// VirtualMachineScaleSetIPConfigurationProperties describes a virtual machine scale set network profile's
// IP configuration properties.
type VirtualMachineScaleSetIPConfigurationProperties struct {
	// Subnet - Specifies the identifier of the subnet.
	Subnet *APIEntityReference `json:"subnet,omitempty"`
	// Primary - Specifies the primary network interface in case the virtual machine has more than 1 network interface.
	Primary *bool `json:"primary,omitempty"`
	// PublicIPAddressConfiguration - The publicIPAddressConfiguration.
	PublicIPAddressConfiguration *VirtualMachineScaleSetPublicIPAddressConfiguration `json:"publicIPAddressConfiguration,omitempty"`
	// PrivateIPAddressVersion - Available from Api-Version 2017-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.  Possible values are: 'IPv4' and 'IPv6'. Possible values include: 'IPv4', 'IPv6'
	PrivateIPAddressVersion IPVersion `json:"privateIPAddressVersion,omitempty"`
	// ApplicationGatewayBackendAddressPools - Specifies an array of references to backend address pools of application gateways. A scale set can reference backend address pools of multiple application gateways. Multiple scale sets cannot use the same application gateway.
	ApplicationGatewayBackendAddressPools *[]SubResource `json:"applicationGatewayBackendAddressPools,omitempty"`
	// ApplicationSecurityGroups - Specifies an array of references to application security group.
	ApplicationSecurityGroups *[]SubResource `json:"applicationSecurityGroups,omitempty"`
	// LoadBalancerBackendAddressPools - Specifies an array of references to backend address pools of load balancers. A scale set can reference backend address pools of one public and one internal load balancer. Multiple scale sets cannot use the same load balancer.
	LoadBalancerBackendAddressPools *[]SubResource `json:"loadBalancerBackendAddressPools,omitempty"`
	// LoadBalancerInboundNatPools - Specifies an array of references to inbound Nat pools of the load balancers. A scale set can reference inbound nat pools of one public and one internal load balancer. Multiple scale sets cannot use the same load balancer
	LoadBalancerInboundNatPools *[]SubResource `json:"loadBalancerInboundNatPools,omitempty"`
}

// VirtualMachineScaleSetIPConfiguration describes a virtual machine scale set network profile's IP
// configuration.
type VirtualMachineScaleSetIPConfiguration struct {
	// Name - The IP configuration name.
	Name                                             *string `json:"name,omitempty"`
	*VirtualMachineScaleSetIPConfigurationProperties `json:"properties,omitempty"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
}

// VirtualMachineScaleSetNetworkConfigurationProperties describes a virtual machine scale set network
// profile's IP configuration.
type VirtualMachineScaleSetNetworkConfigurationProperties struct {
	// Primary - Specifies the primary network interface in case the virtual machine has more than 1 network interface.
	Primary *bool `json:"primary,omitempty"`
	// NetworkSecurityGroup - The network security group.
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`
	// DNSSettings - The dns settings to be applied on the network interfaces.
	DNSSettings *VirtualMachineScaleSetNetworkConfigurationDNSSettings `json:"dnsSettings,omitempty"`
	// IPConfigurations - Specifies the IP configurations of the network interface.
	IPConfigurations *[]VirtualMachineScaleSetIPConfiguration `json:"ipConfigurations,omitempty"`
	// EnableIPForwarding - Whether IP forwarding enabled on this NIC.
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`
	// VirtualNetworkName - TODO: Remove this
	// VirtualNetworkName *string `json:"virtualNetworkName,omitempty"`
}

// VirtualMachineScaleSetNetworkConfiguration describes a virtual machine scale set network profile's
// network configurations.
type VirtualMachineScaleSetNetworkConfiguration struct {
	// Name - The network configuration name.
	Name                                                  *string `json:"name,omitempty"`
	*VirtualMachineScaleSetNetworkConfigurationProperties `json:"properties,omitempty"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
}

type VirtualMachineScaleSetNetworkProfile struct {
	// HealthProbe - A reference to a load balancer probe used to determine the health of an instance in the virtual machine scale set. The reference will be in the form: '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/probes/{probeName}'.
	HealthProbe *APIEntityReference `json:"healthProbe,omitempty"`
	// NetworkInterfaceConfigurations
	NetworkInterfaceConfigurations *[]VirtualMachineScaleSetNetworkConfiguration `json:"networkInterfaceConfigurations,omitempty"`
}

// BootDiagnostics boot Diagnostics is a debugging feature which allows you to view Console Output and
// Screenshot to diagnose VM status. <br><br> You can easily view the output of your console log. <br><br>
// Azure also enables you to see a screenshot of the VM from the hypervisor.
type BootDiagnostics struct {
	// Enabled - Whether boot diagnostics should be enabled on the Virtual Machine.
	Enabled *bool `json:"enabled,omitempty"`
	// StorageURI - Uri of the storage account to use for placing the console output and screenshot.
	StorageURI *string `json:"storageUri,omitempty"`
}

type DiagnosticsProfile struct {
	// BootDiagnostics - Boot Diagnostics is a debugging feature which allows you to view Console Output and Screenshot to diagnose VM status. <br><br> You can easily view the output of your console log. <br><br> Azure also enables you to see a screenshot of the VM from the hypervisor.
	BootDiagnostics *BootDiagnostics `json:"bootDiagnostics,omitempty"`
}

// VirtualMachinePriorityTypes enumerates the values for virtual machine priority types.
type VirtualMachinePriorityTypes string

const (
	Low     VirtualMachinePriorityTypes = "Low"
	Regular VirtualMachinePriorityTypes = "Regular"
)

// VirtualMachineEvictionPolicyTypes enumerates the values for virtual machine eviction policy types.
type VirtualMachineEvictionPolicyTypes string

const (
	Deallocate VirtualMachineEvictionPolicyTypes = "Deallocate"
	Delete     VirtualMachineEvictionPolicyTypes = "Delete"
)

// VaultCertificate describes a single certificate reference in a Key Vault, and where the certificate
// should reside on the VM.
type VaultCertificate struct {
	// CertificateURL - This is the URL of a certificate that has been uploaded to Key Vault as a secret. For adding a secret to the Key Vault, see [Add a key or secret to the key vault](https://docs.microsoft.com/azure/key-vault/key-vault-get-started/#add). In this case, your certificate needs to be It is the Base64 encoding of the following JSON Object which is encoded in UTF-8: <br><br> {<br>  "data":"<Base64-encoded-certificate>",<br>  "dataType":"pfx",<br>  "password":"<pfx-file-password>"<br>}
	CertificateURL *string `json:"certificateUrl,omitempty"`
	// CertificateStore - For Windows VMs, specifies the certificate store on the Virtual Machine to which the certificate should be added. The specified certificate store is implicitly in the LocalMachine account. <br><br>For Linux VMs, the certificate file is placed under the /var/lib/waagent directory, with the file name &lt;UppercaseThumbprint&gt;.crt for the X509 certificate file and &lt;UppercaseThumbprint&gt;.prv for private key. Both of these files are .pem formatted.
	CertificateStore *string `json:"certificateStore,omitempty"`
}

// VaultSecretGroup describes a set of certificates which are all in the same Key Vault.
type VaultSecretGroup struct {
	// SourceVault - The relative URL of the Key Vault containing all of the certificates in VaultCertificates.
	SourceVault *SubResource `json:"sourceVault,omitempty"`
	// VaultCertificates - The list of key vault references in SourceVault which contain certificates.
	VaultCertificates *[]VaultCertificate `json:"vaultCertificates,omitempty"`
}

// VirtualMachineScaleSetOSProfile describes a virtual machine scale set OS profile.
type VirtualMachineScaleSetOSProfile struct {
	// ComputerNamePrefix - Specifies the computer name prefix for all of the virtual machines in the scale set. Computer name prefixes must be 1 to 15 characters long.
	ComputerNamePrefix *string `json:"computerNamePrefix,omitempty"`
	// AdminUsername - Specifies the name of the administrator account. <br><br> **Windows-only restriction:** Cannot end in "." <br><br> **Disallowed values:** "administrator", "admin", "user", "user1", "test", "user2", "test1", "user3", "admin1", "1", "123", "a", "actuser", "adm", "admin2", "aspnet", "backup", "console", "david", "guest", "john", "owner", "root", "server", "sql", "support", "support_388945a0", "sys", "test2", "test3", "user4", "user5". <br><br> **Minimum-length (Linux):** 1  character <br><br> **Max-length (Linux):** 64 characters <br><br> **Max-length (Windows):** 20 characters  <br><br><li> For root access to the Linux VM, see [Using root privileges on Linux virtual machines in Azure](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-use-root-privileges?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json)<br><li> For a list of built-in system users on Linux that should not be used in this field, see [Selecting User Names for Linux on Azure](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-usernames?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json)
	AdminUsername *string `json:"adminUsername,omitempty"`
	// AdminPassword - Specifies the password of the administrator account. <br><br> **Minimum-length (Windows):** 8 characters <br><br> **Minimum-length (Linux):** 6 characters <br><br> **Max-length (Windows):** 123 characters <br><br> **Max-length (Linux):** 72 characters <br><br> **Complexity requirements:** 3 out of 4 conditions below need to be fulfilled <br> Has lower characters <br>Has upper characters <br> Has a digit <br> Has a special character (Regex match [\W_]) <br><br> **Disallowed values:** "abc@123", "P@$$w0rd", "P@ssw0rd", "P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1", "Password22", "iloveyou!" <br><br> For resetting the password, see [How to reset the Remote Desktop service or its login password in a Windows VM](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-windows-reset-rdp?toc=%2fazure%2fvirtual-machines%2fwindows%2ftoc.json) <br><br> For resetting root password, see [Manage users, SSH, and check or repair disks on Azure Linux VMs using the VMAccess Extension](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-using-vmaccess-extension?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json#reset-root-password)
	AdminPassword *string `json:"adminPassword,omitempty"`
	// CustomData - Specifies a base-64 encoded string of custom data. The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes. <br><br> For using cloud-init for your VM, see [Using cloud-init to customize a Linux VM during creation](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-using-cloud-init?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json)
	CustomData *string `json:"customData,omitempty"`
	// WindowsConfiguration - Specifies Windows operating system settings on the virtual machine.
	WindowsConfiguration *WindowsConfiguration `json:"windowsConfiguration,omitempty"`
	// LinuxConfiguration - Specifies the Linux operating system settings on the virtual machine. <br><br>For a list of supported Linux distributions, see [Linux on Azure-Endorsed Distributions](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-endorsed-distros?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json) <br><br> For running non-endorsed distributions, see [Information for Non-Endorsed Distributions](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-create-upload-generic?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json).
	LinuxConfiguration *LinuxConfiguration `json:"linuxConfiguration,omitempty"`
	// Secrets - Specifies set of certificates that should be installed onto the virtual machines in the scale set.
	Secrets *[]VaultSecretGroup `json:"secrets,omitempty"`
	// Bootstrap engine
	OsBootstrapEngine OperatingSystemBootstrapEngine `json:"osbootstrapengine,omitempty"`
}

// DiskCreateOptionTypes enumerates the values for disk create option types.
type DiskCreateOptionTypes string

const (
	// DiskCreateOptionTypesAttach ...
	DiskCreateOptionTypesAttach DiskCreateOptionTypes = "Attach"
	// DiskCreateOptionTypesEmpty ...
	DiskCreateOptionTypesEmpty DiskCreateOptionTypes = "Empty"
	// DiskCreateOptionTypesFromImage ...
	DiskCreateOptionTypesFromImage DiskCreateOptionTypes = "FromImage"
)

// VirtualMachineScaleSetOSDisk describes a virtual machine scale set operating system disk.
type VirtualMachineScaleSetOSDisk struct {
	// Name - The disk name.
	Name *string `json:"name,omitempty"`
	// CreateOption - Specifies how the virtual machines in the scale set should be created.<br><br> The only allowed value is: **FromImage** \u2013 This value is used when you are using an image to create the virtual machine. If you are using a platform image, you also use the imageReference element described above. If you are using a marketplace image, you  also use the plan element previously described. Possible values include: 'DiskCreateOptionTypesFromImage', 'DiskCreateOptionTypesEmpty', 'DiskCreateOptionTypesAttach'
	CreateOption DiskCreateOptionTypes `json:"createOption,omitempty"`
	// DiskSizeGB - Specifies the size of the operating system disk in gigabytes. This element can be used to overwrite the size of the disk in a virtual machine image. <br><br> This value cannot be larger than 1023 GB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
	// OsType - This property allows you to specify the type of the OS that is included in the disk if creating a VM from user-image or a specialized VHD. <br><br> Possible values are: <br><br> **Windows** <br><br> **Linux**. Possible values include: 'Windows', 'Linux'
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// Image - Specifies information about the unmanaged user image to base the scale set on.
	Image *VirtualHardDisk `json:"image,omitempty"`
	// VhdContainers - Specifies the container urls that are used to store operating system disks for the scale set.
	VhdContainers *[]string `json:"vhdContainers,omitempty"`
}

// VirtualMachineScaleSetDataDisk describes a virtual machine scale set data disk.
type VirtualMachineScaleSetDataDisk struct {
	// Name - The disk name.
	Name *string `json:"name,omitempty"`
	// Lun - Specifies the logical unit number of the data disk. This value is used to identify data disks within the VM and therefore must be unique for each data disk attached to a VM.
	Lun *int32 `json:"lun,omitempty"`
	// CreateOption - The create option. Possible values include: 'DiskCreateOptionTypesFromImage', 'DiskCreateOptionTypesEmpty', 'DiskCreateOptionTypesAttach'
	CreateOption DiskCreateOptionTypes `json:"createOption,omitempty"`
	// DiskSizeGB - Specifies the size of an empty data disk in gigabytes. This element can be used to overwrite the size of the disk in a virtual machine image. <br><br> This value cannot be larger than 1023 GB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
	// Image - Specifies information about the unmanaged user image to base the scale set on.
	Image *VirtualHardDisk `json:"image,omitempty"`
}

// VirtualMachineScaleSetStorageProfile describes a virtual machine scale set storage profile.
type VirtualMachineScaleSetStorageProfile struct {
	// ImageReference - Specifies information about the image to use. You can specify information about platform images, marketplace images, or virtual machine images. This element is required when you want to use a platform image, marketplace image, or virtual machine image, but is not used in other creation operations.
	ImageReference *ImageReference `json:"imageReference,omitempty"`
	// OsDisk - Specifies information about the operating system disk used by the virtual machines in the scale set. <br><br> For more information about disks, see [About disks and VHDs for Azure virtual machines](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-windows-about-disks-vhds?toc=%2fazure%2fvirtual-machines%2fwindows%2ftoc.json).
	OsDisk *VirtualMachineScaleSetOSDisk `json:"osDisk,omitempty"`
	// DataDisks - Specifies the parameters that are used to add data disks to the virtual machines in the scale set. <br><br> For more information about disks, see [About disks and VHDs for Azure virtual machines](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-windows-about-disks-vhds?toc=%2fazure%2fvirtual-machines%2fwindows%2ftoc.json).
	DataDisks *[]VirtualMachineScaleSetDataDisk `json:"dataDisks,omitempty"`
	// VMConfigContainerName - Name of the storage container that hosts the VM configuration file
	VmConfigContainerName *string `json:"vmConfigContainerName,omitempty"`
}

// VirtualMachineScaleSetHardwareProfile describes a virtual machine scale set storage profile.
type VirtualMachineScaleSetHardwareProfile struct {
	// VMSize - Specifies the size of the virtual machine.
	VMSize VirtualMachineSizeTypes `json:"vmSize,omitempty"`
	// CustomSize - Specifies cpu/memory information for custom VMSize types.
	CustomSize *VirtualMachineCustomSize `json:"customsize,omitempty"`
}

// VirtualMachineScaleSetVMProfile describes a virtual machine scale set virtual machine profile.
type VirtualMachineScaleSetVMProfile struct {
	// StorageProfile
	StorageProfile *VirtualMachineScaleSetStorageProfile `json:"storageProfile,omitempty"`
	// HardwareProfile
	HardwareProfile *VirtualMachineScaleSetHardwareProfile `json:"hardwareProfile,omitempty"`
	// GuestAgentProfile - Specifies the guest agent settings for the virtual machine.
	GuestAgentProfile *GuestAgentProfile `json:"guestAgentProfile,omitempty"`
	// SecurityProfile - Specifies the security settings for the virtual machine.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// OsProfile
	OsProfile *VirtualMachineScaleSetOSProfile `json:"osProfile,omitempty"`
	// NetworkProfile
	NetworkProfile *VirtualMachineScaleSetNetworkProfile `json:"networkProfile,omitempty"`
	// DiagnosticsProfile - Specifies the boot diagnostic settings state
	DiagnosticsProfile *DiagnosticsProfile `json:"diagnosticsProfile,omitempty"`
	// Priority - Specifies the priority for the virtual machines in the scale set. <br><br>Minimum api-version: 2017-10-30-preview. Possible values include: 'Regular', 'Low'
	Priority VirtualMachinePriorityTypes `json:"priority,omitempty"`
	// EvictionPolicy - Specifies the eviction policy for virtual machines in a low priority scale set. <br><br>Minimum api-version: 2017-10-30-preview. Possible values include: 'Deallocate', 'Delete'
	EvictionPolicy VirtualMachineEvictionPolicyTypes `json:"evictionPolicy,omitempty"`
	// Disable High Availability
	DisableHighAvailability *bool `json:"disableHighAvailability,omitempty"`
}

// ResourceIdentityType enumerates the values for resource identity type.
type ResourceIdentityType string

const (
	// ResourceIdentityTypeNone ...
	ResourceIdentityTypeNone ResourceIdentityType = "None"
	// ResourceIdentityTypeSystemAssigned ...
	ResourceIdentityTypeSystemAssigned ResourceIdentityType = "SystemAssigned"
	// ResourceIdentityTypeSystemAssignedUserAssigned ...
	ResourceIdentityTypeSystemAssignedUserAssigned ResourceIdentityType = "SystemAssigned, UserAssigned"
	// ResourceIdentityTypeUserAssigned ...
	ResourceIdentityTypeUserAssigned ResourceIdentityType = "UserAssigned"
)

// VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue ...
type VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue struct {
	// PrincipalID - READ-ONLY; The principal id of user assigned identity.
	PrincipalID *string `json:"principalId,omitempty"`
	// ClientID - READ-ONLY; The client id of user assigned identity.
	ClientID *string `json:"clientId,omitempty"`
}

// VirtualMachineScaleSetIdentity identity for the virtual machine scale set.
type VirtualMachineScaleSetIdentity struct {
	// PrincipalID - READ-ONLY; The principal id of virtual machine scale set identity. This property will only be provided for a system assigned identity.
	PrincipalID *string `json:"principalId,omitempty"`
	// TenantID - READ-ONLY; The tenant id associated with the virtual machine scale set. This property will only be provided for a system assigned identity.
	TenantID *string `json:"tenantId,omitempty"`
	// Type - The type of identity used for the virtual machine scale set. The type 'SystemAssigned, UserAssigned' includes both an implicitly created identity and a set of user assigned identities. The type 'None' will remove any identities from the virtual machine scale set. Possible values include: 'ResourceIdentityTypeSystemAssigned', 'ResourceIdentityTypeUserAssigned', 'ResourceIdentityTypeSystemAssignedUserAssigned', 'ResourceIdentityTypeNone'
	Type ResourceIdentityType `json:"type,omitempty"`
	// UserAssignedIdentities - The list of user identities associated with the virtual machine scale set. The user identity dictionary key references will be ARM resource ids in the form: '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
	UserAssignedIdentities map[string]*VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue `json:"userAssignedIdentities"`
}

// VirtualMachineScaleSetProperties describes the properties of a Virtual Machine Scale Set.
type VirtualMachineScaleSetProperties struct {
	// VirtualMachineProfile - The virtual machine profile.
	VirtualMachineProfile *VirtualMachineScaleSetVMProfile `json:"virtualMachineProfile,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// VirtualMachineScaleSet
type VirtualMachineScaleSet struct {
	autorest.Response `json:"-"`
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Zones - The virtual machine scale set zones.
	Zones *[]string `json:"zones,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Plan - Specifies information about the marketplace image used to create the virtual machine. This element is only used for marketplace images. Before you can use a marketplace image from an API, you must enable the image for programmatic use.  In the Azure portal, find the marketplace image that you want to use and then click **Want to deploy programmatically, Get Started ->**. Enter any required information and then click **Save**.
	Plan *Plan `json:"plan,omitempty"`

	*VirtualMachineScaleSetProperties `json:"properties,omitempty"`
	// Sku
	Sku *Sku `json:"sku,omitempty"`
	// Identity - The identity of the virtual machine scale set, if configured.
	Identity *VirtualMachineScaleSetIdentity `json:"identity,omitempty"`
}

// OperatingSystemStateTypes enumerates the values for operating system state types.
type OperatingSystemStateTypes string

const (
	// Generalized ...
	Generalized OperatingSystemStateTypes = "Generalized"
	// Specialized ...
	Specialized OperatingSystemStateTypes = "Specialized"
)

// GalleryImageIdentifier this is the gallery Image Definition identifier.
type GalleryImageIdentifier struct {
	// Publisher - The name of the gallery Image Definition publisher.
	Publisher *string `json:"publisher,omitempty"`
	// Offer - The name of the gallery Image Definition offer.
	Offer *string `json:"offer,omitempty"`
	// Sku - The name of the gallery Image Definition SKU.
	Sku *string `json:"sku,omitempty"`
}

// Disallowed describes the disallowed disk types.
type Disallowed struct {
	// DiskTypes - A list of disk types.
	DiskTypes *[]string `json:"diskTypes,omitempty"`
}

// ImagePurchasePlan describes the gallery Image Definition purchase plan. This is used by marketplace
// images.
type ImagePurchasePlan struct {
	// Name - The plan ID.
	Name *string `json:"name,omitempty"`
	// Publisher - The publisher ID.
	Publisher *string `json:"publisher,omitempty"`
	// Product - The product ID.
	Product *string `json:"product,omitempty"`
}

// ResourceRange describes the resource range.
type ResourceRange struct {
	// Min - The minimum number of the resource.
	Min *int32 `json:"min,omitempty"`
	// Max - The maximum number of the resource.
	Max *int32 `json:"max,omitempty"`
}

// RecommendedMachineConfiguration the properties describe the recommended machine configuration for this
// Image Definition. These properties are updatable.
type RecommendedMachineConfiguration struct {
	VCPUs  *ResourceRange `json:"vCPUs,omitempty"`
	Memory *ResourceRange `json:"memory,omitempty"`
}

// ProvisioningState2 enumerates the values for provisioning state 2.
type ProvisioningState2 string

const (
	// ProvisioningState2Creating ...
	ProvisioningState2Creating ProvisioningState2 = "Creating"
	// ProvisioningState2Deleting ...
	ProvisioningState2Deleting ProvisioningState2 = "Deleting"
	// ProvisioningState2Failed ...
	ProvisioningState2Failed ProvisioningState2 = "Failed"
	// ProvisioningState2Migrating ...
	ProvisioningState2Migrating ProvisioningState2 = "Migrating"
	// ProvisioningState2Succeeded ...
	ProvisioningState2Succeeded ProvisioningState2 = "Succeeded"
	// ProvisioningState2Updating ...
	ProvisioningState2Updating ProvisioningState2 = "Updating"
)

// Http Image properties
type HttpImageProperties struct {
	HttpURL string `json:"httpURL,omitempty"`
}

// SFSImage properties
type SFSImageProperties struct {
	Version     string `json:"version,omitempty"`
	ReleaseName string `json:"releasename,omitempty"`
	Parts       int32  `json:"parts,omitempty"`
}

// Local image properties
type LocalImageProperties struct {
	Path string `json:"path,omitempty"`
}

type CloneImageProperties struct {
	CloneSource string `json:"cloneSource,omitempty"`
}

// Azure GalleryImage properties
type AzureGalleryImageProperties struct {
	// SasURL - Specifies the SAS URI for the image
	SasURI string `json:"sasURI,omitempty"`
	// Version - Specifies the version of the platform image or marketplace image used to create the virtual machine.
	Version string `json:"version,omitempty"`
}

// GalleryImageProperties describes the properties of a gallery Image Definition.
type GalleryImageProperties struct {
	// Description - The description of this gallery Image Definition resource. This property is updatable.
	Description *string `json:"description,omitempty"`
	// Eula - The Eula agreement for the gallery Image Definition.
	Eula *string `json:"eula,omitempty"`
	// PrivacyStatementURI - The privacy statement uri.
	PrivacyStatementURI *string `json:"privacyStatementUri,omitempty"`
	// ReleaseNoteURI - The release note uri.
	ReleaseNoteURI *string `json:"releaseNoteUri,omitempty"`
	// OsType - This property allows you to specify the type of the OS that is included in the disk when creating a VM from a managed image. <br><br> Possible values are: <br><br> **Windows** <br><br> **Linux**. Possible values include: 'Windows', 'Linux'
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// OsState - This property allows the user to specify whether the virtual machines created under this image are 'Generalized' or 'Specialized'. Possible values include: 'Generalized', 'Specialized'
	OsState OperatingSystemStateTypes `json:"osState,omitempty"`
	// HyperVGeneration - The hypervisor generation of the Virtual Machine. Applicable to OS disks only. Possible values include: 'V1', 'V2'
	HyperVGeneration common.HyperVGeneration `json:"hyperVGeneration,omitempty"`
	// EndOfLifeDate - The end of life date of the gallery Image Definition. This property can be used for decommissioning purposes. This property is updatable.
	EndOfLifeDate *date.Time                       `json:"endOfLifeDate,omitempty"`
	Identifier    *GalleryImageIdentifier          `json:"identifier,omitempty"`
	Recommended   *RecommendedMachineConfiguration `json:"recommended,omitempty"`
	Disallowed    *Disallowed                      `json:"disallowed,omitempty"`
	PurchasePlan  *ImagePurchasePlan               `json:"purchasePlan,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response. Possible values include: 'ProvisioningState2Creating', 'ProvisioningState2Updating', 'ProvisioningState2Failed', 'ProvisioningState2Succeeded', 'ProvisioningState2Deleting', 'ProvisioningState2Migrating'
	ProvisioningState ProvisioningState2 `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// Container name
	ContainerName *string `json:"containername,omitempty"`
	//Type of source of gal image (sfs/http/local)
	SourceType common.ImageSource `json:"sourceType,omitempty"`
	// CloudInitDataSource - READ-ONLY; The cloud init data source to be used with the image. [NoCloud, Azure]. Default Value – NoCloud. For marketplace images it will be Azure.
	CloudInitDataSource common.CloudInitDataSource `json:"cloudInitDataSource,omitempty"`
}

// GalleryImage specifies information about the gallery Image Definition that you want to create or update.
type GalleryImage struct {
	autorest.Response       `json:"-"`
	*GalleryImageProperties `json:"properties,omitempty"`
	// ID - READ-ONLY; Resource Id
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags
	Tags map[string]*string `json:"tags"`
}

// CachingTypes enumerates the values for caching types.
type CachingTypes string

const (
	// CachingTypesNone ...
	CachingTypesNone CachingTypes = "None"
	// CachingTypesReadOnly ...
	CachingTypesReadOnly CachingTypes = "ReadOnly"
	// CachingTypesReadWrite ...
	CachingTypesReadWrite CachingTypes = "ReadWrite"
)

// StorageAccountTypes enumerates the values for storage account types.
type StorageAccountTypes string

const (
	// StorageAccountTypesPremiumLRS ...
	StorageAccountTypesPremiumLRS StorageAccountTypes = "Premium_LRS"
	// StorageAccountTypesStandardLRS ...
	StorageAccountTypesStandardLRS StorageAccountTypes = "Standard_LRS"
	// StorageAccountTypesStandardSSDLRS ...
	StorageAccountTypesStandardSSDLRS StorageAccountTypes = "StandardSSD_LRS"
	// StorageAccountTypesUltraSSDLRS ...
	StorageAccountTypesUltraSSDLRS StorageAccountTypes = "UltraSSD_LRS"
)

// VirtualMachineImageOSDisk describes an Operating System disk.
type VirtualMachineImageOSDisk struct {
	// OsType - This property allows you to specify the type of the OS that is included in the disk if creating a VM from a custom image. <br><br> Possible values are: <br><br> **Windows** <br><br> **Linux**. Possible values include: 'Windows', 'Linux'
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// OsState - The OS State. Possible values include: 'Generalized', 'Specialized'
	OsState OperatingSystemStateTypes `json:"osState,omitempty"`
	// Snapshot - The snapshot.
	Snapshot *SubResource `json:"snapshot,omitempty"`
	// ManagedDisk - The managedDisk.
	ManagedDisk *SubResource `json:"managedDisk,omitempty"`
	// BlobURI - The Virtual Hard Disk.
	BlobURI *string `json:"blobUri,omitempty"`
	// Caching - Specifies the caching requirements. <br><br> Possible values are: <br><br> **None** <br><br> **ReadOnly** <br><br> **ReadWrite** <br><br> Default: **None for Standard storage. ReadOnly for Premium storage**. Possible values include: 'CachingTypesNone', 'CachingTypesReadOnly', 'CachingTypesReadWrite'
	Caching CachingTypes `json:"caching,omitempty"`
	// DiskSizeGB - Specifies the size of empty data disks in gigabytes. This element can be used to overwrite the name of the disk in a virtual machine image. <br><br> This value cannot be larger than 1023 GB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
	// StorageAccountType - Specifies the storage account type for the managed disk. UltraSSD_LRS cannot be used with OS Disk. Possible values include: 'StorageAccountTypesStandardLRS', 'StorageAccountTypesPremiumLRS', 'StorageAccountTypesStandardSSDLRS', 'StorageAccountTypesUltraSSDLRS'
	StorageAccountType StorageAccountTypes `json:"storageAccountType,omitempty"`
}

// VirtualMachineImageDataDisk describes a data disk.
type VirtualMachineImageDataDisk struct {
	// Lun - Specifies the logical unit number of the data disk. This value is used to identify data disks within the VM and therefore must be unique for each data disk attached to a VM.
	Lun *int32 `json:"lun,omitempty"`
	// Snapshot - The snapshot.
	Snapshot *SubResource `json:"snapshot,omitempty"`
	// ManagedDisk - The managedDisk.
	ManagedDisk *SubResource `json:"managedDisk,omitempty"`
	// BlobURI - The Virtual Hard Disk.
	BlobURI *string `json:"blobUri,omitempty"`
	// Caching - Specifies the caching requirements. <br><br> Possible values are: <br><br> **None** <br><br> **ReadOnly** <br><br> **ReadWrite** <br><br> Default: **None for Standard storage. ReadOnly for Premium storage**. Possible values include: 'CachingTypesNone', 'CachingTypesReadOnly', 'CachingTypesReadWrite'
	Caching CachingTypes `json:"caching,omitempty"`
	// DiskSizeGB - Specifies the size of empty data disks in gigabytes. This element can be used to overwrite the name of the disk in a virtual machine image. <br><br> This value cannot be larger than 1023 GB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
	// StorageAccountType - Specifies the storage account type for the managed disk. NOTE: UltraSSD_LRS can only be used with data disks, it cannot be used with OS Disk. Possible values include: 'StorageAccountTypesStandardLRS', 'StorageAccountTypesPremiumLRS', 'StorageAccountTypesStandardSSDLRS', 'StorageAccountTypesUltraSSDLRS'
	StorageAccountType StorageAccountTypes `json:"storageAccountType,omitempty"`
}

// VirtualMachineImageStorageProfile describes a storage profile.
type VirtualMachineImageStorageProfile struct {
	// OsDisk - Specifies information about the operating system disk used by the virtual machine. <br><br> For more information about disks, see [About disks and VHDs for Azure virtual machines](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-windows-about-disks-vhds?toc=%2fazure%2fvirtual-machines%2fwindows%2ftoc.json).
	OsDisk *VirtualMachineImageOSDisk `json:"osDisk,omitempty"`
	// DataDisks - Specifies the parameters that are used to add a data disk to a virtual machine. <br><br> For more information about disks, see [About disks and VHDs for Azure virtual machines](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-windows-about-disks-vhds?toc=%2fazure%2fvirtual-machines%2fwindows%2ftoc.json).
	DataDisks *[]VirtualMachineImageDataDisk `json:"dataDisks,omitempty"`
	// ZoneResilient - Specifies whether an image is zone resilient or not. Default is false. Zone resilient images can be created only in regions that provide Zone Redundant Storage (ZRS).
	ZoneResilient *bool `json:"zoneResilient,omitempty"`
}

// VirtualMachineImageProperties describes the properties of an Image.
type VirtualMachineImageProperties struct {
	// SourceVirtualMachine - The source virtual machine from which Image is created.
	SourceVirtualMachine *SubResource `json:"sourceVirtualMachine,omitempty"`
	// StorageProfile - Specifies the storage settings for the virtual machine disks.
	StorageProfile *VirtualMachineImageStorageProfile `json:"storageProfile,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// HyperVGeneration - Gets the HyperVGenerationType of the VirtualMachine created from the image. Possible values include: 'HyperVGenerationTypesV1', 'HyperVGenerationTypesV2'
	HyperVGeneration common.HyperVGeneration `json:"hyperVGeneration,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Image the source user image virtual hard disk. The virtual hard disk will be copied before being
// attached to the virtual machine. If SourceImage is provided, the destination virtual hard drive must not
// exist.
type VirtualMachineImage struct {
	autorest.Response              `json:"-"`
	*VirtualMachineImageProperties `json:"properties,omitempty"`
	// ID - READ-ONLY; Resource Id
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags
	Tags map[string]*string `json:"tags"`
}

// BareMetalHostDisk describes a bare metal host disk.
type BareMetalHostDisk struct {
	// Name - Name of the disk
	Name *string `json:"name,omitempty"`
	// DiskSizeGB - Specifies the size of the disk in gigabytes
	DiskSizeGB *uint32 `json:"diskSizeGB,omitempty"`
}

type BareMetalHostStorageProfile struct {
	// Disks
	Disks *[]BareMetalHostDisk `json:"disks,omitempty"`
}

type BareMetalHostNetworkInterface struct {
	// Name
	Name *string `json:"name,omitempty"`
}

type BareMetalHostNetworkProfile struct {
	// NetworkInterfaces
	NetworkInterfaces *[]BareMetalHostNetworkInterface `json:"networkinterfaces,omitempty"`
}

// BareMetalHostSize Specifies cpu/memory information for bare metal hosts.
type BareMetalHostSize struct {
	CpuCount *uint32 `json:"cpucount,omitempty"`
	GpuCount *uint32 `json:"gpucount,omitempty"`
	MemoryMB *uint32 `json:"memorymb,omitempty"`
}

type BareMetalHostHardwareProfile struct {
	MachineSize *BareMetalHostSize `json:"machinesize,omitempty"`
}

// BareMetalHostProperties describes the properties of a Bare Metal Host.
type BareMetalHostProperties struct {
	// StorageProfile
	StorageProfile *BareMetalHostStorageProfile `json:"storageprofile,omitempty"`
	// NetworkProfile
	NetworkProfile *BareMetalHostNetworkProfile `json:"networkprofile,omitempty"`
	// HardwareProfile - Specifies the hardware settings for the bare metal host.
	HardwareProfile *BareMetalHostHardwareProfile `json:"hardwareprofile,omitempty"`
	// SecurityProfile - Specifies the security settings for the bare metal host.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// BareMetalMachine - Specifies information about the bare metal machine.
	// FQDN
	FQDN *string `json:"fqdn,omitempty"`
	// Port
	Port *uint32 `json:"port,omitempty"`
	// AuthorizerPort
	AuthorizerPort *uint32 `json:"authorizerPort,omitempty"`
	// Certificate
	Certificate *string `json:"certificate,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

type BareMetalHost struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Properties
	*BareMetalHostProperties `json:"baremetalhostproperties,omitempty"`
}

type BareMetalMachineImageReference struct {
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
	// Name - Name of the image
	Name *string `json:"name,omitempty"`
}

type BareMetalMachineStorageProfile struct {
	// ImageReference - Specifies information about the image to use.
	ImageReference *BareMetalMachineImageReference `json:"imagereference,omitempty"`
}

type BareMetalMachineOSProfile struct {
	// ComputerName
	ComputerName *string `json:"computername,omitempty"`
	// AdminUsername
	AdminUsername *string `json:"adminusername,omitempty"`
	// AdminPassword
	AdminPassword *string `json:"adminpassword,omitempty"`
	// CustomData Specifies a base-64 encoded string of custom data.
	CustomData *string `json:"customdata,omitempty"`
	// LinuxConfiguration
	LinuxConfiguration *LinuxConfiguration `json:"linuxconfiguration,omitempty"`
}

// BareMetalMachineProperties describes the properties of a Bare Metal Machine.
type BareMetalMachineProperties struct {
	// StorageProfile
	StorageProfile *BareMetalMachineStorageProfile `json:"storageprofile,omitempty"`
	// OsProfile
	OsProfile *BareMetalMachineOSProfile `json:"osprofile,omitempty"`
	// SecurityProfile - Specifies the security settings for the bare metal machine.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// FQDN - Copied from associated BareMetalHost
	FQDN *string `json:"fqdn,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

type BareMetalMachine struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location
	Location *string `json:"location,omitempty"`
	// Properties
	*BareMetalMachineProperties `json:"baremetalmachineproperties,omitempty"`
}

type ExecutionState string

const (
	// ExecutionStateFailed ...
	ExecutionStateFailed ExecutionState = "Failed"
	// ExecutionStateSucceeded ...
	ExecutionStateSucceeded ExecutionState = "Succeeded"
	// ExecutionStateUnknown ...
	ExecutionStateUnknown ExecutionState = "Unknown"
)

// VirtualMachineRunCommandScriptSource describes the script sources for run command.
type VirtualMachineRunCommandScriptSource struct {
	// Script - Specifies the script content to be executed on the VM.
	Script *string `json:"script,omitempty"`
	// ScriptURI - Specifies the script download location.
	ScriptURI *string `json:"scriptUri,omitempty"`
	// CommandID - Specifies a commandId of predefined built-in script.
	CommandID *string `json:"commandId,omitempty"`
}

// RunCommandInputParameter describes the properties of a run command parameter.
type RunCommandInputParameter struct {
	// Name - The run command parameter name.
	Name *string `json:"name,omitempty"`
	// Value - The run command parameter value.
	Value *string `json:"value,omitempty"`
}

// VirtualMachineRunCommandInstanceView the instance view of a virtual machine run command.
type VirtualMachineRunCommandInstanceView struct {
	// ExecutionState - Script execution status. Possible values include: 'ExecutionStateUnknown', 'ExecutionStateFailed', 'ExecutionStateSucceeded'
	ExecutionState ExecutionState `json:"executionState,omitempty"`
	// ExitCode - Exit code returned from script execution.
	ExitCode *int32 `json:"exitCode,omitempty"`
	// Output - Script output stream.
	Output *string `json:"output,omitempty"`
	// Error - Script error stream.
	Error *string `json:"error,omitempty"`
}

// VirtualMachineRunCommandRequest describes the properties of a Virtual Machine run command.
type VirtualMachineRunCommandRequest struct {
	// Source - The source of the run command script.
	Source *VirtualMachineRunCommandScriptSource `json:"source,omitempty"`
	// Parameters - The parameters used by the script.
	Parameters    *[]RunCommandInputParameter `json:"parameters,omitempty"`
	RunAsUser     *string                     `json:"runasuser,omitempty"`
	RunAsPassword *string                     `json:"runaspassword,omitempty"`
}

// VirtualMachineRunCommandResponse
type VirtualMachineRunCommandResponse struct {
	// InstanceView - The virtual machine run command instance view.
	InstanceView *VirtualMachineRunCommandInstanceView `json:"instanceView,omitempty"`
}
