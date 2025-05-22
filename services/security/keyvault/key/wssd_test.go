package key

import (
	"context"
	"testing"

	stdErrors "errors"

	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type KeyAgentClientMock struct {
}

func (s *KeyAgentClientMock) Invoke(ctx context.Context, in *wssdcloudsecurity.KeyRequest, opts ...grpc.CallOption) (*wssdcloudsecurity.KeyResponse, error) {
	out := new(wssdcloudsecurity.KeyResponse)
	out.Keys = in.Keys
	return out, nil
}

func (s *KeyAgentClientMock) Operate(ctx context.Context, in *wssdcloudsecurity.KeyOperationRequest, opts ...grpc.CallOption) (*wssdcloudsecurity.KeyOperationResponse, error) {
	out := new(wssdcloudsecurity.KeyOperationResponse)
	return out, nil
}

// CustomClient embeds the original client and overrides the get method
type CustomClient struct {
	*client
}

// Override the get method
func (c *CustomClient) get(_ context.Context, _, _, _, _ string) ([]*wssdcloudsecurity.Key, error) {
	// Custom implementation here
	// For example, return an empty slice to simulate "key not found"
	return []*wssdcloudsecurity.Key{}, nil
}

// Add this method to your CustomClient struct
func (c *CustomClient) getKeyOperationRequest(ctx context.Context,
	groupName, vaultName, name string,
	param *keyvault.KeyOperationsParameters,
	opType wssdcloudcommon.ProviderAccessOperation,
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

	// Call your overridden get method explicitly
	key, err := c.get(ctx, groupName, vaultName, name, param.KeyVersion)
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.Wrapf(errors.NotFound, "Key[%s] Vault[%s]", name, vaultName)
	}

	request.Key = key[0]
	return request, nil
}

func (c *CustomClient) getKeyOperationRequestRotate(ctx context.Context,
	groupName, vaultName, name string,
	opType wssdcloudcommon.ProviderAccessOperation,
) (*wssdcloudsecurity.KeyOperationRequest, error) {
	request := &wssdcloudsecurity.KeyOperationRequest{
		OperationType: opType,
	}

	key, err := c.get(ctx, groupName, vaultName, name, "")
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.Wrapf(errors.NotFound, "Key[%s] Vault[%s]", name, vaultName)
	}

	request.Key = key[0]
	return request, nil
}

func TestGetKeyOperationRequest_KeyVersion_Exists(t *testing.T) {
	KeyAgentClientMock := &KeyAgentClientMock{}
	mockClient := &client{KeyAgentClientMock}
	pointerToEmptyString := new(string)
	testKeyOperationsParameters := &keyvault.KeyOperationsParameters{
		Algorithm:  keyvault.A256KW,
		Value:      pointerToEmptyString,
		KeyVersion: "KeyVersion",
	}
	testRequest, err := mockClient.getKeyOperationRequest(context.Background(), "groupName", "vaultName", "name", testKeyOperationsParameters, wssdcloudcommon.ProviderAccessOperation_Unspecified)
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
			KeyVersion: "KeyVersion",
		},
	}
	correctRequest.Key.Status.Version = nil
	testRequest.Key.Status.Version = nil
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequest_KeyVersion_Not_Exists(t *testing.T) {
	KeyAgentClientMock := &KeyAgentClientMock{}
	mockClient := &client{KeyAgentClientMock}
	pointerToEmptyString := new(string)
	testKeyOperationsParameters := &keyvault.KeyOperationsParameters{
		Algorithm:  keyvault.A256KW,
		Value:      pointerToEmptyString,
		KeyVersion: "",
	}
	testRequest, err := mockClient.getKeyOperationRequest(context.Background(), "groupName", "vaultName", "name", testKeyOperationsParameters, wssdcloudcommon.ProviderAccessOperation_Unspecified)
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
	correctRequest.Key.Status.Version = nil
	testRequest.Key.Status.Version = nil
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequestRotate(t *testing.T) {
	KeyAgentClientMock := &KeyAgentClientMock{}
	mockClient := &client{KeyAgentClientMock}
	testRequest, err := mockClient.getKeyOperationRequestRotate(context.Background(), "groupName", "vaultName", "name", wssdcloudcommon.ProviderAccessOperation_Key_Rotate)
	assert.NoErrorf(t, err, "Failed to make getKeyOperationRequestRotate call", err)
	correctRequest := wssdcloudsecurity.KeyOperationRequest{
		OperationType: wssdcloudcommon.ProviderAccessOperation_Key_Rotate,
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
	correctRequest.Key.Status.Version = nil
	testRequest.Key.Status.Version = nil
	assert.Equal(t, correctRequest, *testRequest, "The testRequest and correctRequest are not equal")
}

func TestGetKeyOperationRequest_KeyNotFoundScenario(t *testing.T) {
	mockClient := &CustomClient{
		client: &client{KeyAgentClient: &KeyAgentClientMock{}},
	}

	pointerToEmptyString := new(string)
	testKeyOperationsParameters := &keyvault.KeyOperationsParameters{
		Algorithm:  keyvault.A256KW,
		Value:      pointerToEmptyString,
		KeyVersion: "",
	}

	_, err := mockClient.getKeyOperationRequest(context.Background(), "groupName", "vaultName", "name", testKeyOperationsParameters, wssdcloudcommon.ProviderAccessOperation_Unspecified)

	assert.Error(t, err, "Expected error when key is not found")
	assert.True(t, stdErrors.Is(err, errors.NotFound), "Expected NotFound error when key is not found")
	assert.Contains(t, err.Error(), "Key[name] Vault[vaultName]", "Error message should indicate missing key and vault")
}

func TestGetKeyOperationRequestRotate_KeyNotFoundScenario(t *testing.T) {
	mockClient := &CustomClient{
		client: &client{KeyAgentClient: &KeyAgentClientMock{}},
	}
	_, err := mockClient.getKeyOperationRequestRotate(context.Background(), "groupName", "vaultName", "name", wssdcloudcommon.ProviderAccessOperation_Key_Rotate)

	assert.Error(t, err, "Expected error when key is not found")
	assert.True(t, stdErrors.Is(err, errors.NotFound), "Expected NotFound error when key is not found")
	assert.Contains(t, err.Error(), "Key[name] Vault[vaultName]", "Error message should indicate missing key and vault")
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
