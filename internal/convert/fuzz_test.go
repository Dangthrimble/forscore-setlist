package convert

import (
	"os"
	"path/filepath"
	"testing"
)

func fuzzLimits() limits {
	return limits{
		maxEntries:           100,
		maxInputBytes:        64 << 10,
		maxDecompressedBytes: 64 << 10,
		maxDepth:             16,
		maxStringBytes:       1024,
	}
}

func FuzzFromText(f *testing.F) {
	f.Add([]byte("First\n\nSecond\r\n"))
	f.Add([]byte("A & B <C>\n"))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = fromText(data, "fuzz", fuzzLimits())
	})
}

func FuzzFromNative(f *testing.F) {
	fixture, err := os.ReadFile(filepath.Join("testdata", "2026-05-16 Mount Pleasant Baptist.4ss"))
	if err != nil {
		f.Fatal(err)
	}
	f.Add(fixture)
	f.Add([]byte("not a native setlist"))
	f.Add([]byte{0x1f, 0x8b})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = fromNative(data, "fuzz", fuzzLimits())
	})
}
