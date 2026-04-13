package iwmo

import "errors"

// Sentinel errors. Use errors.Is / errors.As to inspect wrapped errors.
var (
	// ErrInvalidMessage is returned when a message fails structural or
	// business-rule validation, or cannot be encoded/decoded.
	ErrInvalidMessage = errors.New("iwmo: invalid message")

	// ErrUnknownMessage is returned by Decode when the BerichtCode in the
	// XML header does not correspond to any known iWMO message type.
	ErrUnknownMessage = errors.New("iwmo: unknown message type")

	// ErrTransport is returned when the underlying HTTP (or custom) sender
	// receives a non-success response or encounters a network error.
	ErrTransport = errors.New("iwmo: transport error")

	// ErrAuthentication is returned when the endpoint responds with HTTP 401
	// or 403, indicating an authentication or authorisation failure.
	ErrAuthentication = errors.New("iwmo: authentication error")

	// ErrConfiguration is returned when the Client is misconfigured, e.g.
	// when neither WithBaseURL nor WithSender was supplied to NewClient.
	ErrConfiguration = errors.New("iwmo: configuration error")
)
