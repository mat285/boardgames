/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

// Unmarshaler is a type that implements `UnmarshalEnv`.
type Unmarshaler interface {
	UnmarshalEnv(vars Vars) error
}
