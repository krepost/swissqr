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
	"math"
	"os"
	"strings"
)

// «If the QR-bill with payment part and receipt or the separate payment part
// with receipt are generated as a PDF document and sent electronically, the A6
// format of the payment part and the receipt on the left must be indicated by
// lines. Each of these lines must bear the scissors symbol “✂” or
// alternatively the instruction “Separate before paying in” above the line
// (outside the payment part). This indicates to the debtor that he or she must
// neatly separate the payment part and receipt if they want to forward the
// QR-bill to their financial institution by post for payment, or settle it at
// the post office counter (branches or branches of partner organisations).»
func Example_boxWithScissors() {
	doc := pdf.New()
	if err := box(doc); err != nil {
		fmt.Println(err)
	}
	if pdfFile, err := os.Create("example_box.pdf"); err != nil {
		fmt.Println(err)
	} else {
		defer pdfFile.Close()
		if err := doc.Encode(pdfFile); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("PDF file successfully created.")

	// Output:
	// PDF file successfully created.
}

func box(doc *pdf.Document) error {
	data := swissqr.Payload{
		Account: swissqr.NewIBANOrDie("CH5800791123000889012"),
		Creditor: swissqr.Entity{
			Name: "Mobile Finanz AG",
			Address: swissqr.StructuredAddress{
				StreetName:     "Bernerstrasse",
				BuildingNumber: "234A",
				PostCode:       "8640",
				TownName:       "Rapperswil",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: swissqr.PaymentAmount{
			Amount:   5005.95,
			Currency: swissqr.CHF,
		},
		UltimateDebtor: swissqr.Entity{
			Name: "Pia-Maria Rutschmann-Schnyder",
			Address: swissqr.StructuredAddress{
				StreetName:     "Grosse Marktgasse",
				BuildingNumber: "28/5",
				PostCode:       "9400",
				TownName:       "Rorschach",
			},
			CountryCode: "CH",
		},
		Reference: swissqr.PaymentReference{
			Number: structref.NewCreditorReferenceOrDie("RF83 1234 5678 9123 4567 8912"),
		},
		AdditionalInformation: swissqr.PaymentInformation{
			UnstructuredMessage: "Beachten Sie unsere Sonderangebotswoche bis 23.02.2017!",
		},
	}
	if err := data.Validate(); err != nil {
		return err
	}

	canvas := doc.NewPage(21.0*pdf.Cm, 29.7*pdf.Cm)
	if err := drawRotatedText(canvas); err != nil {
		return err
	}
	if err := drawAddress(canvas, data); err != nil {
		return err
	}
	if err := drawMainText(canvas, data); err != nil {
		return err
	}
	if err := swissqr.DrawInvoiceWithBorder(canvas, data, "de"); err != nil {
		return err
	}
	canvas.Close()
	return nil
}

func drawRotatedText(canvas *pdf.Canvas) error {
	font, err := canvas.Document().AddFont(pdf.HelveticaBold, pdf.WinAnsiEncoding)
	if err != nil {
		return err
	}
	canvas.Push()
	defer canvas.Pop()
	text := new(pdf.Text)
	text.UseFont(font, 120, 120)
	text.Text("Mobile")
	canvas.Rotate(math.Pi / 2.0)
	canvas.Translate(19.9*pdf.Cm-text.X()/2.0, -110)
	canvas.SetColor(0, 0, 0) // White text colour.
	canvas.DrawText(text)
	return nil
}

func drawAddress(canvas *pdf.Canvas, data swissqr.Payload) error {
	creditor, err := data.Creditor.ToLines()
	if err != nil {
		return err
	}
	debtor, err := data.UltimateDebtor.ToLines()
	if err != nil {
		return err
	}
	font, err := canvas.Document().AddFont(pdf.Helvetica, pdf.WinAnsiEncoding)
	if err != nil {
		return err
	}
	canvas.Push()
	defer canvas.Pop()
	address := new(pdf.Text)
	address.UseFont(font, 6, 8)
	address.Text(strings.Join(creditor, "  •  "))
	width := address.X()
	address.NextLine()
	address.UseFont(font, 10, 12)
	for _, line := range debtor {
		address.NextLine()
		address.Text(line)
	}
	canvas.SetColor(0, 0, 0) // Black text colour.
	canvas.Translate(11.0*pdf.Cm, 22.5*pdf.Cm-address.Y()/2.0)
	canvas.DrawText(address)
	adressSeparator := new(pdf.Path)
	adressSeparator.Move(pdf.Point{0, -3})
	adressSeparator.Line(pdf.Point{width, -3})
	canvas.SetLineWidth(0.5)
	canvas.Stroke(adressSeparator)
	return nil
}

func drawMainText(canvas *pdf.Canvas, data swissqr.Payload) error {
	amount := fmt.Sprintf("%v %v", data.CurrencyAmount.Currency, data.CurrencyAmount.Amount)
	roman, err := canvas.Document().AddFont(pdf.Times, pdf.WinAnsiEncoding)
	if err != nil {
		return err
	}
	bold, err := canvas.Document().AddFont(pdf.TimesBold, pdf.WinAnsiEncoding)
	if err != nil {
		return err
	}
	italics, err := canvas.Document().AddFont(pdf.TimesItalic, pdf.WinAnsiEncoding)
	if err != nil {
		return err
	}
	canvas.Push()
	defer canvas.Pop()
	text := new(pdf.Text)
	text.UseFont(bold, 10, 12)
	text.Text("Rechnung")
	text.NextLine()
	text.NextLine()
	text.UseFont(roman, 10, 12)
	text.Text("Für unsere Dienstleistungen per 1. April 2017: " + amount + ".")
	text.NextLine()
	text.NextLine()
	text.Text("Mit freundlichen Grüssen,")
	text.NextLine()
	text.NextLine()
	text.UseFont(italics, 10, 12)
	text.Text("Ihr Mobile-Team")
	canvas.Translate(5*pdf.Cm, 18*pdf.Cm)
	canvas.DrawText(text)
	return nil
}
