package iwmo

import (
	"errors"
	"os"
	"testing"
)

func TestWMO302_MessageType(t *testing.T) {
	m := &WMO302{}
	if got := m.MessageType(); got != "WMO302" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO302")
	}
}

func TestWMO302_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO302
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     validWMO302(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "invalid BSN",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].Bsn = "000000001"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing Achternaam",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "no VerzoekToewijzingen",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ReferentieAanbieder",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].ReferentieAanbieder = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "Einddatum before Ingangsdatum",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = "2026-12-31"
				m.Clienten[0].VerzoekToewijzingen[0].Einddatum = "2026-01-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "no clients",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].Geboortedatum = "15-01-1980"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Product.Categorie",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Ingangsdatum",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Ingangsdatum format",
			msg: func() *WMO302 {
				m := validWMO302()
				m.Clienten[0].VerzoekToewijzingen[0].Ingangsdatum = "01-05-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Einddatum format",
			msg: func() *WMO302 {
				m := validWMO302()
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
					t.Errorf("expected error code %q in: %v", tt.errCode, ve)
				}
			}
		})
	}
}

func TestWMO302_MarshalUnmarshal(t *testing.T) {
	original := validWMO302()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO302](data)
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
	msg, err := DecodeAs[WMO302](data)
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
	msg, err := DecodeAs[WMO302](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO302]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}
