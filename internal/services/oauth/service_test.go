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

func TestDomainAllowed(t *testing.T) {
	cases := []struct {
		name    string
		email   string
		allowed string
		want    bool
	}{
		{"empty allowlist permits any", "a@example.com", "", true},
		{"exact match", "a@example.com", "example.com", true},
		{"case insensitive", "A@Example.COM", "example.com", true},
		{"one of several", "a@acme.org", "example.com, acme.org", true},
		{"not in list", "a@evil.com", "example.com,acme.org", false},
		{"no at sign", "invalid", "example.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := domainAllowed(tc.email, tc.allowed); got != tc.want {
				t.Errorf("domainAllowed(%q, %q) = %v, want %v", tc.email, tc.allowed, got, tc.want)
			}
		})
	}
}

func TestScopesOrDefault(t *testing.T) {
	if got := scopesOrDefault(""); got != "openid email profile" {
		t.Errorf("expected default scopes, got %q", got)
	}
	if got := scopesOrDefault("openid"); got != "openid" {
		t.Errorf("expected passthrough, got %q", got)
	}
}
