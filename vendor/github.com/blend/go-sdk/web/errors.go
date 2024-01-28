/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"github.com/golang-jwt/jwt"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrSessionIDEmpty is thrown if a session id is empty.
	ErrSessionIDEmpty ex.Class = "auth session id is empty"
	// ErrSecureSessionIDEmpty is an error that is thrown if a given secure session id is invalid.
	ErrSecureSessionIDEmpty ex.Class = "auth secure session id is empty"
	// ErrUnsetViewTemplate is an error that is thrown if a given secure session id is invalid.
	ErrUnsetViewTemplate ex.Class = "view result template is unset"
	// ErrParameterMissing is an error on request validation.
	ErrParameterMissing ex.Class = "parameter is missing"
	// ErrParameterInvalid is an error on request validation.
	ErrParameterInvalid ex.Class = "parameter is invalid"
)

// NewParameterMissingError returns a new parameter missing error.
func NewParameterMissingError(paramName string) error {
	return ex.New(ErrParameterMissing, ex.OptMessagef("%s", paramName))
}

// NewParameterInvalidError returns a new parameter invalid error.
func NewParameterInvalidError(paramName, message string) error {
	return ex.New(ErrParameterMissing, ex.OptMessagef("%q: %s", paramName, message))
}

// IsErrSessionInvalid returns if an error is a session invalid error.
func IsErrSessionInvalid(err error) bool {
	if err == nil {
		return false
	}
	if ex.Is(err, ErrSessionIDEmpty) ||
		ex.Is(err, ErrSecureSessionIDEmpty) ||
		isValidationError(err) {
		return true
	}
	return false
}

func isValidationError(err error) bool {
	_, ok := err.(*jwt.ValidationError)
	return ok
}

// IsErrBadRequest returns if an error is a bad request triggering error.
func IsErrBadRequest(err error) bool {
	return IsErrParameterMissing(err) || IsErrParameterInvalid(err)
}

// IsErrParameterMissing returns if an error is an ErrParameterMissing.
func IsErrParameterMissing(err error) bool {
	if err == nil {
		return false
	}
	return ex.Is(err, ErrParameterMissing)
}

// IsErrParameterInvalid returns if an error is an ErrParameterInvalid.
func IsErrParameterInvalid(err error) bool {
	if err == nil {
		return false
	}
	return ex.Is(err, ErrParameterMissing)
}
