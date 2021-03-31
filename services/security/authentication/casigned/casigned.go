// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package casigned

import (
	"context"
	"sync"
	"time"

	"github.com/microsoft/moc-pkg/pkg/trace"
	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	//log "k8s.io/klog"
)

var once sync.Once

type client struct {
	wssdsecurity.AuthenticationAgentClient
	cloudFQDN string
}

var loginConfig auth.LoginConfig

// UpdateLoginConfig
func UpdateLoginConfig(loginconfig auth.LoginConfig) {
	loginConfig = loginconfig
}

// NewAuthenticationClient creates a client session with the backend wssd agent
func NewAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetAuthenticationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c, subID}, nil
}

func ReLoginOnExpiry(ctx context.Context, loginconfig auth.LoginConfig, group, cloudFQDN string) error {
	authorizer, err := auth.NewAuthorizerForAuth(loginconfig.Token, loginconfig.Certificate, cloudFQDN)
	if err != nil {
		return err
	}

	c, err := NewAuthenticationClient(cloudFQDN, authorizer)
	if err != nil {
		return err
	}
	_, err = c.LoginWithConfig(ctx, group, loginconfig, false)
	return err
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

func RenewRoutine(ctx context.Context, group, server string) {
	renewalAttempt := 0
	// Waiting for a few seconds to avoid spamming short-lived sdk user
	time.Sleep(time.Second * 5)
	var err error
	_, span := trace.NewSpan(ctx, "RenewRoutine")
	defer span.End(err)
	for {
		wssdConfig := auth.WssdConfig{}
		err := marshal.FromJSONFile(auth.GetWssdConfigLocation(), &wssdConfig)
		if err != nil {
			span.Log("Failed to open config file in location %s: %v\n", auth.GetWssdConfigLocation(), err)
			return
		}

		sleepTime, renewalBackoff, expiry, err := renewTime(wssdConfig.ClientCertificate)
		if err != nil {
			span.Log("Failed while calculating certificate renew time %v \n", err)
			return
		}
		span.Log("Waiting for %v to renew cert\n", sleepTime)
		time.Sleep(sleepTime)
		span.Log("Attempting to renew certificate\n")
		err = auth.RenewCertificates(server, auth.GetWssdConfigLocation())
		if err != nil {
			// If certificate is expired, we attempt to re-login with set login config
			if errors.IsExpired(err) {
				span.Log("Certificate Expired, Attemptin re-login %v", err)
				err = ReLoginOnExpiry(ctx, loginConfig, group, server)
				if err == nil {
					span.Log("Re-Login successful")
					renewalAttempt = 0
					continue
				} else {
					span.Log("Re-Login Failure %v", err)
				}
			}
			renewalAttempt += 1
			span.Log("Failed to renew cert: %v. Attempts %d", err, renewalAttempt)
			span.Log("Certificate Expiry %s, Now %s", expiry.UTC().String(), time.Now().UTC().String())
			time.Sleep(renewalBackoff)
			continue
		}
		//reset renewalAttempt after successful renewal
		renewalAttempt = 0
		span.Log("Certificate renewal complete\n")
	}
}

// Get methods invokes the client Get method
func (c *client) LoginWithConfig(ctx context.Context, group string, loginconfig auth.LoginConfig, enableRenewRoutine bool) (*auth.WssdConfig, error) {

	clientCsr, accessFile, err := auth.GenerateClientCsr(loginconfig)
	if err != nil {
		return nil, err
	}

	id := security.Identity{
		Name:        &loginconfig.Name,
		Certificate: &clientCsr,
	}

	clientCert, err := c.Login(ctx, group, &id)
	if err != nil {
		return nil, err
	}
	accessFile.ClientCertificate = *clientCert
	accessFile.ClientCertificateType = auth.CASigned
	accessFile.IdentityName = loginconfig.Name
	auth.PrintAccessFile(accessFile)
	UpdateLoginConfig(loginconfig)
	if enableRenewRoutine {
		once.Do(func() {
			go RenewRoutine(ctx, group, c.cloudFQDN)
		})
	}
	return &accessFile, nil
}

func calculateTime(before, after, now time.Time) (time.Duration, time.Duration) {
	validity := after.Sub(before)
	// renewBackoff is 2% of validity duration
	renewBackoff := time.Duration(float64(validity.Nanoseconds()) * constant.RenewalBackoff)
	// Threshold to renew is 30% of validity
	tresh := time.Duration(float64(validity.Nanoseconds()) * constant.CertificateValidityThreshold)

	treshNotAfter := after.Add(-tresh)
	return treshNotAfter.Sub(now), renewBackoff
}

func renewTime(certificate string) (sleepduration, renewBackoff time.Duration, expiry time.Time, err error) {

	pemCert, err := marshal.FromBase64(certificate)
	if err != nil {
		return
	}

	x509Cert, err := certs.DecodeCertPEM([]byte(pemCert))
	if err != nil {
		return
	}
	sleepduration, renewBackoff = calculateTime(x509Cert.NotBefore, x509Cert.NotAfter, time.Now())
	return sleepduration, renewBackoff, x509Cert.NotAfter, nil
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
