package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 3; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
		for i := 0; i < 3; i++ {
			v, _ := c.Get(Key(strconv.Itoa(i)))
			require.Equal(t, i, v)
		}
		c.Clear()
		for i := 0; i < 3; i++ {
			v, _ := c.Get(Key(strconv.Itoa(i)))
			require.Nil(t, v)
		}
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestCacheDisplacement(t *testing.T) {
	t.Run("simple displacement", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 30; i++ {
			c.Set(Key(strconv.Itoa(i)), i) // [27, 28, 29]
		}

		for i := 0; i < 27; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
			require.Nil(t, val)
		}

		for i := 27; i < 30; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}
	})

	t.Run("displacement oldest", func(t *testing.T) {
		c := NewCache(3)
		for i := 2; i >= 0; i-- {
			c.Set(Key(strconv.Itoa(i)), i) // [0, 1, 2]
		}

		for i := 0; i < 10; i++ {
			if i%2 == 0 {
				c.Get("1")
			} else {
				c.Get("2")
			}
		} // [1, 2, 0]
		c.Set("3", 3) // [3, 1, 2]

		val, ok := c.Get("0")
		require.False(t, ok)
		require.Nil(t, val)

		for i := 1; i <= 3; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}
	})
}

func BenchmarkCacheGet(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			cache := NewCache(size)
			for i := 0; i < size; i++ {
				cache.Set(Key(strconv.Itoa(i)), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cache.Get(Key(strconv.Itoa(rand.Intn(100))))
			}
		})
	}
}

func BenchmarkCacheSet(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			cache := NewCache(size)
			for i := 0; i < size; i++ {
				cache.Set(Key(strconv.Itoa(i)), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v := rand.Intn(100)
				k := Key(strconv.Itoa(v))
				cache.Set(k, v)
			}
		})
	}
}

func BenchmarkCacheClear(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			cache := NewCache(size)
			for i := 0; i < size; i++ {
				cache.Set(Key(strconv.Itoa(i)), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cache.Clear()
			}
		})
	}
}
