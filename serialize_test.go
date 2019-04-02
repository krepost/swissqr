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
	"bytes"
	"github.com/krepost/structref"
	"testing"
)

func TestSerializeAccountNumber(t *testing.T) {
	data := NewIBANOrDie("CH56 0483 5012 3456 7800 9")
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "CH5604835012345678009"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializeEmtpyEntity(t *testing.T) {
	data := Entity{}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "\r\n\r\n\r\n\r\n\r\n\r\n"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializeCombinedAddress(t *testing.T) {
	data := Entity{
		Name: "Test Name",
		Address: CombinedAddress{
			AddressLine1: "Address Line 1",
			AddressLine2: "Address Line 2",
		},
		CountryCode: "CH",
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "K\r\n" +
		"Test Name\r\n" +
		"Address Line 1\r\n" +
		"Address Line 2\r\n" +
		"\r\n" +
		"\r\n" +
		"CH"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializeStructuredAddress(t *testing.T) {
	data := Entity{
		Name: "Test Name",
		Address: StructuredAddress{
			StreetName:     "Street Name",
			BuildingNumber: "17",
			PostCode:       "1234",
			TownName:       "Town Name",
		},
		CountryCode: "CH",
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "S\r\n" +
		"Test Name\r\n" +
		"Street Name\r\n" +
		"17\r\n" +
		"1234\r\n" +
		"Town Name\r\n" +
		"CH"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePaymentAmount(t *testing.T) {
	data := PaymentAmount{
		Amount:   1234.5678,
		Currency: CHF,
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "1234.57\r\nCHF"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePaymentReferenceEmpty(t *testing.T) {
	data := PaymentReference{}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "NON\r\n"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePaymentReferenceESR(t *testing.T) {
	num := "210000000003139471430009017"
	data := PaymentReference{
		Number: structref.NewReferenceNumberOrDie(num),
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "QRR\r\n" + num
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePaymentReferenceISO(t *testing.T) {
	num := "RF8312345678912345678912"
	data := PaymentReference{
		Number: structref.NewCreditorReferenceOrDie(num),
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "SCOR\r\n" + num
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePaymentInformation(t *testing.T) {
	data := PaymentInformation{
		UnstructuredMessage: "Unstructured message",
		StructuredMessage:   BillInformation{CustomerReference: "ref"},
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "Unstructured message\r\nEPD\r\n//S1/20/ref"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializeAlternativeProcedures(t *testing.T) {
	data := AlternativeProcedures{
		AlternativeProcedure{
			Label:     "Label 1",
			Procedure: "Procedure 1",
		},
		AlternativeProcedure{
			Label:     "Label 2",
			Procedure: "Procedure 2",
		},
	}
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "Procedure 1\r\nProcedure 2"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePayloadExample1(t *testing.T) {
	data := examplePayload1
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		t.Errorf("Could not serialize payload: %v", err)
	}
	expected := "SPC\r\n" +
		"0200\r\n" +
		"1\r\n" +
		"CH5800791123000889012\r\n" +
		"S\r\n" +
		"Robert Schneider AG\r\n" +
		"Rue du Lac\r\n" +
		"1268\r\n" +
		"2501\r\n" +
		"Biel\r\n" +
		"CH\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"3949.75\r\n" +
		"CHF\r\n" +
		"S\r\n" +
		"Pia Rutschmann\r\n" +
		"Marktgasse\r\n" +
		"28\r\n" +
		"9400\r\n" +
		"Rorschach\r\n" +
		"CH\r\n" +
		"NON\r\n" +
		"\r\n" +
		"Rechnung Nr. 3139 für Gartenarbeiten und Entsorgung Schnittmaterial\r\n" +
		"EPD\r\n" +
		"\r\n" +
		"\r\n"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestSerializePayloadExample3(t *testing.T) {
	data := examplePayload3
	var buffer bytes.Buffer
	data.Serialize(&buffer)
	expected := "SPC\r\n" +
		"0200\r\n" +
		"1\r\n" +
		"CH3709000000304442225\r\n" +
		"S\r\n" +
		"Salvation Army Foundation Switzerland\r\n" +
		"\r\n" +
		"\r\n" +
		"3000\r\n" +
		"Bern\r\n" +
		"CH\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"CHF\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"\r\n" +
		"NON\r\n" +
		"\r\n" +
		"Donation to the Winterfest Campaign\r\n" +
		"EPD\r\n" +
		"\r\n" +
		"\r\n"
	actual := buffer.String()
	if expected != actual {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}
