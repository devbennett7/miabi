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

package oauth

import "testing"

func TestFirstClaim(t *testing.T) {
	claims := map[string]any{
		"email":      "  std@example.com  ",
		"work_email": "mapped@example.com",
		"num":        42, // non-string ignored
	}
	cases := []struct {
		name  string
		names []string
		want  string
	}{
		{"mapped claim wins", []string{"work_email", "email"}, "mapped@example.com"},
		{"falls back to standard", []string{"", "email"}, "std@example.com"},
		{"trims whitespace", []string{"email"}, "std@example.com"},
		{"missing → empty", []string{"nope"}, ""},
		{"non-string ignored", []string{"num", "email"}, "std@example.com"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := firstClaim(claims, tc.names...); got != tc.want {
				t.Fatalf("firstClaim(%v) = %q, want %q", tc.names, got, tc.want)
			}
		})
	}
}
