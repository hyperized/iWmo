package iwmo_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestWMO304_MessageType(t *testing.T) {
	m := &iwmo.WMO304{}
	if got := m.MessageType(); got != iwmo.MessageTypeWMO304 {
		t.Errorf("MessageType() = %q, want %q", got, iwmo.MessageTypeWMO304)
	}
}

func TestWMO304_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *iwmo.WMO304
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     iwmo.ValidWMO304(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *iwmo.WMO304 {
				m := iwmo.ValidWMO304()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "no retour codes",
			msg: func() *iwmo.WMO304 {
				m := iwmo.ValidWMO304()
				m.RetourCodes = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "retour code with empty Code field",
			msg: func() *iwmo.WMO304 {
				m := iwmo.ValidWMO304()
				m.RetourCodes[0].Code = ""
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
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

func TestWMO304_MarshalUnmarshal(t *testing.T) {
	original := iwmo.ValidWMO304()
	data, err := iwmo.Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	// Verify extended header fields are encoded.
	if !bytes.Contains(data, []byte("<GerefereerdBerichtCode>302</GerefereerdBerichtCode>")) {
		t.Error("encoded output missing GerefereerdBerichtCode")
	}

	decoded, err := iwmo.DecodeAs[iwmo.WMO304](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO304]() error = %v", err)
	}
	if decoded.Header.BerichtCode != "304" {
		t.Errorf("BerichtCode = %q, want 304", decoded.Header.BerichtCode)
	}
	if decoded.Header.GerefereerdBerichtCode != "302" {
		t.Errorf("GerefereerdBerichtCode = %q, want 302", decoded.Header.GerefereerdBerichtCode)
	}
	if len(decoded.RetourCodes) != 1 {
		t.Fatalf("len(RetourCodes) = %d, want 1", len(decoded.RetourCodes))
	}
	if decoded.RetourCodes[0].Code != "0000" {
		t.Errorf("RetourCodes[0].Code = %q, want 0000", decoded.RetourCodes[0].Code)
	}
}

func TestWMO304_FromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo304_valid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO304](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO304]() error = %v", err)
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestWMO304_InvalidFile_FailsValidation(t *testing.T) {
	data, err := os.ReadFile("testdata/wmo304_invalid.xml")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	msg, err := iwmo.DecodeAs[iwmo.WMO304](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO304]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}

func TestWMO304_Validate_MultipleRetourCodes(t *testing.T) {
	m := iwmo.ValidWMO304()
	m.RetourCodes = append(m.RetourCodes, iwmo.RetourCode{
		Code:         "1001",
		Omschrijving: "Waarschuwing: clientgegevens wijken af",
	})
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for multiple valid RetourCodes", err)
	}
}

func TestWMO304_Validate_MultipleRetourCodesOneEmpty(t *testing.T) {
	m := iwmo.ValidWMO304()
	m.RetourCodes = append(m.RetourCodes, iwmo.RetourCode{Code: ""})
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error for empty RetourCode.Code")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	found := false
	for _, e := range ve {
		if e.Field == "RetourCode[1].Code" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for RetourCode[1].Code, got: %v", ve)
	}
}

func TestWMO304_Validate_EmptyMessage(t *testing.T) {
	m := &iwmo.WMO304{}
	err := m.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil for empty WMO304")
	}
	var ve iwmo.ValidationErrors
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(ve) < 3 {
		t.Errorf("expected at least 3 validation errors for empty WMO304, got %d: %v", len(ve), ve)
	}
}

func TestWMO304_Validate_WithoutGerefereerdFields(t *testing.T) {
	// GerefereerdBerichtCode and GerefereerdBerichtIdentificatie are optional.
	m := &iwmo.WMO304{
		Header: iwmo.WMO304Header{
			Header: iwmo.ValidHeaderFixture("304"),
		},
		RetourCodes: []iwmo.RetourCode{
			{Code: "0000", Omschrijving: "OK"},
		},
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil (gerefereerd fields are optional)", err)
	}
}
