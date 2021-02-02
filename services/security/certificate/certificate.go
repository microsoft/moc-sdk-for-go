// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package certificate

import (
	"net"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
)

func getCertificate(cert *wssdcloudsecurity.Certificate) *security.Certificate {
	return &security.Certificate{
		ID:   &cert.Id,
		Name: &cert.Name,
		Cer:  &cert.NewCertificate,
		Attributes: &security.CertificateAttributes{
			NotBefore: &cert.NotBefore,
			Expires:   &cert.NotAfter,
			Statuses:  status.GetStatuses(cert.GetStatus()),
		},
	}
}

func getWssdCertificate(cert *security.Certificate) (*wssdcloudsecurity.Certificate, error) {
	if cert.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Certificate name is missing")
	}
	return &wssdcloudsecurity.Certificate{
		Name: *cert.Name,
	}, nil
}

func getWssdCSR(csr *security.CertificateRequest) (*wssdcloudsecurity.CertificateSigningRequest, string, error) {
	if csr.Name == nil {
		return nil, "", errors.Wrapf(errors.InvalidInput, "CSR name is missing")
	}
	conf := certs.Config{
		CommonName: *csr.Name,
	}
	if csr.Attributes != nil {
		conf.AltNames.DNSNames = *csr.Attributes.DNSNames
		for _, ipStr := range *csr.Attributes.IPs {
			ip, _, err := net.ParseCIDR(ipStr)
			if err != nil {
				return nil, "", errors.Wrapf(errors.InvalidInput, "Invalid Ipaddress %s", ipStr)
			}
			conf.AltNames.IPs = append(conf.AltNames.IPs, ip)
		}
	}
	var key []byte
	var csrRequest []byte
	var err error
	if csr.PrivateKey != nil {
		csrRequest, key, err = certs.GenerateCertificateRequest(&conf, []byte(*csr.PrivateKey))
	} else {
		csrRequest, key, err = certs.GenerateCertificateRequest(&conf, nil)
	}
	if err != nil {
		return nil, "", errors.Wrapf(errors.Failed, "Failed creating certificate Request")
	}
	request := &wssdcloudsecurity.CertificateSigningRequest{
		Name:  *csr.Name,
		Csr:   string(csrRequest),
		Renew: csr.Renew,
	}
	if csr.OldCertificate != nil {
		request.OldCertificate = *csr.OldCertificate
	}
	if csr.CaName != nil {
		request.CaName = *csr.CaName
	}
	return request, string(key), nil
}
