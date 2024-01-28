/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"context"
)

type parameterizedPathKey struct{}

// WithParameterizedPath adds a path with named parameters to a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved.
func WithParameterizedPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, parameterizedPathKey{}, path)
}

// GetParameterizedPath gets a path with named parameters off a context. Useful for
// outbound request aggregation for metrics and tracing when route parameters
// are involved. Relies on OptParameterizedPath being added to a Request.
func GetParameterizedPath(ctx context.Context) string {
	path, _ := ctx.Value(parameterizedPathKey{}).(string)
	return path
}
