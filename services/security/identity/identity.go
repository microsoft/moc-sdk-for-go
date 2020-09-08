// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
)

func getIdentity(id *wssdcloudsecurity.Identity) *security.Identity {
	return &security.Identity{
		ID:    &id.Id,
		Name:  &id.Name,
		Token: &id.Token,
		IdentityProperties: &security.IdentityProperties{
			Statuses: status.GetStatuses(id.GetStatus()),
		},
	}
}

func getWssdIdentity(id *security.Identity) (*wssdcloudsecurity.Identity, error) {
	if id.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Identity name is missing")
	}
	return &wssdcloudsecurity.Identity{
		Name: *id.Name,
	}, nil
}
