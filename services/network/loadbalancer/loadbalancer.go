// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package loadbalancer

import (
	"github.com/microsoft/moc-sdk-for-go/services/network"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

// getWssdLoadBalancer convert our internal representation of a loadbalancer (network.LoadBalancer) to the cloud load balancer protobuf used by wssdcloudagent (wssdnetwork.LoadBalancer)
func getWssdLoadBalancer(networkLB *network.LoadBalancer, group string, apiVersion *wssdcloudcommon.ApiVersion) (wssdCloudLB *wssdcloudnetwork.LoadBalancer, err error) {

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if networkLB.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for Load Balancer")
	}

	wssdCloudLB = &wssdcloudnetwork.LoadBalancer{
		Name:      *networkLB.Name,
		GroupName: group,
	}

	if networkLB.Version != nil {
		if wssdCloudLB.Status == nil {
			wssdCloudLB.Status = status.InitStatus()
		}
		wssdCloudLB.Status.Version.Number = *networkLB.Version
	}

	if networkLB.Location != nil {
		wssdCloudLB.LocationName = *networkLB.Location
	}

	if networkLB.Tags != nil {
		wssdCloudLB.Tags = tags.MapToProto(networkLB.Tags)
	}

	if apiVersion != nil && apiVersion.Major >= 2 {
		err = getLoadBalancerPropertiesV2(networkLB.LoadBalancerPropertiesFormat, wssdCloudLB)
	} else {
		err = getLoadBalancerPropertiesLegacy(networkLB.LoadBalancerPropertiesFormat, wssdCloudLB)
	}
	if err != nil {
		return nil, err
	}

	return wssdCloudLB, nil
}

// Parse the contents of network.LoadBalancerPropertiesFormat for a Legacy LB
// return an error if there is a config error, or pass the formatted values into
// the wssdcloudnetwork.LoadBalancer if they are valid
func getLoadBalancerPropertiesLegacy(lbp *network.LoadBalancerPropertiesFormat,
	wssdCloudLB *wssdcloudnetwork.LoadBalancer) (err error) {
	err = nil
	if lbp == nil {
		// Properties is empty, so there is nothing to parse
		return err
	}
	if lbp.BackendAddressPools != nil && len(*lbp.BackendAddressPools) > 0 {
		bap := *lbp.BackendAddressPools
		if bap[0].Name != nil {
			wssdCloudLB.Backendpoolnames = append(wssdCloudLB.Backendpoolnames, *bap[0].Name)
		}
	}
	if lbp.FrontendIPConfigurations != nil && len(*lbp.FrontendIPConfigurations) > 0 {
		fipc := *lbp.FrontendIPConfigurations
		if fipc[0].FrontendIPConfigurationPropertiesFormat != nil {
			fipcf := fipc[0].FrontendIPConfigurationPropertiesFormat
			if fipcf.Subnet != nil {
				subnet := fipcf.Subnet
				if subnet.ID != nil {
					wssdCloudLB.Networkid = *subnet.ID
				}
			}
			if fipcf.IPAddress != nil {
				wssdCloudLB.FrontendIP = *fipcf.IPAddress
			}
		}
	}
	if lbp.LoadBalancingRules != nil && len(*lbp.LoadBalancingRules) > 0 {
		rules := *lbp.LoadBalancingRules
		for _, rule := range rules {
			if rule.FrontendPort == nil {
				return errors.Wrapf(errors.InvalidInput, "Frontend port not specified")
			}
			if rule.BackendPort == nil {
				return errors.Wrapf(errors.InvalidInput, "Backend port not specified")
			}

			protocol := wssdcloudcommon.Protocol_All
			if string(rule.Protocol) != "" {
				protocol, err = getWssdProtocol(string(rule.Protocol))
				if err != nil {
					return err
				}
			}

			wssdCloudLBRule := &wssdcloudnetwork.LoadBalancingRule{
				FrontendPort: uint32(*rule.FrontendPort),
				BackendPort:  uint32(*rule.BackendPort),
				Protocol:     protocol,
			}
			wssdCloudLB.Loadbalancingrules = append(wssdCloudLB.Loadbalancingrules, wssdCloudLBRule)
		}
	}
	return err
}

// Parse the contents of network.LoadBalancerPropertiesFormat for a V2 LB
// return an error if there is a config error, or pass the formatted values into
// the wssdcloudnetwork.LoadBalancer if they are valid
func getLoadBalancerPropertiesV2(lbp *network.LoadBalancerPropertiesFormat,
	wssdCloudLB *wssdcloudnetwork.LoadBalancer) (err error) {
	err = nil
	if lbp == nil {
		// Properties is empty, so there is nothing to parse
		return err
	}

	// Backend Address Pools
	if lbp.BackendAddressPools != nil {
		for _, bap := range *lbp.BackendAddressPools {
			if bap.Name != nil && *bap.Name != "" {
				wssdCloudLB.BackendAddressPools = append(wssdCloudLB.BackendAddressPools,
					&wssdcloudnetwork.BackendAddressPool{
						Name: *bap.Name,
					})
			}
		}
	}
	// Frontend Ip Configurations
	if lbp.FrontendIPConfigurations != nil {
		for _, feIp := range *lbp.FrontendIPConfigurations {
			wssdCloudLBFIpC := &wssdcloudnetwork.FrontEndIpConfiguration{}
			if feIp.Name != nil && *feIp.Name != "" {
				wssdCloudLBFIpC.Name = *feIp.Name
			} else {
				return errors.Wrapf(errors.InvalidInput, "FrontendIPConfig Name not specified")
			}
			if feIp.FrontendIPConfigurationPropertiesFormat != nil &&
				feIp.FrontendIPConfigurationPropertiesFormat.PublicIPAddress != nil &&
				feIp.FrontendIPConfigurationPropertiesFormat.PublicIPAddress.ID != nil &&
				*feIp.FrontendIPConfigurationPropertiesFormat.PublicIPAddress.ID != "" {
				wssdCloudLBFIpC.PublicIPAddress = &wssdcloudcommon.PublicIPAddressReference{
					ResourceRef: &wssdcloudcommon.ResourceReference{
						Name: *feIp.FrontendIPConfigurationPropertiesFormat.PublicIPAddress.ID,
					},
				}
			} else {
				return errors.Wrapf(errors.InvalidInput, "FrontendIPConfig Public-Ip not specified")
			}
			wssdCloudLB.FrontendIpConfigurations = append(wssdCloudLB.FrontendIpConfigurations, wssdCloudLBFIpC)
		}
	}
	// LoadBalancing Rules
	if lbp.LoadBalancingRules != nil {
		for _, rule := range *lbp.LoadBalancingRules {
			if rule.FrontendPort == nil {
				return errors.Wrapf(errors.InvalidInput, "LB Rule Frontend port not specified")
			}
			if rule.BackendPort == nil {
				return errors.Wrapf(errors.InvalidInput, "LB Rule Backend port not specified")
			}

			protocol := wssdcloudcommon.Protocol_All
			if string(rule.Protocol) != "" {
				protocol, err = getWssdProtocol(string(rule.Protocol))
				if err != nil {
					return err
				}
			}

			if rule.Name == nil || *rule.Name == "" {
				return errors.Wrapf(errors.InvalidInput, "LB Rule Name not specified")
			}

			// Create rule object with required params
			wssdCloudLBRule := &wssdcloudnetwork.LoadBalancingRule{
				Name:         *rule.Name,
				FrontendPort: uint32(*rule.FrontendPort),
				BackendPort:  uint32(*rule.BackendPort),
				Protocol:     protocol,
			}

			// Add optional params
			if rule.IdleTimeoutInMinutes != nil {
				wssdCloudLBRule.IdleTimeoutInMinutes = uint32(*rule.IdleTimeoutInMinutes)
				if wssdCloudLBRule.IdleTimeoutInMinutes < 4 || wssdCloudLBRule.IdleTimeoutInMinutes > 30 {
					return errors.Wrapf(errors.InvalidInput, "LB Rule IdleTimeoutInMinutes %d outside accepted range (4 to 30)", wssdCloudLBRule.IdleTimeoutInMinutes)
				}
			} else {
				wssdCloudLBRule.IdleTimeoutInMinutes = 4
			}
			if string(rule.LoadDistribution) != "" {
				distribution, ok := wssdcloudnetwork.LoadDistribution_value[string(rule.LoadDistribution)]
				if !ok {
					return errors.Wrapf(errors.InvalidInput, "LB Rule Unknown LoadDistribution %s specified", rule.LoadDistribution)
				}
				wssdCloudLBRule.LoadDistribution = wssdcloudnetwork.LoadDistribution(distribution)
			} else {
				wssdCloudLBRule.LoadDistribution = wssdcloudnetwork.LoadDistribution_Default
			}
			if rule.FrontendIPConfiguration != nil && rule.FrontendIPConfiguration.ID != nil {
				wssdCloudLBRule.FrontendIpConfigurationsRef = []*wssdcloudcommon.FrontendIPConfigurationReference{
					{
						ResourceRef: &wssdcloudcommon.ResourceReference{
							Name: *rule.FrontendIPConfiguration.ID,
						},
					},
				}
			}
			if rule.BackendAddressPool != nil && rule.BackendAddressPool.ID != nil {
				wssdCloudLBRule.BackendAddressPoolRef = &wssdcloudcommon.BackendAddressPoolReference{
					ResourceRef: &wssdcloudcommon.ResourceReference{
						Name: *rule.BackendAddressPool.ID,
					},
				}
			}
			if rule.Probe != nil && rule.Probe.ID != nil {
				wssdCloudLBRule.ProbeRef = &wssdcloudcommon.ProbeReference{
					ResourceRef: &wssdcloudcommon.ResourceReference{
						Name: *rule.Probe.ID,
					},
				}
			}
			if rule.EnableFloatingIP != nil {
				wssdCloudLBRule.EnableFloatingIP = *rule.EnableFloatingIP
			}
			if rule.EnableTCPReset != nil {
				wssdCloudLBRule.EnableTcpReset = *rule.EnableTCPReset
			}
			wssdCloudLB.Loadbalancingrules = append(wssdCloudLB.Loadbalancingrules, wssdCloudLBRule)
		}
	}
	// Probes
	if lbp.Probes != nil {
		for _, probe := range *lbp.Probes {
			wssdCloudProbe := &wssdcloudnetwork.Probe{}
			if probe.Name != nil && *probe.Name != "" {
				wssdCloudProbe.Name = *probe.Name
			} else {
				return errors.Wrapf(errors.InvalidInput, "Probe Name not set")
			}
			if probe.Port != nil {
				wssdCloudProbe.Port = uint32(*probe.Port)
			} else {
				return errors.Wrapf(errors.InvalidInput, "Probe Port not set")
			}
			if probe.IntervalInSeconds != nil {
				wssdCloudProbe.IntervalInSeconds = uint32(*probe.IntervalInSeconds)
			} else {
				// Set Default
				wssdCloudProbe.IntervalInSeconds = 15
			}
			if probe.NumberOfProbes != nil {
				wssdCloudProbe.NumberOfProbes = uint32(*probe.NumberOfProbes)
			}
			if string(probe.Protocol) != "" {
				protocolInt, ok := wssdcloudnetwork.ProbeProtocol_value[string(probe.Protocol)]
				if !ok {
					// string not found in has of approved protocols
					return errors.Wrapf(errors.InvalidInput, "Probe Unknown protocol %s specified", probe.Protocol)
				}
				// Convert the int back into the Protocol enum
				wssdCloudProbe.Protocol = wssdcloudnetwork.ProbeProtocol(protocolInt)
			}
			if probe.RequestPath != nil && *probe.RequestPath != "" {
				wssdCloudProbe.RequestPath = &wssdcloudcommon.ProbeRequestPathReference{
					ResourceRef: &wssdcloudcommon.ResourceReference{
						Name: *probe.RequestPath,
					},
				}
			}
			wssdCloudLB.Probes = append(wssdCloudLB.Probes, wssdCloudProbe)
		}
	}
	// Outbound Rules
	if lbp.OutboundRules != nil {
		for _, outRule := range *lbp.OutboundRules {
			wssdCloudOutRule := &wssdcloudnetwork.LoadbalancerOutboundNatRule{}
			if outRule.Name != nil && *outRule.Name != "" {
				wssdCloudOutRule.Name = *outRule.Name
			} else {
				return errors.Wrapf(errors.InvalidInput, "Outbound Rule Name not set")
			}
			if outRule.EnableTCPReset != nil {
				wssdCloudOutRule.EnableTcpReset = *outRule.EnableTCPReset
			}
			if string(outRule.Protocol) != "" {
				wssdCloudOutRule.Protocol, err = getWssdProtocol(string(outRule.Protocol))
				if err != nil {
					return err
				}
			}
			if outRule.FrontendIPConfigurations != nil {
				for _, outfipc := range *outRule.FrontendIPConfigurations {
					if outfipc.ID != nil {
						wssdCloudOutRule.FrontendIpConfigurationsRef = append(
							wssdCloudOutRule.FrontendIpConfigurationsRef,
							&wssdcloudcommon.FrontendIPConfigurationReference{
								ResourceRef: &wssdcloudcommon.ResourceReference{
									Name: *outfipc.ID,
								},
							})
					}
				}
			}
			if outRule.BackendAddressPool != nil && outRule.BackendAddressPool.ID != nil {
				wssdCloudOutRule.BackendAddressPoolRef = &wssdcloudcommon.BackendAddressPoolReference{
					ResourceRef: &wssdcloudcommon.ResourceReference{
						Name: *outRule.BackendAddressPool.ID,
					},
				}
			}
			wssdCloudLB.OutboundNatRules = append(wssdCloudLB.OutboundNatRules, wssdCloudOutRule)
		}
	}
	return err
}

// getLoadBalancer converts the cloud load balancer protobuf returned from wssdcloudagent (wssdcloudnetwork.LoadBalancer) to our internal representation of a loadbalancer (network.LoadBalancer)
func getLoadBalancer(wssdLB *wssdcloudnetwork.LoadBalancer) (networkLB *network.LoadBalancer, err error) {
	networkLB = &network.LoadBalancer{
		Name:     &wssdLB.Name,
		Location: &wssdLB.LocationName,
		ID:       &wssdLB.Id,
		Version:  &wssdLB.Status.Version.Number,
		LoadBalancerPropertiesFormat: &network.LoadBalancerPropertiesFormat{
			Statuses:         status.GetStatuses(wssdLB.GetStatus()),
			ReplicationCount: wssdLB.GetReplicationCount(),
		},
		Tags: tags.ProtoToMap(wssdLB.Tags),
	}

	// V2 BackendAddressPool
	if len(wssdLB.BackendAddressPools) > 0 {
		backendAddressPools := []network.BackendAddressPool{}

		for _, backendPool := range wssdLB.BackendAddressPools {
			if backendPool != nil {
				backendAddressPools = append(backendAddressPools, network.BackendAddressPool{Name: &backendPool.Name})
			}
		}
		networkLB.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools
	} else if len(wssdLB.Backendpoolnames) > 0 {
		backendAddressPools := []network.BackendAddressPool{}

		for _, backendName := range wssdLB.Backendpoolnames {
			if backendName != "" {
				backendAddressPools = append(backendAddressPools, network.BackendAddressPool{Name: &backendName})
			}
		}
		networkLB.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools
	}

	if len(wssdLB.FrontendIP) != 0 || len(wssdLB.Networkid) != 0 {

		frontendipconfigurations := []network.FrontendIPConfiguration{
			{
				FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{},
			},
		}
		if len(wssdLB.FrontendIP) != 0 {
			frontendipconfigurations[0].FrontendIPConfigurationPropertiesFormat.IPAddress = &wssdLB.FrontendIP
		}
		if len(wssdLB.Networkid) != 0 {
			frontendipconfigurations[0].FrontendIPConfigurationPropertiesFormat.Subnet = &network.Subnet{ID: &wssdLB.Networkid}
		}
		networkLB.LoadBalancerPropertiesFormat.FrontendIPConfigurations = &frontendipconfigurations
	}

	// V2 FrontendIpConfigurations
	if len(wssdLB.FrontendIpConfigurations) > 0 {
		frontendipconfigurations := []network.FrontendIPConfiguration{}
		for _, wssdFeIpConf := range wssdLB.FrontendIpConfigurations {
			if wssdFeIpConf != nil {
				feIpConf := network.FrontendIPConfiguration{
					Name: toStringPtr(wssdFeIpConf.Name),
				}
				if wssdFeIpConf.PublicIPAddress != nil && wssdFeIpConf.PublicIPAddress.ResourceRef != nil {
					feIpConf.FrontendIPConfigurationPropertiesFormat = &network.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &network.PublicIPAddress{
							ID: toStringPtr(wssdFeIpConf.PublicIPAddress.ResourceRef.Name),
						},
					}
				}
				frontendipconfigurations = append(frontendipconfigurations, feIpConf)
			}
		}
		networkLB.LoadBalancerPropertiesFormat.FrontendIPConfigurations = &frontendipconfigurations
	}

	// Load Balancing Rules
	if len(wssdLB.Loadbalancingrules) > 0 {
		networkLBRules := []network.LoadBalancingRule{}

		for _, loadbalancingrule := range wssdLB.Loadbalancingrules {
			protocol, err := getNetworkProtocol(loadbalancingrule.Protocol)
			if err != nil {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", loadbalancingrule.Protocol)
			}
			loadDistributionStr, ok := wssdcloudnetwork.LoadDistribution_name[int32(loadbalancingrule.LoadDistribution)]
			if !ok {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown load distribution %s specified", loadbalancingrule.LoadDistribution)
			}
			loadDistribution := network.LoadDistribution(loadDistributionStr)

			networkLBRule := network.LoadBalancingRule{
				Name: toStringPtr(loadbalancingrule.Name),
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					FrontendPort:         toInt32Ptr(int32(loadbalancingrule.FrontendPort)),
					BackendPort:          toInt32Ptr(int32(loadbalancingrule.BackendPort)),
					Protocol:             protocol,
					IdleTimeoutInMinutes: toInt32Ptr(int32(loadbalancingrule.IdleTimeoutInMinutes)),
					EnableFloatingIP:     toBoolPtr(loadbalancingrule.EnableFloatingIP),
					EnableTCPReset:       toBoolPtr(loadbalancingrule.EnableTcpReset),
					LoadDistribution:     loadDistribution,
				},
			}

			if len(loadbalancingrule.FrontendIpConfigurationsRef) > 0 &&
				loadbalancingrule.FrontendIpConfigurationsRef[0] != nil &&
				loadbalancingrule.FrontendIpConfigurationsRef[0].ResourceRef != nil {
				networkLBRule.LoadBalancingRulePropertiesFormat.FrontendIPConfiguration = &network.SubResource{
					ID: toStringPtr(loadbalancingrule.FrontendIpConfigurationsRef[0].ResourceRef.Name),
				}
			}

			if loadbalancingrule.BackendAddressPoolRef != nil &&
				loadbalancingrule.BackendAddressPoolRef.ResourceRef != nil {
				networkLBRule.LoadBalancingRulePropertiesFormat.BackendAddressPool = &network.SubResource{
					ID: toStringPtr(loadbalancingrule.BackendAddressPoolRef.ResourceRef.Name),
				}
			}

			if loadbalancingrule.ProbeRef != nil &&
				loadbalancingrule.ProbeRef.ResourceRef != nil {
				networkLBRule.LoadBalancingRulePropertiesFormat.Probe = &network.SubResource{
					ID: toStringPtr(loadbalancingrule.ProbeRef.ResourceRef.Name),
				}
			}

			networkLBRules = append(networkLBRules, networkLBRule)
		}
		networkLB.LoadBalancerPropertiesFormat.LoadBalancingRules = &networkLBRules
	}

	// V1 Inbound Nate Rules
	if len(wssdLB.InboundNatRules) > 0 {
		networkInboundNatRules := []network.InboundNatRule{}

		for _, wssdInboundNatRule := range wssdLB.InboundNatRules {
			fePort := int32(wssdInboundNatRule.FrontendPort)
			bePort := int32(wssdInboundNatRule.BackendPort)
			var protocol network.TransportProtocol
			if wssdInboundNatRule.Protocol == wssdcloudcommon.Protocol_All {
				protocol = network.TransportProtocolAll
			} else if wssdInboundNatRule.Protocol == wssdcloudcommon.Protocol_Tcp {
				protocol = network.TransportProtocolTCP
			} else if wssdInboundNatRule.Protocol == wssdcloudcommon.Protocol_Udp {
				protocol = network.TransportProtocolUDP
			} else {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", wssdInboundNatRule.Protocol)
			}

			newNetworkInboundNatRule := network.InboundNatRule{
				Name: &wssdInboundNatRule.Name,
				InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
					FrontendPort: &fePort,
					BackendPort:  &bePort,
					Protocol:     protocol,
				},
			}

			networkInboundNatRules = append(networkInboundNatRules, newNetworkInboundNatRule)
		}

		networkLB.InboundNatRules = &networkInboundNatRules
	}

	// Probes
	if len(wssdLB.Probes) > 0 {
		networkProbes := []network.Probe{}
		for _, probe := range wssdLB.Probes {
			networkProbe := network.Probe{
				Name: toStringPtr(probe.Name),
				ProbePropertiesFormat: &network.ProbePropertiesFormat{
					Port:              toInt32Ptr(int32(probe.Port)),
					IntervalInSeconds: toInt32Ptr(int32(probe.IntervalInSeconds)),
					NumberOfProbes:    toInt32Ptr(int32(probe.NumberOfProbes)),
				},
			}
			if probe.RequestPath != nil && probe.RequestPath.ResourceRef != nil {
				networkProbe.ProbePropertiesFormat.RequestPath = toStringPtr(probe.RequestPath.ResourceRef.Name)
			}
			protocol, ok := wssdcloudnetwork.ProbeProtocol_name[int32(probe.Protocol)]
			if !ok {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", probe.Protocol)
			}
			networkProbe.ProbePropertiesFormat.Protocol = network.ProbeProtocol(protocol)
			networkProbes = append(networkProbes, networkProbe)
		}
		networkLB.LoadBalancerPropertiesFormat.Probes = &networkProbes
	}

	// Outbound Nat Rules
	if len(wssdLB.OutboundNatRules) > 0 {
		networkOutNatRules := []network.OutboundRule{}
		for _, outNatRule := range wssdLB.OutboundNatRules {
			networkOutNatRule := network.OutboundRule{
				Name: toStringPtr(outNatRule.Name),
				OutboundRulePropertiesFormat: &network.OutboundRulePropertiesFormat{
					EnableTCPReset: toBoolPtr(outNatRule.EnableTcpReset),
				},
			}
			//Protocol
			protocolStr, ok := wssdcloudcommon.Protocol_name[int32(outNatRule.Protocol)]
			if !ok {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified in outbound nat rule", outNatRule.Protocol)
			}
			networkOutNatRule.OutboundRulePropertiesFormat.Protocol = network.LoadBalancerOutboundRuleProtocol(protocolStr)
			//FrontendIPConfigurations
			if outNatRule.FrontendIpConfigurationsRef != nil &&
				len(outNatRule.FrontendIpConfigurationsRef) > 0 &&
				outNatRule.FrontendIpConfigurationsRef[0].ResourceRef != nil {
				networkOutNatRule.OutboundRulePropertiesFormat.FrontendIPConfigurations = &[]network.SubResource{
					{
						ID: toStringPtr(outNatRule.FrontendIpConfigurationsRef[0].ResourceRef.Name),
					},
				}
			}
			//BackendAddressPool
			if outNatRule.BackendAddressPoolRef != nil &&
				outNatRule.BackendAddressPoolRef.ResourceRef != nil {
				networkOutNatRule.OutboundRulePropertiesFormat.BackendAddressPool = &network.SubResource{
					ID: toStringPtr(outNatRule.BackendAddressPoolRef.ResourceRef.Name),
				}
			}
			networkOutNatRules = append(networkOutNatRules, networkOutNatRule)
		}
		networkLB.LoadBalancerPropertiesFormat.OutboundRules = &networkOutNatRules
	}

	return networkLB, nil
}

func getApiVersion(apiVersion string) (version *wssdcloudcommon.ApiVersion, err error) {

	switch {
	case apiVersion == Version_Default:
		fallthrough
	case apiVersion == Version_1_0:
		return nil, nil
	case apiVersion == Version_2_0:
		version = &wssdcloudcommon.ApiVersion{
			Major: 2,
			Minor: 0,
		}
		return version, nil
	}

	return nil, errors.Wrapf(errors.InvalidVersion, "Apiversion [%s] is unsupported", apiVersion)
}

// toStringPtr returns a pointer to the passed string
func toStringPtr(s string) *string {
	return &s
}

func toInt32Ptr(i int32) *int32 {
	return &i
}

func toBoolPtr(b bool) *bool {
	return &b
}

func getWssdProtocol(protocol string) (wssdcloudcommon.Protocol, error) {
	// Hash lookup where the protocol string is the key and the enum int is the value
	protocolInt, ok := wssdcloudcommon.Protocol_value[protocol]
	if !ok {
		// string not found in has of approved protocols
		return wssdcloudcommon.Protocol_All, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", protocol)
	}
	// Convert the int back into the Protocol enum
	return wssdcloudcommon.Protocol(protocolInt), nil
}

func getNetworkProtocol(wssdProtocol wssdcloudcommon.Protocol) (network.TransportProtocol, error) {
	protocolStr, exists := wssdcloudcommon.Protocol_name[int32(wssdProtocol)]
	if !exists {
		return network.TransportProtocolAll, errors.Wrapf(errors.InvalidInput, "Conversion Error, cannot convert wssd protocol to Network protocol, wssd not found")
	}
	sdnProtcol, exists := network.TransportProtocol_value[protocolStr]
	if !exists {
		return network.TransportProtocolAll, errors.Wrapf(errors.InvalidInput, "Conversion Error, cannot convert wssd protocol to Network protocol, Network not found")
	}

	return sdnProtcol, nil
}
