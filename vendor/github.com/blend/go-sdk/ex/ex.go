/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ error          = (*Ex)(nil)
	_ fmt.Formatter  = (*Ex)(nil)
	_ json.Marshaler = (*Ex)(nil)
)

// New returns a new exception with a call stack.
// Pragma: this violates the rule that you should take interfaces and return
// concrete types intentionally; it is important for the semantics of typed pointers and nil
// for this to return an interface because (*Ex)(nil) != nil, but (error)(nil) == nil.
func New(class interface{}, options ...Option) Exception {
	return NewWithStackDepth(class, DefaultNewStartDepth, options...)
}

// NewWithStackDepth creates a new exception with a given start point of the stack.
func NewWithStackDepth(class interface{}, startDepth int, options ...Option) Exception {
	if class == nil {
		return nil
	}

	var ex *Ex
	switch typed := class.(type) {
	case *Ex:
		if typed == nil {
			return nil
		}
		ex = typed
	case error:
		if typed == nil {
			return nil
		}

		ex = &Ex{
			Class:      typed,
			Inner:      errors.Unwrap(typed),
			StackTrace: Callers(startDepth),
		}
	case string:
		ex = &Ex{
			Class:      Class(typed),
			StackTrace: Callers(startDepth),
		}
	default:
		ex = &Ex{
			Class:      Class(fmt.Sprint(class)),
			StackTrace: Callers(startDepth),
		}
	}
	for _, option := range options {
		option(ex)
	}
	return ex
}

// Ex is an error with a stack trace.
// It also can have an optional cause, it implements `Exception`
type Ex struct {
	// Class disambiguates between errors, it can be used to identify the type of the error.
	Class error
	// Message adds further detail to the error, and shouldn't be used for disambiguation.
	Message string
	// Inner holds the original error in cases where we're wrapping an error with a stack trace.
	Inner error
	// StackTrace is the call stack frames used to create the stack output.
	StackTrace StackTrace
}

// WithMessage sets the exception message.
// Deprecation notice: This method is included as a migraition path from v2, and will be removed after v3.
func (e *Ex) WithMessage(args ...interface{}) Exception {
	e.Message = fmt.Sprint(args...)
	return e
}

// WithMessagef sets the exception message based on a format and arguments.
// Deprecation notice: This method is included as a migration path from v2, and will be removed after v3.
func (e *Ex) WithMessagef(format string, args ...interface{}) Exception {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// WithInner sets the inner ex.
// Deprecation notice: This method is included as a migraition path from v2, and will be removed after v3.
func (e *Ex) WithInner(err error) Exception {
	e.Inner = NewWithStackDepth(err, DefaultNewStartDepth)
	return e
}

// Format allows for conditional expansion in printf statements
// based on the token and flags used.
// 	%+v : class + message + stack
// 	%v, %c : class
// 	%m : message
// 	%t : stack
func (e *Ex) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if e.Class != nil && len(e.Class.Error()) > 0 {
			fmt.Fprint(s, e.Class.Error())
		}
		if len(e.Message) > 0 {
			fmt.Fprint(s, "; "+e.Message)
		}
		if s.Flag('+') && e.StackTrace != nil {
			e.StackTrace.Format(s, verb)
		}
		if e.Inner != nil {
			if typed, ok := e.Inner.(fmt.Formatter); ok {
				fmt.Fprint(s, "\n")
				typed.Format(s, verb)
			} else {
				fmt.Fprintf(s, "\n%v", e.Inner)
			}
		}
		return
	case 'c':
		fmt.Fprint(s, e.Class.Error())
	case 'i':
		if e.Inner != nil {
			if typed, ok := e.Inner.(fmt.Formatter); ok {
				typed.Format(s, verb)
			} else {
				fmt.Fprintf(s, "%v", e.Inner)
			}
		}
	case 'm':
		fmt.Fprint(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}

// Error implements the `error` interface.
// It returns the exception class, without any of the other supporting context like the stack trace.
// To fetch the stack trace, use .String().
func (e *Ex) Error() string {
	return e.Class.Error()
}

// Decompose breaks the exception down to be marshaled into an intermediate format.
func (e *Ex) Decompose() map[string]interface{} {
	values := map[string]interface{}{}
	values["Class"] = e.Class.Error()
	values["Message"] = e.Message
	if e.StackTrace != nil {
		values["StackTrace"] = e.StackTrace.Strings()
	}
	if e.Inner != nil {
		if typed, isTyped := e.Inner.(*Ex); isTyped {
			values["Inner"] = typed.Decompose()
		} else {
			values["Inner"] = e.Inner.Error()
		}
	}
	return values
}

// MarshalJSON is a custom json marshaler.
func (e *Ex) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Decompose())
}

// UnmarshalJSON is a custom json unmarshaler.
func (e *Ex) UnmarshalJSON(contents []byte) error {
	// try first as a string ...
	var class string
	if tryErr := json.Unmarshal(contents, &class); tryErr == nil {
		e.Class = Class(class)
		return nil
	}

	// try an object ...
	values := make(map[string]json.RawMessage)
	if err := json.Unmarshal(contents, &values); err != nil {
		return New(err)
	}

	if class, ok := values["Class"]; ok {
		var classString string
		if err := json.Unmarshal([]byte(class), &classString); err != nil {
			return New(err)
		}
		e.Class = Class(classString)
	}

	if message, ok := values["Message"]; ok {
		if err := json.Unmarshal([]byte(message), &e.Message); err != nil {
			return New(err)
		}
	}

	if inner, ok := values["Inner"]; ok {
		var innerClass string
		if tryErr := json.Unmarshal([]byte(inner), &class); tryErr == nil {
			e.Inner = Class(innerClass)
		}
		innerEx := Ex{}
		if tryErr := json.Unmarshal([]byte(inner), &innerEx); tryErr == nil {
			e.Inner = &innerEx
		}
	}
	if stack, ok := values["StackTrace"]; ok {
		var stackStrings []string
		if err := json.Unmarshal([]byte(stack), &stackStrings); err != nil {
			return New(err)
		}
		e.StackTrace = StackStrings(stackStrings)
	}

	return nil
}

// String returns a fully formed string representation of the ex.
// It's equivalent to calling sprintf("%+v", ex).
func (e *Ex) String() string {
	s := new(bytes.Buffer)
	if e.Class != nil && len(e.Class.Error()) > 0 {
		fmt.Fprintf(s, "%s", e.Class)
	}
	if len(e.Message) > 0 {
		fmt.Fprint(s, " "+e.Message)
	}
	if e.StackTrace != nil {
		fmt.Fprint(s, " "+e.StackTrace.String())
	}
	return s.String()
}

// Unwrap returns the inner error if it exists.
// Enables error chaining and calling errors.Is/As to
// match on inner errors.
func (e *Ex) Unwrap() error {
	return e.Inner
}

// Is returns true if the target error matches the Ex.
// Enables errors.Is on Ex classes when an error
// is wrapped using Ex.
func (e *Ex) Is(target error) bool {
	return Is(e, target)
}

// As delegates to the errors.As to match on the Ex class.
func (e *Ex) As(target interface{}) bool {
	return errors.As(e.Class, target)
}
