# Security Services

The Security Services provide APIs for identity management, certificate management, key vault operations, and role-based access control (RBAC).

## Overview

The MOC SDK Security Services include:

- **Identity** - Manage managed identities for authentication
- **Key Vault** - Store and manage secrets, keys, and certificates
- **Certificates** - Manage certificates for TLS/SSL
- **Roles** - Define RBAC roles and permissions
- **Role Assignments** - Assign roles to identities
- **Authentication** - Handle CA-signed tokens and authentication

## Identity Management

Managed identities provide automatic credential management for applications.

### Creating an Identity

```go
import (
    "context"
    "github.com/microsoft/moc-sdk-for-go/services/security"
    "github.com/microsoft/moc-sdk-for-go/services/security/identity"
)

func createIdentity(identityClient *identity.IdentityClient) error {
    ctx := context.Background()
    
    identitySpec := &security.Identity{
        Name:     stringPtr("app-identity"),
        Location: stringPtr("default"),
        Properties: &security.IdentityProperties{
            TokenExpiry: int64Ptr(3600), // 1 hour token expiry
        },
    }
    
    id, err := identityClient.CreateOrUpdate(ctx, "production", "app-identity", identitySpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created identity: %s\n", *id.Name)
    if id.Properties != nil && id.Properties.ClientID != nil {
        fmt.Printf("  Client ID: %s\n", *id.Properties.ClientID)
    }
    
    return nil
}
```

### Identity Operations

```go
// Get identity
identity, err := identityClient.Get(ctx, "production", "app-identity")

// List all identities
identities, err := identityClient.Get(ctx, "production", "")

// Delete identity
err := identityClient.Delete(ctx, "production", "app-identity")
```

### Using Identity with VMs

```go
vmSpec := &compute.VirtualMachine{
    Name:     stringPtr("app-vm"),
    Location: stringPtr("default"),
    Identity: &compute.Identity{
        Type: compute.ResourceIdentityTypeSystemAssignedUserAssigned,
        UserAssignedIdentities: map[string]*compute.UserAssignedIdentitiesValue{
            "/production/identities/app-identity": {},
        },
    },
    Properties: &compute.VirtualMachineProperties{
        // ... other properties
    },
}
```

## Key Vault

Key vaults securely store secrets, keys, and certificates.

### Creating a Key Vault

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

func createKeyVault(kvClient *keyvault.KeyVaultClient) error {
    ctx := context.Background()
    
    kvSpec := &security.KeyVault{
        Name:     stringPtr("prod-keyvault"),
        Location: stringPtr("default"),
        Properties: &security.KeyVaultProperties{
            EnableSoftDelete: boolPtr(true),
        },
    }
    
    kv, err := kvClient.CreateOrUpdate(ctx, "production", "prod-keyvault", kvSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created key vault: %s\n", *kv.Name)
    return nil
}
```

### Secrets Management

```go
import (
    "os"
    "github.com/microsoft/moc-sdk-for-go/services/security/keyvault/secret"
)

func createSecret(secretClient *secret.SecretClient) error {
    ctx := context.Background()
    
    secretSpec := &security.Secret{
        Name:     stringPtr("database-password"),
        Location: stringPtr("default"),
        Properties: &security.SecretProperties{
            Value: stringPtr(os.Getenv("DB_PASSWORD")), // Use environment variable
        },
    }
    
    sec, err := secretClient.CreateOrUpdate(ctx, "prod-keyvault", "database-password", secretSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created secret: %s\n", *sec.Name)
    return nil
}

// Security Note: Set DB_PASSWORD environment variable before running:
//   export DB_PASSWORD="your-secure-database-password"

func getSecret(secretClient *secret.SecretClient) error {
    ctx := context.Background()
    
    secrets, err := secretClient.Get(ctx, "prod-keyvault", "database-password")
    if err != nil {
        return err
    }
    
    secret := &(*secrets)[0]
    if secret.Properties != nil && secret.Properties.Value != nil {
        fmt.Printf("Secret value: %s\n", *secret.Properties.Value)
    }
    
    return nil
}
```

### Keys Management

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/keyvault/key"

func createKey(keyClient *key.KeyClient) error {
    ctx := context.Background()
    
    keySpec := &security.Key{
        Name:     stringPtr("encryption-key"),
        Location: stringPtr("default"),
        Properties: &security.KeyProperties{
            KeySize: int32Ptr(2048),
            KeyType: security.KeyTypeRSA,
        },
    }
    
    k, err := keyClient.CreateOrUpdate(ctx, "prod-keyvault", "encryption-key", keySpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created key: %s\n", *k.Name)
    return nil
}
```

## Certificate Management

Manage certificates for TLS/SSL and authentication.

### Creating a Certificate

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/certificate"

func createCertificate(certClient *certificate.CertificateClient) error {
    ctx := context.Background()
    
    certSpec := &security.Certificate{
        Name:     stringPtr("web-cert"),
        Location: stringPtr("default"),
        Properties: &security.CertificateProperties{
            // PEM-encoded certificate
            Certificate: stringPtr(`-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKK...
-----END CERTIFICATE-----`),
            // Optional: Private key for certificate
            PrivateKey: stringPtr(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w...
-----END PRIVATE KEY-----`),
        },
    }
    
    cert, err := certClient.CreateOrUpdate(ctx, "production", "web-cert", certSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created certificate: %s\n", *cert.Name)
    return nil
}
```

### Certificate Operations

```go
// Get certificate
cert, err := certClient.Get(ctx, "production", "web-cert")

// List all certificates
certs, err := certClient.Get(ctx, "production", "")

// Delete certificate
err := certClient.Delete(ctx, "production", "web-cert")
```

## Role-Based Access Control (RBAC)

Manage roles and permissions for resources.

### Defining a Role

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/role"

func createRole(roleClient *role.RoleClient) error {
    ctx := context.Background()
    
    roleSpec := &security.Role{
        Name:     stringPtr("vm-operator"),
        Location: stringPtr("default"),
        Properties: &security.RoleProperties{
            Permissions: &[]security.Permission{
                {
                    Actions: &[]string{
                        "Microsoft.Compute/virtualMachines/read",
                        "Microsoft.Compute/virtualMachines/start/action",
                        "Microsoft.Compute/virtualMachines/stop/action",
                    },
                    NotActions: &[]string{},
                },
            },
        },
    }
    
    role, err := roleClient.CreateOrUpdate(ctx, "default", "vm-operator", roleSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created role: %s\n", *role.Name)
    return nil
}
```

### Assigning Roles

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/roleassignment"

func assignRole(raClient *roleassignment.RoleAssignmentClient) error {
    ctx := context.Background()
    
    assignmentSpec := &security.RoleAssignment{
        Name:     stringPtr("app-role-assignment"),
        Location: stringPtr("default"),
        Properties: &security.RoleAssignmentProperties{
            RoleDefinitionID: stringPtr("/default/roles/vm-operator"),
            PrincipalID:      stringPtr("/production/identities/app-identity"),
            Scope:            stringPtr("/production/virtualma chines/*"),
        },
    }
    
    assignment, err := raClient.CreateOrUpdate(ctx, "production", "app-role-assignment", assignmentSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created role assignment: %s\n", *assignment.Name)
    return nil
}
```

### Role Assignment Operations

```go
// Get role assignment
assignment, err := raClient.Get(ctx, "production", "app-role-assignment")

// List all role assignments
assignments, err := raClient.Get(ctx, "production", "")

// Delete role assignment
err := raClient.Delete(ctx, "production", "app-role-assignment")
```

## Authentication

### CA-Signed Authentication

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/authentication"

func authenticateWithCA(authClient *authentication.AuthenticationClient) error {
    ctx := context.Background()
    
    // Get CA-signed token
    token, err := authClient.GetToken(ctx, "default", &security.TokenRequest{
        Identity: stringPtr("/production/identities/app-identity"),
    })
    if err != nil {
        return err
    }
    
    fmt.Printf("Token: %s\n", *token.Token)
    fmt.Printf("Expires: %s\n", token.ExpiresOn)
    
    return nil
}
```

## Security Best Practices

### 1. Use Managed Identities

```go
// ✅ Good: Use managed identity
vmSpec.Identity = &compute.Identity{
    Type: compute.ResourceIdentityTypeSystemAssigned,
}

// ❌ Bad: Hardcode credentials in VM
```

### 2. Store Secrets in Key Vault

```go
// ✅ Good: Store in key vault
secretSpec := &security.Secret{
    Name: stringPtr("api-key"),
    Properties: &security.SecretProperties{
        Value: stringPtr(apiKey),
    },
}

// ❌ Bad: Hardcode in application
const APIKey = "secret-key-123"
```

### 3. Rotate Secrets Regularly

```go
// Update secret value
secretSpec.Properties.Value = stringPtr(newSecretValue)
updatedSecret, err := secretClient.CreateOrUpdate(ctx, vaultName, secretName, secretSpec)
```

### 4. Use Least Privilege

```go
// ✅ Good: Specific permissions
Permissions: &[]security.Permission{
    {
        Actions: &[]string{
            "Microsoft.Compute/virtualMachines/read",
        },
    },
}

// ❌ Bad: Too broad permissions
Actions: &[]string{"*"}
```

### 5. Enable Soft Delete

```go
kvSpec := &security.KeyVault{
    Name:     stringPtr("prod-keyvault"),
    Location: stringPtr("default"),
    Properties: &security.KeyVaultProperties{
        EnableSoftDelete: boolPtr(true), // Enable soft delete
    },
}
```

## Security Patterns

### VM with Managed Identity and Key Vault

```go
// Create identity
identitySpec := &security.Identity{
    Name:     stringPtr("vm-identity"),
    Location: stringPtr("default"),
}
identity, _ := identityClient.CreateOrUpdate(ctx, "production", "vm-identity", identitySpec)

// Create key vault
kvSpec := &security.KeyVault{
    Name:     stringPtr("vm-keyvault"),
    Location: stringPtr("default"),
}
kv, _ := kvClient.CreateOrUpdate(ctx, "production", "vm-keyvault", kvSpec)

// Store secret
secretSpec := &security.Secret{
    Name: stringPtr("db-connection"),
    Properties: &security.SecretProperties{
        Value: stringPtr("Server=db.example.com;Database=prod;"),
    },
}
secret, _ := secretClient.CreateOrUpdate(ctx, "vm-keyvault", "db-connection", secretSpec)

// Create VM with identity
vmSpec := &compute.VirtualMachine{
    Name:     stringPtr("app-vm"),
    Location: stringPtr("default"),
    Identity: &compute.Identity{
        Type: compute.ResourceIdentityTypeSystemAssignedUserAssigned,
        UserAssignedIdentities: map[string]*compute.UserAssignedIdentitiesValue{
            *identity.ID: {},
        },
    },
    Properties: &compute.VirtualMachineProperties{
        // ... VM properties
    },
}
```

## Next Steps

- [Compute Services](compute.md) - Use identities with VMs
- [Network Services](network.md) - Network security
- [Cloud Services](cloud.md) - Resource organization
