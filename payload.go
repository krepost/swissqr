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

// package swissqr creates a QR code for electronic bills as defined in version
// 2.0 of the document “Schweizer Implementation Guidelines QR-Rechnung”, which
// can be downloaded from https://www.paymentstandards.ch/, and version 1.2 of
// the document “Syntaxdefinition der Rechnungsinformationen (S1) bei der QR-
// Rechnung”, which can be downloaded from https://www.swiss-qr-invoice.org/.
package swissqr

import (
	"github.com/almerlucke/go-iban/iban"
	"github.com/krepost/structref"
	"io"
)

const (
	CHF = "CHF"
	EUR = "EUR"
)

// Payload contains all information that can be encoded in the Swiss
// QR Code for invoices.
type Payload struct {
	// IBAN or QR-IBAN of the creditor, according to ISO 13616.
	Account AccountNumber

	// The creditor. Mandatory data group.
	Creditor Entity

	// Information about the ultimate creditor. For future use.
	UltimateCreditor Entity

	// The payment amount in a given currency. Mandatory data group.
	CurrencyAmount PaymentAmount

	// The ultimate debtor. Optional data group.
	UltimateDebtor Entity

	// Reference contains an optional structured payment reference number.
	Reference PaymentReference

	// AdditionalInformation can be used to send additional
	// information to the biller. Optional data group.
	AdditionalInformation PaymentInformation

	// AlternativeProcedureParameters contains a maximum of two
	// entries describing the parameter character chain of the
	// alternative scheme according to the syntax definition in
	// the section on “Alternative procedure” in the Swiss QR standard.
	AlternativeProcedureParameters AlternativeProcedures
}

type qrAddress interface {
	Validate() error
	Serialize(Name string, w io.Writer) error
}

// Entity can represent a creditor or a debtor.
type Entity struct {
	// Mandatory field. Contains first name (optional) and last name
	// or company name.
	Name string

	// Mandatory field. Must contain either a CombinedAddress
	// or a StructuredAddress.
	Address qrAddress

	// Mandatory two-letter country code according to ISO 3166-1.
	CountryCode string
}

// CombinedAdress represents an unstructured address.
type CombinedAddress struct {
	// Optional street name and building number, or PO Box.
	// Maximum 70 characters allowed.
	AddressLine1 string

	// Mandatory field containing postal code and town.
	// Maximum 70 characters allowed.
	AddressLine2 string
}

// StructuredAddress represents an address with structured fields.
type StructuredAddress struct {
	// Optional. May not include house or building number.
	// Maximum 70 characters allowed.
	StreetName string

	// Optional. Maximum 16 characters allowed.
	BuildingNumber string

	// Mandatory. Is always to be entered without a country code.
	// Maximum 16 characters allowed.
	PostCode string

	// Mandatory. Maximum 35 characters allowed.
	TownName string
}

// PaymentAmount contains the amount to be paid.
type PaymentAmount struct {
	// Optional payment amount.
	Amount float64

	// Mandatory payment currency. Only "CHF" and "EUR" are permitted.
	Currency string
}

// Account contains an IBAN or QR-IBAN.
type AccountNumber struct {
	IBAN *iban.IBAN
}

// NewIBANOrDie is a helper function to set the IBAN field in Payload.
// Useful when initializing a Payload struct programmatically.
func NewIBANOrDie(s string) AccountNumber {
	iban, err := iban.NewIBAN(s)
	if err != nil {
		panic(err)
	}
	return AccountNumber{IBAN: iban}
}

// PaymentReference contains either a Swiss ESR reference number,
// or a structured creditor reference according to ISO 11649, or nil.
type PaymentReference struct {
	Number structref.Printer
}

// PaymentInformation includes additional unstructured or coded
// information about the payment. The two fields combined may contain
// at most 140 characters.
type PaymentInformation struct {
	// UnstructuredMessage can be used to indicate the payment purpose
	// or for additional tectual information about payments with a
	// structured reference. Optional field.
	UnstructuredMessage string

	// StructuredInformation contains coded information for automated
	// booking of the payment. Optional field.
	StructuredMessage BillInformation
}

// AlternativeProcedure defines an alternateive payment procedure.
type AlternativeProcedure struct {
	// Label is shown in bold on the QR invoice.
	Label string

	// Procedure describes the alternative payment procedure.
	Procedure string
}

// AlternativeProcedures is a list of alternative payment procedures.
type AlternativeProcedures []AlternativeProcedure
