// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package controlplane

import (
	"github.com/microsoft/moc-sdk-for-go/services/cloud"

	"github.com/microsoft/moc/pkg/convert"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloud "github.com/microsoft/moc/rpc/cloudagent/cloud"
)

// Conversion functions from cloud to wssdcloud
func getWssdControlPlane(cp *cloud.ControlPlaneInfo, location string) (*wssdcloud.ControlPlane, error) {
	if cp.Name == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Name in Configuration")
	}

	if cp.Fqdn == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing FQDN in Configuration")
	}

	if cp.Port == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Port in Configuration")
	}

	controlPlane := &wssdcloud.ControlPlane{
		Name:         *cp.Name,
		Fqdn:         *cp.Fqdn,
		LocationName: location,
		Port:         *cp.Port,
	}

	if cp.Version != nil {
		if controlPlane.Status == nil {
			controlPlane.Status = status.InitStatus()
		}
		controlPlane.Status.Version.Number = *cp.Version
	}

	return controlPlane, nil
}

// Conversion functions from wssdcloud to cloud
func getControlPlane(cp *wssdcloud.ControlPlane) *cloud.ControlPlaneInfo {
	return &cloud.ControlPlaneInfo{
		Name:     &cp.Name,
		Location: &cp.LocationName,
		ControlPlaneProperties: &cloud.ControlPlaneProperties{
			Fqdn:     &cp.Fqdn,
			Port:     &cp.Port,
			Statuses: getControlPlaneStatuses(cp),
		},
		Version: &cp.Status.Version.Number,
	}
}

func getControlPlaneStatuses(cp *wssdcloud.ControlPlane) map[string]*string {
	statuses := status.GetStatuses(cp.GetStatus())
	statuses["State"] = convert.ToStringPtr(cp.GetState().String())
	return statuses
}
