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
	"strings"
)

// ValidateCharacterSet validates that s only contains characters that are
// allowed according to the Swiss Implementation Guidelines for Customer-Bank
// Messages Credit Transfer.
func ValidateCharacterSet(s string) error {
	for _, r := range s {
		if !strings.ContainsRune(validRunes, r) {
			return fmt.Errorf("Rune %#U not allowed in string: %v", r, s)
		}
	}
	return nil
}

// Only these runes are allowed on payment slips in Switzerland.
var validRunes = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	".,:'+-/()? !\"#%&*;<>÷=@_$£[]{}\\`´~" +
	"àáâäçèéêëìíîïñòóôöùúûüýßÀÁÂÄÇÈÉÊËÌÍÎÏÒÓÔÖÙÚÛÜÑ"
