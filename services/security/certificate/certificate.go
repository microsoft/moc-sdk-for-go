// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package certificate

import (
	"net"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
)

func GetCertificate(cert *wssdcloudsecurity.Certificate) *security.Certificate {
	certificateType := cert.Type.String()
	return &security.Certificate{
		ID:      &cert.Id,
		Name:    &cert.Name,
		Cer:     &cert.Certificate,
		Type:    &certificateType,
		Version: &cert.Status.Version.Number,
		Attributes: &security.CertificateAttributes{
			NotBefore: &cert.NotBefore,
			Expires:   &cert.NotAfter,
			Statuses:  status.GetStatuses(cert.GetStatus()),
		},
	}
}

func GetCertificateType(certType string) (wssdcloudsecurity.CertificateType, bool) {
	value, ok := wssdcloudsecurity.CertificateType_value[certType]
	return wssdcloudsecurity.CertificateType(value), ok
}

func GetWssdCertificate(cert *security.Certificate) (*wssdcloudsecurity.Certificate, error) {
	if cert.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Certificate name is missing")
	}
	certType, ok := GetCertificateType(*cert.Type)
	if !ok {
		return nil, errors.Wrapf(errors.InvalidInput, "Invalid certificate type %s", *cert.Type)
	}
	certificate := &wssdcloudsecurity.Certificate{
		Name: *cert.Name,
		Type: certType,
	}

	if cert.Version != nil {
		if certificate.Status == nil {
			certificate.Status = status.InitStatus()
		}
		certificate.Status.Version.Number = *cert.Version
	}
	return certificate, nil
}

func GetMocCSR(csr *security.CertificateRequest) (*wssdcloudsecurity.CertificateSigningRequest, string, error) {
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
		pemKey, err := marshal.FromBase64(*csr.PrivateKey)
		if err != nil {
			return nil, "", err
		}
		csrRequest, key, err = certs.GenerateCertificateRequest(&conf, pemKey)
	} else {
		csrRequest, key, err = certs.GenerateCertificateRequest(&conf, nil)
	}
	if err != nil {
		return nil, "", errors.Wrapf(errors.Failed, "Failed creating certificate Request")
	}
	request := &wssdcloudsecurity.CertificateSigningRequest{
		Name: *csr.Name,
		Csr:  string(csrRequest),
	}
	if csr.OldCertificate != nil {
		request.OldCertificate = *csr.OldCertificate
	}
	if csr.CaName != nil {
		request.CaName = *csr.CaName
	}
	if csr.IsCA != nil {
		request.IsCA = &wrappers.BoolValue{Value: *csr.IsCA}
	}
	return request, string(key), nil
}
