package main

import (
	"fmt"
	"strings"
	"unicode"
)

type parser struct{
	isBasedonRegex bool
}


// Thx to AJ ONeal
// https://groups.google.com/forum/#!topic/golang-nuts/pNwqLyfl2co
const NullStr = rune(0)
func ParseArgs(str string) ([]string, error) {
	var m []string
	var s string

	str = strings.TrimSpace(str) + " "
	lastQuote := NullStr
	isSpace := false

	for i, c := range str {
		switch {
		// If we're ending a quote, break out and skip this character
		case c == lastQuote:
			lastQuote = NullStr

		// If we're in a quote, count this character
		case lastQuote != NullStr:
			s += string(c)

		// If we encounter a quote, enter it and skip this character
		case unicode.In(c, unicode.Quotation_Mark):
			isSpace = false
			lastQuote = c

		// If it's a space, store the string
		case unicode.IsSpace(c):
			if 0 == i || isSpace {
				continue
			}
			isSpace = true
			m = append(m, s)
			s = ""

		default:
			isSpace = false
			s += string(c)
		}
	}
	if lastQuote != NullStr {
		return nil, fmt.Errorf("Quotes did not terminate")
	}
	return m, nil
}
