// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package network

import (
	"github.com/Azure/go-autorest/autorest"
)

type TransportProtocol string

const (
	// TransportProtocolAll
	TransportProtocolAll TransportProtocol = "All"
	// TransportProtocolTCP
	TransportProtocolTCP TransportProtocol = "Tcp"
	// TransportProtocolUDP
	TransportProtocolUDP TransportProtocol = "Udp"
)

// SubResource reference to another subresource.
type SubResource struct {
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// APIEntityReference the API entity reference.
type APIEntityReference struct {
	// ID - The ARM resource id in the form of /subscriptions/{SubscriptionId}/resourceGroups/{ResourceGroupName}/...
	ID *string `json:"id,omitempty"`
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// Deleting ...
	Deleting ProvisioningState = "Deleting"
	// Failed ...
	Failed ProvisioningState = "Failed"
	// Succeeded ...
	Succeeded ProvisioningState = "Succeeded"
	// Updating ...
	Updating ProvisioningState = "Updating"
)

// RouteTablePropertiesFormat route Table resource.
type RouteTablePropertiesFormat struct {
	// Routes - Collection of routes contained within a route table.
	Routes *[]Route `json:"routes,omitempty"`
	// Subnets - READ-ONLY; A collection of references to subnets.
	Subnets *[]Subnet `json:"subnets,omitempty"`
	// ProvisioningState - The provisioning state of the resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// RouteTable route table resource.
type RouteTable struct {
	autorest.Response `json:"-"`
	// RouteTablePropertiesFormat - Properties of the route table.
	*RouteTablePropertiesFormat `json:"properties,omitempty"`
	// Etag - Gets a unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// RoutePropertiesFormat route resource.
type RoutePropertiesFormat struct {
	// AddressPrefix - The destination CIDR to which the route applies.
	AddressPrefix *string `json:"addressPrefix,omitempty"`
	// NextHopIPAddress - The IP address packets should be forwarded to. Next hop values are only allowed in routes where the next hop type is VirtualAppliance.
	NextHopIPAddress *string `json:"nextHopIpAddress,omitempty"`
	// ProvisioningState - The provisioning state of the resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Route is assoicated with a subnet.
type Route struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// RouteProperties - Properties of the route.
	*RoutePropertiesFormat `json:"properties,omitempty"`
}

// IPConfigurationReference
type IPConfigurationReference struct {
	// IPConfigurationID
	IPConfigurationID *string `json:"ID,omitempty"`
}

// Subnet subnet in a virtual network resource.
type Subnet struct {
	autorest.Response `json:"-"`
	// SubnetPropertiesFormat - Properties of the subnet.
	*SubnetPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// Subnet is assoicated with a Virtual Network.
type SubnetPropertiesFormat struct {
	// Cidr for this subnet - IPv4, IPv6
	AddressPrefix *string `json:"addressPrefix,omitempty"`
	// AddressPrefixes - List of address prefixes for the subnet.
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
	// Routes for the subnet
	RouteTable *RouteTable `json:"routeTable,omitempty"`
	// IPConfigurationReferences
	IPConfigurationReferences *[]IPConfigurationReference `json:"ipConfigurationReferences,omitempty"`
	// IPAllocationMethod - The IP address allocation method. Possible values include: 'Static', 'Dynamic'
	IPAllocationMethod IPAllocationMethod `json:"ipAllocationMethod,omitempty"`
	// Vlan
	Vlan    *uint16  `json:"vlan,omitempty"`
	IPPools []IPPool `json:"ippools,omitempty"`
	// NetworkSecurityGroup - The resource reference of the subnet's applied network security group
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`
}

type IPPoolType string

const (
	VM      IPPoolType = "vm"
	VIPPOOL IPPoolType = "vippool"
)

type IPPoolInfo struct {
	// used - no. of ip addresses already allocated from the pool
	Used string `json:"used,omitempty"`
	// available - no. of ip addresses still available in the pool
	Available string `json:"available,omitempty"`
}

// IPPool is associated with a network and represents pool of IP addresses.
type IPPool struct {
	// Name
	Name string `json:"name,omitempty"`
	// Type
	Type IPPoolType `json:"ippooltype,omitempty"`
	// Start - The starting ip address of the pool
	Start string `json:"start,omitempty"`
	// end - The ending ip address of the pool
	End string `json:"end,omitempty"`
	// Auxilliary info associated with an ip pool
	Info *IPPoolInfo `json:"info,omitempty"`
}

// MACRange is associated with MACPool and respresents the start and end addresses.
type MACRange struct {
	// StartMACAddress
	StartMACAddress *string `json:"startmacaddress,omitempty"`
	// EndMACAddress
	EndMACAddress *string `json:"endmacaddress,omitempty"`
}

// MACPoolProperties MAC pool properties.
type MACPoolPropertiesFormat struct {
	// MAC ranges
	Range *MACRange `json:"range,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// MACPool is assoicated with a network and represents pool of MACRanges.
type MACPool struct {
	autorest.Response `json:"-"`
	// MacPoolPropertiesFormat - MAC Pool properties.
	*MACPoolPropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// InterfaceDNSSettings DNS settings of a network interface.
type InterfaceDNSSettings struct {
	// DNSServers - List of DNS servers IP addresses. Use 'AzureProvidedDNS' to switch to azure provided DNS resolution. 'AzureProvidedDNS' value cannot be combined with other IPs, it must be the only value in dnsServers collection.
	DNSServers *[]string `json:"dnsServers,omitempty"`
	// AppliedDNSServers - If the VM that uses this NIC is part of an Availability Set, then this list will have the union of all DNS servers from all NICs that are part of the Availability Set. This property is what is configured on each of those VMs.
	AppliedDNSServers *[]string `json:"appliedDnsServers,omitempty"`
	// InternalDNSNameLabel - Relative DNS name for this NIC used for internal communications between VMs in the same virtual network.
	InternalDNSNameLabel *string `json:"internalDnsNameLabel,omitempty"`
	// InternalFqdn - Fully qualified DNS name supporting internal communications between VMs in the same virtual network.
	InternalFqdn *string `json:"internalFqdn,omitempty"`
	// InternalDomainNameSuffix - Even if internalDnsNameLabel is not specified, a DNS entry is created for the primary NIC of the VM. This DNS name can be constructed by concatenating the VM name with the value of internalDomainNameSuffix.
	InternalDomainNameSuffix *string `json:"internalDomainNameSuffix,omitempty"`
}

// AddressSpace addressSpace contains an array of IP address ranges that can be used by subnets of the
// virtual network.
type AddressSpace struct {
	// AddressPrefixes - A list of address blocks reserved for this virtual network in CIDR notation.
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
}

// FrontendIPConfigurationPropertiesFormat properties of Frontend IP Configuration of the load balancer.
type FrontendIPConfigurationPropertiesFormat struct {
	// InboundNatRules - READ-ONLY; Read only. Inbound rules URIs that use this frontend IP.
	InboundNatRules *[]SubResource `json:"inboundNatRules,omitempty"`
	// InboundNatPools - READ-ONLY; Read only. Inbound pools URIs that use this frontend IP.
	InboundNatPools *[]SubResource `json:"inboundNatPools,omitempty"`
	// OutboundRules - READ-ONLY; Read only. Outbound rules URIs that use this frontend IP.
	OutboundRules *[]SubResource `json:"outboundRules,omitempty"`
	// LoadBalancingRules - READ-ONLY; Gets load balancing rules URIs that use this frontend IP.
	LoadBalancingRules *[]SubResource `json:"loadBalancingRules,omitempty"`
	// PrivateIPAddress - The private IP address of the IP configuration.
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
	// PrivateIPAllocationMethod - The Private IP allocation method. Possible values include: 'Static', 'Dynamic'
	PrivateIPAllocationMethod IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	// PrivateIPAddressVersion - It represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4. Possible values include: 'IPv4', 'IPv6'
	PrivateIPAddressVersion IPVersion `json:"privateIPAddressVersion,omitempty"`
	// Subnet - The reference of the subnet resource.
	Subnet *Subnet `json:"subnet,omitempty"`
	// PublicIPAddress - The reference of the Public IP resource.
	PublicIPAddress *PublicIPAddress `json:"publicIPAddress,omitempty"`
	// PublicIPPrefix - The reference of the Public IP Prefix resource.
	PublicIPPrefix *SubResource `json:"publicIPPrefix,omitempty"`
	// ProvisioningState - Gets the provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// IPAddress - The ip address of the frontend
	IPAddress *string `json:"ipAddress,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// FrontendIPConfiguration frontend IP address of the load balancer.
type FrontendIPConfiguration struct {
	autorest.Response `json:"-"`
	// FrontendIPConfigurationPropertiesFormat - Properties of the load balancer probe.
	*FrontendIPConfigurationPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within the set of frontend IP configurations used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// Zones - A list of availability zones denoting the IP allocated for the resource needs to come from.
	Zones *[]string `json:"zones,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// BackendAddressPoolProperties properties of the backend address pool.
type BackendAddressPoolPropertiesFormat struct {
	// BackendIPConfigurations - READ-ONLY; Gets collection of references to IP addresses defined in network interfaces.
	BackendIPConfigurations *[]InterfaceIPConfiguration `json:"backendIPConfigurations,omitempty"`
	// LoadBalancingRules - READ-ONLY; Gets load balancing rules that use this backend address pool.
	LoadBalancingRules *[]SubResource `json:"loadBalancingRules,omitempty"`
	// OutboundRule - READ-ONLY; Gets outbound rules that use this backend address pool.
	OutboundRule *SubResource `json:"outboundRule,omitempty"`
	// OutboundRules - READ-ONLY; Gets outbound rules that use this backend address pool.
	OutboundRules *[]SubResource `json:"outboundRules,omitempty"`
	// ProvisioningState - Get provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

type BackendAddressPool struct {
	// BackendAddressPoolPropertiesFormat - Properties of load balancer backend address pool.
	*BackendAddressPoolPropertiesFormat `json:"properties,omitempty"`
	// Name - Gets name of the resource that is unique within the set of backend address pools used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// LoadDistribution enumerates the values for load distribution.
type LoadDistribution string

const (
	// LoadDistributionDefault ...
	LoadDistributionDefault LoadDistribution = "Default"
	// LoadDistributionSourceIP ...
	LoadDistributionSourceIP LoadDistribution = "SourceIP"
	// LoadDistributionSourceIPProtocol ...
	LoadDistributionSourceIPProtocol LoadDistribution = "SourceIPProtocol"
)

// LoadBalancingRulePropertiesFormat properties of the load balancer.
type LoadBalancingRulePropertiesFormat struct {
	// FrontendIPConfiguration - A reference to frontend IP addresses.
	FrontendIPConfiguration *SubResource `json:"frontendIPConfiguration,omitempty"`
	// BackendAddressPool - A reference to a pool of DIPs. Inbound traffic is randomly load balanced across IPs in the backend IPs.
	BackendAddressPool *SubResource `json:"backendAddressPool,omitempty"`
	// Probe - The reference of the load balancer probe used by the load balancing rule.
	Probe *SubResource `json:"probe,omitempty"`
	// Protocol - The reference to the transport protocol used by the load balancing rule. Possible values include: 'TransportProtocolUDP', 'TransportProtocolTCP', 'TransportProtocolAll'
	Protocol TransportProtocol `json:"protocol,omitempty"`
	// LoadDistribution - The load distribution policy for this rule. Possible values include: 'LoadDistributionDefault', 'LoadDistributionSourceIP', 'LoadDistributionSourceIPProtocol'
	LoadDistribution LoadDistribution `json:"loadDistribution,omitempty"`
	// FrontendPort - The port for the external endpoint. Port numbers for each rule must be unique within the Load Balancer. Acceptable values are between 0 and 65534. Note that value 0 enables "Any Port".
	FrontendPort *int32 `json:"frontendPort,omitempty"`
	// BackendPort - The port used for internal connections on the endpoint. Acceptable values are between 0 and 65535. Note that value 0 enables "Any Port".
	BackendPort *int32 `json:"backendPort,omitempty"`
	// IdleTimeoutInMinutes - The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// EnableFloatingIP - Configures a virtual machine's endpoint for the floating IP capability required to configure a SQL AlwaysOn Availability Group. This setting is required when using the SQL AlwaysOn Availability Groups in SQL server. This setting can't be changed after you create the endpoint.
	EnableFloatingIP *bool `json:"enableFloatingIP,omitempty"`
	// EnableTCPReset - Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.
	EnableTCPReset *bool `json:"enableTcpReset,omitempty"`
	// DisableOutboundSnat - Configures SNAT for the VMs in the backend pool to use the publicIP address specified in the frontend of the load balancing rule.
	DisableOutboundSnat *bool `json:"disableOutboundSnat,omitempty"`
	// ProvisioningState - Gets the provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// LoadBalancingRule a load balancing rule for a load balancer.
type LoadBalancingRule struct {
	autorest.Response `json:"-"`
	// LoadBalancingRulePropertiesFormat - Properties of load balancer load balancing rule.
	*LoadBalancingRulePropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within the set of load balancing rules used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// ProbeProtocol enumerates the values for probe protocol.
type ProbeProtocol string

const (
	// ProbeProtocolHTTP ...
	ProbeProtocolHTTP ProbeProtocol = "Http"
	// ProbeProtocolHTTPS ...
	ProbeProtocolHTTPS ProbeProtocol = "Https"
	// ProbeProtocolTCP ...
	ProbeProtocolTCP ProbeProtocol = "Tcp"
)

// ProbePropertiesFormat load balancer probe resource.
type ProbePropertiesFormat struct {
	// LoadBalancingRules - READ-ONLY; The load balancer rules that use this probe.
	LoadBalancingRules *[]SubResource `json:"loadBalancingRules,omitempty"`
	// Protocol - The protocol of the end point. If 'Tcp' is specified, a received ACK is required for the probe to be successful. If 'Http' or 'Https' is specified, a 200 OK response from the specifies URI is required for the probe to be successful. Possible values include: 'ProbeProtocolHTTP', 'ProbeProtocolTCP', 'ProbeProtocolHTTPS'
	Protocol ProbeProtocol `json:"protocol,omitempty"`
	// Port - The port for communicating the probe. Possible values range from 1 to 65535, inclusive.
	Port *int32 `json:"port,omitempty"`
	// IntervalInSeconds - The interval, in seconds, for how frequently to probe the endpoint for health status. Typically, the interval is slightly less than half the allocated timeout period (in seconds) which allows two full probes before taking the instance out of rotation. The default value is 15, the minimum value is 5.
	IntervalInSeconds *int32 `json:"intervalInSeconds,omitempty"`
	// NumberOfProbes - The number of probes where if no response, will result in stopping further traffic from being delivered to the endpoint. This values allows endpoints to be taken out of rotation faster or slower than the typical times used in Azure.
	NumberOfProbes *int32 `json:"numberOfProbes,omitempty"`
	// RequestPath - The URI used for requesting health status from the VM. Path is required if a protocol is set to http. Otherwise, it is not allowed. There is no default value.
	RequestPath *string `json:"requestPath,omitempty"`
	// ProvisioningState - Gets the provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Probe a load balancer probe.
type Probe struct {
	autorest.Response `json:"-"`
	// ProbePropertiesFormat - Properties of load balancer probe.
	*ProbePropertiesFormat `json:"properties,omitempty"`
	// Name - Gets name of the resource that is unique within the set of probes used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// InboundNatPoolPropertiesFormat properties of Inbound NAT pool.
type InboundNatPoolPropertiesFormat struct {
	// FrontendIPConfiguration - A reference to frontend IP addresses.
	FrontendIPConfiguration *SubResource `json:"frontendIPConfiguration,omitempty"`
	// Protocol - The reference to the transport protocol used by the inbound NAT pool. Possible values include: 'TransportProtocolUDP', 'TransportProtocolTCP', 'TransportProtocolAll'
	Protocol TransportProtocol `json:"protocol,omitempty"`
	// FrontendPortRangeStart - The first port number in the range of external ports that will be used to provide Inbound Nat to NICs associated with a load balancer. Acceptable values range between 1 and 65534.
	FrontendPortRangeStart *int32 `json:"frontendPortRangeStart,omitempty"`
	// FrontendPortRangeEnd - The last port number in the range of external ports that will be used to provide Inbound Nat to NICs associated with a load balancer. Acceptable values range between 1 and 65535.
	FrontendPortRangeEnd *int32 `json:"frontendPortRangeEnd,omitempty"`
	// BackendPort - The port used for internal connections on the endpoint. Acceptable values are between 1 and 65535.
	BackendPort *int32 `json:"backendPort,omitempty"`
	// IdleTimeoutInMinutes - The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// EnableFloatingIP - Configures a virtual machine's endpoint for the floating IP capability required to configure a SQL AlwaysOn Availability Group. This setting is required when using the SQL AlwaysOn Availability Groups in SQL server. This setting can't be changed after you create the endpoint.
	EnableFloatingIP *bool `json:"enableFloatingIP,omitempty"`
	// EnableTCPReset - Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.
	EnableTCPReset *bool `json:"enableTcpReset,omitempty"`
	// ProvisioningState - Gets the provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// InboundNatPool inbound NAT pool of the load balancer.
type InboundNatPool struct {
	// InboundNatPoolPropertiesFormat - Properties of load balancer inbound nat pool.
	*InboundNatPoolPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within the set of inbound NAT pools used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// LoadBalancerOutboundRuleProtocol enumerates the values for load balancer outbound rule protocol.
type LoadBalancerOutboundRuleProtocol string

const (
	// LoadBalancerOutboundRuleProtocolAll ...
	LoadBalancerOutboundRuleProtocolAll LoadBalancerOutboundRuleProtocol = "All"
	// LoadBalancerOutboundRuleProtocolTCP ...
	LoadBalancerOutboundRuleProtocolTCP LoadBalancerOutboundRuleProtocol = "Tcp"
	// LoadBalancerOutboundRuleProtocolUDP ...
	LoadBalancerOutboundRuleProtocolUDP LoadBalancerOutboundRuleProtocol = "Udp"
)

// OutboundRulePropertiesFormat outbound rule of the load balancer.
type OutboundRulePropertiesFormat struct {
	// AllocatedOutboundPorts - The number of outbound ports to be used for NAT.
	AllocatedOutboundPorts *int32 `json:"allocatedOutboundPorts,omitempty"`
	// FrontendIPConfigurations - The Frontend IP addresses of the load balancer.
	FrontendIPConfigurations *[]SubResource `json:"frontendIPConfigurations,omitempty"`
	// BackendAddressPool - A reference to a pool of DIPs. Outbound traffic is randomly load balanced across IPs in the backend IPs.
	BackendAddressPool *SubResource `json:"backendAddressPool,omitempty"`
	// ProvisioningState - Gets the provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Protocol - The protocol for the outbound rule in load balancer. Possible values include: 'LoadBalancerOutboundRuleProtocolTCP', 'LoadBalancerOutboundRuleProtocolUDP', 'LoadBalancerOutboundRuleProtocolAll'
	Protocol LoadBalancerOutboundRuleProtocol `json:"protocol,omitempty"`
	// EnableTCPReset - Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.
	EnableTCPReset *bool `json:"enableTcpReset,omitempty"`
	// IdleTimeoutInMinutes - The timeout for the TCP idle connection.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// OutboundRule outbound rule of the load balancer.
type OutboundRule struct {
	autorest.Response `json:"-"`
	// OutboundRulePropertiesFormat - Properties of load balancer outbound rule.
	*OutboundRulePropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within the set of outbound rules used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// LoadBalancerPropertiesFormat properties of the load balancer.
type LoadBalancerPropertiesFormat struct {
	// FrontendIPConfigurations - Object representing the frontend IPs to be used for the load balancer.
	FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
	// BackendAddressPools - Collection of backend address pools used by a load balancer.
	BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
	// LoadBalancingRules - Object collection representing the load balancing rules Gets the provisioning.
	LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
	// Probes - Collection of probe objects used in the load balancer.
	Probes *[]Probe `json:"probes,omitempty"`
	// InboundNatRules - Collection of inbound NAT Rules used by a load balancer. Defining inbound NAT rules on your load balancer is mutually exclusive with defining an inbound NAT pool. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an Inbound NAT pool. They have to reference individual inbound NAT rules.
	InboundNatRules *[]InboundNatRule `json:"inboundNatRules,omitempty"`
	// InboundNatPools - Defines an external port range for inbound NAT to a single backend port on NICs associated with a load balancer. Inbound NAT rules are created automatically for each NIC associated with the Load Balancer using an external port from this range. Defining an Inbound NAT pool on your Load Balancer is mutually exclusive with defining inbound Nat rules. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an inbound NAT pool. They have to reference individual inbound NAT rules.
	InboundNatPools *[]InboundNatPool `json:"inboundNatPools,omitempty"`
	// OutboundRules - The outbound rules.
	OutboundRules *[]OutboundRule `json:"outboundRules,omitempty"`
	// ResourceGUID - The resource GUID property of the load balancer resource.
	ResourceGUID *string `json:"resourceGuid,omitempty"`
	// ProvisioningState - Gets the provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// ReplicateionCount
	ReplicationCount uint32 `json:"replicationCount,omitempty"`
}

// LoadBalancer loadBalancer resource.
type LoadBalancer struct {
	// LoadBalancerPropertiesFormat - Properties of load balancer.
	*LoadBalancerPropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// DhcpOptions dhcpOptions contains an array of DNS servers available to VMs deployed in the virtual
// network. Standard DHCP option for a subnet overrides VNET DHCP options.
type DhcpOptions struct {
	// DNSServers - The list of DNS servers IP addresses.
	DNSServers *[]string `json:"dnsServers,omitempty"`
}

// PortForwardingRule defines the structure of a port forwading rule
type PortForwardingRule struct {
	// Type
	Type *string `json:"type,omitempty"`
	// ConnectAddress
	ConnectAddress *string `json:"connectaddress,omitempty"`
	// ConnectPort
	ConnectPort *string `json:"connectport,omitempty"`
	// ListenAddress
	ListenAddress *string `json:"listenaddress,omitempty"`
	// ListenPort
	ListenPort *string `json:"listenport,omitempty"`
}

// VirtualNetworkPropertiesFormat properties of the virtual network.
type VirtualNetworkPropertiesFormat struct {
	// AddressSpace - The AddressSpace that contains an array of IP address ranges that can be used by subnets.
	AddressSpace *AddressSpace `json:"addressSpace,omitempty"`
	// DhcpOptions - The dhcpOptions that contains an array of DNS servers available to VMs deployed in the virtual network.
	DhcpOptions *DhcpOptions `json:"dhcpOptions,omitempty"`
	// Subnets - A list of subnets in a Virtual Network.
	Subnets *[]Subnet `json:"subnets,omitempty"`
	// MACPool name - Name of the associated MAC pool (or leave empty to use the default mac pool)
	MacPoolName *string `json:"macPoolName,omitempty"`
	// ProvisioningState - The provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Port forwarding rules
	PortForwardingRules *[]PortForwardingRule `json:"portforwardingrules,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetwork struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// VirtualNetworkProperties - Properties of the virtual network.
	*VirtualNetworkPropertiesFormat `json:"properties,omitempty"`
}

// IPAllocationMethod enumerates the values for ip allocation method.
type IPAllocationMethod string

const (
	// Dynamic ...
	Dynamic IPAllocationMethod = "Dynamic"
	// Static ...
	Static IPAllocationMethod = "Static"
)

// IPConfigurationPropertiesFormat properties of IP configuration.
type IPConfigurationPropertiesFormat struct {
	// PrivateIPAddress - The private IP address of the IP configuration.
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
	// PrivateIPAllocationMethod - The private IP address allocation method. Possible values include: 'Static', 'Dynamic'
	PrivateIPAllocationMethod IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	// Subnet - The reference of the subnet resource.
	Subnet *Subnet `json:"subnet,omitempty"`
	// PublicIPAddress - The reference of the public IP resource.
	PublicIPAddress *PublicIPAddress `json:"publicIPAddress,omitempty"`
	// ProvisioningState - Gets the provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// IPConfiguration
type IPConfiguration struct {
	// IPConfigurationPropertiesFormat - Properties of the IP configuration.
	*IPConfigurationPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// InboundNatRulePropertiesFormat properties of the inbound NAT rule.
type InboundNatRulePropertiesFormat struct {
	// FrontendIPConfiguration - A reference to frontend IP addresses.
	FrontendIPConfiguration *SubResource `json:"frontendIPConfiguration,omitempty"`
	// BackendIPConfiguration - READ-ONLY; A reference to a private IP address defined on a network interface of a VM. Traffic sent to the frontend port of each of the frontend IP configurations is forwarded to the backend IP.
	BackendIPConfiguration *InterfaceIPConfiguration `json:"backendIPConfiguration,omitempty"`
	// Protocol - The reference to the transport protocol used by the load balancing rule. Possible values include: 'TransportProtocolUDP', 'TransportProtocolTCP', 'TransportProtocolAll'
	Protocol TransportProtocol `json:"protocol,omitempty"`
	// FrontendPort - The port for the external endpoint. Port numbers for each rule must be unique within the Load Balancer. Acceptable values range from 1 to 65534.
	FrontendPort *int32 `json:"frontendPort,omitempty"`
	// BackendPort - The port used for the internal endpoint. Acceptable values range from 1 to 65535.
	BackendPort *int32 `json:"backendPort,omitempty"`
	// IdleTimeoutInMinutes - The timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// EnableFloatingIP - Configures a virtual machine's endpoint for the floating IP capability required to configure a SQL AlwaysOn Availability Group. This setting is required when using the SQL AlwaysOn Availability Groups in SQL server. This setting can't be changed after you create the endpoint.
	EnableFloatingIP *bool `json:"enableFloatingIP,omitempty"`
	// EnableTCPReset - Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.
	EnableTCPReset *bool `json:"enableTcpReset,omitempty"`
	// ProvisioningState - Gets the provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// InboundNatRule inbound NAT rule of the load balancer.
type InboundNatRule struct {
	// InboundNatRulePropertiesFormat - Properties of load balancer inbound nat rule.
	*InboundNatRulePropertiesFormat `json:"properties,omitempty"`
	// Name - Gets name of the resource that is unique within the set of inbound NAT rules used by the load balancer. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Type - READ-ONLY; Type of the resource.
	Type *string `json:"type,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// IPVersion enumerates the values for ip version.
type IPVersion string

const (
	// IPv4 ...
	IPv4 IPVersion = "IPv4"
	// IPv6 ...
	IPv6 IPVersion = "IPv6"
)

// InterfaceIPConfigurationPropertiesFormat properties of IP configuration.
type InterfaceIPConfigurationPropertiesFormat struct {
	// LoadBalancerBackendAddressPools - The reference of LoadBalancerBackendAddressPool resource.
	LoadBalancerBackendAddressPools *[]BackendAddressPool `json:"loadBalancerBackendAddressPools,omitempty"`
	// LoadBalancerInboundNatRules - A list of references of LoadBalancerInboundNatRules.
	LoadBalancerInboundNatRules *[]InboundNatRule `json:"loadBalancerInboundNatRules,omitempty"`
	// PrivateIPAddress - Private IP address of the IP configuration.
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
	// PrivateIPAllocationMethod - The private IP address allocation method. Possible values include: 'Static', 'Dynamic'
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	// PrivateIPAddressVersion - Available from Api-Version 2016-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4. Possible values include: 'IPv4', 'IPv6'
	PrivateIPAddressVersion IPVersion `json:"privateIPAddressVersion,omitempty"`
	// Subnet - Subnet bound to the IP configuration.
	Subnet *APIEntityReference `json:"subnet,omitempty"`
	// PrefixLength
	PrefixLength *string `json:"prefixlength,omitempty"`
	// Gateway
	Gateway *string `json:"gateway,omitempty"`
	// Primary - Gets whether this is a primary customer address on the network interface.
	Primary *bool `json:"primary,omitempty"`
	// PublicIPAddress - Public IP address bound to the IP configuration.
	PublicIPAddress *PublicIPAddress `json:"publicIPAddress,omitempty"`
	// ProvisioningState - The provisioning state of the network interface IP configuration. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// NetworkSecurityGroup - The reference of the NetworkSecurityGroup resource.
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`
}

// InterfaceIPConfiguration iPConfiguration in a network interface.
type InterfaceIPConfiguration struct {
	autorest.Response `json:"-"`
	// InterfaceIPConfigurationPropertiesFormat - Network interface IP configuration properties.
	*InterfaceIPConfigurationPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// SecurityRuleProtocol enumerates the values for security rule protocol.
type SecurityRuleProtocol string

const (
	// SecurityRuleProtocolAsterisk ...
	SecurityRuleProtocolAsterisk SecurityRuleProtocol = "*"
	// SecurityRuleProtocolEsp ...
	SecurityRuleProtocolEsp SecurityRuleProtocol = "Esp"
	// SecurityRuleProtocolIcmp ...
	SecurityRuleProtocolIcmp SecurityRuleProtocol = "Icmp"
	// SecurityRuleProtocolTCP ...
	SecurityRuleProtocolTCP SecurityRuleProtocol = "Tcp"
	// SecurityRuleProtocolUDP ...
	SecurityRuleProtocolUDP SecurityRuleProtocol = "Udp"
)

// ApplicationSecurityGroupPropertiesFormat application security group properties.
type ApplicationSecurityGroupPropertiesFormat struct {
	// ResourceGUID - READ-ONLY; The resource GUID property of the application security group resource. It uniquely identifies a resource, even if the user changes its name or migrate the resource across subscriptions or resource groups.
	ResourceGUID *string `json:"resourceGuid,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state of the application security group resource. Possible values are: 'Succeeded', 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// ApplicationSecurityGroup an application security group in a resource group.
type ApplicationSecurityGroup struct {
	autorest.Response `json:"-"`
	// ApplicationSecurityGroupPropertiesFormat - Properties of the application security group.
	*ApplicationSecurityGroupPropertiesFormat `json:"properties,omitempty"`
	// Etag - READ-ONLY; A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// SecurityRuleAccess enumerates the values for security rule access.
type SecurityRuleAccess string

const (
	// SecurityRuleAccessAllow ...
	SecurityRuleAccessAllow SecurityRuleAccess = "Allow"
	// SecurityRuleAccessDeny ...
	SecurityRuleAccessDeny SecurityRuleAccess = "Deny"
)

// SecurityRuleDirection enumerates the values for security rule direction.
type SecurityRuleDirection string

const (
	// SecurityRuleDirectionInbound ...
	SecurityRuleDirectionInbound SecurityRuleDirection = "Inbound"
	// SecurityRuleDirectionOutbound ...
	SecurityRuleDirectionOutbound SecurityRuleDirection = "Outbound"
)

// SecurityRulePropertiesFormat security rule resource.
type SecurityRulePropertiesFormat struct {
	// Description - A description for this rule. Restricted to 140 chars.
	Description *string `json:"description,omitempty"`
	// Protocol - Network protocol this rule applies to. Possible values include: 'SecurityRuleProtocolTCP', 'SecurityRuleProtocolUDP', 'SecurityRuleProtocolIcmp', 'SecurityRuleProtocolEsp', 'SecurityRuleProtocolAsterisk'
	Protocol SecurityRuleProtocol `json:"protocol,omitempty"`
	// SourcePortRange - The source port or range. Integer or range between 0 and 65535. Asterisk '*' can also be used to match all ports.
	SourcePortRange *string `json:"sourcePortRange,omitempty"`
	// DestinationPortRange - The destination port or range. Integer or range between 0 and 65535. Asterisk '*' can also be used to match all ports.
	DestinationPortRange *string `json:"destinationPortRange,omitempty"`
	// SourceAddressPrefix - The CIDR or source IP range. Asterisk '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used. If this is an ingress rule, specifies where network traffic originates from.
	SourceAddressPrefix *string `json:"sourceAddressPrefix,omitempty"`
	// SourceAddressPrefixes - The CIDR or source IP ranges.
	SourceAddressPrefixes *[]string `json:"sourceAddressPrefixes,omitempty"`
	// SourceApplicationSecurityGroups - The application security group specified as source.
	SourceApplicationSecurityGroups *[]ApplicationSecurityGroup `json:"sourceApplicationSecurityGroups,omitempty"`
	// DestinationAddressPrefix - The destination address prefix. CIDR or destination IP range. Asterisk '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used.
	DestinationAddressPrefix *string `json:"destinationAddressPrefix,omitempty"`
	// DestinationAddressPrefixes - The destination address prefixes. CIDR or destination IP ranges.
	DestinationAddressPrefixes *[]string `json:"destinationAddressPrefixes,omitempty"`
	// DestinationApplicationSecurityGroups - The application security group specified as destination.
	DestinationApplicationSecurityGroups *[]ApplicationSecurityGroup `json:"destinationApplicationSecurityGroups,omitempty"`
	// SourcePortRanges - The source port ranges.
	SourcePortRanges *[]string `json:"sourcePortRanges,omitempty"`
	// DestinationPortRanges - The destination port ranges.
	DestinationPortRanges *[]string `json:"destinationPortRanges,omitempty"`
	// Access - The network traffic is allowed or denied. Possible values include: 'SecurityRuleAccessAllow', 'SecurityRuleAccessDeny'
	Access SecurityRuleAccess `json:"access,omitempty"`
	// Priority - The priority of the rule. The value can be between 100 and 65500. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.
	Priority *uint32 `json:"priority,omitempty"`
	// Direction - The direction of the rule. The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values include: 'SecurityRuleDirectionInbound', 'SecurityRuleDirectionOutbound'
	Direction SecurityRuleDirection `json:"direction,omitempty"`
	// ProvisioningState - The provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// SecurityRule network security rule.

type SecurityRule struct {
	autorest.Response `json:"-"`
	// SecurityRulePropertiesFormat - Properties of the security rule.
	*SecurityRulePropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// SecurityGroupPropertiesFormat network Security Group resource.
type SecurityGroupPropertiesFormat struct {
	// SecurityRules - A collection of security rules of the network security group.
	SecurityRules *[]SecurityRule `json:"securityRules,omitempty"`
	// DefaultSecurityRules - The default security rules of network security group.
	DefaultSecurityRules *[]SecurityRule `json:"defaultSecurityRules,omitempty"`
	// NetworkInterfaces - READ-ONLY; A collection of references to network interfaces.
	NetworkInterfaces *[]Interface `json:"networkInterfaces,omitempty"`
	// Subnets - READ-ONLY; A collection of references to subnets.
	Subnets *[]Subnet `json:"subnets,omitempty"`
	// ResourceGUID - The resource GUID property of the network security group resource.
	ResourceGUID *string `json:"resourceGuid,omitempty"`
	// ProvisioningState - The provisioning state of the network security group resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// SecurityGroup networkSecurityGroup resource.
type SecurityGroup struct {
	autorest.Response `json:"-"`
	// SecurityGroupPropertiesFormat - Properties of the network security group.
	*SecurityGroupPropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// PrivateEndpointProperties properties of the private endpoint.
type PrivateEndpointProperties struct {
	// Subnet - The ID of the subnet from which the private IP will be allocated.
	Subnet *Subnet `json:"subnet,omitempty"`
	// NetworkInterfaces - READ-ONLY; Gets an array of references to the network interfaces created for this private endpoint.
	NetworkInterfaces *[]Interface `json:"networkInterfaces,omitempty"`
	// ProvisioningState - The provisioning state of the private endpoint. Possible values include: 'Succeeded', 'Updating', 'Deleting', 'Failed'
	ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// PrivateEndpoint private endpoint resource.
type PrivateEndpoint struct {
	autorest.Response `json:"-"`
	// PrivateEndpointProperties - Properties of the private endpoint.
	*PrivateEndpointProperties `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// InterfacePropertiesFormat networkInterface properties.
type InterfacePropertiesFormat struct {
	// VirtualMachine - READ-ONLY; The reference of a virtual machine.
	VirtualMachine *SubResource `json:"virtualMachine,omitempty"`
	// PrivateEndpoint - READ-ONLY; A reference to the private endpoint to which the network interface is linked.
	PrivateEndpoint *PrivateEndpoint `json:"privateEndpoint,omitempty"`
	// IPConfigurations - A list of IPConfigurations of the network interface.
	IPConfigurations *[]InterfaceIPConfiguration `json:"ipConfigurations,omitempty"`
	// DNSSettings - The DNS settings in network interface.
	DNSSettings *InterfaceDNSSettings `json:"dnsSettings,omitempty"`
	// MacAddress - The MAC address of the network interface.
	MacAddress *string `json:"macAddress,omitempty"`
	// Primary - Gets whether this is a primary network interface on a virtual machine.
	Primary *bool `json:"primary,omitempty"`
	// EnableAcceleratedNetworking - If the network interface is accelerated networking enabled.
	EnableAcceleratedNetworking *bool `json:"enableAcceleratedNetworking,omitempty"`
	// EnableIPForwarding - Indicates whether IP forwarding is enabled on this network interface.
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`
	// ProvisioningState - The provisioning state of the public IP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// EnableMACSpoofing - enable macspoofing on this nic
	EnableMACSpoofing *bool `json:"enableMACSpoofing,omitempty"`
	// EnableDHCPGuard
	EnableDHCPGuard *bool `json:"enableDHCPGuard,omitempty"`
	// EnableRouterAdvertisementGuard
	EnableRouterAdvertisementGuard *bool `json:"enableRouterAdvertisementGuard,omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type Interface struct {
	// InterfaceProperties - Properties of the network interface.
	*InterfacePropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// PublicIPAddressDNSSettings contains FQDN of the DNS record associated with the public IP address.
type PublicIPAddressDNSSettings struct {
	// DomainNameLabel - Gets or sets the Domain name label.The concatenation of the domain name label and the regionalized DNS zone make up the fully qualified domain name associated with the public IP address. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system.
	DomainNameLabel *string `json:"domainNameLabel,omitempty"`
	// Fqdn - Gets the FQDN, Fully qualified domain name of the A DNS record associated with the public IP. This is the concatenation of the domainNameLabel and the regionalized DNS zone.
	Fqdn *string `json:"fqdn,omitempty"`
	// ReverseFqdn - Gets or Sets the Reverse FQDN. A user-visible, fully qualified domain name that resolves to this public IP address. If the reverseFqdn is specified, then a PTR DNS record is created pointing from the IP address in the in-addr.arpa domain to the reverse FQDN.
	ReverseFqdn *string `json:"reverseFqdn,omitempty"`
}

// IPTag contains the IpTag associated with the object.
type IPTag struct {
	// IPTagType - Gets or sets the ipTag type: Example FirstPartyUsage.
	IPTagType *string `json:"ipTagType,omitempty"`
	// Tag - Gets or sets value of the IpTag associated with the public IP. Example SQL, Storage etc.
	Tag *string `json:"tag,omitempty"`
}

const DefaultIdleTimeoutInMinutes int32 = 4

// PublicIPAddressProperties public IP address properties.
type PublicIPAddressPropertiesFormat struct {
	// PublicIPAllocationMethod - The public IP address allocation method. Possible values include: 'Static', 'Dynamic'
	PublicIPAllocationMethod IPAllocationMethod `json:"publicIPAllocationMethod,omitempty"`
	// PublicIPAddressVersion - The public IP address version. Possible values include: 'IPv4', 'IPv6'
	PublicIPAddressVersion IPVersion `json:"publicIPAddressVersion,omitempty"`
	// IPConfiguration - READ-ONLY; The IP configuration associated with the public IP address.
	IPConfiguration *IPConfiguration `json:"ipConfiguration,omitempty"`
	// DNSSettings - The FQDN of the DNS record associated with the public IP address.
	DNSSettings *PublicIPAddressDNSSettings `json:"dnsSettings,omitempty"`
	// IPTags - The list of tags associated with the public IP address.
	IPTags *[]IPTag `json:"ipTags,omitempty"`
	// IPAddress - The IP address associated with the public IP address resource.
	IPAddress *string `json:"ipAddress,omitempty"`
	// PublicIPPrefix - The Public IP Prefix this Public IP Address should be allocated from.
	PublicIPPrefix *SubResource `json:"publicIPPrefix,omitempty"`
	// IdleTimeoutInMinutes - The idle timeout of the public IP address.
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
	// ProvisioningState - The provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// PublicIPAddress public IP address resource.
type PublicIPAddress struct {
	autorest.Response `json:"-"`
	// PublicIPAddressProperties - Public IP address properties.
	*PublicIPAddressPropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// Zones - A list of availability zones denoting the IP allocated for the resource needs to come from.
	Zones *[]string `json:"zones,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// VipPoolProperties vip pool properties.
type VipPoolPropertiesFormat struct {
	// IPPrefix - The IP Prefix for this Vip Pool
	IPPrefix *string `json:"IPPrefix,omitempty"`
	// StartIP - The starting IP address of this Vip Pool
	StartIP *string `json:"startIP,omitempty"`
	// EndIP - The ending IP address of this Vip Pool
	EndIP *string `json:"endIP,omitempty"`
	// ProvisioningState - The provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// VipPool vip pool resource.
type VipPool struct {
	autorest.Response `json:"-"`
	// VipPoolPropertiesFormat - Vip Pool properties.
	*VipPoolPropertiesFormat `json:"properties,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; Resource name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; Resource type.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Resource tags.
	Tags map[string]*string `json:"tags"`
}

// LogicalSubnet is associated with a Logical Network.
type LogicalSubnetPropertiesFormat struct {
	// CIDR for this subnet - IPv4, IPv6
	AddressPrefix *string `json:"addressPrefix,omitempty"`
	// AddressPrefixes - List of address prefixes for the subnet.
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
	// Routes for the subnet
	RouteTable *RouteTable `json:"routeTable,omitempty"`
	// IPConfiguration References
	IPConfigurationReferences *[]IPConfigurationReference `json:"ipConfigurationReferences,omitempty"`
	// IPAllocationMethod - The IP address allocation method. Possible values include: 'Static', 'Dynamic'
	IPAllocationMethod IPAllocationMethod `json:"ipAllocationMethod,omitempty"`
	// VLAN ID
	Vlan *uint16 `json:"vlan,omitempty"`
	// Pool of IP Addresses
	IPPools []IPPool `json:"ippools,omitempty"`
	// DhcpOptions - The dhcpOptions that contains an array of DNS servers available to VMs deployed in the Logical network.
	DhcpOptions *DhcpOptions `json:"dhcpOptions,omitempty"`
	// Public - Gets whether this is a public subnet on a virtual machine.
	Public *bool `json:"primary,omitempty"`
	// NetworkSecurityGroup - The reference of the NetworkSecurityGroup resource.
	NetworkSecurityGroup *SubResource `json:"networkSecurityGroup,omitempty"`
}

// LogicalSubnet is a subnet in a Logical network resource.
type LogicalSubnet struct {
	autorest.Response `json:"-"`
	// SubnetPropertiesFormat - Properties of the subnet.
	*LogicalSubnetPropertiesFormat `json:"properties,omitempty"`
	// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
	Name *string `json:"name,omitempty"`
	// Etag - A unique read-only string that changes whenever the resource is updated.
	Etag *string `json:"etag,omitempty"`
	// ID - Resource ID.
	ID *string `json:"id,omitempty"`
}

// LogicalNetworkPropertiesFormat properties of the Logical Network.
type LogicalNetworkPropertiesFormat struct {
	// AddressSpace - The AddressSpace that contains an array of IP address ranges that can be used by subnets.
	AddressSpace *AddressSpace `json:"addressSpace,omitempty"`
	// DhcpOptions - The dhcpOptions that contains an array of DNS servers available to VMs deployed in the Logical network.
	DhcpOptions *DhcpOptions `json:"dhcpOptions,omitempty"`
	// Subnets - A list of subnets in a Logical Network.
	Subnets *[]LogicalSubnet `json:"subnets,omitempty"`
	// MACPool name - Name of the associated MAC pool (or leave empty to use the default mac pool)
	MacPoolName *string `json:"macPoolName,omitempty"`
	// ProvisioningState - The provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// NetworkVirtualizationEnabled - Denotes if this lnet can be used as overlay for a vnet
	NetworkVirtualizationEnabled *bool `json:"networkVirtualizationEnabled,omitempty"`
}

// LogicalNetwork defines the structure of an LNET
type LogicalNetwork struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - Resource location.
	Location *string `json:"location,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// LogicalNetworkProperties - Properties of the Logical network.
	*LogicalNetworkPropertiesFormat `json:"properties,omitempty"`
}
