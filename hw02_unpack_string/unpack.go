/*
Package hw02unpackstring unpack string with format <char><number> will be replaced with <char><char><char>... number times, e.g. "a4bc2d5e" -> "aaaabccddddde"
*/
package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var errorInvalidString = errors.New("invalid string")

/*
ErrorStringContainsSeveralDigitsInRow is returned when string contains several digits in a row, e.g. "aaa10b".
*/
var ErrorStringContainsSeveralDigitsInRow = errors.New("string contains several digits in a row")

/*
ErrorStringStartsWithNumber is returned when string starts with number, e.g. "3abc".
*/
var ErrorStringStartsWithNumber = errors.New("string starts with number")

/*
ErrorStringEndsWithBackSlashEscapingNothing is returned when string ends with backslash escaping nothing, e.g. "abc\".
*/
var ErrorStringEndsWithBackSlashEscapingNothing = errors.New("string ends with backslash escaping nothing")

func runeToDigit(r rune) (int, error) {
	if !unicode.IsDigit(r) {
		return 0, fmt.Errorf("rune is not a digit")
	}
	return int(r - '0'), nil
}

/*
Unpack string with format <char><number> will be replaced with <char><char><char>... number times, e.g. "a4bc2d5e" -> "aaaabccddddde"
*/
func Unpack(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}
	var runes = []rune(s)
	if _, err := runeToDigit(runes[0]); err == nil {
		return "", ErrorStringStartsWithNumber
	}
	var sb strings.Builder
	var nextDig int
	var nextDigErr error
	var prevBackSlash bool
	for i, r := range runes {
		nextDig, nextDigErr = 0, errorInvalidString
		if i < len(runes)-1 {
			nextDig, nextDigErr = runeToDigit(runes[i+1]) // if next rune is a digit, will repeat current rune
		}
		if _, err := runeToDigit(r); err == nil && !prevBackSlash {
			if nextDigErr == nil {
				return "", ErrorStringContainsSeveralDigitsInRow
			}
			continue
		}
		if r == '\\' && !prevBackSlash {
			if i == len(runes)-1 {
				return "", ErrorStringEndsWithBackSlashEscapingNothing
			}
			prevBackSlash = true
			continue
		}		
		if nextDigErr == nil {
			sb.WriteString(strings.Repeat(string(r), nextDig))
		} else {
			sb.WriteString(string(r))
		}
		prevBackSlash = false
	}
	return sb.String(), nil
}
