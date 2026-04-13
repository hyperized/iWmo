package iwmo

import (
	"errors"
	"os"
	"testing"
)

func TestWMO315_MessageType(t *testing.T) {
	m := &WMO315{}
	if got := m.MessageType(); got != "WMO315" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO315")
	}
}

func TestWMO315_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO315
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     validWMO315(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "missing Achternaam",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Bsn = "123456789"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing StatusCode",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusCode = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing StatusDatum",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusDatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid StatusDatum",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusDatum = "12-05-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "no statusmeldingen",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Statusmeldingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Statusmeldingen[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "no clients",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *WMO315 {
				m := validWMO315()
				m.Clienten[0].Geboortedatum = "15-01-1980"
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

func TestWMO315_MarshalUnmarshal(t *testing.T) {
	original := validWMO315()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO315](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO315]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "315" {
		t.Errorf("BerichtCode = %q, want 315", decoded.Header.BerichtCode)
	}
	if len(decoded.Clienten) != 1 {
		t.Fatalf("len(Clienten) = %d, want 1", len(decoded.Clienten))
	}
	sm := decoded.Clienten[0].Statusmeldingen[0]
	if sm.StatusCode != "01" {
		t.Errorf("StatusCode = %q, want 01", sm.StatusCode)
	}
	if sm.Commentaar != "Zorg gestart" {
		t.Errorf("Commentaar = %q, want Zorg gestart", sm.Commentaar)
	}
}

func TestWMO315_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo315_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := DecodeAs[WMO315](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO315]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}
