// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package authentication

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/pkg/constant"
	"github.com/microsoft/moc-sdk-for-go/services/security"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/fs"
	"github.com/microsoft/moc/pkg/marshal"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
func newAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetAuthenticationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c, subID}, nil
}

func reLoginOnExpiry(ctx context.Context, loginconfig auth.LoginConfig, group, cloudFQDN string) error {
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

func renewRoutine(ctx context.Context, group, server string, logger logr.Logger) {
	renewalAttempt := 0
	// Waiting for a few seconds to avoid spamming short-lived sdk user
	time.Sleep(time.Second * 5)
	for {
		wssdConfig := auth.WssdConfig{}
		err := marshal.FromJSONFile(auth.GetWssdConfigLocation(), &wssdConfig)
		if err != nil {
			logger.Error(err, "Failed to open config file", "location", auth.GetWssdConfigLocation())
			return
		}

		sleepTime, renewalBackoff, expiry, err := renewTime(wssdConfig.ClientCertificate)
		if err != nil {
			logger.Error(err, "Failed while calculating certificate renew time")
			return
		}
		logger.Info("Waiting to renew certificate", "sleepTime", sleepTime)
		time.Sleep(sleepTime)
		logger.Info("Attempting to renew certificate")
		err = auth.RenewCertificates(server, auth.GetWssdConfigLocation())
		if err != nil {
			// If certificate is expired, we attempt to re-login with set login config
			if errors.IsExpired(err) {
				logger.Error(err, "Certificate expired, attempting re-login")
				err = reLoginOnExpiry(ctx, loginConfig, group, server)
				if err == nil {
					logger.Info("Re-login successful")
					renewalAttempt = 0
					continue
				} else {
					logger.Error(err, "Re-login failure")
				}
			}
			renewalAttempt++
			logger.Error(err, "Failed to renew certificate", "attempts", renewalAttempt)
			logger.Info("Certificate expiry details", "expiry", expiry.UTC().String(), "now", time.Now().UTC().String())
			time.Sleep(renewalBackoff)
			continue
		}
		// Reset renewalAttempt after successful renewal
		renewalAttempt = 0
		logger.Info("Certificate renewal complete")
	}
}

// Get methods invokes the client Get method
func (c *client) LoginWithConfig(ctx context.Context, group string, loginconfig auth.LoginConfig, enableRenewRoutine bool) (*auth.WssdConfig, error) {
	logger := log.FromContext(ctx) // Retrieve the logger from context
	if logger.GetSink() == nil {
		logger = logr.Discard() // Use a no-op logger to avoid panics
	}
	logger.Info("Generating client CSR")
	clientCsr, accessFile, err := auth.GenerateClientCsr(loginconfig)
	if err != nil {
		logger.Error(err, "Failed to generate client CSR")
		return nil, err
	}

	id := security.Identity{
		Name:        &loginconfig.Name,
		Certificate: &clientCsr,
	}

	logger.Info("Attempting to log in with client CSR")
	clientCert, err := c.Login(ctx, group, &id)
	if err != nil {
		logger.Error(err, "Login failed")
		return nil, err
	}
	accessFile.ClientCertificate = *clientCert
	accessFile.ClientCertificateType = auth.CASigned
	accessFile.IdentityName = loginconfig.Name

	logger.Info("Printing access file")
	if err := auth.PrintAccessFile(accessFile); err != nil {
		logger.Error(err, "PrintAccessFile failed")
		return &accessFile, errors.Wrap(err, "PrintAccessFile failed")
	}

	logger.Info("Setting file permissions for WSSD config location")
	if err := fs.Chmod(auth.GetWssdConfigLocation(), 0600); err != nil {
		logger.Error(err, "Failed to set file permissions")
		return &accessFile, err
	}
	UpdateLoginConfig(loginconfig)
	if enableRenewRoutine {
		logger.Info("Starting renew routine")
		once.Do(func() {
			go renewRoutine(ctx, group, c.cloudFQDN, logger)
		})
	}
	logger.Info("LoginWithConfig completed successfully")
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
