// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package constant

import (
	"time"
)

const (
	DefaultServerContextTimeout          = 10 * time.Minute
	CertificateValidityThreshold float64 = (30.0 / 100.0)
	RenewalBackoff               float64 = (2.0 / 100.0)
	OsRegistrationStatus         string  = "osRegistrationStatus"
	OsVersion                    string  = "osVersion"
)
