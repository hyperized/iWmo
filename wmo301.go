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
func (m *WMO301) MessageType() string { return MessageTypeWMO301 }

// Validate checks all header fields and every client and assignment record.
func (m *WMO301) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != berichtCodeWMO301 {
		errs = append(errs, ValidationError{
			Field: msgFieldHeaderBerichtCode, Code: codeWrongCode,
			Message: "BerichtCode moet 301 zijn voor WMO301",
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

	if len(cl.Toewijzingen) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Toewijzing", Code: codeRequired,
			Message: "ten minste één Toewijzing-element is verplicht",
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
			Field: prefix + ".ToewijzingNummer", Code: codeRequired,
			Message: msgToewijzingNummerRequired,
		})
	}

	if tw.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: codeRequired,
			Message: msgProductCategorieRequired,
		})
	}

	if tw.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: codeRequired,
			Message: msgProductCodeRequired,
		})
	}

	if tw.Toewijzingsdatum != "" && !ValidateDate(tw.Toewijzingsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Toewijzingsdatum", Code: codeInvalidDate,
			Message: "Toewijzingsdatum moet de notatie JJJJ-MM-DD hebben",
		})
	}

	if tw.Ingangsdatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: codeRequired,
			Message: msgIngangsdatumRequired,
		})
	} else if !ValidateDate(tw.Ingangsdatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Ingangsdatum", Code: codeInvalidDate,
			Message: msgIngangsdatumFormat,
		})
	}

	if tw.Einddatum != "" {
		if !ValidateDate(tw.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeInvalidDate,
				Message: msgEinddatumFormat,
			})
		} else if tw.Ingangsdatum != "" && !ValidatePeriod(tw.Ingangsdatum, tw.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeInvalidPeriod,
				Message: msgEinddatumAfterIngangsdatum,
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
