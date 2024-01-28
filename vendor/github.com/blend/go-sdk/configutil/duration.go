/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"time"
)

var (
	_ DurationSource = (*Duration)(nil)
)

// Duration implements value provider.
//
// If the value is zero, a nil is returned by the implementation indicating
// the value was not present.
//
// If you want 0 to be a valid value, you must use DurationPtr.
type Duration time.Duration

// Duration returns the value for a constant.
func (dc Duration) Duration(_ context.Context) (*time.Duration, error) {
	if dc > 0 {
		value := time.Duration(dc)
		return &value, nil
	}
	return nil, nil
}
