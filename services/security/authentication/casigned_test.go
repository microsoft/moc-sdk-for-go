// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 license.
package authentication

import (
	"context"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/rpc/cloudagent/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
	"time"
)

type MockClient struct {
	mock.Mock
}

var accessLog = ctrl.Log.WithName("cloudaccess")

func Test_CalculateTime(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * -10))
	after := now.Add(time.Duration(time.Second * 10))
	actual, renew := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*4) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*4), actual)
	}
	if renew != time.Duration(time.Millisecond*400) {
		t.Errorf("Wrong renewbackoff time returned Expected %s Actual %s", time.Duration(time.Millisecond*400), renew)
	}
}

func Test_CalculateTimeNegative(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * -20))
	after := now.Add(time.Duration(time.Second * -10))
	actual, renew := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*-13) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*-13), actual)
	}
	if renew != time.Duration(time.Millisecond*200) {
		t.Errorf("Wrong renewbackoff time returned Expected %s Actual %s", time.Duration(time.Millisecond*200), renew)
	}
}

func Test_CalculateTimeAfter(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * 10))
	after := now.Add(time.Duration(time.Second * 30))
	actual, renew := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*24) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*24), actual)
	}
	if renew != time.Duration(time.Millisecond*400) {
		t.Errorf("Wrong renewbackoff time returned Expected %s Actual %s", time.Duration(time.Millisecond*400), renew)
	}
}

func (m *MockClient) Login(ctx context.Context, request *security.AuthenticationRequest, opts ...grpc.CallOption) (*security.AuthenticationResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*security.AuthenticationResponse), args.Error(1)
}

func Test_LoginWithConfigWithLoggerFromContext(t *testing.T) {
	mockClient := new(MockClient)
	ctx := context.Background()
	ctx = log.IntoContext(ctx, accessLog) // Add the logger to the context
	group := "test-group"
	loginConfig := auth.LoginConfig{
		Name:        "test-identity",
		Token:       "dGVzdC10b2tlbg==",
		Certificate: "dGVzdC1jZXJ0aWZpY2F0ZQ==",
	}
	enableRenewRoutine := true

	// Mock the Login method
	mockClient.On("Login", ctx, mock.AnythingOfType("*security.AuthenticationRequest")).Return(&security.AuthenticationResponse{
		Token: "mock-client-cert",
	}, nil)

	c := &client{
		AuthenticationAgentClient: mockClient,
		cloudFQDN:                 "test-cloudFQDN",
	}

	// Act
	result, err := c.LoginWithConfig(ctx, group, loginConfig, enableRenewRoutine)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-client-cert", result.ClientCertificate)
	mockClient.AssertExpectations(t)
}

func Test_LoginWithConfigWithoutLoggerFromContext(t *testing.T) {
	mockClient := new(MockClient)
	ctx := context.Background()
	group := "test-group"
	loginConfig := auth.LoginConfig{
		Name:        "test-identity",
		Token:       "dGVzdC10b2tlbg==",
		Certificate: "dGVzdC1jZXJ0aWZpY2F0ZQ==",
	}
	enableRenewRoutine := true

	// Mock the Login method
	mockClient.On("Login", ctx, mock.AnythingOfType("*security.AuthenticationRequest")).Return(&security.AuthenticationResponse{
		Token: "mock-client-cert",
	}, nil)

	c := &client{
		AuthenticationAgentClient: mockClient,
		cloudFQDN:                 "test-cloudFQDN",
	}

	// Act
	result, err := c.LoginWithConfig(ctx, group, loginConfig, enableRenewRoutine)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-client-cert", result.ClientCertificate)
	mockClient.AssertExpectations(t)
}
