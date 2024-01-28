/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blend/go-sdk/ansi"
)

// FormatLabels formats the output of labels as a string.
// Field keys will be printed in alphabetic order.
func FormatLabels(tf TextFormatter, keyColor ansi.Color, labels Labels) string {
	var keys []string
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%s", tf.Colorize(key, keyColor), labels[key]))
	}
	return strings.Join(values, " ")
}
