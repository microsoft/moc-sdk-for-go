// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package constant

import (
	"sync"
	"time"
)

const (
	DefaultServerContextTimeout          = 10 * time.Minute
	CertificateValidityThreshold float64 = (30.0 / 100.0)
	RenewalBackoff               float64 = (2.0 / 100.0)
)

type ClientOpts struct {
	NoExitOnConnFailure bool
}

var clientOptsSingleton *ClientOpts
var once sync.Once

//If GetClientOpts called before SetClientOpts, it will return a default ClientOpts
func GetClientOpts() *ClientOpts {
	once.Do(func() {
		clientOptsSingleton = &ClientOpts{}
	})
	return clientOptsSingleton
}

func SetClientOpts(clientOpts *ClientOpts) {
	once.Do(func() {
		clientOptsSingleton = clientOpts
	})
}
