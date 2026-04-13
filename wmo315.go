package iwmo

import (
	"encoding/xml"
	"fmt"
)

// WMO315 represents a Statusmelding (status update) message sent from a care
// provider (zorgaanbieder) to a municipality (gemeente). It reports the current
// delivery status of one or more care assignments.
type WMO315 struct {
	XMLName  xml.Name       `xml:"Bericht"`
	Header   Header         `xml:"Header"`
	Clienten []WMO315Client `xml:"Client"`
}

// MessageType returns "WMO315".
func (m *WMO315) MessageType() string { return "WMO315" }

// Validate checks all header fields and every client and statusmelding record.
func (m *WMO315) Validate() error {
	var errs ValidationErrors
	errs = append(errs, validateHeader(m.Header)...)
	if m.Header.BerichtCode != "315" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 315 for WMO315",
		})
	}
	if len(m.Clienten) == 0 {
		errs = append(errs, ValidationError{
			Field: "Client", Code: "REQUIRED",
			Message: "at least one Client element is required",
		})
	}
	for i, cl := range m.Clienten {
		f := fmt.Sprintf("Client[%d]", i)
		errs = append(errs, validateWMO315Client(f, cl)...)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateWMO315Client(prefix string, cl WMO315Client) ValidationErrors {
	var errs ValidationErrors
	if !ValidateBSN(cl.Bsn) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Bsn", Code: "INVALID_BSN",
			Message: "BSN must be 9 digits passing the elfproef (11-check)",
		})
	}
	if cl.Naam.Achternaam == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Naam.Achternaam", Code: "REQUIRED",
			Message: "Achternaam is required",
		})
	}
	if cl.Geboortedatum != "" && !ValidateDate(cl.Geboortedatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Geboortedatum", Code: "INVALID_DATE",
			Message: "Geboortedatum must be formatted YYYY-MM-DD",
		})
	}
	if len(cl.Statusmeldingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Statusmelding", Code: "REQUIRED",
			Message: "at least one Statusmelding element is required",
		})
	}
	for j, sm := range cl.Statusmeldingen {
		sf := fmt.Sprintf("%s.Statusmelding[%d]", prefix, j)
		errs = append(errs, validateStatusmeldingRecord(sf, sm)...)
	}
	return errs
}

func validateStatusmeldingRecord(prefix string, sm StatusmeldingRecord) ValidationErrors {
	var errs ValidationErrors
	if sm.ToewijzingNummer == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".ToewijzingNummer", Code: "REQUIRED",
			Message: "ToewijzingNummer is required",
		})
	}
	if sm.StatusCode == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusCode", Code: "REQUIRED",
			Message: "StatusCode is required",
		})
	}
	if sm.StatusDatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusDatum", Code: "REQUIRED",
			Message: "StatusDatum is required",
		})
	} else if !ValidateDate(sm.StatusDatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusDatum", Code: "INVALID_DATE",
			Message: "StatusDatum must be formatted YYYY-MM-DD",
		})
	}
	return errs
}

// WMO315Client is a client record within a WMO315 message.
type WMO315Client struct {
	// Bsn is the client's Burgerservicenummer.
	Bsn string `xml:"Bsn"`
	// Naam holds the client's name components.
	Naam Naam `xml:"Naam"`
	// Geboortedatum is the client's date of birth in YYYY-MM-DD format.
	Geboortedatum string `xml:"Geboortedatum,omitempty"`
	// Statusmeldingen contains one or more status update records.
	Statusmeldingen []StatusmeldingRecord `xml:"Statusmelding"`
}

// StatusmeldingRecord is a single status update record within a WMO315 message.
type StatusmeldingRecord struct {
	// ToewijzingNummer references the care assignment being reported on.
	ToewijzingNummer string `xml:"ToewijzingNummer"`
	// StatusCode is the status code indicating the delivery state.
	StatusCode string `xml:"StatusCode"`
	// StatusDatum is the date of this status update (YYYY-MM-DD).
	StatusDatum string `xml:"StatusDatum"`
	// Commentaar is an optional free-text remark.
	Commentaar string `xml:"Commentaar,omitempty"`
}
