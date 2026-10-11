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

package node

import "testing"

func TestPublicIP(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"public v4", "203.0.113.7", "203.0.113.7"},
		{"public v4 with port", "203.0.113.7:54122", "203.0.113.7"},
		{"public v6", "2606:4700:4700::1111", "2606:4700:4700::1111"},
		{"private 10/8", "10.0.0.4", ""},
		{"private 192.168", "192.168.1.20", ""},
		{"cgnat-ish private 172.16", "172.16.5.5", ""},
		{"loopback", "127.0.0.1", ""},
		{"link-local", "169.254.10.1", ""},
		{"unspecified", "0.0.0.0", ""},
		{"empty", "", ""},
		{"garbage", "not-an-ip", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := PublicIP(tc.in); got != tc.want {
				t.Errorf("PublicIP(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
