// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	admin_pb "github.com/microsoft/moc/rpc/cloudagent/admin"
	cadmin_pb "github.com/microsoft/moc/rpc/common/admin"
)

// GetLogClient returns the log client to communicate with the wssdcloud agent
func GetLogClient(serverAddress *string, authorizer auth.Authorizer) (admin_pb.LogAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get LogClient. Failed to dial: %v", err)
	}

	return admin_pb.NewLogAgentClient(conn), nil
}

// GetRecoveryClient returns the log client to communicate with the wssdcloud agent
func GetRecoveryClient(serverAddress *string, authorizer auth.Authorizer) (cadmin_pb.RecoveryAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get RecoveryClient. Failed to dial: %v", err)
	}

	return cadmin_pb.NewRecoveryAgentClient(conn), nil
}

// GetDebugClient returns the log client to communicate with the wssdcloud agent
func GetDebugClient(serverAddress *string, authorizer auth.Authorizer) (cadmin_pb.DebugAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get DebugClient. Failed to dial: %v", err)
	}

	return cadmin_pb.NewDebugAgentClient(conn), nil
}

// GetVersionClient returns the wssdcloudagent version
func GetVersionClient(serverAddress *string, authorizer auth.Authorizer) (cadmin_pb.VersionAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get VersionClient. Failed to dial: %v", err)
	}

	return cadmin_pb.NewVersionAgentClient(conn), nil
}
