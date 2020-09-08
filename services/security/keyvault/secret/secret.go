// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package secret

import (
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getSecret(sec *wssdcloudsecurity.Secret, vaultName string) *keyvault.Secret {
	value := string(sec.Value)
	return &keyvault.Secret{
		ID:      &sec.Id,
		Name:    &sec.Name,
		Value:   &value,
		Version: &sec.Status.Version.Number,
		SecretProperties: &keyvault.SecretProperties{
			FileName:  &sec.Filename,
			VaultName: &vaultName,
			Statuses:  status.GetStatuses(sec.GetStatus()),
		},
	}
}

func getWssdSecret(sec *keyvault.Secret, opType wssdcloudcommon.Operation) (*wssdcloudsecurity.Secret, error) {
	if sec.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Keyvault Secret name is missing")
	}
	if sec.VaultName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Keyvault name is missing")
	}
	secret := &wssdcloudsecurity.Secret{
		Name:      *sec.Name,
		VaultName: *sec.VaultName,
	}

	if sec.Version != nil {
		if secret.Status == nil {
			secret.Status = status.InitStatus()
		}
		secret.Status.Version.Number = *sec.Version
	}

	if sec.Value == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Secrets Value is empty")
	}
	if opType == wssdcloudcommon.Operation_POST {
		secret.Value = []byte(*sec.Value)
	}

	return secret, nil
}
