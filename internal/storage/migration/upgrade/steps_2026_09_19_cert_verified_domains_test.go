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

package upgrade

import "testing"

func TestFirstUncoveredName(t *testing.T) {
	verified := []certDomainRow{{Name: "example.com", Verified: true}}
	unverified := []certDomainRow{{Name: "example.com"}}
	banned := []certDomainRow{{Name: "example.com", Verified: true, Banned: true}}

	cases := []struct {
		name    string
		cert    certRow
		domains []certDomainRow
		want    string
	}{
		{"covered by a verified domain", certRow{CommonName: "app.example.com"}, verified, ""},
		{"the domain itself", certRow{CommonName: "example.com"}, verified, ""},
		{"wildcard", certRow{CommonName: "*.example.com"}, verified, ""},
		// The hijack: registered but never proven.
		{"unverified domain", certRow{CommonName: "app.example.com"}, unverified, "app.example.com"},
		{"banned domain", certRow{CommonName: "app.example.com"}, banned, "app.example.com"},
		{"no domains at all", certRow{CommonName: "app.example.com"}, nil, "app.example.com"},
		// A single bad SAN condemns the certificate — it is the name that would be served.
		{"bad SAN", certRow{CommonName: "app.example.com", DNSNames: `["app.example.com","console.other.com"]`}, verified, "console.other.com"},
		{"good SANs", certRow{CommonName: "app.example.com", DNSNames: `["app.example.com","api.example.com"]`}, verified, ""},
		// Unparseable SANs must not silently pass the common name alone.
		{"broken SAN json", certRow{CommonName: "app.example.com", DNSNames: `[oops`}, verified, ""},
		{"empty cert", certRow{}, verified, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := firstUncoveredName(c.cert, c.domains); got != c.want {
				t.Fatalf("firstUncoveredName = %q, want %q", got, c.want)
			}
		})
	}
}
