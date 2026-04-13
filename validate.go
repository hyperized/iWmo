package iwmo

import (
	"fmt"
	"strings"
	"time"
)

// ValidationError represents a single field-level validation failure.
type ValidationError struct {
	// Field is the dot-separated path to the offending field (e.g. "Client[0].Bsn").
	Field string
	// Code is a machine-readable error code (e.g. "REQUIRED", "INVALID_BSN").
	Code string
	// Message is a human-readable description of the violation.
	Message string
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("field %s (%s): %s", e.Field, e.Code, e.Message)
}

// ValidationErrors is a slice of ValidationError returned by Validate() when
// one or more violations are found. It implements the error interface so it
// can be inspected with errors.As.
type ValidationErrors []ValidationError

// Error returns all validation error messages joined by "; ".
func (ve ValidationErrors) Error() string {
	msgs := make([]string, len(ve))
	for i, e := range ve {
		msgs[i] = e.Error()
	}

	return strings.Join(msgs, "; ")
}

// parseDate parses a YYYY-MM-DD string. It returns the parsed time and true on
// success, or the zero time and false if the string is not a valid date.
func parseDate(s string) (time.Time, bool) {
	ts, err := time.Parse("2006-01-02", s)
	return ts, err == nil
}

// BSN length and elfproef constants.
const (
	bsnLength    = 9 // A BSN is exactly 9 digits.
	bsnLastIndex = 8 // Index of the last digit (uses weight -1).
)

// ValidateBSN validates a Dutch Burgerservicenummer (BSN) using the elfproef
// (11-proof) algorithm.
//
// A valid BSN is exactly 9 ASCII digits. The weighted sum is computed using
// multipliers [9, 8, 7, 6, 5, 4, 3, 2] for digits 0–7 and -1 for digit 8.
// The sum must be positive and divisible by 11.
func ValidateBSN(bsn string) bool {
	if len(bsn) != bsnLength {
		return false
	}

	sum := 0

	for i, ch := range bsn {
		if ch < '0' || ch > '9' {
			return false
		}

		d := int(ch - '0')

		if i < bsnLastIndex {
			sum += d * (bsnLength - i)
		} else {
			sum -= d
		}
	}

	return sum > 0 && sum%11 == 0
}

// ValidateDate reports whether s is a valid ISO 8601 date in YYYY-MM-DD format.
func ValidateDate(s string) bool {
	_, ok := parseDate(s)
	return ok
}

// ValidatePeriod reports whether begin is on or before end.
// Both strings must be valid YYYY-MM-DD dates; false is returned if either is
// unparseable.
func ValidatePeriod(begin, end string) bool {
	b, ok := parseDate(begin)
	if !ok {
		return false
	}

	e, ok := parseDate(end)
	if !ok {
		return false
	}

	return !b.After(e)
}

// validateHeader validates the common header fields shared by all iWMO messages.
// It is called from each message type's Validate() method.
func validateHeader(h Header) ValidationErrors {
	var errs ValidationErrors

	if h.BerichtCode == "" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtCode", Code: "REQUIRED",
			Message: "BerichtCode is required",
		})
	}

	if h.BerichtVersie == "" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtVersie", Code: "REQUIRED",
			Message: "BerichtVersie is required",
		})
	}

	if h.Afzender == "" {
		errs = append(errs, ValidationError{
			Field: "Header.Afzender", Code: "REQUIRED",
			Message: "Afzender is required",
		})
	}

	if h.Ontvanger == "" {
		errs = append(errs, ValidationError{
			Field: "Header.Ontvanger", Code: "REQUIRED",
			Message: "Ontvanger is required",
		})
	}

	if h.BerichtIdentificatie == "" {
		errs = append(errs, ValidationError{
			Field: "Header.BerichtIdentificatie", Code: "REQUIRED",
			Message: "BerichtIdentificatie is required",
		})
	}

	if h.DagtekeningBericht == "" {
		errs = append(errs, ValidationError{
			Field: "Header.DagtekeningBericht", Code: "REQUIRED",
			Message: "DagtekeningBericht is required",
		})
	} else if !ValidateDate(h.DagtekeningBericht) {
		errs = append(errs, ValidationError{
			Field: "Header.DagtekeningBericht", Code: "INVALID_DATE",
			Message: "DagtekeningBericht must be formatted YYYY-MM-DD",
		})
	}

	return errs
}
