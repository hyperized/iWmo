package iwmo

import (
	"errors"
	"os"
	"testing"
)

func TestWMO303_MessageType(t *testing.T) {
	m := &WMO303{}
	if got := m.MessageType(); got != "WMO303" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO303")
	}
}

func TestWMO303_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO303
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     validWMO303(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "missing Achternaam",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Bsn = "123456789"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "declaratieperiode end before begin",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = "2026-04-30"
				m.Clienten[0].Declaratieperiode.Einddatum = "2026-04-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "missing Bedrag",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Bedrag = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "prestatie period reversed",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = "2026-04-30"
				m.Clienten[0].Prestaties[0].Einddatum = "2026-04-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "no clients",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Geboortedatum = "15-01-1980"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "no prestaties",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Categorie in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Begindatum in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Begindatum format in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = "01-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Einddatum in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Einddatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Einddatum format in prestatie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Einddatum = "30-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Omvang.Volume",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Volume = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Omvang.Eenheid",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Eenheid = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Omvang.Frequentie",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Frequentie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing declaratieperiode Begindatum",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid declaratieperiode Begindatum format",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = "01-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing declaratieperiode Einddatum",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Declaratieperiode.Einddatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid declaratieperiode Einddatum format",
			msg: func() *WMO303 {
				m := validWMO303()
				m.Clienten[0].Declaratieperiode.Einddatum = "30-04-2026"
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

func TestWMO303_MarshalUnmarshal(t *testing.T) {
	original := validWMO303()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := DecodeAs[WMO303](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO303]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "303" {
		t.Errorf("BerichtCode = %q, want 303", decoded.Header.BerichtCode)
	}
	if len(decoded.Clienten) != 1 {
		t.Fatalf("len(Clienten) = %d, want 1", len(decoded.Clienten))
	}
	cl := decoded.Clienten[0]
	if cl.Declaratieperiode.Begindatum != "2026-04-01" {
		t.Errorf("Declaratieperiode.Begindatum = %q, want 2026-04-01", cl.Declaratieperiode.Begindatum)
	}
	if len(cl.Prestaties) != 1 {
		t.Fatalf("len(Prestaties) = %d, want 1", len(cl.Prestaties))
	}
	if cl.Prestaties[0].Bedrag != "1600.00" {
		t.Errorf("Bedrag = %q, want 1600.00", cl.Prestaties[0].Bedrag)
	}
}

func TestWMO303_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo303_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := DecodeAs[WMO303](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO303]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}
