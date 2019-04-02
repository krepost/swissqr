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
	"strings"
)

const (
	paymentPartInformation = iota
	receiptPartInformation
)

// Paragraph represents a header followed by text lines on the payment slip.
type Paragraph struct {
	Heading string
	Lines   []string
}

// AmountSectionData represents the payment amount on the payment slip.
type AmountSectionData struct {
	CurrencyHeading string
	CurrencyValue   string
	AmountHeading   string
	AmountValue     string
}

// TitleSectionData represents the title sections on the payment slip.
type TitleSectionData struct {
	PaymentPart string
	Receipt     string
}

// AmountSection returns the payment amount.
func AmountSection(p Payload, language string) (AmountSectionData, error) {
	if err := p.Validate(); err != nil {
		return AmountSectionData{}, err
	}
	if err := checkLanguage(language); err != nil {
		return AmountSectionData{}, err
	}
	amt := AmountSectionData{
		CurrencyHeading: headings[currency][language],
		CurrencyValue:   p.CurrencyAmount.Currency,
		AmountHeading:   headings[amount][language],
	}
	if p.CurrencyAmount.Amount > 0.0 {
		s := fmt.Sprintf("%.2f", p.CurrencyAmount.Amount)
		dot := strings.IndexByte(s, '.')
		if dot > 0 {
			firstSpace := dot % 3
			spaced := s[:firstSpace]
			for j := firstSpace; j < dot; j += 3 {
				if spaced != "" {
					spaced = spaced + " "
				}
				spaced = spaced + s[j:j+3]
			}
			spaced = spaced + s[dot:]
			s = spaced
		}
		amt.AmountValue = s
	}
	return amt, nil
}

// TitleSection returns the titles of the receipt and payment parts.
func TitleSection(p Payload, language string) (TitleSectionData, error) {
	if err := p.Validate(); err != nil {
		return TitleSectionData{}, err
	}
	if err := checkLanguage(language); err != nil {
		return TitleSectionData{}, err
	}
	return TitleSectionData{
		PaymentPart: headings[paymentPart][language],
		Receipt:     headings[receipt][language],
	}, nil
}

// InformationSection returns a slice of paragraphs to be rendered on the
// payment slip. The value of info (paymentPartInformation or
// receiptPartInformation) determines if information for the payment part
// or receipt part should be returned.
func InformationSection(p Payload, language string,
	width float64, info int) ([]Paragraph, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := checkLanguage(language); err != nil {
		return nil, err
	}
	lines := []string{p.Account.IBAN.PrintCode}
	if payableTo, err := p.Creditor.ToLines(); err != nil {
		return nil, err
	} else {
		lines = append(lines, payableTo...)
	}
	sections := []Paragraph{Paragraph{
		Heading: headings[accountPayableTo][language],
		Lines:   reflowAtSpace(lines, width),
	}}
	switch v := p.Reference.Number.(type) {
	case *structref.ReferenceNumber, *structref.CreditorReference:
		lines := []string{v.PrintFormat()}
		sections = append(sections, Paragraph{
			Heading: headings[reference][language],
			Lines:   reflowAtSpace(lines, width),
		})
	}
	if info == paymentPartInformation {
		lines := []string{}
		if s := p.AdditionalInformation.UnstructuredMessage; len(s) > 0 {
			lines = append(lines, s)
		}
		if s := p.AdditionalInformation.StructuredMessage.ToString(); len(s) > 0 {
			lines = append(lines, s)
		}
		if len(lines) > 0 {
			sections = append(sections, Paragraph{
				Heading: headings[additionalInformation][language],
				Lines:   reflowAtSpace(lines, width),
			})
		}
	}
	// The debtor section is mandatory and should be drawn as
	// a box by the client if the debtor information is empty.
	if lines, err := p.UltimateDebtor.ToLines(); err != nil {
		return nil, err
	} else {
		if len(lines) == 0 {
			sections = append(sections, Paragraph{
				Heading: headings[payableByNameAddress][language],
				Lines:   []string{},
			})
		} else {
			sections = append(sections, Paragraph{
				Heading: headings[payableBy][language],
				Lines:   reflowAtSpace(lines, width),
			})
		}
	}
	return sections, nil
}

// BorderText returns the text that is to be printed above the QR invoice.
func BorderText(language string) (string, error) {
	if err := checkLanguage(language); err != nil {
		return "", err
	}
	return headings[pleaseSeparate][language], nil
}

// ToLines converts an Entity to a set of lines suitable for display
// on a payment slip. It is assumed that the Entity is valid.
func (e Entity) ToLines() ([]string, error) {
	lines := []string{}
	if len(e.Name) > 0 {
		lines = append(lines, e.Name)
		switch addr := e.Address.(type) {
		case CombinedAddress:
			if addr.AddressLine1 != "" {
				lines = append(lines, addr.AddressLine1)
			}
			lines = append(lines, addr.AddressLine2)
		case StructuredAddress:
			if len(addr.StreetName) > 0 {
				line := addr.StreetName
				if len(addr.BuildingNumber) > 0 {
					line = line + " " + addr.BuildingNumber
				}
				lines = append(lines, line)
			}
			lines = append(lines, addr.PostCode+" "+addr.TownName)
		}
	}
	return lines, nil
}
