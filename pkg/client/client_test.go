package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/rpc/testagent"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Test_AuthenticationClientConnectionsLeak(t *testing.T) {
	ClearConnectionCache()
	server := "localhost"
	port := "9005"
	address := server + ":" + port
	tlsCert, certPem, _ := getClientCert(t)
	creds := getTlsCreds(t, []tls.Certificate{tlsCert}, [][]byte{certPem})
	grpcServer := getGrpcServer(t, creds)
	go startHelloServer(grpcServer, address)
	defer grpcServer.Stop()

	time.Sleep((time.Second * 3))

	tlsCertClient, _, _ := getSignedCert(t, tlsCert)
	authorizer, err := auth.NewAuthorizerFromInput(tlsCertClient, certPem, server)
	assert.NoErrorf(t, err, "Failed to create TLS Credentials", err)

	initialConnections, err := getActiveConnectionOnPort(port)
	assert.NoErrorf(t, err, "Failed to get number of active connections", err)
	for i := 0; i < 1000; i++ {
		_, err = getAuthConnection(&address, authorizer)
		assert.NoErrorf(t, err, "Failed to create CASignedAuth client", err)
	}

	waitTime := time.Now().Add(time.Minute * 10)
	for time.Now().Before(waitTime) {
		newConnections, err := getActiveConnectionOnPort(port)
		assert.NoErrorf(t, err, "Failed to get number of active connections", err)
		if newConnections == initialConnections {
			break
		}
		fmt.Println("Waiting for connections to close, current connections: ", newConnections)
		time.Sleep(time.Minute * 1)
	}
}

type TestTlsServer struct {
}

func (s *TestTlsServer) PingHello(ctx context.Context, in *testagent.Hello) (*testagent.Hello, error) {
	return &testagent.Hello{Name: "Hello From the Server!" + in.Name}, nil
}

func startHelloServer(grpcServer *grpc.Server, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tlsServer := TestTlsServer{}
	testagent.RegisterHelloAgentServer(grpcServer, &tlsServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func getGrpcServer(t *testing.T, creds credentials.TransportCredentials) *grpc.Server {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)
	return grpcServer
}

func getTlsCreds(t *testing.T, tlsCert []tls.Certificate, certPems [][]byte) credentials.TransportCredentials {
	certPool := x509.NewCertPool()
	for _, certPem := range certPems {
		ok := certPool.AppendCertsFromPEM(certPem)
		assert.True(t, ok, "Failed setting up cert pool")
	}
	return credentials.NewTLS(&tls.Config{
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		Certificates:             tlsCert,
		ClientCAs:                certPool,
	})
}

func getSignedCert(t *testing.T, signer tls.Certificate) (tls.Certificate, []byte, []byte) {
	caConfig := certs.CAConfig{
		RootSigner: &signer,
	}
	caAuth, err := certs.NewCertificateAuthority(&caConfig)
	assert.NoErrorf(t, err, "Error creation CA Auth: %v", err)
	conf := certs.Config{
		CommonName:   "Test Cert",
		Organization: []string{"microsoft"},
	}
	conf.AltNames.DNSNames = []string{"Test Cert"}
	csr, keyClientPem, err := certs.GenerateCertificateRequest(&conf, nil)
	assert.NoErrorf(t, err, "Error creation in CSR: %v", err)
	signConf := certs.SignConfig{Offset: time.Second * 5}
	clientCertPem, err := caAuth.SignRequest(csr, nil, &signConf)
	assert.NoErrorf(t, err, "Error signing CSR: %v", err)
	tlsCert, err := tls.X509KeyPair(clientCertPem, keyClientPem)
	return tlsCert, clientCertPem, keyClientPem
}

func getClientCert(t *testing.T) (tls.Certificate, []byte, []byte) {
	cert, key, err := certs.GenerateClientCertificate("test")
	certPem := certs.EncodeCertPEM(cert)
	keyPem := certs.EncodePrivateKeyPEM(key)
	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	assert.NoErrorf(t, err, "Failed to get tls cert", err)
	return tlsCert, certPem, keyPem
}

func getActiveConnectionOnPort(port string) (count int32, err error) {
	cmd := exec.Command("netstat", "-an", "-t")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// Parse output and count active ports
	lines := strings.Split(string(output), "\n")
	count = 0
	for _, line := range lines {
		if strings.Contains(line, ":"+port) {
			count++
		}
	}

	return
}
