// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package scanner

import (
	"unicode"
)

// IsIdentifier returns whether a string is an identifier
func IsIdentifier(text string) bool {

	if IsKeyword(text) {
		return false
	}

	for i, r := range text {
		if i == 0 {
			if !IsIdentStart(r) {
				return false
			}
		} else {
			if !IsIdentContinue(r) {
				return false
			}
		}
	}

	return true
}

// IsIdentStart returns whether a rune can be the start of an identifier
func IsIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// IsIdentContinue returns whether a rune can be in the middle of an identifier
func IsIdentContinue(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
