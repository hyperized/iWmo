package iwmo

// Exported message-type identifiers returned by [Message.MessageType].
// These are part of the public API surface and stable across releases.
const (
	// MessageTypeWMO301 is the value returned by (*WMO301).MessageType().
	MessageTypeWMO301 = "WMO301"
	// MessageTypeWMO302 is the value returned by (*WMO302).MessageType().
	MessageTypeWMO302 = "WMO302"
	// MessageTypeWMO303 is the value returned by (*WMO303).MessageType().
	MessageTypeWMO303 = "WMO303"
	// MessageTypeWMO304 is the value returned by (*WMO304).MessageType().
	MessageTypeWMO304 = "WMO304"
	// MessageTypeWMO305 is the value returned by (*WMO305).MessageType().
	MessageTypeWMO305 = "WMO305"
	// MessageTypeWMO315 is the value returned by (*WMO315).MessageType().
	MessageTypeWMO315 = "WMO315"
)

// WMO305 Mutatiecode values. These are the only valid values for
// [Mutatie.Mutatiecode]; anything else fails validation with codeInvalidValue.
const (
	// MutatiecodeStart marks the beginning of care delivery.
	MutatiecodeStart = "01"
	// MutatiecodeWijziging marks a change to an existing care delivery.
	MutatiecodeWijziging = "02"
	// MutatiecodeStop marks the end of care delivery.
	MutatiecodeStop = "03"
)

// BerichtCode values used in the XML <Header><BerichtCode> element. Distinct
// from the MessageType* constants above: BerichtCode is the wire format
// identifier ("301"), MessageType is the Go-facing label ("WMO301").
const (
	berichtCodeWMO301 = "301"
	berichtCodeWMO302 = "302"
	berichtCodeWMO303 = "303"
	berichtCodeWMO304 = "304"
	berichtCodeWMO305 = "305"
	berichtCodeWMO315 = "315"
)

// ValidationError.Code values. Unexported because they shape the API surface
// but should be treated as opaque by callers — match against them via the
// package-level constants if interop is needed.
const (
	codeRequired      = "REQUIRED"
	codeInvalidDate   = "INVALID_DATE"
	codeInvalidPeriod = "INVALID_PERIOD"
	codeInvalidBSN    = "INVALID_BSN"
	codeInvalidValue  = "INVALID_VALUE"
	codeWrongCode     = "WRONG_CODE"
)

// Repeated user-facing validation message fragments. Kept private; the
// Message field on ValidationError is documentation, not API.
const (
	msgFieldClient                = "Client"
	msgFieldHeaderBerichtCode     = "Header.BerichtCode"
	msgClientRequired             = "ten minste één Client-element is verplicht"
	msgBSNInvalid                 = "BSN moet 9 cijfers zijn en voldoen aan de elfproef"
	msgAchternaamRequired         = "Achternaam is verplicht"
	msgGeboortedatumFormat        = "Geboortedatum moet de notatie JJJJ-MM-DD hebben"
	msgToewijzingNummerRequired   = "ToewijzingNummer is verplicht"
	msgProductCategorieRequired   = "Product.Categorie is verplicht"
	msgProductCodeRequired        = "Product.Code is verplicht"
	msgIngangsdatumRequired       = "Ingangsdatum is verplicht"
	msgIngangsdatumFormat         = "Ingangsdatum moet de notatie JJJJ-MM-DD hebben"
	msgEinddatumFormat            = "Einddatum moet de notatie JJJJ-MM-DD hebben"
	msgEinddatumAfterIngangsdatum = "Einddatum moet op of na Ingangsdatum liggen"
	msgEinddatumAfterBegindatum   = "Einddatum moet op of na Begindatum liggen"
	msgBegindatumFormat           = "Begindatum moet de notatie JJJJ-MM-DD hebben"
	msgGeslachtInvalid            = "Geslacht moet leeg zijn of één van 0, 1, 2 of 9"
)
