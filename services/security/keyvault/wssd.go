// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package keyvault

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudsecurity.KeyVaultAgentClient
}

// NewKeyVaultClientN- creates a client session with the backend wssdcloud agent
func newKeyVaultClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetKeyVaultClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.KeyVault, error) {
	request, err := getKeyVaultRequest(wssdcloudcommon.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.KeyVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	return getKeyVaultsFromResponse(response, group), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.KeyVault) (*security.KeyVault, error) {
	request, err := getKeyVaultRequest(wssdcloudcommon.Operation_POST, group, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.KeyVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}

	vault := getKeyVaultsFromResponse(response, group)

	if len(*vault) == 0 {
		return nil, fmt.Errorf("[KeyVault][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*vault)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vault, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vault) == 0 {
		return fmt.Errorf("Keyvault [%s] not found", name)
	}

	request, err := getKeyVaultRequest(wssdcloudcommon.Operation_DELETE, group, name, &(*vault)[0])
	if err != nil {
		return err
	}
	_, err = c.KeyVaultAgentClient.Invoke(ctx, request)
	services.HandleGRPCError(err)
	return err
}

func getKeyVaultsFromResponse(response *wssdcloudsecurity.KeyVaultResponse, group string) *[]security.KeyVault {
	vaults := []security.KeyVault{}
	for _, keyvaults := range response.GetKeyVaults() {
		vaults = append(vaults, *(getKeyVault(keyvaults, group)))
	}

	return &vaults
}

func getKeyVaultRequest(opType wssdcloudcommon.Operation, group, name string, vault *security.KeyVault) (*wssdcloudsecurity.KeyVaultRequest, error) {
	request := &wssdcloudsecurity.KeyVaultRequest{
		OperationType: opType,
		KeyVaults:     []*wssdcloudsecurity.KeyVault{},
	}

	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}

	wssdkeyvault := &wssdcloudsecurity.KeyVault{
		Name:      name,
		GroupName: group,
	}

	var err error
	if vault != nil {
		wssdkeyvault, err = getWssdKeyVault(vault, group)
		if err != nil {
			return nil, err
		}
	}
	request.KeyVaults = append(request.KeyVaults, wssdkeyvault)
	return request, nil
}
