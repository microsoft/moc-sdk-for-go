// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualmachine

import (
	"context"
	"fmt"

	"github.com/microsoft/moc-sdk-for-go/services/compute"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"
	prototags "github.com/microsoft/moc/pkg/tags"
	"github.com/microsoft/moc/pkg/validations"
	wssdcloudproto "github.com/microsoft/moc/rpc/common"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	wssdcloudcompute "github.com/microsoft/moc/rpc/cloudagent/compute"
)

type client struct {
	wssdcloudcompute.VirtualMachineAgentClient
}

// newVirtualMachineClient - creates a client session with the backend wssdcloud agent
func newVirtualMachineClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetVirtualMachineClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineFromResponse(response, group), nil
}

// Get
func (c *client) get(ctx context.Context, group, name string) ([]*wssdcloudcompute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_GET, group, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetVirtualMachines(), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_POST, group, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vms := c.getVirtualMachineFromResponse(response, group)
	if len(*vms) == 0 {
		return nil, fmt.Errorf("Creation of Virtual Machine failed to unknown reason.")
	}

	return &(*vms)[0], nil
}

// Hydrate
func (c *client) Hydrate(ctx context.Context, group, name string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_HYDRATE, group, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vms := c.getVirtualMachineFromResponse(response, group)
	if len(*vms) == 0 {
		return nil, fmt.Errorf("hydration of Virtual Machine failed to unknown reason")
	}

	return &(*vms)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vm, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vm) == 0 {
		return fmt.Errorf("Virtual Machine [%s] not found", name)
	}

	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_DELETE, group, name, &(*vm)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineAgentClient.Invoke(ctx, request)

	return err
}

// Query
func (c *client) Query(ctx context.Context, group, query string) (*[]compute.VirtualMachine, error) {
	vms, err := c.Get(ctx, group, "")
	if err != nil {
		return nil, err
	}

	filteredBytes, err := config.MarshalOutput(*vms, query, "json")
	if err != nil {
		return nil, err
	}

	err = marshal.FromJSONBytes(filteredBytes, vms)
	if err != nil {
		return nil, err
	}

	return vms, nil
}

// Stop
func (c *client) Stop(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Stop, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// Start
func (c *client) Start(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Start, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// Pause
func (c *client) Pause(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Pause, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// Save
func (c *client) Save(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Save, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// RemoveIso
func (c *client) RemoveIsoDisk(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Remove_Iso_Disk, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// RepairGuestAgent
func (c *client) RepairGuestAgent(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcloudproto.ProviderAccessOperation_VirtualMachine_Repair_Guest_Agent, group, name)
	if err != nil {
		return
	}

	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

// RunCommand
func (c *client) RunCommand(ctx context.Context, group, name string, request *compute.VirtualMachineRunCommandRequest) (response *compute.VirtualMachineRunCommandResponse, err error) {
	mocRequest, err := c.getVirtualMachineRunCommandRequest(ctx, group, name, request)
	if err != nil {
		return
	}

	mocResponse, err := c.VirtualMachineAgentClient.RunCommand(ctx, mocRequest)
	if err != nil {
		return
	}
	response, err = c.getVirtualMachineRunCommandResponse(mocResponse)
	return
}

// Get
func (c *client) Validate(ctx context.Context, group, name string) error {
	request, err := c.getVirtualMachineRequest(wssdcloudproto.Operation_VALIDATE, group, name, nil)
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Precheck(ctx context.Context, group string, vms []*compute.VirtualMachine) (bool, error) {
	request, err := c.getVirtualMachinePrecheckRequest(group, vms)
	if err != nil {
		return false, err
	}
	response, err := c.VirtualMachineAgentClient.Precheck(ctx, request)
	if err != nil {
		return false, err
	}
	return getVirtualMachinePrecheckResponse(response)
}

func getVirtualMachinePrecheckResponse(response *wssdcloudcompute.VirtualMachinePrecheckResponse) (bool, error) {
	var err error = nil
	result := response.GetResult().GetValue()
	if !result {
		err = errors.New(response.GetError())
	}
	return result, err
}

func (c *client) getVirtualMachinePrecheckRequest(group string, vms []*compute.VirtualMachine) (*wssdcloudcompute.VirtualMachinePrecheckRequest, error) {
	request := &wssdcloudcompute.VirtualMachinePrecheckRequest{
		VirtualMachines: []*wssdcloudcompute.VirtualMachine{},
	}
	for _, vm := range vms {
		if vm != nil {
			mocvm, err := c.getWssdVirtualMachine(vm, group)
			if err != nil {
				return nil, err
			}
			request.VirtualMachines = append(request.VirtualMachines, mocvm)
		}
	}
	return request, nil
}

// Private methods
func (c *client) getVirtualMachineRunCommandRequest(ctx context.Context, group, name string, request *compute.VirtualMachineRunCommandRequest) (mocRequest *wssdcloudcompute.VirtualMachineRunCommandRequest, err error) {
	vms, err := c.get(ctx, group, name)
	if err != nil {
		return
	}

	if len(vms) != 1 {
		err = errors.Wrapf(errors.InvalidInput, "Multiple Virtual Machines found in group %s with name %s", group, name)
		return
	}
	vm := vms[0]

	var params []*wssdcloudproto.VirtualMachineRunCommandInputParameter
	if request.Parameters != nil {
		params = make([]*wssdcloudproto.VirtualMachineRunCommandInputParameter, len(*request.Parameters))
		for i, param := range *request.Parameters {
			tmp := &wssdcloudproto.VirtualMachineRunCommandInputParameter{
				Name:  *param.Name,
				Value: *param.Value,
			}
			params[i] = tmp
		}
	}

	var scriptSource wssdcloudproto.VirtualMachineRunCommandScriptSource
	if request.Source.Script != nil {
		scriptSource.Script = *request.Source.Script
	}
	if request.Source.ScriptURI != nil {
		scriptSource.ScriptURI = *request.Source.ScriptURI
	}
	if request.Source.CommandID != nil {
		scriptSource.CommandID = *request.Source.CommandID
	}

	mocRequest = &wssdcloudcompute.VirtualMachineRunCommandRequest{
		VirtualMachine:            vm,
		RunCommandInputParameters: params,
		Source:                    &scriptSource,
	}

	if request.RunAsUser != nil {
		mocRequest.RunAsUser = *request.RunAsUser
	}
	if request.RunAsPassword != nil {
		mocRequest.RunAsPassword = *request.RunAsPassword
	}
	return
}

func (c *client) getVirtualMachineRunCommandResponse(mocResponse *wssdcloudcompute.VirtualMachineRunCommandResponse) (*compute.VirtualMachineRunCommandResponse, error) {
	var executionState compute.ExecutionState
	switch mocResponse.GetInstanceView().ExecutionState {
	case wssdcloudproto.VirtualMachineRunCommandExecutionState_ExecutionState_UNKNOWN:
		executionState = compute.ExecutionStateUnknown
	case wssdcloudproto.VirtualMachineRunCommandExecutionState_ExecutionState_SUCCEEDED:
		executionState = compute.ExecutionStateSucceeded
	case wssdcloudproto.VirtualMachineRunCommandExecutionState_ExecutionState_FAILED:
		executionState = compute.ExecutionStateFailed
	default:
		return nil, errors.Wrapf(errors.NotSupported, "Unknown execution state reported for virtual machine run command")
	}

	instanceView := &compute.VirtualMachineRunCommandInstanceView{
		ExecutionState: executionState,
		ExitCode:       &mocResponse.GetInstanceView().ExitCode,
		Output:         &mocResponse.GetInstanceView().Output,
		Error:          &mocResponse.GetInstanceView().Error,
	}

	response := &compute.VirtualMachineRunCommandResponse{
		InstanceView: instanceView,
	}
	return response, nil
}

func (c *client) getVirtualMachineFromResponse(response *wssdcloudcompute.VirtualMachineResponse, group string) *[]compute.VirtualMachine {
	vms := []compute.VirtualMachine{}
	for _, vm := range response.GetVirtualMachines() {
		vms = append(vms, *(c.getVirtualMachine(vm)))
	}

	return &vms
}

func (c *client) getVirtualMachineRequest(opType wssdcloudproto.Operation, group, name string, vmss *compute.VirtualMachine) (*wssdcloudcompute.VirtualMachineRequest, error) {
	request := &wssdcloudcompute.VirtualMachineRequest{
		OperationType:   opType,
		VirtualMachines: []*wssdcloudcompute.VirtualMachine{},
	}
	var err error
	wssdvm := &wssdcloudcompute.VirtualMachine{
		Name:      name,
		GroupName: group,
	}
	if vmss != nil {
		err = c.virtualMachineValidations(opType, vmss)
		if err != nil {
			return nil, err
		}
		wssdvm, err = c.getWssdVirtualMachine(vmss, group)
		if err != nil {
			return nil, err
		}
	}
	request.VirtualMachines = append(request.VirtualMachines, wssdvm)
	return request, nil
}

func (c *client) getVirtualMachineOperationRequest(ctx context.Context,
	opType wssdcloudproto.ProviderAccessOperation,
	group, name string) (request *wssdcloudcompute.VirtualMachineOperationRequest, err error) {

	vms, err := c.get(ctx, group, name)
	if err != nil {
		return
	}

	request = &wssdcloudcompute.VirtualMachineOperationRequest{
		OperationType:   opType,
		VirtualMachines: vms,
	}
	return
}

func getComputeTags(tags *wssdcloudproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcloudproto.Tags {
	return prototags.MapToProto(tags)
}

func (c *client) virtualMachineValidations(opType wssdcloudproto.Operation, vmss *compute.VirtualMachine) error {
	if vmss.OsProfile == nil {
		return nil
	}
	if vmss.OsProfile.ProxyConfiguration != nil && opType == wssdcloudproto.Operation_POST {
		if vmss.OsProfile.ProxyConfiguration.HttpProxy != nil && *vmss.OsProfile.ProxyConfiguration.HttpProxy != "" {
			_, err := validations.ValidateProxyURL(*vmss.OsProfile.ProxyConfiguration.HttpProxy)
			if err != nil {
				return err
			}
		}
		if vmss.OsProfile.ProxyConfiguration.HttpsProxy != nil && *vmss.OsProfile.ProxyConfiguration.HttpsProxy != "" {
			_, err := validations.ValidateProxyURL(*vmss.OsProfile.ProxyConfiguration.HttpsProxy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *client) GetHyperVVmId(ctx context.Context, group, name string) (*compute.VirtualMachineHyperVVmId, error) {
	vm, err := c.get(ctx, group, name)
	if err != nil {
		return nil, err
	}

	mocResponse, err := c.VirtualMachineAgentClient.GetHyperVVmId(ctx, vm[0])
	if err != nil {
		return nil, err
	}
	response := &compute.VirtualMachineHyperVVmId{
		HyperVVmId: &mocResponse.HyperVVmId,
	}

	return response, nil
}

func (c *client) GetHostNodeName(ctx context.Context, group, name string) (*compute.VirtualMachineHostNodeName, error) {
	vm, err := c.get(ctx, group, name)
	if err != nil {
		return nil, err
	}

	mocResponse, err := c.VirtualMachineAgentClient.GetHostNodeName(ctx, vm[0])
	if err != nil {
		return nil, err
	}
	response := &compute.VirtualMachineHostNodeName{
		HostNodeName: &mocResponse.HostNodeName,
	}

	return response, nil
}
