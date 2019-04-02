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
	"fmt"
	"github.com/krepost/structref"
	"io"
	"strings"
)

// Serialize serializes the payload data to w in a form that can
// be encoded in a Swiss QR Code. The payload is validated before
// serialization.
func (p Payload) Serialize(w io.Writer) error {
	if err := p.Validate(); err != nil {
		return err
	}
	io.WriteString(w, "SPC\r\n0200\r\n1\r\n") // Header.
	if err := p.Account.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.Creditor.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.UltimateCreditor.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.CurrencyAmount.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.UltimateDebtor.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.Reference.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.AdditionalInformation.Serialize(w); err != nil {
		return err
	}
	io.WriteString(w, "\r\n")
	if err := p.AlternativeProcedureParameters.Serialize(w); err != nil {
		return err
	}
	return nil
}

// Serialize serializes an account record.
// It is assumed that the record is valid.
func (a AccountNumber) Serialize(w io.Writer) error {
	_, err := io.WriteString(w, a.IBAN.Code)
	return err
}

// Serialize serializes an entity record.
// It is assumed that the record is valid.
func (e Entity) Serialize(w io.Writer) error {
	if e.Name == "" {
		_, err := io.WriteString(w, "\r\n\r\n\r\n\r\n\r\n\r\n")
		return err
	}
	if err := e.Address.Serialize(e.Name, w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "\r\n"+e.CountryCode); err != nil {
		return err
	}
	return nil
}

// Serialize serializes a combined address record.
// It is assumed that the record is valid.
func (ca CombinedAddress) Serialize(name string, w io.Writer) error {
	_, err := io.WriteString(w, strings.Join([]string{
		"K", name, ca.AddressLine1, ca.AddressLine2, "", "",
	}, "\r\n"))
	return err
}

// Serialize serializes a structured address record.
// It is assumed that the record is valid.
func (sa StructuredAddress) Serialize(name string, w io.Writer) error {
	_, err := io.WriteString(w, strings.Join([]string{
		"S", name, sa.StreetName, sa.BuildingNumber, sa.PostCode, sa.TownName,
	}, "\r\n"))
	return err
}

// Serialize serializes a payment amount record.
// It is assumed that the record is valid.
func (pa PaymentAmount) Serialize(w io.Writer) error {
	s := ""
	if pa.Amount > 0.0 {
		s = s + fmt.Sprintf("%.2f", pa.Amount)
	}
	s = s + "\r\n" + pa.Currency
	_, err := io.WriteString(w, s)
	return err
}

// Serialize serializes a payment reference record.
// It is assumed that the record is valid.
func (pr PaymentReference) Serialize(w io.Writer) error {
	var err error = nil
	switch ref := pr.Number.(type) {
	case *structref.ReferenceNumber:
		_, err = io.WriteString(w, "QRR\r\n"+ref.DigitalFormat())
	case *structref.CreditorReference:
		_, err = io.WriteString(w, "SCOR\r\n"+ref.DigitalFormat())
	case nil:
		_, err = io.WriteString(w, "NON\r\n")
	}
	return err
}

// Serialize serializes additional payment information record.
// It is assumed that the record is valid.
func (pi PaymentInformation) Serialize(w io.Writer) error {
	s := pi.UnstructuredMessage + "\r\nEPD\r\n" + pi.StructuredMessage.ToString()
	_, err := io.WriteString(w, s)
	return err
}

// Serialize serializes alternative procedure parameters.
// It is assumed that the parameters are valid.
func (vec AlternativeProcedures) Serialize(w io.Writer) error {
	apJoin := []string{"", ""}
	for i, ap := range vec {
		apJoin[i] = ap.Procedure
	}
	_, err := io.WriteString(w, strings.Join(apJoin, "\r\n"))
	return err
}
