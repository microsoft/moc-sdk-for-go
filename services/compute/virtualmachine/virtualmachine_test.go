// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/certs"
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
func Test_getVirtualMachineAvailabilitySetProfile(t *testing.T) {
	return
}
func Test_getVirtualMachineAvailabilityZoneProfile(t *testing.T) {
	return
}

func Test_getWssdVirtualMachineProxyConfiguration(t *testing.T) {
	proxy := NewProxy()
	defer proxy.Target.Close()
	HttpProxy := proxy.Target.URL
	HttpsProxy := proxy.Target.URL
	NoProxy := []string{"localhost", "127.0.0.1", ".svc", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.0.0.0/8", ".corp.microsoft.com", ".masd.stbtest.microsoft.com"}

	caCert, _, err := certs.GenerateClientCertificate("ValidCertificate")
	if err != nil {
		t.Fatalf(err.Error())
	}
	certBytes := certs.EncodeCertPEM(caCert)
	TrustedCa := string(certBytes)

	proxyConfig := &compute.ProxyConfiguration{
		HttpProxy:  &HttpProxy,
		HttpsProxy: &HttpsProxy,
		NoProxy:    &NoProxy,
		TrustedCa:  &TrustedCa,
	}

	wssdcloudclient := client{}
	config := wssdcloudclient.getWssdVirtualMachineProxyConfiguration(proxyConfig)

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

	config = wssdcloudclient.getWssdVirtualMachineProxyConfiguration(nil)

	if config != nil && err != nil {
		t.Fatalf("Test_getWssdVirtualMachineProxyConfiguration test case failed: Expected output to be nil since input passed was nil")
	}
}

// Proxy is a simple proxy server for unit tests.
type Proxy struct {
	Target *httptest.Server
}

// NewProxy creates a new proxy server for unit tests.
func NewProxy() *Proxy {
	target := httptest.NewServer(http.DefaultServeMux)
	return &Proxy{Target: target}
}
