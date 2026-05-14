package iwmo

// Shared test fixtures used by both internal (package iwmo) and external
// (package iwmo_test) tests. The external package consumes these via the
// re-exports in export_test.go. Kept in the internal package so internal
// tests like client_test.go and codec_test.go can use them without an
// import cycle.

func validHeader(code string) Header {
	return Header{
		BerichtCode:          code,
		BerichtVersie:        "3.2",
		Afzender:             "0363",
		Ontvanger:            "12345678",
		BerichtIdentificatie: "MSG-TEST-001",
		DagtekeningBericht:   "2026-04-12",
	}
}

const validBSN = "123456782" // elfproef: sum=154, 154%11=0

func validWMO301() *WMO301 {
	return &WMO301{
		Header: validHeader("301"),
		Clienten: []WMO301Client{
			{
				Bsn:           validBSN,
				Naam:          Naam{Voornamen: "Jan", Tussenvoegsels: "van", Achternaam: "Janssen"},
				Geboortedatum: "1980-01-15",
				Toewijzingen: []Toewijzing{
					{
						ToewijzingNummer: "12345",
						Product:          Product{Categorie: "03", Code: "H532"},
						Ingangsdatum:     "2026-05-01",
						Einddatum:        "2026-12-31",
					},
				},
			},
		},
	}
}

func validWMO302() *WMO302 {
	return &WMO302{
		Header: validHeader("302"),
		Clienten: []WMO302Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				VerzoekToewijzingen: []VerzoekToewijzing{
					{
						ReferentieAanbieder: "REF-001",
						Product:             Product{Categorie: "03", Code: "H532"},
						Ingangsdatum:        "2026-05-01",
					},
				},
			},
		},
	}
}

func validWMO303() *WMO303 {
	return &WMO303{
		Header: validHeader("303"),
		Clienten: []WMO303Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Declaratieperiode: Declaratieperiode{
					Begindatum: "2026-04-01",
					Einddatum:  "2026-04-30",
				},
				Prestaties: []Prestatie{
					{
						ToewijzingNummer: "12345",
						Product:          Product{Categorie: "03", Code: "H532"},
						Begindatum:       "2026-04-01",
						Einddatum:        "2026-04-30",
						Omvang:           Omvang{Volume: "32", Eenheid: "uur", Frequentie: "maand"},
						Bedrag:           "1600.00",
					},
				},
			},
		},
	}
}

func validWMO304() *WMO304 {
	return &WMO304{
		Header: WMO304Header{
			Header:                          validHeader("304"),
			GerefereerdBerichtCode:          "302",
			GerefereerdBerichtIdentificatie: "MSG-TEST-001",
		},
		RetourCodes: []RetourCode{
			{Code: "0000", Omschrijving: "Bericht in goede orde ontvangen"},
		},
	}
}

func validWMO305() *WMO305 {
	return &WMO305{
		Header: validHeader("305"),
		Clienten: []WMO305Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Mutaties: []Mutatie{
					{
						ToewijzingNummer: "12345",
						Mutatiedatum:     "2026-04-12",
						Mutatiecode:      MutatiecodeStart,
						Begindatum:       "2026-05-01",
					},
				},
			},
		},
	}
}

func validWMO315() *WMO315 {
	return &WMO315{
		Header: validHeader("315"),
		Clienten: []WMO315Client{
			{
				Bsn:  validBSN,
				Naam: Naam{Achternaam: "Janssen"},
				Statusmeldingen: []StatusmeldingRecord{
					{
						ToewijzingNummer: "12345",
						StatusCode:       "01",
						StatusDatum:      "2026-05-01",
						Commentaar:       "Zorg gestart",
					},
				},
			},
		},
	}
}
