/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import "net/http"

// GetContentType gets the content type out of a header collection.
func GetContentType(header http.Header) string {
	if header != nil {
		return header.Get(HeaderContentType)
	}
	return ""
}
