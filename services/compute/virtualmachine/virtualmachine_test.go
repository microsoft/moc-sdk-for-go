// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"io/ioutil"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
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
func Test_getWssdVirtualMachineProxyConfiguration(t *testing.T) {
	HttpProxy := "http://akse2e:akse2e@skyproxy.ceccloud1.selfhost.corp.microsoft.com:3128"
	HttpsProxy := "http://akse2e:akse2e@skyproxy.ceccloud1.selfhost.corp.microsoft.com:3128"
	NoProxy := []string{"localhost", "127.0.0.1", ".svc", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.0.0.0/8", ".corp.microsoft.com", ".masd.stbtest.microsoft.com"}

	caCert, err := ioutil.ReadFile("proxy.crt")
	if err != nil {
		t.Fatalf(err.Error())
	}
	TrustedCa = string(caCert)

	proxyConfig := &compute.ProxyConfiguration{
		HttpProxy:  &HttpProxy,
		HttpsProxy: &HttpsProxy,
		NoProxy:    &NoProxy,
		TrustedCa:  &TrustedCa,
	}

	wssdcloudclient := client{}
	config, _ := wssdcloudclient.getWssdVirtualMachineProxyConfiguration(proxyConfig, wssdcloudproto.Operation_POST)

	if config.HttpProxy != HttpProxy {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: HttpProxy does not match")
	}

	if config.HttpsProxy != HttpsProxy {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: HttpsProxy does not match")
	}

	if len(config.NoProxy) != len(NoProxy) {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: NoProxy does not match")
	}

	if config.TrustedCa != TrustedCa {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: TrustedCa does not match")
	}

	config, err := wssdcloudclient.getWssdVirtualMachineProxyConfiguration(nil, wssdcloudproto.Operation_POST)

	if config != nil && err != nil {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: Expected output to be nil since input passed was nil")
	}

	HttpProxy = "https://akse2e:akse2e@skyproxy.ceccloud1.selfhost.corp.microsoft.com:3128"
	proxyConfig.HttpProxy = &HttpProxy
	config, err = wssdcloudclient.getWssdVirtualMachineProxyConfiguration(proxyConfig, wssdcloudproto.Operation_POST)

	if err == nil && config != nil {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: Expected proxy config validation to fail but passed")
	}
}
