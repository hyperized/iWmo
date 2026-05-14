package iwmo_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestHeader_MarshalUnmarshal(t *testing.T) {
	// Wrap Header in a struct with a root element for round-trip testing.
	type wrapper struct {
		XMLName xml.Name    `xml:"Root"`
		Header  iwmo.Header `xml:"Header"`
	}
	original := wrapper{
		Header: iwmo.Header{
			BerichtCode:          "301",
			BerichtVersie:        "3.2",
			Afzender:             "0363",
			Ontvanger:            "12345678",
			BerichtIdentificatie: "MSG-001",
			DagtekeningBericht:   "2026-04-12",
			XsltVersie:           "1.0",
			XsdVersie:            "3.2",
		},
	}
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("xml.Marshal error = %v", err)
	}
	var decoded wrapper
	if err = xml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal error = %v", err)
	}
	h := decoded.Header
	if h.BerichtCode != "301" {
		t.Errorf("BerichtCode = %q, want %q", h.BerichtCode, "301")
	}
	if h.BerichtVersie != "3.2" {
		t.Errorf("BerichtVersie = %q, want %q", h.BerichtVersie, "3.2")
	}
	if h.Afzender != "0363" {
		t.Errorf("Afzender = %q, want %q", h.Afzender, "0363")
	}
	if h.Ontvanger != "12345678" {
		t.Errorf("Ontvanger = %q, want %q", h.Ontvanger, "12345678")
	}
	if h.BerichtIdentificatie != "MSG-001" {
		t.Errorf("BerichtIdentificatie = %q, want %q", h.BerichtIdentificatie, "MSG-001")
	}
	if h.DagtekeningBericht != "2026-04-12" {
		t.Errorf("DagtekeningBericht = %q, want %q", h.DagtekeningBericht, "2026-04-12")
	}
}

// TestHeader_XsltAndXsdVersionsRoundTrip ensures the optional XsltVersie and
// XsdVersie fields survive an encode/decode cycle. These fields are required
// for spec-compliant integration testing with some gemeente endpoints.
func TestHeader_XsltAndXsdVersionsRoundTrip(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name    `xml:"Root"`
		Header  iwmo.Header `xml:"Header"`
	}
	original := wrapper{
		Header: iwmo.Header{
			BerichtCode:          "301",
			BerichtVersie:        "3.2",
			Afzender:             "0363",
			Ontvanger:            "12345678",
			BerichtIdentificatie: "MSG-001",
			DagtekeningBericht:   "2026-04-12",
			XsltVersie:           "iwmo-3.2-xslt",
			XsdVersie:            "iwmo-3.2-xsd",
		},
	}
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("xml.Marshal: %v", err)
	}
	var decoded wrapper
	if err := xml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal: %v", err)
	}
	if decoded.Header.XsltVersie != "iwmo-3.2-xslt" {
		t.Errorf("XsltVersie = %q, want iwmo-3.2-xslt", decoded.Header.XsltVersie)
	}
	if decoded.Header.XsdVersie != "iwmo-3.2-xsd" {
		t.Errorf("XsdVersie = %q, want iwmo-3.2-xsd", decoded.Header.XsdVersie)
	}
}

func TestNaam_MarshalUnmarshal(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name  `xml:"Root"`
		Naam    iwmo.Naam `xml:"Naam"`
	}
	original := wrapper{
		Naam: iwmo.Naam{
			Voornamen:      "Jan",
			Tussenvoegsels: "van",
			Achternaam:     "Janssen",
		},
	}
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("xml.Marshal error = %v", err)
	}
	var decoded wrapper
	if err = xml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal error = %v", err)
	}
	n := decoded.Naam
	if n.Voornamen != "Jan" {
		t.Errorf("Voornamen = %q, want %q", n.Voornamen, "Jan")
	}
	if n.Tussenvoegsels != "van" {
		t.Errorf("Tussenvoegsels = %q, want %q", n.Tussenvoegsels, "van")
	}
	if n.Achternaam != "Janssen" {
		t.Errorf("Achternaam = %q, want %q", n.Achternaam, "Janssen")
	}
}

func TestNaam_OmitemptyTussenvoegsels(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name  `xml:"Root"`
		Naam    iwmo.Naam `xml:"Naam"`
	}
	n := wrapper{Naam: iwmo.Naam{Achternaam: "Smit"}}
	data, err := xml.Marshal(n)
	if err != nil {
		t.Fatalf("xml.Marshal error = %v", err)
	}
	// Tussenvoegsels is omitempty, so it should not appear in output.
	xmlStr := string(data)
	if strings.Contains(xmlStr, "<Tussenvoegsels>") {
		t.Error("expected <Tussenvoegsels> to be omitted when empty")
	}
}

func TestProduct_MarshalUnmarshal(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name     `xml:"Root"`
		Product iwmo.Product `xml:"Product"`
	}
	original := wrapper{Product: iwmo.Product{Categorie: "03", Code: "H532"}}
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("xml.Marshal error = %v", err)
	}
	var decoded wrapper
	if err = xml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal error = %v", err)
	}
	if decoded.Product.Categorie != "03" {
		t.Errorf("Categorie = %q, want %q", decoded.Product.Categorie, "03")
	}
	if decoded.Product.Code != "H532" {
		t.Errorf("Code = %q, want %q", decoded.Product.Code, "H532")
	}
}

func TestOmvang_MarshalUnmarshal(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name    `xml:"Root"`
		Omvang  iwmo.Omvang `xml:"Omvang"`
	}
	original := wrapper{Omvang: iwmo.Omvang{Volume: "8", Eenheid: "uur", Frequentie: "week"}}
	data, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("xml.Marshal error = %v", err)
	}
	var decoded wrapper
	if err = xml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("xml.Unmarshal error = %v", err)
	}
	o := decoded.Omvang
	if o.Volume != "8" {
		t.Errorf("Volume = %q, want %q", o.Volume, "8")
	}
	if o.Eenheid != "uur" {
		t.Errorf("Eenheid = %q, want %q", o.Eenheid, "uur")
	}
	if o.Frequentie != "week" {
		t.Errorf("Frequentie = %q, want %q", o.Frequentie, "week")
	}
}
