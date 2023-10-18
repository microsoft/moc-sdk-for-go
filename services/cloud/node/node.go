// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package node

import (
	"fmt"
	"strconv"

	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
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
		Tags:    generateNodeTags(nd),
	}
}

func getNodeStatuses(node *wssdcloud.Node) map[string]*string {
	statuses := status.GetStatuses(node.GetStatus())
	statuses["RunningState"] = convert.ToStringPtr(node.GetRunningState().String())
	statuses["Info"] = convert.ToStringPtr(node.GetInfo().String())
	return statuses
}

func generateNodeTags(node *wssdcloud.Node) map[string]*string {
	tags := make(map[string]*string)
	populateOsRegistrationStatusTag(tags, node)
	populateOsVersionTag(tags, node)

	testStr := "test"
	tags[constant.OsVersion+constant.OsVersion] = &testStr

	fmt.Println("Printing tags - ", tags)
	fmt.Printf("Printing node object - %+v", node)

	if len(tags) > 0 {
		return tags
	}
	return nil
}

func populateOsRegistrationStatusTag(tags map[string]*string, node *wssdcloud.Node) {
	if node.Info != nil && node.Info.OsInfo != nil && node.Info.OsInfo.OsRegistrationStatus != nil {
		osRegistrationStatus := strconv.Itoa(int(node.Info.OsInfo.OsRegistrationStatus.Status))
		tags[constant.OsRegistrationStatus] = &osRegistrationStatus
	}
}

func populateOsVersionTag(tags map[string]*string, node *wssdcloud.Node) {
	if node.Info != nil && node.Info.OsInfo != nil {
		tags[constant.OsVersion] = &node.Info.OsInfo.Osversion
	}
}
