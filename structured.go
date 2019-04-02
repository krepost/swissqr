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
	"regexp"
	"strings"
	"time"
)

// BillInformation contains structured bill information that
// can be added to the Swiss QR invoice. All fields are optional.
type BillInformation struct {
	// InvoiceNumber is free text.
	InvoiceNumber string

	// InvoiceDate contains one date.
	InvoiceDate dates

	// CustomerReference is free text.
	CustomerReference string

	// VATNumber contains the UID with “CHE” prefix, without separators,
	// and without the  MWST/TVA/IVA/VAT suffix.
	VATNumber string

	// VATDates contains either the date of service
	// or the start and end date of service.
	VATDates dates

	// VATRates contains either the VAT rate for the invoice or a list
	// of rates and amounts.
	VATRates TaxRates

	// VATImportTaxRates contains a list of rates and amounts applied
	// during import.
	VATImportTaxRates TaxRates

	// Conditions contains a list of payment conditions.
	Conditions PaymentConditions
}

// dates is an internal struct that encodes either one date
// or a start date with and end date.
type dates struct {
	Date time.Time
	End  time.Time
}

// OneDate is a helper function creating a date.
func OneDate(year int, month time.Month, day int) dates {
	return dates{Date: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// StartAndEndDate is a helper function creating a date interval.
func StartAndEndDate(
	startYear int, startMonth time.Month, startDay int,
	endYear int, endMonth time.Month, endDay int) dates {
	return dates{
		Date: time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, time.UTC),
		End:  time.Date(endYear, endMonth, endDay, 0, 0, 0, 0, time.UTC),
	}
}

// TaxRate represents a tax rate with an optional amount. For an amount
// applicable to the entire invoice amount, set Amount to zero.
type TaxRate struct {
	RatePercent float64
	Amount      float64
}

// TaxRates is a list of tax rates.
type TaxRates []TaxRate

// PaymentCondition represents a discount applied if the invoice
// is paid within a specified number of days. For a payment condition
// of the form “payable within n days”, set DiscountPercent to 0
// and NumberOfDays to n.
type PaymentCondition struct {
	DiscountPercent float64
	NumberOfDays    int
}

// PaymentConditions is a list of payment conditions.
type PaymentConditions []PaymentCondition

// Validate valides a given BillInformation.
func (bi BillInformation) Validate() error {
	if err := ValidateCharacterSet(bi.InvoiceNumber); err != nil {
		return err
	}
	if !bi.InvoiceDate.End.IsZero() {
		return fmt.Errorf("Invoice date may not have an end date: %v", bi.InvoiceDate.End)
	}
	if err := ValidateCharacterSet(bi.CustomerReference); err != nil {
		return err
	}
	if match, _ := regexp.MatchString("^[0-9]*$", bi.VATNumber); !match {
		return fmt.Errorf("VAT number may only contain digits 0-9: %v", bi.VATNumber)
	}
	if !bi.VATDates.Date.IsZero() {
		if !bi.VATDates.End.IsZero() {
			if !bi.VATDates.End.After(bi.VATDates.Date) {
				return fmt.Errorf("End date must come after start date: %v versus %v",
					bi.VATDates.Date, bi.VATDates.End)
			}
		}
	}
	for _, rate := range bi.VATRates {
		if rate.Amount < 0.0 {
			return fmt.Errorf("VAT amount may not be negative: %v", rate)
		}
		if rate.RatePercent < 0.0 {
			return fmt.Errorf("VAT tax rate may not be negative: %v", rate)
		}
	}
	for _, rate := range bi.VATImportTaxRates {
		if rate.Amount < 0.0 {
			return fmt.Errorf("VAT amount may not be negative: %v", rate)
		}
		if rate.RatePercent < 0.0 {
			return fmt.Errorf("VAT tax rate may not be negative: %v", rate)
		}
	}
	for _, condition := range bi.Conditions {
		if condition.DiscountPercent < 0.0 {
			return fmt.Errorf("Discount may not be negative: %v", condition)
		}
		if condition.NumberOfDays < 0 {
			return fmt.Errorf("Number of days may not be negative: %v", condition)
		}
	}
	return nil
}

// ToString converts a given BillInformation to a string that can be added
// to a Swiss QR invoice. It is assumed that the parameters are valid.
func (bi BillInformation) ToString() string {
	result := ""
	if bi.InvoiceNumber != "" {
		result = result + "/10/" + bi.InvoiceNumber
	}
	if s := bi.InvoiceDate.ToString(); s != "" {
		result = result + "/11/" + s
	}
	if bi.CustomerReference != "" {
		result = result + "/20/" + bi.CustomerReference
	}
	if bi.VATNumber != "" {
		result = result + "/30/" + bi.VATNumber
	}
	if s := bi.VATDates.ToString(); s != "" {
		result = result + "/31/" + s
	}
	if s := bi.VATRates.ToString(); s != "" {
		result = result + "/32/" + s
	}
	if s := bi.VATImportTaxRates.ToString(); s != "" {
		result = result + "/33/" + s
	}
	if s := bi.Conditions.ToString(); s != "" {
		result = result + "/40/" + s
	}
	if result != "" {
		result = "//S1" + result
	}
	return result
}

// ToString converts a date or a date interval to a string
// that can be added to a structured bill information.
func (d dates) ToString() string {
	if d.Date.IsZero() {
		return ""
	}
	if d.End.IsZero() {
		return d.Date.Format("060102")
	}
	return d.Date.Format("060102") + d.End.Format("060102")
}

// ToString converts a list of tax rates to a string
// that can be added to a structured bill information.
func (t TaxRates) ToString() string {
	fields := []string{}
	for _, rate := range t {
		if rate.Amount > 0 {
			fields = append(fields, fmt.Sprintf("%g:%g", rate.RatePercent, rate.Amount))
		} else {
			fields = append(fields, fmt.Sprintf("%g", rate.RatePercent))
		}
	}
	return strings.Join(fields, ";")
}

// ToString converts a list of payment conditions to a string
// that can be added to a structured bill information.
func (c PaymentConditions) ToString() string {
	fields := []string{}
	for _, cond := range c {
		fields = append(fields, fmt.Sprintf("%g:%d", cond.DiscountPercent, cond.NumberOfDays))
	}
	return strings.Join(fields, ";")
}
