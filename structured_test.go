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

import (
	"strings"
	"testing"
	"time"
)

func TestEmptyStructuredMessage(t *testing.T) {
	msg := BillInformation{}
	if err := msg.Validate(); err != nil {
		t.Errorf("Expected no error; got %v", err)
		return
	}
	expected := ""
	actual := msg.ToString()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

// Example 1 from version 1.2 (2018-11-23) of the document
// “Syntaxdefinition der Rechnungsinformationen (S1) bei der QR-Rechnung”
func TestStructuredExample1(t *testing.T) {
	msg := BillInformation{
		InvoiceNumber:     "10201409",
		InvoiceDate:       OneDate(2019, time.May, 12),
		CustomerReference: "1400.000-53",
		VATNumber:         "106017086",
		VATDates:          OneDate(2018, time.May, 8),
		VATRates: TaxRates{
			TaxRate{RatePercent: 7.7},
		},
		Conditions: PaymentConditions{
			PaymentCondition{DiscountPercent: 2, NumberOfDays: 10},
			PaymentCondition{DiscountPercent: 0, NumberOfDays: 30},
		},
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Expected no error; got %v", err)
		return
	}
	expected := "//S1/10/10201409/11/190512/20/1400.000-53/30/106017086/31/180508/32/7.7/40/2:10;0:30"
	actual := msg.ToString()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

// Example 2 from version 1.2 (2018-11-23) of the document
// “Syntaxdefinition der Rechnungsinformationen (S1) bei der QR-Rechnung”
func TestStructuredExample2(t *testing.T) {
	msg := BillInformation{
		InvoiceNumber: "10104",
		InvoiceDate:   OneDate(2018, time.February, 28),
		VATNumber:     "395856455",
		VATDates:      StartAndEndDate(2018, time.February, 26, 2018, time.February, 27),
		VATRates: TaxRates{
			TaxRate{RatePercent: 3.7, Amount: 400.19},
			TaxRate{RatePercent: 7.7, Amount: 553.39},
			TaxRate{RatePercent: 0, Amount: 14},
		},
		Conditions: PaymentConditions{
			PaymentCondition{DiscountPercent: 0, NumberOfDays: 30},
		},
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Expected no error; got %v", err)
		return
	}
	expected := "//S1/10/10104/11/180228/30/395856455/31/180226180227/32/3.7:400.19;7.7:553.39;0:14/40/0:30"
	actual := msg.ToString()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

// Example 3 from version 1.2 (2018-11-23) of the document
// “Syntaxdefinition der Rechnungsinformationen (S1) bei der QR-Rechnung”
func TestStructuredExample3(t *testing.T) {
	msg := BillInformation{
		InvoiceNumber:     "4031202511",
		InvoiceDate:       OneDate(2018, time.January, 7),
		CustomerReference: "61257233.4",
		VATNumber:         "105493567",
		VATRates: TaxRates{
			TaxRate{RatePercent: 8, Amount: 49.82},
		},
		VATImportTaxRates: TaxRates{
			TaxRate{RatePercent: 2.5, Amount: 14.85},
		},
		Conditions: PaymentConditions{
			PaymentCondition{DiscountPercent: 0, NumberOfDays: 30},
		},
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Expected no error; got %v", err)
		return
	}
	expected := "//S1/10/4031202511/11/180107/20/61257233.4/30/105493567/32/8:49.82/33/2.5:14.85/40/0:30"
	actual := msg.ToString()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestValidateStructuredMessage(t *testing.T) {
	var testdata = []struct {
		structured BillInformation
		message    string
	}{
		{
			structured: BillInformation{},
			message:    "",
		},
		{
			structured: BillInformation{
				InvoiceNumber: "№17",
			},
			message: "Rune U+2116 '№' not allowed",
		},
		{
			structured: BillInformation{
				InvoiceDate: OneDate(2017, time.May, 17),
			},
			message: "",
		},
		{
			structured: BillInformation{
				InvoiceDate: StartAndEndDate(2017, time.May, 17, 2017, time.May, 21),
			},
			message: "Invoice date may not have an end date",
		},
		{
			structured: BillInformation{
				CustomerReference: "Best customer…",
			},
			message: "Rune U+2026 '…' not allowed",
		},
		{
			structured: BillInformation{
				VATNumber: "106017086",
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATNumber: "CHE-106.017.086",
			},
			message: "VAT number may only contain digits 0-9",
		},
		{
			structured: BillInformation{
				VATDates: OneDate(2017, time.May, 17),
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATDates: StartAndEndDate(2017, time.May, 17, 2017, time.May, 21),
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATDates: StartAndEndDate(2017, time.May, 21, 2017, time.May, 17),
			},
			message: "End date must come after start date",
		},
		{
			structured: BillInformation{
				VATRates: TaxRates{TaxRate{RatePercent: 3.0}},
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATRates: TaxRates{
					TaxRate{RatePercent: 3.0, Amount: 10.0},
					TaxRate{RatePercent: 5.0, Amount: 20.0},
				},
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATRates: TaxRates{TaxRate{RatePercent: -3.0}},
			},
			message: "VAT tax rate may not be negative",
		},
		{
			structured: BillInformation{
				VATRates: TaxRates{
					TaxRate{RatePercent: 5.0, Amount: -20.0},
				},
			},
			message: "VAT amount may not be negative",
		},
		{
			structured: BillInformation{
				VATImportTaxRates: TaxRates{TaxRate{RatePercent: 3.0}},
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATImportTaxRates: TaxRates{
					TaxRate{RatePercent: 3.0, Amount: 10.0},
					TaxRate{RatePercent: 5.0, Amount: 20.0},
				},
			},
			message: "",
		},
		{
			structured: BillInformation{
				VATImportTaxRates: TaxRates{TaxRate{RatePercent: -3.0}},
			},
			message: "VAT tax rate may not be negative",
		},
		{
			structured: BillInformation{
				VATRates: TaxRates{
					TaxRate{RatePercent: 5.0, Amount: -20.0},
				},
			},
			message: "VAT amount may not be negative",
		},
		{
			structured: BillInformation{
				Conditions: PaymentConditions{
					PaymentCondition{NumberOfDays: 30},
				},
			},
			message: "",
		},
		{
			structured: BillInformation{
				Conditions: PaymentConditions{
					PaymentCondition{DiscountPercent: 2, NumberOfDays: 10},
					PaymentCondition{DiscountPercent: 0, NumberOfDays: 60},
				},
			},
			message: "",
		},
		{
			structured: BillInformation{
				Conditions: PaymentConditions{
					PaymentCondition{DiscountPercent: 3, NumberOfDays: 15},
					PaymentCondition{DiscountPercent: 0.5, NumberOfDays: 45},
					PaymentCondition{DiscountPercent: 0, NumberOfDays: 90},
				},
			},
			message: "",
		},
		{
			structured: BillInformation{
				Conditions: PaymentConditions{
					PaymentCondition{NumberOfDays: -30},
				},
			},
			message: "Number of days may not be negative",
		},
		{
			structured: BillInformation{
				Conditions: PaymentConditions{
					PaymentCondition{DiscountPercent: -2, NumberOfDays: 10},
				},
			},
			message: "Discount may not be negative",
		},
	}
	for i, data := range testdata {
		err := data.structured.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Item %v: expected error; got no error.", i)
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}
