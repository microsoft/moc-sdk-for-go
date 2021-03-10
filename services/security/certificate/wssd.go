// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package certificate

import (
	"context"
	"fmt"

	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	log "k8s.io/klog"
)

type client struct {
	wssdcloudsecurity.CertificateAgentClient
}

// NewCertificateClientN- creates a client session with the backend wssdcloud agent
func newCertificateClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetCertificateClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.Certificate, error) {
	request, err := getCertificateRequest(name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.CertificateAgentClient.Get(ctx, request)
	if err != nil {
		return nil, err
	}
	return getCertificatesFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.Certificate) (*security.Certificate, error) {
	request, err := getCertificateRequest(name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.CertificateAgentClient.CreateOrUpdate(ctx, request)
	if err != nil {
		log.Errorf("[Certificate] Create failed with error %v", err)
		return nil, err
	}

	cert := getCertificatesFromResponse(response)

	if len(*cert) == 0 {
		return nil, fmt.Errorf("[Certificate][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*cert)[0]), err
}

// Sign
func (c *client) Sign(ctx context.Context, group, name string, csr *security.CertificateRequest) (*security.Certificate, string, error) {
	csr.OldCertificate = nil
	request, key, err := getCSRRequest(name, csr)
	if err != nil {
		return nil, "", err
	}
	response, err := c.CertificateAgentClient.Sign(ctx, request)
	if err != nil {
		log.Errorf("[Certificate] Create failed with error %v", err)
		return nil, "", err
	}

	cert := getCertificatesFromResponse(response)

	if len(*cert) == 0 {
		return nil, "", fmt.Errorf("[Certificate][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*cert)[0]), string(key), err
}

// CreateOrUpdate
func (c *client) Renew(ctx context.Context, group, name string, csr *security.CertificateRequest) (*security.Certificate, string, error) {
	if csr.OldCertificate == nil || len(*csr.OldCertificate) == 0 {
		return nil, "", errors.Wrapf(errors.NotFound, "[Certificate] Renew missing oldCert field")
	}

	request, key, err := getCSRRequest(name, csr)
	if err != nil {
		return nil, "", err
	}
	response, err := c.CertificateAgentClient.Renew(ctx, request)
	if err != nil {
		log.Errorf("[Certificate] Create failed with error %v", err)
		return nil, "", err
	}

	cert := getCertificatesFromResponse(response)

	if len(*cert) == 0 {
		return nil, "", fmt.Errorf("[Certificate][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*cert)[0]), string(key), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	cert, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*cert) == 0 {
		return fmt.Errorf("Certificate [%s] not found", name)
	}

	request, err := getCertificateRequest(name, &(*cert)[0])
	if err != nil {
		return err
	}
	_, err = c.CertificateAgentClient.Delete(ctx, request)
	return err
}

func getCertificatesFromResponse(response *wssdcloudsecurity.CertificateResponse) *[]security.Certificate {
	certs := []security.Certificate{}
	for _, certificates := range response.GetCertificates() {
		certs = append(certs, *(GetCertificate(certificates)))
	}

	return &certs
}

func getCertificateRequest(name string, cert *security.Certificate) (*wssdcloudsecurity.CertificateRequest, error) {
	request := &wssdcloudsecurity.CertificateRequest{
		Certificates: []*wssdcloudsecurity.Certificate{},
	}
	wssdcertificate := &wssdcloudsecurity.Certificate{
		Name: name,
	}

	var err error
	if cert != nil {
		wssdcertificate, err = GetWssdCertificate(cert)
		if err != nil {
			return nil, err
		}
	}
	request.Certificates = append(request.Certificates, wssdcertificate)
	return request, nil
}

func getCSRRequest(name string, csr *security.CertificateRequest) (*wssdcloudsecurity.CSRRequest, string, error) {
	request := &wssdcloudsecurity.CSRRequest{
		CSRs: []*wssdcloudsecurity.CertificateSigningRequest{},
	}
	wssdcsr := &wssdcloudsecurity.CertificateSigningRequest{
		Name: name,
	}

	var err error
	var key string
	if csr != nil {
		wssdcsr, key, err = GetMocCSR(csr)
		if err != nil {
			return nil, "", err
		}
	}
	request.CSRs = append(request.CSRs, wssdcsr)
	return request, key, nil
}
