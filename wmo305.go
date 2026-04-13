package iwmo

import (
	"encoding/xml"
	"fmt"
)

// WMO305 represents a Mutatie (care delivery change) message sent from a care
// provider (zorgaanbieder) to a municipality (gemeente). A Mutatie reports the
// start, change, or end of care delivery for an assigned client.
type WMO305 struct {
	XMLName  xml.Name       `xml:"Bericht"`
	Header   Header         `xml:"Header"`
	Clienten []WMO305Client `xml:"Client"`
}

// MessageType returns "WMO305".
func (m *WMO305) MessageType() string { return "WMO305" }

// Validate checks all header fields and every client and mutatie record.
func (m *WMO305) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != "305" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 305 for WMO305",
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
		errs = append(errs, validateWMO305Client(f, cl)...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateWMO305Client(prefix string, cl WMO305Client) ValidationErrors {
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

	if len(cl.Mutaties) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatie", Code: "REQUIRED",
			Message: "at least one Mutatie element is required",
		})
	}

	for j, mu := range cl.Mutaties {
		mf := fmt.Sprintf("%s.Mutatie[%d]", prefix, j)
		errs = append(errs, validateMutatie(mf, mu)...)
	}

	return errs
}

func validateMutatie(prefix string, mu Mutatie) ValidationErrors {
	var errs ValidationErrors
	if mu.ToewijzingNummer == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".ToewijzingNummer", Code: "REQUIRED",
			Message: "ToewijzingNummer is required",
		})
	}

	if mu.Mutatiedatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiedatum", Code: "REQUIRED",
			Message: "Mutatiedatum is required",
		})
	} else if !ValidateDate(mu.Mutatiedatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiedatum", Code: "INVALID_DATE",
			Message: "Mutatiedatum must be formatted YYYY-MM-DD",
		})
	}

	switch mu.Mutatiecode {
	case "01":
		if mu.Begindatum == "" {
			errs = append(errs, ValidationError{
				Field: prefix + ".Begindatum", Code: "REQUIRED",
				Message: "Begindatum is required for start (Mutatiecode 01)",
			})
		}
	case "02":
		if mu.Product == nil {
			errs = append(errs, ValidationError{
				Field: prefix + ".Product", Code: "REQUIRED",
				Message: "Product is required for wijziging (Mutatiecode 02)",
			})
		}
	case "03":
		if mu.Einddatum == "" {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "REQUIRED",
				Message: "Einddatum is required for stop (Mutatiecode 03)",
			})
		}
	case "":
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiecode", Code: "REQUIRED",
			Message: "Mutatiecode is required",
		})
	default:
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiecode", Code: "INVALID_VALUE",
			Message: "Mutatiecode must be 01 (start), 02 (wijziging), or 03 (stop)",
		})
	}

	if mu.Begindatum != "" && !ValidateDate(mu.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: "INVALID_DATE",
			Message: "Begindatum must be formatted YYYY-MM-DD",
		})
	}

	if mu.Einddatum != "" && !ValidateDate(mu.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: "INVALID_DATE",
			Message: "Einddatum must be formatted YYYY-MM-DD",
		})
	}

	if mu.Begindatum != "" && mu.Einddatum != "" && !ValidatePeriod(mu.Begindatum, mu.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: "INVALID_PERIOD",
			Message: "Einddatum must be on or after Begindatum",
		})
	}

	return errs
}

// WMO305Client is a client record within a WMO305 message.
type WMO305Client struct {
	// Bsn is the client's Burgerservicenummer.
	Bsn string `xml:"Bsn"`
	// Naam holds the client's name components.
	Naam Naam `xml:"Naam"`
	// Geboortedatum is the client's date of birth in YYYY-MM-DD format.
	Geboortedatum string `xml:"Geboortedatum,omitempty"`
	// Mutaties contains one or more care delivery change records.
	Mutaties []Mutatie `xml:"Mutatie"`
}

// Mutatie is a single care delivery change record within a WMO305 message.
type Mutatie struct {
	// ToewijzingNummer references the original care assignment.
	ToewijzingNummer string `xml:"ToewijzingNummer"`
	// Mutatiedatum is the date the change takes effect (YYYY-MM-DD).
	Mutatiedatum string `xml:"Mutatiedatum"`
	// Mutatiecode indicates the type of change:
	//   "01" = Start (beginning of care delivery)
	//   "02" = Wijziging (change to an existing delivery)
	//   "03" = Stop (end of care delivery)
	Mutatiecode string `xml:"Mutatiecode"`
	// Product identifies the care product, required for Wijziging.
	Product *Product `xml:"Product,omitempty"`
	// Begindatum is the care start date for this mutatie (YYYY-MM-DD).
	Begindatum string `xml:"Begindatum,omitempty"`
	// Einddatum is the care end date for this mutatie (YYYY-MM-DD).
	Einddatum string `xml:"Einddatum,omitempty"`
	// Commentaar is a free-text remark.
	Commentaar string `xml:"Commentaar,omitempty"`
}
