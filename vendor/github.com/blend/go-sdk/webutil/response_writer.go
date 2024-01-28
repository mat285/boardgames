/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import "net/http"

// ResponseWriter is a response writer that also returns the written status.
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	StatusCode() int
	InnerResponse() http.ResponseWriter
}
