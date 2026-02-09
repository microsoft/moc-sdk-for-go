# Admin Services

The Admin Services provide APIs for system administration, monitoring, and maintenance operations.

## Overview

The MOC SDK Admin Services include:

- **Version** - Query system version information
- **Health** - Monitor system health and status
- **Logging** - Configure and retrieve logs
- **Recovery** - Backup and recovery operations
- **Validation** - Validate system configuration

## Version Management

Query system and component version information.

### Getting Version Information

```go
import (
    "context"
    "github.com/microsoft/moc-sdk-for-go/services/admin/version"
)

func getVersion(versionClient *version.VersionClient) error {
    ctx := context.Background()
    
    ver, err := versionClient.Get(ctx)
    if err != nil {
        return err
    }
    
    fmt.Printf("System Version: %s\n", *ver.Version)
    fmt.Printf("Build Number: %s\n", *ver.BuildNumber)
    fmt.Printf("Commit: %s\n", *ver.Commit)
    
    return nil
}
```

## Health Monitoring

Monitor the health status of the MOC service and components.

### Checking System Health

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/health"

func checkHealth(healthClient *health.HealthClient) error {
    ctx := context.Background()
    
    health, err := healthClient.Get(ctx, "default", "")
    if err != nil {
        return err
    }
    
    for _, h := range *health {
        fmt.Printf("Component: %s\n", *h.Name)
        if h.Properties != nil {
            fmt.Printf("  Status: %s\n", h.Properties.Status)
            if h.Properties.Message != nil {
                fmt.Printf("  Message: %s\n", *h.Properties.Message)
            }
        }
    }
    
    return nil
}
```

## Logging

Configure and retrieve system logs.

### Configuring Logging

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/logging"

func configureLogging(loggingClient *logging.LoggingClient) error {
    ctx := context.Background()
    
    loggingSpec := &admin.Logging{
        Name:     stringPtr("system-logging"),
        Location: stringPtr("default"),
        Properties: &admin.LoggingProperties{
            Level:      stringPtr("INFO"),
            MaxSizeMB:  int32Ptr(1000),
            MaxBackups: int32Ptr(5),
        },
    }
    
    log, err := loggingClient.CreateOrUpdate(ctx, "default", "system-logging", loggingSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Configured logging: %s\n", *log.Name)
    return nil
}
```

### Retrieving Logs

```go
func getLogs(loggingClient *logging.LoggingClient) error {
    ctx := context.Background()
    
    logs, err := loggingClient.Get(ctx, "default", "system-logging")
    if err != nil {
        return err
    }
    
    log := &(*logs)[0]
    if log.Properties != nil && log.Properties.Content != nil {
        fmt.Printf("Logs:\n%s\n", *log.Properties.Content)
    }
    
    return nil
}
```

## Recovery Operations

Perform backup and recovery operations.

### Creating a Backup

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/recovery"

func createBackup(recoveryClient *recovery.RecoveryClient) error {
    ctx := context.Background()
    
    backupSpec := &admin.Recovery{
        Name:     stringPtr("daily-backup"),
        Location: stringPtr("default"),
        Properties: &admin.RecoveryProperties{
            BackupPath: stringPtr("/backups/daily"),
        },
    }
    
    backup, err := recoveryClient.CreateOrUpdate(ctx, "default", "daily-backup", backupSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created backup: %s\n", *backup.Name)
    return nil
}
```

### Restoring from Backup

```go
func restoreBackup(recoveryClient *recovery.RecoveryClient) error {
    ctx := context.Background()
    
    restoreSpec := &admin.Recovery{
        Name:     stringPtr("restore-operation"),
        Location: stringPtr("default"),
        Properties: &admin.RecoveryProperties{
            BackupPath:  stringPtr("/backups/daily/backup-20240130.tar"),
            RestorePath: stringPtr("/restore"),
        },
    }
    
    restore, err := recoveryClient.CreateOrUpdate(ctx, "default", "restore-operation", restoreSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Restore initiated: %s\n", *restore.Name)
    return nil
}
```

## Validation

Validate system configuration and resource placement.

### Validating Configuration

```go
import "github.com/microsoft/moc-sdk-for-go/services/admin/validation"

func validateConfiguration(validationClient *validation.ValidationClient) error {
    ctx := context.Background()
    
    validationSpec := &admin.Validation{
        Name:     stringPtr("system-validation"),
        Location: stringPtr("default"),
        Properties: &admin.ValidationProperties{
            ValidationType: stringPtr("configuration"),
        },
    }
    
    result, err := validationClient.CreateOrUpdate(ctx, "default", "system-validation", validationSpec)
    if err != nil {
        return err
    }
    
    if result.Properties != nil {
        if *result.Properties.Status == "Valid" {
            fmt.Println("Configuration is valid")
        } else {
            fmt.Printf("Validation failed: %s\n", *result.Properties.Message)
        }
    }
    
    return nil
}
```

## Administrative Patterns

### System Health Check

```go
func performHealthCheck(healthClient *health.HealthClient) error {
    ctx := context.Background()
    
    health, err := healthClient.Get(ctx, "default", "")
    if err != nil {
        return fmt.Errorf("health check failed: %w", err)
    }
    
    unhealthyComponents := []string{}
    for _, h := range *health {
        if h.Properties != nil && h.Properties.Status != "Healthy" {
            unhealthyComponents = append(unhealthyComponents, *h.Name)
        }
    }
    
    if len(unhealthyComponents) > 0 {
        return fmt.Errorf("unhealthy components: %v", unhealthyComponents)
    }
    
    fmt.Println("All components healthy")
    return nil
}
```

### Automated Backup

```go
func setupDailyBackup(recoveryClient *recovery.RecoveryClient) error {
    ctx := context.Background()
    
    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("backup-%s", timestamp)
    
    backupSpec := &admin.Recovery{
        Name:     stringPtr(backupName),
        Location: stringPtr("default"),
        Properties: &admin.RecoveryProperties{
            BackupPath: stringPtr(fmt.Sprintf("/backups/%s", timestamp)),
        },
    }
    
    backup, err := recoveryClient.CreateOrUpdate(ctx, "default", backupName, backupSpec)
    if err != nil {
        return fmt.Errorf("backup failed: %w", err)
    }
    
    fmt.Printf("Backup completed: %s\n", *backup.Name)
    return nil
}
```

### Log Rotation

```go
func rotateLogsIfNeeded(loggingClient *logging.LoggingClient) error {
    ctx := context.Background()
    
    logs, err := loggingClient.Get(ctx, "default", "system-logging")
    if err != nil {
        return err
    }
    
    log := &(*logs)[0]
    if log.Properties != nil && log.Properties.CurrentSizeMB != nil {
        if *log.Properties.CurrentSizeMB > 900 { // 90% of max size
            // Trigger rotation
            log.Properties.Rotate = boolPtr(true)
            _, err := loggingClient.CreateOrUpdate(ctx, "default", "system-logging", log)
            if err != nil {
                return err
            }
            fmt.Println("Log rotation triggered")
        }
    }
    
    return nil
}
```

## Best Practices

### 1. Regular Health Checks

```go
// Schedule periodic health checks
ticker := time.NewTicker(5 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    if err := performHealthCheck(healthClient); err != nil {
        log.Printf("Health check failed: %v", err)
        // Alert administrator
    }
}
```

### 2. Automated Backups

```go
// Daily backup at midnight
ticker := time.NewTicker(24 * time.Hour)
defer ticker.Stop()

for range ticker.C {
    if err := setupDailyBackup(recoveryClient); err != nil {
        log.Printf("Backup failed: %v", err)
        // Alert administrator
    }
}
```

### 3. Log Management

```go
// Configure appropriate log levels for production
loggingSpec.Properties.Level = stringPtr("WARN") // Production
// vs.
loggingSpec.Properties.Level = stringPtr("DEBUG") // Development
```

### 4. Version Tracking

```go
// Track version for compatibility
ver, _ := versionClient.Get(ctx)
minRequiredVersion := "1.0.0"

if *ver.Version < minRequiredVersion {
    log.Printf("Warning: Version %s is below minimum required %s", *ver.Version, minRequiredVersion)
}
```

## Next Steps

- [Cloud Services](cloud.md) - Infrastructure organization
- [Troubleshooting](../troubleshooting.md) - Common issues
- [CI/CD](../development/ci-cd.md) - Automation pipelines
