package key

import (
	"context"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
)

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

func TestGetKeyOperationRequest_Updated(t *testing.T) {
	mockClient := &client{nil}

	tests := []struct {
		name          string
		groupName     string
		vaultName     string
		keyName       string
		param         *keyvault.KeyOperationsParameters
		opType        wssdcloudcommon.ProviderAccessOperation
		expectError   bool
		errorContains string
	}{
		{
			name:      "Valid request with all parameters",
			groupName: "testGroup",
			vaultName: "testVault",
			keyName:   "testKey",
			param: &keyvault.KeyOperationsParameters{
				Algorithm:  keyvault.A256CBC,
				Value:      func() *string { s := "testValue"; return &s }(),
				KeyVersion: "123",
			},
			opType:      wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError: false,
		},
		{
			name:      "Valid request without key version",
			groupName: "testGroup",
			vaultName: "testVault",
			keyName:   "testKey",
			param: &keyvault.KeyOperationsParameters{
				Algorithm:  keyvault.A256CBC,
				Value:      func() *string { s := "testValue"; return &s }(),
				KeyVersion: "",
			},
			opType:      wssdcloudcommon.ProviderAccessOperation_Key_Decrypt,
			expectError: false,
		},
		{
			name:          "Empty vault name",
			groupName:     "testGroup",
			vaultName:     "",
			keyName:       "testKey",
			param:         &keyvault.KeyOperationsParameters{Algorithm: keyvault.A256CBC, Value: func() *string { s := "testValue"; return &s }()},
			opType:        wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError:   true,
			errorContains: "Keyvault name is missing",
		},
		{
			name:          "Empty key name",
			groupName:     "testGroup",
			vaultName:     "testVault",
			keyName:       "",
			param:         &keyvault.KeyOperationsParameters{Algorithm: keyvault.A256CBC, Value: func() *string { s := "testValue"; return &s }()},
			opType:        wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError:   true,
			errorContains: "Key name is missing",
		},
		{
			name:          "Nil parameters",
			groupName:     "testGroup",
			vaultName:     "testVault",
			keyName:       "testKey",
			param:         nil,
			opType:        wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError:   true,
			errorContains: "Missing KeyOperationsParameters",
		},
		{
			name:      "Nil value in parameters",
			groupName: "testGroup",
			vaultName: "testVault",
			keyName:   "testKey",
			param: &keyvault.KeyOperationsParameters{
				Algorithm: keyvault.A256CBC,
				Value:     nil,
			},
			opType:        wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError:   true,
			errorContains: "Missing Value to be operated on",
		},
		{
			name:      "Invalid algorithm",
			groupName: "testGroup",
			vaultName: "testVault",
			keyName:   "testKey",
			param: &keyvault.KeyOperationsParameters{
				Algorithm: "INVALID_ALGORITHM",
				Value:     func() *string { s := "testValue"; return &s }(),
			},
			opType:        wssdcloudcommon.ProviderAccessOperation_Key_Encrypt,
			expectError:   true,
			errorContains: "Invalid Algorithm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := mockClient.getKeyOperationRequest(
				context.Background(),
				tt.groupName,
				tt.vaultName,
				tt.keyName,
				tt.param,
				tt.opType,
			)

			if tt.expectError {
				assert.Error(t, err, "Expected error for test case: %s", tt.name)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				}
				assert.Nil(t, request, "Request should be nil when error occurs")
			} else {
				assert.NoError(t, err, "Expected no error for test case: %s", tt.name)
				assert.NotNil(t, request, "Request should not be nil for valid inputs")

				// Validate the request structure
				assert.Equal(t, tt.opType, request.OperationType, "Operation type should match")
				assert.Equal(t, *tt.param.Value, request.Data, "Data should match param value")

				// Validate the key structure
				assert.NotNil(t, request.Key, "Key should not be nil")
				assert.Equal(t, tt.keyName, request.Key.Name, "Key name should match")
				assert.Equal(t, tt.vaultName, request.Key.VaultName, "Vault name should match")
				assert.Equal(t, tt.groupName, request.Key.GroupName, "Group name should match")
				assert.Equal(t, tt.param.KeyVersion, request.Key.KeyVersion, "Key version should match")

				// Validate algorithm conversion
				expectedAlgo, algoErr := getMOCAlgorithm(tt.param.Algorithm)
				assert.NoError(t, algoErr, "Algorithm conversion should not fail")
				assert.Equal(t, expectedAlgo, request.Algorithm, "Algorithm should be correctly converted")
			}
		})
	}
}

func TestGetKeyOperationRequestRotate(t *testing.T) {
	mockClient := &client{nil}

	tests := []struct {
		name          string
		groupName     string
		vaultName     string
		keyName       string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid request with all parameters",
			groupName:   "testGroup",
			vaultName:   "testVault",
			keyName:     "testKey",
			expectError: false,
		},
		{
			name:        "Valid request with different names",
			groupName:   "anotherGroup",
			vaultName:   "anotherVault",
			keyName:     "anotherKey",
			expectError: false,
		},
		{
			name:          "Empty vault name",
			groupName:     "testGroup",
			vaultName:     "",
			keyName:       "testKey",
			expectError:   true,
			errorContains: "Keyvault name is missing",
		},
		{
			name:          "Empty key name",
			groupName:     "testGroup",
			vaultName:     "testVault",
			keyName:       "",
			expectError:   true,
			errorContains: "Key name is missing",
		},
		{
			name:          "Both vault name and key name empty",
			groupName:     "testGroup",
			vaultName:     "",
			keyName:       "",
			expectError:   true,
			errorContains: "Keyvault name is missing", // First validation error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := mockClient.getKeyOperationRequestRotate(
				context.Background(),
				tt.groupName,
				tt.vaultName,
				tt.keyName,
			)

			if tt.expectError {
				assert.Error(t, err, "Expected error for test case: %s", tt.name)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				}
				assert.Nil(t, request, "Request should be nil when error occurs")
			} else {
				assert.NoError(t, err, "Expected no error for test case: %s", tt.name)
				assert.NotNil(t, request, "Request should not be nil for valid inputs")

				// Validate the request structure
				assert.Equal(t, wssdcloudcommon.ProviderAccessOperation_Key_Rotate, request.OperationType, "Operation type should be Key_Rotate")

				// Validate the key structure
				assert.NotNil(t, request.Key, "Key should not be nil")
				assert.Equal(t, tt.keyName, request.Key.Name, "Key name should match")
				assert.Equal(t, tt.vaultName, request.Key.VaultName, "Vault name should match")
				assert.Equal(t, tt.groupName, request.Key.GroupName, "Group name should match")
				assert.Equal(t, "", request.Key.KeyVersion, "Key version should be empty for rotate operation")

				// Validate default key properties from getWssdKeyByVaultName
				assert.Equal(t, wssdcloudcommon.JsonWebKeyType_EC, request.Key.Type, "Key type should be default EC")
				assert.Equal(t, wssdcloudcommon.KeySize_K_UNKNOWN, request.Key.Size, "Key size should be default K_UNKNOWN")
				assert.NotNil(t, request.Key.KeyOps, "KeyOps should not be nil")
				assert.Equal(t, 0, len(request.Key.KeyOps), "KeyOps should be empty slice")
				assert.NotNil(t, request.Key.Status, "Status should not be nil")
			}
		})
	}
}
