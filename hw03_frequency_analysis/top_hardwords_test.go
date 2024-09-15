package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTop10HardWords(t *testing.T) {
	text := "Hello,world! It's an example. It's a... complex wo...rd 'example'.! Multy-hyphen '----' - is a word too."

	t.Run("simple positive test with hard words", func(t *testing.T) {
		expected := []string{
			"It's",         // 2
			"'----'",       // 1
			"'example'.!",  // 1
			"-",            // 1
			"Hello,world!", // 1
			"Multy-hyphen", // 1
			"a",            // 1
			"a...",         // 1
			"an",           // 1
			"complex",      // 1
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("hard positive test with hard words", func(t *testing.T) {
		expected := []string{
			"a",            // 2
			"example",      // 2
			"it's",         // 2
			"----",         // 1
			"an",           // 1
			"complex",      // 1
			"hello,world",  // 1
			"is",           // 1
			"multy-hyphen", // 1
			"too",          // 1
		}
		require.Equal(t, expected, Top10Hard(text, true))
	})
}
