// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"context"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/stretchr/testify/assert"
)

type fakeDiskUpdateService struct {
	Service

	calls     int
	group     string
	vmName    string
	dataDisks []compute.DataDisk
	operation compute.VirtualMachineDiskOperation
}

func (f *fakeDiskUpdateService) UpdateDisks(
	_ context.Context,
	group string,
	vmName string,
	dataDisks []compute.DataDisk,
	operation compute.VirtualMachineDiskOperation,
) (*compute.VirtualMachine, error) {
	f.calls++
	f.group = group
	f.vmName = vmName
	f.dataDisks = dataDisks
	f.operation = operation
	return &compute.VirtualMachine{}, nil
}

func TestDiskOperationsUseUpdateDisks(t *testing.T) {
	const (
		group    = "grp"
		vmName   = "vm0"
		diskName = "disk-a"
	)

	tests := []struct {
		name          string
		run           func(*VirtualMachineClient) error
		wantOperation compute.VirtualMachineDiskOperation
	}{
		{
			name: "attach",
			run: func(client *VirtualMachineClient) error {
				return client.DiskAttach(context.Background(), group, vmName, diskName)
			},
			wantOperation: compute.VirtualMachineDiskOperationAttach,
		},
		{
			name: "detach",
			run: func(client *VirtualMachineClient) error {
				return client.DiskDetach(context.Background(), group, vmName, diskName)
			},
			wantOperation: compute.VirtualMachineDiskOperationDetach,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeDiskUpdateService{}
			client := &VirtualMachineClient{internal: service}

			err := tt.run(client)

			assert.NoError(t, err)
			assert.Equal(t, 1, service.calls)
			assert.Equal(t, group, service.group)
			assert.Equal(t, vmName, service.vmName)
			assert.Equal(t, tt.wantOperation, service.operation)
			if assert.Len(t, service.dataDisks, 1) &&
				assert.NotNil(t, service.dataDisks[0].Vhd) &&
				assert.NotNil(t, service.dataDisks[0].Vhd.URI) {
				assert.Equal(t, diskName, *service.dataDisks[0].Vhd.URI)
			}
		})
	}
}
