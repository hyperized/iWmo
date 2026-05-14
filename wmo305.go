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
func (m *WMO305) MessageType() string { return MessageTypeWMO305 }

// Validate checks all header fields and every client and mutatie record.
func (m *WMO305) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != berichtCodeWMO305 {
		errs = append(errs, ValidationError{
			Field: msgFieldHeaderBerichtCode, Code: codeWrongCode,
			Message: "BerichtCode moet 305 zijn voor WMO305",
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

	if len(cl.Mutaties) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatie", Code: codeRequired,
			Message: "ten minste één Mutatie-element is verplicht",
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
			Field: prefix + ".ToewijzingNummer", Code: codeRequired,
			Message: msgToewijzingNummerRequired,
		})
	}

	if mu.Mutatiedatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiedatum", Code: codeRequired,
			Message: "Mutatiedatum is verplicht",
		})
	} else if !ValidateDate(mu.Mutatiedatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiedatum", Code: codeInvalidDate,
			Message: "Mutatiedatum moet de notatie JJJJ-MM-DD hebben",
		})
	}

	switch mu.Mutatiecode {
	case MutatiecodeStart:
		if mu.Begindatum == "" {
			errs = append(errs, ValidationError{
				Field: prefix + ".Begindatum", Code: codeRequired,
				Message: "Begindatum is verplicht bij start (Mutatiecode 01)",
			})
		}
	case MutatiecodeWijziging:
		if mu.Product == nil {
			errs = append(errs, ValidationError{
				Field: prefix + ".Product", Code: codeRequired,
				Message: "Product is verplicht bij wijziging (Mutatiecode 02)",
			})
		}
	case MutatiecodeStop:
		if mu.Einddatum == "" {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeRequired,
				Message: "Einddatum is verplicht bij stop (Mutatiecode 03)",
			})
		}
	case "":
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiecode", Code: codeRequired,
			Message: "Mutatiecode is verplicht",
		})
	default:
		errs = append(errs, ValidationError{
			Field: prefix + ".Mutatiecode", Code: codeInvalidValue,
			Message: "Mutatiecode moet 01 (start), 02 (wijziging) of 03 (stop) zijn",
		})
	}

	if mu.Begindatum != "" && !ValidateDate(mu.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: codeInvalidDate,
			Message: msgBegindatumFormat,
		})
	}

	if mu.Einddatum != "" && !ValidateDate(mu.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: codeInvalidDate,
			Message: msgEinddatumFormat,
		})
	}

	if mu.Begindatum != "" && mu.Einddatum != "" && !ValidatePeriod(mu.Begindatum, mu.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: codeInvalidPeriod,
			Message: msgEinddatumAfterBegindatum,
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
