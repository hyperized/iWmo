package iwmo

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestWMO304_MessageType(t *testing.T) {
	m := &WMO304{}
	if got := m.MessageType(); got != "WMO304" {
		t.Errorf("MessageType() = %q, want %q", got, "WMO304")
	}
}

func TestWMO304_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *WMO304
		wantErr bool
		errCode string
	}{
		{
			name:    "valid message",
			msg:     validWMO304(),
			wantErr: false,
		},
		{
			name: "wrong BerichtCode",
			msg: func() *WMO304 {
				m := validWMO304()
				m.Header.BerichtCode = "301"
				return m
			}(),
			wantErr: true,
			errCode: "WRONG_CODE",
		},
		{
			name: "no retour codes",
			msg: func() *WMO304 {
				m := validWMO304()
				m.RetourCodes = nil
				return m
			}(),
			wantErr: true,
			errCode: "REQUIRED",
		},
		{
			name: "retour code with empty Code field",
			msg: func() *WMO304 {
				m := validWMO304()
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

func TestWMO304_MarshalUnmarshal(t *testing.T) {
	original := validWMO304()
	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	// Verify extended header fields are encoded.
	if !bytes.Contains(data, []byte("<GerefereerdBerichtCode>302</GerefereerdBerichtCode>")) {
		t.Error("encoded output missing GerefereerdBerichtCode")
	}

	decoded, err := DecodeAs[WMO304](data)
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
	msg, err := DecodeAs[WMO304](data)
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
	msg, err := DecodeAs[WMO304](data)
	if err != nil {
		t.Fatalf("DecodeAs[WMO304]() error = %v", err)
	}
	if err := msg.Validate(); err == nil {
		t.Error("Validate() returned nil for invalid message, want error")
	}
}
