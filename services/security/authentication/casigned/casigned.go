// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package casigned

import (
	"context"
	"fmt"
	"time"

	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	//log "k8s.io/klog"
)

type client struct {
	wssdsecurity.AuthenticationAgentClient
	cloudFQDN string
}

// NewAuthenticationClient creates a client session with the backend wssd agent
func NewAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetAuthenticationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c, subID}, nil
}

// Login
func (c *client) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	request := getAuthenticationRequest(identity)
	response, err := c.AuthenticationAgentClient.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return &response.Token, nil
}

func RenewRoutine(group, server string) {
	renewalAttempt := 0
	// Waiting for a few seconds to avoid spamming short-lived sdk user
	time.Sleep(time.Second * 5)
	for {
		wssdConfig := auth.WssdConfig{}
		err := marshal.FromJSONFile(auth.GetWssdConfigLocation(), &wssdConfig)
		if err != nil {
			fmt.Printf("Failed to open config file in location %s: %v\n", auth.GetWssdConfigLocation(), err)
			panic("Failed to open config file")
		}

		sleepTime, err := renewTime(wssdConfig.ClientCertificate)
		if err != nil {
			fmt.Printf("Failed to find sleep time for cert %v \n", err)
			panic("Failed to find renew time for certificate")
		}
		fmt.Printf("Waiting for %v to renew cert\n", sleepTime)
		time.Sleep(sleepTime)
		fmt.Printf("Attempting to renew certificate\n")
		err = auth.RenewCertificates(server, auth.GetWssdConfigLocation())
		if err != nil {
			if renewalAttempt > 5 {
				fmt.Printf("Failed to renew cert to %v \n", err)
				panic("Failed to renew certificate")
			}
			if errors.IsExpired(err) {
				fmt.Printf("Certificate Expired %v \n", err)
				panic("Certificate has expired")
			}
			renewalAttempt += 1
			// Wait for 30 seconds before next renewal attempt
			time.Sleep(time.Second * 30)
			continue
		}
		//reset renewalAttempt after successful renewal
		renewalAttempt = 0
		fmt.Printf("Certificate renewal complete\n")
	}
}

// Get methods invokes the client Get method
func (c *client) LoginWithConfig(group string, loginconfig auth.LoginConfig) (*auth.WssdConfig, error) {

	clientCsr, accessFile, err := auth.GenerateClientCsr(loginconfig)
	if err != nil {
		return nil, err
	}

	id := security.Identity{
		Name:        &loginconfig.Name,
		Certificate: &clientCsr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	clientCert, err := c.Login(ctx, group, &id)
	if err != nil {
		return nil, err
	}
	accessFile.ClientCertificate = *clientCert
	accessFile.ClientCertificateType = auth.CASigned
	accessFile.IdentityName = loginconfig.Name
	auth.PrintAccessFile(accessFile)
	go RenewRoutine(group, c.cloudFQDN)
	return &accessFile, nil
}

func calculateTime(before, after, now time.Time) time.Duration {
	validity := after.Sub(before)
	// Threshold to renew is 30% of validity
	tresh := time.Duration(float64(validity.Nanoseconds()) * constant.CertificateValidityThreshold)

	treshNotAfter := after.Add(-tresh)
	return treshNotAfter.Sub(now)
}

func renewTime(certificate string) (duration time.Duration, err error) {

	pemCert, err := marshal.FromBase64(certificate)
	if err != nil {
		return
	}

	x509Cert, err := certs.DecodeCertPEM([]byte(pemCert))
	if err != nil {
		return
	}
	return calculateTime(x509Cert.NotBefore, x509Cert.NotAfter, time.Now()), nil
}

func getAuthenticationRequest(identity *security.Identity) *wssdsecurity.AuthenticationRequest {
	certs := map[string]string{"": *identity.Certificate}
	request := &wssdsecurity.AuthenticationRequest{
		Identity: &wssdsecurity.Identity{
			Name:         *identity.Name,
			Certificates: certs,
		},
	}
	return request
}
