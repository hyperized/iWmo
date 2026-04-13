package iwmo

import (
	"encoding/xml"
	"fmt"
)

// WMO303 represents a Declaratie (billing declaration) message sent from a
// care provider (zorgaanbieder) to a municipality (gemeente).
type WMO303 struct {
	XMLName  xml.Name       `xml:"Bericht"`
	Header   Header         `xml:"Header"`
	Clienten []WMO303Client `xml:"Client"`
}

// MessageType returns "WMO303".
func (m *WMO303) MessageType() string { return "WMO303" }

// Validate checks all header fields and every client and prestatie record.
func (m *WMO303) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != "303" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 303 for WMO303",
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
		errs = append(errs, validateWMO303Client(f, cl)...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateWMO303Client(prefix string, cl WMO303Client) ValidationErrors {
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

	if cl.Declaratieperiode.Begindatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Begindatum", Code: "REQUIRED",
			Message: "Declaratieperiode.Begindatum is required",
		})
	} else if !ValidateDate(cl.Declaratieperiode.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Begindatum", Code: "INVALID_DATE",
			Message: "Declaratieperiode.Begindatum must be formatted YYYY-MM-DD",
		})
	}

	if cl.Declaratieperiode.Einddatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Einddatum", Code: "REQUIRED",
			Message: "Declaratieperiode.Einddatum is required",
		})
	} else if !ValidateDate(cl.Declaratieperiode.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Einddatum", Code: "INVALID_DATE",
			Message: "Declaratieperiode.Einddatum must be formatted YYYY-MM-DD",
		})
	}

	if cl.Declaratieperiode.Begindatum != "" && cl.Declaratieperiode.Einddatum != "" {
		if !ValidatePeriod(cl.Declaratieperiode.Begindatum, cl.Declaratieperiode.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Declaratieperiode", Code: "INVALID_PERIOD",
				Message: "Declaratieperiode.Einddatum must be on or after Begindatum",
			})
		}
	}

	if len(cl.Prestaties) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Prestatie", Code: "REQUIRED",
			Message: "at least one Prestatie element is required",
		})
	}

	for j, p := range cl.Prestaties {
		pf := fmt.Sprintf("%s.Prestatie[%d]", prefix, j)
		errs = append(errs, validatePrestatie(pf, p)...)
	}

	return errs
}

func validatePrestatie(prefix string, p Prestatie) ValidationErrors {
	var errs ValidationErrors
	if p.ToewijzingNummer == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".ToewijzingNummer", Code: "REQUIRED",
			Message: "ToewijzingNummer is required",
		})
	}

	if p.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: "REQUIRED",
			Message: "Product.Categorie is required",
		})
	}

	if p.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: "REQUIRED",
			Message: "Product.Code is required",
		})
	}

	if p.Begindatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: "REQUIRED",
			Message: "Begindatum is required",
		})
	} else if !ValidateDate(p.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: "INVALID_DATE",
			Message: "Begindatum must be formatted YYYY-MM-DD",
		})
	}

	if p.Einddatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: "REQUIRED",
			Message: "Einddatum is required",
		})
	} else if !ValidateDate(p.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: "INVALID_DATE",
			Message: "Einddatum must be formatted YYYY-MM-DD",
		})
	}

	if p.Begindatum != "" && p.Einddatum != "" {
		if !ValidatePeriod(p.Begindatum, p.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "INVALID_PERIOD",
				Message: "Einddatum must be on or after Begindatum",
			})
		}
	}

	if p.Omvang.Volume == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Volume", Code: "REQUIRED",
			Message: "Omvang.Volume is required",
		})
	}

	if p.Omvang.Eenheid == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Eenheid", Code: "REQUIRED",
			Message: "Omvang.Eenheid is required",
		})
	}

	if p.Omvang.Frequentie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Frequentie", Code: "REQUIRED",
			Message: "Omvang.Frequentie is required",
		})
	}

	if p.Bedrag == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Bedrag", Code: "REQUIRED",
			Message: "Bedrag is required",
		})
	}

	return errs
}

// WMO303Client is a client record within a WMO303 message.
type WMO303Client struct {
	// Bsn is the client's Burgerservicenummer.
	Bsn string `xml:"Bsn"`
	// Naam holds the client's name components.
	Naam Naam `xml:"Naam"`
	// Geboortedatum is the client's date of birth in YYYY-MM-DD format.
	Geboortedatum string `xml:"Geboortedatum,omitempty"`
	// Declaratieperiode is the billing period covered by this declaration.
	Declaratieperiode Declaratieperiode `xml:"Declaratieperiode"`
	// Prestaties contains the individual care delivery records being declared.
	Prestaties []Prestatie `xml:"Prestatie"`
}

// Declaratieperiode is the billing period within a WMO303 client record.
type Declaratieperiode struct {
	// Begindatum is the start of the billing period (YYYY-MM-DD).
	Begindatum string `xml:"Begindatum"`
	// Einddatum is the end of the billing period (YYYY-MM-DD).
	Einddatum string `xml:"Einddatum"`
}

// Prestatie is a single care delivery record being declared in a WMO303 message.
type Prestatie struct {
	// ToewijzingNummer references the original care assignment.
	ToewijzingNummer string `xml:"ToewijzingNummer"`
	// Product identifies the care product delivered.
	Product Product `xml:"Product"`
	// Begindatum is the start of the delivered care period (YYYY-MM-DD).
	Begindatum string `xml:"Begindatum"`
	// Einddatum is the end of the delivered care period (YYYY-MM-DD).
	Einddatum string `xml:"Einddatum"`
	// Omvang specifies the volume, unit, and frequency delivered.
	Omvang Omvang `xml:"Omvang"`
	// Bedrag is the declared amount as a string (e.g. "1600.00").
	Bedrag string `xml:"Bedrag"`
}
