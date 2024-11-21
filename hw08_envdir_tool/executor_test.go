package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple run", func(t *testing.T) {
		ret := RunCmd([]string{"echo", "hello"}, nil)
		require.Zero(t, ret)
	})
}
