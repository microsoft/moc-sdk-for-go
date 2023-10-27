// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"io/ioutil"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
)

func Test_getWssdVirtualMachine(t *testing.T) {
}

func Test_getWssdVirtualMachineStorageConfiguration(t *testing.T) {}

func Test_getWssdVirtualMachineStorageConfigurationOsDisk(t *testing.T) {}

func Test_getWssdVirtualMachineStorageConfigurationDataDisks(t *testing.T) {}

func Test_getWssdVirtualMachineNetworkConfiguration(t *testing.T) {}

func Test_getWssdVirtualMachineOSSSHPublicKeys(t *testing.T) {}
func Test_getWssdVirtualMachineOSConfiguration(t *testing.T) {}

func Test_getVirtualMachine(t *testing.T)                        {}
func Test_getVirtualMachineStorageProfile(t *testing.T)          {}
func Test_getVirtualMachineStorageProfileOsDisk(t *testing.T)    {}
func Test_getVirtualMachineStorageProfileDataDisks(t *testing.T) {}
func Test_getVirtualMachineNetworkProfile(t *testing.T)          {}
func Test_getVirtualMachineOSProfile(t *testing.T)               {}

func Test_getWssdVirtualMachineHttpProxyConfiguration(t *testing.T) {
	httpProxyConfig := compute.OSProfile.ProxyConfiguration

	httpProxyConfig.HttpProxy = "http://ubuntu:ubuntu@192.168.200.40:3128"
	httpProxyConfig.HttpsProxy = "http://ubuntu:ubuntu@192.168.200.40:3128"
	httpProxyConfig.NoProxy = []string{"localhost", "127.0.0.1", ".svc", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.0.0.0/8", ".corp.microsoft.com", ".masd.stbtest.microsoft.com"}

	caCert, err := ioutil.ReadFile("proxy.crt")
	if err != nil {
		t.Fatalf(err.Error())
	}
	httpProxyConfig.TrustedCa = string(caCert)

	wssdcloudclient := client{}
	config := wssdcloudclient.getWssdVirtualMachineHttpProxyConfiguration(httpProxyConfig)

	if config.HttpProxy != httpProxyConfig.HttpProxy {
		t.Fatalf("Test_getWssdVirtualMachineHttpProxyConfiguration test case failed: HttpProxy does not match")
	}

	if config.HttpsProxy != httpProxyConfig.HttpsProxy {
		t.Fatalf("Test_getWssdVirtualMachineHttpProxyConfiguration test case failed: HttpsProxy does not match")
	}

	if config.NoProxy != httpProxyConfig.NoProxy {
		t.Fatalf("Test_getWssdVirtualMachineHttpProxyConfiguration test case failed: NoProxy does not match")
	}

	if config.TrustedCa != httpProxyConfig.TrustedCa {
		t.Fatalf("Test_getWssdVirtualMachineHttpProxyConfiguration test case failed: TrustedCa does not match")
	}
}
