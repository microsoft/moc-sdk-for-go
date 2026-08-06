// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"context"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/stretchr/testify/assert"
)

// fakeDiskAttachService is a minimal Service implementation that lets us drive
// DiskAttach without a live cloud agent. Only Get and CreateOrUpdate are used; the
// embedded Service interface satisfies the rest of the contract (they must not be
// called by these tests).
type fakeDiskAttachService struct {
	Service
	vm     *compute.VirtualMachine
	getErr error

	createOrUpdateCalls int
	lastUpdated         *compute.VirtualMachine
}

func (f *fakeDiskAttachService) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	vms := []compute.VirtualMachine{*f.vm}
	return &vms, nil
}

func (f *fakeDiskAttachService) CreateOrUpdate(ctx context.Context, group, name string, vm *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	f.createOrUpdateCalls++
	f.lastUpdated = vm
	return vm, nil
}

func vmWithDataDisks(name string, diskURIs ...string) *compute.VirtualMachine {
	disks := make([]compute.DataDisk, 0, len(diskURIs))
	for i := range diskURIs {
		uri := diskURIs[i]
		disks = append(disks, compute.DataDisk{Vhd: &compute.VirtualHardDisk{URI: &uri}})
	}
	n := name
	return &compute.VirtualMachine{
		Name: &n,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				DataDisks: &disks,
			},
		},
	}
}

func dataDiskURIs(vm *compute.VirtualMachine) []string {
	uris := []string{}
	for _, d := range *vm.StorageProfile.DataDisks {
		uris = append(uris, *d.Vhd.URI)
	}
	return uris
}

// Test_DiskAttach_Idempotency verifies that DiskAttach is a clean no-op success when the
// disk is already attached to the target VM (regression guard for the CSI retry / stranded
// VolumeAttachment bug), while still attaching disks that are not yet present.
func Test_DiskAttach_Idempotency(t *testing.T) {
	ctx := context.Background()
	const (
		group  = "grp"
		vmName = "vm0"
		disk   = "disk-a"
	)

	t.Run("disk already attached to this VM is an idempotent no-op success", func(t *testing.T) {
		fake := &fakeDiskAttachService{vm: vmWithDataDisks(vmName, disk)}
		c := &VirtualMachineClient{internal: fake}

		err := c.DiskAttach(ctx, group, vmName, disk)

		assert.NoError(t, err, "re-attaching an already-attached disk must succeed, not return AlreadyExists")
		assert.Equal(t, 0, fake.createOrUpdateCalls, "no store update should be issued when the disk is already attached")
	})

	t.Run("disk not attached is appended and persisted", func(t *testing.T) {
		fake := &fakeDiskAttachService{vm: vmWithDataDisks(vmName)}
		c := &VirtualMachineClient{internal: fake}

		err := c.DiskAttach(ctx, group, vmName, disk)

		assert.NoError(t, err)
		assert.Equal(t, 1, fake.createOrUpdateCalls, "a new disk must be persisted via CreateOrUpdate")
		if assert.NotNil(t, fake.lastUpdated) {
			assert.Contains(t, dataDiskURIs(fake.lastUpdated), disk, "the new disk must be present in the persisted VM")
		}
	})

	t.Run("attaching a different disk preserves the already-attached one", func(t *testing.T) {
		existing := "disk-existing"
		fake := &fakeDiskAttachService{vm: vmWithDataDisks(vmName, existing)}
		c := &VirtualMachineClient{internal: fake}

		err := c.DiskAttach(ctx, group, vmName, disk)

		assert.NoError(t, err)
		assert.Equal(t, 1, fake.createOrUpdateCalls)
		if assert.NotNil(t, fake.lastUpdated) {
			uris := dataDiskURIs(fake.lastUpdated)
			assert.Contains(t, uris, existing, "the existing disk must be preserved")
			assert.Contains(t, uris, disk, "the new disk must be added")
		}
	})
}
