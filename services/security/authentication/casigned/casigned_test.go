// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 license.
package casigned

import (
	"testing"
	"time"
)

func Test_CalculateTime(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * -10))
	after := now.Add(time.Duration(time.Second * 10))
	actual := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*4) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*4), actual)
	}
}

func Test_CalculateTimeNegative(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * -20))
	after := now.Add(time.Duration(time.Second * -10))
	actual := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*-13) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*-13), actual)
	}
}

func Test_CalculateTimeAfter(t *testing.T) {
	now := time.Now()
	before := now.Add(time.Duration(time.Second * 10))
	after := now.Add(time.Duration(time.Second * 30))
	actual := calculateTime(before, after, now)
	if actual != time.Duration(time.Second*24) {
		t.Errorf("Wrong wait time returned Expected %s Actual %s", time.Duration(time.Second*24), actual)
	}
}
