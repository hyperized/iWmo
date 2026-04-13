package iwmo

import (
	"encoding/xml"
	"fmt"
)

// WMO302 represents a Verzoek om Toewijzing (care assignment request) message
// sent from a care provider (zorgaanbieder) to a municipality (gemeente).
type WMO302 struct {
	XMLName  xml.Name       `xml:"Bericht"`
	Header   Header         `xml:"Header"`
	Clienten []WMO302Client `xml:"Client"`
}

// MessageType returns "WMO302".
func (m *WMO302) MessageType() string { return "WMO302" }

// Validate checks all header fields and every client and request record.
func (m *WMO302) Validate() error {
	var errs ValidationErrors
	errs = append(errs, validateHeader(m.Header)...)
	if m.Header.BerichtCode != "302" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "WRONG_CODE",
			Message: "BerichtCode must be 302 for WMO302",
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
		errs = append(errs, validateWMO302Client(f, cl)...)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateWMO302Client(prefix string, cl WMO302Client) ValidationErrors {
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
	if len(cl.VerzoekToewijzingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".VerzoekToewijzing", Code: "REQUIRED",
			Message: "at least one VerzoekToewijzing element is required",
		})
	}
	for j, vt := range cl.VerzoekToewijzingen {
		vf := fmt.Sprintf("%s.VerzoekToewijzing[%d]", prefix, j)
		errs = append(errs, validateVerzoekToewijzing(vf, vt)...)
	}
	return errs
}

func validateVerzoekToewijzing(prefix string, vt VerzoekToewijzing) ValidationErrors {
	var errs ValidationErrors
	if vt.ReferentieAanbieder == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".ReferentieAanbieder", Code: "REQUIRED",
			Message: "ReferentieAanbieder is required",
		})
	}
	if vt.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: "REQUIRED",
			Message: "Product.Categorie is required",
		})
	}
	if vt.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: "REQUIRED",
			Message: "Product.Code is required",
		})
	}
	if vt.Ingangsdatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: "REQUIRED",
			Message: "Ingangsdatum is required",
		})
	} else if !ValidateDate(vt.Ingangsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: "INVALID_DATE",
			Message: "Ingangsdatum must be formatted YYYY-MM-DD",
		})
	}
	if vt.Einddatum != "" {
		if !ValidateDate(vt.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "INVALID_DATE",
				Message: "Einddatum must be formatted YYYY-MM-DD",
			})
		} else if vt.Ingangsdatum != "" && !ValidatePeriod(vt.Ingangsdatum, vt.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: "INVALID_PERIOD",
				Message: "Einddatum must be on or after Ingangsdatum",
			})
		}
	}
	return errs
}

// WMO302Client is a client record within a WMO302 message.
type WMO302Client struct {
	// Bsn is the client's Burgerservicenummer.
	Bsn string `xml:"Bsn"`
	// Naam holds the client's name components.
	Naam Naam `xml:"Naam"`
	// Geboortedatum is the client's date of birth in YYYY-MM-DD format.
	Geboortedatum string `xml:"Geboortedatum,omitempty"`
	// Geslacht is the client's gender code.
	Geslacht string `xml:"Geslacht,omitempty"`
	// VerzoekToewijzingen contains one or more care assignment requests.
	VerzoekToewijzingen []VerzoekToewijzing `xml:"VerzoekToewijzing"`
}

// VerzoekToewijzing is a single care assignment request record within a WMO302 message.
type VerzoekToewijzing struct {
	// ReferentieAanbieder is the care provider's reference for this request.
	ReferentieAanbieder string `xml:"ReferentieAanbieder"`
	// Product identifies the requested care product.
	Product Product `xml:"Product"`
	// Ingangsdatum is the requested start date (YYYY-MM-DD).
	Ingangsdatum string `xml:"Ingangsdatum"`
	// Einddatum is the optional requested end date (YYYY-MM-DD).
	Einddatum string `xml:"Einddatum,omitempty"`
	// Omvang specifies the requested volume, unit, and frequency.
	Omvang *Omvang `xml:"Omvang,omitempty"`
	// Commentaar is a free-text remark.
	Commentaar string `xml:"Commentaar,omitempty"`
}
