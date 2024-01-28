/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// PercentDifference computes the percentage difference between two values.
// The formula is (v2-v1)/v1.
func PercentDifference(v1, v2 float64) float64 {
	if v1 == 0 {
		return 0
	}
	return (v2 - v1) / v1
}
