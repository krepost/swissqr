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

// Third example in Appendix A of the Swiss QR standard.
func Example_armeeDuSalut() {
	data := swissqr.Payload{
		Account: swissqr.NewIBANOrDie("CH3709000000304442225"),
		Creditor: swissqr.Entity{
			Name: "Fondation Armée du salut suisse",
			Address: swissqr.StructuredAddress{
				PostCode: "3000",
				TownName: "Berne",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: swissqr.PaymentAmount{Currency: swissqr.CHF},
		AdditionalInformation: swissqr.PaymentInformation{
			UnstructuredMessage: "Don pour l'action Fête Hiver",
		},
	}
	if err := data.Validate(); err != nil {
		fmt.Println("Unexpected error:", err)
	}

	doc := pdf.New()
	canvas := doc.NewPage(21.0*pdf.Cm, 10.5*pdf.Cm) // QR invoice size.
	if err := swissqr.DrawInvoiceWithScissors(canvas, data, "fr"); err != nil {
		fmt.Println("Unexpected error:", err)
	}
	canvas.Close()

	if pdfFile, err := os.Create("example_salut.pdf"); err != nil {
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
