// Package iwmo provides a Go client library for iWMO v3.2, the Dutch
// information standard governing message exchange between municipalities
// (gemeenten) and care providers (zorgaanbieders) under the Wet
// maatschappelijke ondersteuning (WMO).
//
// # Supported message types
//
//   - WMO301: Toewijzing (care assignment, gemeente → zorgaanbieder)
//   - WMO302: Verzoek om Toewijzing (care assignment request, zorgaanbieder → gemeente)
//   - WMO303: Declaratie (billing declaration, zorgaanbieder → gemeente)
//   - WMO304: Retourbericht (acknowledgement, bidirectional)
//   - WMO305: Mutatie (care delivery change, zorgaanbieder → gemeente)
//   - WMO315: Statusmelding (status update, zorgaanbieder → gemeente)
//
// All messages are XML governed by XSD schemas published by iStandaarden.
// See https://www.istandaarden.nl/iwmo for the full specification.
//
// # Basic Usage
//
//	client, err := iwmo.NewClient(
//	    iwmo.WithBaseURL("https://example.gemeente.nl/iwmo"),
//	    iwmo.WithAGBCode("12345678"),
//	    iwmo.WithGemeenteCode("0363"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	msg := &iwmo.WMO302{
//	    Header: iwmo.Header{
//	        BerichtCode:          "302",
//	        BerichtVersie:        "3.2",
//	        Afzender:             "12345678",
//	        Ontvanger:            "0363",
//	        BerichtIdentificatie: "MSG-001",
//	        DagtekeningBericht:   "2026-04-12",
//	    },
//	    Clienten: []iwmo.WMO302Client{...},
//	}
//	retour, err := client.SendVerzoekToewijzing(ctx, msg)
//
// # Validation
//
// Every message type implements [Message], whose Validate method performs
// structural and business-rule checks. Errors are returned as [ValidationErrors],
// a slice of [ValidationError] values that can be inspected with errors.As.
//
// BSN (Burgerservicenummer) validation uses the elfproef (11-proof) algorithm.
// Dates must be formatted as YYYY-MM-DD (ISO 8601).
//
// # Transport
//
// The client uses HTTP(S) by default. Alternative transports (e.g., VECOZO
// message bus, file-based exchange) can be plugged in via [WithSender].
package iwmo
