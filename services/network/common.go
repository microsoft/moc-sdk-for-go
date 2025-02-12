package network

import (
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
)

// Returns cloudagent representation of AdvancedNetworkPolicy from sdk representation
func GetWssdAdvancedNetworkPolicies(policies *[]AdvancedNetworkPolicy) (wssdPolicies []*wssdcommonproto.AdvancedNetworkPolicy) {

	if policies == nil || len(*policies) == 0 {
		return nil
	}

	wssdPolicies = []*wssdcommonproto.AdvancedNetworkPolicy{}

	for _, policy := range *policies {
		wssdPolicy := &wssdcommonproto.AdvancedNetworkPolicy{
			Type:    getWssdPolicyType(policy.Type),
			Enabled: policy.Enabled,
		}
		wssdPolicies = append(wssdPolicies, wssdPolicy)
	}
	return wssdPolicies
}

// Returns sdk representation of AdvancedNetworkPolicy from cloudagent representation
func GetNetworkAdvancedNetworkPolicies(wssdPolicies []*wssdcommonproto.AdvancedNetworkPolicy) (policies []AdvancedNetworkPolicy) {

	if len(wssdPolicies) == 0 {
		return nil
	}

	policies = []AdvancedNetworkPolicy{}

	for _, wssdPolicy := range wssdPolicies {
		policy := AdvancedNetworkPolicy{
			Type:    getNetworkPolicyType(wssdPolicy.Type),
			Enabled: wssdPolicy.Enabled,
		}
		policies = append(policies, policy)
	}
	return policies
}

// Converts policy type from sdk to cloudagent representation
func getWssdPolicyType(policyType NetworkPolicyType) wssdcommonproto.NetworkPolicyType {
	switch policyType {
	case NetworkPolicyType_SDN:
		return wssdcommonproto.NetworkPolicyType_SDN
	default:
		return wssdcommonproto.NetworkPolicyType_INVALID
	}
}

// Converts policy type from cloudagent to sdk representation
func getNetworkPolicyType(wssdPolicyType wssdcommonproto.NetworkPolicyType) NetworkPolicyType {
	switch wssdPolicyType {
	case wssdcommonproto.NetworkPolicyType_SDN:
		return NetworkPolicyType_SDN
	default:
		return NetworkPolicyType_Invalid
	}
}
