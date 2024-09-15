package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestTop10Simple(t *testing.T) {
	text := "Hello, world! This is an example."

	t.Run("simple positive test", func(t *testing.T) {
		expected := []string{
			"Hello,",   // 1
			"This",     // 1
			"an",       // 1
			"example.", // 1
			"is",       // 1
			"world!",   // 1

		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("hard positive test with simple words", func(t *testing.T) {
		expected := []string{
			"an",      // 1
			"example", // 1
			"hello",   // 1
			"is",      // 1
			"this",    // 1
			"world",   // 1
		}
		require.Equal(t, expected, Top10Hard(text, true))
	})
}

func TestTop10SimpleCyrillic(t *testing.T) {
	text := "Всем привет! Это пример текста."

	t.Run("simple positive test with cyrillic", func(t *testing.T) {
		expected := []string{
			"Всем",    // 1
			"Это",     // 1
			"привет!", // 1
			"пример",  // 1
			"текста.", // 1
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("hard positive test with cyrillic", func(t *testing.T) {
		expected := []string{
			"всем",   // 1
			"привет", // 1
			"пример", // 1
			"текста", // 1
			"это",    // 1
		}
		require.Equal(t, expected, Top10Hard(text, true))
	})
}
