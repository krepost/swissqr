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
	"github.com/krepost/structref"
	"github.com/krepost/swissqr"
	"os"
	"time"
)

// Second example in Appendix A of the Swiss QR standard.
func Example_schneiderOrder() {
	data := swissqr.Payload{
		Account: swissqr.NewIBANOrDie("CH4431999123000889012"),
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
			Amount:   1949.75,
			Currency: swissqr.CHF,
		},
		UltimateDebtor: swissqr.Entity{
			Name: "Pia-Maria Rutschmann-Schnyder",
			Address: swissqr.StructuredAddress{
				StreetName:     "Grosse Marktgasse",
				BuildingNumber: "28",
				PostCode:       "9400",
				TownName:       "Rorschach",
			},
			CountryCode: "CH",
		},
		Reference: swissqr.PaymentReference{
			Number: structref.NewReferenceNumberOrDie("21 00000 00003 13947 14300 09017"),
		},
		AdditionalInformation: swissqr.PaymentInformation{
			UnstructuredMessage: "Order dated 18.06.2020",
			StructuredMessage: swissqr.BillInformation{
				InvoiceNumber:     "10201409",
				InvoiceDate:       swissqr.OneDate(2019, time.May, 12),
				CustomerReference: "140.000-53",
				VATNumber:         "106017086",
				VATDates:          swissqr.OneDate(2018, time.May, 8),
				VATRates: swissqr.TaxRates{
					swissqr.TaxRate{RatePercent: 7.7},
				},
				Conditions: swissqr.PaymentConditions{
					swissqr.PaymentCondition{DiscountPercent: 2, NumberOfDays: 10},
					swissqr.PaymentCondition{DiscountPercent: 0, NumberOfDays: 30},
				},
			},
		},
		AlternativeProcedureParameters: swissqr.AlternativeProcedures{
			swissqr.AlternativeProcedure{"Name AV1", "UV;UltraPay005;12345"},
			swissqr.AlternativeProcedure{"Name AV2", "XY;XYService;54321"},
		},
	}
	if err := data.Validate(); err != nil {
		fmt.Println("Unexpected error:", err)
	}

	doc := pdf.New()
	canvas := doc.NewPage(21.0*pdf.Cm, 29.7*pdf.Cm)
	if err := swissqr.DrawInvoiceWithBorder(canvas, data, "en"); err != nil {
		fmt.Println("Unexpected error:", err)
	}
	canvas.Close()

	if pdfFile, err := os.Create("example_order.pdf"); err != nil {
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
