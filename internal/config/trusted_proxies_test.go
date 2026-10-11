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

package config

import (
	"slices"
	"testing"
)

func TestTrustedProxiesParsing(t *testing.T) {
	t.Setenv("MIABI_TRUSTED_PROXIES", " 10.62.0.0/16, ,192.0.2.7,")
	if got, want := New().TrustedProxies, []string{"10.62.0.0/16", "192.0.2.7"}; !slices.Equal(got, want) {
		t.Errorf("TrustedProxies = %v, want %v", got, want)
	}
	t.Setenv("MIABI_TRUSTED_PROXIES", "")
	if got := New().TrustedProxies; len(got) != 0 {
		t.Errorf("unset TrustedProxies = %v, want empty", got)
	}
}
