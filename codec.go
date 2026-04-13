package iwmo

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// Encode marshals msg to XML and prepends the standard XML declaration.
// The returned bytes are ready for transmission or storage.
func Encode(msg Message) ([]byte, error) {
	b, err := xml.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	return append([]byte(xml.Header), b...), nil
}

// sniffBerichtCode performs a byte-level scan to extract the BerichtCode value
// from raw XML, avoiding a full unmarshal on the happy path.
func sniffBerichtCode(data []byte) string {
	start := bytes.Index(data, []byte("<BerichtCode>"))
	if start < 0 {
		return ""
	}
	start += len("<BerichtCode>")
	end := bytes.Index(data[start:], []byte("</BerichtCode>"))
	if end < 0 {
		return ""
	}
	return string(bytes.TrimSpace(data[start : start+end]))
}

// Decode sniffs the BerichtCode in the XML to determine the concrete message
// type, then fully decodes into that type and returns it as a Message.
//
// For known BerichtCodes the XML is parsed only once. For unknown or absent
// codes a fallback unmarshal is performed to surface a proper XML parse error.
//
// Returns [ErrUnknownMessage] if the BerichtCode does not match any known type.
// Returns [ErrInvalidMessage] if the XML cannot be parsed.
func Decode(data []byte) (Message, error) {
	switch sniffBerichtCode(data) {
	case "301":
		return DecodeAs[WMO301](data)
	case "302":
		return DecodeAs[WMO302](data)
	case "303":
		return DecodeAs[WMO303](data)
	case "304":
		return DecodeAs[WMO304](data)
	case "305":
		return DecodeAs[WMO305](data)
	case "315":
		return DecodeAs[WMO315](data)
	default:
		// Code not found or unrecognised; fall back to a full unmarshal to
		// surface a proper XML parse error or report the unknown code.
		var peek berichtPeek
		if err := xml.Unmarshal(data, &peek); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
		}
		return nil, fmt.Errorf("%w: BerichtCode %q", ErrUnknownMessage, peek.Header.BerichtCode)
	}
}

// DecodeAs unmarshals data into a newly allocated *T without inspecting the
// BerichtCode. Use this when the caller already knows the expected message type.
//
// T should be one of [WMO301], [WMO302], [WMO303], [WMO304], [WMO305], or [WMO315].
func DecodeAs[T any](data []byte) (*T, error) {
	var target T
	if err := xml.Unmarshal(data, &target); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
	}
	return &target, nil
}
