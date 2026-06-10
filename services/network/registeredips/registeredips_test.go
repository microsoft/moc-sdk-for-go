// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package registeredips

import (
	"encoding/json"
	"reflect"
	"testing"

	wssdcloudnetwork "github.com/microsoft/moc/rpc/cloudagent/network"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
	"gopkg.in/yaml.v3"
)

func TestIPUpdateErrorCode_String(t *testing.T) {
	cases := map[IPUpdateErrorCode]string{
		IPUpdateUnknown:          "Unknown",
		IPUpdateInvalidFormat:    "InvalidFormat",
		IPUpdateOutOfRange:       "OutOfRange",
		IPUpdateSubnetNotFound:   "SubnetNotFound",
		IPUpdateAlreadyAllocated: "AlreadyAllocated",
		IPUpdateNoPoolsInSubnet:  "NoPoolsInSubnet",
		IPUpdateErrorCode(99):    "IPUpdateErrorCode(99)",
	}
	for code, want := range cases {
		if got := code.String(); got != want {
			t.Errorf("String(%d)=%q want %q", int32(code), got, want)
		}
	}
}

func TestIPUpdateErrorCode_JSONRoundTrip(t *testing.T) {
	for _, code := range []IPUpdateErrorCode{
		IPUpdateUnknown,
		IPUpdateInvalidFormat,
		IPUpdateOutOfRange,
		IPUpdateSubnetNotFound,
		IPUpdateAlreadyAllocated,
		IPUpdateNoPoolsInSubnet,
	} {
		data, err := json.Marshal(code)
		if err != nil {
			t.Fatalf("Marshal(%v): %v", code, err)
		}
		var got IPUpdateErrorCode
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal(%s): %v", data, err)
		}
		if got != code {
			t.Errorf("round-trip %v -> %s -> %v", code, data, got)
		}
	}
}

func TestIPUpdateErrorCode_UnmarshalNumeric(t *testing.T) {
	var c IPUpdateErrorCode
	if err := json.Unmarshal([]byte("4"), &c); err != nil {
		t.Fatalf("Unmarshal(4): %v", err)
	}
	if c != IPUpdateAlreadyAllocated {
		t.Errorf("Unmarshal(4)=%v want %v", c, IPUpdateAlreadyAllocated)
	}
}

func TestIPUpdateErrorCode_UnmarshalUnknownString(t *testing.T) {
	var c IPUpdateErrorCode
	if err := json.Unmarshal([]byte(`"Bogus"`), &c); err == nil {
		t.Errorf("expected error for unknown string, got %v", c)
	}
}

func TestIPUpdateErrorCode_MarshalYAML(t *testing.T) {
	data, err := yaml.Marshal(IPUpdateOutOfRange)
	if err != nil {
		t.Fatalf("yaml.Marshal: %v", err)
	}
	if string(data) != "OutOfRange\n" {
		t.Errorf("yaml=%q want %q", data, "OutOfRange\n")
	}
}

func TestSubnetsFromProto_Logical(t *testing.T) {
	in := []*wssdcloudnetwork.LogicalSubnetIPUpdate{
		{SubnetName: "s1", RegisteredIPAddresses: []string{"10.0.0.1", "10.0.0.2"}},
		nil, // skipped
		{SubnetName: "s2", RegisteredIPAddresses: nil},
	}
	got := SubnetsFromProto(in)
	want := []SubnetRegisteredIPs{
		{SubnetName: "s1", RegisteredIPAddresses: []string{"10.0.0.1", "10.0.0.2"}},
		{SubnetName: "s2", RegisteredIPAddresses: []string{}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got=%+v want=%+v", got, want)
	}
	// Cleared subnets must produce non-nil empty slice.
	if got[1].RegisteredIPAddresses == nil {
		t.Errorf("cleared subnet got nil slice, want non-nil empty slice")
	}
}

func TestSubnetsFromProto_Virtual(t *testing.T) {
	in := []*wssdcloudnetwork.VirtualSubnetIPUpdate{
		{SubnetName: "vs1", RegisteredIPAddresses: []string{"192.168.0.5"}},
	}
	got := SubnetsFromProto(in)
	if len(got) != 1 || got[0].SubnetName != "vs1" || len(got[0].RegisteredIPAddresses) != 1 {
		t.Errorf("got=%+v", got)
	}
}

func TestSubnetsFromProto_Empty(t *testing.T) {
	if got := SubnetsFromProto([]*wssdcloudnetwork.LogicalSubnetIPUpdate(nil)); got != nil {
		t.Errorf("nil-input got=%+v want nil", got)
	}
}

func TestFailuresFromProto(t *testing.T) {
	in := []*wssdcloudcommon.IPAddressUpdateFailure{
		{SubnetName: "s1", IPAddress: "10.0.0.99", Code: wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_OUT_OF_RANGE, Error: "out of range"},
		nil, // skipped
		{SubnetName: "s2", IPAddress: "", Code: wssdcloudcommon.IPUpdateErrorCode_IP_UPDATE_NO_POOLS_IN_SUBNET, Error: "no pools"},
	}
	got := FailuresFromProto(in)
	want := []IPAddressUpdateFailure{
		{SubnetName: "s1", IPAddress: "10.0.0.99", Code: IPUpdateOutOfRange, Error: "out of range"},
		{SubnetName: "s2", IPAddress: "", Code: IPUpdateNoPoolsInSubnet, Error: "no pools"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got=%+v want=%+v", got, want)
	}
}

func TestFailuresFromProto_Empty(t *testing.T) {
	if got := FailuresFromProto(nil); got != nil {
		t.Errorf("nil-input got=%+v want nil", got)
	}
}

func TestSubnetRegisteredIPs_JSONTags(t *testing.T) {
	v := SubnetRegisteredIPs{SubnetName: "s1", RegisteredIPAddresses: []string{"10.0.0.1"}}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want := `{"subnetName":"s1","registeredIPAddresses":["10.0.0.1"]}`
	if string(data) != want {
		t.Errorf("got=%s want=%s", data, want)
	}
}

func TestIPAddressUpdateFailure_JSONTags(t *testing.T) {
	v := IPAddressUpdateFailure{SubnetName: "s1", IPAddress: "10.0.0.99", Code: IPUpdateOutOfRange, Error: "boom"}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want := `{"subnetName":"s1","ipAddress":"10.0.0.99","code":"OutOfRange","error":"boom"}`
	if string(data) != want {
		t.Errorf("got=%s want=%s", data, want)
	}
}
