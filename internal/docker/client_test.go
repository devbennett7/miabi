// Copyright (c) 2026 Jonas Kaninda.
//
// This file is part of Miabi.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package docker

import (
	"errors"
	"testing"
)

func TestLineWriterSplitsLines(t *testing.T) {
	var got []string
	w := &lineWriter{stream: "stdout", sink: func(l LogLine) error {
		got = append(got, l.Text)
		return nil
	}}
	// Write in chunks that split a line across calls.
	_, _ = w.Write([]byte("hello\nwor"))
	_, _ = w.Write([]byte("ld\nthird\n"))

	want := []string{"hello", "world", "third"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines %v, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLineWriterStopsOnSinkError(t *testing.T) {
	w := &lineWriter{stream: "stdout", sink: func(l LogLine) error {
		return errors.New("client gone")
	}}
	_, err := w.Write([]byte("a\n"))
	if !errors.Is(err, errSinkStop) {
		t.Fatalf("expected errSinkStop, got %v", err)
	}
}

func TestSanitizeTail(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", "100"},
		{"all", "5000"},
		{"garbage", "100"},
		{"0", "100"},
		{"-5", "100"},
		{"50", "50"},
		{"5000", "5000"},
		{"999999", "5000"},
	}
	for _, c := range cases {
		if got := sanitizeTail(c.in); got != c.want {
			t.Errorf("sanitizeTail(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
