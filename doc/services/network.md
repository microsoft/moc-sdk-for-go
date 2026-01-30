# Network Services

The Network Services provide APIs for managing virtual networks, load balancers, network interfaces, and network security.

## Overview

The MOC SDK Network Services include:

- **Virtual Networks** - Create and manage virtual networks and subnets
- **Network Interfaces** - Manage VM network interfaces (NICs)
- **Load Balancers** - Configure load balancing for VMs
- **Public IP Addresses** - Manage public IP addresses
- **Network Security Groups** - Define network security rules
- **Logical Networks** - Configure underlying network fabric
- **MAC Pools** - Manage MAC address pools
- **VIP Pools** - Manage virtual IP address pools

## Virtual Networks

Virtual networks provide network isolation and connectivity for VMs.

### Creating a Virtual Network

```go
import (
    "context"
    "github.com/microsoft/moc-sdk-for-go/services/network"
    "github.com/microsoft/moc-sdk-for-go/services/network/virtualnetwork"
)

func createVNet(vnetClient *virtualnetwork.VirtualNetworkClient) error {
    ctx := context.Background()
    
    vnetSpec := &network.VirtualNetwork{
        Name:     stringPtr("prod-vnet"),
        Location: stringPtr("default"),
        Properties: &network.VirtualNetworkPropertiesFormat{
            AddressSpace: &network.AddressSpace{
                AddressPrefixes: &[]string{"10.0.0.0/16"},
            },
            Subnets: &[]network.Subnet{
                {
                    Name: stringPtr("default-subnet"),
                    Properties: &network.SubnetPropertiesFormat{
                        AddressPrefix: stringPtr("10.0.1.0/24"),
                    },
                },
                {
                    Name: stringPtr("web-subnet"),
                    Properties: &network.SubnetPropertiesFormat{
                        AddressPrefix: stringPtr("10.0.2.0/24"),
                    },
                },
            },
        },
    }
    
    vnet, err := vnetClient.CreateOrUpdate(ctx, "production", "prod-vnet", vnetSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created VNet: %s\n", *vnet.Name)
    return nil
}
```

### Virtual Network Operations

```go
// Get virtual network
vnet, err := vnetClient.Get(ctx, "production", "prod-vnet")

// List all virtual networks
vnets, err := vnetClient.Get(ctx, "production", "")

// Update virtual network
updatedVNet, err := vnetClient.CreateOrUpdate(ctx, "production", "prod-vnet", vnetSpec)

// Delete virtual network
err := vnetClient.Delete(ctx, "production", "prod-vnet")
```

## Network Interfaces

Network interfaces connect VMs to virtual networks.

### Creating a Network Interface

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/networkinterface"

func createNIC(nicClient *networkinterface.NetworkInterfaceClient) error {
    ctx := context.Background()
    
    nicSpec := &network.NetworkInterface{
        Name:     stringPtr("web-nic-01"),
        Location: stringPtr("default"),
        Properties: &network.NetworkInterfacePropertiesFormat{
            IPConfigurations: &[]network.NetworkInterfaceIPConfiguration{
                {
                    Name: stringPtr("ipconfig1"),
                    Properties: &network.NetworkInterfaceIPConfigurationPropertiesFormat{
                        Subnet: &network.Subnet{
                            ID: stringPtr("/production/virtualnetworks/prod-vnet/subnets/web-subnet"),
                        },
                        PrivateIPAllocationMethod: network.IPAllocationMethodDynamic,
                    },
                },
            },
        },
    }
    
    nic, err := nicClient.CreateOrUpdate(ctx, "production", "web-nic-01", nicSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created NIC: %s\n", *nic.Name)
    return nil
}
```

### Network Interface Operations

```go
// Get network interface
nic, err := nicClient.Get(ctx, "production", "web-nic-01")

// Get MAC address
if nic.Properties != nil && nic.Properties.MacAddress != nil {
    fmt.Printf("MAC Address: %s\n", *nic.Properties.MacAddress)
}

// List all network interfaces
nics, err := nicClient.Get(ctx, "production", "")

// Delete network interface
err := nicClient.Delete(ctx, "production", "web-nic-01")
```

### Static IP Configuration

```go
nicSpec := &network.NetworkInterface{
    Name:     stringPtr("web-nic-01"),
    Location: stringPtr("default"),
    Properties: &network.NetworkInterfacePropertiesFormat{
        IPConfigurations: &[]network.NetworkInterfaceIPConfiguration{
            {
                Name: stringPtr("ipconfig1"),
                Properties: &network.NetworkInterfaceIPConfigurationPropertiesFormat{
                    Subnet: &network.Subnet{
                        ID: stringPtr("/production/virtualnetworks/prod-vnet/subnets/web-subnet"),
                    },
                    PrivateIPAddress:          stringPtr("10.0.2.10"),
                    PrivateIPAllocationMethod: network.IPAllocationMethodStatic,
                },
            },
        },
    },
}
```

## Load Balancers

Load balancers distribute traffic across multiple VMs.

### Creating a Load Balancer

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/loadbalancer"

func createLoadBalancer(lbClient *loadbalancer.LoadBalancerClient) error {
    ctx := context.Background()
    
    lbSpec := &network.LoadBalancer{
        Name:     stringPtr("web-lb"),
        Location: stringPtr("default"),
        Properties: &network.LoadBalancerPropertiesFormat{
            FrontendIPConfigurations: &[]network.FrontendIPConfiguration{
                {
                    Name: stringPtr("frontend-1"),
                    Properties: &network.FrontendIPConfigurationPropertiesFormat{
                        PrivateIPAddress:          stringPtr("10.0.1.10"),
                        PrivateIPAllocationMethod: network.IPAllocationMethodStatic,
                    },
                },
            },
            BackendAddressPools: &[]network.BackendAddressPool{
                {
                    Name: stringPtr("backend-pool"),
                },
            },
            LoadBalancingRules: &[]network.LoadBalancingRule{
                {
                    Name: stringPtr("http-rule"),
                    Properties: &network.LoadBalancingRulePropertiesFormat{
                        FrontendPort: int32Ptr(80),
                        BackendPort:  int32Ptr(80),
                        Protocol:     network.TransportProtocolTCP,
                    },
                },
            },
        },
    }
    
    lb, err := lbClient.CreateOrUpdate(ctx, "production", "web-lb", lbSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created load balancer: %s\n", *lb.Name)
    return nil
}
```

### Load Balancer Operations

```go
// Get load balancer
lb, err := lbClient.Get(ctx, "production", "web-lb")

// List all load balancers
lbs, err := lbClient.Get(ctx, "production", "")

// Delete load balancer
err := lbClient.Delete(ctx, "production", "web-lb")
```

## Public IP Addresses

Manage public IP addresses for external connectivity.

### Creating a Public IP

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/publicipaddress"

func createPublicIP(pipClient *publicipaddress.PublicIPAddressClient) error {
    ctx := context.Background()
    
    pipSpec := &network.PublicIPAddress{
        Name:     stringPtr("web-public-ip"),
        Location: stringPtr("default"),
        Properties: &network.PublicIPAddressPropertiesFormat{
            PublicIPAllocationMethod: network.IPAllocationMethodStatic,
        },
    }
    
    pip, err := pipClient.CreateOrUpdate(ctx, "production", "web-public-ip", pipSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created public IP: %s\n", *pip.Properties.IPAddress)
    return nil
}
```

### Public IP Operations

```go
// Get public IP
pip, err := pipClient.Get(ctx, "production", "web-public-ip")

// Get IP address
if pip.Properties != nil && pip.Properties.IPAddress != nil {
    fmt.Printf("Public IP: %s\n", *pip.Properties.IPAddress)
}

// Delete public IP
err := pipClient.Delete(ctx, "production", "web-public-ip")
```

## Network Security Groups

Define network security rules for traffic filtering.

### Creating a Network Security Group

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/networksecuritygroup"

func createNSG(nsgClient *networksecuritygroup.NetworkSecurityGroupClient) error {
    ctx := context.Background()
    
    nsgSpec := &network.NetworkSecurityGroup{
        Name:     stringPtr("web-nsg"),
        Location: stringPtr("default"),
        Properties: &network.NetworkSecurityGroupPropertiesFormat{
            SecurityRules: &[]network.SecurityRule{
                {
                    Name: stringPtr("allow-http"),
                    Properties: &network.SecurityRulePropertiesFormat{
                        Protocol:                 network.SecurityRuleProtocolTCP,
                        SourcePortRange:          stringPtr("*"),
                        DestinationPortRange:     stringPtr("80"),
                        SourceAddressPrefix:      stringPtr("*"),
                        DestinationAddressPrefix: stringPtr("*"),
                        Access:                   network.SecurityRuleAccessAllow,
                        Priority:                 int32Ptr(100),
                        Direction:                network.SecurityRuleDirectionInbound,
                    },
                },
                {
                    Name: stringPtr("allow-https"),
                    Properties: &network.SecurityRulePropertiesFormat{
                        Protocol:                 network.SecurityRuleProtocolTCP,
                        SourcePortRange:          stringPtr("*"),
                        DestinationPortRange:     stringPtr("443"),
                        SourceAddressPrefix:      stringPtr("*"),
                        DestinationAddressPrefix: stringPtr("*"),
                        Access:                   network.SecurityRuleAccessAllow,
                        Priority:                 int32Ptr(110),
                        Direction:                network.SecurityRuleDirectionInbound,
                    },
                },
            },
        },
    }
    
    nsg, err := nsgClient.CreateOrUpdate(ctx, "production", "web-nsg", nsgSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created NSG: %s\n", *nsg.Name)
    return nil
}
```

### Applying NSG to Subnet

```go
vnetSpec.Properties.Subnets = &[]network.Subnet{
    {
        Name: stringPtr("web-subnet"),
        Properties: &network.SubnetPropertiesFormat{
            AddressPrefix: stringPtr("10.0.2.0/24"),
            NetworkSecurityGroup: &network.NetworkSecurityGroup{
                ID: stringPtr("/production/networksecuritygroups/web-nsg"),
            },
        },
    },
}
```

## Logical Networks

Configure the underlying network fabric.

### Creating a Logical Network

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/logicalnetwork"

func createLogicalNetwork(lnClient *logicalnetwork.LogicalNetworkClient) error {
    ctx := context.Background()
    
    lnSpec := &network.LogicalNetwork{
        Name:     stringPtr("fabric-network"),
        Location: stringPtr("default"),
        Properties: &network.LogicalNetworkPropertiesFormat{
            Subnets: &[]network.LogicalSubnet{
                {
                    Name: stringPtr("subnet-1"),
                    Properties: &network.LogicalSubnetPropertiesFormat{
                        AddressPrefix: stringPtr("192.168.1.0/24"),
                        Vlan:          int32Ptr(100),
                    },
                },
            },
        },
    }
    
    ln, err := lnClient.CreateOrUpdate(ctx, "default", "fabric-network", lnSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created logical network: %s\n", *ln.Name)
    return nil
}
```

## MAC and VIP Pools

### Managing MAC Pools

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/macpool"

func createMACPool(macClient *macpool.MacPoolClient) error {
    ctx := context.Background()
    
    macSpec := &network.MacPool{
        Name:     stringPtr("mac-pool-1"),
        Location: stringPtr("default"),
        Properties: &network.MacPoolPropertiesFormat{
            Range: &network.MacRange{
                StartMacAddress: stringPtr("00-15-5D-00-00-00"),
                EndMacAddress:   stringPtr("00-15-5D-FF-FF-FF"),
            },
        },
    }
    
    mac, err := macClient.CreateOrUpdate(ctx, "default", "mac-pool-1", macSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created MAC pool: %s\n", *mac.Name)
    return nil
}
```

### Managing VIP Pools

```go
import "github.com/microsoft/moc-sdk-for-go/services/network/vippool"

func createVIPPool(vipClient *vippool.VipPoolClient) error {
    ctx := context.Background()
    
    vipSpec := &network.VipPool{
        Name:     stringPtr("vip-pool-1"),
        Location: stringPtr("default"),
        Properties: &network.VipPoolPropertiesFormat{
            StartIP: stringPtr("10.0.100.10"),
            EndIP:   stringPtr("10.0.100.50"),
        },
    }
    
    vip, err := vipClient.CreateOrUpdate(ctx, "default", "vip-pool-1", vipSpec)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created VIP pool: %s\n", *vip.Name)
    return nil
}
```

## Common Networking Patterns

### VM with Public IP

```go
// Create public IP
pipSpec := &network.PublicIPAddress{
    Name:     stringPtr("vm-public-ip"),
    Location: stringPtr("default"),
    Properties: &network.PublicIPAddressPropertiesFormat{
        PublicIPAllocationMethod: network.IPAllocationMethodStatic,
    },
}
pip, _ := pipClient.CreateOrUpdate(ctx, "production", "vm-public-ip", pipSpec)

// Create NIC with public IP
nicSpec := &network.NetworkInterface{
    Name:     stringPtr("vm-nic"),
    Location: stringPtr("default"),
    Properties: &network.NetworkInterfacePropertiesFormat{
        IPConfigurations: &[]network.NetworkInterfaceIPConfiguration{
            {
                Name: stringPtr("ipconfig1"),
                Properties: &network.NetworkInterfaceIPConfigurationPropertiesFormat{
                    Subnet: &network.Subnet{
                        ID: stringPtr("/production/virtualnetworks/prod-vnet/subnets/default"),
                    },
                    PrivateIPAllocationMethod: network.IPAllocationMethodDynamic,
                    PublicIPAddress: &network.PublicIPAddress{
                        ID: pip.ID,
                    },
                },
            },
        },
    },
}
```

### Load Balanced VMs

```go
// Add NIC to load balancer backend pool
nicSpec.Properties.IPConfigurations[0].Properties.LoadBalancerBackendAddressPools = &[]network.BackendAddressPool{
    {
        ID: stringPtr("/production/loadbalancers/web-lb/backendAddressPools/backend-pool"),
    },
}
```

## Next Steps

- [Compute Services](compute.md) - Create VMs with networking
- [Network Examples](../examples/network-setup.md) - Detailed examples
- [Security Services](security.md) - Network security
