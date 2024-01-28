/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import "net/http"

// Middleware is a wrapping function that takes a handler and returns a handler.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// NestMiddleware nests middleware steps.
func NestMiddleware(action http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	if len(middleware) == 0 {
		return action
	}

	var nest = func(a, b Middleware) Middleware {
		if b == nil {
			return a
		}
		return func(inner http.HandlerFunc) http.HandlerFunc {
			return a(b(inner))
		}
	}

	var outer Middleware
	for _, step := range middleware {
		outer = nest(step, outer)
	}
	return outer(action)
}
