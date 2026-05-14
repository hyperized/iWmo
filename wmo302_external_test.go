package iwmo_test

import (
	"errors"
	"os"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestWMO302_MessageType(t *testing.T) {
	m := &iwmo.WMO302{}
	if got := m.MessageType(); got != iwmo.MessageTypeWMO302 {
		t.Errorf("MessageType() = %q, want %q", got, iwmo.MessageTypeWMO302)
	}
}

func TestWMO302_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *iwmo.WMO302
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     iwmo.ValidWMO302(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "invalid BSN",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].Bsn = "000000001"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing Achternaam",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "no VerzoekToewijzingen",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ReferentieAanbieder",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].ReferentieAanbieder = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "Einddatum before Ingangsdatum",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = "2026-12-31"
				m.Clienten[0].VerzoekToewijzingen[0].Einddatum = "2026-01-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "no clients",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].Geboortedatum = "15-01-1980"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Geslacht",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].Geslacht = "X"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_VALUE",
		},
		{
			name: "missing Product.Categorie",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Ingangsdatum",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Ingangsdatum format",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = "01-05-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Einddatum format",
			msg: func() *iwmo.WMO302 {
				m := iwmo.ValidWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Einddatum = "31-12-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				var ve iwmo.ValidationErrors
				if !errors.As(err, &ve) {
					t.Fatalf("expected ValidationErrors, got %T", err)
				}
				found := false
				for _, e := range ve {
					if e.Code == tt.errCode {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error code %q in: %v", tt.errCode, ve)
				}
			}
		})
	}
}

func TestWMO302_MarshalUnmarshal(t *testing.T) {
	original := iwmo.ValidWMO302()
	data, err := iwmo.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := iwmo.DecodeAs[iwmo.WMO302](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO302]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "302" {
		t.Errorf("BerichtCode = %q, want 302", decoded.Header.BerichtCode)
	}
	if len(decoded.Clienten) != 1 {
		t.Fatalf("len(Clienten) = %d, want 1", len(decoded.Clienten))
	}
	vt := decoded.Clienten[0].VerzoekToewijzingen[0]
	if vt.ReferentieAanbieder != "REF-001" {
		t.Errorf("ReferentieAanbieder = %q, want REF-001", vt.ReferentieAanbieder)
	}
}

func TestWMO302_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo302_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO302](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO302]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO302_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo302_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO302](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO302]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}

func TestWMO302_Validate_MultipleClients(t *testing.T) {
	m := iwmo.ValidWMO302()
	m.Clienten = append(m.Clienten, iwmo.WMO302Client{
		Bsn:  "900212640",
		Naam: iwmo.Naam{Achternaam: "De Vries"},
		VerzoekToewijzingen: []iwmo.VerzoekToewijzing{
			{
				ReferentieAanbieder: "REF-002",
				Product:             iwmo.Product{Categorie: "03", Code: "H533"},
				Ingangsdatum:        "2026-06-01",
			},
		},
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for valid multi-client message", err)
	}
}

func TestWMO302_Validate_EmptyMessage(t *testing.T) {
	m := &iwmo.WMO302{}
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil for empty WMO302")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(ve) < 3 {
		t.Errorf("expected at least 3 validation errors for empty WMO302, got %d: %v", len(ve), ve)
	}
}

func TestWMO302_Validate_EinddatumWithInvalidIngangsdatum(t *testing.T) {
	// When Ingangsdatum is invalid and Einddatum is valid, validation
	// produces both INVALID_DATE for Ingangsdatum and INVALID_PERIOD
	// because ValidatePeriod fails on the unparseable Ingangsdatum.
	m := iwmo.ValidWMO302()
	m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = "bad"
	m.Clienten[0].VerzoekToewijzingen[0].Einddatum = "2026-12-31"
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error for invalid Ingangsdatum")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	foundDate := false
	foundPeriod := false
	for _, e := range ve {
		if e.Code == "INVALID_DATE" {
			foundDate = true
		}
		if e.Code == "INVALID_PERIOD" {
			foundPeriod = true
		}
	}
	if !foundDate {
		t.Error("expected INVALID_DATE for invalid Ingangsdatum")
	}
	if !foundPeriod {
		t.Error("expected INVALID_PERIOD when Ingangsdatum is unparseable and Einddatum is valid")
	}
}

func TestWMO302_Validate_WithOptionalOmvang(t *testing.T) {
	m := iwmo.ValidWMO302()
	m.Clienten[0].VerzoekToewijzingen[0].Omvang = &iwmo.Omvang{
		Volume: "4", Eenheid: "uur", Frequentie: "week",
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for message with Omvang", err)
	}
}

func TestWMO302_Validate_WithValidEinddatum(t *testing.T) {
	m := iwmo.ValidWMO302()
	m.Clienten[0].VerzoekToewijzingen[0].Einddatum = "2026-12-31"
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}
