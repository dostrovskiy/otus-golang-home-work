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

func runeToNumber(r rune) (int, error) {
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
	if _, err := runeToNumber(runes[0]); err == nil {
		return "", ErrorStringStartsWithNumber
	}
	var sb strings.Builder
	var nextNum int
	var nextErr error
	for i, r := range runes {
		nextNum, nextErr = 0, errorInvalidString
		if i < len(runes)-1 {
			nextNum, nextErr = runeToNumber(runes[i+1])
		}
		if _, err := runeToNumber(r); err == nil {
			if nextErr == nil {
				return "", ErrorStringContainsSeveralDigitsInRow
			}
			continue
		}		
		if nextErr == nil {
			sb.WriteString(strings.Repeat(string(r), nextNum))
		} else {
			sb.WriteString(string(r))
		}
	}
	return sb.String(), nil
}
