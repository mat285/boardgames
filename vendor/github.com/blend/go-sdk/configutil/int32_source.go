/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// Int32Source is a type that can return a value.
type Int32Source interface {
	// Int32 should return an int32 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int32(ctx context.Context) (*int32, error)
}
