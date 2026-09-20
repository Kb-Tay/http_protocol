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
