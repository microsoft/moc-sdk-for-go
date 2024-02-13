// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package availabilityset

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
	wssdcommon "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
)

var id = "id"
var name = "avset1"
var group = "testGroup1"
var fdcount int32 = 2
var location = "mocLocation"
var a1Name = "a1"
var a1Group = "ag1"
var a2Name = "a2"
var a2Group = "ag2"
var provisionoingstate = "CREATED"
var health = "OK"
var wssdvms = []*compute.CloudSubResource{{Name: &a1Name, GroupName: &a1Group}, {Name: &a2Name, GroupName: &a2Group}}
var rpcvms = []*wssdcommon.CloudSubResource{{Name: a1Name, GroupName: a1Group}, {Name: a2Name, GroupName: a2Group}}
var wssdstatus = map[string]*string{
	"ProvisionState": &provisionoingstate,
	"HealthState":    &health,
}
var rpcstatus = wssdcommon.Status{
	Health:             &wssdcommon.Health{CurrentState: 1},
	ProvisioningStatus: &wssdcommon.ProvisionStatus{CurrentState: 2},
	Version:            &wssdcommon.Version{Number: "123"},
}

func Test_getRpcAvailabilitySet(t *testing.T) {
	result, err := getRpcAvailabilitySet(nil, "testGroup")
	assert.Error(t, err)
	assert.Nil(t, result)

	avset := compute.AvailabilitySet{
		Name:                     &name,
		PlatformFaultDomainCount: &fdcount,
		Location:                 &location,
		VirtualMachines:          wssdvms,
		Statuses:                 wssdstatus,
	}

	result, err = getRpcAvailabilitySet(&avset, "group1")
	assert.Error(t, err)

	avset.ID = &id
	result, err = getRpcAvailabilitySet(&avset, "group1")
	assert.Nil(t, err)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, fdcount, result.PlatformFaultDomainCount)
	assert.Equal(t, location, result.LocationName)
	assert.EqualValues(t, "group1", result.GroupName)
	assert.Equal(t, 2, len(result.VirtualMachines))
}

func Test_getWssdAvailabilitySet(t *testing.T) {
	result, err := getWssdAvailabilitySet(nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	avset := wssdcloudcompute.AvailabilitySet{
		Name:                     name,
		PlatformFaultDomainCount: fdcount,
		LocationName:             location,
		GroupName:                "group1",
		VirtualMachines:          rpcvms,
		Status:                   &rpcstatus,
	}

	result, err = getWssdAvailabilitySet(&avset)
	assert.Nil(t, err)
	assert.EqualValues(t, name, *result.Name)
	assert.EqualValues(t, fdcount, *result.PlatformFaultDomainCount)
	assert.EqualValues(t, location, *result.Location)
	assert.EqualValues(t, 2, len(result.VirtualMachines))
}

func Test_getRpcWssdSubResources(t *testing.T) {
	a1Name := "a1"
	a1Group := "ag1"
	a2Name := "a2"
	a2Group := "ag2"
	a := []*compute.CloudSubResource{{Name: &a1Name, GroupName: &a1Group}, {Name: &a2Name, GroupName: &a2Group}}
	b := getRpcSubResources(a)
	assert.Equal(t, 2, len(b))
	assert.EqualValues(t, a1Name, b[0].Name)
	assert.EqualValues(t, a1Group, b[0].GroupName)
	assert.EqualValues(t, a2Name, b[1].Name)
	assert.EqualValues(t, a2Group, b[1].GroupName)

	b = getRpcSubResources([]*compute.CloudSubResource{})
	assert.Equal(t, 0, len(b))

	b = getRpcSubResources([]*compute.CloudSubResource{nil})
	assert.Equal(t, 1, len(b))
	assert.Nil(t, b[0])
}
