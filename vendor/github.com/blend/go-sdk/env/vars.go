/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/reflectutil"
	"github.com/blend/go-sdk/stringutil"
)

// New returns a new env var set.
//
/* By default, it is empty. In order to populate it with the current
runtime environment variables, you need to pass in options:
    vars := env.New(env.OptEnviron(os.Environ()...))
*/
func New(opts ...Option) Vars {
	vars := make(Vars)
	for _, opt := range opts {
		opt(vars)
	}
	return vars
}

// Vars is a set of environment variables.
type Vars map[string]string

// Get gets a variable as a string.
// It mirrors os.Getenv.
func (ev Vars) Get(envVar string) string {
	if value, ok := ev[envVar]; ok {
		return value
	}
	return ""
}

// Set sets a value for a key.
func (ev Vars) Set(envVar, value string) {
	ev[envVar] = value
}

// Restore resets an environment variable to it's environment value.
func (ev Vars) Restore(key string) {
	ev[key] = os.Getenv(key)
}

// Delete removes a key from the set.
func (ev Vars) Delete(key string) {
	delete(ev, key)
}

// String returns a string value for a given key, with an optional default vaule.
func (ev Vars) String(envVar string, defaults ...string) string {
	if value, hasValue := ev[envVar]; hasValue {
		return value
	}
	for _, defaultValue := range defaults {
		if defaultValue != "" {
			return defaultValue
		}
	}
	return ""
}

// CSV returns a string array for a given string var.
func (ev Vars) CSV(envVar string, defaults ...string) []string {
	if value, hasValue := ev[envVar]; hasValue && len(value) > 0 {
		return strings.Split(value, ",")
	}
	return defaults
}

// Bool returns a boolean value for a key, defaulting to false.
// Valid "truthy" values are `true`, `yes`, and `1`.
// Everything else is false, including `REEEEEEEEEEEEEEE`.
func (ev Vars) Bool(envVar string, defaults ...bool) bool {
	if value, hasValue := ev[envVar]; hasValue {
		boolValue, err := stringutil.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return false
}

// Int returns an integer value for a given key.
func (ev Vars) Int(envVar string, defaults ...int) (int, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return parsedValue, nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustInt returns an integer value for a given key and panics if it is malformed.
func (ev Vars) MustInt(envVar string, defaults ...int) int {
	value, err := ev.Int(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Int32 returns an integer value for a given key.
func (ev Vars) Int32(envVar string, defaults ...int32) (int32, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return int32(parsedValue), nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustInt32 returns an integer value for a given key and panics if it is malformed.
func (ev Vars) MustInt32(envVar string, defaults ...int32) int32 {
	value, err := ev.Int32(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Int64 returns an int64 value for a given key.
func (ev Vars) Int64(envVar string, defaults ...int64) (int64, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return parsedValue, nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustInt64 returns an int64 value for a given key and panics if it is malformed.
func (ev Vars) MustInt64(envVar string, defaults ...int64) int64 {
	value, err := ev.Int64(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint32 returns an uint32 value for a given key.
func (ev Vars) Uint32(envVar string, defaults ...uint32) (uint32, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return uint32(parsedValue), nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustUint32 returns an uint32 value for a given key and panics if it is malformed.
func (ev Vars) MustUint32(envVar string, defaults ...uint32) uint32 {
	value, err := ev.Uint32(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Uint64 returns an uint64 value for a given key.
func (ev Vars) Uint64(envVar string, defaults ...uint64) (uint64, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return parsedValue, nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustUint64 returns an uint64 value for a given key and panics if it is malformed.
func (ev Vars) MustUint64(envVar string, defaults ...uint64) uint64 {
	value, err := ev.Uint64(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Float32 returns an float32 value for a given key.
func (ev Vars) Float32(envVar string, defaults ...float32) (float32, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return float32(parsedValue), nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustFloat32 returns an float64 value for a given key and panics if it is malformed.
func (ev Vars) MustFloat32(envVar string, defaults ...float32) float32 {
	value, err := ev.Float32(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Float64 returns an float64 value for a given key.
func (ev Vars) Float64(envVar string, defaults ...float64) (float64, error) {
	if value, hasValue := ev[envVar]; hasValue {
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, ex.New(err, ex.OptMessagef("var: %q", envVar))
		}
		return parsedValue, nil
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustFloat64 returns an float64 value for a given key and panics if it is malformed.
func (ev Vars) MustFloat64(envVar string, defaults ...float64) float64 {
	value, err := ev.Float64(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Duration returns a duration value for a given key.
func (ev Vars) Duration(envVar string, defaults ...time.Duration) (time.Duration, error) {
	if value, hasValue := ev[envVar]; hasValue {
		return time.ParseDuration(value)
	}
	for _, defaultValue := range defaults {
		if defaultValue > 0 {
			return defaultValue, nil
		}
	}
	return 0, nil
}

// MustDuration returnss a duration value for a given key and panics if malformed.
func (ev Vars) MustDuration(envVar string, defaults ...time.Duration) time.Duration {
	value, err := ev.Duration(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Bytes returns a []byte value for a given key.
func (ev Vars) Bytes(envVar string, defaults ...[]byte) []byte {
	if value, hasValue := ev[envVar]; hasValue && len(value) > 0 {
		return []byte(value)
	}
	for _, defaultValue := range defaults {
		if len(defaultValue) > 0 {
			return defaultValue
		}
	}
	return nil
}

// Base64 returns a []byte value for a given key whose value is encoded in base64.
func (ev Vars) Base64(envVar string, defaults ...[]byte) ([]byte, error) {
	if value, hasValue := ev[envVar]; hasValue && len(value) > 0 {
		return base64.StdEncoding.DecodeString(value)
	}
	for _, defaultValue := range defaults {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
	}
	return nil, nil
}

// MustBase64 returns a []byte value for a given key encoded with base64, and panics if malformed.
func (ev Vars) MustBase64(envVar string, defaults ...[]byte) []byte {
	value, err := ev.Base64(envVar, defaults...)
	if err != nil {
		panic(err)
	}
	return value
}

// Has returns if a key is present in the set.
func (ev Vars) Has(envVar string) bool {
	_, hasKey := ev[envVar]
	return hasKey
}

// HasAll returns if all of the given vars are present in the set.
func (ev Vars) HasAll(envVars ...string) bool {
	if len(envVars) == 0 {
		return false
	}
	for _, envVar := range envVars {
		if !ev.Has(envVar) {
			return false
		}
	}
	return true
}

// HasAny returns if any of the given vars are present in the set.
func (ev Vars) HasAny(envVars ...string) bool {
	for _, envVar := range envVars {
		if ev.Has(envVar) {
			return true
		}
	}
	return false
}

// Require enforces that a given set of environment variables are present.
func (ev Vars) Require(keys ...string) error {
	for _, key := range keys {
		if !ev.Has(key) {
			return fmt.Errorf("the following environment variables are required: `%s`", strings.Join(keys, ","))
		}
	}
	return nil
}

// Must enforces that a given set of environment variables are present and panics
// if they're not present.
func (ev Vars) Must(keys ...string) {
	for _, key := range keys {
		if !ev.Has(key) {
			panic(fmt.Sprintf("the following environment variables are required: `%s`", strings.Join(keys, ",")))
		}
	}
}

// Union returns the union of the two sets, other replacing conflicts.
func (ev Vars) Union(other Vars) Vars {
	newSet := New()
	for key, value := range ev {
		newSet[key] = value
	}
	for key, value := range other {
		newSet[key] = value
	}
	return newSet
}

// Vars returns all the vars stored in the env var set.
func (ev Vars) Vars() []string {
	var envVars = make([]string, len(ev))
	var index int
	for envVar := range ev {
		envVars[index] = envVar
		index++
	}
	return envVars
}

// Raw returns a raw KEY=VALUE form of the vars.
func (ev Vars) Raw() []string {
	var raw []string
	for key, value := range ev {
		raw = append(raw, fmt.Sprintf("%s=%s", key, value))
	}
	return raw
}

// ReadInto sets an object based on the fields in the env vars set.
func (ev Vars) ReadInto(obj interface{}) error {
	if typed, isTyped := obj.(Unmarshaler); isTyped {
		return typed.UnmarshalEnv(ev)
	}
	return reflectutil.PatchStrings(ReflectTagName, ev, obj)
}

// Expand calls os.Expand with the variable set as the environment value resolver.
func (ev Vars) Expand(value string) string {
	return os.Expand(value, ev.Get)
}

// --------------------------------------------------------------------------------
// Service Specific helpers
// --------------------------------------------------------------------------------

// ServiceEnv is a common environment variable for the services environment.
// Common values include "dev", "ci", "sandbox", "preprod", "beta", and "prod".
func (ev Vars) ServiceEnv(defaults ...string) string {
	return ev.String(VarServiceEnv, defaults...)
}

// IsProduction returns if the ServiceEnv is a production environment.
func (ev Vars) IsProduction() bool {
	return IsProduction(ev.ServiceEnv())
}

// IsProdlike returns if the ServiceEnv is "prodlike".
func (ev Vars) IsProdlike() bool {
	return IsProdlike(ev.ServiceEnv())
}

// IsDev returns if the ServiceEnv is the local development environment.
func (ev Vars) IsDev() bool {
	return IsDev(ev.ServiceEnv())
}

// IsDevlike returns if the ServiceEnv is strictly the inverse of `IsProdlike`.
func (ev Vars) IsDevlike() bool {
	return !IsProdlike(ev.ServiceEnv())
}

// ServiceName is a common environment variable for the service's name.
func (ev Vars) ServiceName(defaults ...string) string {
	return ev.String(VarServiceName, defaults...)
}

// Hostname is a common environment variable for the machine's hostname.
func (ev Vars) Hostname(defaults ...string) string {
	return ev.String(VarHostname, defaults...)
}

// Version is a common environment variable for the service version.
func (ev Vars) Version(defaults ...string) string {
	return ev.String(VarVersion, defaults...)
}
