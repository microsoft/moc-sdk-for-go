// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package networksecuritygroup

import (
	"context"
	"fmt"
	"strings"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/moc/pkg/tags"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudnetwork.NetworkSecurityGroupAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newNetworkSecurityGroupClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetNetworkSecurityGroupClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get network security groups by name.  If name is nil, get all network security groups
func (c *client) Get(ctx context.Context, group, name string) (*[]network.SecurityGroup, error) {

	request, err := c.getNetworkSecurityGroupRequestByName(wssdcloudcommon.Operation_GET, group, name)
	if err != nil {
		return nil, err
	}

	response, err := c.NetworkSecurityGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	nsgs, err := c.getNetworkSecurityGroupsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return nsgs, nil

}

// CreateOrUpdate creates a network security group if it does not exist, or updates an existing network security group
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, inputNSG *network.SecurityGroup) (*network.SecurityGroup, error) {

	if inputNSG == nil || inputNSG.SecurityGroupPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Security Group Properties")
	}

	request, err := c.getNetworkSecurityGroupRequest(wssdcloudcommon.Operation_POST, group, name, inputNSG)
	if err != nil {
		return nil, err
	}
	response, err := c.NetworkSecurityGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	nsgs, err := c.getNetworkSecurityGroupsFromResponse(response)
	if err != nil {
		return nil, err
	}

	return &(*nsgs)[0], nil
}

// Delete a network security group
func (c *client) Delete(ctx context.Context, group, name string) error {
	nsgs, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*nsgs) == 0 {
		return fmt.Errorf("Network Security Group [%s] not found", name)
	}

	request, err := c.getNetworkSecurityGroupRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*nsgs)[0])
	if err != nil {
		return err
	}
	_, err = c.NetworkSecurityGroupAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

func (c *client) getNetworkSecurityGroupRequestByName(opType wssdcloudcommon.Operation, group, name string) (*wssdcloudnetwork.NetworkSecurityGroupRequest, error) {
	networkNSG := network.SecurityGroup{
		Name: &name,
	}
	return c.getNetworkSecurityGroupRequest(opType, group, name, &networkNSG)
}

// getNetworkSecurityGroupRequest converts our internal representation of a network security group (network.SecurityGroup) into a protobuf request (wssdcloudnetwork.NetworkSecurityGroupRequest) that can be sent to wssdcloudagent
func (c *client) getNetworkSecurityGroupRequest(opType wssdcloudcommon.Operation, group, name string, networkNSG *network.SecurityGroup) (*wssdcloudnetwork.NetworkSecurityGroupRequest, error) {

	if networkNSG == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.NetworkSecurityGroupRequest{
		OperationType:         opType,
		NetworkSecurityGroups: []*wssdcloudnetwork.NetworkSecurityGroup{},
	}
	var err error

	wssdCloudNSG, err := getWssdNetworkSecurityGroup(networkNSG, group)
	if err != nil {
		return nil, err
	}

	request.NetworkSecurityGroups = append(request.NetworkSecurityGroups, wssdCloudNSG)
	return request, nil
}

// getNetworkSecurityGroupsFromResponse converts a protobuf response from wssdcloudagent (wssdcloudnetwork.NetworkSecurityGroupResponse) to out internal representation of a network security group (network.SecurityGroup)
func (c *client) getNetworkSecurityGroupsFromResponse(response *wssdcloudnetwork.NetworkSecurityGroupResponse) (*[]network.SecurityGroup, error) {
	networkdNSGs := []network.SecurityGroup{}

	for _, wssdCloudNSG := range response.GetNetworkSecurityGroups() {
		networkNSG, err := getNetworkSecurityGroup(wssdCloudNSG)
		if err != nil {
			return nil, err
		}

		networkdNSGs = append(networkdNSGs, *networkNSG)
	}

	return &networkdNSGs, nil
}

// getWssdNetworkSecurityGroup converts our internal representation of a networksecuritygroup (network.SecurityGroup) to the cloud network security group protobuf used by wssdcloudagent (wssdnetwork.NetworkSecurityGroup)
func getWssdNetworkSecurityGroup(networkNSG *network.SecurityGroup, group string) (wssdCloudNSG *wssdcloudnetwork.NetworkSecurityGroup, err error) {

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	if networkNSG.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for Network Security Group")
	}

	wssdCloudNSG = &wssdcloudnetwork.NetworkSecurityGroup{
		Name:      *networkNSG.Name,
		GroupName: group,
	}

	if networkNSG.Location != nil {
		wssdCloudNSG.LocationName = *networkNSG.Location
	}

	if networkNSG.Tags != nil {
		wssdCloudNSG.Tags = tags.MapToProto(networkNSG.Tags)
	}

	if networkNSG.SecurityGroupPropertiesFormat != nil {
		nsgRules, err := getWssdNetworkSecurityGroupRules(networkNSG.SecurityRules)
		if err != nil {
			return nil, err
		}
		wssdCloudNSG.Networksecuritygrouprules = nsgRules
	}

	return wssdCloudNSG, nil
}

// getWssdNetworkSecurityGroupRules converts our internal representation of a networksecuritygroup rule (network.SecurityRule) to the cloud network security group rule protobuf used by wssdcloudagent (wssdnetwork.NetworkSecurityGroupRule)
func getWssdNetworkSecurityGroupRules(securityRules *[]network.SecurityRule) (wssdNSGRules []*wssdcloudnetwork.NetworkSecurityGroupRule, err error) {
	if securityRules == nil || len(*securityRules) <= 0 {
		return
	}

	for _, rule := range *securityRules {
		if rule.SecurityRulePropertiesFormat == nil {
			continue
		}

		wssdCloudNSGRule := &wssdcloudnetwork.NetworkSecurityGroupRule{}

		if rule.Name == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Network Security Rule name not specified")
		}
		wssdCloudNSGRule.Name = *rule.Name

		if rule.Description != nil {
			wssdCloudNSGRule.Description = *rule.Description
		}

		if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolAsterisk)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_All
		} else if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolTCP)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_Tcp
		} else if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolUDP)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_Udp
		} else {
			return nil, errors.Wrapf(errors.InvalidInput, "Unknown Protocol %s specified", rule.Protocol)
		}

		if rule.SourceAddressPrefix != nil {
			wssdCloudNSGRule.SourceAddressPrefix = *rule.SourceAddressPrefix
		}

		if rule.DestinationAddressPrefix != nil {
			wssdCloudNSGRule.DestinationAddressPrefix = *rule.DestinationAddressPrefix
		}

		if rule.SourcePortRange != nil {
			wssdCloudNSGRule.SourcePortRange = *rule.SourcePortRange
		}

		if rule.DestinationPortRange != nil {
			wssdCloudNSGRule.DestinationPortRange = *rule.DestinationPortRange
		}

		if strings.EqualFold(string(rule.Access), string(network.SecurityRuleAccessAllow)) {
			wssdCloudNSGRule.Action = wssdcloudnetwork.Action_Allow
		} else if strings.EqualFold(string(rule.Access), string(network.SecurityRuleAccessDeny)) {
			wssdCloudNSGRule.Action = wssdcloudnetwork.Action_Deny
		} else {
			return nil, errors.Wrapf(errors.InvalidInput, "Unknown Access %s specified", rule.Access)
		}

		if strings.EqualFold(string(rule.Direction), string(network.SecurityRuleDirectionInbound)) {
			wssdCloudNSGRule.Direction = wssdcloudnetwork.Direction_Inbound
		} else if strings.EqualFold(string(rule.Direction), string(network.SecurityRuleDirectionOutbound)) {
			wssdCloudNSGRule.Direction = wssdcloudnetwork.Direction_Outbound
		} else {
			return nil, errors.Wrapf(errors.InvalidInput, "Unknown Direction %s specified", rule.Access)
		}

		if rule.Priority != nil && isValidPriority(*rule.Priority) {
			wssdCloudNSGRule.Priority = uint32(*rule.Priority)
		} else {
			wssdCloudNSGRule.Priority = 4096 // TODO: what should be the default value?
		}

		wssdNSGRules = append(wssdNSGRules, wssdCloudNSGRule)
	}
	return
}

func isValidPriority(priority uint32) bool {
	return priority >= 100 && priority <= 4096
}

// getNetworkSecurityGroup converts the cloud network security group protobuf returned from wssdcloudagent (wssdcloudnetwork.NetworkSecurityGroup) to our internal representation of a networksecuritygroup (network.SecurityGroup)
func getNetworkSecurityGroup(wssdNSG *wssdcloudnetwork.NetworkSecurityGroup) (networkNSG *network.SecurityGroup, err error) {
	networkNSG = &network.SecurityGroup{
		Name:     &wssdNSG.Name,
		Location: &wssdNSG.LocationName,
		ID:       &wssdNSG.Id,
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			Statuses: status.GetStatuses(wssdNSG.GetStatus()),
		},
	}

	if len(wssdNSG.Networksecuritygrouprules) > 0 {
		networkNSGRules := []network.SecurityRule{}

		for _, rule := range wssdNSG.Networksecuritygrouprules {
			name := rule.Name
			description := rule.Description
			protocol := network.SecurityRuleProtocolAsterisk
			action := network.SecurityRuleAccessDeny
			priority := uint32(rule.GetPriority())

			if rule.Protocol == wssdcloudcommon.Protocol_All {
				protocol = network.SecurityRuleProtocolAsterisk
			} else if rule.Protocol == wssdcloudcommon.Protocol_Tcp {
				protocol = network.SecurityRuleProtocolTCP
			} else if rule.Protocol == wssdcloudcommon.Protocol_Udp {
				protocol = network.SecurityRuleProtocolUDP
			} else {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown Protocol %s specified", rule.Protocol)
			}

			if rule.Action == wssdcloudnetwork.Action_Allow {
				action = network.SecurityRuleAccessAllow
			} else if rule.Action == wssdcloudnetwork.Action_Deny {
				action = network.SecurityRuleAccessDeny
			} else {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown Access %s specified", rule.Action)
			}

			networkNSGRules = append(networkNSGRules, network.SecurityRule{
				Name: &name,
				SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
					Description:              &description,
					Protocol:                 protocol,
					SourceAddressPrefix:      &rule.SourceAddressPrefix,
					DestinationAddressPrefix: &rule.DestinationAddressPrefix,
					SourcePortRange:          &rule.SourcePortRange,
					DestinationPortRange:     &rule.DestinationPortRange,
					Access:                   action,
					Priority:                 &priority,
				},
			})
		}
		networkNSG.SecurityGroupPropertiesFormat.SecurityRules = &networkNSGRules
	}

	return networkNSG, nil
}
