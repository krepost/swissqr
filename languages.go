// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swissqr

import "fmt"

const (
	paymentPart = iota
	accountPayableTo
	reference
	additionalInformation
	currency
	amount
	receipt
	acceptancePoint
	pleaseSeparate
	payableBy
	payableByNameAddress
	inFavourOf
	dateFormat
)

// headings contains all invoice-related strings that require localization.
// All but the date format are taken from the Swiss QR Invoice standard.
var headings = map[int]map[string]string{
	paymentPart: {
		"de": "Zahlteil",
		"fr": "Section paiement",
		"it": "Sezione pagamento",
		"en": "Payment part",
	},
	accountPayableTo: {
		"de": "Konto / Zahlbar an",
		"fr": "Compte / Payable à",
		"it": "Conto / Pagabile a",
		"en": "Account / Payable to",
	},
	reference: {
		"de": "Referenz",
		"fr": "Référence",
		"it": "Riferimento",
		"en": "Reference",
	},
	additionalInformation: {
		"de": "Zusätzliche Informationen",
		"fr": "Informations supplémentaires",
		"it": "Informazioni supplementari",
		"en": "Additional information",
	},
	currency: {
		"de": "Währung",
		"fr": "Monnaie",
		"it": "Valuta",
		"en": "Currency",
	},
	amount: {
		"de": "Betrag",
		"fr": "Montant",
		"it": "Importo",
		"en": "Amount",
	},
	receipt: {
		"de": "Empfangsschein",
		"fr": "Récépissé",
		"it": "Ricevuta",
		"en": "Receipt",
	},
	acceptancePoint: {
		"de": "Annahmestelle",
		"fr": "Point de dépôt",
		"it": "Punto di accettazione",
		"en": "Acceptance point",
	},
	pleaseSeparate: {
		"de": "Vor der Einzahlung abzutrennen",
		"fr": "A détacher avant le versement",
		"it": "De staccare prima del versamento",
		"en": "Separate before paying in",
	},
	payableBy: {
		"de": "Zahlbar durch",
		"fr": "Payable par",
		"it": "Pagabile da",
		"en": "Payable by",
	},
	payableByNameAddress: {
		"de": "Zahlbar durch (Name/Adresse)",
		"fr": "Payable par (nom/adresse)",
		"it": "Pagabile da (nome/indirizzo)",
		"en": "Payable by (name/address)",
	},
	inFavourOf: {
		"de": "Zugunsten",
		"fr": "En faveur de",
		"it": "A favore di",
		"en": "In favour of",
	},
	dateFormat: {
		"de": "02.01.2006",
		"fr": "02.01.2006",
		"it": "02.01.2006",
		"en": "2006-01-02",
	},
}

// checkLanguage returns nil if language is supported.
func checkLanguage(language string) error {
	for _, heading := range headings {
		if _, ok := heading[language]; !ok {
			return fmt.Errorf("Unsupported langauge: %v", language)
		}
	}
	return nil
}
