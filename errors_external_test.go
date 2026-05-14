package iwmo_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hyperized/iwmo"
)

func TestSentinelErrors_Is(t *testing.T) {
	tests := []struct {
		name    string
		wrapped error
		target  error
	}{
		{"ErrInvalidMessage wrapped with fmt.Errorf", fmt.Errorf("context: %w", iwmo.ErrInvalidMessage), iwmo.ErrInvalidMessage},
		{"ErrUnknownMessage wrapped with fmt.Errorf", fmt.Errorf("context: %w", iwmo.ErrUnknownMessage), iwmo.ErrUnknownMessage},
		{"ErrTransport wrapped with fmt.Errorf", fmt.Errorf("context: %w", iwmo.ErrTransport), iwmo.ErrTransport},
		{"ErrAuthentication wrapped with fmt.Errorf", fmt.Errorf("context: %w", iwmo.ErrAuthentication), iwmo.ErrAuthentication},
		{"ErrConfiguration wrapped with fmt.Errorf", fmt.Errorf("context: %w", iwmo.ErrConfiguration), iwmo.ErrConfiguration},
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
		iwmo.ErrInvalidMessage,
		iwmo.ErrUnknownMessage,
		iwmo.ErrTransport,
		iwmo.ErrAuthentication,
		iwmo.ErrConfiguration,
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
		{iwmo.ErrInvalidMessage, "iwmo: invalid message"},
		{iwmo.ErrUnknownMessage, "iwmo: unknown message type"},
		{iwmo.ErrTransport, "iwmo: transport error"},
		{iwmo.ErrAuthentication, "iwmo: authentication error"},
		{iwmo.ErrConfiguration, "iwmo: configuration error"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("error string = %q, want %q", got, tt.want)
			}
		})
	}
}
