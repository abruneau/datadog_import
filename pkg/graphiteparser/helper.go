package graphiteparser

import (
	"math"
)

func isDecimalDigit(s string) bool {
	return len(s) == 1 && s[0] >= '0' && s[0] <= '9'
}

func isOctalDigit(s string) bool {
	return len(s) == 1 && s[0] >= '0' && s[0] <= '7'
}

func isHexDigit(s string) bool {
	return len(s) == 1 && (s[0] >= '0' && s[0] <= '9' || s[0] >= 'a' && s[0] <= 'f' || s[0] >= 'A' && s[0] <= 'F')
}

func isIdentifierStart(ch string) bool {
	// Very simplified version.
	return ch == "$" || ch == "_" || ch == "\\" || (ch >= "a" && ch <= "z") || (ch >= "A" && ch <= "Z")
}

// isPunctuator returns true if ch is one of the punctuation characters.
func isPunctuator(ch string) bool {
	switch ch {
	case ".", "(", ")", ",", "{", "}":
		return true
	}
	return false
}

// isFinite is a helper to check if a float64 is finite.
func isFinite(num float64) bool {
	return !math.IsInf(num, 0) && !math.IsNaN(num)
}
