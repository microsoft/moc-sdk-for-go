// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	admin_pb "github.com/microsoft/moc/rpc/cloudagent/admin"
)

// GetLogClient returns the log client to communicate with the wssdcloud agent
func GetLogClient(serverAddress *string, authorizer auth.Authorizer) (admin_pb.LogAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get LogClient. Failed to dial: %v", err)
	}

	return admin_pb.NewLogAgentClient(conn), nil
}
