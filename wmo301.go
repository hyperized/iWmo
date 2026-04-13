package iwmo

import (
	"encoding/xml"
	"fmt"
)

// WMO301 represents a Toewijzing (care assignment) message sent from a
// municipality (gemeente) to a care provider (zorgaanbieder).
type WMO301 struct {
	XMLName  xml.Name       `xml:"Bericht"`
	Header   Header         `xml:"Header"`
	Clienten []WMO301Client `xml:"Client"`
}

// MessageType returns "WMO301".
func (m *WMO301) MessageType() string { return "WMO301" }

// Validate checks all header fields and every client and assignment record.
func (m *WMO301) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != "301" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 301 for WMO301",
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
		errs = append(errs, validateWMO301Client(f, cl)...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateWMO301Client(prefix string, cl WMO301Client) ValidationErrors {
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

	if len(cl.Toewijzingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Toewijzing", Code: "REQUIRED",
			Message: "at least one Toewijzing element is required",
		})
	}

	for j, tw := range cl.Toewijzingen {
		tf := fmt.Sprintf("%s.Toewijzing[%d]", prefix, j)
		errs = append(errs, validateToewijzing(tf, tw)...)
	}

	return errs
}

func validateToewijzing(prefix string, tw Toewijzing) ValidationErrors {
	var errs ValidationErrors

	if tw.ToewijzingNummer == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".ToewijzingNummer", Code: "REQUIRED",
			Message: "ToewijzingNummer is required",
		})
	}

	if tw.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: "REQUIRED",
			Message: "Product.Categorie is required",
		})
	}

	if tw.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: "REQUIRED",
			Message: "Product.Code is required",
		})
	}

	if tw.Toewijzingsdatum != "" && !ValidateDate(tw.Toewijzingsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Toewijzingsdatum", Code: "INVALID_DATE",
			Message: "Toewijzingsdatum must be formatted YYYY-MM-DD",
		})
	}

	if tw.Ingangsdatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: "REQUIRED",
			Message: "Ingangsdatum is required",
		})
	} else if !ValidateDate(tw.Ingangsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: "INVALID_DATE",
			Message: "Ingangsdatum must be formatted YYYY-MM-DD",
		})
	}

	if tw.Einddatum != "" {
		if !ValidateDate(tw.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "INVALID_DATE",
				Message: "Einddatum must be formatted YYYY-MM-DD",
			})
		} else if tw.Ingangsdatum != "" && !ValidatePeriod(tw.Ingangsdatum, tw.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "INVALID_PERIOD",
				Message: "Einddatum must be on or after Ingangsdatum",
			})
		}
	}

	return errs
}

// WMO301Client is a client record within a WMO301 message.
type WMO301Client struct {
	// Bsn is the client's Burgerservicenummer (9-digit Dutch social security number).
	Bsn string `xml:"Bsn"`
	// Naam holds the client's name components.
	Naam Naam `xml:"Naam"`
	// Geboortedatum is the client's date of birth in YYYY-MM-DD format.
	Geboortedatum string `xml:"Geboortedatum,omitempty"`
	// Geslacht is the client's gender code (0=unknown, 1=male, 2=female, 9=unspecified).
	Geslacht string `xml:"Geslacht,omitempty"`
	// Toewijzingen contains one or more care assignment records.
	Toewijzingen []Toewijzing `xml:"Toewijzing"`
}

// Toewijzing is a single care assignment record within a WMO301 message.
type Toewijzing struct {
	// ToewijzingNummer is the unique assignment number assigned by the gemeente.
	ToewijzingNummer string `xml:"ToewijzingNummer"`
	// ReferentieAanbieder is the care provider's own reference for this assignment.
	ReferentieAanbieder string `xml:"ReferentieAanbieder,omitempty"`
	// Product identifies the care product.
	Product Product `xml:"Product"`
	// Toewijzingsdatum is the date the assignment was made (YYYY-MM-DD).
	Toewijzingsdatum string `xml:"Toewijzingsdatum,omitempty"`
	// Ingangsdatum is the start date of the care period (YYYY-MM-DD).
	Ingangsdatum string `xml:"Ingangsdatum"`
	// Einddatum is the optional end date of the care period (YYYY-MM-DD).
	Einddatum string `xml:"Einddatum,omitempty"`
	// Omvang specifies volume, unit, and frequency if applicable.
	Omvang *Omvang `xml:"Omvang,omitempty"`
	// RedenWijziging is the reason-for-change code.
	RedenWijziging string `xml:"RedenWijziging,omitempty"`
	// Commentaar is a free-text remark.
	Commentaar string `xml:"Commentaar,omitempty"`
}
