// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	log "k8s.io/klog"

	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
)

const (
	debugModeTLS     = "WSSD_DEBUG_MODE"
	ServerPort   int = 55000
	AuthPort     int = 65000
)

var (
	mux             sync.Mutex
	connectionCache map[string]*grpc.ClientConn
)

func clientConnOptionsInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			grpcConnFailure := errors.IsGRPCUnavailable(err) || errors.IsGRPCDeadlineExceeded(err)
			exitProcess := !constant.GetClientOpts().NoExitOnConnFailure
			if grpcConnFailure && exitProcess {
				log.Fatalf("Communication with cloud agent failed. Exiting Process with error: %+v\n", err)
			}
		}
		return err
	}
}

func init() {
	connectionCache = map[string]*grpc.ClientConn{}
}

func ClearConnectionCache() {
	connectionCache = map[string]*grpc.ClientConn{}
}

// Returns nil if debug mode is on; err if it is not
func isDebugMode() error {
	debugEnv := strings.ToLower(os.Getenv(debugModeTLS))
	if debugEnv == "on" {
		return nil
	}
	return fmt.Errorf("Debug Mode not set")
}

func getServerEndpoint(serverAddress *string) string {
	if strings.Contains(*serverAddress, ":") {
		return *serverAddress
	}
	return fmt.Sprintf("%s:%d", *serverAddress, ServerPort)
}

func getAuthServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, AuthPort)
}

func getDefaultDialOption(authorizer auth.Authorizer) []grpc.DialOption {
	var opts []grpc.DialOption

	// Debug Mode allows us to talk to wssdagent without a proper handshake
	// This means we can debug and test wssdagent without generating certs
	// and having proper tokens

	// Check if debug mode is on
	if ok := isDebugMode(); ok == nil {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	}

	opts = append(opts, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}))

	return opts
}

func isValidConnections(conn *grpc.ClientConn) bool {

	switch conn.GetState() {
	case connectivity.TransientFailure:
		fallthrough
	case connectivity.Shutdown:
		return false
	default:
		return true
	}
}

func getClientConnection(serverAddress *string, authorizer auth.Authorizer) (*grpc.ClientConn, error) {
	mux.Lock()
	defer mux.Unlock()
	endpoint := getServerEndpoint(serverAddress)

	conn, ok := connectionCache[endpoint]
	if ok {
		if isValidConnections(conn) {
			return conn, nil
		}
		conn.Close()
	}

	opts := getDefaultDialOption(authorizer)
	opts = append(opts, grpc.WithUnaryInterceptor(clientConnOptionsInterceptor()))
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	connectionCache[endpoint] = conn

	return conn, nil
}

func getAuthConnection(serverAddress *string, authorizer auth.Authorizer) (*grpc.ClientConn, error) {
	mux.Lock()
	defer mux.Unlock()
	endpoint := getAuthServerEndpoint(serverAddress)

	conn, ok := connectionCache[endpoint]
	if ok {
		return conn, nil
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	opts = append(opts, grpc.WithPerRPCCredentials(authorizer.WithRPCAuthorization()))

	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	connectionCache[endpoint] = conn

	return conn, nil
}
