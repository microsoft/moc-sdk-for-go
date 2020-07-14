// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package security

import (
	"github.com/Azure/go-autorest/autorest"
	uuid "github.com/satori/go.uuid"
)

// Reference: github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault/models.go
// SkuName enumerates the values for sku name.
type SkuName string

const (
	// Premium ...
	Premium SkuName = "premium"
	// Standard ...
	Standard SkuName = "standard"
)

// Sku SKU details
type Sku struct {
	// Family - SKU family name
	Family *string `json:"family,omitempty"`
	// Name - SKU name to specify whether the key vault is a standard vault or a premium vault. Possible values include: 'Standard', 'Premium'
	Name SkuName `json:"name,omitempty"`
}

// KeyPermissions enumerates the values for key permissions.
type KeyPermissions string

const (
	// KeyPermissionsBackup ...
	KeyPermissionsBackup KeyPermissions = "backup"
	// KeyPermissionsCreate ...
	KeyPermissionsCreate KeyPermissions = "create"
	// KeyPermissionsDecrypt ...
	KeyPermissionsDecrypt KeyPermissions = "decrypt"
	// KeyPermissionsDelete ...
	KeyPermissionsDelete KeyPermissions = "delete"
	// KeyPermissionsEncrypt ...
	KeyPermissionsEncrypt KeyPermissions = "encrypt"
	// KeyPermissionsGet ...
	KeyPermissionsGet KeyPermissions = "get"
	// KeyPermissionsImport ...
	KeyPermissionsImport KeyPermissions = "import"
	// KeyPermissionsList ...
	KeyPermissionsList KeyPermissions = "list"
	// KeyPermissionsPurge ...
	KeyPermissionsPurge KeyPermissions = "purge"
	// KeyPermissionsRecover ...
	KeyPermissionsRecover KeyPermissions = "recover"
	// KeyPermissionsRestore ...
	KeyPermissionsRestore KeyPermissions = "restore"
	// KeyPermissionsSign ...
	KeyPermissionsSign KeyPermissions = "sign"
	// KeyPermissionsUnwrapKey ...
	KeyPermissionsUnwrapKey KeyPermissions = "unwrapKey"
	// KeyPermissionsUpdate ...
	KeyPermissionsUpdate KeyPermissions = "update"
	// KeyPermissionsVerify ...
	KeyPermissionsVerify KeyPermissions = "verify"
	// KeyPermissionsWrapKey ...
	KeyPermissionsWrapKey KeyPermissions = "wrapKey"
)

// SecretPermissions enumerates the values for secret permissions.
type SecretPermissions string

const (
	// SecretPermissionsBackup ...
	SecretPermissionsBackup SecretPermissions = "backup"
	// SecretPermissionsDelete ...
	SecretPermissionsDelete SecretPermissions = "delete"
	// SecretPermissionsGet ...
	SecretPermissionsGet SecretPermissions = "get"
	// SecretPermissionsList ...
	SecretPermissionsList SecretPermissions = "list"
	// SecretPermissionsPurge ...
	SecretPermissionsPurge SecretPermissions = "purge"
	// SecretPermissionsRecover ...
	SecretPermissionsRecover SecretPermissions = "recover"
	// SecretPermissionsRestore ...
	SecretPermissionsRestore SecretPermissions = "restore"
	// SecretPermissionsSet ...
	SecretPermissionsSet SecretPermissions = "set"
)

// CertificatePermissions enumerates the values for certificate permissions.
type CertificatePermissions string

const (
	// Backup ...
	Backup CertificatePermissions = "backup"
	// Create ...
	Create CertificatePermissions = "create"
	// Delete ...
	Delete CertificatePermissions = "delete"
	// Deleteissuers ...
	Deleteissuers CertificatePermissions = "deleteissuers"
	// Get ...
	Get CertificatePermissions = "get"
	// Getissuers ...
	Getissuers CertificatePermissions = "getissuers"
	// Import ...
	Import CertificatePermissions = "import"
	// List ...
	List CertificatePermissions = "list"
	// Listissuers ...
	Listissuers CertificatePermissions = "listissuers"
	// Managecontacts ...
	Managecontacts CertificatePermissions = "managecontacts"
	// Manageissuers ...
	Manageissuers CertificatePermissions = "manageissuers"
	// Purge ...
	Purge CertificatePermissions = "purge"
	// Recover ...
	Recover CertificatePermissions = "recover"
	// Restore ...
	Restore CertificatePermissions = "restore"
	// Setissuers ...
	Setissuers CertificatePermissions = "setissuers"
	// Update ...
	Update CertificatePermissions = "update"
)

// StoragePermissions enumerates the values for storage permissions.
type StoragePermissions string

const (
	// StoragePermissionsBackup ...
	StoragePermissionsBackup StoragePermissions = "backup"
	// StoragePermissionsDelete ...
	StoragePermissionsDelete StoragePermissions = "delete"
	// StoragePermissionsDeletesas ...
	StoragePermissionsDeletesas StoragePermissions = "deletesas"
	// StoragePermissionsGet ...
	StoragePermissionsGet StoragePermissions = "get"
	// StoragePermissionsGetsas ...
	StoragePermissionsGetsas StoragePermissions = "getsas"
	// StoragePermissionsList ...
	StoragePermissionsList StoragePermissions = "list"
	// StoragePermissionsListsas ...
	StoragePermissionsListsas StoragePermissions = "listsas"
	// StoragePermissionsPurge ...
	StoragePermissionsPurge StoragePermissions = "purge"
	// StoragePermissionsRecover ...
	StoragePermissionsRecover StoragePermissions = "recover"
	// StoragePermissionsRegeneratekey ...
	StoragePermissionsRegeneratekey StoragePermissions = "regeneratekey"
	// StoragePermissionsRestore ...
	StoragePermissionsRestore StoragePermissions = "restore"
	// StoragePermissionsSet ...
	StoragePermissionsSet StoragePermissions = "set"
	// StoragePermissionsSetsas ...
	StoragePermissionsSetsas StoragePermissions = "setsas"
	// StoragePermissionsUpdate ...
	StoragePermissionsUpdate StoragePermissions = "update"
)

// Permissions permissions the identity has for keys, secrets, certificates and storage.
type Permissions struct {
	// Keys - Permissions to keys
	Keys *[]KeyPermissions `json:"keys,omitempty"`
	// Secrets - Permissions to secrets
	Secrets *[]SecretPermissions `json:"secrets,omitempty"`
	// Certificates - Permissions to certificates
	Certificates *[]CertificatePermissions `json:"certificates,omitempty"`
	// Storage - Permissions to storage accounts
	Storage *[]StoragePermissions `json:"storage,omitempty"`
}

// AccessPolicyEntry an identity that have access to the key vault. All identities in the array must use
// the same tenant ID as the key vault's tenant ID.
type AccessPolicyEntry struct {
	// TenantID - The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.
	TenantID *uuid.UUID `json:"tenantId,omitempty"`
	// ObjectID - The object ID of a user, service principal or security group in the Azure Active Directory tenant for the vault. The object ID must be unique for the list of access policies.
	ObjectID *string `json:"objectId,omitempty"`
	// ApplicationID -  Application ID of the client making request on behalf of a principal
	ApplicationID *uuid.UUID `json:"applicationId,omitempty"`
	// Permissions - Permissions the identity has for keys, secrets and certificates.
	Permissions *Permissions `json:"permissions,omitempty"`
}

// KeyVaultProperties properties of the vault
type KeyVaultProperties struct {
	// TenantID - The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.
	TenantID *uuid.UUID `json:"tenantId,omitempty"`
	// Sku - SKU details
	Sku *Sku `json:"sku,omitempty"`
	// AccessPolicies - An array of 0 to 16 identities that have access to the key vault. All identities in the array must use the same tenant ID as the key vault's tenant ID. When `createMode` is set to `recover`, access policies are not required. Otherwise, access policies are required.
	AccessPolicies *[]AccessPolicyEntry `json:"accessPolicies,omitempty"`
	// VaultURI - The URI of the vault for performing operations on keys and secrets.
	VaultURI *string `json:"vaultUri,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// KeyVault resource information with extended details.
type KeyVault struct {
	autorest.Response `json:"-"`
	// KeyVaultProperties - Properties of the vault
	*KeyVaultProperties `json:"properties,omitempty"`
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
}

// CertificateAttributes the certificate management attributes
type CertificateAttributes struct {
	// Enabled - Determines whether the object is enabled
	Enabled *bool `json:"enabled,omitempty"`
	// NotBefore - Not before date in seconds since 1970-01-01T00:00:00Z
	NotBefore *int64 `json:"nbf,omitempty"`
	// Expires - Expiry date in seconds since 1970-01-01T00:00:00Z
	Expires *int64 `json:"exp,omitempty"`
	// Created - READ-ONLY; Creation time in seconds since 1970-01-01T00:00:00Z
	Created *int64 `json:"created,omitempty"`
	// Updated - READ-ONLY; Last updated time in seconds since 1970-01-01T00:00:00Z
	Updated *int64 `json:"updated,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Certificate a certificate consists of a certificate (X509) plus its attributes.
type Certificate struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The certificate id
	ID *string `json:"id,omitempty"`
	// Name - The certificate name
	Name *string `json:"name,omitempty"`
	// X509Thumbprint - READ-ONLY; Thumbprint of the certificate. (a URL-encoded base64 string)
	X509Thumbprint *string `json:"x5t,omitempty"`
	// Cer - CER contents of x509 certificate.
	Cer *[]byte `json:"cer,omitempty"`
	// Type - The content type of the certificate
	Type *string `json:"contentType,omitempty"`
	// Attributes - The certificate attributes.
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	// Tags - Application-specific metadata in the form of key-value pairs
	Tags map[string]*string `json:"tags"`
}

// IdentityProperties defines the structure of a Security Item
type IdentityProperties struct {
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Identity defines the structure of a identity
type Identity struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Token
	Token *string `json:"token,omitempty"`
	// Certificate
	Certificate *[]byte `json:"certificate,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*IdentityProperties `json:"properties,omitempty"`
}
