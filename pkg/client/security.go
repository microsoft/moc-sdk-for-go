// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	"google.golang.org/grpc"
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
	security_pb "github.com/microsoft/moc/rpc/cloudagent/security"
)

// GetKeyVaultClient returns the keyvault client to communicate with the wssdagent
func GetKeyVaultClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.KeyVaultAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get KeyVaultClient. Failed to dial: %v", err)
	}

	return security_pb.NewKeyVaultAgentClient(conn), nil
}

// GetSecretClient returns the secret client to communicate with the wssdagent
func GetSecretClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.SecretAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get SecretClient. Failed to dial: %v", err)
	}

	return security_pb.NewSecretAgentClient(conn), nil
}

// GetKeyClient returns the secret client to communicate with the wssdagent
func GetKeyClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.KeyAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get KeyClient. Failed to dial: %v", err)
	}

	return security_pb.NewKeyAgentClient(conn), nil
}

// GetCertificateClient returns the secret client to communicate with the wssdagent
func GetCertificateClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.CertificateAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get CertificateClient. Failed to dial: %v", err)
	}

	return security_pb.NewCertificateAgentClient(conn), nil
}

// GetIdentityClient returns the secret client to communicate with the wssdagent
func GetIdentityClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.IdentityAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get IdentityClient. Failed to dial: %v", err)
	}

	return security_pb.NewIdentityAgentClient(conn), nil
}

// GetRoleClient returns the role client to communicate with the wssdagent
func GetRoleClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.RoleAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get RoleClient. Failed to dial: %v", err)
	}

	return security_pb.NewRoleAgentClient(conn), nil
}

// GetRoleAssignmentClient returns the roleAssignment client to communicate with the wssdagent
func GetRoleAssignmentClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.RoleAssignmentAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get RoleAssignmentClient. Failed to dial: %v", err)
	}

	return security_pb.NewRoleAssignmentAgentClient(conn), nil
}

// GetAuthenticationClient returns the secret client to communicate with the wssdagent
func GetAuthenticationClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.AuthenticationAgentClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	opts = append(opts, grpc.WithPerRPCCredentials(authorizer.WithRPCAuthorization()))

	conn, err := grpc.Dial(getAuthServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get AuthenticationClient. Failed to dial: %v", err)
	}

	return security_pb.NewAuthenticationAgentClient(conn), nil
}
