package iwmo_test

import (
	"testing"

	"github.com/hyperized/iwmo"
)

func BenchmarkEncode(b *testing.B) {
	msgs := []struct {
		name string
		msg  iwmo.Message
	}{
		{"WMO301", iwmo.ValidWMO301()},
		{"WMO302", iwmo.ValidWMO302()},
		{"WMO303", iwmo.ValidWMO303()},
		{"WMO304", iwmo.ValidWMO304()},
		{"WMO305", iwmo.ValidWMO305()},
		{"WMO315", iwmo.ValidWMO315()},
	}
	for _, tt := range msgs {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				_, err := iwmo.Encode(tt.msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	msgs := []struct {
		name string
		msg  iwmo.Message
	}{
		{"WMO301", iwmo.ValidWMO301()},
		{"WMO302", iwmo.ValidWMO302()},
		{"WMO303", iwmo.ValidWMO303()},
		{"WMO304", iwmo.ValidWMO304()},
		{"WMO305", iwmo.ValidWMO305()},
		{"WMO315", iwmo.ValidWMO315()},
	}
	for _, tt := range msgs {
		data, err := iwmo.Encode(tt.msg)
		if err != nil {
			b.Fatalf("Encode(%s): %v", tt.name, err)
		}

		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				_, err := iwmo.Decode(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeAs(b *testing.B) {
	benchDecodeAs[iwmo.WMO301](b, "WMO301", iwmo.ValidWMO301())
	benchDecodeAs[iwmo.WMO302](b, "WMO302", iwmo.ValidWMO302())
	benchDecodeAs[iwmo.WMO303](b, "WMO303", iwmo.ValidWMO303())
	benchDecodeAs[iwmo.WMO304](b, "WMO304", iwmo.ValidWMO304())
	benchDecodeAs[iwmo.WMO305](b, "WMO305", iwmo.ValidWMO305())
	benchDecodeAs[iwmo.WMO315](b, "WMO315", iwmo.ValidWMO315())
}

func benchDecodeAs[T any](b *testing.B, name string, msg iwmo.Message) {
	b.Helper()

	data, err := iwmo.Encode(msg)
	if err != nil {
		b.Fatalf("Encode(%s): %v", name, err)
	}

	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, err := iwmo.DecodeAs[T](data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
