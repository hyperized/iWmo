package iwmo

import (
	"encoding/xml"
	"fmt"
)

// Encode marshals msg to XML and prepends the standard XML declaration.
// The returned bytes are ready for transmission or storage.
func Encode(msg Message) ([]byte, error) {
	b, err := xml.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidMessage, err)
	}

	buf := make([]byte, 0, len(xml.Header)+len(b))
	buf = append(buf, xml.Header...)
	buf = append(buf, b...)

	return buf, nil
}

// Decode inspects the BerichtCode in the XML header to determine the concrete
// message type, then fully decodes into that type and returns it as a Message.
//
// For known BerichtCodes the XML is parsed twice: once to sniff the code, once
// to decode into the concrete struct. For unknown or absent codes, the sniff
// parse is the only one performed.
//
// Returns [ErrUnknownMessage] if the BerichtCode does not match any known type.
// Returns [ErrInvalidMessage] if the XML cannot be parsed.
func Decode(data []byte) (Message, error) {
	var peek berichtPeek
	if err := xml.Unmarshal(data, &peek); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidMessage, err)
	}

	switch peek.Header.BerichtCode {
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
		return nil, fmt.Errorf("%w: %w", ErrInvalidMessage, err)
	}

	return &target, nil
}
