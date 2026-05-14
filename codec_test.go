package iwmo

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"
)

// badMsg implements Message but causes xml.Marshal to fail by returning an
// error from MarshalXML. Used to exercise the error path in Encode.
type badMsg struct{}

func (b *badMsg) MessageType() string                                       { return "BAD" }
func (b *badMsg) Validate() error                                           { return nil }
func (b *badMsg) MarshalXML(_ *xml.Encoder, _ xml.StartElement) error {
	return errors.New("intentional xml marshal error")
}

func TestEncode_WMO301(t *testing.T) {
	msg := validWMO301()
	data, err := Encode(msg)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if !bytes.Contains(data, []byte("<?xml")) {
		t.Error("encoded output missing XML declaration")
	}
	if !bytes.Contains(data, []byte("<Bericht>")) {
		t.Error("encoded output missing <Bericht> root element")
	}
	if !bytes.Contains(data, []byte("<BerichtCode>301</BerichtCode>")) {
		t.Error("encoded output missing BerichtCode 301")
	}
}

func TestDecodeAs_WMO301_RoundTrip(t *testing.T) {
	original := validWMO301()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO301](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO301]() error = %v", err)
	}
	if decoded.Header.BerichtCode != original.Header.BerichtCode {
		t.Errorf("BerichtCode = %q, want %q", decoded.Header.BerichtCode, original.Header.BerichtCode)
	}
	if decoded.Header.BerichtIdentificatie != original.Header.BerichtIdentificatie {
		t.Errorf("BerichtIdentificatie = %q, want %q",
			decoded.Header.BerichtIdentificatie, original.Header.BerichtIdentificatie)
	}
	if len(decoded.Clienten) != len(original.Clienten) {
		t.Errorf("len(Clienten) = %d, want %d", len(decoded.Clienten), len(original.Clienten))
	}
	if len(decoded.Clienten) > 0 {
		if decoded.Clienten[0].Bsn != original.Clienten[0].Bsn {
			t.Errorf("Clienten[0].Bsn = %q, want %q",
				decoded.Clienten[0].Bsn, original.Clienten[0].Bsn)
		}
	}
}

func TestDecode_ByBerichtCode(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		want string
	}{
		{"WMO301", validWMO301(), "WMO301"},
		{"WMO302", validWMO302(), "WMO302"},
		{"WMO303", validWMO303(), "WMO303"},
		{"WMO304", validWMO304(), "WMO304"},
		{"WMO305", validWMO305(), "WMO305"},
		{"WMO315", validWMO315(), "WMO315"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Encode(tt.msg)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			got, err := Decode(data)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			if got.MessageType() != tt.want {
				t.Errorf("MessageType() = %q, want %q", got.MessageType(), tt.want)
			}
		})
	}
}

func TestDecode_UnknownCode(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Bericht>
  <Header>
    <BerichtCode>999</BerichtCode>
    <BerichtVersie>3.2</BerichtVersie>
    <Afzender>0363</Afzender>
    <Ontvanger>12345678</Ontvanger>
    <BerichtIdentificatie>MSG-001</BerichtIdentificatie>
    <DagtekeningBericht>2026-04-12</DagtekeningBericht>
  </Header>
</Bericht>`
	_, err := Decode([]byte(xml))
	if err == nil {
		t.Fatal("Decode() error = nil, want ErrUnknownMessage")
	}
	if !errors.Is(err, ErrUnknownMessage) {
		t.Errorf("errors.Is(err, ErrUnknownMessage) = false, got: %v", err)
	}
}

func TestDecode_MalformedXML(t *testing.T) {
	_, err := Decode([]byte("not xml at all"))
	if err == nil {
		t.Fatal("Decode() error = nil, want error for malformed XML")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestDecodeAs_MalformedXML(t *testing.T) {
	_, err := DecodeAs[WMO301]([]byte("not xml"))
	if err == nil {
		t.Fatal("DecodeAs[WMO301]() error = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestEncode_MarshalError(t *testing.T) {
	_, err := Encode(&badMsg{})
	if err == nil {
		t.Fatal("Encode() error = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestDecode_TruncatedBerichtCodeTag(t *testing.T) {
	// Truncated XML: the <BerichtCode> element is never closed, so
	// xml.Unmarshal fails and Decode returns ErrInvalidMessage.
	data := []byte(`<?xml version="1.0"?><Bericht><Header><BerichtCode>301`)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode() error = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestDecode_EmptyBerichtCode(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Bericht>
  <Header>
    <BerichtCode></BerichtCode>
  </Header>
</Bericht>`)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode() error = nil, want ErrUnknownMessage for empty BerichtCode")
	}
	if !errors.Is(err, ErrUnknownMessage) {
		t.Errorf("errors.Is(err, ErrUnknownMessage) = false, got: %v", err)
	}
}

func TestDecode_MissingBerichtCode(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Bericht>
  <Header>
    <BerichtVersie>3.2</BerichtVersie>
  </Header>
</Bericht>`)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode() error = nil, want ErrUnknownMessage for missing BerichtCode")
	}
	if !errors.Is(err, ErrUnknownMessage) {
		t.Errorf("errors.Is(err, ErrUnknownMessage) = false, got: %v", err)
	}
}

func TestDecode_EmptyXML(t *testing.T) {
	_, err := Decode([]byte(""))
	if err == nil {
		t.Fatal("Decode() error = nil, want error for empty input")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

func TestEncode_AllMessageTypes(t *testing.T) {
	// Verify Encode succeeds for each message type and output contains
	// the XML declaration and correct BerichtCode.
	tests := []struct {
		name string
		msg  Message
		code string
	}{
		{"WMO301", validWMO301(), "301"},
		{"WMO302", validWMO302(), "302"},
		{"WMO303", validWMO303(), "303"},
		{"WMO304", validWMO304(), "304"},
		{"WMO305", validWMO305(), "305"},
		{"WMO315", validWMO315(), "315"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Encode(tt.msg)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			if !bytes.Contains(data, []byte("<?xml")) {
				t.Error("encoded output missing XML declaration")
			}
			if !bytes.Contains(data, []byte("<BerichtCode>"+tt.code+"</BerichtCode>")) {
				t.Errorf("encoded output missing BerichtCode %s", tt.code)
			}
		})
	}
}

func TestDecodeAs_WrongType(t *testing.T) {
	// Encode a WMO301 but decode as WMO302; this should succeed at the XML
	// level but the fields won't match the expected structure.
	data, err := Encode(validWMO301())
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	msg, err := DecodeAs[WMO302](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO302]() error = %v (XML is valid, just wrong type)", err)
	}
	// The decoded message should have BerichtCode "301" (from the XML)
	// even though we decoded it as WMO302.
	if msg.Header.BerichtCode != "301" {
		t.Errorf("BerichtCode = %q, want 301", msg.Header.BerichtCode)
	}
}

// Fixtures (validHeader, validBSN, validWMO301, …) live in fixtures_test.go
// and are re-exported via export_test.go for the external test package.
