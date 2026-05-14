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
func (m *WMO315) MessageType() string { return MessageTypeWMO315 }

// Validate checks all header fields and every client and statusmelding record.
func (m *WMO315) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != berichtCodeWMO315 {
		errs = append(errs, ValidationError{
			Field: msgFieldHeaderBerichtCode, Code: codeWrongCode,
			Message: "BerichtCode moet 315 zijn voor WMO315",
		})
	}

	if len(m.Clienten) == 0 {
		errs = append(errs, ValidationError{
			Field: msgFieldClient, Code: codeRequired,
			Message: msgClientRequired,
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
			Field: prefix + ".Bsn", Code: codeInvalidBSN,
			Message: msgBSNInvalid,
		})
	}

	if cl.Naam.Achternaam == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Naam.Achternaam", Code: codeRequired,
			Message: msgAchternaamRequired,
		})
	}

	if cl.Geboortedatum != "" && !ValidateDate(cl.Geboortedatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Geboortedatum", Code: codeInvalidDate,
			Message: msgGeboortedatumFormat,
		})
	}

	if len(cl.Statusmeldingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Statusmelding", Code: codeRequired,
			Message: "ten minste één Statusmelding-element is verplicht",
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
			Field: prefix + ".ToewijzingNummer", Code: codeRequired,
			Message: msgToewijzingNummerRequired,
		})
	}

	if sm.StatusCode == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusCode", Code: codeRequired,
			Message: "StatusCode is verplicht",
		})
	}

	if sm.StatusDatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusDatum", Code: codeRequired,
			Message: "StatusDatum is verplicht",
		})
	} else if !ValidateDate(sm.StatusDatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".StatusDatum", Code: codeInvalidDate,
			Message: "StatusDatum moet de notatie JJJJ-MM-DD hebben",
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
