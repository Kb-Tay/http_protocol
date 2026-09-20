package utils

import "unicode"

func IsAlphabetic(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) {
			return false
		}
	}

	return true
}

func IsNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsNumber(c) {
			return false
		}
	}

	return true
}
