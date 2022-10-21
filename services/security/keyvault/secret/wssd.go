// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package secret

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

type client struct {
	wssdcloudsecurity.SecretAgentClient
}

// NewSecretClient - creates a client session with the backend wssdcloud agent
func newSecretClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetSecretClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name, vaultName string) (*[]keyvault.Secret, error) {
	request, err := getSecretRequest(wssdcloudcommon.Operation_GET, name, vaultName, group, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getSecretsFromResponse(response, vaultName), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *keyvault.Secret) (*keyvault.Secret, error) {
	err := c.validate(ctx, group, name, sg)
	if err != nil {
		return nil, err
	}
	request, err := getSecretRequest(wssdcloudcommon.Operation_POST, name, *sg.VaultName, group, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, errors.Wrapf(err, "Secrets Create failed")
	}

	sec := getSecretsFromResponse(response, *sg.VaultName)

	if len(*sec) == 0 {
		return nil, fmt.Errorf("[Secret][Create] Unexpected error: Creating a secret returned no result")
	}

	return &((*sec)[0]), err
}

func (c *client) validate(ctx context.Context, group, name string, sg *keyvault.Secret) (err error) {
	if sg == nil || sg.VaultName == nil || sg.Value == nil {
		return errors.Wrapf(errors.InvalidInput, "Invalid Configuration")
	}
	if len(*sg.VaultName) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Missing Vault Name")
	}

	if sg.Name == nil {
		sg.Name = &name
	}
	return nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name, vaultName string) error {
	secret, err := c.Get(ctx, group, name, vaultName)
	if err != nil {
		return err
	}
	if len(*secret) == 0 {
		return fmt.Errorf("Keysecret [%s] not found", name)
	}

	request, err := getSecretRequest(wssdcloudcommon.Operation_DELETE, name, vaultName, group, &(*secret)[0])
	if err != nil {
		return err
	}
	_, err = c.SecretAgentClient.Invoke(ctx, request)
	return err
}

func getSecretsFromResponse(response *wssdcloudsecurity.SecretResponse, vaultName string) *[]keyvault.Secret {
	Secrets := []keyvault.Secret{}
	for _, secrets := range response.GetSecrets() {
		Secrets = append(Secrets, *(getSecret(secrets, vaultName)))
	}

	return &Secrets
}

func getSecretRequest(opType wssdcloudcommon.Operation, name, vaultName, groupName string, sec *keyvault.Secret) (*wssdcloudsecurity.SecretRequest, error) {

	request := &wssdcloudsecurity.SecretRequest{
		OperationType: opType,
		Secrets:       []*wssdcloudsecurity.Secret{},
	}
	if sec != nil {
		secret, err := getWssdSecret(groupName, sec, opType)
		if err != nil {
			return nil, err
		}
		request.Secrets = append(request.Secrets, secret)
	} else if len(name) > 0 {
		request.Secrets = append(request.Secrets,
			&wssdcloudsecurity.Secret{
				Name:      name,
				VaultName: vaultName,
				GroupName: groupName,
			})
	} else {
		request.Secrets = append(request.Secrets,
			&wssdcloudsecurity.Secret{
				VaultName: vaultName,
				GroupName: groupName,
			})
	}

	return request, nil
}
