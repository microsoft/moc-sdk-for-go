// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package moc

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc-sdk-for-go/services/security/certificate"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	//log "k8s.io/klog"
)

type client struct {
	wssdsecurity.RenewalAgentClient
}

// NewRenewClient creates a client session with the backend wssd agent
func NewRenewClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetRenewClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func fromBase64(cert, key string) (pemCert, pemKey []byte, err error) {
	pemCert, err = marshal.FromBase64(cert)
	if err != nil {
		return
	}
	pemKey, err = marshal.FromBase64(key)
	if err != nil {
		return
	}
	return
}

func renewRequired(x509Cert *x509.Certificate) bool {
	validity := x509Cert.NotAfter.Sub(x509Cert.NotBefore)

	// Threshold to renew is 30% of validity
	tresh := time.Duration(float64(validity.Nanoseconds()) * constant.CertificateValidityThreshold)

	treshNotAfter := time.Now().Add(tresh)
	if x509Cert.NotAfter.After(treshNotAfter) {
		return false
	}
	return true
}

func (c *client) RenewConfig(ctx context.Context, wssdConfig *auth.WssdConfig) (newConfig *auth.WssdConfig, renewed bool, err error) {
	renewed = false
	pemCert, _, err := fromBase64(wssdConfig.ClientCertificate, wssdConfig.ClientKey)
	if err != nil {
		return
	}

	x509Cert, err := certs.DecodeCertPEM([]byte(pemCert))
	if err != nil {
		return
	}

	if !renewRequired(x509Cert) {
		return wssdConfig, renewed, nil
	}

	csr := &security.CertificateRequest{
		Name:           &x509Cert.Issuer.CommonName,
		OldCertificate: &wssdConfig.ClientCertificate,
	}

	newCert, newKey, err := c.Renew(ctx, "", csr)
	if err != nil {
		if errors.IsNotSupported(err) {
			return wssdConfig, renewed, nil
		}
		return
	}
	fmt.Println("Renewd babyyyy")
	newConfig = &auth.WssdConfig{
		CloudCertificate:  wssdConfig.CloudCertificate,
		ClientCertificate: marshal.ToBase64(*newCert.Cer),
		ClientKey:         marshal.ToBase64(string(newKey)),
	}
	renewed = true
	return
}

func (c *client) Renew(ctx context.Context, group string, csr *security.CertificateRequest) (*security.Certificate, string, error) {

	if csr.OldCertificate == nil || len(*csr.OldCertificate) == 0 {
		return nil, "", errors.Wrapf(errors.NotFound, "[Certificate] Renew missing oldCert field")
	}

	request, key, err := getCSRRequest(*csr.Name, csr)
	if err != nil {
		return nil, "", err
	}
	response, err := c.RenewalAgentClient.RenewCertificate(ctx, request)
	if err != nil {
		return nil, "", errors.Wrapf(err, "[Certificate] Create failed with error")
	}
	return getCertificatesFromResponse(response), string(key), err
}

func getCertificatesFromResponse(response *wssdsecurity.RenewResponse) *security.Certificate {
	return certificate.GetCertificate(response.GetCertificate())
}

func getCSRRequest(name string, csr *security.CertificateRequest) (*wssdsecurity.RenewRequest, string, error) {
	wssdcsr := &wssdsecurity.CertificateSigningRequest{
		Name: name,
	}

	var err error
	var key string
	if csr != nil {
		wssdcsr, key, err = certificate.GetMocCSR(csr)
		if err != nil {
			return nil, "", err
		}
	}
	request := &wssdsecurity.RenewRequest{
		CSR: wssdcsr,
	}
	return request, key, nil
}
