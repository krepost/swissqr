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
	"errors"
	"fmt"
	"github.com/krepost/structref"
)

// Validate validates the payload and returns nil on success.
func (p Payload) Validate() error {
	if err := p.Account.Validate(); err != nil {
		return err
	}
	if err := p.Creditor.Validate(); err != nil {
		return err
	}
	if err := p.UltimateCreditor.Validate(); err != nil {
		return err
	}
	if err := p.CurrencyAmount.Validate(); err != nil {
		return err
	}
	if err := p.UltimateDebtor.Validate(); err != nil {
		return err
	}
	if err := p.Reference.Validate(); err != nil {
		return err
	}
	if err := p.AdditionalInformation.Validate(); err != nil {
		return err
	}
	if err := p.AlternativeProcedureParameters.Validate(); err != nil {
		return err
	}
	// The Creditor field must not be empty, but this is not checked by the
	// Validate() method since other Entity fields can be empty.
	if p.Creditor.Name == "" {
		return errors.New("No creditor name specified.")
	}
	// The UltimateCreditor field is reserved for future use.
	if p.UltimateCreditor.Name != "" {
		return errors.New("UltimateCreditor is currently not supported.")
	}
	// If a QR-IBAN is used, Reference must contain a QRReference code.
	// Otherwise, either no reference or a Creditor Reference must be used.
	// A QR-IBAN has a bank clearing number (first five digits of the IBAN
	// itself, after country code and check sum digits) between 30000 and 31999.
	if p.Account.IBAN.BBAN[:2] == "30" || p.Account.IBAN.BBAN[:2] == "31" {
		if p.Reference.Number == nil {
			return fmt.Errorf("QR Reference number required for QR-IBAN: %v",
				p.Account.IBAN.PrintCode)
		}
		if _, ok := p.Reference.Number.(*structref.ReferenceNumber); !ok {
			return fmt.Errorf("QR Reference number required for QR-IBAN: %v",
				p.Account.IBAN.PrintCode)
		}
	} else {
		if p.Reference.Number != nil {
			if _, ok := p.Reference.Number.(*structref.ReferenceNumber); ok {
				return fmt.Errorf("QR Reference not allowed for IBAN: %v",
					p.Account.IBAN.PrintCode)
			}
		}
	}
	return nil
}

// Validate validates an Account
func (a AccountNumber) Validate() error {
	if a.IBAN == nil {
		return errors.New("No account specified")
	}
	if a.IBAN.CountryCode != "CH" && a.IBAN.CountryCode != "LI" {
		return fmt.Errorf("Only CH and LI accounts allowed: %v", a.IBAN.PrintCode)
	}
	return nil
}

// Validate validates an Entity
func (e Entity) Validate() error {
	// Empty record is allowed.
	if e.Name == "" && e.Address == nil && e.CountryCode == "" {
		return nil
	}

	// Name is mandatory for non-empty records.
	if e.Name == "" {
		return errors.New("Name must be specified.")
	}
	if len(e.Name) > 70 {
		return fmt.Errorf("Maximum name length is 70 characters: %v", e.Name)
	}
	if err := ValidateCharacterSet(e.Name); err != nil {
		return err
	}

	// Country code is mandatory.
	if e.CountryCode == "" {
		return fmt.Errorf("Country code must be specified for name: %v", e.Name)
	}
	if len([]rune(e.CountryCode)) > 2 {
		return fmt.Errorf("Country should be given as two-letter code: %v", e.CountryCode)
	}
	if _, found := countryCodes[e.CountryCode]; !found {
		return fmt.Errorf("Invalid country code: %v", e.CountryCode)
	}

	// Check address type and validate recursively.
	switch a := e.Address.(type) {
	case CombinedAddress, StructuredAddress:
		return a.Validate()
	default:
		return fmt.Errorf("Unsupported address type: %T", a)
	}

	return nil
}

// Validate validates a CombinedAddress.
func (ca CombinedAddress) Validate() error {
	// Combined address mode.
	if ca.AddressLine2 == "" {
		return fmt.Errorf("Address line 2 must be set for address: %v", ca)
	}
	if err := ValidateCharacterSet(ca.AddressLine1); err != nil {
		return err
	}
	if err := ValidateCharacterSet(ca.AddressLine2); err != nil {
		return err
	}
	if len(ca.AddressLine1) > 70 {
		return fmt.Errorf("Maximum address line length is 70 characters: %v", ca.AddressLine1)
	}
	if len(ca.AddressLine2) > 70 {
		return fmt.Errorf("Maximum address line length is 70 characters: %v", ca.AddressLine2)
	}
	return nil
}

// Validate validates a StructuredAddress.
func (sa StructuredAddress) Validate() error {
	if sa.PostCode == "" || sa.TownName == "" {
		return fmt.Errorf("Must specify post code and town in address: %v", sa)
	}
	if err := ValidateCharacterSet(sa.StreetName); err != nil {
		return err
	}
	if err := ValidateCharacterSet(sa.BuildingNumber); err != nil {
		return err
	}
	if err := ValidateCharacterSet(sa.PostCode); err != nil {
		return err
	}
	if err := ValidateCharacterSet(sa.TownName); err != nil {
		return err
	}
	if len(sa.StreetName) > 70 {
		return fmt.Errorf("Maximum street name length is 70 characters: %v", sa.StreetName)
	}
	if len(sa.BuildingNumber) > 16 {
		return fmt.Errorf("Maximum building number length is 16 characters: %v", sa.BuildingNumber)
	}
	if len(sa.PostCode) > 16 {
		return fmt.Errorf("Maximum post code length is 16 characters: %v", sa.PostCode)
	}
	if len(sa.TownName) > 35 {
		return fmt.Errorf("Maximum town name length is 35 characters: %v", sa.TownName)
	}
	return nil
}

// Validate validates a PaymentAmount.
func (pa PaymentAmount) Validate() error {
	if pa.Currency != "CHF" && pa.Currency != "EUR" {
		return fmt.Errorf("Currency must be CHF or EUR: %v", pa.Currency)
	}
	if pa.Amount < 0.0 {
		return fmt.Errorf("Amount cannot be negative: %v", pa.Amount)
	}
	if len(fmt.Sprintf("%.2f", pa.Amount)) > 12 {
		return fmt.Errorf("Amount too large: %v", pa.Amount)
	}
	return nil
}

// Validate validates a payment reference.
func (r PaymentReference) Validate() error {
	switch v := r.Number.(type) {
	case nil, *structref.CreditorReference, *structref.ReferenceNumber:
		return nil
	default:
		return fmt.Errorf("Unknown reference type: %T", v)
	}
}

// Validate validates additional payment information.
func (pi PaymentInformation) Validate() error {
	if err := ValidateCharacterSet(pi.UnstructuredMessage); err != nil {
		return err
	}
	if err := pi.StructuredMessage.Validate(); err != nil {
		return err
	}
	if s := pi.UnstructuredMessage + pi.StructuredMessage.ToString(); len(s) > 140 {
		return fmt.Errorf("Maximum combined length is 140: %v", s)
	}
	return nil
}

// Validate validates alternative payment procedures.
func (vec AlternativeProcedures) Validate() error {
	if len(vec) > 2 {
		return fmt.Errorf("Maximum two alternate payment schemes allowed: %v", vec)
	}
	for _, ap := range vec {
		if err := ValidateCharacterSet(ap.Label); err != nil {
			return err
		}
		if err := ValidateCharacterSet(ap.Procedure); err != nil {
			return err
		}
		if ap.Label == "" {
			return fmt.Errorf("No label specified: %v", ap)
		}
		if ap.Procedure == "" {
			return fmt.Errorf("No procedure specified: %v", ap)
		}
		if len(ap.Procedure) > 100 {
			return fmt.Errorf("Maximum field length is 100 characters: %v", ap)
		}
	}
	return nil
}

var countryCodes = map[string]bool{
	"AD": true,
	"AE": true,
	"AF": true,
	"AG": true,
	"AI": true,
	"AL": true,
	"AM": true,
	"AO": true,
	"AQ": true,
	"AR": true,
	"AS": true,
	"AT": true,
	"AU": true,
	"AW": true,
	"AX": true,
	"AZ": true,
	"BA": true,
	"BB": true,
	"BD": true,
	"BE": true,
	"BF": true,
	"BG": true,
	"BH": true,
	"BI": true,
	"BJ": true,
	"BL": true,
	"BM": true,
	"BN": true,
	"BO": true,
	"BQ": true,
	"BR": true,
	"BS": true,
	"BT": true,
	"BV": true,
	"BW": true,
	"BY": true,
	"BZ": true,
	"CA": true,
	"CC": true,
	"CD": true,
	"CF": true,
	"CG": true,
	"CH": true,
	"CI": true,
	"CK": true,
	"CL": true,
	"CM": true,
	"CN": true,
	"CO": true,
	"CR": true,
	"CU": true,
	"CV": true,
	"CW": true,
	"CX": true,
	"CY": true,
	"CZ": true,
	"DE": true,
	"DJ": true,
	"DK": true,
	"DM": true,
	"DO": true,
	"DZ": true,
	"EC": true,
	"EE": true,
	"EG": true,
	"EH": true,
	"ER": true,
	"ES": true,
	"ET": true,
	"FI": true,
	"FJ": true,
	"FK": true,
	"FM": true,
	"FO": true,
	"FR": true,
	"GA": true,
	"GB": true,
	"GD": true,
	"GE": true,
	"GF": true,
	"GG": true,
	"GH": true,
	"GI": true,
	"GL": true,
	"GM": true,
	"GN": true,
	"GP": true,
	"GQ": true,
	"GR": true,
	"GS": true,
	"GT": true,
	"GU": true,
	"GW": true,
	"GY": true,
	"HK": true,
	"HM": true,
	"HN": true,
	"HR": true,
	"HT": true,
	"HU": true,
	"ID": true,
	"IE": true,
	"IL": true,
	"IM": true,
	"IN": true,
	"IO": true,
	"IQ": true,
	"IR": true,
	"IS": true,
	"IT": true,
	"JE": true,
	"JM": true,
	"JO": true,
	"JP": true,
	"KE": true,
	"KG": true,
	"KH": true,
	"KI": true,
	"KM": true,
	"KN": true,
	"KP": true,
	"KR": true,
	"KW": true,
	"KY": true,
	"KZ": true,
	"LA": true,
	"LB": true,
	"LC": true,
	"LI": true,
	"LK": true,
	"LR": true,
	"LS": true,
	"LT": true,
	"LU": true,
	"LV": true,
	"LY": true,
	"MA": true,
	"MC": true,
	"MD": true,
	"ME": true,
	"MF": true,
	"MG": true,
	"MH": true,
	"MK": true,
	"ML": true,
	"MM": true,
	"MN": true,
	"MO": true,
	"MP": true,
	"MQ": true,
	"MR": true,
	"MS": true,
	"MT": true,
	"MU": true,
	"MV": true,
	"MW": true,
	"MX": true,
	"MY": true,
	"MZ": true,
	"NA": true,
	"NC": true,
	"NE": true,
	"NF": true,
	"NG": true,
	"NI": true,
	"NL": true,
	"NO": true,
	"NP": true,
	"NR": true,
	"NU": true,
	"NZ": true,
	"OM": true,
	"PA": true,
	"PE": true,
	"PF": true,
	"PG": true,
	"PH": true,
	"PK": true,
	"PL": true,
	"PM": true,
	"PN": true,
	"PR": true,
	"PS": true,
	"PT": true,
	"PW": true,
	"PY": true,
	"QA": true,
	"RE": true,
	"RO": true,
	"RS": true,
	"RU": true,
	"RW": true,
	"SA": true,
	"SB": true,
	"SC": true,
	"SD": true,
	"SE": true,
	"SG": true,
	"SH": true,
	"SI": true,
	"SJ": true,
	"SK": true,
	"SL": true,
	"SM": true,
	"SN": true,
	"SO": true,
	"SR": true,
	"SS": true,
	"ST": true,
	"SV": true,
	"SX": true,
	"SY": true,
	"SZ": true,
	"TC": true,
	"TD": true,
	"TF": true,
	"TG": true,
	"TH": true,
	"TJ": true,
	"TK": true,
	"TL": true,
	"TM": true,
	"TN": true,
	"TO": true,
	"TR": true,
	"TT": true,
	"TV": true,
	"TW": true,
	"TZ": true,
	"UA": true,
	"UG": true,
	"UM": true,
	"US": true,
	"UY": true,
	"UZ": true,
	"VA": true,
	"VC": true,
	"VE": true,
	"VG": true,
	"VI": true,
	"VN": true,
	"VU": true,
	"WF": true,
	"WS": true,
	"YE": true,
	"YT": true,
	"ZA": true,
	"ZM": true,
	"ZW": true,
}
