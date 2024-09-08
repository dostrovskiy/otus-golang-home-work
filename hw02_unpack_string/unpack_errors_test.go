package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpackErrorStringStartsWithNumber(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "0aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrorStringStartsWithNumber), "actual error %q", err)
		})
	}
}

func TestUnpackErrorStringContainsSeveralDigitsInRow(t *testing.T) {
	invalidStrings := []string{"abc500", "象4象形00文字", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrorStringContainsSeveralDigitsInRow), "actual error %q", err)
		})
	}
}

func TestUnpackErrorStringEndsWithBackSlashEscapingNothing(t *testing.T) {
	invalidStrings := []string{`abc\`, `abc\\\`, `\\\\\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrorStringEndsWithBackSlashEscapingNothing), "actual error %q", err)
		})
	}
}
