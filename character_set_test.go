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

func TestValidCharacters(t *testing.T) {
	testData := []string{
		"Pia Rutschmann",
		"Grossmünster [#5]",
		"Señor",
		"Fußschweiß",
		"peter@muster.ch",
	}
	for _, s := range testData {
		if err := ValidateCharacterSet(s); err != nil {
			t.Errorf("Expected no error; got %v", err)
		}
	}
}

func TestInvalidCharacters(t *testing.T) {
	testData := []struct{ s, invalid string }{
		{"sær", "U+00E6 'æ'"},
		{"øl", "U+00F8 'ø'"},
		{"Григорий", "U+0413 'Г'"},
	}
	for _, data := range testData {
		if err := ValidateCharacterSet(data.s); err == nil {
			t.Errorf("Expected error due to invalid character %v; got no error.", data.invalid)
		} else {
			if !strings.Contains(err.Error(), data.invalid) {
				t.Errorf("Expected error due to invalid character; got %v", err)
			}
		}
	}
}
