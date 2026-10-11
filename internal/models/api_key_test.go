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

package models

import "testing"

func TestIsRegistryOnly(t *testing.T) {
	cases := []struct {
		scopes []string
		want   bool
	}{
		{nil, false}, // empty defaults to read (general)
		{[]string{ScopeRead}, false},
		{[]string{ScopeWrite}, false},
		{[]string{ScopeAll}, false},
		{[]string{ScopeRegistryRead}, true},
		{[]string{ScopeRegistryWrite}, true},
		{[]string{ScopeRegistryRead, ScopeRegistryWrite}, true},
		{[]string{ScopeRegistryWrite, ScopeRead}, false}, // mixed → general access too
	}
	for _, tc := range cases {
		k := &APIKey{Scopes: tc.scopes}
		if got := k.IsRegistryOnly(); got != tc.want {
			t.Errorf("IsRegistryOnly(%v) = %v, want %v", tc.scopes, got, tc.want)
		}
	}
}
