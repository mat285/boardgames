/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import "context"

// Checker is a type that can be checked for SLA status.
type Checker interface {
	Check(context.Context) error
}

var (
	_ Checker = (*CheckerFunc)(nil)
)

// CheckerFunc implements Checker.
type CheckerFunc func(context.Context) error

// Check implements Checker.
func (cf CheckerFunc) Check(ctx context.Context) error {
	return cf(ctx)
}
