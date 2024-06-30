// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package publicipaddress

import (
	"context"
	"fmt"
	"strings"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/network"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudnetwork.PublicIPAddressAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newPublicIPAddressAgentClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetPublicIPAddressAgentClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get a public IP address by name.  If name is nil, get all public IP addresses
func (c *client) Get(ctx context.Context, group, name string) (*[]network.PublicIPAddress, error) {
	// Validate the inputs, but skip the name validation
	// because we want to get all public IP addresses if name is nil
	err := validateGroupAndNameInputs(group, name, true)
	if err != nil {
		return nil, err
	}

	request, err := c.getPublicIPAddressRequestByName(wssdcloudcommon.Operation_GET, group, name)
	if err != nil {
		return nil, err
	}

	response, err := c.PublicIPAddressAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	pips, err := c.getPublicIPAddressesFromResponse(response)
	if err != nil {
		return nil, err
	}

	return pips, nil
}

// CreateOrUpdate creates a public IP address if it does not exist, or updates an existing public IP address
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, inputPip *network.PublicIPAddress) (*network.PublicIPAddress, error) {
	err := validateInputsAndSetDefaults(group, name, inputPip)
	if err != nil {
		return nil, err
	}

	request, err := c.getPublicIPAddressRequest(wssdcloudcommon.Operation_POST, group, name, inputPip)
	if err != nil {
		return nil, err
	}

	response, err := c.PublicIPAddressAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	pips, err := c.getPublicIPAddressesFromResponse(response)
	if err != nil {
		return nil, err
	}

	if len(*pips) == 0 {
		return nil, fmt.Errorf("[PublicIPAddress][CreateOrUpdate] Public IP address [%s] not created or updated", name)
	}

	return &(*pips)[0], nil
}

// Delete a public IP address
func (c *client) Delete(ctx context.Context, group, name string) error {
	err := validateGroupAndNameInputs(group, name, false)
	if err != nil {
		return err
	}

	pips, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*pips) == 0 {
		return fmt.Errorf("[PublicIPAddress][Delete] Public IP address [%s] not found", name)
	}

	request, err := c.getPublicIPAddressRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*pips)[0])
	if err != nil {
		return err
	}
	_, err = c.PublicIPAddressAgentClient.Invoke(ctx, request)

	return err
}

func (c *client) Precheck(ctx context.Context, group string, pips []*network.PublicIPAddress) (bool, error) {
	request, err := getPublicIPAddressPrecheckRequest(group, pips)
	if err != nil {
		return false, err
	}

	response, err := c.PublicIPAddressAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getPublicIPAddressPrecheckResponse(response)
}

func getPublicIPAddressPrecheckRequest(group string, publicIPAddresses []*network.PublicIPAddress) (*wssdcloudnetwork.PublicIPAddressPrecheckRequest, error) {
	request := &wssdcloudnetwork.PublicIPAddressPrecheckRequest{}
	protoPips := make([]*wssdcloudnetwork.PublicIPAddress, 0, len(publicIPAddresses))

	for _, pip := range publicIPAddresses {
		// can public IP address ever be nil here? what would be the meaning of that?
		if pip != nil {
			protoPip, err := getWssdPublicIPAddress(pip, group)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert PublicIPAddress to Protobuf representation")
			}
			protoPips = append(protoPips, protoPip)
		}
	}

	request.PublicIPAddresses = protoPips
	return request, nil
}

func getPublicIPAddressPrecheckResponse(response *wssdcloudnetwork.PublicIPAddressPrecheckResponse) (bool, error) {
	result := response.GetResult().GetValue()
	if !result {
		return result, errors.New(response.GetError())
	}
	return result, nil
}

// validateInputsAndSetDefaults validates the input parameters and sets default values for a PublicIPAddress object.
func validateInputsAndSetDefaults(group, name string, inputPip *network.PublicIPAddress) error {
	err := validateGroupAndNameInputs(group, name, false)
	if err != nil {
		return err
	}

	if inputPip == nil || inputPip.PublicIPAddressPropertiesFormat == nil {
		return errors.Wrapf(errors.InvalidConfiguration, "Missing public IP address properties")
	}

	if inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes == nil {
		// Set default if not set
		var defaultValue int32 = network.DefaultIdleTimeoutInMinutes
		inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes = &defaultValue
	} else {
		if *inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes < 4 || *inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes > 30 {
			return errors.Wrapf(errors.InvalidInput, "Idle timeout in minute needs to be between 4 and 30 minutes")
		}
	}

	// In case of case insensitive comparison of IP address version, we need to convert the input to lower case
	// because CLI may send the input in upper case
	if len(inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion) == 0 ||
		strings.EqualFold(string(inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion), string(network.IPv4)) {
		// Set default if not set and make sure that the input is 'IPv4'
		inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion = network.IPv4
	} else {
		return errors.Wrapf(errors.InvalidInput, "Public IP address for IPv4 is only supported")
	}

	// In case of case insensitive comparison of IP allocation method, we need to convert the input to lower case
	// because CLI may send the input in upper case
	if len(inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod) == 0 ||
		strings.EqualFold(string(inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod), string(network.Dynamic)) {
		// Set default if not set and make sure that the input is 'Dynamic'
		inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod = network.Dynamic
	} else {
		return errors.Wrapf(errors.InvalidInput, "Public IP address for dynamic allocation is only supported")
	}

	return nil
}

// validateGroupAndNameInputs validates the input parameters for group and name.
// It checks if the group is specified and, if bypassName is false, it also checks if the name is specified.
// If any of the required parameters are missing, it returns an error.
func validateGroupAndNameInputs(group, name string, bypassName bool) error {
	if len(group) == 0 {
		return errors.Wrapf(errors.InvalidInput, "MOC Group not specified")
	}

	if !bypassName && len(name) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Public IP address name not specified")
	}

	return nil
}

// getPublicIPAddressRequestByName creates a PublicIPAddressRequest for a given operation type, resource group, and public IP address name.
// It initializes a PublicIPAddress object with the provided name and calls getPublicIPAddressRequest to generate the request.
func (c *client) getPublicIPAddressRequestByName(opType wssdcloudcommon.Operation, group, name string) (*wssdcloudnetwork.PublicIPAddressRequest, error) {
	pip := network.PublicIPAddress{
		Name: &name,
	}
	return c.getPublicIPAddressRequest(opType, group, name, &pip)
}

// getPublicIPAddressRequest converts our internal representation of a public IP address (network.PublicIPAddress) into a protobuf request (wssdcloudnetwork.PublicIPAddressRequest) that can be sent to wssdcloudagent
func (c *client) getPublicIPAddressRequest(opType wssdcloudcommon.Operation, group, name string, pip *network.PublicIPAddress) (*wssdcloudnetwork.PublicIPAddressRequest, error) {

	if pip == nil {
		return nil, errors.InvalidInput
	}

	request := &wssdcloudnetwork.PublicIPAddressRequest{
		OperationType:     opType,
		PublicIPAddresses: []*wssdcloudnetwork.PublicIPAddress{},
	}

	wssdCloudPip, err := getWssdPublicIPAddress(pip, group)
	if err != nil {
		return nil, err
	}

	request.PublicIPAddresses = append(request.PublicIPAddresses, wssdCloudPip)
	return request, nil
}

// GetPublicIPAddressesFromResponse converts a protobuf response from wssdcloudagent (wssdcloudnetwork.PublicIPAddressResponse) to out internal representation of a public IP address (network.PublicIPAddress)
func (c *client) getPublicIPAddressesFromResponse(response *wssdcloudnetwork.PublicIPAddressResponse) (*[]network.PublicIPAddress, error) {
	networkPips := []network.PublicIPAddress{}

	for _, wssdCloudPip := range response.GetPublicIPAddresses() {
		networkPip, err := getPublicIPAddress(wssdCloudPip)
		if err != nil {
			return nil, err
		}

		networkPips = append(networkPips, *networkPip)
	}

	return &networkPips, nil
}
