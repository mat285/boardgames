/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

var (
	_ Float64Source = (*Float64Func)(nil)
)

// Float64Func is a float value source from a commandline flag.
type Float64Func func(context.Context) (*float64, error)

// Float64 returns an invocation of the function.
func (vf Float64Func) Float64(ctx context.Context) (*float64, error) {
	return vf(ctx)
}
