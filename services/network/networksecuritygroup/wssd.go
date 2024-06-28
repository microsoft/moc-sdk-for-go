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
func (c *client) Get(ctx context.Context, location, name string) (*[]network.SecurityGroup, error) {

	request, err := c.getNetworkSecurityGroupRequestByName(wssdcloudcommon.Operation_GET, location, name)
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
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, inputNSG *network.SecurityGroup) (*network.SecurityGroup, error) {

	if inputNSG == nil || inputNSG.SecurityGroupPropertiesFormat == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Network Security Group Properties")
	}

	if inputNSG.SecurityGroupPropertiesFormat.SecurityRules != nil {
		nameMap := map[string]bool{}
		for _, item := range *inputNSG.SecurityGroupPropertiesFormat.SecurityRules {
			_, alreadyExists := nameMap[*item.Name]
			if alreadyExists {
				return nil, errors.Wrapf(errors.InvalidConfiguration, "Network Security Group Rules cannot have duplicate names")
			}
			nameMap[name] = true
		}
	}

	if inputNSG.SecurityGroupPropertiesFormat.DefaultSecurityRules != nil {
		nameMap := map[string]bool{}
		for _, item := range *inputNSG.SecurityGroupPropertiesFormat.DefaultSecurityRules {
			_, alreadyExists := nameMap[*item.Name]
			if alreadyExists {
				return nil, errors.Wrapf(errors.InvalidConfiguration, "Network Security Group Default Rules cannot have duplicate names")
			}
			nameMap[name] = true
		}
	}

	request, err := c.getNetworkSecurityGroupRequest(wssdcloudcommon.Operation_POST, location, name, inputNSG)
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
func (c *client) Delete(ctx context.Context, location, name string) error {
	nsgs, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*nsgs) == 0 {
		return fmt.Errorf("Network Security Group [%s] not found", name)
	}

	request, err := c.getNetworkSecurityGroupRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*nsgs)[0])
	if err != nil {
		return err
	}
	_, err = c.NetworkSecurityGroupAgentClient.Invoke(ctx, request)

	return err
}

func (c *client) Precheck(ctx context.Context, location string, networkSecurityGroups []*network.SecurityGroup) (bool, error) {
	request, err := getNetworkSecurityGroupPrecheckRequest(location, networkSecurityGroups)
	if err != nil {
		return false, err
	}
	response, err := c.NetworkSecurityGroupAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getNetworkSecurityGroupPrecheckResponse(response)
}

func getNetworkSecurityGroupPrecheckRequest(location string, networkSecurityGroups []*network.SecurityGroup) (*wssdcloudnetwork.NetworkSecurityGroupPrecheckRequest, error) {
	request := &wssdcloudnetwork.NetworkSecurityGroupPrecheckRequest{}

	protoNSGs := make([]*wssdcloudnetwork.NetworkSecurityGroup, 0, len(networkSecurityGroups))

	for _, nsg := range networkSecurityGroups {
		// can nsg ever be nil here? what would be the meaning of that?
		if nsg != nil {
			protoNSG, err := getWssdNetworkSecurityGroup(nsg, location)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert NetworkSecurityGroup to Protobuf representation")
			}
			protoNSGs = append(protoNSGs, protoNSG)
		}
	}

	request.NetworkSecurityGroups = protoNSGs
	return request, nil
}

func getNetworkSecurityGroupPrecheckResponse(response *wssdcloudnetwork.NetworkSecurityGroupPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}

func (c *client) getNetworkSecurityGroupRequestByName(opType wssdcloudcommon.Operation, location, name string) (*wssdcloudnetwork.NetworkSecurityGroupRequest, error) {
	networkNSG := network.SecurityGroup{
		Name: &name,
	}
	return c.getNetworkSecurityGroupRequest(opType, location, name, &networkNSG)
}

// getNetworkSecurityGroupRequest converts our internal representation of a network security group (network.SecurityGroup) into a protobuf request (wssdcloudnetwork.NetworkSecurityGroupRequest) that can be sent to wssdcloudagent
func (c *client) getNetworkSecurityGroupRequest(opType wssdcloudcommon.Operation, location, name string, networkNSG *network.SecurityGroup) (*wssdcloudnetwork.NetworkSecurityGroupRequest, error) {

	if networkNSG == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.NetworkSecurityGroupRequest{
		OperationType:         opType,
		NetworkSecurityGroups: []*wssdcloudnetwork.NetworkSecurityGroup{},
	}
	var err error

	wssdCloudNSG, err := getWssdNetworkSecurityGroup(networkNSG, location)
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
func getWssdNetworkSecurityGroup(networkNSG *network.SecurityGroup, location string) (wssdCloudNSG *wssdcloudnetwork.NetworkSecurityGroup, err error) {

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if networkNSG.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for Network Security Group")
	}

	wssdCloudNSG = &wssdcloudnetwork.NetworkSecurityGroup{
		Name:         *networkNSG.Name,
		LocationName: location,
	}

	if networkNSG.Tags != nil {
		wssdCloudNSG.Tags = tags.MapToProto(networkNSG.Tags)
	}

	if networkNSG.SecurityGroupPropertiesFormat != nil {
		nsgRules, err := getWssdNetworkSecurityGroupRules(networkNSG.SecurityRules, false)
		if err != nil {
			return nil, err
		}
		defaultNsgRules, err := getWssdNetworkSecurityGroupRules(networkNSG.DefaultSecurityRules, true)
		if err != nil {
			return nil, err
		}
		wssdCloudNSG.Networksecuritygrouprules = append(nsgRules, defaultNsgRules...)
	}

	return wssdCloudNSG, nil
}

// getWssdNetworkSecurityGroupRules converts our internal representation of a networksecuritygroup rule (network.SecurityRule) to the cloud network security group rule protobuf used by wssdcloudagent (wssdnetwork.NetworkSecurityGroupRule)
func getWssdNetworkSecurityGroupRules(securityRules *[]network.SecurityRule, isDefault bool) (wssdNSGRules []*wssdcloudnetwork.NetworkSecurityGroupRule, err error) {
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
		wssdCloudNSGRule.IsDefaultRule = isDefault

		if rule.Description != nil {
			wssdCloudNSGRule.Description = *rule.Description
		}

		if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolAsterisk)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_All
		} else if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolTCP)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_Tcp
		} else if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolUDP)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_Udp
		} else if strings.EqualFold(string(rule.Protocol), string(network.SecurityRuleProtocolIcmp)) {
			wssdCloudNSGRule.Protocol = wssdcloudcommon.Protocol_Icmpv4
		} else {
			return nil, errors.Wrapf(errors.InvalidInput, "Unknown Protocol %s specified", rule.Protocol)
		}

		if rule.SourceAddressPrefix != nil {
			wssdCloudNSGRule.SourceAddressPrefix = *rule.SourceAddressPrefix
		} else if rule.SourceAddressPrefixes != nil {
			concatRule := ""
			for _, prefix := range *rule.SourceAddressPrefixes {
				concatRule += prefix
			}
			wssdCloudNSGRule.SourceAddressPrefix = concatRule
		}

		if rule.DestinationAddressPrefix != nil {
			wssdCloudNSGRule.DestinationAddressPrefix = *rule.DestinationAddressPrefix
		} else if rule.DestinationAddressPrefixes != nil {
			concatRule := ""
			for _, prefix := range *rule.DestinationAddressPrefixes {
				concatRule += prefix
			}
			wssdCloudNSGRule.DestinationAddressPrefix = concatRule
		}

		if rule.SourcePortRange != nil {
			wssdCloudNSGRule.SourcePortRange = *rule.SourcePortRange
		} else if rule.SourcePortRanges != nil {
			concatRule := ""
			for _, prefix := range *rule.SourcePortRanges {
				concatRule += prefix
			}
			wssdCloudNSGRule.SourcePortRange = concatRule
		}

		if rule.DestinationPortRange != nil {
			wssdCloudNSGRule.DestinationPortRange = *rule.DestinationPortRange
		} else if rule.DestinationPortRanges != nil {
			concatRule := ""
			for _, prefix := range *rule.DestinationPortRanges {
				concatRule += prefix
			}
			wssdCloudNSGRule.DestinationPortRange = concatRule
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
			wssdCloudNSGRule.Priority = 4096 // Max for Azure, which expects 100 to 4096
		}

		wssdNSGRules = append(wssdNSGRules, wssdCloudNSGRule)
	}
	return
}

func isValidPriority(priority uint32) bool {
	return priority >= 100 && priority <= 65500
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

	if wssdNSG.Tags != nil {
		networkNSG.Tags = tags.ProtoToMap(wssdNSG.Tags)
	}

	if len(wssdNSG.Networksecuritygrouprules) > 0 {
		networkNSGRules := []network.SecurityRule{}
		networkDefaultNSGRules := []network.SecurityRule{}

		for _, rule := range wssdNSG.Networksecuritygrouprules {
			name := rule.Name
			description := rule.Description
			protocol := network.SecurityRuleProtocolAsterisk
			action := network.SecurityRuleAccessDeny
			direction := network.SecurityRuleDirectionInbound
			priority := uint32(rule.GetPriority())

			if rule.Protocol == wssdcloudcommon.Protocol_All {
				protocol = network.SecurityRuleProtocolAsterisk
			} else if rule.Protocol == wssdcloudcommon.Protocol_Tcp {
				protocol = network.SecurityRuleProtocolTCP
			} else if rule.Protocol == wssdcloudcommon.Protocol_Udp {
				protocol = network.SecurityRuleProtocolUDP
			} else if rule.Protocol == wssdcloudcommon.Protocol_Icmpv4 {
				protocol = network.SecurityRuleProtocolIcmp
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

			if rule.Direction == wssdcloudnetwork.Direction_Inbound {
				direction = network.SecurityRuleDirectionInbound
			} else if rule.Direction == wssdcloudnetwork.Direction_Outbound {
				direction = network.SecurityRuleDirectionOutbound
			} else {
				return nil, errors.Wrapf(errors.InvalidInput, "Unknown Direction %s specified", rule.Direction)
			}

			securityRule := network.SecurityRule{
				Name: &name,
				SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
					Description:              &description,
					Protocol:                 protocol,
					SourceAddressPrefix:      &rule.SourceAddressPrefix,
					DestinationAddressPrefix: &rule.DestinationAddressPrefix,
					SourcePortRange:          &rule.SourcePortRange,
					DestinationPortRange:     &rule.DestinationPortRange,
					Access:                   action,
					Direction:                direction,
					Priority:                 &priority,
				},
			}

			if rule.IsDefaultRule {
				networkDefaultNSGRules = append(networkDefaultNSGRules, securityRule)
			} else {
				networkNSGRules = append(networkNSGRules, securityRule)
			}
		}
		networkNSG.SecurityGroupPropertiesFormat.SecurityRules = &networkNSGRules
		networkNSG.SecurityGroupPropertiesFormat.DefaultSecurityRules = &networkDefaultNSGRules
	}

	return networkNSG, nil
}
