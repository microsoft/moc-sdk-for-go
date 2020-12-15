// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package certificate

import (
	"github.com/microsoft/moc-sdk-for-go/services/security"
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
