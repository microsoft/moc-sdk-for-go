// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"strings"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/cloud"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
	"github.com/microsoft/moc/rpc/common"
)

var (
	name                        = "test"
	Id                          = "1234"
	Version                     = "1234"
	Location                    = "mylocation"
	Fqdn                        = "fqdn"
	Port                  int32 = 1234
	AuthorizePort         int32 = 5678
	NodeAgentAuthMode           = cloud.NodeAgentCertificateAuth
	NodeAgentAuthModeWssd       = wssdcloud.NodeAgentAuthenticationMode_Certificate
	RunningState                = wssdcloud.NodeState_Active
	Info                        = common.NodeInfo{Name: name}
)

func Test_getWssdNode(t *testing.T) {
	grp := &cloud.Node{
		Name: &name,
		ID:   &Id,
		NodeProperties: &cloud.NodeProperties{
			FQDN:                        &Fqdn,
			Port:                        &Port,
			AuthorizerPort:              &AuthorizePort,
			NodeAgentAuthenticationMode: &NodeAgentAuthMode,
		},
	}
	wssdcloudNode, err := getWssdNode(grp, Location)
	if err != nil {
		t.Errorf("error while attempting to get wssdcloudControlPlane, %v", err)
	}

	if *grp.Name != wssdcloudNode.Name {
		t.Errorf("Name doesnt match post conversion")
	}
	if Location != wssdcloudNode.LocationName {
		t.Errorf("Location doesnt match")
	}
	if *grp.FQDN != wssdcloudNode.Fqdn {
		t.Errorf("FQDN doesnt match post conversion")
	}
	if *grp.Port != wssdcloudNode.Port {
		t.Errorf("Port doesnt match post conversion")
	}
	if *grp.AuthorizerPort != wssdcloudNode.AuthorizerPort {
		t.Errorf("AuthorizerPort doesnt match")
	}
	if NodeAgentAuthModeWssd != wssdcloudNode.NodeagentAuthenticationMode {
		t.Errorf("NodeAgentAuthenticationMode doesnt match")
	}
}

func Test_checkCloudNodeAgentAuthMode(t *testing.T) {
	tests := []struct {
		name                  string
		expectedNodeAuthAgent cloud.NodeAgentAuthenticationModeType
		actualNodeAuthAgent   wssdcloud.NodeAgentAuthenticationMode
		iscloudAgentAuthNil   bool
	}{
		{
			name:                  "with certificate auth selected",
			expectedNodeAuthAgent: cloud.NodeAgentCertificateAuth,
			actualNodeAuthAgent:   wssdcloud.NodeAgentAuthenticationMode_Certificate,
			iscloudAgentAuthNil:   false,
		},
		{
			name:                  "with poptoken auth selected",
			expectedNodeAuthAgent: cloud.NodeAgentPopTokenAuth,
			actualNodeAuthAgent:   wssdcloud.NodeAgentAuthenticationMode_PopToken,
			iscloudAgentAuthNil:   false,
		},
		{
			name:                  "with nil set, default to certificate auth",
			expectedNodeAuthAgent: cloud.NodeAgentCertificateAuth,
			actualNodeAuthAgent:   wssdcloud.NodeAgentAuthenticationMode_Certificate,
			iscloudAgentAuthNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var expectedNodeAuthAgent *cloud.NodeAgentAuthenticationModeType = nil
			if !tt.iscloudAgentAuthNil {
				expectedNodeAuthAgent = &tt.expectedNodeAuthAgent
			}

			grp := &cloud.Node{
				Name: &name,
				ID:   &Id,
				NodeProperties: &cloud.NodeProperties{
					FQDN:                        &Fqdn,
					Port:                        &Port,
					AuthorizerPort:              &AuthorizePort,
					NodeAgentAuthenticationMode: expectedNodeAuthAgent,
				},
			}
			wssdcloudNode, err := getWssdNode(grp, Location)
			if err != nil {
				t.Errorf("error while attempting to get wssdcloudNode, %v", err)
			}

			if tt.actualNodeAuthAgent != wssdcloudNode.NodeagentAuthenticationMode {
				t.Errorf("NodeAgentAuthenticationMode doesnt match")
			}
		})
	}

}

func Test_getNode(t *testing.T) {
	wssdcloudNode := &wssdcloud.Node{
		Name:                        name,
		Id:                          Id,
		LocationName:                Location,
		Fqdn:                        Fqdn,
		Port:                        Port,
		AuthorizerPort:              AuthorizePort,
		NodeagentAuthenticationMode: NodeAgentAuthModeWssd,
		RunningState:                RunningState,
		Info:                        &Info,
		Status:                      &common.Status{Version: &common.Version{Number: Version}},
	}
	grp := getNode(wssdcloudNode)
	if *grp.Location != wssdcloudNode.LocationName {
		t.Errorf("Location doesnt match post conversion")
	}
	if *grp.Name != wssdcloudNode.Name {
		t.Errorf("Name doesnt match post conversion")
	}
	if *grp.Version != Version {
		t.Errorf("Name doesnt match post conversion")
	}
	if *grp.NodeProperties.FQDN != wssdcloudNode.Fqdn {
		t.Errorf("Fqdn doesnt match post conversion")
	}
	if *grp.NodeProperties.Port != wssdcloudNode.Port {
		t.Errorf("Port doesnt match post conversion")
	}
	if *grp.NodeProperties.AuthorizerPort != wssdcloudNode.AuthorizerPort {
		t.Errorf("Authorizer port doesnt match post conversion")
	}
	if *grp.NodeProperties.NodeAgentAuthenticationMode != NodeAgentAuthMode {
		t.Errorf("Node agent auth mode doesnt match post conversion")
	}
	if *grp.NodeProperties.AuthorizerPort != wssdcloudNode.AuthorizerPort {
		t.Errorf("Authorizer port doesnt match post conversion")
	}
	if *grp.Statuses["RunningState"] != RunningState.String() {
		t.Errorf("Statues runningState doesnt match post conversion")
	}
	if !strings.Contains(*grp.Statuses["Info"], name) {
		t.Errorf("Statues info doesnt match post conversion")
	}

}

func Test_checkWssdNodeAgentAuthMode(t *testing.T) {
	tests := []struct {
		name                  string
		actualNodeAuthAgent   cloud.NodeAgentAuthenticationModeType
		expectedNodeAuthAgent wssdcloud.NodeAgentAuthenticationMode
	}{
		{
			name:                  "with certificate auth selected",
			actualNodeAuthAgent:   cloud.NodeAgentCertificateAuth,
			expectedNodeAuthAgent: wssdcloud.NodeAgentAuthenticationMode_Certificate,
		},
		{
			name:                  "with poptoken auth selected",
			actualNodeAuthAgent:   cloud.NodeAgentPopTokenAuth,
			expectedNodeAuthAgent: wssdcloud.NodeAgentAuthenticationMode_PopToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wssdcloudNode := &wssdcloud.Node{
				Name:                        name,
				Id:                          Id,
				LocationName:                Location,
				Fqdn:                        Fqdn,
				Port:                        Port,
				AuthorizerPort:              AuthorizePort,
				NodeagentAuthenticationMode: tt.expectedNodeAuthAgent,
				RunningState:                RunningState,
				Info:                        &Info,
				Status:                      &common.Status{Version: &common.Version{Number: Version}},
			}
			grp := getNode(wssdcloudNode)

			if tt.actualNodeAuthAgent != *grp.NodeAgentAuthenticationMode {
				t.Errorf("NodeAgentAuthenticationMode doesnt match")
			}
		})
	}

}
