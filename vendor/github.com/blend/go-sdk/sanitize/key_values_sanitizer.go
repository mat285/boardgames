/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sanitize

// KeyValuesSanitizer is a type that can sanitize http header or query key values.
//
// Values are passed from `map[string][]string` typically,
// hence the `key` and `values...` parameters.
//
// The `SanitizeValue` function should return the modified or
// sanitized value for each of the input values, for a given key.
type KeyValuesSanitizer interface {
	SanitizeKeyValues(key string, values ...string) []string
}

// KeyValuesSanitizerFunc is a function implementation of ValueSanitizer.
type KeyValuesSanitizerFunc func(key string, values ...string) []string

// SanitizeKeyValues implements `KeyValuesSanitizer`.
func (vsf KeyValuesSanitizerFunc) SanitizeKeyValues(key string, values ...string) []string {
	return vsf(key, values...)
}

// DefaultKeyValuesSanitizerFunc is the default value sanitizer.
//
// For any given key's values it will simply return nil, implying that
// the key was one of the banned keys and we should completely omit
// the values.
func DefaultKeyValuesSanitizerFunc(_ string, _ ...string) []string {
	return nil
}
