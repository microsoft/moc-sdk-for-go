// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/network"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	"github.com/microsoft/moc/rpc/common"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

const (
	// supported API Versions for Load Balancers
	Version_Default = ""
	Version_1_0     = "1.0"
	Version_2_0     = "2.0"
)

type client struct {
	wssdcloudnetwork.LoadBalancerAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newLoadBalancerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetLoadBalancerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get load balancers by name.  If name is nil, get all load balancers
func (c *client) GetWithVersion(ctx context.Context, group, name, apiVersion string) (*[]network.LoadBalancer, error) {
	return c.internalGetWithVersion(ctx, group, name, apiVersion)
}

func (c *client) Get(ctx context.Context, group, name string) (*[]network.LoadBalancer, error) {
	return c.internalGetWithVersion(ctx, group, name, Version_Default)
}

func (c *client) internalGetWithVersion(ctx context.Context, group, name, apiVersion string) (*[]network.LoadBalancer, error) {
	request, err := c.getLoadBalancerRequestByName(wssdcloudcommon.Operation_GET, group, name, apiVersion)
	if err != nil {
		return nil, err
	}

	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	lbs, err := c.getLoadBalancersFromResponse(response)
	if err != nil {
		return nil, err
	}

	return lbs, nil
}

// CreateOrUpdate creates a load balancer if it does not exist, or updates an existing load balancer
func (c *client) CreateOrUpdateWithVersion(ctx context.Context, group, name string, inputLB *network.LoadBalancer, apiVersion string) (*network.LoadBalancer, error) {
	return c.internalCreateOrUpdateWithVersion(ctx, group, name, inputLB, apiVersion)
}

func (c *client) CreateOrUpdate(ctx context.Context, group, name string, inputLB *network.LoadBalancer) (*network.LoadBalancer, error) {
	return c.internalCreateOrUpdateWithVersion(ctx, group, name, inputLB, Version_Default)
}

func (c *client) internalCreateOrUpdateWithVersion(ctx context.Context, group, name string, inputLB *network.LoadBalancer, apiVersion string) (*network.LoadBalancer, error) {

	if inputLB == nil || inputLB.LoadBalancerPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Load Balancer Properties")
	}

	request, err := c.getLoadBalancerRequest(wssdcloudcommon.Operation_POST, group, name, inputLB, apiVersion)
	if err != nil {
		return nil, err
	}
	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	lbs, err := c.getLoadBalancersFromResponse(response)
	if err != nil {
		return nil, err
	}

	return &(*lbs)[0], nil
}

// Delete a load balancer
func (c *client) DeleteWithVersion(ctx context.Context, group, name, apiVersion string) error {
	return c.internalDeleteWithVersion(ctx, group, name, apiVersion)
}

func (c *client) Delete(ctx context.Context, group, name string) error {
	return c.internalDeleteWithVersion(ctx, group, name, Version_Default)
}

func (c *client) internalDeleteWithVersion(ctx context.Context, group, name, apiVersion string) error {
	lbs, err := c.GetWithVersion(ctx, group, name, apiVersion)
	if err != nil {
		return err
	}
	if len(*lbs) == 0 {
		return fmt.Errorf("Load Balancer [%s] not found", name)
	}

	request, err := c.getLoadBalancerRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*lbs)[0], apiVersion)
	if err != nil {
		return err
	}
	_, err = c.LoadBalancerAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

func (c *client) Precheck(ctx context.Context, group string, loadBalancers []*network.LoadBalancer) (bool, error) {
	request, err := getLoadBalancerPrecheckRequest(group, loadBalancers)
	if err != nil {
		return false, err
	}
	response, err := c.LoadBalancerAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getLoadBalancerPrecheckResponse(response)
}

func getLoadBalancerPrecheckRequest(group string, loadBalancers []*network.LoadBalancer) (*wssdcloudnetwork.LoadBalancerPrecheckRequest, error) {
	request := &wssdcloudnetwork.LoadBalancerPrecheckRequest{}

	protoLoadBalancers := make([]*wssdcloudnetwork.LoadBalancer, 0, len(loadBalancers))

	for _, lb := range loadBalancers {
		// can lb ever be nil here? what would be the meaning of that?
		if lb != nil {
			protoLB, err := getWssdLoadBalancer(lb, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert LoadBalancer to Protobuf representation")
			}
			protoLoadBalancers = append(protoLoadBalancers, protoLB)
		}
	}

	request.LoadBalancers = protoLoadBalancers
	return request, nil
}

func getLoadBalancerPrecheckResponse(response *wssdcloudnetwork.LoadBalancerPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}

func (c *client) getLoadBalancerRequestByName(opType wssdcloudcommon.Operation, group, name, apiVersion string) (*wssdcloudnetwork.LoadBalancerRequest, error) {
	networkLB := network.LoadBalancer{
		Name: &name,
	}
	return c.getLoadBalancerRequest(opType, group, name, &networkLB, apiVersion)
}

// getLoadBalancerRequest converts our internal representation of a load balancer (network.LoadBalancer) into a protobuf request (wssdcloudnetwork.LoadBalancerRequest) that can be sent to wssdcloudagent
func (c *client) getLoadBalancerRequest(opType wssdcloudcommon.Operation, group, name string, networkLB *network.LoadBalancer, apiVersion string) (*wssdcloudnetwork.LoadBalancerRequest, error) {

	var err error
	var version *common.ApiVersion

	if version, err = getApiVersion(apiVersion); err != nil {
		return nil, err
	}

	if networkLB == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.LoadBalancerRequest{
		OperationType: opType,
		LoadBalancers: []*wssdcloudnetwork.LoadBalancer{},
		Version:       version,
	}

	wssdCloudLB, err := getWssdLoadBalancer(networkLB, group)
	if err != nil {
		return nil, err
	}

	request.LoadBalancers = append(request.LoadBalancers, wssdCloudLB)
	return request, nil
}

// getLoadBalancersFromResponse converts a protobuf response from wssdcloudagent (wssdcloudnetwork.LoadBalancerResponse) to out internal representation of a load balancer (network.LoadBalancer)
func (c *client) getLoadBalancersFromResponse(response *wssdcloudnetwork.LoadBalancerResponse) (*[]network.LoadBalancer, error) {
	networkLBs := []network.LoadBalancer{}

	for _, wssdCloudLB := range response.GetLoadBalancers() {
		networkLB, err := getLoadBalancer(wssdCloudLB)
		if err != nil {
			return nil, err
		}

		networkLBs = append(networkLBs, *networkLB)
	}

	return &networkLBs, nil
}

// getWssdLoadBalancer convert our internal representation of a loadbalancer (network.LoadBalancer) to the cloud load balancer protobuf used by wssdcloudagent (wssdnetwork.LoadBalancer)
func getWssdLoadBalancer(networkLB *network.LoadBalancer, group string) (wssdCloudLB *wssdcloudnetwork.LoadBalancer, err error) {

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

	if networkLB.LoadBalancerPropertiesFormat != nil {
		lbp := networkLB.LoadBalancerPropertiesFormat
		// Backend Address Pool. Populate the string array and the BackendAddressPool object
		if lbp.BackendAddressPools != nil && len(*lbp.BackendAddressPools) > 0 {
			bap := *lbp.BackendAddressPools
			if bap[0].Name != nil {
				wssdCloudLB.Backendpoolnames = append(wssdCloudLB.Backendpoolnames, *bap[0].Name)
			}
			// Populate the new object
			for _, bap := range *lbp.BackendAddressPools {
				if bap.Name != nil {
					wssdCloudLB.BackendAddressPools = append(wssdCloudLB.BackendAddressPools,
						&wssdcloudnetwork.BackendAddressPool{
							Name: *bap.Name,
						})
				}
			}

		}
		// Frontend Ip Configurations
		if lbp.FrontendIPConfigurations != nil && len(*lbp.FrontendIPConfigurations) > 0 {
			fipc := *lbp.FrontendIPConfigurations
			if fipc[0].FrontendIPConfigurationPropertiesFormat != nil {
				fipcf := fipc[0].FrontendIPConfigurationPropertiesFormat
				// New FrontendIpConfigurations
				wssdCloudLBFIpC := &wssdcloudnetwork.FrontEndIpConfiguration{}
				if fipcf.Subnet != nil {
					subnet := fipcf.Subnet
					if subnet.ID != nil {
						wssdCloudLB.Networkid = *subnet.ID
					}
				}
				if fipcf.IPAddress != nil {
					wssdCloudLB.FrontendIP = *fipcf.IPAddress
				}
				if fipc[0].Name != nil {
					wssdCloudLBFIpC.Name = *fipc[0].Name
				}
				if fipcf.PublicIPAddress != nil {
					if fipcf.PublicIPAddress.Name != nil {
						wssdCloudLBFIpC.PublicIPAddress = &wssdcloudcommon.PublicIPAddressReference{
							ResourceRef: &wssdcloudcommon.ResourceReference{
								Name: *fipcf.PublicIPAddress.Name,
							},
						}
					}
					if fipcf.PublicIPAddress.Type != nil {
						allocMethod, ok := wssdcloudcommon.IPAllocationMethod_value[*(fipcf.PublicIPAddress.Type)]
						if !ok {
							return nil, errors.Wrapf(errors.InvalidInput, "Unknown Allocation Method %s specified", *(fipcf.PublicIPAddress.Type))
						}
						wssdCloudLBFIpC.AllocationMethod = wssdcloudcommon.IPAllocationMethod(allocMethod)
					}
				}

				wssdCloudLB.FrontendIpConfigurations = append(wssdCloudLB.FrontendIpConfigurations, wssdCloudLBFIpC)
			}
		}
		// LoadBalancing Rules
		if lbp.LoadBalancingRules != nil && len(*lbp.LoadBalancingRules) > 0 {
			rules := *lbp.LoadBalancingRules
			for _, rule := range rules {
				if rule.FrontendPort == nil {
					return nil, errors.Wrapf(errors.InvalidInput, "Frontend port not specified")
				}
				if rule.BackendPort == nil {
					return nil, errors.Wrapf(errors.InvalidInput, "Backend port not specified")
				}

				protocol := wssdcloudcommon.Protocol_All
				if string(rule.Protocol) != "" {
					protocol, err = getWssdProtocol(string(rule.Protocol))
					if err != nil {
						return nil, err
					}
				}

				// Create rule object with required params
				wssdCloudLBRule := &wssdcloudnetwork.LoadBalancingRule{
					FrontendPort: uint32(*rule.FrontendPort),
					BackendPort:  uint32(*rule.BackendPort),
					Protocol:     protocol,
				}

				// Add optional params
				if rule.IdleTimeoutInMinutes != nil {
					wssdCloudLBRule.IdleTimeoutInMinutes = uint32(*rule.IdleTimeoutInMinutes)
				} else {
					wssdCloudLBRule.IdleTimeoutInMinutes = 4
				}
				if string(rule.LoadDistribution) != "" {
					distribution, ok := wssdcloudnetwork.LoadDistribution_value[string(rule.LoadDistribution)]
					if !ok {
						return nil, errors.Wrapf(errors.InvalidInput, "Unknown LoadDistribution %s specified", rule.LoadDistribution)
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
				if rule.Name != nil {
					wssdCloudLBRule.Name = *rule.Name
				}
				wssdCloudLB.Loadbalancingrules = append(wssdCloudLB.Loadbalancingrules, wssdCloudLBRule)
			}
		}
		if lbp.Probes != nil && len(*lbp.Probes) > 0 {
			probes := *lbp.Probes
			for _, probe := range probes {
				wssdCloudProbe := &wssdcloudnetwork.Probe{}
				if probe.Name != nil {
					wssdCloudProbe.Name = *probe.Name
				}
				if probe.Port != nil {
					wssdCloudProbe.Port = uint32(*probe.Port)
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
						return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", probe.Protocol)
					}
					// Convert the int back into the Protocol enum
					wssdCloudProbe.Protocol = wssdcloudnetwork.ProbeProtocol(protocolInt)
				}
				if probe.RequestPath != nil {
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
		if lbp.OutboundRules != nil && len(*lbp.OutboundRules) > 0 {
			OutboundRules := *lbp.OutboundRules
			for _, outRule := range OutboundRules {
				wssdCloudOutRule := &wssdcloudnetwork.LoadbalancerOutboundNatRule{}
				if outRule.Name != nil {
					wssdCloudOutRule.Name = *outRule.Name
				}
				if outRule.EnableTCPReset != nil {
					wssdCloudOutRule.EnableTcpReset = *outRule.EnableTCPReset
				}
				if string(outRule.Protocol) != "" {
					wssdCloudOutRule.Protocol, err = getWssdProtocol(string(outRule.Protocol))
					if err != nil {
						return nil, err
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
	}

	return wssdCloudLB, nil
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
					FrontendIPConfigurationPropertiesFormat: &network.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &network.PublicIPAddress{
							Type: toStringPtr(wssdcloudcommon.IPAllocationMethod_name[int32(wssdFeIpConf.AllocationMethod)]),
						},
					},
				}
				if wssdFeIpConf.PublicIPAddress != nil && wssdFeIpConf.PublicIPAddress.ResourceRef != nil {
					feIpConf.FrontendIPConfigurationPropertiesFormat.PublicIPAddress.Name =
						toStringPtr(wssdFeIpConf.PublicIPAddress.ResourceRef.Name)
				}
				frontendipconfigurations = append(frontendipconfigurations, feIpConf)
			}
		}
		networkLB.LoadBalancerPropertiesFormat.FrontendIPConfigurations = &frontendipconfigurations
	}

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

			networkLBRules = append(networkLBRules, network.LoadBalancingRule{
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
			})
		}
		networkLB.LoadBalancerPropertiesFormat.LoadBalancingRules = &networkLBRules
	}

	if len(wssdLB.InboundNatRules) > 0 {
		networkInboundNatRules := []network.InboundNatRule{}

		for _, wssdInboundNatRule := range wssdLB.InboundNatRules {
			fePort := int32(wssdInboundNatRule.FrontendPort)
			bePort := int32(wssdInboundNatRule.BackendPort)
			protocol := network.TransportProtocolAll
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
