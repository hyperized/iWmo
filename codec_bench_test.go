package iwmo

import "testing"

func BenchmarkEncode(b *testing.B) {
	msgs := []struct {
		name string
		msg  Message
	}{
		{"WMO301", validWMO301()},
		{"WMO302", validWMO302()},
		{"WMO303", validWMO303()},
		{"WMO304", validWMO304()},
		{"WMO305", validWMO305()},
		{"WMO315", validWMO315()},
	}
	for _, tt := range msgs {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				_, err := Encode(tt.msg)
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
		msg  Message
	}{
		{"WMO301", validWMO301()},
		{"WMO302", validWMO302()},
		{"WMO303", validWMO303()},
		{"WMO304", validWMO304()},
		{"WMO305", validWMO305()},
		{"WMO315", validWMO315()},
	}
	for _, tt := range msgs {
		data, err := Encode(tt.msg)
		if err != nil {
			b.Fatalf("Encode(%s): %v", tt.name, err)
		}

		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				_, err := Decode(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeAs(b *testing.B) {
	benchDecodeAs[WMO301](b, "WMO301", validWMO301())
	benchDecodeAs[WMO302](b, "WMO302", validWMO302())
	benchDecodeAs[WMO303](b, "WMO303", validWMO303())
	benchDecodeAs[WMO304](b, "WMO304", validWMO304())
	benchDecodeAs[WMO305](b, "WMO305", validWMO305())
	benchDecodeAs[WMO315](b, "WMO315", validWMO315())
}

func benchDecodeAs[T any](b *testing.B, name string, msg Message) {
	b.Helper()

	data, err := Encode(msg)
	if err != nil {
		b.Fatalf("Encode(%s): %v", name, err)
	}

	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, err := DecodeAs[T](data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
