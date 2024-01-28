/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

// HasPrefixCaseless returns if a corpus has a prefix regardless of casing.
func HasPrefixCaseless(corpus, prefix string) bool {
	corpusLen := len(corpus)
	prefixLen := len(prefix)

	if corpusLen < prefixLen {
		return false
	}

	for x := 0; x < prefixLen; x++ {
		charCorpus := uint(corpus[x])
		charPrefix := uint(prefix[x])

		if charCorpus-LowerA <= LowerDiff {
			charCorpus = charCorpus - 0x20
		}

		if charPrefix-LowerA <= LowerDiff {
			charPrefix = charPrefix - 0x20
		}
		if charCorpus != charPrefix {
			return false
		}
	}
	return true
}
