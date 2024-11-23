package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("simple read env from dir", func(t *testing.T) {
		envs, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}, envs)
	})

	t.Run("special read env from dir", func(t *testing.T) {
		envs, err := ReadDir("testdata/envspec")
		require.NoError(t, err)
		require.Equal(t, Environment{
			"NBSP": EnvValue{"spaces and tabs should be removed at the end of the line," +
				" but e.g. nbsp is allowed: \u00A0\u00A0", false},
		}, envs)
	})

	t.Run("error env name", func(t *testing.T) {
		_, err := ReadDir("testdata/enverr")
		require.Truef(t, errors.Is(err, ErrInvalidEnvName), "actual error %q", err)
	})
}
