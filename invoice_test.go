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
	"reflect"
	"testing"
)

func TestTitleSectionExample1De(t *testing.T) {
	expected := TitleSectionData{
		PaymentPart: "Zahlteil",
		Receipt:     "Empfangsschein",
	}
	actual, err := TitleSection(examplePayload1, "de")
	if err != nil {
		t.Errorf("Could not create title section: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestTitleSectionExample3En(t *testing.T) {
	expected := TitleSectionData{
		PaymentPart: "Payment part",
		Receipt:     "Receipt",
	}
	actual, err := TitleSection(examplePayload3, "en")
	if err != nil {
		t.Errorf("Could not create invoice text: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestAmountSectionExample1De(t *testing.T) {
	expected := AmountSectionData{
		CurrencyHeading: "Währung",
		CurrencyValue:   "CHF",
		AmountHeading:   "Betrag",
		AmountValue:     "3 949.75",
	}
	actual, err := AmountSection(examplePayload1, "de")
	if err != nil {
		t.Errorf("Could not create title section: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestAmountSectionExample3En(t *testing.T) {
	expected := AmountSectionData{
		CurrencyHeading: "Currency",
		CurrencyValue:   "CHF",
		AmountHeading:   "Amount",
		AmountValue:     "",
	}
	actual, err := AmountSection(examplePayload3, "en")
	if err != nil {
		t.Errorf("Could not create invoice text: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestInformationSectionExample1De(t *testing.T) {
	expected := []Paragraph{
		Paragraph{
			Heading: "Konto / Zahlbar an",
			Lines: []string{
				"CH58 0079 1123 0008 8901 2",
				"Robert Schneider AG",
				"Rue du Lac 1268",
				"2501 Biel",
			},
		},
		Paragraph{
			Heading: "Zusätzliche Informationen",
			Lines: []string{
				"Rechnung Nr. 3139 für Gartenarbeiten und",
				"Entsorgung Schnittmaterial",
			},
		},
		Paragraph{
			Heading: "Zahlbar durch",
			Lines: []string{
				"Pia Rutschmann",
				"Marktgasse 28",
				"9400 Rorschach",
			},
		},
	}
	actual, err := InformationSection(examplePayload1, "de",
		8.5*28.35/10.0, // 8.5cm × 28.35 pt/cm ÷ 10pt font size.
		paymentPartInformation)
	if err != nil {
		t.Errorf("Could not create invoice text: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestInformationSectionExample3En(t *testing.T) {
	expected := []Paragraph{
		Paragraph{
			Heading: "Account / Payable to",
			Lines: []string{
				"CH37 0900 0000 3044 4222 5",
				"Salvation Army Foundation Switzerland",
				"3000 Bern",
			},
		},
		Paragraph{Heading: "Payable by (name/address)", Lines: []string{}},
	}
	actual, err := InformationSection(examplePayload3, "en",
		8.5*28.35/10.0, // 8.5cm × 28.35 pt/cm ÷ 10pt font size.
		receiptPartInformation)
	if err != nil {
		t.Errorf("Could not create invoice text: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected:\n\n%#v\n\nGot:\n\n%#v\n\n", expected, actual)
	}
}

func TestEntityToLines(t *testing.T) {
	entity := Entity{
		Name: "Pia-Maria Rutschmann-Schnyder",
		Address: StructuredAddress{
			StreetName:     "Grosse Marktgasse",
			BuildingNumber: "28",
			PostCode:       "9400",
			TownName:       "Rorschach",
		},
		CountryCode: "CH",
	}
	expected := []string{
		"Pia-Maria Rutschmann-Schnyder",
		"Grosse Marktgasse 28",
		"9400 Rorschach",
	}
	actual, err := entity.ToLines()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestNoBuildingNumber(t *testing.T) {
	entity := Entity{
		Name: "Pia-Maria Rutschmann-Schnyder",
		Address: StructuredAddress{
			StreetName: "Grosse Marktgasse",
			PostCode:   "9400",
			TownName:   "Rorschach",
		},
		CountryCode: "CH",
	}
	expected := []string{
		"Pia-Maria Rutschmann-Schnyder",
		"Grosse Marktgasse",
		"9400 Rorschach",
	}
	actual, err := entity.ToLines()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestNoStreetAddress(t *testing.T) {
	entity := Entity{
		Name: "Pia-Maria Rutschmann-Schnyder",
		Address: StructuredAddress{
			PostCode: "9400",
			TownName: "Rorschach",
		},
		CountryCode: "CH",
	}
	expected := []string{
		"Pia-Maria Rutschmann-Schnyder",
		"9400 Rorschach",
	}
	actual, err := entity.ToLines()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestCombinedAddress(t *testing.T) {
	entity := Entity{
		Name: "Pia-Maria Rutschmann-Schnyder",
		Address: CombinedAddress{
			AddressLine1: "Grosse Marktgasse 28",
			AddressLine2: "9400 Rorschach",
		},
	}
	expected := []string{
		"Pia-Maria Rutschmann-Schnyder",
		"Grosse Marktgasse 28",
		"9400 Rorschach",
	}
	actual, err := entity.ToLines()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}
