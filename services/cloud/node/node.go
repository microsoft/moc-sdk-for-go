// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	"github.com/microsoft/moc/pkg/convert"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdNode(nd *cloud.Node, location string) (*wssdcloud.Node, error) {
	if nd.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	if nd.FQDN == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing FQDN in Configuration")
	}

	if nd.Port == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Port in Configuration")
	}

	if nd.AuthorizerPort == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing AuthorizrPort in Configuration")
	}

	node := &wssdcloud.Node{
		Name:           *nd.Name,
		Fqdn:           *nd.FQDN,
		LocationName:   location,
		Port:           *nd.Port,
		AuthorizerPort: *nd.AuthorizerPort,
	}

	if nd.Version != nil {
		if node.Status == nil {
			node.Status = status.InitStatus()
		}
		node.Status.Version.Number = *nd.Version
	}

	if nd.Certificate != nil {
		node.Certificate = *nd.Certificate
	}

	return node, nil
}

// Conversion functions from wssdcloud to cloud
func getNode(nd *wssdcloud.Node) *cloud.Node {
	return &cloud.Node{
		Name:     &nd.Name,
		Location: &nd.LocationName,
		NodeProperties: &cloud.NodeProperties{
			FQDN:           &nd.Fqdn,
			Port:           &nd.Port,
			AuthorizerPort: &nd.AuthorizerPort,
			Certificate:    &nd.Certificate,
			Statuses:       getNodeStatuses(nd),
		},
		Version: &nd.Status.Version.Number,
		Tags:    getNodeTags(nd),
	}
}

func getNodeStatuses(node *wssdcloud.Node) map[string]*string {
	statuses := status.GetStatuses(node.GetStatus())
	statuses["RunningState"] = convert.ToStringPtr(node.GetRunningState().String())
	statuses["Info"] = convert.ToStringPtr(node.GetInfo().String())
	return statuses
}

func getNodeTags(node *wssdcloud.Node) map[string]*string {
	tags := make(map[string]*string)
	if node.Info != nil {
		if node.Info.Capability != nil {
			if node.Info.Capability.OsInfo != nil {
				registrationStatus := string(int32(node.Info.Capability.OsInfo.OsRegistrationStatus))
				tags["registrationStatus"] = &registrationStatus
			}
		}
	}
	return tags
}
