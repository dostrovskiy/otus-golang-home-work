package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("simple copy", func(t *testing.T) {
		err := Copy("testdata/tolstoy_anna_karenina.txt", "part.txt", 1000000, 1000000)
		require.NoError(t, err)
		require.FileExists(t, "part.txt")
		fi, err := os.Stat("part.txt")
		require.NoError(t, err)
		require.Equal(t, int64(1000000), fi.Size(), "actual error %d", fi.Size())
		err = os.Remove("part.txt")
		require.NoError(t, err)
	})

	t.Run("unsupported file error", func(t *testing.T) {
		fi, err := os.Create("zero_sized.txt")
		require.NoError(t, err)
		fi.Close()
		err = Copy("zero_sized.txt", "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
		err = os.Remove("zero_sized.txt")
		require.NoError(t, err)
	})

	t.Run("offset exceeds file size error", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 6618, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("no such file error", func(t *testing.T) {
		err := Copy("no_such_file.txt", "out.txt", 0, 0)
		require.Error(t, err)
	})
}
