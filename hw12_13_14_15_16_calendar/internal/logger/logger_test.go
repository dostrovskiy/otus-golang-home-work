package logger

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("simple logging", func(t *testing.T) {
		out := os.Stdout
		defer func() { os.Stdout = out }()

		r, w, err := os.Pipe()
		if err != nil {
			t.Errorf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		logger := New("WARN")
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		w.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		if err != nil {
			t.Errorf("failed to read from pipe: %v", err)
		}

		require.Equal(t, "warn message\nerror message\n", buf.String())
	})
}
