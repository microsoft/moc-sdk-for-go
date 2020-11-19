// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package controlplane

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdControlPlane(nd *cloud.ControlPlaneInfo, location string) (*wssdcloud.ControlPlane, error) {
	if nd.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	if nd.Fqdn == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing FQDN in Configuration")
	}

	if nd.Port == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Port in Configuration")
	}

	controlPlane := &wssdcloud.ControlPlane{
		Name:         *nd.Name,
		Fqdn:         *nd.Fqdn,
		LocationName: location,
		Port:         *nd.Port,
	}

	if nd.Version != nil {
		if controlPlane.Status == nil {
			controlPlane.Status = status.InitStatus()
		}
		controlPlane.Status.Version.Number = *nd.Version
	}

	return controlPlane, nil
}

// Conversion functions from wssdcloud to cloud
func getControlPlane(nd *wssdcloud.ControlPlane) *cloud.ControlPlaneInfo {
	return &cloud.ControlPlaneInfo{
		Name:     &nd.Name,
		Location: &nd.LocationName,
		ControlPlaneProperties: &cloud.ControlPlaneProperties{
			Fqdn:     &nd.Fqdn,
			Port:     &nd.Port,
			Statuses: getControlPlaneStatuses(nd),
		},
		Version: &nd.Status.Version.Number,
	}
}

func getControlPlaneStatuses(controlPlane *wssdcloud.ControlPlane) map[string]*string {
	statuses := status.GetStatuses(controlPlane.GetStatus())
	return statuses
}
