# Tests

> Last updated: 2026-01-11

## Test Framework

- **Framework**: Go's built-in `testing` package
- **Assertions**: `github.com/stretchr/testify/assert`
- **Mocking**: `github.com/stretchr/testify/mock`
- **Comparison**: `github.com/google/go-cmp/cmp`

## Running Tests

```bash
# Run all tests
make test
# or
GOARCH=amd64 go test -v ./...

# Run unit tests only (client and security packages)
make unittest
# or
GOARCH=amd64 go test -v ./pkg/client/...
GOARCH=amd64 go test -v ./services/security/...

# Run specific package tests
go test -v ./services/compute/virtualmachine/...

# Run specific test
go test -v ./services/security/authentication/... -run Test_LoginWithConfig
```

## Test Organization

Tests are co-located with source files using `*_test.go` naming:

```
services/compute/virtualmachine/
├── virtualmachine.go
├── virtualmachine_test.go    # SDK type conversion tests
├── wssd.go
└── wssd_test.go              # gRPC layer tests

services/security/authentication/
├── casigned.go
└── casigned_test.go          # Auth timing/mock tests
```

## Test Files

| File | Coverage |
|------|----------|
| `pkg/client/client_test.go` | Connection management |
| `services/cloud/group/group_test.go` | Group type conversions |
| `services/cloud/location/location_test.go` | Location type conversions |
| `services/cloud/node/node_test.go` | Node type conversions |
| `services/cloud/zone/zone_test.go` | Zone type conversions |
| `services/compute/availabilityset/availabilityset_test.go` | AvailabilitySet tests |
| `services/compute/placementgroup/placementgroup_test.go` | PlacementGroup tests |
| `services/compute/virtualmachine/virtualmachine_test.go` | VM type tests |
| `services/compute/virtualmachine/wssd_test.go` | VM gRPC tests |
| `services/network/networkinterface/wssd_test.go` | NIC gRPC tests |
| `services/security/authentication/casigned_test.go` | Auth flow tests |
| `services/security/identity/identity_test.go` | Identity tests |
| `services/security/keyvault/key/wssd_test.go` | KeyVault key tests |
| `services/security/role/role_test.go` | Role tests |
| `services/security/roleassignment/roleassignment_test.go` | RoleAssignment tests |

## Writing New Tests

### Unit Test Pattern

```go
package virtualmachine

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test_ConvertToSdkType(t *testing.T) {
    // Arrange
    protoVM := &wssdcompute.VirtualMachine{
        Name: "test-vm",
    }

    // Act
    result := getVirtualMachine(protoVM, "test-group")

    // Assert
    assert.NotNil(t, result)
    assert.Equal(t, "test-vm", *result.Name)
}
```

### Mock Client Pattern

```go
import (
    "github.com/stretchr/testify/mock"
)

type MockClient struct {
    mock.Mock
}

func (m *MockClient) Login(ctx context.Context, req *security.AuthenticationRequest, opts ...grpc.CallOption) (*security.AuthenticationResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*security.AuthenticationResponse), args.Error(1)
}

func Test_LoginWithMock(t *testing.T) {
    mockClient := new(MockClient)
    mockClient.On("Login", mock.Anything, mock.AnythingOfType("*security.AuthenticationRequest")).
        Return(&security.AuthenticationResponse{Token: "test"}, nil)

    // ... use mockClient

    mockClient.AssertExpectations(t)
}
```

## Test Notes

- Tests require `GOARCH=amd64` when running on different architectures
- Integration tests require a running MOC cloud agent
- Unit tests focus on type conversions and business logic
- Mock clients used for gRPC layer testing without backend
