package iwmo_test

import (
	"errors"
	"os"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestWMO303_MessageType(t *testing.T) {
	m := &iwmo.WMO303{}
	if got := m.MessageType(); got != iwmo.MessageTypeWMO303 {
		t.Errorf("MessageType() = %q, want %q", got, iwmo.MessageTypeWMO303)
	}
}

func TestWMO303_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *iwmo.WMO303
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     iwmo.ValidWMO303(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "missing Achternaam",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Naam.Achternaam = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid BSN",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Bsn = "123456789"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_BSN",
		},
		{
			name: "invalid Geslacht",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Geslacht = "5"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_VALUE",
		},
		{
			name: "declaratieperiode end before begin",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = "2026-04-30"
				m.Clienten[0].Declaratieperiode.Einddatum = "2026-04-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "missing Bedrag",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Bedrag = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "prestatie period reversed",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = "2026-04-30"
				m.Clienten[0].Prestaties[0].Einddatum = "2026-04-01"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_PERIOD",
		},
		{
			name: "no clients",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Geboortedatum",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Geboortedatum = "15-01-1980"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "no prestaties",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing ToewijzingNummer in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].ToewijzingNummer = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Categorie in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Product.Categorie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Product.Code in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Product.Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Begindatum in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Begindatum format in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Begindatum = "01-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Einddatum in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Einddatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid Einddatum format in prestatie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Einddatum = "30-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing Omvang.Volume",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Volume = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Omvang.Eenheid",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Eenheid = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing Omvang.Frequentie",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Prestaties[0].Omvang.Frequentie = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "missing declaratieperiode Begindatum",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid declaratieperiode Begindatum format",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Declaratieperiode.Begindatum = "01-04-2026"
				return m
			}(),
			wantErr: true,
			errCode: "INVALID_DATE",
		},
		{
			name: "missing declaratieperiode Einddatum",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
				m.Clienten[0].Declaratieperiode.Einddatum = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "invalid declaratieperiode Einddatum format",
			msg: func() *iwmo.WMO303 {
				m := iwmo.ValidWMO303()
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

func TestWMO303_MarshalUnmarshal(t *testing.T) {
	original := iwmo.ValidWMO303()
	data, err := iwmo.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	decoded, err := iwmo.DecodeAs[iwmo.WMO303](data)
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
	msg, err := iwmo.DecodeAs[iwmo.WMO303](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO303]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO303_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo303_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO303](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO303]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}

func TestWMO303_Validate_MultipleClients(t *testing.T) {
	m := iwmo.ValidWMO303()
	m.Clienten = append(m.Clienten, iwmo.WMO303Client{
		Bsn:  "900212640",
		Naam: iwmo.Naam{Achternaam: "De Vries"},
		Declaratieperiode: iwmo.Declaratieperiode{
			Begindatum: "2026-05-01",
			Einddatum:  "2026-05-31",
		},
		Prestaties: []iwmo.Prestatie{
			{
				ToewijzingNummer: "67890",
				Product:          iwmo.Product{Categorie: "03", Code: "H533"},
				Begindatum:       "2026-05-01",
				Einddatum:        "2026-05-31",
				Omvang:           iwmo.Omvang{Volume: "16", Eenheid: "uur", Frequentie: "maand"},
				Bedrag:           "800.00",
			},
		},
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for valid multi-client message", err)
	}
}

func TestWMO303_Validate_EmptyMessage(t *testing.T) {
	m := &iwmo.WMO303{}
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil for empty WMO303")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(ve) < 3 {
		t.Errorf("expected at least 3 validation errors for empty WMO303, got %d: %v", len(ve), ve)
	}
}

func TestWMO303_Validate_MultiplePrestaties(t *testing.T) {
	m := iwmo.ValidWMO303()
	m.Clienten[0].Prestaties = append(m.Clienten[0].Prestaties, iwmo.Prestatie{
		ToewijzingNummer: "67890",
		Product:          iwmo.Product{Categorie: "03", Code: "H533"},
		Begindatum:       "2026-04-01",
		Einddatum:        "2026-04-15",
		Omvang:           iwmo.Omvang{Volume: "8", Eenheid: "uur", Frequentie: "maand"},
		Bedrag:           "400.00",
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for message with multiple prestaties", err)
	}
}

func TestWMO303_Validate_BothDeclatieperiodeDatesInvalid(t *testing.T) {
	m := iwmo.ValidWMO303()
	m.Clienten[0].Declaratieperiode.Begindatum = "bad"
	m.Clienten[0].Declaratieperiode.Einddatum = "also-bad"
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want errors for both invalid dates")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	dateErrors := 0
	for _, e := range ve {
		if e.Code == "INVALID_DATE" {
			dateErrors++
		}
	}
	if dateErrors < 2 {
		t.Errorf("expected at least 2 INVALID_DATE errors, got %d", dateErrors)
	}
}

func TestWMO303_Validate_PrestatieWithInvalidBegindatum(t *testing.T) {
	// When Begindatum is invalid and Einddatum is valid, validation
	// produces both INVALID_DATE for Begindatum and INVALID_PERIOD
	// because ValidatePeriod fails on the unparseable Begindatum.
	m := iwmo.ValidWMO303()
	m.Clienten[0].Prestaties[0].Begindatum = "bad"
	m.Clienten[0].Prestaties[0].Einddatum = "2026-04-30"
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error for invalid Begindatum")
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
		t.Error("expected INVALID_DATE for invalid Begindatum")
	}
	if !foundPeriod {
		t.Error("expected INVALID_PERIOD when Begindatum is unparseable and Einddatum is valid")
	}
}
