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
	"math"

	"github.com/krepost/gopdf/pdf"
)

// DrawInvoice draws a standard Swiss QR Invoice on the given canvas, starting
// at the current position as the lower left corner of the invoice. The invoice
// is localized to the given language. The size of the invoice is “DIN A6/5
// Querformat”, i.e., 210 mm wide and 105 mm high. It is the responsibility
// of the caller to make sure that the invoice area in the PDF is clear.
func DrawInvoice(canvas *pdf.Canvas, data Payload, language string) error {
	canvas.Push()
	defer canvas.Pop()
	// drawSupportLines(canvas) // Only for debugging.
	invoice, err := setupForNewInvoice(canvas, data, language)
	if err != nil {
		return err
	}
	if err = invoice.drawReceiptPart(); err != nil {
		return err
	}
	if err = invoice.drawPaymentPart(); err != nil {
		return err
	}
	return nil
}

// DrawInvoiceWithBorder draws a standard Swiss QR Invoice on the given canvas,
// starting at the current position as the lower left corner of the invoice.
// The invoice is localized to the given language. The size of the invoice is
// “DIN A6/5 Querformat”, i.e., 210 mm wide and 105 mm high. A border is drawn
// around the invoice and a text indicating that the invoice part is to be
// removed from the rest of the document is printed above the border. It is
// the responsibility of the caller to make sure that the invoice area in the
// PDF is clear.
func DrawInvoiceWithBorder(canvas *pdf.Canvas, data Payload, language string) error {
	canvas.Push()
	defer canvas.Pop()
	// drawSupportLines(canvas) // Only for debugging.
	invoice, err := setupForNewInvoice(canvas, data, language)
	if err != nil {
		return err
	}
	if err = invoice.drawReceiptPart(); err != nil {
		return err
	}
	if err = invoice.drawPaymentPart(); err != nil {
		return err
	}
	if err = invoice.drawBorderWithText(); err != nil {
		return err
	}
	return nil
}

// DrawInvoiceWithScissors draws a standard Swiss QR Invoice on the given
// canvas, starting at the current position as the lower left corner of the
// invoice. The invoice is localized to the given language. The size of the
// invoice is “DIN A6/5 Querformat”, i.e., 210 mm wide and 105 mm high. A line
// separating the receipt part and the payment part is drawn, and a scissor
// symbol indicating that the receipt part is to be removed is added next to
// the line. It is the responsibility of the caller to make sure that the
// invoice area in the PDF is clear.
func DrawInvoiceWithScissors(canvas *pdf.Canvas, data Payload, language string) error {
	canvas.Push()
	defer canvas.Pop()
	// drawSupportLines(canvas) // Only for debugging.
	invoice, err := setupForNewInvoice(canvas, data, language)
	if err != nil {
		return err
	}
	if err = invoice.drawReceiptPart(); err != nil {
		return err
	}
	if err = invoice.drawPaymentPart(); err != nil {
		return err
	}
	if err = invoice.drawSeparatorWithScissors(); err != nil {
		return err
	}
	return nil
}

// drawBorderWithText draws a solid black border on top of the QR invoice as
// well as between the receipt part and the payment part. A text indicating
// that the payment part should be detached from the rest of the paper is also
// drawn. It is assumed that the current point is at the lower left corner of
// the invoice area.
func (i *pdfInvoice) drawBorderWithText() error {
	i.canvas.Push()
	defer i.canvas.Pop()
	path := new(pdf.Path)
	path.Move(pdf.Point{0, 10.5 * pdf.Cm})
	path.Line(pdf.Point{21.0 * pdf.Cm, 10.5 * pdf.Cm})
	path.Move(pdf.Point{6.2 * pdf.Cm, 0})
	path.Line(pdf.Point{6.2 * pdf.Cm, 10.5 * pdf.Cm})
	i.canvas.SetStrokeColor(0, 0, 0)
	i.canvas.SetLineWidth(1.0)
	i.canvas.Stroke(path)
	if sep, err := BorderText(i.language); err != nil {
		return err
	} else {
		text := new(pdf.Text)
		text.UseFont(i.textFont, 6, 7)
		text.Text(sep)
		i.canvas.Translate(10.5*pdf.Cm-text.X()/2.0, 10.5*pdf.Cm+3)
		i.canvas.DrawText(text)
	}
	return nil
}

// drawSeparatorWithScissors draws a solid black border between the receipt
// part and the payment part of the QR invoice. A scissors symbol is drawn next
// to the separator line. It is assumed that the current point is at the lower
// left corner of the invoice area.
func (i *pdfInvoice) drawSeparatorWithScissors() error {
	i.canvas.Push()
	defer i.canvas.Pop()
	path := new(pdf.Path)
	path.Move(pdf.Point{6.2 * pdf.Cm, 0})
	path.Line(pdf.Point{6.2 * pdf.Cm, 10.5 * pdf.Cm})
	i.canvas.SetStrokeColor(0, 0, 0)
	i.canvas.SetLineWidth(1.0)
	i.canvas.Stroke(path)
	doc := i.canvas.Document()
	dingbats, err := doc.AddFont(pdf.ZapfDingbats, pdf.StandardEncoding)
	if err != nil {
		return err
	}
	rotated := new(pdf.Text)
	rotated.UseFont(dingbats, 20, 25)
	rotated.Text("✂")
	i.canvas.Rotate(math.Pi / 2.0)
	i.canvas.Translate(5.0*pdf.Cm, -6.2*pdf.Cm)
	i.canvas.DrawText(rotated)
	return nil
}

// drawSupportLines draws a grid on canvas to help positioning elements.
// This grid follows the layout given in “Style Guide QR-Rechnung”.
func drawSupportLines(canvas *pdf.Canvas) {
	canvas.Push()
	defer canvas.Pop()
	path := new(pdf.Path)

	path.Rectangle(pdf.Rectangle{
		pdf.Point{0.5 * pdf.Cm, 9.3 * pdf.Cm},
		pdf.Point{5.7 * pdf.Cm, 10.0 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{0.5 * pdf.Cm, 3.7 * pdf.Cm},
		pdf.Point{5.7 * pdf.Cm, 9.3 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{0.5 * pdf.Cm, 2.3 * pdf.Cm},
		pdf.Point{5.7 * pdf.Cm, 3.7 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{0.5 * pdf.Cm, 0.5 * pdf.Cm},
		pdf.Point{5.7 * pdf.Cm, 2.3 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{6.7 * pdf.Cm, 9.3 * pdf.Cm},
		pdf.Point{11.8 * pdf.Cm, 10.0 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{6.7 * pdf.Cm, 4.2 * pdf.Cm},
		pdf.Point{11.3 * pdf.Cm, 8.8 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{6.7 * pdf.Cm, 1.5 * pdf.Cm},
		pdf.Point{11.8 * pdf.Cm, 3.7 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{11.8 * pdf.Cm, 1.5 * pdf.Cm},
		pdf.Point{20.5 * pdf.Cm, 10.0 * pdf.Cm}})

	path.Rectangle(pdf.Rectangle{
		pdf.Point{6.7 * pdf.Cm, 0.5 * pdf.Cm},
		pdf.Point{20.5 * pdf.Cm, 1.5 * pdf.Cm}})

	canvas.SetColor(0.8, 0.8, 0.8)       // Light grey.
	canvas.SetStrokeColor(0.8, 0.8, 0.8) // Light grey.
	canvas.Stroke(path)
}

type pdfInvoice struct {
	canvas    *pdf.Canvas
	titleFont *pdf.Font
	textFont  *pdf.Font
	data      Payload
	language  string
}

type layoutOptions struct {
	headerSize pdf.Unit
	textSize   pdf.Unit
	leading    pdf.Unit
	topLeft    pdf.Point
	maxHeight  pdf.Unit
	maxWidth   pdf.Unit
	boxSize    pdf.Point
}

func setupForNewInvoice(canvas *pdf.Canvas, data Payload, language string) (*pdfInvoice, error) {
	canvas.SetColor(0, 0, 0)
	canvas.SetStrokeColor(0, 0, 0)
	canvas.SetLineWidth(0.75)
	invoice := &pdfInvoice{
		canvas:   canvas,
		data:     data,
		language: language,
	}
	doc := canvas.Document()
	if f, err := doc.AddFont(pdf.Helvetica, pdf.WinAnsiEncoding); err != nil {
		return nil, err
	} else {
		invoice.textFont = f
	}
	if f, err := doc.AddFont(pdf.HelveticaBold, pdf.WinAnsiEncoding); err != nil {
		return nil, err
	} else {
		invoice.titleFont = f
	}
	return invoice, nil
}

func (i *pdfInvoice) drawReceiptPart() error {
	if title, err := TitleSection(i.data, i.language); err != nil {
		return err
	} else {
		i.canvas.Push()
		i.canvas.Translate(0.5*pdf.Cm, 10.0*pdf.Cm-11) // Font size 11.
		text := new(pdf.Text)
		text.UseFont(i.titleFont, 11, 13)
		text.Text(title.Receipt)
		i.canvas.DrawText(text)
		i.canvas.Pop()
	}

	if amt, err := AmountSection(i.data, i.language); err != nil {
		return err
	} else {
		i.drawAmount(amt, layoutOptions{
			headerSize: 6,
			textSize:   8,
			leading:    9,
			topLeft:    pdf.Point{0.5 * pdf.Cm, 3.7 * pdf.Cm},
			maxHeight:  1.4 * pdf.Cm,
			maxWidth:   5.2 * pdf.Cm,
			boxSize:    pdf.Point{3.0 * pdf.Cm, 1.0 * pdf.Cm},
		})
	}

	if info, err := InformationSection(i.data, i.language,
		5.2*28.35/8.0, // 5.2cm × 28.35 pt/cm ÷ 8pt font size.
		receiptPartInformation); err != nil {
		return err
	} else {
		err = i.drawParagraphs(info, layoutOptions{
			headerSize: 6,
			textSize:   8,
			leading:    9,
			topLeft:    pdf.Point{0.5 * pdf.Cm, 9.3 * pdf.Cm},
			maxHeight:  5.6 * pdf.Cm,
			boxSize:    pdf.Point{5.2 * pdf.Cm, 2.0 * pdf.Cm},
		})
		if err != nil {
			return err
		}
	}

	i.canvas.Push()
	text := new(pdf.Text)
	text.UseFont(i.titleFont, 6, 9)
	text.Text(headings[acceptancePoint][i.language])
	i.canvas.Translate(5.7*pdf.Cm-text.X(), 2.3*pdf.Cm-6)
	i.canvas.DrawText(text)
	i.canvas.Pop()

	return nil
}

func (i *pdfInvoice) drawPaymentPart() error {
	if title, err := TitleSection(i.data, i.language); err != nil {
		return err
	} else {
		i.canvas.Push()
		i.canvas.Translate(6.7*pdf.Cm, 10.0*pdf.Cm-11) // Font size 11.
		text := new(pdf.Text)
		text.UseFont(i.titleFont, 11, 13)
		text.Text(title.PaymentPart)
		i.canvas.DrawText(text)
		i.canvas.Pop()
	}

	if amt, err := AmountSection(i.data, i.language); err != nil {
		return err
	} else {
		i.drawAmount(amt, layoutOptions{
			headerSize: 8,
			textSize:   10,
			leading:    11,
			topLeft:    pdf.Point{6.7 * pdf.Cm, 3.7 * pdf.Cm},
			maxHeight:  2.2 * pdf.Cm,
			maxWidth:   5.1 * pdf.Cm,
			boxSize:    pdf.Point{4.0 * pdf.Cm, 1.5 * pdf.Cm},
		})
	}

	if qrImage, err := CreateQR(i.data); err != nil {
		return err
	} else {
		// 46×46 mm image; at least 5 mm margin.
		// Payment part starts at 61.5 mm indent.
		i.canvas.DrawImage(qrImage, pdf.Rectangle{
			pdf.Point{6.7 * pdf.Cm, 4.3 * pdf.Cm},
			pdf.Point{11.3 * pdf.Cm, 8.9 * pdf.Cm}})
	}

	if section, err := InformationSection(i.data, i.language,
		8.5*28.35/10.0, // 8.5cm × 28.35 pt/cm ÷ 10pt font size.
		paymentPart); err != nil {
		return err
	} else {
		err = i.drawParagraphs(section, layoutOptions{
			headerSize: 8,
			textSize:   10,
			leading:    11,
			topLeft:    pdf.Point{11.9 * pdf.Cm, 10.0 * pdf.Cm},
			maxHeight:  8.5 * pdf.Cm,
			boxSize:    pdf.Point{6.5 * pdf.Cm, 2.5 * pdf.Cm},
		})
		if err != nil {
			return err
		}
	}

	i.canvas.Push()
	text := new(pdf.Text)
	for _, ap := range i.data.AlternativeProcedureParameters {
		text.UseFont(i.titleFont, 7, 8)
		text.Text(ap.Label + ": ")
		text.UseFont(i.textFont, 7, 8)
		// 13cm × 28.35 pt/cm ÷ 7pt font size is total width.
		remainingWidth := (13.8*28.35 - text.X()) / 7.0
		text.Text(shortenToWidth(ap.Procedure, float64(remainingWidth)))
		text.NextLine()
	}
	i.canvas.Translate(6.7*pdf.Cm, 1.5*pdf.Cm-8)
	i.canvas.DrawText(text)
	i.canvas.Pop()

	return nil
}

// drawAmount draws a payment amount, or an empty box if there is no amount.
func (i *pdfInvoice) drawAmount(amt AmountSectionData, layout layoutOptions) error {
	i.canvas.Push()
	defer i.canvas.Pop()
	i.canvas.Translate(layout.topLeft.X, layout.topLeft.Y-layout.headerSize)
	text := new(pdf.Text)
	text.UseFont(i.titleFont, layout.headerSize, layout.leading)
	text.Text(amt.CurrencyHeading)
	columnSeparation := text.X() + layout.headerSize
	text.NextLineOffset(columnSeparation, 0)
	text.Text(amt.AmountHeading)
	totalHeaderWidth := text.X()
	text.NextLineOffset(-columnSeparation, -layout.leading)
	text.UseFont(i.textFont, layout.textSize, layout.leading)
	text.Text(amt.CurrencyValue)
	if amt.AmountValue != "" {
		text.NextLineOffset(columnSeparation, 0)
		text.Text(amt.AmountValue)
	} else {
		topLeft := pdf.Point{layout.maxWidth - layout.boxSize.X, layout.headerSize}
		if totalHeaderWidth+layout.boxSize.X > layout.maxWidth {
			topLeft.Y = text.Y() + layout.leading - 5
		}
		path := new(pdf.Path)
		drawCorners(path, pdf.Rectangle{
			Min: pdf.Point{topLeft.X, topLeft.Y - layout.boxSize.Y},
			Max: pdf.Point{topLeft.X + layout.boxSize.X, topLeft.Y},
		})
		i.canvas.Stroke(path)
	}
	i.canvas.DrawText(text)
	return nil
}

// drawParagraphs draws a slice of paragraphs following the layout options.
// If a paragraph is empty, a box is drawn instead.
func (i *pdfInvoice) drawParagraphs(section []Paragraph, layout layoutOptions) error {
	i.canvas.Push()
	defer i.canvas.Pop()
	i.canvas.Translate(layout.topLeft.X, layout.topLeft.Y-layout.headerSize)
	text := new(pdf.Text)
	firstLine := true
	for _, s := range section {
		if firstLine {
			firstLine = false
		} else {
			text.NextLine()
			text.NextLineOffset(0, -3)
		}
		text.UseFont(i.titleFont, layout.headerSize, layout.leading)
		text.Text(s.Heading)
		if len(s.Lines) > 0 {
			text.UseFont(i.textFont, layout.textSize, layout.leading)
			for _, line := range s.Lines {
				text.NextLine()
				text.Text(line)
			}
		} else {
			path := new(pdf.Path)
			drawCorners(path, pdf.Rectangle{
				Min: pdf.Point{0, text.Y() - 5 - layout.boxSize.Y},
				Max: pdf.Point{layout.boxSize.X, text.Y() - 5},
			})
			i.canvas.Stroke(path)
			text.NextLineOffset(0, -layout.boxSize.Y-3)
		}
	}

	if -text.Y() > layout.maxHeight { // -text.Y() is the text height.
		return errors.New("Invoice text height too large.")
	}
	i.canvas.DrawText(text)
	return nil
}

// drawCorners draws corner marks around the given box.
func drawCorners(path *pdf.Path, box pdf.Rectangle) {
	// According to the standard, the corner marks should be 3 mm long.
	size := 0.3 * pdf.Cm
	// Lower left corner
	path.Move(pdf.Point{box.Min.X + size, box.Min.Y})
	path.Line(pdf.Point{box.Min.X, box.Min.Y})
	path.Line(pdf.Point{box.Min.X, box.Min.Y + size})
	// Upper left corner
	path.Move(pdf.Point{box.Min.X, box.Max.Y - size})
	path.Line(pdf.Point{box.Min.X, box.Max.Y})
	path.Line(pdf.Point{box.Min.X + size, box.Max.Y})
	// Upper right corner
	path.Move(pdf.Point{box.Max.X - size, box.Max.Y})
	path.Line(pdf.Point{box.Max.X, box.Max.Y})
	path.Line(pdf.Point{box.Max.X, box.Max.Y - size})
	// Lower right corner
	path.Move(pdf.Point{box.Max.X, box.Min.Y + size})
	path.Line(pdf.Point{box.Max.X, box.Min.Y})
	path.Line(pdf.Point{box.Max.X - size, box.Min.Y})
}
