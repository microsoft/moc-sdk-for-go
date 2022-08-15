// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package keyvault

import (
	"encoding/json"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
)

type SecretProperties struct {
	// VaultName
	VaultName *string `json:"vaultname"`
	// FileName
	FileName *string `json:"filename"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Secret defines the structure of a secret
type Secret struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The Azure Resource Manager resource ID for the key vault.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the key vault.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type of the key vault.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The supported Azure location where the key vault should be created.
	Location *string `json:"location,omitempty"`
	// Tags - The tags that will be assigned to the key vault.
	Tags map[string]*string `json:"tags"`
	// Value
	Value *string `json:"value"`
	// Properties - Properties of the vault
	*SecretProperties `json:"properties,omitempty"`
}

// JSONWebKeyType enumerates the values for json web key type.
type JSONWebKeyType string

const (
	// EC Elliptic Curve.
	EC JSONWebKeyType = "EC"
	// ECHSM Elliptic Curve with a private key which is not exportable from the HSM.
	ECHSM JSONWebKeyType = "EC-HSM"
	// Oct Octet sequence (used to represent symmetric keys)
	Oct JSONWebKeyType = "oct"
	// RSA RSA (https://tools.ietf.org/html/rfc3447)
	RSA JSONWebKeyType = "RSA"
	// RSAHSM RSA with a private key which is not exportable from the HSM.
	RSAHSM JSONWebKeyType = "RSA-HSM"
	// AES Advanced Encrytion Standard.
	AES JSONWebKeyType = "AES"
)

// JSONWebKeyCurveName enumerates the values for json web key curve name.
type JSONWebKeyCurveName string

const (
	// P256 The NIST P-256 elliptic curve, AKA SECG curve SECP256R1.
	P256 JSONWebKeyCurveName = "P-256"
	// P256K The SECG SECP256K1 elliptic curve.
	P256K JSONWebKeyCurveName = "P-256K"
	// P384 The NIST P-384 elliptic curve, AKA SECG curve SECP384R1.
	P384 JSONWebKeyCurveName = "P-384"
	// P521 The NIST P-521 elliptic curve, AKA SECG curve SECP521R1.
	P521 JSONWebKeyCurveName = "P-521"
)

// KeyProperties properties of the key pair backing a certificate.
type KeyProperties struct {
	// Exportable - Indicates if the private key can be exported.
	Exportable *bool `json:"exportable,omitempty"`
	// KeyType - The type of key pair to be used for the certificate. Possible values include: 'EC', 'ECHSM', 'RSA', 'RSAHSM', 'Oct', 'AES'
	KeyType JSONWebKeyType         `json:"kty,omitempty"`
	KeyOps  *[]JSONWebKeyOperation `json:"key_ops,omitempty"`
	// KeySize - The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`
	// ReuseKey - Indicates if the same key pair will be used on certificate renewal.
	ReuseKey *bool `json:"reuse_key,omitempty"`
	// Curve - Elliptic curve name. For valid values, see JsonWebKeyCurveName. Possible values include: 'P256', 'P384', 'P521', 'P256K'
	Curve JSONWebKeyCurveName `json:"crv,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
	// KeyRotationFrequencyInSeconds - Configures key rotation frequency.
	KeyRotationFrequencyInSeconds *int64 `json:"keyRotationFrequencyInSeconds,omitempty"`
}

// KeyOperationResult the key operation result.
type KeyOperationResult struct {
	autorest.Response `json:"-"`
	// Kid - READ-ONLY; Key identifier
	Kid *string `json:"kid,omitempty"`
	// Result - READ-ONLY; a URL-encoded base64 string
	Result *string `json:"value,omitempty"`
}

// JSONWebKeyEncryptionAlgorithm enumerates the values for json web key encryption algorithm.
type JSONWebKeyEncryptionAlgorithm string

const (
	// RSA15 ...
	RSA15 JSONWebKeyEncryptionAlgorithm = "RSA1_5"
	// RSAOAEP ...
	RSAOAEP JSONWebKeyEncryptionAlgorithm = "RSA-OAEP"
	// RSAOAEP256 ...
	RSAOAEP256 JSONWebKeyEncryptionAlgorithm = "RSA-OAEP-256"
	// A256KW AES Key Wrap with 256 bit key-encryption key
	A256KW JSONWebKeyEncryptionAlgorithm = "A256KW"
)

// KeyOperationsParameters the key operations parameters.
type KeyOperationsParameters struct {
	// Algorithm - algorithm identifier. Possible values include: 'RSAOAEP', 'RSAOAEP256', 'RSA15', 'A256KW'
	Algorithm JSONWebKeyEncryptionAlgorithm `json:"alg,omitempty"`
	// Value - a URL-encoded base64 string
	Value *string `json:"value,omitempty"`
}

// KeyRestoreParameters the key restore parameters.
type KeyRestoreParameters struct {
	// KeyBundleBackup - The backup blob associated with a key bundle. (a URL-encoded base64 string)
	KeyBundleBackup *string `json:"value,omitempty"`
}

// JSONWebKeyOperation enumerates the values for json web key operation.
type JSONWebKeyOperation string

const (
	// Decrypt ...
	Decrypt JSONWebKeyOperation = "decrypt"
	// Encrypt ...
	Encrypt JSONWebKeyOperation = "encrypt"
	// Sign ...
	Sign JSONWebKeyOperation = "sign"
	// UnwrapKey ...
	UnwrapKey JSONWebKeyOperation = "unwrapKey"
	// Verify ...
	Verify JSONWebKeyOperation = "verify"
	// WrapKey ...
	WrapKey JSONWebKeyOperation = "wrapKey"
)

// DeletionRecoveryLevel enumerates the values for deletion recovery level.
type DeletionRecoveryLevel string

const (
	// Purgeable ...
	Purgeable DeletionRecoveryLevel = "Purgeable"
	// Recoverable ...
	Recoverable DeletionRecoveryLevel = "Recoverable"
	// RecoverableProtectedSubscription ...
	RecoverableProtectedSubscription DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"
	// RecoverablePurgeable ...
	RecoverablePurgeable DeletionRecoveryLevel = "Recoverable+Purgeable"
)

// KeyAttributes the attributes of a key managed by the key vault service.
type KeyAttributes struct {
	// RecoveryLevel - READ-ONLY; Reflects the deletion recovery level currently in effect for keys in the current vault. If it contains 'Purgeable' the key can be permanently deleted by a privileged user; otherwise, only the system can purge the key, at the end of the retention interval. Possible values include: 'Purgeable', 'RecoverablePurgeable', 'Recoverable', 'RecoverableProtectedSubscription'
	RecoveryLevel DeletionRecoveryLevel `json:"recoveryLevel,omitempty"`
	// Enabled - Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// NotBefore - Not before date in UTC.
	NotBefore *date.UnixTime `json:"nbf,omitempty"`
	// Expires - Expiry date in UTC.
	Expires *date.UnixTime `json:"exp,omitempty"`
	// Created - READ-ONLY; Creation time in UTC.
	Created *date.UnixTime `json:"created,omitempty"`
	// Updated - READ-ONLY; Last updated time in UTC.
	Updated *date.UnixTime `json:"updated,omitempty"`
}

// Key defines the structure of a secret
type Key struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The Azure Resource Manager resource ID for the key vault.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the key vault.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type of the key vault.
	Type *string `json:"type,omitempty"`
	// Version
	Version *string `json:"version,omitempty"`
	// Location - The supported Azure location where the key vault should be created.
	Location *string `json:"location,omitempty"`
	// Tags - The tags that will be assigned to the key vault.
	Tags map[string]*string `json:"tags"`
	// Value
	Value *string `json:"value"`
	// Properties - Properties of the vault
	*KeyProperties `json:"properties,omitempty"`
}

// KeyWrappingAlgorithm enumerates the values for private key wrapping info for key import export operations.
type KeyWrappingAlgorithm string

const (
	// CKM_RSA_AES_KEY_WRAP
	CKM_RSA_AES_KEY_WRAP KeyWrappingAlgorithm = "CKM_RSA_AES_KEY_WRAP"
	// NO_KEY_WRAP
	NO_KEY_WRAP KeyWrappingAlgorithm = "NO_KEY_WRAP"
)

// Defines private key wrapping infor for key import/export operations
type PrivateKeyWrappingInfo struct {
	// KeyName - READ-ONLY; Name of the wrapping key
	KeyName *string `json:"key-name,omitempty"`
	// PublicKey - Public key of the wrapping key
	PublicKey *string `json:"public-key,omitempty"`
	// KeyWrappingAlgorithm - Key wrapping algorithm
	KeyWrappingAlgorithm *KeyWrappingAlgorithm `json:"enc,omitempty"`
}

// Defines the parameters for key import/export operations.
// Key import/export operations expect Json string equivalent of this struct in Key.Value field
type KeyImportExportValue struct {
	// Version - READ-ONLY; Version of the KeyImportExportValue JSON.
	Version string `json:"version,omitempty"`
	// PublicKey
	PublicKey *string `json:"public-key,omitempty"`
	// PrivateKey
	PrivateKey *string `json:"private-key,omitempty"`
	// PrivateKeyWrappingInfo
	PrivateKeyWrappingInfo *PrivateKeyWrappingInfo `json:"private-key-wrapping-info,omitempty"`
}

// JSONWebKeySignatureAlgorithm enumerates the values for json web key signature algorithm.
//https://docs.microsoft.com/en-us/dotnet/api/microsoft.azure.keyvault.models.keysignparameters?view=azure-dotnet-legacy&viewFallbackFrom=azure-dotnet
type JSONWebKeySignatureAlgorithm string

const (
	ES256  JSONWebKeySignatureAlgorithm = "ES256"
	ES256K JSONWebKeySignatureAlgorithm = "ES256K"
	ES384  JSONWebKeySignatureAlgorithm = "ES384"
	ES512  JSONWebKeySignatureAlgorithm = "ES512"
	PS256  JSONWebKeySignatureAlgorithm = "PS256"
	PS384  JSONWebKeySignatureAlgorithm = "PS384"
	PS512  JSONWebKeySignatureAlgorithm = "PS512"
	RS256  JSONWebKeySignatureAlgorithm = "RS256"
	RS384  JSONWebKeySignatureAlgorithm = "RS384"
	RS512  JSONWebKeySignatureAlgorithm = "RS512"
	RSNULL JSONWebKeySignatureAlgorithm = "RSNULL"
)

// KeySignParameters the key signing parameters.
type KeySignParameters struct {
	// Algorithm - Signing algorithm identifier. Possible values include: 'PS256', 'RS384', etc (see JSONWebKeySignatureAlgorithm)
	Algorithm JSONWebKeySignatureAlgorithm `json:"alg,omitempty"`
	// Value - a URL-encoded base64 string
	Value *string `json:"value,omitempty"`
}

// KeyVerifyParameters the key signing parameters.
type KeyVerifyParameters struct {
	// Algorithm - algorithm identifier. Possible values include: 'PS256', 'RS384', etc (see JSONWebKeySignatureAlgorithm)
	Algorithm JSONWebKeySignatureAlgorithm `json:"alg,omitempty"`
	// Value - a URL-encoded base64 string of the digest used for signing
	Digest *string `json:"digest,omitempty"`

	// Value - a URL-encoded base64 string of the signature to be verifed
	Signature *string `json:"value,omitempty"`
}

// KeyVerifyResult the key verify operation result.
//https://docs.microsoft.com/en-us/dotnet/api/microsoft.azure.keyvault.models.keyverifyresult?view=azure-dotnet-legacy
type KeyVerifyResult struct {
	autorest.Response `json:"-"` //Sethbe what is this???do we need it
	Value             *bool      `json:"value,omitempty"`
}

func getKeyVault(vault *wssdcloudsecurity.KeyVault, group string) *security.KeyVault {
	return &security.KeyVault{
		ID:      &vault.Id,
		Name:    &vault.Name,
		Version: &vault.Status.Version.Number,
		//	Source : &vault.Source,
		KeyVaultProperties: &security.KeyVaultProperties{
			Statuses: status.GetStatuses(vault.GetStatus()),
		},
	}
}

func getWssdKeyVault(vault *security.KeyVault, group string) (*wssdcloudsecurity.KeyVault, error) {
	if vault.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Keyvault name is missing")
	}
	if len(group) == 0 {
		return nil, errors.Wrapf(errors.InvalidGroup, "Group not specified")
	}
	keyvault := &wssdcloudsecurity.KeyVault{
		GroupName: group,
		Name:      *vault.Name,
		//	Source: *vault.Source,
		Secrets: []*wssdcloudsecurity.Secret{},
	}

	if vault.Version != nil {
		if keyvault.Status == nil {
			keyvault.Status = status.InitStatus()
		}
		keyvault.Status.Version.Number = *vault.Version
	}

	return keyvault, nil
}

func GetKeyImportExportJsonValue(publicKey, privateKey, wrappingKeyName, wrappingPubKey *string, keyWrappingAlgo *KeyWrappingAlgorithm) (string, error) {
	keyImportExportValue := KeyImportExportValue{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		PrivateKeyWrappingInfo: &PrivateKeyWrappingInfo{
			KeyName:              wrappingKeyName,
			PublicKey:            wrappingPubKey,
			KeyWrappingAlgorithm: keyWrappingAlgo,
		},
	}
	jsonBytes, err := json.Marshal(keyImportExportValue)
	return string(jsonBytes), err
}
