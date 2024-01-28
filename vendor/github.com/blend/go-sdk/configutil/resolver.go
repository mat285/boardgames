/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import "context"

// Resolver is a type that can be resolved.
type Resolver interface {
	Resolve(context.Context) error
}
