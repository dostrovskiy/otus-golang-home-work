package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("unsupported file error", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("no such file error", func(t *testing.T) {
		err := Copy("no_such_file.txt", "out.txt", 0, 0)
		require.ErrorContains(t, err, "no such file or directory")
	})
}
