package iwmo

import (
	"errors"
	"testing"
)

func TestValidateBSN(t *testing.T) {
	tests := []struct {
		bsn  string
		want bool
	}{
		// Valid BSNs (elfproef passes)
		{"123456782", true},  // sum=154, 154%11=0
		{"900212640", true},  // sum=132, 132%11=0
		{"987654321", false}, // sum=283, 283%11=8 → invalid
		// Invalid: wrong length
		{"12345678", false},   // 8 digits
		{"1234567820", false}, // 10 digits
		{"", false},
		// Invalid: non-digit characters
		{"12345678a", false},
		{"12345 782", false},
		// Invalid: all zeros (sum = 0, not > 0)
		{"000000000", false},
		// Invalid: wrong check digit
		{"123456789", false}, // sum=147, 147%11=4
		{"999999999", false}, // sum=9*(-1+sum of 9*(9+8+7+6+5+4+3+2))=... invalid
	}
	for _, tt := range tests {
		t.Run(tt.bsn, func(t *testing.T) {
			if got := ValidateBSN(tt.bsn); got != tt.want {
				t.Errorf("ValidateBSN(%q) = %v, want %v", tt.bsn, got, tt.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"2026-04-12", true},
		{"2000-01-01", true},
		{"1900-12-31", true},
		{"", false},
		{"2026-13-01", false}, // month 13
		{"2026-04-31", false}, // April has 30 days
		{"04-12-2026", false}, // wrong format
		{"2026/04/12", false},
		{"20260412", false},
		{"not-a-date", false},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := ValidateDate(tt.s); got != tt.want {
				t.Errorf("ValidateDate(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestValidatePeriod(t *testing.T) {
	tests := []struct {
		begin string
		end   string
		want  bool
	}{
		{"2026-01-01", "2026-12-31", true},
		{"2026-04-12", "2026-04-12", true}, // same day is valid
		{"2026-12-31", "2026-01-01", false}, // end before begin
		{"2026-01-01", "not-a-date", false},
		{"not-a-date", "2026-01-01", false},
		{"", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.begin+"_"+tt.end, func(t *testing.T) {
			if got := ValidatePeriod(tt.begin, tt.end); got != tt.want {
				t.Errorf("ValidatePeriod(%q, %q) = %v, want %v", tt.begin, tt.end, got, tt.want)
			}
		})
	}
}

func TestValidateHeader(t *testing.T) {
	validHeader := Header{
		BerichtCode:          "301",
		BerichtVersie:        "3.2",
		Afzender:             "0363",
		Ontvanger:            "12345678",
		BerichtIdentificatie: "MSG-001",
		DagtekeningBericht:   "2026-04-12",
	}

	tests := []struct {
		name    string
		modify  func(Header) Header
		wantErr bool
	}{
		{
			name:    "valid header",
			modify:  func(h Header) Header { return h },
			wantErr: false,
		},
		{
			name:    "missing BerichtCode",
			modify:  func(h Header) Header { h.BerichtCode = ""; return h },
			wantErr: true,
		},
		{
			name:    "missing BerichtVersie",
			modify:  func(h Header) Header { h.BerichtVersie = ""; return h },
			wantErr: true,
		},
		{
			name:    "missing Afzender",
			modify:  func(h Header) Header { h.Afzender = ""; return h },
			wantErr: true,
		},
		{
			name:    "missing Ontvanger",
			modify:  func(h Header) Header { h.Ontvanger = ""; return h },
			wantErr: true,
		},
		{
			name:    "missing BerichtIdentificatie",
			modify:  func(h Header) Header { h.BerichtIdentificatie = ""; return h },
			wantErr: true,
		},
		{
			name:    "missing DagtekeningBericht",
			modify:  func(h Header) Header { h.DagtekeningBericht = ""; return h },
			wantErr: true,
		},
		{
			name:    "invalid DagtekeningBericht format",
			modify:  func(h Header) Header { h.DagtekeningBericht = "12-04-2026"; return h },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.modify(validHeader)
			errs := validateHeader(h)
			if (len(errs) > 0) != tt.wantErr {
				t.Errorf("validateHeader errors=%v, wantErr=%v", errs, tt.wantErr)
			}
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	ve := ValidationErrors{
		{Field: "Foo", Code: "REQUIRED", Message: "Foo is required"},
		{Field: "Bar", Code: "INVALID", Message: "Bar is invalid"},
	}
	got := ve.Error()
	want := "field Foo (REQUIRED): Foo is required; field Bar (INVALID): Bar is invalid"
	if got != want {
		t.Errorf("ValidationErrors.Error() = %q, want %q", got, want)
	}
}

func TestValidationErrors_As(t *testing.T) {
	ve := ValidationErrors{
		{Field: "Bsn", Code: "INVALID_BSN", Message: "invalid BSN"},
	}
	var got ValidationErrors
	if !errors.As(ve, &got) {
		t.Fatal("errors.As returned false for ValidationErrors")
	}
	if len(got) != 1 {
		t.Errorf("len(got) = %d, want 1", len(got))
	}
}
