// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package identity

import (
	"runtime"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
)

var (
	id                         = "1234"
	name                       = "test-identity"
	token                      = "TUIOWYTIUTI"
	tokenExpiry          int64 = 30
	tokenExpiryInSeconds int64 = 30
	revoked                    = false
	location                   = "TestLocation"
	version                    = "1.1.34"
	authType                   = auth.CASigned
	clientType                 = auth.ExternalClient
	cloudFqdn                  = ""
	cloudPort            int32 = 5000
	cloudAuthPort        int32 = 5000

	LoginFilePathEmpty           = ""
	LoginFilePathAbsoluteWindows = "C:\\AksHci\\ClusterStorage\\nodeToCloudLogin.yaml"
	LoginFilePathAbsoluteLinux   = "/home/usr/AksHci/ClusterStorage/nodeToCloudLogin.yaml"
	LoginFilepathRelative        = ".\\nodeToCloudLogin.yaml"

	expectedIdenityAutoRotateDisabled = security.Identity{
		ID:                   &id,
		Name:                 &name,
		Token:                &token,
		TokenExpiry:          &tokenExpiry,
		TokenExpiryInSeconds: &tokenExpiryInSeconds,
		Revoked:              revoked,
		Location:             &location,
		Version:              &version,
		AutoRotate:           false,
		LoginFilePath:        &LoginFilePathEmpty,
		IdentityProperties: &security.IdentityProperties{
			ClientType:    clientType,
			CloudFqdn:     &cloudFqdn,
			CloudPort:     &cloudPort,
			CloudAuthPort: &cloudAuthPort,
		},
	}
	expectedIdenityAutoRotateEnabledEmptyPath = security.Identity{
		ID:                   &id,
		Name:                 &name,
		Token:                &token,
		TokenExpiry:          &tokenExpiry,
		TokenExpiryInSeconds: &tokenExpiryInSeconds,
		Revoked:              revoked,
		Location:             &location,
		Version:              &version,
		AutoRotate:           true,
		LoginFilePath:        &LoginFilePathEmpty,
		IdentityProperties: &security.IdentityProperties{
			ClientType:    clientType,
			CloudFqdn:     &cloudFqdn,
			CloudPort:     &cloudPort,
			CloudAuthPort: &cloudAuthPort,
		},
	}

	expectedIdenityAutoRotateEnabledAbsolutePath = security.Identity{
		ID:                   &id,
		Name:                 &name,
		Token:                &token,
		TokenExpiry:          &tokenExpiry,
		TokenExpiryInSeconds: &tokenExpiryInSeconds,
		Revoked:              revoked,
		Location:             &location,
		Version:              &version,
		AutoRotate:           true,
		// LoginFilePath:        &LoginFilePathAbsoluteWindows,
		IdentityProperties: &security.IdentityProperties{
			ClientType:    clientType,
			CloudFqdn:     &cloudFqdn,
			CloudPort:     &cloudPort,
			CloudAuthPort: &cloudAuthPort,
		},
	}

	invalidNameIdentity = security.Identity{
		ID:                   &id,
		Token:                &token,
		TokenExpiry:          &tokenExpiry,
		TokenExpiryInSeconds: &tokenExpiryInSeconds,
		Revoked:              revoked,
		Location:             &location,
		Version:              &version,
		AutoRotate:           false,
		LoginFilePath:        &LoginFilePathEmpty,
		IdentityProperties: &security.IdentityProperties{
			ClientType:    clientType,
			CloudFqdn:     &cloudFqdn,
			CloudPort:     &cloudPort,
			CloudAuthPort: &cloudAuthPort,
		},
	}
	invalidLoginFilePathIdentity = security.Identity{
		ID:                   &id,
		Name:                 &name,
		Token:                &token,
		TokenExpiry:          &tokenExpiry,
		TokenExpiryInSeconds: &tokenExpiryInSeconds,
		Revoked:              revoked,
		Location:             &location,
		Version:              &version,
		AutoRotate:           true,
		LoginFilePath:        &LoginFilepathRelative,
		IdentityProperties: &security.IdentityProperties{
			ClientType:    clientType,
			CloudFqdn:     &cloudFqdn,
			CloudPort:     &cloudPort,
			CloudAuthPort: &cloudAuthPort,
		},
	}
)

func Test_getWssdIdeneityValid(t *testing.T) {
	var err error
	_, err = getWssdIdentity(&expectedIdenityAutoRotateDisabled)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = getWssdIdentity(&expectedIdenityAutoRotateEnabledEmptyPath)
	if err != nil {
		t.Errorf(err.Error())
	}
	if runtime.GOOS == "windows" {
		expectedIdenityAutoRotateEnabledAbsolutePath.LoginFilePath = &LoginFilePathAbsoluteWindows
	} else {
		expectedIdenityAutoRotateEnabledAbsolutePath.LoginFilePath = &LoginFilePathAbsoluteLinux
	}
	_, err = getWssdIdentity(&expectedIdenityAutoRotateEnabledAbsolutePath)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_getWssdIdentityInvalidNoName(t *testing.T) {
	var err error
	_, err = getWssdIdentity(&invalidNameIdentity)
	if err == nil {
		t.Errorf("ERROR: getWssdIdentity DID NOT throw Identity name is missing error")
	}
}

func Test_getWssdIdentityRelativelyPath(t *testing.T) {
	var err error
	_, err = getWssdIdentity(&invalidLoginFilePathIdentity)
	if err == nil {
		t.Errorf("ERROR: getWssdIdentity DID NOT throw Identity Loginfile must be absolute filepath error")
	}
}
