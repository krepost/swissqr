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
	"reflect"
	"testing"
)

func TestReflowAtRune(t *testing.T) {
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit,",
		"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}
	width := 5.0 * 28.35 / 10.0 // 5cm × 28.35 pt/cm ÷ 10pt font size.
	expected := []string{
		"Lorem ipsum dolor sit amet, co",
		"nsectetur adipiscing elit,",
		"sed do eiusmod tempor incididu",
		"nt ut labore et dolore magna ali",
		"qua.",
	}
	actual := reflowAtRune(lines, width)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestReflowAtSpace(t *testing.T) {
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit,",
		"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}
	width := 5.0 * 28.35 / 10.0 // 5cm × 28.35 pt/cm ÷ 10pt font size.
	expected := []string{
		"Lorem ipsum dolor sit amet,",
		"consectetur adipiscing elit,",
		"sed do eiusmod tempor",
		"incididunt ut labore et dolore",
		"magna aliqua.",
	}
	actual := reflowAtSpace(lines, width)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestReflowAtSpaceWithLongWord(t *testing.T) {
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit,",
		"sed do eiusmodtemporincididuntutlaboreetdoloremagna aliqua.",
	}
	width := 5.0 * 28.35 / 10.0 // 5cm × 28.35 pt/cm ÷ 10pt font size.
	expected := []string{
		"Lorem ipsum dolor sit amet,",
		"consectetur adipiscing elit,",
		"sed do eiusmodtemporincididun",
		"tutlaboreetdoloremagna aliqua.",
	}
	actual := reflowAtSpace(lines, width)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestReflowAtSpaceWithOnlyLongWord(t *testing.T) {
	lines := []string{
		"Eiusmodtemporincididuntutlaboreetdoloremagna.",
	}
	width := 5.0 * 28.35 / 10.0 // 5cm × 28.35 pt/cm ÷ 10pt font size.
	expected := []string{
		"Eiusmodtemporincididuntutlabo",
		"reetdoloremagna.",
	}
	actual := reflowAtSpace(lines, width)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}
