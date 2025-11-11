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
	return c.internalPrecheckWithVersion(ctx, group, loadBalancers, Version_Default)
}

func (c *client) PrecheckWithVersion(ctx context.Context, group string, loadBalancers []*network.LoadBalancer, apiVersion string) (bool, error) {
	return c.internalPrecheckWithVersion(ctx, group, loadBalancers, apiVersion)
}

func (c *client) internalPrecheckWithVersion(ctx context.Context, group string, loadBalancers []*network.LoadBalancer, apiVersion string) (bool, error) {
	request, err := getLoadBalancerPrecheckRequest(group, loadBalancers, apiVersion)
	if err != nil {
		return false, err
	}
	response, err := c.LoadBalancerAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getLoadBalancerPrecheckResponse(response)
}

func getLoadBalancerPrecheckRequest(group string, loadBalancers []*network.LoadBalancer, apiVersion string) (*wssdcloudnetwork.LoadBalancerPrecheckRequest, error) {
	request := &wssdcloudnetwork.LoadBalancerPrecheckRequest{}

	version, err := getApiVersion(apiVersion)
	if err != nil {
		return nil, err
	}

	protoLoadBalancers := make([]*wssdcloudnetwork.LoadBalancer, 0, len(loadBalancers))

	for _, lb := range loadBalancers {
		// can lb ever be nil here? what would be the meaning of that?
		if lb != nil {
			protoLB, err := getWssdLoadBalancer(lb, group, version)
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

	wssdCloudLB, err := getWssdLoadBalancer(networkLB, group, version)
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
