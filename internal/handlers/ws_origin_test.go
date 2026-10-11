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

package handlers

import (
	"net/http"
	"testing"
)

func req(origin, host string) *http.Request {
	r := &http.Request{Host: host, Header: http.Header{}}
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	return r
}

func TestAllowWSOrigin(t *testing.T) {
	allow := allowWSOrigin([]string{"https://panel.example.com"})
	cases := []struct {
		name   string
		origin string
		host   string
		want   bool
	}{
		{"no origin (go client)", "", "panel.example.com", true},
		{"same-origin", "https://panel.example.com", "panel.example.com", true},
		{"allowlisted origin", "https://panel.example.com", "api.internal", true},
		{"allowlisted trailing slash", "https://panel.example.com/", "api.internal", true},
		{"cross-site origin rejected", "https://evil.example.com", "panel.example.com", false},
		{"malformed origin rejected", "not-a-url", "panel.example.com", false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := allow(req(tt.origin, tt.host)); got != tt.want {
				t.Errorf("allow(origin=%q host=%q) = %v, want %v", tt.origin, tt.host, got, tt.want)
			}
		})
	}
}

func TestAllowWSOriginWildcardDisablesCheck(t *testing.T) {
	allow := allowWSOrigin([]string{"*"})
	if !allow(req("https://evil.example.com", "panel.example.com")) {
		t.Error("wildcard allowlist should permit any origin (dev)")
	}
}
