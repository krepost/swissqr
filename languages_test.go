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

import "testing"

func TestLanguageSupported(t *testing.T) {
	err := checkLanguage("en")
	if err != nil {
		t.Error("Expected language en to be supported.")
	}
}

func TestLanguageNotSupported(t *testing.T) {
	err := checkLanguage("sv")
	if err == nil {
		t.Error("Expected language sv to not be supported.")
	}
}

func TestLanguageLookup(t *testing.T) {
	var languageTests = []struct {
		id       int
		language string
		expected string
	}{
		{currency, "de", "Währung"},
		{accountPayableTo, "fr", "Compte / Payable à"},
		{amount, "it", "Importo"},
		{payableByNameAddress, "en", "Payable by (name/address)"},
	}
	for _, testCase := range languageTests {
		actual, found := headings[testCase.id][testCase.language]
		if !found {
			t.Errorf("Not found: %v,%v", testCase.id, testCase.language)
		}
		if testCase.expected != actual {
			t.Errorf("Expected %v, got %v", testCase.expected, actual)
		}
	}
}
