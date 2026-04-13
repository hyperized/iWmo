package iwmo

import (
	"encoding/xml"
	"strconv"
)

// WMO304 represents a Retourbericht (acknowledgement/return message). It is
// sent in response to any of the other iWMO message types in either direction
// and contains one or more return codes that indicate acceptance or rejection.
type WMO304 struct {
	XMLName     xml.Name     `xml:"Bericht"`
	Header      WMO304Header `xml:"Header"`
	RetourCodes []RetourCode `xml:"RetourCode,omitempty"`
}

// MessageType returns "WMO304".
func (m *WMO304) MessageType() string { return "WMO304" }

// Validate checks all header fields and that at least one RetourCode is present.
func (m *WMO304) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header.Header)...)

	if m.Header.BerichtCode != "304" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 304 for WMO304",
		})
	}

	if len(m.RetourCodes) == 0 {
		errs = append(errs, ValidationError{
			Field: "RetourCode", Code: "REQUIRED",
			Message: "at least one RetourCode element is required",
		})
	}

	for i, rc := range m.RetourCodes {
		if rc.Code == "" {
			errs = append(errs, ValidationError{
				Field: "RetourCode[" + strconv.Itoa(i) + "].Code", Code: "REQUIRED",
				Message: "RetourCode.Code is required",
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// WMO304Header extends [Header] with fields that identify the original message
// being acknowledged.
type WMO304Header struct {
	Header

	// GerefereerdBerichtCode is the BerichtCode of the message being acknowledged.
	GerefereerdBerichtCode string `xml:"GerefereerdBerichtCode,omitempty"`
	// GerefereerdBerichtIdentificatie is the BerichtIdentificatie of the
	// message being acknowledged.
	GerefereerdBerichtIdentificatie string `xml:"GerefereerdBerichtIdentificatie,omitempty"`
}

// RetourCode carries a status code and optional description for one aspect of
// the acknowledged message.
type RetourCode struct {
	// Code is the machine-readable return code (e.g. "0000" for success).
	Code string `xml:"Code"`
	// Omschrijving is the optional human-readable description.
	Omschrijving string `xml:"Omschrijving,omitempty"`
}
