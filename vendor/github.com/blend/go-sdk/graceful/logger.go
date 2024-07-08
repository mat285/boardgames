/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package graceful

// Logger is a type that can be used as a graceful process logger.
type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// MaybeInfof calls the logger infof method if the logger is set.
func MaybeInfof(log Logger, format string, args ...interface{}) {
	if log != nil {
		log.Infof(format, args...)
	}
}

// MaybeDebugf  calls the logger debugf method if the logger is set.
func MaybeDebugf(log Logger, format string, args ...interface{}) {
	if log != nil {
		log.Debugf(format, args...)
	}
}

// MaybeErrorf  calls the logger errorf method if the logger is set.
func MaybeErrorf(log Logger, format string, args ...interface{}) {
	if log != nil {
		log.Errorf(format, args...)
	}
}
