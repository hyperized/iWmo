//go:build integration

package iwmo

import (
	"context"
	"os"
	"testing"
	"time"
)

// Integration tests require a running iWMO-compatible endpoint.
//
// Configure via environment variables before running:
//
//	IWMO_BASE_URL   - base URL of the test endpoint (required)
//	IWMO_AGB_CODE   - AGB-Z code of the test care provider (required)
//	IWMO_GEMEENTE   - municipality code of the test gemeente (required)
//
// Run with:
//
//	go test -tags integration -v ./...
func newIntegrationClient(t *testing.T) *Client {
	t.Helper()
	baseURL := os.Getenv("IWMO_BASE_URL")
	if baseURL == "" {
		t.Skip("IWMO_BASE_URL not set; skipping integration test")
	}
	agb := os.Getenv("IWMO_AGB_CODE")
	gem := os.Getenv("IWMO_GEMEENTE")
	c, err := NewClient(
		WithBaseURL(baseURL),
		WithAGBCode(agb),
		WithGemeenteCode(gem),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func integrationCtx(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func TestIntegration_SendVerzoekToewijzing(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := integrationCtx(t)
	defer cancel()

	msg := validWMO302()
	retour, err := c.SendVerzoekToewijzing(ctx, msg)
	if err != nil {
		t.Fatalf("SendVerzoekToewijzing: %v", err)
	}
	if len(retour.RetourCodes) == 0 {
		t.Error("expected at least one RetourCode in response")
	}
	t.Logf("WMO304 RetourCodes: %+v", retour.RetourCodes)
}

func TestIntegration_SendDeclaratie(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := integrationCtx(t)
	defer cancel()

	msg := validWMO303()
	retour, err := c.SendDeclaratie(ctx, msg)
	if err != nil {
		t.Fatalf("SendDeclaratie: %v", err)
	}
	if len(retour.RetourCodes) == 0 {
		t.Error("expected at least one RetourCode in response")
	}
}

func TestIntegration_SendMutatie(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := integrationCtx(t)
	defer cancel()

	msg := validWMO305()
	retour, err := c.SendMutatie(ctx, msg)
	if err != nil {
		t.Fatalf("SendMutatie: %v", err)
	}
	if len(retour.RetourCodes) == 0 {
		t.Error("expected at least one RetourCode in response")
	}
}

func TestIntegration_SendStatusmelding(t *testing.T) {
	c := newIntegrationClient(t)
	ctx, cancel := integrationCtx(t)
	defer cancel()

	msg := validWMO315()
	retour, err := c.SendStatusmelding(ctx, msg)
	if err != nil {
		t.Fatalf("SendStatusmelding: %v", err)
	}
	if len(retour.RetourCodes) == 0 {
		t.Error("expected at least one RetourCode in response")
	}
}
