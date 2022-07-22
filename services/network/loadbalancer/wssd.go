// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package loadbalancer

import (
	"context"
	"fmt"
	"strings"

	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/network"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
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
func (c *client) Get(ctx context.Context, group, name string) (*[]network.LoadBalancer, error) {

	request, err := c.getLoadBalancerRequestByName(wssdcloudcommon.Operation_GET, group, name)
	if err != nil {
		return nil, err
	}

	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	lbs, err := c.getLoadBalancersFromResponse(response)
	if err != nil {
		return nil, err
	}

	return lbs, nil

}

// CreateOrUpdate creates a load balancer if it does not exist, or updates an existing load balancer
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, inputLB *network.LoadBalancer) (*network.LoadBalancer, error) {

	if inputLB == nil || inputLB.LoadBalancerPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Load Balancer Properties")
	}

	request, err := c.getLoadBalancerRequest(wssdcloudcommon.Operation_POST, group, name, inputLB)
	if err != nil {
		return nil, err
	}
	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	lbs, err := c.getLoadBalancersFromResponse(response)
	if err != nil {
		return nil, err
	}

	return &(*lbs)[0], nil
}

// Delete a load balancer
func (c *client) Delete(ctx context.Context, group, name string) error {
	lbs, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*lbs) == 0 {
		return fmt.Errorf("Load Balancer [%s] not found", name)
	}

	request, err := c.getLoadBalancerRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*lbs)[0])
	if err != nil {
		return err
	}
	_, err = c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return err
	}

	return err
}

func (c *client) getLoadBalancerRequestByName(opType wssdcloudcommon.Operation, group, name string) (*wssdcloudnetwork.LoadBalancerRequest, error) {
	networkLB := network.LoadBalancer{
		Name: &name,
	}
	return c.getLoadBalancerRequest(opType, group, name, &networkLB)
}

// getLoadBalancerRequest converts our internal representation of a load balancer (network.LoadBalancer) into a protobuf request (wssdcloudnetwork.LoadBalancerRequest) that can be sent to wssdcloudagent
func (c *client) getLoadBalancerRequest(opType wssdcloudcommon.Operation, group, name string, networkLB *network.LoadBalancer) (*wssdcloudnetwork.LoadBalancerRequest, error) {

	if networkLB == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.LoadBalancerRequest{
		OperationType: opType,
		LoadBalancers: []*wssdcloudnetwork.LoadBalancer{},
	}
	var err error

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

	if networkLB.LoadBalancerPropertiesFormat != nil {
		lbp := networkLB.LoadBalancerPropertiesFormat
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
					return nil, errors.Wrapf(errors.InvalidInput, "Frontend port not specified")
				}
				if rule.BackendPort == nil {
					return nil, errors.Wrapf(errors.InvalidInput, "Backend port not specified")
				}

				protocol := wssdcloudcommon.Protocol_All

				if strings.EqualFold(string(rule.Protocol), string(network.TransportProtocolAll)) {
					protocol = wssdcloudcommon.Protocol_All
				} else if strings.EqualFold(string(rule.Protocol), string(network.TransportProtocolTCP)) {
					protocol = wssdcloudcommon.Protocol_Tcp
				} else if strings.EqualFold(string(rule.Protocol), string(network.TransportProtocolUDP)) {
					protocol = wssdcloudcommon.Protocol_Udp
				} else {
					return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", rule.Protocol)
				}

				wssdCloudLBRule := &wssdcloudnetwork.LoadBalancingRule{
					FrontendPort: uint32(*rule.FrontendPort),
					BackendPort:  uint32(*rule.BackendPort),
					Protocol:     protocol,
				}
				wssdCloudLB.Loadbalancingrules = append(wssdCloudLB.Loadbalancingrules, wssdCloudLBRule)
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

	if len(wssdLB.Backendpoolnames) > 0 {
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

	if len(wssdLB.Loadbalancingrules) > 0 {
		networkLBRules := []network.LoadBalancingRule{}

		for _, loadbalancingrule := range wssdLB.Loadbalancingrules {
			frontendport := int32(loadbalancingrule.FrontendPort)
			backendport := int32(loadbalancingrule.BackendPort)
			protocol := network.TransportProtocolAll

			if loadbalancingrule.Protocol == wssdcloudcommon.Protocol_All {
				protocol = network.TransportProtocolAll
			} else if loadbalancingrule.Protocol == wssdcloudcommon.Protocol_Tcp {
				protocol = network.TransportProtocolTCP
			} else if loadbalancingrule.Protocol == wssdcloudcommon.Protocol_Udp {
				protocol = network.TransportProtocolUDP
			} else {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown protocol %s specified", loadbalancingrule.Protocol)
			}
			networkLBRules = append(networkLBRules, network.LoadBalancingRule{
				LoadBalancingRulePropertiesFormat: &network.LoadBalancingRulePropertiesFormat{
					FrontendPort: &frontendport,
					BackendPort:  &backendport,
					Protocol:     protocol,
				},
			})
		}
		networkLB.LoadBalancerPropertiesFormat.LoadBalancingRules = &networkLBRules
	}

	return networkLB, nil
}
