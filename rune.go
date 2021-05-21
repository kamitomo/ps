package ps

import (
	"strings"
	"unicode"
)

func isSymbol(r rune) bool {
	return strings.ContainsRune(`+-*/<>=&?.@_#$:*`, r)
}

func isLetterDigitSymbol(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || isSymbol(r)
}
