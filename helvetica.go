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

import "strings"

func reflowAtRune(lines []string, maxWidth float64) []string {
	reflowed := []string{}
	for _, line := range lines {
		if stringWidth(line) < maxWidth {
			reflowed = append(reflowed, line)
		} else {
			currentLine := ""
			currentWidth := 0.0
			for _, r := range line {
				w := 1.0 // Default.
				if m, ok := regularWX[r]; ok {
					w = m
				}
				if currentWidth+w > maxWidth {
					reflowed = append(reflowed, currentLine)
					currentLine = ""
					currentWidth = 0.0
				}
				currentWidth = currentWidth + w
				currentLine = currentLine + string(r)
			}
			reflowed = append(reflowed, currentLine)
		}
	}
	return reflowed
}

func reflowAtSpace(lines []string, maxWidth float64) []string {
	reflowed := []string{}
	spaceWidth, _ := regularWX[' ']
	for _, line := range lines {
		if stringWidth(line) < maxWidth {
			reflowed = append(reflowed, line)
		} else {
			currentLine := ""
			currentWidth := 0.0
			for _, word := range strings.Fields(line) {
				w := stringWidth(word)
				if w > maxWidth {
					// Switch to “reflowAtRune” algorithm.
					if currentLine != "" {
						currentWidth = currentWidth + spaceWidth
						currentLine = currentLine + " "
					}
					for _, r := range word {
						w := 1.0 // Default.
						if m, ok := regularWX[r]; ok {
							w = m
						}
						if currentWidth+w > maxWidth {
							reflowed = append(reflowed, currentLine)
							currentLine = ""
							currentWidth = 0.0
						}
						currentWidth = currentWidth + w
						currentLine = currentLine + string(r)
					}
				} else {
					// Check if we should break before this word.
					// If current line is empty, we know that
					// word will fit (since w < maxWidth here).
					if currentLine != "" {
						if currentWidth+spaceWidth+w > maxWidth {
							reflowed = append(reflowed, currentLine)
							currentLine = ""
							currentWidth = 0.0
						} else {
							currentWidth = currentWidth + spaceWidth
							currentLine = currentLine + " "
						}
					}
					currentWidth = currentWidth + w
					currentLine = currentLine + word
				}
			}
			reflowed = append(reflowed, currentLine)
		}
	}
	return reflowed
}

func shortenToWidth(line string, maxWidth float64) string {
	shortened := ""
	suffix := "…"
	suffixWidth := stringWidth(suffix)
	currentWidth := 0.0
	for _, r := range line {
		w := 1.0 // Default.
		if m, ok := regularWX[r]; ok {
			w = m
		}
		if currentWidth+suffixWidth+w > maxWidth {
			return shortened + suffix
		}
		currentWidth = currentWidth + w
		shortened = shortened + string(r)
	}
	return line
}

// Returns the width of s when printed in Helvetica font,
// relative to the point size used.
func stringWidth(s string) float64 {
	multiplier := 0.0
	for _, r := range s {
		if m, ok := regularWX[r]; ok {
			multiplier = multiplier + m
		} else {
			multiplier = multiplier + 1.0
		}
	}
	return multiplier
}

// Glyph widths for Helvetica font, relative to point size.
var regularWX = map[rune]float64{
	'\u0020': 0.278, // space
	'\u0021': 0.278, // exclam
	'\u0022': 0.355, // quotedbl
	'\u0023': 0.556, // numbersign
	'\u0024': 0.556, // dollar
	'\u0025': 0.889, // percent
	'\u0026': 0.667, // ampersand
	'\u0027': 0.191, // quotesingle
	'\u0028': 0.333, // parenleft
	'\u0029': 0.333, // parenright
	'\u002A': 0.389, // asterisk
	'\u002B': 0.584, // plus
	'\u002C': 0.278, // comma
	'\u002D': 0.333, // hyphen
	'\u002E': 0.278, // period
	'\u002F': 0.278, // slash
	'\u0030': 0.556, // zero
	'\u0031': 0.556, // one
	'\u0032': 0.556, // two
	'\u0033': 0.556, // three
	'\u0034': 0.556, // four
	'\u0035': 0.556, // five
	'\u0036': 0.556, // six
	'\u0037': 0.556, // seven
	'\u0038': 0.556, // eight
	'\u0039': 0.556, // nine
	'\u003A': 0.278, // colon
	'\u003B': 0.278, // semicolon
	'\u003C': 0.584, // less
	'\u003D': 0.584, // equal
	'\u003E': 0.584, // greater
	'\u003F': 0.556, // question
	'\u0040': 1.015, // at
	'\u0041': 0.667, // A
	'\u0042': 0.667, // B
	'\u0043': 0.722, // C
	'\u0044': 0.722, // D
	'\u0045': 0.667, // E
	'\u0046': 0.611, // F
	'\u0047': 0.778, // G
	'\u0048': 0.722, // H
	'\u0049': 0.278, // I
	'\u004A': 0.500, // J
	'\u004B': 0.667, // K
	'\u004C': 0.556, // L
	'\u004D': 0.833, // M
	'\u004E': 0.722, // N
	'\u004F': 0.778, // O
	'\u0050': 0.667, // P
	'\u0051': 0.778, // Q
	'\u0052': 0.722, // R
	'\u0053': 0.667, // S
	'\u0054': 0.611, // T
	'\u0055': 0.722, // U
	'\u0056': 0.667, // V
	'\u0057': 0.944, // W
	'\u0058': 0.667, // X
	'\u0059': 0.667, // Y
	'\u005A': 0.611, // Z
	'\u005B': 0.278, // bracketleft
	'\u005C': 0.278, // backslash
	'\u005D': 0.278, // bracketright
	'\u005E': 0.469, // asciicircum
	'\u005F': 0.556, // underscore
	'\u0060': 0.333, // grave
	'\u0061': 0.556, // a
	'\u0062': 0.556, // b
	'\u0063': 0.500, // c
	'\u0064': 0.556, // d
	'\u0065': 0.556, // e
	'\u0066': 0.278, // f
	'\u0067': 0.556, // g
	'\u0068': 0.556, // h
	'\u0069': 0.222, // i
	'\u006A': 0.222, // j
	'\u006B': 0.500, // k
	'\u006C': 0.222, // l
	'\u006D': 0.833, // m
	'\u006E': 0.556, // n
	'\u006F': 0.556, // o
	'\u0070': 0.556, // p
	'\u0071': 0.556, // q
	'\u0072': 0.333, // r
	'\u0073': 0.500, // s
	'\u0074': 0.278, // t
	'\u0075': 0.556, // u
	'\u0076': 0.500, // v
	'\u0077': 0.722, // w
	'\u0078': 0.500, // x
	'\u0079': 0.500, // y
	'\u007A': 0.500, // z
	'\u007B': 0.334, // braceleft
	'\u007C': 0.260, // bar
	'\u007D': 0.334, // braceright
	'\u007E': 0.584, // asciitilde
	'\u00A1': 0.333, // exclamdown
	'\u00A2': 0.556, // cent
	'\u00A3': 0.556, // sterling
	'\u00A4': 0.556, // currency
	'\u00A5': 0.556, // yen
	'\u00A6': 0.260, // brokenbar
	'\u00A7': 0.556, // section
	'\u00A8': 0.333, // dieresis
	'\u00A9': 0.737, // copyright
	'\u00AA': 0.370, // ordfeminine
	'\u00AB': 0.556, // guillemotleft
	'\u00AC': 0.584, // logicalnot
	'\u00AE': 0.737, // registered
	'\u00AF': 0.333, // macron
	'\u00B0': 0.400, // degree
	'\u00B1': 0.584, // plusminus
	'\u00B2': 0.333, // twosuperior
	'\u00B3': 0.333, // threesuperior
	'\u00B4': 0.333, // acute
	'\u00B5': 0.556, // mu
	'\u00B6': 0.537, // paragraph
	'\u00B7': 0.278, // periodcentered
	'\u00B8': 0.333, // cedilla
	'\u00B9': 0.333, // onesuperior
	'\u00BA': 0.365, // ordmasculine
	'\u00BB': 0.556, // guillemotright
	'\u00BC': 0.834, // onequarter
	'\u00BD': 0.834, // onehalf
	'\u00BE': 0.834, // threequarters
	'\u00BF': 0.611, // questiondown
	'\u00C0': 0.667, // Agrave
	'\u00C1': 0.667, // Aacute
	'\u00C2': 0.667, // Acircumflex
	'\u00C3': 0.667, // Atilde
	'\u00C4': 0.667, // Adieresis
	'\u00C5': 0.667, // Aring
	'\u00C6': 1.000, // AE
	'\u00C7': 0.722, // Ccedilla
	'\u00C8': 0.667, // Egrave
	'\u00C9': 0.667, // Eacute
	'\u00CA': 0.667, // Ecircumflex
	'\u00CB': 0.667, // Edieresis
	'\u00CC': 0.278, // Igrave
	'\u00CD': 0.278, // Iacute
	'\u00CE': 0.278, // Icircumflex
	'\u00CF': 0.278, // Idieresis
	'\u00D0': 0.722, // Eth
	'\u00D1': 0.722, // Ntilde
	'\u00D2': 0.778, // Ograve
	'\u00D3': 0.778, // Oacute
	'\u00D4': 0.778, // Ocircumflex
	'\u00D5': 0.778, // Otilde
	'\u00D6': 0.778, // Odieresis
	'\u00D7': 0.584, // multiply
	'\u00D8': 0.778, // Oslash
	'\u00D9': 0.722, // Ugrave
	'\u00DA': 0.722, // Uacute
	'\u00DB': 0.722, // Ucircumflex
	'\u00DC': 0.722, // Udieresis
	'\u00DD': 0.667, // Yacute
	'\u00DE': 0.667, // Thorn
	'\u00DF': 0.611, // germandbls
	'\u00E0': 0.556, // agrave
	'\u00E1': 0.556, // aacute
	'\u00E2': 0.556, // acircumflex
	'\u00E3': 0.556, // atilde
	'\u00E4': 0.556, // adieresis
	'\u00E5': 0.556, // aring
	'\u00E6': 0.889, // ae
	'\u00E7': 0.500, // ccedilla
	'\u00E8': 0.556, // egrave
	'\u00E9': 0.556, // eacute
	'\u00EA': 0.556, // ecircumflex
	'\u00EB': 0.556, // edieresis
	'\u00EC': 0.278, // igrave
	'\u00ED': 0.278, // iacute
	'\u00EE': 0.278, // icircumflex
	'\u00EF': 0.278, // idieresis
	'\u00F0': 0.556, // eth
	'\u00F1': 0.556, // ntilde
	'\u00F2': 0.556, // ograve
	'\u00F3': 0.556, // oacute
	'\u00F4': 0.556, // ocircumflex
	'\u00F5': 0.556, // otilde
	'\u00F6': 0.556, // odieresis
	'\u00F7': 0.584, // divide
	'\u00F8': 0.611, // oslash
	'\u00F9': 0.556, // ugrave
	'\u00FA': 0.556, // uacute
	'\u00FB': 0.556, // ucircumflex
	'\u00FC': 0.556, // udieresis
	'\u00FD': 0.500, // yacute
	'\u00FE': 0.556, // thorn
	'\u00FF': 0.500, // ydieresis
	'\u0100': 0.667, // Amacron
	'\u0101': 0.556, // amacron
	'\u0102': 0.667, // Abreve
	'\u0103': 0.556, // abreve
	'\u0104': 0.667, // Aogonek
	'\u0105': 0.556, // aogonek
	'\u0106': 0.722, // Cacute
	'\u0107': 0.500, // cacute
	'\u010C': 0.722, // Ccaron
	'\u010D': 0.500, // ccaron
	'\u010E': 0.722, // Dcaron
	'\u010F': 0.643, // dcaron
	'\u0110': 0.722, // Dcroat
	'\u0111': 0.556, // dcroat
	'\u0112': 0.667, // Emacron
	'\u0113': 0.556, // emacron
	'\u0116': 0.667, // Edotaccent
	'\u0117': 0.556, // edotaccent
	'\u0118': 0.667, // Eogonek
	'\u0119': 0.556, // eogonek
	'\u011A': 0.667, // Ecaron
	'\u011B': 0.556, // ecaron
	'\u011E': 0.778, // Gbreve
	'\u011F': 0.556, // gbreve
	'\u0122': 0.778, // Gcommaaccent
	'\u0123': 0.556, // gcommaaccent
	'\u012A': 0.278, // Imacron
	'\u012B': 0.278, // imacron
	'\u012E': 0.278, // Iogonek
	'\u012F': 0.222, // iogonek
	'\u0130': 0.278, // Idotaccent
	'\u0131': 0.278, // dotlessi
	'\u0136': 0.667, // Kcommaaccent
	'\u0137': 0.500, // kcommaaccent
	'\u0139': 0.556, // Lacute
	'\u013A': 0.222, // lacute
	'\u013B': 0.556, // Lcommaaccent
	'\u013C': 0.222, // lcommaaccent
	'\u013D': 0.556, // Lcaron
	'\u013E': 0.299, // lcaron
	'\u0141': 0.556, // Lslash
	'\u0142': 0.222, // lslash
	'\u0143': 0.722, // Nacute
	'\u0144': 0.556, // nacute
	'\u0145': 0.722, // Ncommaaccent
	'\u0146': 0.556, // ncommaaccent
	'\u0147': 0.722, // Ncaron
	'\u0148': 0.556, // ncaron
	'\u014C': 0.778, // Omacron
	'\u014D': 0.556, // omacron
	'\u0150': 0.778, // Ohungarumlaut
	'\u0151': 0.556, // ohungarumlaut
	'\u0152': 1.000, // OE
	'\u0153': 0.944, // oe
	'\u0154': 0.722, // Racute
	'\u0155': 0.333, // racute
	'\u0156': 0.722, // Rcommaaccent
	'\u0157': 0.333, // rcommaaccent
	'\u0158': 0.722, // Rcaron
	'\u0159': 0.333, // rcaron
	'\u015A': 0.667, // Sacute
	'\u015B': 0.500, // sacute
	'\u015E': 0.667, // Scedilla
	'\u015F': 0.500, // scedilla
	'\u0160': 0.667, // Scaron
	'\u0161': 0.500, // scaron
	'\u0162': 0.611, // Tcommaaccent
	'\u0163': 0.278, // tcommaaccent
	'\u0164': 0.611, // Tcaron
	'\u0165': 0.317, // tcaron
	'\u016A': 0.722, // Umacron
	'\u016B': 0.556, // umacron
	'\u016E': 0.722, // Uring
	'\u016F': 0.556, // uring
	'\u0170': 0.722, // Uhungarumlaut
	'\u0171': 0.556, // uhungarumlaut
	'\u0172': 0.722, // Uogonek
	'\u0173': 0.556, // uogonek
	'\u0178': 0.667, // Ydieresis
	'\u0179': 0.611, // Zacute
	'\u017A': 0.500, // zacute
	'\u017B': 0.611, // Zdotaccent
	'\u017C': 0.500, // zdotaccent
	'\u017D': 0.611, // Zcaron
	'\u017E': 0.500, // zcaron
	'\u0192': 0.556, // florin
	'\u0218': 0.667, // Scommaaccent
	'\u0219': 0.500, // scommaaccent
	'\u02C6': 0.333, // circumflex
	'\u02C7': 0.333, // caron
	'\u02D8': 0.333, // breve
	'\u02D9': 0.333, // dotaccent
	'\u02DA': 0.333, // ring
	'\u02DB': 0.333, // ogonek
	'\u02DC': 0.333, // tilde
	'\u02DD': 0.333, // hungarumlaut
	'\u2013': 0.556, // endash
	'\u2014': 1.000, // emdash
	'\u2018': 0.222, // quoteleft
	'\u2019': 0.222, // quoteright
	'\u201A': 0.222, // quotesinglbase
	'\u201C': 0.333, // quotedblleft
	'\u201D': 0.333, // quotedblright
	'\u201E': 0.333, // quotedblbase
	'\u2020': 0.556, // dagger
	'\u2021': 0.556, // daggerdbl
	'\u2022': 0.350, // bullet
	'\u2026': 1.000, // ellipsis
	'\u2030': 1.000, // perthousand
	'\u2039': 0.333, // guilsinglleft
	'\u203A': 0.333, // guilsinglright
	'\u2044': 0.167, // fraction
	'\u20AC': 0.556, // Euro
	'\u2122': 1.000, // trademark
	'\u2202': 0.476, // partialdiff
	'\u2206': 0.612, // Delta
	'\u2211': 0.600, // summation
	'\u2212': 0.584, // minus
	'\u221A': 0.453, // radical
	'\u2260': 0.549, // notequal
	'\u2264': 0.549, // lessequal
	'\u2265': 0.549, // greaterequal
	'\u25CA': 0.471, // lozenge
	'\uF6C3': 0.250, // commaaccent
	'\uFB01': 0.500, // fi
	'\uFB02': 0.500, // fl
}

// Glyph widths for Helvetica-Bold font, relative to point size.
var boldWX = map[rune]float64{
	'\u0020': 0.278, // space
	'\u0021': 0.333, // exclam
	'\u0022': 0.474, // quotedbl
	'\u0023': 0.556, // numbersign
	'\u0024': 0.556, // dollar
	'\u0025': 0.889, // percent
	'\u0026': 0.722, // ampersand
	'\u0027': 0.238, // quotesingle
	'\u0028': 0.333, // parenleft
	'\u0029': 0.333, // parenright
	'\u002A': 0.389, // asterisk
	'\u002B': 0.584, // plus
	'\u002C': 0.278, // comma
	'\u002D': 0.333, // hyphen
	'\u002E': 0.278, // period
	'\u002F': 0.278, // slash
	'\u0030': 0.556, // zero
	'\u0031': 0.556, // one
	'\u0032': 0.556, // two
	'\u0033': 0.556, // three
	'\u0034': 0.556, // four
	'\u0035': 0.556, // five
	'\u0036': 0.556, // six
	'\u0037': 0.556, // seven
	'\u0038': 0.556, // eight
	'\u0039': 0.556, // nine
	'\u003A': 0.333, // colon
	'\u003B': 0.333, // semicolon
	'\u003C': 0.584, // less
	'\u003D': 0.584, // equal
	'\u003E': 0.584, // greater
	'\u003F': 0.611, // question
	'\u0040': 0.975, // at
	'\u0041': 0.722, // A
	'\u0042': 0.722, // B
	'\u0043': 0.722, // C
	'\u0044': 0.722, // D
	'\u0045': 0.667, // E
	'\u0046': 0.611, // F
	'\u0047': 0.778, // G
	'\u0048': 0.722, // H
	'\u0049': 0.278, // I
	'\u004A': 0.556, // J
	'\u004B': 0.722, // K
	'\u004C': 0.611, // L
	'\u004D': 0.833, // M
	'\u004E': 0.722, // N
	'\u004F': 0.778, // O
	'\u0050': 0.667, // P
	'\u0051': 0.778, // Q
	'\u0052': 0.722, // R
	'\u0053': 0.667, // S
	'\u0054': 0.611, // T
	'\u0055': 0.722, // U
	'\u0056': 0.667, // V
	'\u0057': 0.944, // W
	'\u0058': 0.667, // X
	'\u0059': 0.667, // Y
	'\u005A': 0.611, // Z
	'\u005B': 0.333, // bracketleft
	'\u005C': 0.278, // backslash
	'\u005D': 0.333, // bracketright
	'\u005E': 0.584, // asciicircum
	'\u005F': 0.556, // underscore
	'\u0060': 0.333, // grave
	'\u0061': 0.556, // a
	'\u0062': 0.611, // b
	'\u0063': 0.556, // c
	'\u0064': 0.611, // d
	'\u0065': 0.556, // e
	'\u0066': 0.333, // f
	'\u0067': 0.611, // g
	'\u0068': 0.611, // h
	'\u0069': 0.278, // i
	'\u006A': 0.278, // j
	'\u006B': 0.556, // k
	'\u006C': 0.278, // l
	'\u006D': 0.889, // m
	'\u006E': 0.611, // n
	'\u006F': 0.611, // o
	'\u0070': 0.611, // p
	'\u0071': 0.611, // q
	'\u0072': 0.389, // r
	'\u0073': 0.556, // s
	'\u0074': 0.333, // t
	'\u0075': 0.611, // u
	'\u0076': 0.556, // v
	'\u0077': 0.778, // w
	'\u0078': 0.556, // x
	'\u0079': 0.556, // y
	'\u007A': 0.500, // z
	'\u007B': 0.389, // braceleft
	'\u007C': 0.280, // bar
	'\u007D': 0.389, // braceright
	'\u007E': 0.584, // asciitilde
	'\u00A1': 0.333, // exclamdown
	'\u00A2': 0.556, // cent
	'\u00A3': 0.556, // sterling
	'\u00A4': 0.556, // currency
	'\u00A5': 0.556, // yen
	'\u00A6': 0.280, // brokenbar
	'\u00A7': 0.556, // section
	'\u00A8': 0.333, // dieresis
	'\u00A9': 0.737, // copyright
	'\u00AA': 0.370, // ordfeminine
	'\u00AB': 0.556, // guillemotleft
	'\u00AC': 0.584, // logicalnot
	'\u00AE': 0.737, // registered
	'\u00AF': 0.333, // macron
	'\u00B0': 0.400, // degree
	'\u00B1': 0.584, // plusminus
	'\u00B2': 0.333, // twosuperior
	'\u00B3': 0.333, // threesuperior
	'\u00B4': 0.333, // acute
	'\u00B5': 0.611, // mu
	'\u00B6': 0.556, // paragraph
	'\u00B7': 0.278, // periodcentered
	'\u00B8': 0.333, // cedilla
	'\u00B9': 0.333, // onesuperior
	'\u00BA': 0.365, // ordmasculine
	'\u00BB': 0.556, // guillemotright
	'\u00BC': 0.834, // onequarter
	'\u00BD': 0.834, // onehalf
	'\u00BE': 0.834, // threequarters
	'\u00BF': 0.611, // questiondown
	'\u00C0': 0.722, // Agrave
	'\u00C1': 0.722, // Aacute
	'\u00C2': 0.722, // Acircumflex
	'\u00C3': 0.722, // Atilde
	'\u00C4': 0.722, // Adieresis
	'\u00C5': 0.722, // Aring
	'\u00C6': 1.000, // AE
	'\u00C7': 0.722, // Ccedilla
	'\u00C8': 0.667, // Egrave
	'\u00C9': 0.667, // Eacute
	'\u00CA': 0.667, // Ecircumflex
	'\u00CB': 0.667, // Edieresis
	'\u00CC': 0.278, // Igrave
	'\u00CD': 0.278, // Iacute
	'\u00CE': 0.278, // Icircumflex
	'\u00CF': 0.278, // Idieresis
	'\u00D0': 0.722, // Eth
	'\u00D1': 0.722, // Ntilde
	'\u00D2': 0.778, // Ograve
	'\u00D3': 0.778, // Oacute
	'\u00D4': 0.778, // Ocircumflex
	'\u00D5': 0.778, // Otilde
	'\u00D6': 0.778, // Odieresis
	'\u00D7': 0.584, // multiply
	'\u00D8': 0.778, // Oslash
	'\u00D9': 0.722, // Ugrave
	'\u00DA': 0.722, // Uacute
	'\u00DB': 0.722, // Ucircumflex
	'\u00DC': 0.722, // Udieresis
	'\u00DD': 0.667, // Yacute
	'\u00DE': 0.667, // Thorn
	'\u00DF': 0.611, // germandbls
	'\u00E0': 0.556, // agrave
	'\u00E1': 0.556, // aacute
	'\u00E2': 0.556, // acircumflex
	'\u00E3': 0.556, // atilde
	'\u00E4': 0.556, // adieresis
	'\u00E5': 0.556, // aring
	'\u00E6': 0.889, // ae
	'\u00E7': 0.556, // ccedilla
	'\u00E8': 0.556, // egrave
	'\u00E9': 0.556, // eacute
	'\u00EA': 0.556, // ecircumflex
	'\u00EB': 0.556, // edieresis
	'\u00EC': 0.278, // igrave
	'\u00ED': 0.278, // iacute
	'\u00EE': 0.278, // icircumflex
	'\u00EF': 0.278, // idieresis
	'\u00F0': 0.611, // eth
	'\u00F1': 0.611, // ntilde
	'\u00F2': 0.611, // ograve
	'\u00F3': 0.611, // oacute
	'\u00F4': 0.611, // ocircumflex
	'\u00F5': 0.611, // otilde
	'\u00F6': 0.611, // odieresis
	'\u00F7': 0.584, // divide
	'\u00F8': 0.611, // oslash
	'\u00F9': 0.611, // ugrave
	'\u00FA': 0.611, // uacute
	'\u00FB': 0.611, // ucircumflex
	'\u00FC': 0.611, // udieresis
	'\u00FD': 0.556, // yacute
	'\u00FE': 0.611, // thorn
	'\u00FF': 0.556, // ydieresis
	'\u0100': 0.722, // Amacron
	'\u0101': 0.556, // amacron
	'\u0102': 0.722, // Abreve
	'\u0103': 0.556, // abreve
	'\u0104': 0.722, // Aogonek
	'\u0105': 0.556, // aogonek
	'\u0106': 0.722, // Cacute
	'\u0107': 0.556, // cacute
	'\u010C': 0.722, // Ccaron
	'\u010D': 0.556, // ccaron
	'\u010E': 0.722, // Dcaron
	'\u010F': 0.743, // dcaron
	'\u0110': 0.722, // Dcroat
	'\u0111': 0.611, // dcroat
	'\u0112': 0.667, // Emacron
	'\u0113': 0.556, // emacron
	'\u0116': 0.667, // Edotaccent
	'\u0117': 0.556, // edotaccent
	'\u0118': 0.667, // Eogonek
	'\u0119': 0.556, // eogonek
	'\u011A': 0.667, // Ecaron
	'\u011B': 0.556, // ecaron
	'\u011E': 0.778, // Gbreve
	'\u011F': 0.611, // gbreve
	'\u0122': 0.778, // Gcommaaccent
	'\u0123': 0.611, // gcommaaccent
	'\u012A': 0.278, // Imacron
	'\u012B': 0.278, // imacron
	'\u012E': 0.278, // Iogonek
	'\u012F': 0.278, // iogonek
	'\u0130': 0.278, // Idotaccent
	'\u0131': 0.278, // dotlessi
	'\u0136': 0.722, // Kcommaaccent
	'\u0137': 0.556, // kcommaaccent
	'\u0139': 0.611, // Lacute
	'\u013A': 0.278, // lacute
	'\u013B': 0.611, // Lcommaaccent
	'\u013C': 0.278, // lcommaaccent
	'\u013D': 0.611, // Lcaron
	'\u013E': 0.400, // lcaron
	'\u0141': 0.611, // Lslash
	'\u0142': 0.278, // lslash
	'\u0143': 0.722, // Nacute
	'\u0144': 0.611, // nacute
	'\u0145': 0.722, // Ncommaaccent
	'\u0146': 0.611, // ncommaaccent
	'\u0147': 0.722, // Ncaron
	'\u0148': 0.611, // ncaron
	'\u014C': 0.778, // Omacron
	'\u014D': 0.611, // omacron
	'\u0150': 0.778, // Ohungarumlaut
	'\u0151': 0.611, // ohungarumlaut
	'\u0152': 1.000, // OE
	'\u0153': 0.944, // oe
	'\u0154': 0.722, // Racute
	'\u0155': 0.389, // racute
	'\u0156': 0.722, // Rcommaaccent
	'\u0157': 0.389, // rcommaaccent
	'\u0158': 0.722, // Rcaron
	'\u0159': 0.389, // rcaron
	'\u015A': 0.667, // Sacute
	'\u015B': 0.556, // sacute
	'\u015E': 0.667, // Scedilla
	'\u015F': 0.556, // scedilla
	'\u0160': 0.667, // Scaron
	'\u0161': 0.556, // scaron
	'\u0162': 0.611, // Tcommaaccent
	'\u0163': 0.333, // tcommaaccent
	'\u0164': 0.611, // Tcaron
	'\u0165': 0.389, // tcaron
	'\u016A': 0.722, // Umacron
	'\u016B': 0.611, // umacron
	'\u016E': 0.722, // Uring
	'\u016F': 0.611, // uring
	'\u0170': 0.722, // Uhungarumlaut
	'\u0171': 0.611, // uhungarumlaut
	'\u0172': 0.722, // Uogonek
	'\u0173': 0.611, // uogonek
	'\u0178': 0.667, // Ydieresis
	'\u0179': 0.611, // Zacute
	'\u017A': 0.500, // zacute
	'\u017B': 0.611, // Zdotaccent
	'\u017C': 0.500, // zdotaccent
	'\u017D': 0.611, // Zcaron
	'\u017E': 0.500, // zcaron
	'\u0192': 0.556, // florin
	'\u0218': 0.667, // Scommaaccent
	'\u0219': 0.556, // scommaaccent
	'\u02C6': 0.333, // circumflex
	'\u02C7': 0.333, // caron
	'\u02D8': 0.333, // breve
	'\u02D9': 0.333, // dotaccent
	'\u02DA': 0.333, // ring
	'\u02DB': 0.333, // ogonek
	'\u02DC': 0.333, // tilde
	'\u02DD': 0.333, // hungarumlaut
	'\u2013': 0.556, // endash
	'\u2014': 1.000, // emdash
	'\u2018': 0.278, // quoteleft
	'\u2019': 0.278, // quoteright
	'\u201A': 0.278, // quotesinglbase
	'\u201C': 0.500, // quotedblleft
	'\u201D': 0.500, // quotedblright
	'\u201E': 0.500, // quotedblbase
	'\u2020': 0.556, // dagger
	'\u2021': 0.556, // daggerdbl
	'\u2022': 0.350, // bullet
	'\u2026': 1.000, // ellipsis
	'\u2030': 1.000, // perthousand
	'\u2039': 0.333, // guilsinglleft
	'\u203A': 0.333, // guilsinglright
	'\u2044': 0.167, // fraction
	'\u20AC': 0.556, // Euro
	'\u2122': 1.000, // trademark
	'\u2202': 0.494, // partialdiff
	'\u2206': 0.612, // Delta
	'\u2211': 0.600, // summation
	'\u2212': 0.584, // minus
	'\u221A': 0.549, // radical
	'\u2260': 0.549, // notequal
	'\u2264': 0.549, // lessequal
	'\u2265': 0.549, // greaterequal
	'\u25CA': 0.494, // lozenge
	'\uF6C3': 0.250, // commaaccent
	'\uFB01': 0.611, // fi
	'\uFB02': 0.611, // fl
}
