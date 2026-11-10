package utils

import "unicode"

func MixCase(s string) string {
	runes := []rune(s)
	letterIndex := 0
	for i, r := range runes {
		if !unicode.IsLetter(r) {
			continue
		}
		if letterIndex%2 == 0 {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
		letterIndex++
	}
	return string(runes)
}
