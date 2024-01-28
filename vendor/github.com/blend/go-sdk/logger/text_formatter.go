/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import "github.com/blend/go-sdk/ansi"

// TextFormatter is a type that can format text output.
type TextFormatter interface {
	Colorize(string, ansi.Color) string
}
