package iwmo

import "encoding/xml"

// Message is the interface implemented by every iWMO message type.
type Message interface {
	// MessageType returns the message code string, e.g. "WMO301".
	MessageType() string
	// Validate performs structural and business-rule checks.
	// It returns nil when the message is valid, or a [ValidationErrors]
	// value (which satisfies error) listing all violations found.
	Validate() error
}

// Header holds the fields common to every iWMO message.
type Header struct {
	// BerichtCode identifies the message type ("301", "302", …).
	BerichtCode string `xml:"BerichtCode"`
	// BerichtVersie is the iWMO version string, e.g. "3.2".
	BerichtVersie string `xml:"BerichtVersie"`
	// Afzender is the sender's identification code.
	Afzender string `xml:"Afzender"`
	// Ontvanger is the receiver's identification code.
	Ontvanger string `xml:"Ontvanger"`
	// BerichtIdentificatie is the unique message identifier.
	BerichtIdentificatie string `xml:"BerichtIdentificatie"`
	// DagtekeningBericht is the message date in YYYY-MM-DD format.
	DagtekeningBericht string `xml:"DagtekeningBericht"`
	// XsltVersie is the optional XSLT version string.
	XsltVersie string `xml:"XsltVersie,omitempty"`
	// XsdVersie is the optional XSD version string.
	XsdVersie string `xml:"XsdVersie,omitempty"`
}

// Naam holds the components of a person's name.
type Naam struct {
	// Voornamen holds the person's first name(s).
	Voornamen string `xml:"Voornamen,omitempty"`
	// Tussenvoegsels holds name prefixes such as "van", "de", "van der".
	Tussenvoegsels string `xml:"Tussenvoegsels,omitempty"`
	// Achternaam is the family name (required).
	Achternaam string `xml:"Achternaam"`
}

// Product identifies a care product by its category and product code.
type Product struct {
	// Categorie is the product category code (e.g. "03").
	Categorie string `xml:"Categorie"`
	// Code is the product code within the category.
	Code string `xml:"Code"`
}

// Omvang specifies the volume, unit, and frequency of care delivery.
type Omvang struct {
	// Volume is the quantity, e.g. "8".
	Volume string `xml:"Volume"`
	// Eenheid is the unit, e.g. "uur" (hour) or "dagdeel" (half-day).
	Eenheid string `xml:"Eenheid"`
	// Frequentie is the frequency period, e.g. "week" or "maand" (month).
	Frequentie string `xml:"Frequentie"`
}

// berichtPeek is used internally by Decode to sniff BerichtCode before
// fully decoding into a concrete message type.
type berichtPeek struct {
	XMLName xml.Name `xml:"Bericht"`
	Header  struct {
		BerichtCode string `xml:"BerichtCode"`
	} `xml:"Header"`
}
