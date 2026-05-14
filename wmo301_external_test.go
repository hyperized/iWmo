package iwmo_test

import (
	"errors"
	"os"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestWMO301_MessageType(t *testing.T) {
	m := &iwmo.WMO301{}
	if got := m.MessageType(); got != iwmo.MessageTypeWMO301 {
		t.Errorf("MessageType() = %q, want %q", got, iwmo.MessageTypeWMO301)
	}
}

func TestWMO301_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *iwmo.WMO301
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     iwmo.ValidWMO301(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Header.BerichtCode = "302"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "no clients",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Bsn = "123456789" // fails elfproef
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing Achternaam",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Geboortedatum = "not-a-date"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Geslacht",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Geslacht = "7" // not in allow-list
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_VALUE",
		},
		{
			name: "no toewijzingen",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Ingangsdatum",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Ingangsdatum",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = "bad"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "Einddatum before Ingangsdatum",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = "2026-12-31"
				m.Clienten[0].Toewijzingen[0].Einddatum = "2026-01-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "invalid Einddatum format",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Einddatum = "31-12-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Product.Categorie",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Toewijzingsdatum format",
			msg: func() *iwmo.WMO301 {
				m := iwmo.ValidWMO301()
				m.Clienten[0].Toewijzingen[0].Toewijzingsdatum = "01-04-2026"
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
					t.Errorf("expected error code %q in validation errors: %v", tt.errCode, ve)
				}
			}
		})
	}
}

func TestWMO301_MarshalUnmarshal(t *testing.T) {
	original := iwmo.ValidWMO301()
	data, err := iwmo.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := iwmo.DecodeAs[iwmo.WMO301](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO301]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "301" {
		t.Errorf("BerichtCode = %q, want 301", decoded.Header.BerichtCode)
	}
	if len(decoded.Clienten) != 1 {
		t.Fatalf("len(Clienten) = %d, want 1", len(decoded.Clienten))
	}
	cl := decoded.Clienten[0]
	if cl.Bsn != iwmo.ValidBSN {
		t.Errorf("Bsn = %q, want %q", cl.Bsn, iwmo.ValidBSN)
	}
	if cl.Naam.Achternaam != "Janssen" {
		t.Errorf("Achternaam = %q, want Janssen", cl.Naam.Achternaam)
	}
	if len(cl.Toewijzingen) != 1 {
		t.Fatalf("len(Toewijzingen) = %d, want 1", len(cl.Toewijzingen))
	}
	tw := cl.Toewijzingen[0]
	if tw.ToewijzingNummer != "12345" {
		t.Errorf("ToewijzingNummer = %q, want 12345", tw.ToewijzingNummer)
	}
	if tw.Product.Categorie != "03" {
		t.Errorf("Product.Categorie = %q, want 03", tw.Product.Categorie)
	}
}

func TestWMO301_DecodeFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo301_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO301](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO301]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO301_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo301_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO301](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO301]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}

func TestWMO301_Validate_MultipleClients(t *testing.T) {
	m := iwmo.ValidWMO301()
	m.Clienten = append(m.Clienten, iwmo.WMO301Client{
		Bsn:           "900212640",
		Naam:          iwmo.Naam{Achternaam: "De Vries"},
		Geboortedatum: "1990-06-15",
		Toewijzingen: []iwmo.Toewijzing{
			{
				ToewijzingNummer: "67890",
				Product:          iwmo.Product{Categorie: "03", Code: "H533"},
				Ingangsdatum:     "2026-06-01",
				Einddatum:        "2026-12-31",
			},
		},
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for valid multi-client message", err)
	}
}

func TestWMO301_Validate_MultipleClientsOneInvalid(t *testing.T) {
	m := iwmo.ValidWMO301()
	m.Clienten = append(m.Clienten, iwmo.WMO301Client{
		Bsn:  "000000000", // invalid BSN
		Naam: iwmo.Naam{Achternaam: ""},
	})
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error for invalid second client")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	// Second client should have errors referencing Client[1]
	foundClient1 := false
	for _, e := range ve {
		if len(e.Field) > 9 && e.Field[:9] == "Client[1]" {
			foundClient1 = true
			break
		}
	}
	if !foundClient1 {
		t.Errorf("expected errors for Client[1], got: %v", ve)
	}
}

func TestWMO301_Validate_EmptyMessage(t *testing.T) {
	m := &iwmo.WMO301{}
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil for empty WMO301")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	// Should have header errors + WRONG_CODE + REQUIRED (no clients)
	if len(ve) < 3 {
		t.Errorf("expected at least 3 validation errors for empty WMO301, got %d: %v", len(ve), ve)
	}
}

func TestWMO301_Validate_ToewijzingWithOptionalOmvang(t *testing.T) {
	m := iwmo.ValidWMO301()
	m.Clienten[0].Toewijzingen[0].Omvang = &iwmo.Omvang{
		Volume: "8", Eenheid: "uur", Frequentie: "week",
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for toewijzing with Omvang", err)
	}
}

func TestWMO301_Validate_EinddatumWithInvalidIngangsdatum(t *testing.T) {
	// When Ingangsdatum is invalid and Einddatum is valid, validation
	// produces both INVALID_DATE for Ingangsdatum and INVALID_PERIOD
	// because ValidatePeriod fails on the unparseable Ingangsdatum.
	m := iwmo.ValidWMO301()
	m.Clienten[0].Toewijzingen[0].Ingangsdatum = "bad"
	m.Clienten[0].Toewijzingen[0].Einddatum = "2026-12-31"
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error for invalid Ingangsdatum")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	foundIngangsdatumErr := false
	foundPeriodErr := false
	for _, e := range ve {
		if e.Code == "INVALID_DATE" && e.Field == "Client[0].Toewijzing[0].Ingangsdatum" {
			foundIngangsdatumErr = true
		}
		if e.Code == "INVALID_PERIOD" {
			foundPeriodErr = true
		}
	}
	if !foundIngangsdatumErr {
		t.Error("expected INVALID_DATE error for Ingangsdatum")
	}
	if !foundPeriodErr {
		t.Error("expected INVALID_PERIOD when Ingangsdatum is unparseable and Einddatum is valid")
	}
}
