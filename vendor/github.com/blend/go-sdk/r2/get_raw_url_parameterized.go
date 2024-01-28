/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"strings"
)

// GetRawURLParameterized gets a URL string with named route parameters in place of
// the raw path for a request. Useful for outbound request aggregation for
// metrics and tracing when route parameters are involved.
// Relies on the request's context storing the parameterized path, otherwise will default
// to returning the request `URL`'s `String()`.
func GetRawURLParameterized(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	url := req.URL
	path := GetParameterizedPath(req.Context())
	if path == "" {
		return url.String()
	}

	// Stripped down version of "net/url" `URL.String()`
	var buf strings.Builder
	if url.Scheme != "" {
		buf.WriteString(url.Scheme)
		buf.WriteByte(':')
	}
	if url.Scheme != "" || url.Host != "" {
		if url.Host != "" || url.Path != "" {
			buf.WriteString("//")
		}
		if host := url.Host; host != "" {
			buf.WriteString(host)
		}
	}
	if !strings.HasPrefix(path, "/") {
		buf.WriteByte('/')
	}
	buf.WriteString(path)
	return buf.String()
}
