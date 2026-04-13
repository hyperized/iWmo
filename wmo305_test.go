package iwmo

import (
	"errors"
	"os"
	"testing"
)

func TestWMO305_MessageType(t *testing.T) {
	m := &WMO305{}
	if got := m.MessageType(); got != "WMO305" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO305")
	}
}

func TestWMO305_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO305
		wantErr bool
		errCode string
	}{
		{
			name:    "valid start mutatie",
			msg:     validWMO305(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "missing Achternaam",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "no mutaties for client",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer in mutatie",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Bsn = "000000000"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing Mutatiedatum",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiedatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Mutatiecode",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "99"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_VALUE",
		},
		{
			name: "missing Mutatiecode",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "valid wijziging mutatie (code 02)",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "02"
				m.Clienten[0].Mutaties[0].Product = &Product{Categorie: "03", Code: "H532"}
				return m
			}(),
			wantErr: false,
		},
		{
			name: "valid stop mutatie (code 03)",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "03"
				m.Clienten[0].Mutaties[0].Einddatum = "2026-12-31"
				return m
			}(),
			wantErr: false,
		},
		{
			name: "code 01 missing Begindatum",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "01"
				m.Clienten[0].Mutaties[0].Begindatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "code 02 missing Product",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "02"
				m.Clienten[0].Mutaties[0].Product = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "code 03 missing Einddatum",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "03"
				m.Clienten[0].Mutaties[0].Einddatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "period reversed",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Begindatum = "2026-12-31"
				m.Clienten[0].Mutaties[0].Einddatum = "2026-01-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "no clients",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Geboortedatum = "15-01-1980"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Mutatiedatum format",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiedatum = "12-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Begindatum format",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Begindatum = "01-05-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "invalid Einddatum format",
			msg: func() *WMO305 {
				m := validWMO305()
				m.Clienten[0].Mutaties[0].Mutatiecode = "03"
				m.Clienten[0].Mutaties[0].Einddatum = "31-12-2026"
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

func TestWMO305_MarshalUnmarshal(t *testing.T) {
	original := validWMO305()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO305](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO305]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "305" {
		t.Errorf("BerichtCode = %q, want 305", decoded.Header.BerichtCode)
	}
	if len(decoded.Clienten) != 1 {
		t.Fatalf("len(Clienten) = %d, want 1", len(decoded.Clienten))
	}
	mu := decoded.Clienten[0].Mutaties[0]
	if mu.Mutatiecode != "01" {
		t.Errorf("Mutatiecode = %q, want 01", mu.Mutatiecode)
	}
}

func TestWMO305_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo305_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := DecodeAs[WMO305](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO305]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO305_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo305_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := DecodeAs[WMO305](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO305]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}
