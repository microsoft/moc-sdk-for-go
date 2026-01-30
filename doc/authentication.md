# Authentication and Authorization

This guide explains how to authenticate with the MOC service using the SDK. The MOC SDK supports multiple authentication methods through the `auth.Authorizer` interface.

## Authentication Overview

The MOC SDK uses the `auth.Authorizer` from the `github.com/microsoft/moc` package to handle authentication. All client constructors require an authorizer:

```go
import "github.com/microsoft/moc/pkg/auth"

// Every client requires an authorizer
client, err := virtualmachine.NewVirtualMachineClient(cloudFQDN, authorizer)
```

## Authentication Methods

### 1. Certificate-Based Authentication (Recommended)

Certificate-based authentication is the recommended method for production environments. It uses mutual TLS (mTLS) for secure communication.

#### Basic Certificate Authentication

```go
import "github.com/microsoft/moc/pkg/auth"

func createAuthorizer() (auth.Authorizer, error) {
    authorizer, err := auth.NewAuthorizerFromCertificate(
        "/path/to/client-cert.pem",     // Client certificate
        "/path/to/client-key.pem",      // Client private key
        "/path/to/ca-cert.pem",         // CA certificate
        "",                              // Password (empty if not encrypted)
    )
    if err != nil {
        return nil, err
    }
    return authorizer, nil
}
```

#### Certificate with Password

If your private key is encrypted:

```go
authorizer, err := auth.NewAuthorizerFromCertificate(
    "/path/to/client-cert.pem",
    "/path/to/client-key.pem",
    "/path/to/ca-cert.pem",
    "your-key-password",
)
```

#### Certificate Requirements

- **Client Certificate**: X.509 certificate signed by a trusted CA
- **Private Key**: RSA or ECDSA private key (can be password-protected)
- **CA Certificate**: Root CA certificate that signed the client certificate

**Generate Certificates:**

```bash
# Generate CA private key
openssl genrsa -out ca-key.pem 4096

# Generate CA certificate
openssl req -new -x509 -days 365 -key ca-key.pem -out ca-cert.pem

# Generate client private key
openssl genrsa -out client-key.pem 4096

# Generate client certificate signing request
openssl req -new -key client-key.pem -out client.csr

# Sign client certificate with CA
openssl x509 -req -days 365 -in client.csr -CA ca-cert.pem \
    -CAkey ca-key.pem -CAcreateserial -out client-cert.pem
```

### 2. Token-Based Authentication

Token-based authentication uses JWT tokens for authentication:

```go
func createTokenAuthorizer() (auth.Authorizer, error) {
    token := "your-jwt-token"
    
    authorizer, err := auth.NewAuthorizerFromToken(
        token,
        "/path/to/ca-cert.pem",
    )
    if err != nil {
        return nil, err
    }
    return authorizer, nil
}
```

### 3. Environment-Based Authentication

Load credentials from environment variables:

```go
func createAuthorizerFromEnv() (auth.Authorizer, error) {
    // Expects these environment variables:
    // - MOC_CLIENT_CERT: path to client certificate
    // - MOC_CLIENT_KEY: path to client key
    // - MOC_CA_CERT: path to CA certificate
    
    authorizer, err := auth.NewAuthorizerFromEnvironment()
    if err != nil {
        return nil, err
    }
    return authorizer, nil
}
```

Set the environment variables:

```bash
export MOC_CLIENT_CERT=/path/to/client-cert.pem
export MOC_CLIENT_KEY=/path/to/client-key.pem
export MOC_CA_CERT=/path/to/ca-cert.pem
```

### 4. Configuration File Authentication

Load credentials from a configuration file:

```go
import "github.com/spf13/viper"

func createAuthorizerFromConfig() (auth.Authorizer, error) {
    viper.SetConfigName("auth")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("/etc/moc/")
    viper.AddConfigPath(".")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    authorizer, err := auth.NewAuthorizerFromCertificate(
        viper.GetString("client_cert"),
        viper.GetString("client_key"),
        viper.GetString("ca_cert"),
        viper.GetString("key_password"),
    )
    return authorizer, err
}
```

**auth.yaml:**

```yaml
client_cert: /path/to/client-cert.pem
client_key: /path/to/client-key.pem
ca_cert: /path/to/ca-cert.pem
key_password: ""
```

## Debug Mode (Development Only)

For development and testing, you can disable TLS authentication:

```go
import "os"

// Set debug mode
os.Setenv("WSSD_DEBUG_MODE", "on")

// Or using viper
viper.Set("Debug", true)
```

**⚠️ Warning:** Never use debug mode in production! It disables all security checks.

## Authorization Scopes

The MOC service uses role-based access control (RBAC). Ensure your authenticated identity has the required permissions:

### Common Roles

- **Reader**: Read-only access to resources
- **Contributor**: Read and write access to resources
- **Owner**: Full access including permission management

### Checking Permissions

```go
import "github.com/microsoft/moc-sdk-for-go/services/security/roleassignment"

func checkPermissions(cloudFQDN string, authorizer auth.Authorizer) error {
    roleClient, err := roleassignment.NewRoleAssignmentClient(cloudFQDN, authorizer)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    assignments, err := roleClient.Get(ctx, "default", "")
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d role assignments\n", len(*assignments))
    return nil
}
```

## Server Endpoints

The SDK connects to two ports on the MOC server:

- **Port 55000**: Main service port (default)
- **Port 65000**: Authentication port

### Custom Ports

Specify custom ports in the server address:

```go
// With custom port
cloudFQDN := "moc-server.example.com:55001"

// Default port (55000) is used if not specified
cloudFQDN := "moc-server.example.com"
```

## TLS Configuration

### Custom TLS Settings

For advanced TLS configuration:

```go
import (
    "crypto/tls"
    "github.com/microsoft/moc/pkg/auth"
)

func createCustomAuthorizer() (auth.Authorizer, error) {
    // Load certificates
    cert, err := tls.LoadX509KeyPair(
        "/path/to/client-cert.pem",
        "/path/to/client-key.pem",
    )
    if err != nil {
        return nil, err
    }
    
    // Create authorizer with custom TLS config
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
        // Add more custom settings as needed
    }
    
    authorizer := auth.NewAuthorizerFromTLSConfig(tlsConfig)
    return authorizer, nil
}
```

### Certificate Validation

The SDK validates server certificates against the provided CA certificate. To skip validation (not recommended):

```go
tlsConfig := &tls.Config{
    Certificates:       []tls.Certificate{cert},
    InsecureSkipVerify: true, // DON'T DO THIS IN PRODUCTION
}
```

## Connection Security

### Keep-Alive Settings

The SDK uses keep-alive to maintain connections:

```go
// Default keep-alive settings (configured automatically):
// - Time: 1 minute
// - Timeout: 20 seconds
// - PermitWithoutStream: true
```

These settings are configured in the client and don't need manual configuration.

### Connection Caching

The SDK automatically caches and reuses connections. To clear the cache:

```go
import "github.com/microsoft/moc-sdk-for-go/pkg/client"

// Clear all cached connections
client.ClearConnectionCache()
```

## Best Practices

### 1. Secure Credential Storage

```go
// ✅ Good: Load from secure storage
authorizer, err := auth.NewAuthorizerFromEnvironment()

// ❌ Bad: Hardcode credentials
authorizer, err := auth.NewAuthorizerFromCertificate(
    "/tmp/cert.pem", // Don't hardcode paths
    "/tmp/key.pem",
    "/tmp/ca.pem",
    "", // Don't hardcode passwords - use secure storage
)
```

### 2. Rotate Certificates Regularly

Implement certificate rotation:

```go
func rotateCredentials(oldAuth auth.Authorizer) (auth.Authorizer, error) {
    // Load new certificates
    newAuth, err := auth.NewAuthorizerFromCertificate(
        "/path/to/new-cert.pem",
        "/path/to/new-key.pem",
        "/path/to/ca-cert.pem",
        "",
    )
    if err != nil {
        return oldAuth, err // Keep using old auth on error
    }
    
    // Clear connection cache to use new credentials
    client.ClearConnectionCache()
    
    return newAuth, nil
}
```

### 3. Handle Authentication Errors

```go
authorizer, err := createAuthorizer()
if err != nil {
    log.Printf("Authentication failed: %v", err)
    
    // Check specific error types
    if strings.Contains(err.Error(), "certificate") {
        log.Println("Certificate error - check certificate paths and validity")
    } else if strings.Contains(err.Error(), "permission denied") {
        log.Println("Permission error - check RBAC roles")
    }
    
    return err
}
```

### 4. Use Least Privilege

Request only the permissions your application needs:

```go
// If you only need read access, use read-only credentials
// Configure this with your MOC administrator
```

## Troubleshooting

### "certificate verify failed"

**Problem:** Server certificate validation failed.

**Solution:**
```bash
# Verify CA certificate is correct
openssl verify -CAfile ca-cert.pem client-cert.pem

# Check certificate expiration
openssl x509 -in client-cert.pem -noout -dates
```

### "tls: bad certificate"

**Problem:** Server rejected your client certificate.

**Solution:**
- Ensure client certificate is signed by the trusted CA
- Check certificate hasn't expired
- Verify private key matches certificate

### "connection refused"

**Problem:** Cannot connect to MOC server.

**Solution:**
```bash
# Test network connectivity
nc -zv moc-server.example.com 55000

# Check if debug mode is needed
export WSSD_DEBUG_MODE=on
```

### "permission denied"

**Problem:** Authenticated but not authorized.

**Solution:**
- Check RBAC role assignments
- Verify identity has required permissions
- Contact MOC administrator to grant access

## Next Steps

- [Client Usage Guide](client-usage.md) - Configure clients with authentication
- [Security Services](services/security.md) - Manage identities and permissions
- [Advanced Topics](advanced/grpc-communication.md) - Deep dive into security
