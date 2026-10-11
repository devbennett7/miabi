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

func TestSplitImageRef(t *testing.T) {
	cases := []struct {
		ref       string
		wantImage string
		wantTag   string
	}{
		{"nginx:latest", "nginx", "latest"},
		{"nginx", "nginx", ""},
		{"ghcr.io/org/app:v2", "ghcr.io/org/app", "v2"},
		{"registry:5000/app:1.2.3", "registry:5000/app", "1.2.3"}, // port in host, tag at end
		{"registry:5000/app", "registry:5000/app", ""},            // port in host, no tag
		{"repo@sha256:abc", "repo@sha256:abc", ""},                // digest pin
		{"", "", ""},
	}
	for _, c := range cases {
		t.Run(c.ref, func(t *testing.T) {
			img, tag := SplitImageRef(c.ref)
			if img != c.wantImage || tag != c.wantTag {
				t.Errorf("SplitImageRef(%q) = (%q,%q), want (%q,%q)", c.ref, img, tag, c.wantImage, c.wantTag)
			}
		})
	}
}

// Round-trip: composing then splitting yields the original repo + a tag.
func TestSplitComposeRoundTrip(t *testing.T) {
	img, tag := SplitImageRef(ComposeImageRef("nginx", "v2"))
	if img != "nginx" || tag != "v2" {
		t.Errorf("round-trip = (%q,%q), want (nginx,v2)", img, tag)
	}
	// No tag composes to :latest and splits back to latest.
	img, tag = SplitImageRef(ComposeImageRef("nginx", ""))
	if img != "nginx" || tag != "latest" {
		t.Errorf("round-trip default = (%q,%q), want (nginx,latest)", img, tag)
	}
}
