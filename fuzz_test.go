package iwmo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hyperized/iwmo"
)

// FuzzDecode exercises the XML parser against arbitrary byte sequences. The
// only contract: Decode must never panic, regardless of input. Seeded with
// every fixture in testdata/ so coverage starts from valid-shaped XML.
func FuzzDecode(f *testing.F) {
	matches, err := filepath.Glob("testdata/*.xml")
	if err != nil {
		f.Fatalf("Glob testdata/*.xml: %v", err)
	}
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			f.Fatalf("ReadFile %s: %v", p, err)
		}
		f.Add(data)
	}
	// Minimal seeds covering the error branches in Decode.
	f.Add([]byte(""))
	f.Add([]byte("not xml"))
	f.Add([]byte(`<?xml version="1.0"?><Bericht><Header><BerichtCode>999</BerichtCode></Header></Bericht>`))

	f.Fuzz(func(_ *testing.T, data []byte) {
		// Result is intentionally discarded; we only care that no input causes a panic.
		_, _ = iwmo.Decode(data)
	})
}

// FuzzValidateBSN exercises the elfproef implementation against arbitrary
// strings, ensuring it never panics on malformed UTF-8, control characters,
// or length boundaries.
func FuzzValidateBSN(f *testing.F) {
	seeds := []string{
		"123456782", // valid
		"900212640", // valid
		"123456789", // wrong checksum
		"000000000", // all zeros
		"",          // empty
		"abcdefghi", // letters
		"12345678a", // mixed
		"12345678٠", // non-ASCII digit
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(_ *testing.T, s string) {
		_ = iwmo.ValidateBSN(s)
	})
}
