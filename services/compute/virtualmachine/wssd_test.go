package virtualmachine

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
	"github.com/stretchr/testify/assert"
)

func Test_VirtualMachineValidations(t *testing.T) {
	wssdcloudclient := client{}
	httpProxy := ""
	httpsProxy := ""

	// empty url
	vm := &compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile: &compute.OSProfile{
				ProxyConfiguration: &compute.ProxyConfiguration{
					HttpProxy:  &httpProxy,
					HttpsProxy: &httpsProxy,
				},
			},
		},
	}

	err := wssdcloudclient.virtualMachineValidations(wssdcloudproto.Operation_POST, vm)

	if err != nil {
		t.Fatalf("Test_VirtualMachineValidations failed: Error should be nil")
	}

	// Invalid url
	httpProxy = "//skyproxy.ceccloud1.selfhost:3128"
	vm = &compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile: &compute.OSProfile{
				ProxyConfiguration: &compute.ProxyConfiguration{
					HttpProxy:  &httpProxy,
					HttpsProxy: &httpsProxy,
				},
			},
		},
	}
	err = wssdcloudclient.virtualMachineValidations(wssdcloudproto.Operation_POST, vm)
	expectedErrorMsg := "Invalid proxy URL. The URL scheme should be http or https: Invalid Input"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)

	//nil value
	vm = &compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile: &compute.OSProfile{
				ProxyConfiguration: &compute.ProxyConfiguration{
					HttpProxy:  nil,
					HttpsProxy: nil,
				},
			},
		},
	}

	err = wssdcloudclient.virtualMachineValidations(wssdcloudproto.Operation_POST, vm)
	if err != nil {
		t.Fatalf("Test_VirtualMachineValidations failed: Error should be nil for nil value")
	}

	// Invalid Http URI and nil Https
	httpProxy = "https"
	vm = &compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile: &compute.OSProfile{
				ProxyConfiguration: &compute.ProxyConfiguration{
					HttpProxy:  &httpProxy,
					HttpsProxy: nil,
				},
			},
		},
	}
	err = wssdcloudclient.virtualMachineValidations(wssdcloudproto.Operation_POST, vm)
	expectedErrorMsg = "parse \"https\": invalid URI for request: Invalid Input"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)

	// Valid Https URI and nil Http
	httpsProxy = "http://skyproxy.ceccloud1.selfhost.corp.microsoft.com:3128"
	vm = &compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile: &compute.OSProfile{
				ProxyConfiguration: &compute.ProxyConfiguration{
					HttpProxy:  nil,
					HttpsProxy: &httpsProxy,
				},
			},
		},
	}

	err = wssdcloudclient.virtualMachineValidations(wssdcloudproto.Operation_POST, vm)
	if err != nil {
		t.Fatalf("Test_VirtualMachineValidations failed: Error should be nil")
	}
}
