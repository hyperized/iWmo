package iwmo_test

import (
	"errors"
	"os"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestWMO315_MessageType(t *testing.T) {
	m := &iwmo.WMO315{}
	if got := m.MessageType(); got != iwmo.MessageTypeWMO315 {
		t.Errorf("MessageType() = %q, want %q", got, iwmo.MessageTypeWMO315)
	}
}

func TestWMO315_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *iwmo.WMO315
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     iwmo.ValidWMO315(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "missing Achternaam",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Bsn = "123456789"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "missing StatusCode",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusCode = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing StatusDatum",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusDatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid StatusDatum",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Statusmeldingen[0].StatusDatum = "12-05-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "no statusmeldingen",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Statusmeldingen = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten[0].Statusmeldingen[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "no clients",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *iwmo.WMO315 {
				m := iwmo.ValidWMO315()
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

func TestWMO315_MarshalUnmarshal(t *testing.T) {
	original := iwmo.ValidWMO315()
	data, err := iwmo.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := iwmo.DecodeAs[iwmo.WMO315](data)
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
	msg, err := iwmo.DecodeAs[iwmo.WMO315](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO315]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO315_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo315_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO315](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO315]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}

func TestWMO315_Validate_MultipleClients(t *testing.T) {
	m := iwmo.ValidWMO315()
	m.Clienten = append(m.Clienten, iwmo.WMO315Client{
		Bsn:  "900212640",
		Naam: iwmo.Naam{Achternaam: "De Vries"},
		Statusmeldingen: []iwmo.StatusmeldingRecord{
			{
				ToewijzingNummer: "67890",
				StatusCode:       "02",
				StatusDatum:      "2026-06-01",
			},
		},
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for valid multi-client message", err)
	}
}

func TestWMO315_Validate_EmptyMessage(t *testing.T) {
	m := &iwmo.WMO315{}
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil for empty WMO315")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(ve) < 3 {
		t.Errorf("expected at least 3 validation errors for empty WMO315, got %d: %v", len(ve), ve)
	}
}

func TestWMO315_Validate_MultipleStatusmeldingen(t *testing.T) {
	m := iwmo.ValidWMO315()
	m.Clienten[0].Statusmeldingen = append(m.Clienten[0].Statusmeldingen, iwmo.StatusmeldingRecord{
		ToewijzingNummer: "12345",
		StatusCode:       "02",
		StatusDatum:      "2026-06-15",
		Commentaar:       "Zorg gewijzigd",
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for multiple valid statusmeldingen", err)
	}
}

func TestWMO315_Validate_WithoutCommentaar(t *testing.T) {
	m := iwmo.ValidWMO315()
	m.Clienten[0].Statusmeldingen[0].Commentaar = ""
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil (Commentaar is optional)", err)
	}
}
