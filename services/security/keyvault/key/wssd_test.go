package key

import (
	"context"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type testKeyAgentClient struct {
}

func (s *testKeyAgentClient) Invoke(ctx context.Context, in *wssdcloudsecurity.KeyRequest, opts ...grpc.CallOption) (*wssdcloudsecurity.KeyResponse, error) {
	out := new(wssdcloudsecurity.KeyResponse)
	out.Keys = in.Keys
	return out, nil
}

func (s *testKeyAgentClient) Operate(ctx context.Context, in *wssdcloudsecurity.KeyOperationRequest, opts ...grpc.CallOption) (*wssdcloudsecurity.KeyOperationResponse, error) {
	out := new(wssdcloudsecurity.KeyOperationResponse)
	return out, nil
}

func TestGetKeyOperationRequest_keyID_Exists(t *testing.T) {
	testKeyAgentClient := &testKeyAgentClient{}
	mockClient := &client{testKeyAgentClient}
	pointerToEmptyString := new(string)
	testKeyOperationsParameters := &keyvault.KeyOperationsParameters{
		Algorithm: keyvault.A256KW,
		Value:     pointerToEmptyString,
	}
	testRequest, err := mockClient.getKeyOperationRequest(context.Background(), "groupName", "vaultName", "name", "keyID", testKeyOperationsParameters, wssdcloudcommon.ProviderAccessOperation_Unspecified)
	assert.NoErrorf(t, err, "Failed to make getKeyOperationRequest call", err)
	algo, err := getMOCAlgorithm(testKeyOperationsParameters.Algorithm)
	assert.NoErrorf(t, err, "Failed to make getMOCAlgorithm call", err)
	correctRequest := wssdcloudsecurity.KeyOperationRequest{
		OperationType: wssdcloudcommon.ProviderAccessOperation_Unspecified,
		Data:          "",
		Algorithm:     algo,
		Key: &wssdcloudsecurity.Key{
			Name:       "name",
			VaultName:  "vaultName",
			GroupName:  "groupName",
			Type:       wssdcloudcommon.JsonWebKeyType_EC,
			Size:       wssdcloudcommon.KeySize_K_UNKNOWN,
			KeyOps:     []wssdcloudcommon.KeyOperation{},
			Status:     status.InitStatus(),
			KeyVersion: "keyID",
		},
	}
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequest_keyID_Not_Exists(t *testing.T) {
	testKeyAgentClient := &testKeyAgentClient{}
	mockClient := &client{testKeyAgentClient}
	pointerToEmptyString := new(string)
	testKeyOperationsParameters := &keyvault.KeyOperationsParameters{
		Algorithm: keyvault.A256KW,
		Value:     pointerToEmptyString,
	}
	request, err := mockClient.getKeyOperationRequest(context.Background(), "groupName", "vaultName", "name", "", testKeyOperationsParameters, wssdcloudcommon.ProviderAccessOperation_Unspecified)
	assert.NoErrorf(t, err, "Failed to make getKeyOperationRequest call", err)
	algo, err := getMOCAlgorithm(testKeyOperationsParameters.Algorithm)
	assert.NoErrorf(t, err, "Failed to make getMOCAlgorithm call", err)
	correctRequest := wssdcloudsecurity.KeyOperationRequest{
		OperationType: wssdcloudcommon.ProviderAccessOperation_Unspecified,
		Data:          "",
		Algorithm:     algo,
		Key: &wssdcloudsecurity.Key{
			Name:       "name",
			VaultName:  "vaultName",
			GroupName:  "groupName",
			Type:       wssdcloudcommon.JsonWebKeyType_EC,
			Size:       wssdcloudcommon.KeySize_K_UNKNOWN,
			KeyOps:     []wssdcloudcommon.KeyOperation{},
			Status:     status.InitStatus(),
			KeyVersion: "",
		},
	}
	assert.Equal(t, correctRequest, *request, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequestRotate_keyID_Exists(t *testing.T) {
	testKeyAgentClient := &testKeyAgentClient{}
	mockClient := &client{testKeyAgentClient}
	testRequest, err := mockClient.getKeyOperationRequestRotate(context.Background(), "groupName", "vaultName", "name", "keyID", wssdcloudcommon.ProviderAccessOperation_Unspecified)
	assert.NoErrorf(t, err, "Failed to make getKeyOperationRequestRotate call", err)
	correctRequest := wssdcloudsecurity.KeyOperationRequest{
		OperationType: wssdcloudcommon.ProviderAccessOperation_Unspecified,
		Key: &wssdcloudsecurity.Key{
			Name:       "name",
			VaultName:  "vaultName",
			GroupName:  "groupName",
			Type:       wssdcloudcommon.JsonWebKeyType_EC,
			Size:       wssdcloudcommon.KeySize_K_UNKNOWN,
			KeyOps:     []wssdcloudcommon.KeyOperation{},
			Status:     status.InitStatus(),
			KeyVersion: "keyID",
		},
	}
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequestRotate_keyID_Not_Exists(t *testing.T) {
	testKeyAgentClient := &testKeyAgentClient{}
	mockClient := &client{testKeyAgentClient}
	testRequest, err := mockClient.getKeyOperationRequestRotate(context.Background(), "groupName", "vaultName", "name", "", wssdcloudcommon.ProviderAccessOperation_Unspecified)
	assert.NoErrorf(t, err, "Failed to make getKeyOperationRequestRotate call", err)
	correctRequest := wssdcloudsecurity.KeyOperationRequest{
		OperationType: wssdcloudcommon.ProviderAccessOperation_Unspecified,
		Key: &wssdcloudsecurity.Key{
			Name:       "name",
			VaultName:  "vaultName",
			GroupName:  "groupName",
			Type:       wssdcloudcommon.JsonWebKeyType_EC,
			Size:       wssdcloudcommon.KeySize_K_UNKNOWN,
			KeyOps:     []wssdcloudcommon.KeyOperation{},
			Status:     status.InitStatus(),
			KeyVersion: "",
		},
	}
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestEncryptValidation_invalidAlgorithm(t *testing.T) {
	mockClient := &client{nil}
	err := mockClient.isSupportedEncryptionAlgorithm(keyvault.A256KW)

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestEncryptValidation_validAlgorithm(t *testing.T) {
	mockClient := &client{nil}
	err := mockClient.isSupportedEncryptionAlgorithm(keyvault.A256CBC)

	if err != nil {
		t.Errorf("Unexpected error  %+v", err)
	}
}

func TestWrapValidation_invalidAlgorithm(t *testing.T) {
	mockClient := &client{nil}
	err := mockClient.isSupportedWrapAlgorithm(keyvault.A256CBC)

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestWrapValidation_validAlgorithm(t *testing.T) {
	mockClient := &client{nil}
	err := mockClient.isSupportedWrapAlgorithm(keyvault.A256KW)

	if err != nil {
		t.Errorf("Unexpected error  %+v", err)
	}
}
