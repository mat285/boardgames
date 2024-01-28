/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

var (
	_ Int64Source = (*Int64)(nil)
)

// Int64 implements value provider.
//
// Note: Int64 treats 0 as unset, if 0 is a valid value you must use configutil.Int64Ptr.
type Int64 int64

// Int64 returns the value for a constant.
func (i Int64) Int64(_ context.Context) (*int64, error) {
	if i > 0 {
		value := int64(i)
		return &value, nil
	}
	return nil, nil
}
