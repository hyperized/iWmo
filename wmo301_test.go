package iwmo

import (
	"errors"
	"os"
	"testing"
)

func TestWMO301_MessageType(t *testing.T) {
	m := &WMO301{}
	if got := m.MessageType(); got != "WMO301" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO301")
	}
}

func TestWMO301_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO301
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     validWMO301(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Header.BerichtCode = "302"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "no clients",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Bsn = "123456789" // fails elfproef
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing Achternaam",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Geboortedatum = "not-a-date"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "no toewijzingen",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Ingangsdatum",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Ingangsdatum",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = "bad"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "Einddatum before Ingangsdatum",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Ingangsdatum = "2026-12-31"
				m.Clienten[0].Toewijzingen[0].Einddatum = "2026-01-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "invalid Einddatum format",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Einddatum = "31-12-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Product.Categorie",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code",
			msg: func() *WMO301 {
				m := validWMO301()
				m.Clienten[0].Toewijzingen[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Toewijzingsdatum format",
			msg: func() *WMO301 {
				m := validWMO301()
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
				var ve ValidationErrors
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
	original := validWMO301()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO301](data)
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
	if cl.Bsn != validBSN {
		t.Errorf("Bsn = %q, want %q", cl.Bsn, validBSN)
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
	msg, err := DecodeAs[WMO301](data)
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
	msg, err := DecodeAs[WMO301](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO301]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}
