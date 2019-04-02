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
	"github.com/krepost/structref"
	"time"
)

// Examples from the Swiss QR Standard.
var (
	examplePayload1 = Payload{
		Account: NewIBANOrDie("CH5800791123000889012"),
		Creditor: Entity{
			Name: "Robert Schneider AG",
			Address: StructuredAddress{
				StreetName:     "Rue du Lac",
				BuildingNumber: "1268",
				PostCode:       "2501",
				TownName:       "Biel",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: PaymentAmount{
			Amount:   3949.75,
			Currency: CHF,
		},
		UltimateDebtor: Entity{
			Name: "Pia Rutschmann",
			Address: StructuredAddress{
				StreetName:     "Marktgasse",
				BuildingNumber: "28",
				PostCode:       "9400",
				TownName:       "Rorschach",
			},
			CountryCode: "CH",
		},
		AdditionalInformation: PaymentInformation{
			UnstructuredMessage: "Rechnung Nr. 3139 für Gartenarbeiten und Entsorgung Schnittmaterial",
		},
	}

	examplePayload2 = Payload{
		Account: NewIBANOrDie("CH4431999123000889012"),
		Creditor: Entity{
			Name: "Robert Schneider AG",
			Address: StructuredAddress{
				StreetName:     "Rue du Lac",
				BuildingNumber: "1268",
				PostCode:       "2501",
				TownName:       "Biel",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: PaymentAmount{
			Amount:   1949.75,
			Currency: CHF,
		},
		UltimateDebtor: Entity{
			Name: "Pia-Maria Rutschmann-Schnyder",
			Address: StructuredAddress{
				StreetName:     "Grosse Marktgasse",
				BuildingNumber: "28",
				PostCode:       "9400",
				TownName:       "Rorschach",
			},
			CountryCode: "CH",
		},
		Reference: PaymentReference{
			Number: structref.NewReferenceNumberOrDie("210000000003139471430009017"),
		},
		AdditionalInformation: PaymentInformation{
			UnstructuredMessage: "Auftrag vom 18.06.2020",
			StructuredMessage: BillInformation{
				InvoiceNumber:     "10201409",
				InvoiceDate:       OneDate(2019, time.May, 12),
				CustomerReference: "140.000-53",
				VATNumber:         "106017086",
				VATDates:          OneDate(2018, time.May, 8),
				VATRates: TaxRates{
					TaxRate{RatePercent: 7.7},
				},
				Conditions: PaymentConditions{
					PaymentCondition{DiscountPercent: 2, NumberOfDays: 10},
					PaymentCondition{DiscountPercent: 0, NumberOfDays: 30},
				},
			},
		},
		AlternativeProcedureParameters: []AlternativeProcedure{
			AlternativeProcedure{
				Label:     "Name AV1",
				Procedure: "UV;UltraPay005;12345",
			},
			AlternativeProcedure{
				Label:     "Name AV2",
				Procedure: "XY;XYService;54321",
			},
		},
	}

	examplePayload3 = Payload{
		Account: NewIBANOrDie("CH3709000000304442225"),
		Creditor: Entity{
			Name: "Salvation Army Foundation Switzerland",
			Address: StructuredAddress{
				PostCode: "3000",
				TownName: "Bern",
			},
			CountryCode: "CH",
		},
		CurrencyAmount: PaymentAmount{Currency: CHF},
		AdditionalInformation: PaymentInformation{
			UnstructuredMessage: "Donation to the Winterfest Campaign",
		},
	}
)
