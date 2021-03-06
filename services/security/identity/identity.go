// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getIdentity(id *wssdcloudsecurity.Identity) *security.Identity {
	clitype := security.ExternalClient
	if id.ClientType == wssdcloudcommon.ClientType_CONTROLPLANE {
		clitype = security.ControlPlane
	} else if id.ClientType == wssdcloudcommon.ClientType_NODE {
		clitype = security.Node
	}

	return &security.Identity{
		ID:          &id.Id,
		Name:        &id.Name,
		Token:       &id.Token,
		TokenExpiry: &id.TokenExpiry,
		Location:    &id.LocationName,
		Version:     &id.Status.Version.Number,
		AuthType:    auth.AuthTypeToLoginType(id.AuthType),
		IdentityProperties: &security.IdentityProperties{
			Statuses:      status.GetStatuses(id.GetStatus()),
			ClientType:    clitype,
			CloudFqdn:     &id.CloudFqdn,
			CloudPort:     &id.CloudPort,
			CloudAuthPort: &id.CloudAuthPort,
		},
	}
}

func getWssdIdentity(id *security.Identity) (*wssdcloudsecurity.Identity, error) {
	if id.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Identity name is missing")
	}

	wssdidentity := &wssdcloudsecurity.Identity{
		Name: *id.Name,
	}

	if id.TokenExpiry != nil {
		wssdidentity.TokenExpiry = *id.TokenExpiry
	}

	if id.Location != nil { // WIll need to do error checking if location not set !!!!s
		wssdidentity.LocationName = *id.Location
	}

	if id.Version != nil {
		if wssdidentity.Status == nil {
			wssdidentity.Status = status.InitStatus()
		}
		wssdidentity.Status.Version.Number = *id.Version
	}

	clitype := wssdcloudcommon.ClientType_EXTERNALCLIENT
	if id.IdentityProperties != nil {
		if id.IdentityProperties.ClientType == security.ControlPlane {
			clitype = wssdcloudcommon.ClientType_CONTROLPLANE
		} else if id.IdentityProperties.ClientType == security.Node {
			clitype = wssdcloudcommon.ClientType_NODE
		}

		if id.IdentityProperties.CloudFqdn != nil {
			wssdidentity.CloudFqdn = *id.CloudFqdn
		}

		if id.IdentityProperties.CloudPort != nil {
			wssdidentity.CloudPort = *id.CloudPort
		}

		if id.IdentityProperties.CloudAuthPort != nil {
			wssdidentity.CloudAuthPort = *id.CloudAuthPort
		}
	}

	wssdidentity.ClientType = clitype

	return wssdidentity, nil
}
