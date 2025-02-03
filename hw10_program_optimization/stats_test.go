//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("domain error", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(data), "\\")
		require.Error(t, err)
	})

	t.Run("json error", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString("{"), "gov")
		require.Error(t, err)
	})
}

func BenchmarkGetDomainStat(b *testing.B) {
	getData := func(size int) []byte {
		data := make([]byte, 0)
		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(b, err)
		defer r.Close()
		require.Equal(b, 1, len(r.File))
		
		f, err := r.File[0].Open()
		require.NoError(b, err)

		scanner := bufio.NewScanner(f)
		for i := 0; i < size; i++ {
			if !scanner.Scan() {
				scanner = bufio.NewScanner(f)
			}
			b := append(scanner.Bytes(), []byte("\n")...)
			if i == size-1 {
				b = b[:len(b)-1]
			}
			data = append(data, b...)
		}
		return data
	}

	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.StopTimer()
			data := getData(size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				_, err := GetDomainStat(bytes.NewReader(data), "gov")
				require.NoError(b, err)
			}
		})
	}
}
