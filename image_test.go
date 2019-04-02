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
)

func TestCreateImageFromValidData(t *testing.T) {
	exampleCorrect := Payload{
		Account: NewIBANOrDie("CH5604835012345678009"),
		Creditor: Entity{
			Name:        "Test Creditor",
			Address:     CombinedAddress{AddressLine2: "Test Address"},
			CountryCode: "CH",
		},
		CurrencyAmount: PaymentAmount{Currency: CHF},
	}
	img, err := CreateQR(exampleCorrect)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	rect := img.Bounds()
	if rect.Max.X-rect.Min.X != 1086 || rect.Max.Y-rect.Min.Y != 1086 {
		t.Errorf("Unexpected image size: %v", rect)
	}
}

func TestCreateImageFromInvalidData(t *testing.T) {
	exampleNoCreditor := Payload{
		Account:        NewIBANOrDie("CH5800791123000889012"),
		CurrencyAmount: PaymentAmount{Currency: EUR},
	}
	_, err := CreateQR(exampleNoCreditor)
	if err == nil {
		t.Error("Expected error due to missing creditor")
	}
	if !strings.HasPrefix(err.Error(), "No creditor name specified") {
		t.Errorf("Expected error due to no creditor name; got %v", err)
	}
}
