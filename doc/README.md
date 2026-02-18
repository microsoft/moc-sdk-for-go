# MOC SDK for Go - Documentation

Welcome to the comprehensive documentation for the Microsoft MOC (Management of Cloud) SDK for Go. This SDK provides a Go client library for managing cloud infrastructure resources through gRPC-based APIs.

## Table of Contents

### Getting Started
- [Installation Guide](installation.md) - Set up the SDK in your project
- [Getting Started](getting-started.md) - Quick start guide with basic examples
- [Authentication](authentication.md) - Configure authentication and authorization
- [Client Usage](client-usage.md) - Initialize and configure SDK clients

### Architecture & Design
- [Architecture Overview](architecture.md) - SDK architecture and design principles
- [API Reference](api-reference.md) - Complete API reference guide

### Service Documentation

The MOC SDK provides clients for managing various cloud infrastructure resources:

- **[Compute Services](services/compute.md)** - Virtual machines, scale sets, images, hosts, and availability sets
- **[Network Services](services/network.md)** - Virtual networks, load balancers, network interfaces, and public IPs
- **[Storage Services](services/storage.md)** - Virtual hard disks and storage containers
- **[Security Services](services/security.md)** - Identity, certificates, key vault, roles, and RBAC
- **[Cloud Services](services/cloud.md)** - Locations, zones, nodes, and resource groups
- **[Admin Services](services/admin.md)** - Version management, health monitoring, logging, and validation

### Code Examples

- [Virtual Machine Management](examples/vm-management.md) - Create, manage, and delete VMs

### Development

- [Building from Source](development/building.md) - Build the SDK locally
- [Contributing Guide](development/contributing.md) - Contribute to the project

### Additional Resources

- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [FAQ](faq.md) - Frequently asked questions

## Quick Links

- [GitHub Repository](https://github.com/microsoft/moc-sdk-for-go)
- [MOC Protocol Definitions](https://github.com/microsoft/moc)
- [Contributing Guidelines](../CONTRIBUTING.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)
- [Security Policy](../SECURITY.md)

## Overview

The MOC SDK for Go is a client library that enables Go applications to interact with Microsoft's cloud infrastructure management services. It provides:

- **Type-safe Go interfaces** for all cloud resource operations
- **gRPC-based communication** for efficient and reliable API calls
- **Connection management** with automatic pooling and caching
- **Flexible authentication** supporting multiple authorization methods
- **Resource organization** using groups and hierarchical naming
- **Comprehensive error handling** with detailed error types
- **Cross-platform support** including Windows DLL wrapper

## Key Features

- ✅ Complete CRUD operations for compute, network, storage, and security resources
- ✅ Lifecycle management (start, stop, pause, save) for virtual machines
- ✅ Advanced networking with load balancers, security groups, and VIP pools
- ✅ Identity and access management with RBAC support
- ✅ Key vault integration for secrets and certificate management
- ✅ Resource placement and availability set management
- ✅ Health monitoring and validation capabilities
- ✅ Debug mode for development and testing

## System Requirements

- Go 1.25.0 or later
- Access to a MOC backend service
- Valid authentication credentials
- For Windows DLL: mingw-w64 compiler

## Getting Help

If you encounter issues or have questions:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review the [FAQ](faq.md)
3. Search existing [GitHub Issues](https://github.com/microsoft/moc-sdk-for-go/issues)
4. Open a new issue with detailed information

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

---

**Next Steps:** Start with the [Installation Guide](installation.md) to set up the SDK in your project.
