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
	// XML has the opening <BerichtCode> tag but no closing tag; sniffBerichtCode
	// returns "" and Decode falls back to full unmarshal which fails.
	data := []byte(`<?xml version="1.0"?><Bericht><Header><BerichtCode>301`)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode() error = nil, want error")
	}
	if !errors.Is(err, ErrInvalidMessage) {
		t.Errorf("errors.Is(err, ErrInvalidMessage) = false, got: %v", err)
	}
}

// helpers shared across test files

func validHeader(code string) Header {
	return Header{
		BerichtCode:          code,
		BerichtVersie:        "3.2",
		Afzender:             "0363",
		Ontvanger:            "12345678",
		BerichtIdentificatie: "MSG-TEST-001",
		DagtekeningBericht:   "2026-04-12",
	}
}

const validBSN = "123456782" // elfproef: sum=154, 154%11=0

func validWMO301() *WMO301 {
	return &WMO301{
		Header: validHeader("301"),
		Clienten: []WMO301Client{
			{
				Bsn:           validBSN,
				Naam:          Naam{Voornamen: "Jan", Tussenvoegsels: "van", Achternaam: "Janssen"},
				Geboortedatum: "1980-01-15",
				Toewijzingen: []Toewijzing{
					{
						ToewijzingNummer: "12345",
						Product:          Product{Categorie: "03", Code: "H532"},
						Ingangsdatum:     "2026-05-01",
						Einddatum:        "2026-12-31",
					},
				},
			},
		},
	}
}

func validWMO302() *WMO302 {
	return &WMO302{
		Header: validHeader("302"),
		Clienten: []WMO302Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				VerzoekToewijzingen: []VerzoekToewijzing{
					{
						ReferentieAanbieder: "REF-001",
						Product:             Product{Categorie: "03", Code: "H532"},
						Ingangsdatum:        "2026-05-01",
					},
				},
			},
		},
	}
}

func validWMO303() *WMO303 {
	return &WMO303{
		Header: validHeader("303"),
		Clienten: []WMO303Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Declaratieperiode: Declaratieperiode{
					Begindatum: "2026-04-01",
					Einddatum:  "2026-04-30",
				},
				Prestaties: []Prestatie{
					{
						ToewijzingNummer: "12345",
						Product:          Product{Categorie: "03", Code: "H532"},
						Begindatum:       "2026-04-01",
						Einddatum:        "2026-04-30",
						Omvang:           Omvang{Volume: "32", Eenheid: "uur", Frequentie: "maand"},
						Bedrag:           "1600.00",
					},
				},
			},
		},
	}
}

func validWMO304() *WMO304 {
	return &WMO304{
		Header: WMO304Header{
			Header: validHeader("304"),
			GerefereerdBerichtCode:          "302",
			GerefereerdBerichtIdentificatie: "MSG-TEST-001",
		},
		RetourCodes: []RetourCode{
			{Code: "0000", Omschrijving: "Bericht in goede orde ontvangen"},
		},
	}
}

func validWMO305() *WMO305 {
	return &WMO305{
		Header: validHeader("305"),
		Clienten: []WMO305Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Mutaties: []Mutatie{
					{
						ToewijzingNummer: "12345",
						Mutatiedatum:     "2026-04-12",
						Mutatiecode:      "01",
						Begindatum:       "2026-05-01",
					},
				},
			},
		},
	}
}

func validWMO315() *WMO315 {
	return &WMO315{
		Header: validHeader("315"),
		Clienten: []WMO315Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Statusmeldingen: []StatusmeldingRecord{
					{
						ToewijzingNummer: "12345",
						StatusCode:       "01",
						StatusDatum:      "2026-05-01",
						Commentaar:       "Zorg gestart",
					},
				},
			},
		},
	}
}
