// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package key

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services"
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudsecurity.KeyAgentClient
}

// NewKeyClient - creates a client session with the backend wssdcloud agent
func newKeyClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetKeyClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, vaultName, name string) (*[]keyvault.Key, error) {
	request, err := getKeyRequestByVaultName(wssdcloudcommon.Operation_GET, group, vaultName, name)
	if err != nil {
		return nil, err
	}
	response, err := c.KeyAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}
	return getKeysFromResponse(response, vaultName)
}

// get
func (c *client) get(ctx context.Context, group, vaultName, name string) ([]*wssdcloudsecurity.Key, error) {
	request, err := getKeyRequestByVaultName(wssdcloudcommon.Operation_GET, group, vaultName, name)
	if err != nil {
		return nil, err
	}
	response, err := c.KeyAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, err
	}

	return response.GetKeys(), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, vaultName, name string,
	param *keyvault.Key) (*keyvault.Key, error) {
	err := c.validate(ctx, group, vaultName, name, param)
	if err != nil {
		return nil, err
	}
	if param.KeySize == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Invalid KeySize - Missing")
	}
	request, err := getKeyRequest(wssdcloudcommon.Operation_POST, group, vaultName, name, param)
	if err != nil {
		return nil, err
	}
	response, err := c.KeyAgentClient.Invoke(ctx, request)
	if err != nil {
		services.HandleGRPCError(err)

		return nil, errors.Wrapf(err, "Keys Create failed")
	}

	sec, err := getKeysFromResponse(response, vaultName)
	if err != nil {
		return nil, err
	}

	if len(*sec) == 0 {
		return nil, fmt.Errorf("[Key][Create] Unexpected error: Creating a key returned no result")
	}
	return &((*sec)[0]), err
}

func (c *client) validate(ctx context.Context, group, vaultName, name string, param *keyvault.Key) (err error) {
	if param == nil {
		return errors.Wrapf(errors.InvalidInput, "Invalid Configuration")
	}

	if len(vaultName) == 0 {
		errors.Wrapf(errors.InvalidInput, "Keyvault name is missing")
	}
	if len(name) == 0 {
		errors.Wrapf(errors.InvalidInput, "Keyvault name is missing")
	}

	return nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name, vaultName string) error {
	key, err := c.Get(ctx, group, vaultName, name)
	if err != nil {
		return err
	}
	if len(*key) == 0 {
		return fmt.Errorf("Keykey [%s] not found", name)
	}

	request, err := getKeyRequest(wssdcloudcommon.Operation_DELETE, group, vaultName, name, &(*key)[0])
	if err != nil {
		return err
	}
	_, err = c.KeyAgentClient.Invoke(ctx, request)
	services.HandleGRPCError(err)
	return err
}

func (c *client) Encrypt(ctx context.Context, group, vaultName, name string, param *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	request, err := c.getKeyOperationRequest(ctx, group, vaultName, name, param, wssdcloudcommon.KeyOperation_ENCRYPT)
	if err != nil {
		return
	}
	response, err := c.KeyAgentClient.Operate(ctx, request)
	if err != nil {
		return
	}
	result, err = getDataFromResponse(response)
	return
}

func (c *client) Decrypt(ctx context.Context, group, vaultName, name string, param *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	request, err := c.getKeyOperationRequest(ctx, group, vaultName, name, param, wssdcloudcommon.KeyOperation_DECRYPT)
	if err != nil {
		return
	}
	response, err := c.KeyAgentClient.Operate(ctx, request)
	if err != nil {
		return
	}
	result, err = getDataFromResponse(response)
	return
}

func (c *client) WrapKey(ctx context.Context, group, vaultName, name string, param *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	request, err := c.getKeyOperationRequest(ctx, group, vaultName, name, param, wssdcloudcommon.KeyOperation_WRAPKEY)
	if err != nil {
		return
	}
	response, err := c.KeyAgentClient.Operate(ctx, request)
	if err != nil {
		return
	}
	result, err = getDataFromResponse(response)
	return
}

func (c *client) UnwrapKey(ctx context.Context, group, vaultName, name string, param *keyvault.KeyOperationsParameters) (result *keyvault.KeyOperationResult, err error) {
	request, err := c.getKeyOperationRequest(ctx, group, vaultName, name, param, wssdcloudcommon.KeyOperation_UNWRAPKEY)
	if err != nil {
		return
	}
	response, err := c.KeyAgentClient.Operate(ctx, request)
	if err != nil {
		return
	}
	result, err = getDataFromResponse(response)
	return
}

func getKeysFromResponse(response *wssdcloudsecurity.KeyResponse, vaultName string) (*[]keyvault.Key, error) {
	tmp := []keyvault.Key{}
	for _, keys := range response.GetKeys() {
		tmpKey, err1 := getKey(keys, vaultName)
		if err1 != nil {
			return nil, err1
		}
		tmp = append(tmp, tmpKey)
	}

	return &tmp, nil
}

func getKeyRequestByVaultName(opType wssdcloudcommon.Operation, groupName, vaultName, name string) (*wssdcloudsecurity.KeyRequest, error) {
	request := &wssdcloudsecurity.KeyRequest{
		OperationType: opType,
		Keys:          []*wssdcloudsecurity.Key{},
	}
	key, err := getWssdKeyByVaultName(name, groupName, vaultName, opType)
	if err != nil {
		return nil, err
	}
	request.Keys = append(request.Keys, key)
	return request, nil
}

func getKeyRequest(opType wssdcloudcommon.Operation, groupName, vaultName, name string, param *keyvault.Key) (*wssdcloudsecurity.KeyRequest, error) {
	request := &wssdcloudsecurity.KeyRequest{
		OperationType: opType,
		Keys:          []*wssdcloudsecurity.Key{},
	}
	key, err := getWssdKey(name, param, groupName, vaultName, opType)
	if err != nil {
		return nil, err
	}
	request.Keys = append(request.Keys, key)
	return request, nil
}

func getDataFromResponse(response *wssdcloudsecurity.KeyOperationResponse) (result *keyvault.KeyOperationResult, err error) {
	result = &keyvault.KeyOperationResult{
		Result: &response.Data,
	}
	return result, nil
}

func (c *client) getKeyOperationRequest(ctx context.Context,
	groupName, vaultName, name string,
	param *keyvault.KeyOperationsParameters,
	opType wssdcloudcommon.KeyOperation,
) (*wssdcloudsecurity.KeyOperationRequest, error) {

	if param == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing KeyOperationsParameters")
	}

	if param.Value == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Value to be operated on")
	}

	algo, err := getMOCAlgorithm(param.Algorithm)
	if err != nil {
		return nil, err
	}
	request := &wssdcloudsecurity.KeyOperationRequest{
		OperationType: opType,
		Data:          *param.Value,
		Algorithm:     algo,
	}

	key, err := c.get(ctx, groupName, vaultName, name)
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.Wrapf(errors.NotFound, "Key[%s] Vault[%s]", name, vaultName)
	}

	request.Key = key[0]
	return request, nil
}
