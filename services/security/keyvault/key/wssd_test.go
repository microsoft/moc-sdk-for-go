package key

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
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
