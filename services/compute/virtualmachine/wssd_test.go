package virtualmachine

import (
	"testing"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"
)

func Test_VirtualMachineValidations(t *testing.T) {
	wssdcloudclient := client{}
	httpProxy := ""
	httpsProxy := ""

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
}
