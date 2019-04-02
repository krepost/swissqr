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
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"github.com/boombuler/barcode"
	barcode_qr "github.com/boombuler/barcode/qr"
)

// CreateQR creates a QR code image from the given payload data. The image is
// 1086×1086 pixels, which amounts to 46×46 mm at 600 dpi, and includes the
// Swiss cross as required by the Swiss payment QR code standard.
func CreateQR(data Payload) (image.Image, error) {
	var buffer bytes.Buffer
	if err := data.Serialize(&buffer); err != nil {
		return nil, err
	}
	qrCode, err := barcode_qr.Encode(buffer.String(), barcode_qr.M, barcode_qr.Unicode)
	if err != nil {
		return nil, err
	}
	qrCode, err = barcode.Scale(qrCode, 1086, 1086) // 46×46 mm at 600 dpi.
	if err != nil {
		return nil, err
	}
	swissQr := image.NewGray16(image.Rect(0, 0, 1086, 1086))
	draw.Draw(swissQr, swissQr.Bounds(), qrCode, image.ZP, draw.Src)
	// Draw 166×166 pixels Swiss cross at the center of the QR code. This
	// produces the symbol which is published at www.paymentstandards.ch.
	swissCross := []struct {
		rect image.Rectangle
		img  image.Image
	}{
		{image.Rect(460, 460, 626, 626), &image.Uniform{color.White}},
		{image.Rect(472, 472, 614, 614), &image.Uniform{color.Black}},
		{image.Rect(496, 526, 590, 554), &image.Uniform{color.White}},
		{image.Rect(528, 494, 558, 586), &image.Uniform{color.White}},
	}
	for _, elem := range swissCross {
		draw.Draw(swissQr, elem.rect, elem.img, image.ZP, draw.Src)
	}
	return swissQr, nil
}
