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
func (m *WMO302) MessageType() string { return MessageTypeWMO302 }

// Validate checks all header fields and every client and request record.
func (m *WMO302) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != berichtCodeWMO302 {
		errs = append(errs, ValidationError{
			Field: msgFieldHeaderBerichtCode, Code: codeWrongCode,
			Message: "BerichtCode moet 302 zijn voor WMO302",
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

	errs = append(errs, validateGeslacht(prefix+".Geslacht", cl.Geslacht)...)

	if len(cl.VerzoekToewijzingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".VerzoekToewijzing", Code: codeRequired,
			Message: "ten minste één VerzoekToewijzing-element is verplicht",
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
			Field: prefix + ".ReferentieAanbieder", Code: codeRequired,
			Message: "ReferentieAanbieder is verplicht",
		})
	}

	if vt.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: codeRequired,
			Message: msgProductCategorieRequired,
		})
	}

	if vt.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: codeRequired,
			Message: msgProductCodeRequired,
		})
	}

	if vt.Ingangsdatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: codeRequired,
			Message: msgIngangsdatumRequired,
		})
	} else if !ValidateDate(vt.Ingangsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: codeInvalidDate,
			Message: msgIngangsdatumFormat,
		})
	}

	if vt.Einddatum != "" {
		if !ValidateDate(vt.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeInvalidDate,
				Message: msgEinddatumFormat,
			})
		} else if vt.Ingangsdatum != "" && !ValidatePeriod(vt.Ingangsdatum, vt.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeInvalidPeriod,
				Message: msgEinddatumAfterIngangsdatum,
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
