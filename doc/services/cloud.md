# Cloud Services

The Cloud Services provide APIs for managing cloud infrastructure organization including locations, zones, nodes, and resource groups.

## Overview

The MOC SDK Cloud Services include:

- **Locations** - Manage geographical locations
- **Zones** - Configure availability zones
- **Nodes** - Manage compute nodes/hosts
- **Groups** - Organize resources into logical groups

## Locations

Locations represent geographical regions or data centers.

### Working with Locations

```go
import (
    "context"
    "github.com/microsoft/moc-sdk-for-go/services/cloud"
    "github.com/microsoft/moc-sdk-for-go/services/cloud/location"
)

func listLocations(locationClient *location.LocationClient) error {
    ctx := context.Background()
    
    locations, err := locationClient.Get(ctx, "", "")
    if err != nil {
        return err
    }
    
    for _, loc := range *locations {
        fmt.Printf("Location: %s\n", *loc.Name)
        if loc.Properties != nil {
            fmt.Printf("  Display Name: %s\n", *loc.Properties.DisplayName)
        }
    }
    
    return nil
}
```

### Creating a Location

```go
func createLocation(locationClient *location.LocationClient) error {
    ctx := context.Background()
    
    locationSpec := &cloud.Location{
        Name: stringPtr("eastus"),
        Properties: &cloud.LocationProperties{
            DisplayName: stringPtr("East US"),
        },
    }
    
    loc, err := locationClient.CreateOrUpdate(ctx, "", "eastus", locationSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created location: %s\n", *loc.Name)
    return nil
}
```

## Zones

Availability zones provide fault isolation and high availability.

### Working with Zones

```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/zone"

func listZones(zoneClient *zone.ZoneClient) error {
    ctx := context.Background()
    
    zones, err := zoneClient.Get(ctx, "eastus", "")
    if err != nil {
        return err
    }
    
    for _, z := range *zones {
        fmt.Printf("Zone: %s\n", *z.Name)
        if z.Properties != nil && z.Properties.ZoneNumber != nil {
            fmt.Printf("  Zone Number: %d\n", *z.Properties.ZoneNumber)
        }
    }
    
    return nil
}
```

### Creating a Zone

```go
func createZone(zoneClient *zone.ZoneClient) error {
    ctx := context.Background()
    
    zoneSpec := &cloud.Zone{
        Name:     stringPtr("zone-1"),
        Location: stringPtr("eastus"),
        Properties: &cloud.ZoneProperties{
            ZoneNumber: int32Ptr(1),
        },
    }
    
    z, err := zoneClient.CreateOrUpdate(ctx, "eastus", "zone-1", zoneSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created zone: %s\n", *z.Name)
    return nil
}
```

### Zone Operations

```go
// Get specific zone
zone, err := zoneClient.Get(ctx, "eastus", "zone-1")

// List all zones in a location
zones, err := zoneClient.Get(ctx, "eastus", "")

// Delete zone
err := zoneClient.Delete(ctx, "eastus", "zone-1")
```

## Nodes

Nodes represent physical or virtual compute hosts in the infrastructure.

### Working with Nodes

```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/node"

func listNodes(nodeClient *node.NodeClient) error {
    ctx := context.Background()
    
    nodes, err := nodeClient.Get(ctx, "eastus", "")
    if err != nil {
        return err
    }
    
    for _, n := range *nodes {
        fmt.Printf("Node: %s\n", *n.Name)
        if n.Properties != nil {
            if n.Properties.FQDN != nil {
                fmt.Printf("  FQDN: %s\n", *n.Properties.FQDN)
            }
            if n.Properties.Port != nil {
                fmt.Printf("  Port: %d\n", *n.Properties.Port)
            }
        }
        if n.Statuses != nil {
            fmt.Printf("  Status: %s\n", *n.Statuses)
        }
    }
    
    return nil
}
```

### Creating a Node

```go
func createNode(nodeClient *node.NodeClient) error {
    ctx := context.Background()
    
    nodeSpec := &cloud.Node{
        Name:     stringPtr("node-01"),
        Location: stringPtr("eastus"),
        Properties: &cloud.NodeProperties{
            FQDN: stringPtr("node01.example.com"),
            Port: int32Ptr(55000),
        },
    }
    
    n, err := nodeClient.CreateOrUpdate(ctx, "eastus", "node-01", nodeSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created node: %s\n", *n.Name)
    return nil
}
```

### Node Operations

```go
// Get specific node
node, err := nodeClient.Get(ctx, "eastus", "node-01")

// Query nodes with filter
nodes, err := nodeClient.Query(ctx, "eastus", "status eq 'Ready'")

// Delete node
err := nodeClient.Delete(ctx, "eastus", "node-01")
```

## Resource Groups

Resource groups organize related resources for management and access control.

### Working with Groups

```go
import "github.com/microsoft/moc-sdk-for-go/services/cloud/group"

func listGroups(groupClient *group.GroupClient) error {
    ctx := context.Background()
    
    groups, err := groupClient.Get(ctx, "")
    if err != nil {
        return err
    }
    
    for _, g := range *groups {
        fmt.Printf("Group: %s\n", *g.Name)
        if g.Location != nil {
            fmt.Printf("  Location: %s\n", *g.Location)
        }
    }
    
    return nil
}
```

### Creating a Resource Group

```go
func createGroup(groupClient *group.GroupClient) error {
    ctx := context.Background()
    
    groupSpec := &cloud.Group{
        Name:     stringPtr("production"),
        Location: stringPtr("eastus"),
        Tags: map[string]*string{
            "environment": stringPtr("production"),
            "cost-center":  stringPtr("engineering"),
        },
    }
    
    g, err := groupClient.CreateOrUpdate(ctx, "production", groupSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created resource group: %s\n", *g.Name)
    return nil
}
```

### Group Operations

```go
// Get specific group
group, err := groupClient.Get(ctx, "production")

// List all groups
groups, err := groupClient.Get(ctx, "")

// Delete group
err := groupClient.Delete(ctx, "production")
```

## Cloud Organization Patterns

### Multi-Location Setup

```go
// Create locations
locations := []string{"eastus", "westus", "centralus"}

for _, locName := range locations {
    locationSpec := &cloud.Location{
        Name: stringPtr(locName),
        Properties: &cloud.LocationProperties{
            DisplayName: stringPtr(locName),
        },
    }
    
    _, err := locationClient.CreateOrUpdate(ctx, "", locName, locationSpec)
    if err != nil {
        log.Printf("Failed to create location %s: %v", locName, err)
    }
}
```

### High Availability with Zones

```go
// Create 3 availability zones for a location
for i := 1; i <= 3; i++ {
    zoneName := fmt.Sprintf("zone-%d", i)
    zoneSpec := &cloud.Zone{
        Name:     stringPtr(zoneName),
        Location: stringPtr("eastus"),
        Properties: &cloud.ZoneProperties{
            ZoneNumber: int32Ptr(int32(i)),
        },
    }
    
    _, err := zoneClient.CreateOrUpdate(ctx, "eastus", zoneName, zoneSpec)
    if err != nil {
        log.Printf("Failed to create zone %d: %v", i, err)
    }
}
```

### Organizing Resources by Environment

```go
// Create resource groups for different environments
environments := []string{"development", "staging", "production"}

for _, env := range environments {
    groupSpec := &cloud.Group{
        Name:     stringPtr(env),
        Location: stringPtr("eastus"),
        Tags: map[string]*string{
            "environment": stringPtr(env),
        },
    }
    
    _, err := groupClient.CreateOrUpdate(ctx, env, groupSpec)
    if err != nil {
        log.Printf("Failed to create group %s: %v", env, err)
    }
}
```

### Node Pool Management

```go
// Create a pool of compute nodes
for i := 1; i <= 5; i++ {
    nodeName := fmt.Sprintf("node-%02d", i)
    nodeSpec := &cloud.Node{
        Name:     stringPtr(nodeName),
        Location: stringPtr("eastus"),
        Properties: &cloud.NodeProperties{
            FQDN: stringPtr(fmt.Sprintf("%s.example.com", nodeName)),
            Port: int32Ptr(55000),
        },
    }
    
    _, err := nodeClient.CreateOrUpdate(ctx, "eastus", nodeName, nodeSpec)
    if err != nil {
        log.Printf("Failed to create node %s: %v", nodeName, err)
    }
}
```

## Tagging Strategy

### Resource Tagging

```go
// Tag resources for organization
tags := map[string]*string{
    "environment":  stringPtr("production"),
    "cost-center":  stringPtr("engineering"),
    "owner":        stringPtr("team-platform"),
    "project":      stringPtr("infrastructure"),
    "created-by":   stringPtr("automation"),
}

groupSpec := &cloud.Group{
    Name:     stringPtr("production"),
    Location: stringPtr("eastus"),
    Tags:     tags,
}
```

### Query by Tags

```go
// List all production resources
groups, err := groupClient.Get(ctx, "")
if err != nil {
    return err
}

for _, g := range *groups {
    if g.Tags != nil {
        if env, ok := g.Tags["environment"]; ok && *env == "production" {
            fmt.Printf("Production group: %s\n", *g.Name)
        }
    }
}
```

## Best Practices

### 1. Logical Resource Organization

```go
// ✅ Good: Organize by environment and workload
// Groups: prod-web, prod-db, dev-web, dev-db

// ❌ Bad: Single group for everything
// Group: all-resources
```

### 2. Use Availability Zones

```go
// ✅ Good: Distribute VMs across zones
vmSpec1.Location = stringPtr("eastus")
vmSpec1.Zones = &[]string{"zone-1"}

vmSpec2.Location = stringPtr("eastus")
vmSpec2.Zones = &[]string{"zone-2"}

// ❌ Bad: All VMs in one zone
```

### 3. Consistent Naming Convention

```go
// ✅ Good: Clear naming pattern
// Locations: eastus, westus, centralus
// Zones: zone-1, zone-2, zone-3
// Nodes: node-01, node-02, node-03
// Groups: prod-web, prod-db

// ❌ Bad: Inconsistent naming
// east, westregion, center-us-1
```

### 4. Tag Everything

```go
// ✅ Good: Comprehensive tagging
Tags: map[string]*string{
    "environment":  stringPtr("production"),
    "cost-center":  stringPtr("engineering"),
    "owner":        stringPtr("team-platform"),
    "project":      stringPtr("web-app"),
}

// ❌ Bad: No tags
Tags: nil
```

## Next Steps

- [Compute Services](compute.md) - Deploy VMs to locations/zones
- [Network Services](network.md) - Network organization
- [Admin Services](admin.md) - Infrastructure management
