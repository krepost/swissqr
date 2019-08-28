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
	"io/ioutil"
	"testing"

	"github.com/krepost/gopdf/pdf"
)

func TestExample1FromStandard(t *testing.T) {
	data := examplePayload1
	if err := data.Validate(); err != nil {
		t.Error(err)
	}
	doc := pdf.New()
	canvas := doc.NewPage(21.0*pdf.Cm, 29.7*pdf.Cm)
	if err := DrawInvoice(canvas, data, "de"); err != nil {
		t.Error(err)
	}
	canvas.Close()
	if err := doc.Encode(ioutil.Discard); err != nil {
		t.Error(err)
	}
}

func TestExample2FromStandard(t *testing.T) {
	data := examplePayload2
	if err := data.Validate(); err != nil {
		t.Error(err)
	}

	doc := pdf.New()
	canvas := doc.NewPage(21.0*pdf.Cm, 10.5*pdf.Cm)
	if err := DrawInvoice(canvas, data, "fr"); err != nil {
		t.Error(err)
	}
	canvas.Close()
	if err := doc.Encode(ioutil.Discard); err != nil {
		t.Error(err)
	}
}
