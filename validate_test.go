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
	"github.com/krepost/structref"
	"io"
	"strings"
	"testing"
)

var (
	minimalCorrectEntity = Entity{
		Name:        "Test Creditor",
		Address:     CombinedAddress{AddressLine2: "Test Address"},
		CountryCode: "CH",
	}

	minimalCorrectPayload = Payload{
		Account:        NewIBANOrDie("CH5604835012345678009"),
		Creditor:       minimalCorrectEntity,
		CurrencyAmount: PaymentAmount{Currency: CHF},
	}
)

func TestValidateMinimalCorrect(t *testing.T) {
	if err := minimalCorrectEntity.Validate(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := minimalCorrectPayload.Validate(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateCreditorPresent(t *testing.T) {
	payload := minimalCorrectPayload
	payload.Creditor = Entity{}
	err := payload.Validate()
	// While an Entity field can be empty in general, the Creditor field cannot.
	if err == nil || err.Error() != "No creditor name specified." {
		t.Error("Expected error due to no creditor.")
	}
}

func TestValidateNoUltimateCreditorPresent(t *testing.T) {
	payload := minimalCorrectPayload
	payload.UltimateCreditor = payload.Creditor
	err := payload.Validate()
	// The UltimateCreditor field is currently not supported.
	if err == nil || err.Error() != "UltimateCreditor is currently not supported." {
		t.Error("Expected error due to ultimate creditor.")
	}
}

func TestValidateCrossFieldDependencies(t *testing.T) {
	var testdata = []struct {
		account   AccountNumber
		reference PaymentReference
		message   string
	}{
		// If a QR-IBAN is used, the Reference field must contain a QRReference
		// code. A QR-IBAN has a bank clearing number (first five digits of the
		// IBAN itself, after country code and check sum digits) between 30000 and
		// 31999.
		{
			account:   NewIBANOrDie("CH4431999123000889012"),
			reference: PaymentReference{},
			message:   "QR Reference number required for QR-IBAN",
		},
		{
			account: NewIBANOrDie("CH4431999123000889012"),
			reference: PaymentReference{
				Number: structref.NewCreditorReferenceOrDie("RF8312345678912345678912"),
			},
			message: "QR Reference number required for QR-IBAN",
		},
		{
			account: NewIBANOrDie("CH4431999123000889012"),
			reference: PaymentReference{
				Number: structref.NewReferenceNumberOrDie("210000000003139471430009017"),
			},
			message: "",
		},
		{
			// If a regular IBAN is used, the Reference field must be either empty or
			// contain a creditor reference code.
			account:   NewIBANOrDie("CH5604835012345678009"),
			reference: PaymentReference{},
			message:   "",
		},
		{
			account: NewIBANOrDie("CH5604835012345678009"),
			reference: PaymentReference{
				Number: structref.NewCreditorReferenceOrDie("RF8312345678912345678912"),
			},
			message: "",
		},
		{
			account: NewIBANOrDie("CH5604835012345678009"),
			reference: PaymentReference{
				Number: structref.NewReferenceNumberOrDie("210000000003139471430009017"),
			},
			message: "QR Reference not allowed for IBAN",
		},
	}
	for i, data := range testdata {
		payload := minimalCorrectPayload
		payload.Account = data.account
		payload.Reference = data.reference
		err := payload.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

// Now follow table-driven tests for validation of parts of struct Payload.

type testAddress struct{}

func (t testAddress) Serialize(name string, w io.Writer) error { return nil }

func (t testAddress) Validate() error { return nil }

func TestValidateEntity(t *testing.T) {
	var testdata = []struct {
		entity  Entity
		message string
	}{
		{
			entity:  Entity{},
			message: "",
		},
		{
			entity:  Entity{Name: "Name"},
			message: "Country code must be specified for name",
		},
		{
			entity:  Entity{Name: "Name", CountryCode: "CH"},
			message: "Unsupported address type",
		},
		{
			entity: Entity{
				Name:        "Name",
				Address:     testAddress{},
				CountryCode: "CH",
			},
			message: "Unsupported address type",
		},
		{
			entity: Entity{
				Name:        "Name",
				Address:     StructuredAddress{PostCode: "8000", TownName: "Zürich"},
				CountryCode: "CH",
			},
			message: "",
		},
		{
			entity: Entity{
				Name:        "Name",
				Address:     CombinedAddress{AddressLine2: "8000 Zürich"},
				CountryCode: "CH",
			},
			message: "",
		},
		{
			entity: Entity{
				Name:        "Næjm",
				Address:     CombinedAddress{AddressLine2: "8000 Zürich"},
				CountryCode: "CH",
			},
			message: "Rune U+00E6 'æ' not allowed",
		},
		{
			entity: Entity{
				Name:        strings.Repeat("Name", 18),
				Address:     CombinedAddress{AddressLine2: "8000 Zürich"},
				CountryCode: "CH",
			},
			message: "Maximum name length is 70 characters",
		},
		{
			entity: Entity{
				Name:        "Name",
				Address:     CombinedAddress{AddressLine2: "8000 Zürich"},
				CountryCode: "CHE",
			},
			message: "Country should be given as two-letter code",
		},
		{
			entity: Entity{
				Name:        "Name",
				Address:     CombinedAddress{AddressLine2: "8000 Zürich"},
				CountryCode: "рф",
			},
			message: "Invalid country code",
		},
	}
	for i, data := range testdata {
		err := data.entity.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidateStructuredAddress(t *testing.T) {
	var testdata = []struct {
		address StructuredAddress
		message string
	}{
		{
			address: StructuredAddress{},
			message: "Must specify post code and town in address",
		},
		{
			address: StructuredAddress{PostCode: "code", TownName: "town"},
			message: "",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "no",
				PostCode:       "code",
				TownName:       "town",
			},
			message: "",
		},
		{
			address: StructuredAddress{
				StreetName:     "strĳt",
				BuildingNumber: "no",
				PostCode:       "code",
				TownName:       "town",
			},
			message: "Rune U+0133 'ĳ' not allowed",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "nø",
				PostCode:       "code",
				TownName:       "town",
			},
			message: "Rune U+00F8 'ø' not allowed",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "no",
				PostCode:       "cœde",
				TownName:       "town",
			},
			message: "Rune U+0153 'œ' not allowed",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "no",
				PostCode:       "code",
				TownName:       "tøwn",
			},
			message: "Rune U+00F8 'ø' not allowed",
		},
		{
			address: StructuredAddress{
				StreetName:     strings.Repeat("street", 12),
				BuildingNumber: "no",
				PostCode:       "code",
				TownName:       "town",
			},
			message: "Maximum street name length is 70 characters",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "12345678901234567",
				PostCode:       "code",
				TownName:       "town",
			},
			message: "Maximum building number length is 16 characters",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "no",
				PostCode:       "12345678901234567",
				TownName:       "town",
			},
			message: "Maximum post code length is 16 characters",
		},
		{
			address: StructuredAddress{
				StreetName:     "street",
				BuildingNumber: "no",
				PostCode:       "code",
				TownName:       strings.Repeat("town", 9),
			},
			message: "Maximum town name length is 35 characters",
		},
	}
	for i, data := range testdata {
		err := data.address.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidateCombinedAddress(t *testing.T) {
	var testdata = []struct {
		address CombinedAddress
		message string
	}{
		{
			address: CombinedAddress{},
			message: "Address line 2 must be set for address",
		},
		{
			address: CombinedAddress{AddressLine1: "Line 1"},
			message: "Address line 2 must be set for address",
		},
		{
			address: CombinedAddress{AddressLine2: "Line 2"},
			message: "",
		},
		{
			address: CombinedAddress{
				AddressLine1: "Line 1",
				AddressLine2: "Line 2",
			},
			message: "",
		},
		{
			address: CombinedAddress{
				AddressLine1: "Line¹",
				AddressLine2: "Line 2",
			},
			message: "Rune U+00B9 '¹' not allowed",
		},
		{
			address: CombinedAddress{
				AddressLine1: "Line 1",
				AddressLine2: "Line²",
			},
			message: "Rune U+00B2 '²' not allowed",
		},
		{
			address: CombinedAddress{
				AddressLine1: strings.Repeat("Line 1", 11),
				AddressLine2: strings.Repeat("Line 2", 11),
			},
			message: "",
		},
		{
			address: CombinedAddress{
				AddressLine1: strings.Repeat("Line 1", 10),
				AddressLine2: strings.Repeat("Line 2", 12),
			},
			message: "Maximum address line length is 70 characters",
		},
		{
			address: CombinedAddress{
				AddressLine1: strings.Repeat("Line 1", 12),
				AddressLine2: strings.Repeat("Line 2", 10),
			},
			message: "Maximum address line length is 70 characters",
		},
	}
	for i, data := range testdata {
		err := data.address.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidateAccount(t *testing.T) {
	var testdata = []struct {
		account AccountNumber
		message string
	}{
		{
			account: AccountNumber{},
			message: "No account specified",
		},
		{
			account: NewIBANOrDie("CH5604835012345678009"),
			message: "",
		},
		{
			account: NewIBANOrDie("DE91100000000123456789"),
			message: "Only CH and LI accounts allowed",
		},
	}
	for i, data := range testdata {
		err := data.account.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidatePaymentAmount(t *testing.T) {
	var testdata = []struct {
		amount  PaymentAmount
		message string
	}{
		{
			amount:  PaymentAmount{},
			message: "Currency must be CHF or EUR",
		},
		{
			amount:  PaymentAmount{Currency: "SEK"},
			message: "Currency must be CHF or EUR",
		},
		{
			amount:  PaymentAmount{Amount: 17.0},
			message: "Currency must be CHF or EUR",
		},
		{
			amount:  PaymentAmount{Currency: CHF},
			message: "",
		},
		{
			amount:  PaymentAmount{Currency: EUR},
			message: "",
		},
		{
			amount:  PaymentAmount{Currency: CHF, Amount: 17.0},
			message: "",
		},
		{
			amount:  PaymentAmount{Currency: CHF, Amount: -17.0},
			message: "Amount cannot be negative",
		},
		{
			amount:  PaymentAmount{Currency: CHF, Amount: 1234567890.0},
			message: "Amount too large",
		},
	}
	for i, data := range testdata {
		err := data.amount.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

// fakeReference implements interface structref.Printer.
type fakeReference struct{}

func (r fakeReference) DigitalFormat() string { return "" }
func (r fakeReference) PrintFormat() string   { return "" }

func TestValidatePaymentReference(t *testing.T) {
	var testdata = []struct {
		ref     PaymentReference
		message string
	}{
		{
			ref:     PaymentReference{},
			message: "",
		},
		{
			ref:     PaymentReference{Number: fakeReference{}},
			message: "Unknown reference type",
		},
		{
			ref: PaymentReference{
				Number: structref.NewReferenceNumberOrDie("210000000003139471430009017"),
			},
			message: "",
		},
		{
			ref: PaymentReference{
				Number: structref.NewCreditorReferenceOrDie("RF8312345678912345678912"),
			},
			message: "",
		},
	}
	for i, data := range testdata {
		err := data.ref.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidatePaymentInformation(t *testing.T) {
	var testdata = []struct {
		info    PaymentInformation
		message string
	}{
		{
			info:    PaymentInformation{},
			message: "",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: "Unstructured",
			},
			message: "",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: "Øystein",
			},
			message: "Rune U+00D8 'Ø' not allowed",
		},
		{
			info: PaymentInformation{
				StructuredMessage: BillInformation{CustomerReference: "ref"},
			},
			message: "",
		},
		{
			info: PaymentInformation{
				StructuredMessage: BillInformation{CustomerReference: "ræf"},
			},
			message: "Rune U+00E6 'æ' not allowed",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: "Unstructured",
				StructuredMessage:   BillInformation{CustomerReference: "ref"},
			},
			message: "",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: strings.Repeat("A long message.", 7),
				StructuredMessage:   BillInformation{CustomerReference: "ref"},
			},
			message: "",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: "Unstructured",
				StructuredMessage:   BillInformation{CustomerReference: "ref"},
			},
			message: "",
		},
		{
			info: PaymentInformation{
				UnstructuredMessage: strings.Repeat("A long message.", 9),
				StructuredMessage:   BillInformation{CustomerReference: "ref"},
			},
			message: "Maximum combined length is 140",
		},
	}
	for i, data := range testdata {
		err := data.info.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}

func TestValidateAlternativeProcedures(t *testing.T) {
	var testdata = []struct {
		procedures AlternativeProcedures
		message    string
	}{
		{
			procedures: AlternativeProcedures{},
			message:    "",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "X"},
			},
			message: "No procedure specified",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Procedure: "Procedure X"},
			},
			message: "No label specified",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "X", Procedure: "Procedure X"},
			},
			message: "",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "X", Procedure: "Procedure X"},
				AlternativeProcedure{Label: "Y", Procedure: "Procedure Y"},
			},
			message: "",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "Ü", Procedure: "Procedure X"},
				AlternativeProcedure{Label: "Ĳ", Procedure: "Procedure Y"},
			},
			message: "Rune U+0132 'Ĳ' not allowed",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "X", Procedure: "Prøcedure X"},
			},
			message: "Rune U+00F8 'ø' not allowed",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{Label: "X", Procedure: "Procedure X"},
				AlternativeProcedure{Label: "Y", Procedure: "Procedure Y"},
				AlternativeProcedure{Label: "Z", Procedure: "Procedure Z"},
			},
			message: "Maximum two alternate payment schemes allowed",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{
					Label:     "Long",
					Procedure: strings.Repeat("Procedure", 11),
				},
			},
			message: "",
		},
		{
			procedures: AlternativeProcedures{
				AlternativeProcedure{
					Label:     "Long",
					Procedure: strings.Repeat("Procedure", 12),
				},
			},
			message: "Maximum field length is 100 characters",
		},
	}
	for i, data := range testdata {
		err := data.procedures.Validate()
		if data.message == "" {
			if err != nil {
				t.Errorf("Item %v: expected no error; got %v", i, err)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error; got no error.")
			} else if !strings.Contains(err.Error(), data.message) {
				t.Errorf("Item %v: expected error %#v, got: %v", i, data.message, err)
			}
		}
	}
}
