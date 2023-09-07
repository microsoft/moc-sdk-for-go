// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

package diagnostics

import "context"

type key int

const CorrelationId key = 1

func GetCorrelationId(ctx context.Context) string {
	id, ok := ctx.Value(CorrelationId).(string)
	if ok {
		return id
	} else {
		return ""
	}
}

func NewContextWithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, CorrelationId, correlationId)
}
