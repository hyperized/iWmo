# iwmo

Go client library for **iWMO v3.2**, the Dutch information standard for XML message exchange between municipalities (*gemeenten*) and care providers (*zorgaanbieders*) under the [Wet maatschappelijke ondersteuning](https://www.rijksoverheid.nl/onderwerpen/zorg-en-ondersteuning-thuis/wmo-2015) (WMO).

Zero runtime dependencies -- stdlib only.

```
go get github.com/hyperized/iwmo
```

## Message types

| Code   | Name                  | Direction                   | Purpose                      |
|--------|-----------------------|-----------------------------|------------------------------|
| WMO301 | Toewijzing            | Gemeente -> Zorgaanbieder   | Care assignment              |
| WMO302 | Verzoek om Toewijzing | Zorgaanbieder -> Gemeente   | Care assignment request      |
| WMO303 | Declaratie            | Zorgaanbieder -> Gemeente   | Billing declaration          |
| WMO304 | Retourbericht         | Bidirectional               | Acknowledgement              |
| WMO305 | Mutatie               | Zorgaanbieder -> Gemeente   | Care delivery change         |
| WMO315 | Statusmelding         | Zorgaanbieder -> Gemeente   | Status update                |

All messages conform to XSD schemas published by [iStandaarden](https://www.istandaarden.nl/iwmo).

## Usage

### Sending a message

```go
package main

import (
	"context"
	"log"

	"github.com/hyperized/iwmo"
)

func main() {
	client, err := iwmo.NewClient(
		iwmo.WithBaseURL("https://example.gemeente.nl/iwmo"),
		iwmo.WithAGBCode("12345678"),
		iwmo.WithGemeenteCode("0363"),
	)
	if err != nil {
		log.Fatal(err)
	}

	msg := &iwmo.WMO302{
		Header: iwmo.Header{
			BerichtCode:          "302",
			BerichtVersie:        "3.2",
			Afzender:             "12345678",
			Ontvanger:            "0363",
			BerichtIdentificatie: "MSG-001",
			DagtekeningBericht:   "2026-04-12",
		},
		Clienten: []iwmo.WMO302Client{
			{
				Bsn:            "123456782",
				Naam:           iwmo.Naam{Voornamen: "Jan", Achternaam: "Janssen"},
				Geboortedatum:  "1980-01-15",
				Geslacht:       "1",
				VerzoekToewijzingen: []iwmo.VerzoekToewijzing{
					{
						ReferentieAanbieder: "REF-001",
						Product:             iwmo.Product{Categorie: "03", Code: "H532"},
						Ingangsdatum:        "2026-05-01",
						Einddatum:           "2026-12-31",
						Omvang:              iwmo.Omvang{Volume: "8", Eenheid: "uur", Frequentie: "week"},
					},
				},
			},
		},
	}

	retour, err := client.SendVerzoekToewijzing(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received WMO304 with %d retour codes", len(retour.RetourCodes))
}
```

### Decoding incoming XML

```go
// Auto-detect message type from BerichtCode
msg, err := iwmo.Decode(xmlBytes)

// Decode into a known type
toewijzing, err := iwmo.DecodeAs[iwmo.WMO301](xmlBytes)
```

### Encoding to XML

```go
xmlBytes, err := iwmo.Encode(msg)
```

### Validation

Every message type implements `Validate()`, which checks structural and business-rule constraints. Errors are returned as `ValidationErrors`:

```go
if err := msg.Validate(); err != nil {
	var ve iwmo.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			fmt.Printf("%s (%s): %s\n", e.Field, e.Code, e.Message)
		}
	}
}
```

Standalone validation helpers are also available:

```go
iwmo.ValidateBSN("123456782")          // true  (elfproef/11-proof)
iwmo.ValidateDate("2026-05-01")        // true  (YYYY-MM-DD)
iwmo.ValidatePeriod("2026-05-01", "2026-12-31") // true  (begin <= end)
```

### Custom transport

The default transport uses HTTP POST. Provide a custom `Sender` for alternative backends (VECOZO, file-based exchange, etc.):

```go
client, err := iwmo.NewClient(
	iwmo.WithSender(myVecozoSender),
	iwmo.WithAGBCode("12345678"),
)
```

The `Sender` interface:

```go
type Sender interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}
```

## Error handling

The library defines four sentinel errors for use with `errors.Is`:

| Error                | When                                                   |
|----------------------|--------------------------------------------------------|
| `ErrInvalidMessage`  | Validation failure or XML encoding/decoding error      |
| `ErrUnknownMessage`  | `Decode` encounters an unrecognized `BerichtCode`      |
| `ErrTransport`       | Network error or non-2xx HTTP response                 |
| `ErrAuthentication`  | HTTP 401 or 403 from the endpoint                      |

## Testing

```sh
# Unit tests
go test ./...

# With race detector
go test -race ./...

# Integration tests (requires env vars)
IWMO_BASE_URL=https://... IWMO_AGB_CODE=12345678 IWMO_GEMEENTE=0363 \
  go test -tags integration -v ./...
```

## License

See [LICENSE](LICENSE) for details.
