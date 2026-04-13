package iwmo

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors_Is(t *testing.T) {
	tests := []struct {
		name    string
		wrapped error
		target  error
	}{
		{"ErrInvalidMessage wrapped with fmt.Errorf", fmt.Errorf("context: %w", ErrInvalidMessage), ErrInvalidMessage},
		{"ErrUnknownMessage wrapped with fmt.Errorf", fmt.Errorf("context: %w", ErrUnknownMessage), ErrUnknownMessage},
		{"ErrTransport wrapped with fmt.Errorf", fmt.Errorf("context: %w", ErrTransport), ErrTransport},
		{"ErrAuthentication wrapped with fmt.Errorf", fmt.Errorf("context: %w", ErrAuthentication), ErrAuthentication},
		{"ErrConfiguration wrapped with fmt.Errorf", fmt.Errorf("context: %w", ErrConfiguration), ErrConfiguration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.wrapped, tt.target) {
				t.Errorf("errors.Is(%v, %v) = false, want true", tt.wrapped, tt.target)
			}
		})
	}
}

func TestSentinelErrors_Distinct(t *testing.T) {
	sentinels := []error{
		ErrInvalidMessage,
		ErrUnknownMessage,
		ErrTransport,
		ErrAuthentication,
		ErrConfiguration,
	}
	for i, a := range sentinels {
		for j, b := range sentinels {
			if i == j {
				continue
			}
			if errors.Is(a, b) {
				t.Errorf("errors.Is(%v, %v) = true, want false (sentinels should be distinct)", a, b)
			}
		}
	}
}

func TestSentinelErrors_Messages(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrInvalidMessage, "iwmo: invalid message"},
		{ErrUnknownMessage, "iwmo: unknown message type"},
		{ErrTransport, "iwmo: transport error"},
		{ErrAuthentication, "iwmo: authentication error"},
		{ErrConfiguration, "iwmo: configuration error"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("error string = %q, want %q", got, tt.want)
			}
		})
	}
}
