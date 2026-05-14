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
func (m *WMO303) MessageType() string { return MessageTypeWMO303 }

// Validate checks all header fields and every client and prestatie record.
func (m *WMO303) Validate() error {
	var errs ValidationErrors

	errs = append(errs, validateHeader(m.Header)...)

	if m.Header.BerichtCode != berichtCodeWMO303 {
		errs = append(errs, ValidationError{
			Field: msgFieldHeaderBerichtCode, Code: codeWrongCode,
			Message: "BerichtCode moet 303 zijn voor WMO303",
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

	if cl.Declaratieperiode.Begindatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Begindatum", Code: codeRequired,
			Message: "Declaratieperiode.Begindatum is verplicht",
		})
	} else if !ValidateDate(cl.Declaratieperiode.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Begindatum", Code: codeInvalidDate,
			Message: "Declaratieperiode.Begindatum moet de notatie JJJJ-MM-DD hebben",
		})
	}

	if cl.Declaratieperiode.Einddatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Einddatum", Code: codeRequired,
			Message: "Declaratieperiode.Einddatum is verplicht",
		})
	} else if !ValidateDate(cl.Declaratieperiode.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Declaratieperiode.Einddatum", Code: codeInvalidDate,
			Message: "Declaratieperiode.Einddatum moet de notatie JJJJ-MM-DD hebben",
		})
	}

	if cl.Declaratieperiode.Begindatum != "" && cl.Declaratieperiode.Einddatum != "" {
		if !ValidatePeriod(cl.Declaratieperiode.Begindatum, cl.Declaratieperiode.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Declaratieperiode", Code: codeInvalidPeriod,
				Message: "Declaratieperiode.Einddatum moet op of na Begindatum liggen",
			})
		}
	}

	if len(cl.Prestaties) == 0 {
		errs = append(errs, ValidationError{
			Field: prefix + ".Prestatie", Code: codeRequired,
			Message: "ten minste één Prestatie-element is verplicht",
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
			Field: prefix + ".ToewijzingNummer", Code: codeRequired,
			Message: msgToewijzingNummerRequired,
		})
	}

	if p.Product.Categorie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Categorie", Code: codeRequired,
			Message: msgProductCategorieRequired,
		})
	}

	if p.Product.Code == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Product.Code", Code: codeRequired,
			Message: msgProductCodeRequired,
		})
	}

	if p.Begindatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: codeRequired,
			Message: "Begindatum is verplicht",
		})
	} else if !ValidateDate(p.Begindatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Begindatum", Code: codeInvalidDate,
			Message: msgBegindatumFormat,
		})
	}

	if p.Einddatum == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: codeRequired,
			Message: "Einddatum is verplicht",
		})
	} else if !ValidateDate(p.Einddatum) {
		errs = append(errs, ValidationError{
			Field: prefix + ".Einddatum", Code: codeInvalidDate,
			Message: msgEinddatumFormat,
		})
	}

	if p.Begindatum != "" && p.Einddatum != "" {
		if !ValidatePeriod(p.Begindatum, p.Einddatum) {
			errs = append(errs, ValidationError{
				Field: prefix + ".Einddatum", Code: codeInvalidPeriod,
				Message: msgEinddatumAfterBegindatum,
			})
		}
	}

	if p.Omvang.Volume == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Volume", Code: codeRequired,
			Message: "Omvang.Volume is verplicht",
		})
	}

	if p.Omvang.Eenheid == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Eenheid", Code: codeRequired,
			Message: "Omvang.Eenheid is verplicht",
		})
	}

	if p.Omvang.Frequentie == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Omvang.Frequentie", Code: codeRequired,
			Message: "Omvang.Frequentie is verplicht",
		})
	}

	if p.Bedrag == "" {
		errs = append(errs, ValidationError{
			Field: prefix + ".Bedrag", Code: codeRequired,
			Message: "Bedrag is verplicht",
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
	// Geslacht is the client's gender code (0=unknown, 1=male, 2=female, 9=unspecified).
	Geslacht string `xml:"Geslacht,omitempty"`
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
	// Required for Prestatie (unlike WMO301 Toewijzing where Omvang is optional).
	Omvang Omvang `xml:"Omvang"`
	// Bedrag is the declared amount as a string (e.g. "1600.00").
	Bedrag string `xml:"Bedrag"`
}
