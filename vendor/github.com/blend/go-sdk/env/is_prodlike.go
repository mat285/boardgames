/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

// IsProduction returns if the environment is production.
func IsProduction(serviceEnv string) bool {
	switch serviceEnv {
	case ServiceEnvPreprod, ServiceEnvProd:
		return true
	default:
		return false
	}
}

// IsProdlike returns if the environment is prodlike.
func IsProdlike(serviceEnv string) bool {
	switch serviceEnv {
	case ServiceEnvDev, ServiceEnvCI, ServiceEnvTest, ServiceEnvSandbox:
		return false
	default:
		return true
	}
}
