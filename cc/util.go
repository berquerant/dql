package cc

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseFloat(x string) (float64, error) { return strconv.ParseFloat(x, 64) }

func ParseInt(x string) (int, error) { return strconv.Atoi(x) }

func IsDigit(x rune) bool { return unicode.IsDigit(x) }

func IsSpace(x rune) bool { return unicode.IsSpace(x) }

func IsIdentTail(x rune) bool {
	return (unicode.IsLetter(x) || unicode.IsDigit(x)) || strings.ContainsRune("_", x) && !strings.ContainsRune("(){}[],:;+-*/=<>&|^~", x)
}

func IsIdentHead(x rune) bool {
	return unicode.IsLetter(x) || strings.ContainsRune("_", x) && !strings.ContainsRune("(){}[],:;+-*/=<>&|^~", x)
}
