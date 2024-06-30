// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package publicipaddress

import (
	"context"
	"fmt"

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

	err := ValidateInputs(ctx, group, name, inputPip)
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

	return &(*pips)[0], nil
}

// Delete a public IP address
func (c *client) Delete(ctx context.Context, group, name string) error {
	pips, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*pips) == 0 {
		return fmt.Errorf("Public IP address [%s] not found", name)
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

// Validate inputs
func ValidateInputs(ctx context.Context, group, name string, inputPip *network.PublicIPAddress) error {

	if len(group) == 0 {
		return errors.Wrapf(errors.InvalidInput, "MOC Group not specified")
	}

	if len(name) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Public IP address name not specified")
	}

	if inputPip == nil || inputPip.PublicIPAddressPropertiesFormat == nil {
		return errors.Wrapf(errors.InvalidConfiguration, "Missing public IP address properties")
	}

	if inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes == nil {
		//return errors.Wrapf(errors.InvalidInput, "Idle timeout in minute needs to be between 4 and 30 minutes")
		// Set default if not set
		*inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes = 4
	} else {
		if *inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes < 4 || *inputPip.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes > 30 {
			return errors.Wrapf(errors.InvalidInput, "Idle timeout in minute needs to be between 4 and 30 minutes")
		}
	}

	if len(inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion) == 0 {
		inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion = network.IPv4
	} else {
		if inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion != network.IPv4 && inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion != network.IPv6 {
			return errors.Wrapf(errors.InvalidInput, "Invalid IP address version for public IP address")
		}

		if inputPip.PublicIPAddressPropertiesFormat.PublicIPAddressVersion != network.IPv4 {
			return errors.Wrapf(errors.InvalidInput, "Public IP address for IPv4 is only supported")
		}
	}

	if len(inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod) == 0 {
		inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod = network.Dynamic
	} else {
		if inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod != network.Dynamic && inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod != network.Static {
			return errors.Wrapf(errors.InvalidInput, "Invalid IP address allocation for public IP address")
		}

		if inputPip.PublicIPAddressPropertiesFormat.PublicIPAllocationMethod == network.Static {

			return errors.Wrapf(errors.InvalidInput, "Public IP address for dynamic allocation is only supported")
			// // Enable these when the static allocation is supported
			// // if *inputPip.PublicIPAddressPropertiesFormat.IPAddress == "" {
			// // 	return errors.Wrapf(errors.InvalidInput, "Public IP address must be specified for static allocation")
			// // }
			// // parsedIP := net.ParseIP(*inputPip.PublicIPAddressPropertiesFormat.IPAddress)
			// // if parsedIP == nil || parsedIP.To4() == nil {
			// // 	return errors.Wrapf(errors.InvalidInput, "Public IP address must be valid IPv4 address format")
			// // }
		}
	}

	return nil
}

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
	var err error

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

// GetWssdPublicIPAddress converts our internal representation of a PublicIPAddress (network.PublicIPAddress) to the cloud public IP address protobuf used by wssdcloudagent (wssdnetwork.PublicIPAddress)

/*func getWssdPublicIPAddress(networkPip *network.PublicIPAddress, location string) (wssdCloudPip *wssdcloudnetwork.PublicIPAddress, err error) {

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if networkPip.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name for public IP Address")
	}

	wssdCloudPip = &wssdcloudnetwork.PublicIPAddress{
		Name:         *networkPip.Name,
		LocationName: location,
		IpAddress:    *networkPip.IPAddress,
		//Allocation: networkPip.PublicIPAllocationMethod,
		DomainNameLabel: *networkPip.DNSSettings.DomainNameLabel,
		ReverseFqdn:     *networkPip.DNSSettings.ReverseFqdn,
	}

	if networkPip.PublicIPAllocationMethod != nil {
		wssdCloudPip.Tags = tags.MapToProto(networkPip.Tags)
	}

	if networkPip.PublicIPAddressVersion != nil {
		wssdCloudPip.Tags = tags.MapToProto(networkPip.Tags)
	}

	if networkPip.Tags != nil {
		wssdCloudPip.Tags = tags.MapToProto(networkPip.Tags)
	}

	return wssdCloudPip, nil
}

// GetPublicIPAddress converts the cloud public IP address protobuf returned from wssdcloudagent (wssdcloudnetwork.PublicIPAddress) to our internal representation of a public IP address (network.PublicIPAddress)
func getPublicIPAddress(wssdPip *wssdcloudnetwork.PublicIPAddress) (networkPip *network.PublicIPAddress, err error) {
	networkPip = &network.PublicIPAddress{
		Name:     &wssdPip.Name,
		Location: &wssdPip.LocationName,
		ID:       &wssdPip.Id,
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			Statuses: status.GetStatuses(wssdPip.GetStatus()),
		},
	}

	if wssdPip.Tags != nil {
		networkPip.Tags = tags.ProtoToMap(wssdPip.Tags)
	}

	return networkPip, nil
}
*/
