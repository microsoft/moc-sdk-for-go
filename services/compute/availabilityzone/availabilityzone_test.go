// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityzone

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcommon "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
)

var id = "id"
var name = "avzone1"
var group = "testGroup1"
var location = "mocLocation"
var a1Name = "a1"
var a1Group = "ag1"
var a2Name = "a2"
var a2Group = "ag2"
var provisionoingstate = "CREATED"
var health = "OK"
var wssdnodes = []string {"node1", "node2"}
var wssdvms = []*compute.VirtualMachineReference{{Name: &a1Name, GroupName: &a1Group}, {Name: &a2Name, GroupName: &a2Group}}
var rpcvms = []*wssdcloudcompute.VirtualMachineReference{{Name: a1Name, GroupName: a1Group}, {Name: a2Name, GroupName: a2Group}}
var wssdstatus = map[string]*string{
	"ProvisionState": &provisionoingstate,
	"HealthState":    &health,
}
var rpcstatus = wssdcommon.Status{
	Health:             &wssdcommon.Health{CurrentState: 1},
	ProvisioningStatus: &wssdcommon.ProvisionStatus{CurrentState: 2},
	Version:            &wssdcommon.Version{Number: "123"},
}

func Test_getRpcAvailabilityZone(t *testing.T) {
	result, err := getRpcAvailabilityZone(nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	avzone := compute.AvailabilityZone{
		Name:                     &name,
		Location:                 &location,
		VirtualMachines:          wssdvms,
		Statuses:                 wssdstatus,
		Nodes:                    wssdnodes,
	}

	result, err = getRpcAvailabilityZone(&avzone)
	assert.Error(t, err)

	avzone.ID = &id
	result, err = getRpcAvailabilityZone(&avzone)
	assert.Nil(t, err)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, location, result.LocationName)
	assert.Equal(t, 2, len(result.VirtualMachines))
}


func Test_getWssdAvailabilityZone(t *testing.T) {
	result, err := getWssdAvailabilityZone(nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	avset := wssdcloudcompute.AvailabilityZone{
		Name:                     name,
		LocationName:             location,
		VirtualMachines:          rpcvms,
		Status:                   &rpcstatus,
		Nodes:                    wssdnodes,
	}

	result, err = getWssdAvailabilityZone(&avzone)
	assert.Nil(t, err)
	assert.EqualValues(t, name, *result.Name)
	assert.EqualValues(t, location, *result.Location)
	assert.EqualValues(t, 2, len(result.VirtualMachines))
}

func Test_getRpcWssdVirtualMachineReference(t *testing.T) {
	a1Name := "a1"
	a1Group := "ag1"
	a2Name := "a2"
	a2Group := "ag2"
	a := []*compute.VirtualMachineReference{{Name: &a1Name, GroupName: &a1Group}, {Name: &a2Name, GroupName: &a2Group}}
	b := getRpcVirtualMachineReferences(a)
	assert.Equal(t, 2, len(b))
	assert.EqualValues(t, a1Name, b[0].Name)
	assert.EqualValues(t, a1Group, b[0].GroupName)
	assert.EqualValues(t, a2Name, b[1].Name)
	assert.EqualValues(t, a2Group, b[1].GroupName)

	b = getRpcVirtualMachineReferences([]*compute.VirtualMachineReference{})
	assert.Equal(t, 0, len(b))

	b = getRpcVirtualMachineReferences([]*compute.VirtualMachineReference{nil})
	assert.Equal(t, 1, len(b))
	assert.Nil(t, b[0])
}