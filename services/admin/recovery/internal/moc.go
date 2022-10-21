// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	mocclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc/pkg/auth"
	mocadmin "github.com/microsoft/moc/rpc/common/admin"
)

type client struct {
	mocadmin.RecoveryAgentClient
}

// NewRecoveryClient - creates a client session with the backend moc agent
func NewRecoveryClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mocclient.GetRecoveryClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Backup
func (c *client) Backup(ctx context.Context, path string, configFilePath string, storeType string) error {
	request := getRecoveryRequest(mocadmin.Operation_BACKUP, path, configFilePath, storeType)
	_, err := c.RecoveryAgentClient.Invoke(ctx, request)
	return err
}

// Restore
func (c *client) Restore(ctx context.Context, path string, configFilePath string, storeType string) error {
	request := getRecoveryRequest(mocadmin.Operation_RESTORE, path, configFilePath, storeType)
	_, err := c.RecoveryAgentClient.Invoke(ctx, request)
	return err
}

func getRecoveryRequest(operation mocadmin.Operation, path string, configFilePath string, storeType string) *mocadmin.RecoveryRequest {
	return &mocadmin.RecoveryRequest{OperationType: operation, Path: path, ConfigFilePath: configFilePath, StoreType: storeType}
}
