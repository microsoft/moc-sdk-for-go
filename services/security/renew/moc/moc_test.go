// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 license.
package moc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

var key *rsa.PrivateKey

func init() {
	key, _ = rsa.GenerateKey(rand.Reader, 2048)
}

func Test_GetCertRenewRequired(t *testing.T) {
	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: "test",
		},
		NotBefore: now,
		NotAfter:  now.Add(time.Second * 10),
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		t.Errorf("Failed creating certificate %v", err)
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		t.Errorf("Failed parsing certificate %v", err)
	}
	if renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:false Actual:true")
	}
}

func Test_GetCertRenewRequiredExpired(t *testing.T) {
	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: "test",
		},
		NotBefore: now.Add(-(time.Second * 10)),
		NotAfter:  now,
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		t.Errorf("Failed creating certificate %v", err)
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		t.Errorf("Failed parsing certificate %v", err)
	}
	if !renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:true Actual:false")
	}
}

func Test_GetCertRenewRequiredBeforeThreshold(t *testing.T) {
	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: "test",
		},
		NotBefore: now.Add(-(time.Second * 6)),
		NotAfter:  now.Add(time.Second * 4),
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		t.Errorf("Failed creating certificate %v", err)
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		t.Errorf("Failed parsing certificate %v", err)
	}
	if renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:false Actual:true")
	}
}

func Test_GetCertRenewRequiredAfterThreshold(t *testing.T) {
	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: "test",
		},
		NotBefore: now.Add(-(time.Second * 8)),
		NotAfter:  now.Add(time.Second * 2),
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		t.Errorf("Failed creating certificate %v", err)
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		t.Errorf("Failed parsing certificate %v", err)
	}
	if !renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:true Actual:false")
	}
}

func Test_GetCertRenewRequiredDelay(t *testing.T) {
	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: "test",
		},
		NotBefore: now.Add(-(time.Second * 6)),
		NotAfter:  now.Add(time.Second * 4),
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		t.Errorf("Failed creating certificate %v", err)
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		t.Errorf("Failed parsing certificate %v", err)
	}
	if renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:false Actual:true")
	}
	time.Sleep(time.Second * 2)
	if !renewRequired(x509Cert) {
		t.Errorf("RenewRequired Expected:true Actual:false")
	}
}
