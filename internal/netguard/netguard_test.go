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

package netguard

import "testing"

func TestValidateURL(t *testing.T) {
	Configure(false) // block private ranges
	blocked := []string{
		"http://127.0.0.1/hook",
		"https://169.254.169.254/latest/meta-data",
		"http://10.0.0.5:8080/x",
		"http://192.168.1.10/x",
		"http://[::1]/x",
		"ftp://example.com/x",
		"http://",
	}
	for _, u := range blocked {
		if err := ValidateURL(u); err == nil {
			t.Errorf("ValidateURL(%q) = nil, want error", u)
		}
	}
	allowed := []string{
		"https://example.com/hook",
		"http://1.2.3.4:9000/x",
	}
	for _, u := range allowed {
		if err := ValidateURL(u); err != nil {
			t.Errorf("ValidateURL(%q) = %v, want nil", u, err)
		}
	}
}

func TestAllowPrivateToggle(t *testing.T) {
	Configure(true) // allow private (homelab)
	defer Configure(false)
	if err := ValidateURL("http://192.168.1.10/x"); err != nil {
		t.Errorf("private target should be allowed when configured: %v", err)
	}
	// Loopback / link-local remain blocked even when private is allowed.
	if err := ValidateURL("http://127.0.0.1/x"); err == nil {
		t.Error("loopback must stay blocked")
	}
	if err := ValidateURL("http://169.254.169.254/x"); err == nil {
		t.Error("link-local metadata must stay blocked")
	}
}
