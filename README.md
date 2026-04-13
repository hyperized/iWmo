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

The library defines five sentinel errors for use with `errors.Is`:

| Error                | When                                                   |
|----------------------|--------------------------------------------------------|
| `ErrInvalidMessage`  | Validation failure or XML encoding/decoding error      |
| `ErrUnknownMessage`  | `Decode` encounters an unrecognized `BerichtCode`      |
| `ErrTransport`       | Network error or non-2xx HTTP response                 |
| `ErrAuthentication`  | HTTP 401 or 403 from the endpoint                      |
| `ErrConfiguration`   | `NewClient` called without required options            |

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

## Benchmarks

Measured on Apple M2 Pro, Go 1.26, 5 runs (`go test -bench=. -benchmem -count=5`).

### Encode (struct to XML)

| Message | ns/op | B/op | allocs/op |
|---------|------:|-----:|----------:|
| WMO301  | 3,686 | 6,496 | 14 |
| WMO302  | 3,228 | 6,240 | 14 |
| WMO303  | 4,268 | 6,752 | 14 |
| WMO304  | 2,483 | 5,792 | 12 |
| WMO305  | 3,103 | 5,792 | 12 |
| WMO315  | 2,920 | 5,856 | 12 |

### Decode (XML to struct, auto-detect via BerichtCode)

| Message | ns/op  | B/op   | allocs/op |
|---------|-------:|-------:|----------:|
| WMO301  | 27,540 | 15,264 | 397 |
| WMO302  | 23,064 | 12,976 | 329 |
| WMO303  | 32,290 | 17,888 | 474 |
| WMO304  | 19,781 | 10,224 | 253 |
| WMO305  | 21,946 | 12,336 | 314 |
| WMO315  | 22,144 | 12,304 | 314 |

### DecodeAs (XML to struct, known type, single-pass)

| Message | ns/op  | B/op  | allocs/op |
|---------|-------:|------:|----------:|
| WMO301  | 15,001 | 8,424 | 221 |
| WMO302  | 12,529 | 7,096 | 181 |
| WMO303  | 17,679 | 9,920 | 265 |
| WMO304  | 10,449 | 5,400 | 137 |
| WMO305  | 11,908 | 6,744 | 173 |
| WMO315  | 12,082 | 6,696 | 173 |

`DecodeAs` is ~1.8x faster than `Decode` because it skips the header-sniffing unmarshal pass.

Run benchmarks yourself:

```sh
go test -bench=. -benchmem ./...
```

## License

See [LICENSE](LICENSE) for details.
