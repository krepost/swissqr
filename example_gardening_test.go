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

package swissqr_test

import (
	"bitbucket.org/krepost/gopdf/pdf"
	"fmt"
	"github.com/krepost/swissqr"
	"os"
)

// First example in Appendix A of the Swiss QR standard.
func Example_gardeningInvoice() {
	data := swissqr.Payload{
		Account: swissqr.NewIBANOrDie("CH58 0079 1123 0008 8901 2"),
		Creditor: swissqr.Entity{
			Name: "Robert Schneider AG",
			Address: swissqr.StructuredAddress{
				StreetName:     "Rue du Lac",
				BuildingNumber: "1268",
				PostCode:       "2501",
				TownName:       "Biel",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: swissqr.PaymentAmount{
			Amount:   3949.75,
			Currency: swissqr.CHF,
		},
		UltimateDebtor: swissqr.Entity{
			Name: "Pia Rutschmann",
			Address: swissqr.StructuredAddress{
				StreetName:     "Marktgasse",
				BuildingNumber: "28",
				PostCode:       "9400",
				TownName:       "Rorschach",
			},
			CountryCode: "CH",
		},
		AdditionalInformation: swissqr.PaymentInformation{
			UnstructuredMessage: "Bill no. 3139 for gardening work and disposal of waste material",
		},
	}
	if err := data.Validate(); err != nil {
		fmt.Println("Unexpected error:", err)
	}

	doc := pdf.New()
	canvas := doc.NewPage(21.0*pdf.Cm, 10.5*pdf.Cm) // QR invoice size.
	if err := swissqr.DrawInvoiceWithScissors(canvas, data, "en"); err != nil {
		fmt.Println("Unexpected error:", err)
	}
	canvas.Close()

	if pdfFile, err := os.Create("example_gardening.pdf"); err != nil {
		fmt.Println("Unexpected error:", err)
	} else {
		defer pdfFile.Close()
		if err := doc.Encode(pdfFile); err != nil {
			fmt.Println("Unexpected error:", err)
		}
	}

	fmt.Println("PDF file successfully created.")

	// Output:
	// PDF file successfully created.
}
